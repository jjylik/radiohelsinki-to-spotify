package spotifyclient

import (
	"context"
	"jjylik/radiohelsinki-to-spotify/playlist"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const sessionCookieName = "sessionid"

// FIXME: Hack, leaks memory
var clients sync.Map

var auth *spotifyauth.Authenticator
var spotifyAuthState = uuid.New().String()

func init() {
	redirect_url, ok := os.LookupEnv("REDIRECT_URL")
	if !ok {
		panic("REDIRECT_URL not set")
	}
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirect_url), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistModifyPrivate, spotifyauth.ScopePlaylistReadPrivate, spotifyauth.ScopePlaylistReadCollaborative))
}

func RegisterAuthHandlers() {
	http.HandleFunc("/auth", (func(w http.ResponseWriter, r *http.Request) {
		url := auth.AuthURL(spotifyAuthState)
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}))
	http.HandleFunc("/callback", completeAuth)
}

func GetClient(r *http.Request) *spotify.Client {
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		if client, ok := clients.Load(cookie.Value); ok {
			return client.(*spotify.Client)
		}
	}
	return nil
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), spotifyAuthState, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
		return
	}
	if st := r.FormValue("state"); st != spotifyAuthState {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, spotifyAuthState)
		return
	}
	sessionToken := uuid.New().String()
	token := &http.Cookie{
		Name:   sessionCookieName,
		Value:  sessionToken,
		MaxAge: 3600, // Magic magic... https://developer.spotify.com/documentation/ios/guides/token-swap-and-refresh/
	}
	http.SetCookie(w, token)
	client := spotify.New(auth.Client(r.Context(), tok))
	clients.Store(sessionToken, client)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func GetPlaylistNames(client *spotify.Client, userID string) ([]string, error) {
	spotifyPlaylists, err := client.GetPlaylistsForUser(context.Background(), userID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var playlistNames []string
	for _, spotifyPlaylist := range spotifyPlaylists.Playlists {
		playlistNames = append(playlistNames, spotifyPlaylist.Name)
	}
	return playlistNames, nil
}

func CreatePlaylists(client *spotify.Client, userID string, newPlaylists []playlist.Playlist) error {
	for _, playlist := range newPlaylists {
		newPlaylist, err := client.CreatePlaylistForUser(context.Background(), userID, playlist.Name, playlist.Name, false, false)
		if err != nil {
			log.Fatal(err)
			return err
		}
		trackIDs := []spotify.ID{}
		for _, track := range playlist.Tracks {
			query := track.Name + " artist:" + track.Artist
			searchResult, err := client.Search(context.Background(), query, spotify.SearchTypeTrack)
			if err != nil {
				log.Fatal(err)
				return err
			}
			foundTrack := searchResult.Tracks.Tracks
			if len(foundTrack) == 0 {
				log.Println("Could not find track:", track.Name, track.Artist)
			} else {
				trackIDs = append(trackIDs, foundTrack[0].ID)
			}
		}
		_, err = client.AddTracksToPlaylist(context.Background(), newPlaylist.ID, trackIDs...)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

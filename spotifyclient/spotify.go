package spotifyclient

import (
	"context"
	"jjylik/radiohelsinki-to-spotify/playlist"
	"log"
	"net/http"
	"os"
	"strings"
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
		log.Println(err)
		return
	}
	if st := r.FormValue("state"); st != spotifyAuthState {
		http.NotFound(w, r)
		log.Printf("State mismatch: %s != %s\n", st, spotifyAuthState)
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
		log.Println(err)
		return nil, err
	}
	var playlistNames []string
	for _, spotifyPlaylist := range spotifyPlaylists.Playlists {
		playlistNames = append(playlistNames, spotifyPlaylist.Name)
	}
	return playlistNames, nil
}

func AppendToPlaylist(client *spotify.Client, userID string, playlist playlist.Playlist) error {
	spotifyPlaylists, err := client.GetPlaylistsForUser(context.Background(), userID)
	if err != nil {
		log.Println(err)
		return err
	}
	var playlistID spotify.ID
	for _, spotifyPlaylist := range spotifyPlaylists.Playlists {
		if strings.EqualFold(spotifyPlaylist.Name, playlist.Name) {
			playlistID = spotifyPlaylist.ID
			break
		}
	}
	if playlistID == "" {
		newPlaylist, err := client.CreatePlaylistForUser(context.Background(), userID, playlist.Name, playlist.Name, false, false)
		if err != nil {
			log.Println(err)
			return err
		}
		playlistID = newPlaylist.ID
	}
	return addTracksToPlaylist(playlist, playlistID, client)
}

func CreatePlaylists(client *spotify.Client, userID string, newPlaylists []playlist.Playlist) error {
	for _, playlist := range newPlaylists {
		newPlaylist, err := client.CreatePlaylistForUser(context.Background(), userID, playlist.Name, playlist.Name, false, false)
		if err != nil {
			log.Println(err)
			return err
		}
		err = addTracksToPlaylist(playlist, newPlaylist.ID, client)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func addTracksToPlaylist(playlist playlist.Playlist, spotifyID spotify.ID, client *spotify.Client) error {
	trackIDs := []spotify.ID{}
	for _, track := range playlist.Tracks {
		query := track.Name + " artist:" + track.Artist
		searchResult, err := client.Search(context.Background(), query, spotify.SearchTypeTrack)
		if err != nil {
			log.Println(err)
			return err
		}
		foundTrack := searchResult.Tracks.Tracks
		if len(foundTrack) == 0 {
			log.Println("Could not find track:", track.Name, track.Artist)
		} else {
			trackIDs = append(trackIDs, foundTrack[0].ID)
		}
	}
	_, err := client.AddTracksToPlaylist(context.Background(), spotifyID, trackIDs...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

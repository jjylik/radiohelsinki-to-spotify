package main

import (
	"context"
	"jjylik/radiohelsinki-to-spotify/parser"
	"jjylik/radiohelsinki-to-spotify/playlist"
	"jjylik/radiohelsinki-to-spotify/spotifyclient"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	spotifyclient.RegisterAuthHandlers()
	http.HandleFunc("/playlists", (func(w http.ResponseWriter, r *http.Request) {
		if client := getSpotifyClientOrWriteError(w, r); client != nil {
			err := createMissingPlaylists(client)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}))
	http.HandleFunc("/custom-playlist", (func(w http.ResponseWriter, r *http.Request) {
		if client := getSpotifyClientOrWriteError(w, r); client != nil {
			playlist, err := parser.ParseFromWebForm(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			user, err := client.CurrentUser(context.Background())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			err = spotifyclient.AppendToPlaylist(client, user.ID, *playlist)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("Done!"))
		}
	}))

	log.Println("Server running")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getSpotifyClientOrWriteError(w http.ResponseWriter, r *http.Request) *spotify.Client {
	client := spotifyclient.GetClient(r)
	if client == nil {
		http.Error(w, "Invalid session, please login!", http.StatusBadRequest)
		return nil
	}
	return client
}

func createMissingPlaylists(client *spotify.Client) error {
	playlistFromRadioHelsinki, err := parser.ParsePlaylistFromRadioHelsinki()
	if err != nil {
		return err
	}
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		return err
	}
	existingPlaylists, err := spotifyclient.GetPlaylistNames(client, user.ID)
	if err != nil {
		return err
	}
	newPlaylists := playlist.GetNewPlaylists(playlistFromRadioHelsinki, existingPlaylists)
	spotifyclient.CreatePlaylists(client, user.ID, newPlaylists)
	log.Println("Done!")
	return nil
}

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
		if client := spotifyclient.GetClient(r); client != nil {
			err := createMissingPlaylists(client)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Invalid session, please login!", http.StatusBadRequest)
		}
	}))
	log.Println("Server running")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
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

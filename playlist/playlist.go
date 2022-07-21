package playlist

import (
	"log"
)

type Playlist struct {
	Name   string
	Tracks []Track
}
type Track struct {
	Name   string
	Artist string
}

func GetNewPlaylists(playlists *[]Playlist, existingPlaylistNames []string) []Playlist {
	var results []Playlist
	for _, playlist := range *playlists {
		exists := false
		for _, existingPlaylistName := range existingPlaylistNames {
			if playlist.Name == existingPlaylistName {
				exists = true
				break
			}
		}
		if !exists {
			log.Println("Found new playlist", playlist.Name)
			results = append(results, playlist)
		}
	}
	return results
}

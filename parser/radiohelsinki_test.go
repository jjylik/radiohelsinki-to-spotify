package parser

import (
	"io"
	"os"
	"path"
	"testing"
)

func getInputHtmlReader() io.Reader {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	reader, err := os.Open(path.Join(pwd, "test-fixtures/radiohelsinki.html"))
	if err != nil {
		panic(err)
	}
	return reader
}

func TestGetPlayListCount(t *testing.T) {
	playlists, _ := ParsePlaylists(getInputHtmlReader())
	if len(*playlists) != 10 {
		t.Errorf("expected 10 playlists, got %d", len(*playlists))
	}
}

func TestGetPlaylistSongName(t *testing.T) {
	playlists, _ := ParsePlaylists(getInputHtmlReader())
	firstSong := (*playlists)[0].Tracks[0].Name
	if firstSong != "C'est Cette Commette" {
		t.Errorf("expected C'est Cette Commette, got %s", firstSong)
	}
}

func TestGetPlaylistSongArtist(t *testing.T) {
	playlists, _ := ParsePlaylists(getInputHtmlReader())
	firstSong := (*playlists)[1].Tracks[0].Artist
	if firstSong != "ANDREW BIRD" {
		t.Errorf("expected ANDREW BIRD, got %s", firstSong)
	}
}

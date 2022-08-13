package parser

import (
	"errors"
	"jjylik/radiohelsinki-to-spotify/playlist"
	"net/http"
	"strings"
)

func ParseFromWebForm(r *http.Request) (*playlist.Playlist, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	var tracks []playlist.Track
	playlistTitle := r.Form.Get("name")
	rawSongs := r.Form.Get("tracks")
	artistTrackPairs := strings.Split(rawSongs, "\n")
	if len(artistTrackPairs) > 30 {
		return nil, errors.New("too many tracks")
	}
	for _, rawPairs := range artistTrackPairs {
		pair := strings.Split(strings.TrimSpace(rawPairs), ":")
		if len(pair) != 2 {
			return nil, errors.New("invalid input")
		}
		artist := pair[0]
		track := pair[1]
		tracks = append(tracks, playlist.Track{Artist: artist, Name: track})
	}
	if len(strings.TrimSpace(playlistTitle)) == 0 {
		playlistTitle = "Radio Helsinki"
	}
	p := playlist.Playlist{Name: playlistTitle, Tracks: tracks}
	return &p, nil
}

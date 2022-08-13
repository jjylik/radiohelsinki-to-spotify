package parser

import (
	"io"
	"jjylik/radiohelsinki-to-spotify/playlist"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParsePlaylistFromRadioHelsinki() (*[]playlist.Playlist, error) {
	res, err := http.Get("http://www.radiohelsinki.fi/ohjelma/henri-pulkkinen/")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return ParsePlaylists(res.Body)
}

func ParsePlaylists(document io.Reader) (*[]playlist.Playlist, error) {
	doc, err := goquery.NewDocumentFromReader(document)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var playlists []playlist.Playlist
	doc.Find(".playlist").Each(func(i int, pl *goquery.Selection) {
		name := pl.Find(".date").Text()
		var tracks []playlist.Track
		pl.Find(".songs").Find("li").Each(func(j int, songNode *goquery.Selection) {
			title := songNode.Find(".song").Text()
			artist := songNode.Find(".artist").Text()
			tracks = append(tracks, playlist.Track{Name: strings.TrimSpace(title), Artist: strings.TrimSpace(artist)})
		})
		playlists = append(playlists, playlist.Playlist{Name: "RH " + strings.TrimSpace(name), Tracks: tracks})
	})
	return &playlists, nil
}

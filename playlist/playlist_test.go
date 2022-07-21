package playlist

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetNewPlaylists(t *testing.T) {
	var tests = []struct {
		all      []Playlist
		existing []string
		want     []Playlist
	}{
		{[]Playlist{{Name: "1"}}, []string{"1"}, []Playlist{}},
		{[]Playlist{}, []string{"1"}, []Playlist{}},
		{[]Playlist{{Name: "1"}}, []string{"2"}, []Playlist{{Name: "1"}}},
		{[]Playlist{{Name: "1"}, {Name: "2"}}, []string{"2", "3"}, []Playlist{{Name: "1"}}},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v,%v", tt.all, tt.existing)
		t.Run(testname, func(t *testing.T) {
			ans := GetNewPlaylists(&tt.all, tt.existing)
			if (len(ans) != 0 && len(tt.want) != 0) && !cmp.Equal(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}

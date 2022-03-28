package netease

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/chyroc/dl/pkgs/parse/adapter/common"
)

func Parse(url string) (req common.MusicRequest, err error) {
	idType, id, err := parseID(url)
	if err != nil {
		return
	}

	switch idType {
	case "album":
		req = NewAlbumRequest(id)
	case "artist":
		req = NewArtistRequest(id)
	case "program":
		req = NewDJRequest(id)
	case "djradio":
		req = NewDjradioRequest(id)
	case "playlist", "discover/toplist":
		req = NewPlaylistRequest(id)
	case "song":
		req = NewSongRequest(id)
	}

	return
}

var rulPattern = regexp.MustCompile(`(?m)\/(song|artist|album|djradio|program|playlist|discover\/toplist)\?id=(\d+)`)

func parseID(url string) (string, int, error) {
	matched := rulPattern.FindStringSubmatch(url)
	if len(matched) < 3 {
		return "", 0, fmt.Errorf("could not parse the url: %s", url)
	}

	id, err := strconv.Atoi(matched[2])
	if err != nil {
		return "", 0, fmt.Errorf("could not parse the url: %s", url)
	}
	return matched[1], id, nil
}

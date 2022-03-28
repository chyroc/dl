package netease

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/chyroc/dl/internal/parse/adapter/common"
)

var UrlPattern = regexp.MustCompile("/(song|artist|album|playlist)\\?id=(\\d+)")

func Parse(url string) (req common.MusicRequest, err error) {
	matched, ok := UrlPattern.FindStringSubmatch(url), UrlPattern.MatchString(url)
	if !ok || len(matched) < 3 {
		err = fmt.Errorf("could not parse the url: %s", url)
		return
	}

	id, err := strconv.Atoi(matched[2])
	if err != nil {
		return
	}

	switch matched[1] {
	case "song":
		req = NewSongRequest(id)
	case "artist":
		req = NewArtistRequest(id)
	case "album":
		req = NewAlbumRequest(id)
	case "playlist":
		req = NewPlaylistRequest(id)
	}

	return
}

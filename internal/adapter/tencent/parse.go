package tencent

import (
	"fmt"
	"regexp"

	"github.com/chyroc/dl/internal/adapter/common"
)

const (
	UrlPattern = "/(song|songDetail|singer|album|playsquare)/(\\w+)\\.html"
)

func Parse(url string) (req common.MusicRequest, err error) {
	re := regexp.MustCompile(UrlPattern)
	matched, ok := re.FindStringSubmatch(url), re.MatchString(url)
	if !ok || len(matched) < 3 {
		err = fmt.Errorf("could not parse the url: %s", url)
		return
	}

	switch matched[1] {
	case "song", "songDetail":
		req = NewSongRequest(matched[2])
	case "singer":
		req = NewSingerRequest(matched[2])
	case "album":
		req = NewAlbumRequest(matched[2])
	case "playsquare":
		req = NewPlaylistRequest(matched[2])
	}

	return
}

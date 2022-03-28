package netease

import (
	"fmt"
	"strings"
	"time"

	"github.com/chyroc/dl/internal/helper"
	"github.com/chyroc/dl/internal/resource"
)

type Artist struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Album struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PicUrl      string `json:"picUrl"`
	PublishTime int64  `json:"publishTime"`
}

type SongUrl struct {
	Id   int    `json:"id"`
	Code int    `json:"code"`
	Url  string `json:"url"`
}

type Song struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Artist      []Artist `json:"ar"`
	Album       Album    `json:"al"`
	Position    int      `json:"no"`
	PublishTime int64    `json:"publishTime"`
}

type TrackId struct {
	Id int `json:"id"`
}

type Playlist struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	TrackIds []TrackId `json:"trackIds"`
}

func (s *Song) Extract() *resource.MP3 {
	title, album := strings.TrimSpace(s.Name), strings.TrimSpace(s.Album.Name)
	publishTime := time.Unix(0, s.PublishTime*1000*1000)
	year, track := fmt.Sprintf("%d", publishTime.Year()), fmt.Sprintf("%d", s.Position)
	coverImage := s.Album.PicUrl

	artistList := make([]string, 0, len(s.Artist))
	for _, ar := range s.Artist {
		artistList = append(artistList, strings.TrimSpace(ar.Name))
	}
	artist := strings.Join(artistList, "/")

	fileName := helper.TrimInvalidFilePathChars(fmt.Sprintf("%s - %s.mp3", strings.Join(artistList, " "), title))
	tag := resource.Tag{
		Title:      title,
		Artist:     artist,
		Album:      album,
		Year:       year,
		Track:      track,
		CoverImage: coverImage,
	}

	return &resource.MP3{
		FileName: fileName,
		Tag:      tag,
	}
}

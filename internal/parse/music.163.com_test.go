package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/chyroc/dl/internal/resource"
	"github.com/chyroc/go-assert"
)

func Test_music163Com(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	// python main.py http://music.163.com/#/program?id=1369232209 # dj
	t.Run("单曲", func(t *testing.T) {
		res, err := parse.NewMusic163Com().Parse("https://music.163.com/#/song?id=1843572582")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("张靖中 双火 minor - 是以不去.mp3", res.Title())
	})

	t.Run("专辑", func(t *testing.T) {
		res, err := parse.NewMusic163Com().Parse("https://music.163.com/#/album?id=132874562")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("幸存者的负罪感", res.Title())
		as.Len(res.(resource.MP3ChapterResource).Chapters(), 11)
		as.Equal("王以太 艾热 AIR - 末 (Intro).mp3", res.(resource.MP3ChapterResource).Chapters()[0].Title())
	})

	t.Run("歌单", func(t *testing.T) {
		res, err := parse.NewMusic163Com().Parse("https://music.163.com/#/playlist?id=156934569")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("这些充满『强烈画面感』的音乐", res.Title())
		as.Len(res.(resource.MP3ChapterResource).Chapters(), 433)
		as.Equal("胡伟立 - 勇往直前.mp3", res.(resource.MP3ChapterResource).Chapters()[0].Title())
	})

	t.Run("榜单", func(t *testing.T) {
		res, err := parse.NewMusic163Com().Parse("https://music.163.com/#/discover/toplist?id=1978921795")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("云音乐电音榜", res.Title())
		as.True(len(res.(resource.MP3ChapterResource).Chapters()) > 10)
		as.Equal("TYSM - Normal No More.mp3", res.(resource.MP3ChapterResource).Chapters()[0].Title())
	})

	t.Run("艺术家", func(t *testing.T) {
		res, err := parse.NewMusic163Com().Parse("https://music.163.com/#/artist?id=905705")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("剪影姐", res.Title())
		as.True(len(res.(resource.MP3ChapterResource).Chapters()) > 10)
		as.Equal("剪影姐 - 怎样.mp3", res.(resource.MP3ChapterResource).Chapters()[0].Title())
	})

	t.Run("电台", func(t *testing.T) {
		res, err := parse.NewMusic163Com().Parse("https://music.163.com/#/djradio?id=970764541")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("剪影姐", res.Title())
		as.True(len(res.(resource.MP3ChapterResource).Chapters()) > 10)
		as.Equal("剪影姐 - 怎样.mp3", res.(resource.MP3ChapterResource).Chapters()[0].Title())
	})

	t.Run("DJ", func(t *testing.T) {
		res, err := parse.NewMusic163Com().Parse("https://music.163.com/#/program?id=2499619396")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("三亩的音乐电台 - 神奇的心理疗法.mp3", res.Title())
	})
}

package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/chyroc/go-assert"
)

func Test_music163Com(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	as.Run("", func(a *assert.Assertions) {
		res, err := parse.NewMusic163Com().Parse("https://music.163.com/#/song?id=1843572582")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("张靖中 双火 minor - 是以不去", res.Title())
	})

	as.Run("", func(a *assert.Assertions) {
		res, err := parse.NewMusic163Com().Parse("https://music.163.com/#/album?id=132874562")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("幸存者的负罪感", res.Title())
	})
	as.Run("", func(a *assert.Assertions) {
		res, err := parse.NewMusic163Com().Parse("https://music.163.com/#/playlist?id=156934569")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("这些充满『强烈画面感』的音乐", res.Title())
	})
}

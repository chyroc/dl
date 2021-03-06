package parse_test

import (
	"testing"

	"github.com/chyroc/dl/pkgs/parse"
	"github.com/chyroc/dl/pkgs/resource"
	"github.com/chyroc/go-assert"
)

func Test_wwwMissevanCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	as.Run("", func(as *assert.Assertions) {
		res, err := parse.NewWwwMissevanCom().Parse("https://www.missevan.com/sound/player?id=1303686")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("【超酥助眠】新一助眠电台", res.Title())
		as.Equal("[1274558] 【声控】男友音耳边数羊，伴你入眠.mp3", res.(resource.ChapterResource).Chapters()[0].Title())
	})
}

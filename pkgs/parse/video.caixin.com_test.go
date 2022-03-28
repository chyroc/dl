package parse_test

import (
	"testing"

	"github.com/chyroc/dl/pkgs/parse"
	"github.com/chyroc/go-assert"
)

func Test_VideoCaixinCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	as.Run("", func(as *assert.Assertions) {
		res, err := parse.NewVideoCaixinCom().Parse("https://video.caixin.com/2021-09-09/101770028.html")
		as.Nil(err)
		as.Equal("【云起·智来】中旅王斌：传统旅游业数字化转型，关键在“五化”.mp4", res.Title())
	})
}

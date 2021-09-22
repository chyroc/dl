package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/chyroc/go-assert"
)

func Test_yQqCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	{
		res, err := parse.NewYQqCom().Parse("https://y.qq.com/n/yqq/song/002Zkt5S2z8JZx.html")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("...", res.Title())
	}
}

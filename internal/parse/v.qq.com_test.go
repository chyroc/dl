package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/chyroc/go-assert"
)

func Test_vQqCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	{
		res, err := parse.NewVQqCom().Parse("https://v.qq.com/txp/iframe/player.html?vid=j0822mqey5h")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("黄田天后古庙.mp4", res.Title())
	}
}

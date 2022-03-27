package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/chyroc/go-assert"
)

func Test_www333tttCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	as.Run("", func(as *assert.Assertions) {
		res, err := parse.NewWww333tttCom().Parse("http://www.333ttt.com/up/yy6182865.html")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("可惜没如果-林俊杰.mp3", res.Title())
	})
}

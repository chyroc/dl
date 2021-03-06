package parse_test

import (
	"testing"

	"github.com/chyroc/dl/pkgs/parse"
	"github.com/chyroc/go-assert"
)

func Test_a36krCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	as.Run("", func(as *assert.Assertions) {
		res, err := parse.NewA36krCom().Parse("https://36kr.com/video/1673124114052352")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("猪价持续下跌，什么是猪周期？.mp4", res.Title())
	})
}

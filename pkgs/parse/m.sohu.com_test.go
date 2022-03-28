package parse_test

import (
	"testing"

	"github.com/chyroc/dl/pkgs/parse"
	"github.com/chyroc/go-assert"
)

func Test_mSohuCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	as.Run("", func(as *assert.Assertions) {
		res, err := parse.NewMSohuCom().Parse("https://m.sohu.com/a/490513509_120538293")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("好吃的月饼#私藏美食大分享_120538293.mp4", res.Title())
	})
}

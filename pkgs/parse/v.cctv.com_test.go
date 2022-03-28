package parse_test

import (
	"testing"

	"github.com/chyroc/dl/pkgs/parse"
	"github.com/chyroc/go-assert"
)

func Test_vCctvCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	as.Run("", func(as *assert.Assertions) {
		res, err := parse.NewVCctvCom().Parse("https://v.cctv.com/2021/09/17/VIDERZvtKr1arx2zGkZprwqR210917.shtml")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("王冰冰和撒贝宁可以说毫无默契了.mp4", res.Title())
	})
}

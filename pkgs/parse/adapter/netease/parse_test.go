package netease

import (
	"testing"

	"github.com/chyroc/go-assert"
)

func Test_parseID(t *testing.T) {
	as := assert.New(t)
	as.Nil(nil)

	t.Run("", func(t *testing.T) {
		a, b, err := parseID("https://music.163.com/#/discover/toplist?id=1978921795")
		as.Nil(err)
		as.Equal("discover/toplist", a)
		as.Equal(1978921795, b)
	})
}

package parse_test

import (
	"testing"

	"github.com/chyroc/dl/pkgs/parse"
	"github.com/chyroc/go-assert"
)

func Test_wwwHanjukankanCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	{
		res, err := parse.NewWwwHanjukankanCom().Parse("...")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("...", res.Title())
	}
}


package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/chyroc/go-assert"
)

func Test_wwwSztvComCn(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	{
		res, err := parse.NewWwwSztvComCn().Parse("...")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("...", res.Title())
	}
}

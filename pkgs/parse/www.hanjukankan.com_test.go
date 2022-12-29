package parse_test

import (
	"testing"

	"github.com/chyroc/go-assert"

	"github.com/chyroc/dl/pkgs/parse"
)

func Test_wwwHanjukankanCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	{
		res, err := parse.NewWwwHanjukankanCom().Parse("https://www.hanjukankan.com/movie/index159.html")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("请回答1988", res.Title())
	}

	t.Run("", func(t *testing.T) {
		as.Equal("01", parse.FormatIndex(1))
		as.Equal("10", parse.FormatIndex(10))
		as.Equal("101", parse.FormatIndex(101))
	})
}

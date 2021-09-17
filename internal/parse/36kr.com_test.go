package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/stretchr/testify/assert"
)

func Test_a36krCom(t *testing.T) {
	as := assert.New(t)

	{
		res, err := parse.NewA36krCom().Parse("https://36kr.com/video/1402287552562048")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("想换iPhone13，官网天猫京东拼多多买有啥不一样？", res.Title())
	}
}

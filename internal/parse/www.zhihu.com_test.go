package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/stretchr/testify/assert"
)

func Test_wwwZhihuCom(t *testing.T) {
	as := assert.New(t)

	{
		res, err := parse.NewWwwZhihuCom().Parse("https://www.zhihu.com/zvideo/1385301246845173760")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("【答主推荐】这是一位不好好创业就要回家「继承」太平洋的宝藏视频答主", res.Title())
	}
}

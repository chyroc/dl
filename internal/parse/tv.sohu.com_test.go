package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/stretchr/testify/assert"
)

func Test_tvSohuCom(t *testing.T) {
	as := assert.New(t)

	{
		res, err := parse.NewTvSohuCom().Parse("https://tv.sohu.com/v/MjAyMTA5MTYvbjYwMTA0NzczNC5zaHRtbA==.html")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("考古笔记挖出一座城，竟改变中国三千年前历史？_288294918", res.Title())
	}
}

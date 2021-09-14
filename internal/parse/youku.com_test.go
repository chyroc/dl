package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/stretchr/testify/assert"
)

func Test_youkuCom(t *testing.T) {
	as := assert.New(t)

	{
		res, err := parse.NewYoukuCom().Parse("https://v.youku.com/v_show/id_XNDU1MTg1NjM2OA==")
		as.Nil(err)
		as.Equal("天津话《乡村爱情12》王木生VS二奎，俩大不孝子让赵本山操碎了心", res.Title())
	}
	{
		res, err := parse.NewYoukuCom().Parse("https://v.youku.com/v_show/id_XNTgwNTgzNTYwNA==.html?spm=some")
		as.Nil(err)
		as.Equal("程序员那么可爱：姜逸城吃醋三大症状 ！", res.Title())
	}
}

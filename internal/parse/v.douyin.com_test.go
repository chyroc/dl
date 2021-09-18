package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/stretchr/testify/assert"
)

func Test_vDouyinCom(t *testing.T) {
	as := assert.New(t)

	{
		res, err := parse.NewVDouyinCom().Parse("https://v.douyin.com/dAAcx4R/")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("风鹰铠甲：这波我在大气层#铠甲勇士 @飞天小帅帅 _111561509837", res.Title())
	}
	{
		res, err := parse.NewVDouyinCom().Parse("https://www.iesdouyin.com/share/video/7006244137951333662/?region=CN")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("风鹰铠甲：这波我在大气层#铠甲勇士 @飞天小帅帅 _111561509837", res.Title())
	}
	{
		res, err := parse.NewVDouyinCom().Parse("https://www.douyin.com/video/7006244137951333662?previous_page")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("风鹰铠甲：这波我在大气层#铠甲勇士 @飞天小帅帅 _111561509837", res.Title())
	}
}

package parse_test

import (
	"testing"

	"github.com/chyroc/dl/pkgs/parse"
	"github.com/chyroc/dl/pkgs/resource"
	"github.com/chyroc/go-assert"
)

func Test_mpWeixinQqCom(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		res, err := parse.NewMpWeixinQqCom().Parse("https://mp.weixin.qq.com/s/1Li7fFUu49XQoo6-dNwnmQ")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("永恒的丰碑——纪念黄田战役77周年！", res.Title())
		as.Equal("《地名宝安》丰碑.mp4", res.(resource.ChapterResource).Chapters()[0].Title())
	})

	t.Run("", func(t *testing.T) {
		res, err := parse.NewMpWeixinQqCom().Parse("https://mp.weixin.qq.com/s/xDNx_qGrrUoHK_HXBiKI5w")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("千年历史，别样航城，五个故事带你读懂航城！", res.Title())
		as.Len(res.(resource.ChapterResource).Chapters(), 3)
		as.Equal("黄田天后古庙.mp4", res.(resource.ChapterResource).Chapters()[0].Title())
		as.Equal("黄田革命纪念碑.mp4", res.(resource.ChapterResource).Chapters()[1].Title())
		as.Equal("草围蛋家文化.mp4", res.(resource.ChapterResource).Chapters()[2].Title())
	})
}

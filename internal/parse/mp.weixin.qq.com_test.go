package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/chyroc/dl/internal/resource"
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
}

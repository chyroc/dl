package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/chyroc/dl/internal/resource"
	"github.com/chyroc/go-assert"
)

func Test_wwwSztvComCn(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	t.Run("", func(t *testing.T) {
		res, err := parse.NewWwwSztvComCn().Parse("https://www.sztv.com.cn/ysz/yszlm/mt/szylshpd/78504152.shtml")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("《寻“宝”百年——宝安区“三微”展播第23集‖留取丹心照汗青——走近黄田革命烈士纪念碑》", res.Title())
		as.Len(res.(resource.ChapterResource).Chapters(), 1)
		as.Equal("留取丹心照汗青——走近黄田革命烈士纪念碑。.mp4", res.(resource.ChapterResource).Chapters()[0].Title())
	})
}

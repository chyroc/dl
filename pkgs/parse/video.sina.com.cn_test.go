package parse_test

import (
	"testing"

	"github.com/chyroc/dl/pkgs/parse"
	"github.com/chyroc/go-assert"
)

func Test_VideoSinaComCn(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))
	as.Nil(nil)

	as.Run("", func(as *assert.Assertions) {
		uri := "http://video.sina.com.cn/p/ent/doc/2016-10-14/225965380865.html"
		res, err := parse.NewVideoSinaComCn().Parse(uri)
		as.Nil(err)
		as.Equal("视频：辣眼睛！《爸爸4》田亮扭动唱歌穿秋裤.mp4", res.Title())
	})

	as.Run("", func(as *assert.Assertions) {
		res, err := parse.NewVideoSinaComCn().Parse("http://video.sina.com.cn/p/ent/doc/2018-02-07/090568002248.html")
		as.Nil(err)
		as.Equal("视频：高圆圆包场支持赵又廷新片 看到吻戏笑出声.mp4", res.Title())
	})
}

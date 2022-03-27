package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/chyroc/go-assert"
)

func Test_haokanBaiduCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	as.Run("", func(as *assert.Assertions) {
		res, err := parse.NewHaokanBaiduCom().Parse("https://haokan.baidu.com/v?vid=6518154682148026126")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("迷你世界 灾难模拟器 上帝控制着灾难 不爱护环境的人 接受惩罚吧,游戏,沙盒游戏,好看视频.mp4", res.Title())
	})
	as.Run("", func(as *assert.Assertions) {
		res, err := parse.NewHaokanBaiduCom().Parse("https://haokan.baidu.com/v?vid=7249594116085322255")
		as.Nil(err)
		as.Equal("蜡笔小新：小新一家来采茶园玩，猜茶环节小新胜出，主办方亏大了,动漫,日本动漫,好看视频.mp4", res.Title())
	})
}

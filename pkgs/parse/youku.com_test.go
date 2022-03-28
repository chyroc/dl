package parse

import (
	"testing"

	"github.com/chyroc/go-assert"
)

func Test_YoukuCom(t *testing.T) {
	as := assert.New(t, assert.WithFailRerun(5))

	t.Run("", func(t *testing.T) {
		x := ""
		x, _ = (&vYoukuCom{}).getVideoID("https://v.youku.com/v_show/id_XNDU1MTg1NjM2OA")
		as.Equal("XNDU1MTg1NjM2OA", x)

		x, _ = (&vYoukuCom{}).getVideoID("https://v.youku.com/v_show/id_XNDU1MTg1NjM2OA==")
		as.Equal("XNDU1MTg1NjM2OA", x)

		x, _ = (&vYoukuCom{}).getVideoID("https://v.youku.com/v_show/id_XNDU1MTg1NjM2OA==.html")
		as.Equal("XNDU1MTg1NjM2OA", x)
	})

	as.Run("1", func(as *assert.Assertions) {
		res, err := NewVYoukuCom().Parse("https://v.youku.com/v_show/id_XNDU1MTg1NjM2OA==")
		as.Nil(err)
		as.Equal("天津话《乡村爱情12》王木生VS二奎，俩大不孝子让赵本山操碎了心.mp4", res.Title())
	})

	as.Run("2", func(as *assert.Assertions) {
		res, err := NewVYoukuCom().Parse("https://v.youku.com/v_show/id_XNTgwNTgzNTYwNA==.html?spm=some")
		as.Nil(err)
		as.Equal("程序员那么可爱：姜逸城吃醋三大症状 ！.mp4", res.Title())
	})
}

package parse_test

import (
	"testing"

	"github.com/chyroc/dl/internal/parse"
	"github.com/chyroc/go-assert"
)

func Test_open163Com(t *testing.T) {
	t.Skip()

	as := assert.New(t, assert.WithFailRerun(5))

	as.Run("", func(as *assert.Assertions) {
		res, err := parse.NewOpen163Com().Parse("https://open.163.com/newview/movie/free?pid=HFD3PMIPO")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("如何提高阅读效率，读更多本书？", res.Title())
	})

	as.Run("", func(as *assert.Assertions) {
		res, err := parse.NewOpen163Com().Parse("https://open.163.com/movie/2010/6/D/6/M6TCSIN1U_M6TCSTQD6.html")
		as.Nil(err)
		as.NotNil(res)
		as.Equal("麻省理工学院公开课：计算机科学及编程导论", res.Title())
		m := res.MultiDownload()
		as.Len(m, 24)
		as.Equal("[M6TCSTQD6] 1_课程目标，数据类型，运算，变量", m[0].Title())
		as.Equal("[M6TCT9E81] 13_动态规划,重叠的子问题,最优子结构", m[12].Title())
		as.Equal("[M6TCTH5OC] 24_计算机科学家都做什么", m[23].Title())
	})
}

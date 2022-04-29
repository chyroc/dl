package identify

import (
	"github.com/chyroc/dl/pkgs/parse"
)

func init() {
	register(parse.NewA36krCom())
	register(parse.NewHaokanBaiduCom())
	register(parse.NewM3u8())
	register(parse.NewMBggeeCom())
	register(parse.NewMSohuCom())
	register(parse.NewMobileWeiboCn())
	register(parse.NewMpWeixinQqCom())
	register(parse.NewMusic163Com())
	register(parse.NewNewsCctvCom())
	register(parse.NewOpen163Com())
	register(parse.NewTvSohuCom())
	register(parse.NewVCctvCom())
	register(parse.NewVDouyinCom())
	register(parse.NewVQqCom())
	register(parse.NewVYoukuCom())
	register(parse.NewVideoCaixinCom())
	register(parse.NewVideoSinaComCn())
	register(parse.NewWww333tttCom())
	register(parse.NewWwwMissevanCom())
	register(parse.NewWwwSztvComCn())
	register(parse.NewWwwZhihuCom())
	register(parse.NewYQqCom())
}
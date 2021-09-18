package parse

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/download"
)

func NewVDouyinCom() Parser {
	return &vDouyinCom{}
}

type vDouyinCom struct{}

func (r *vDouyinCom) Kind() string {
	return "v.douyin.com,www.iesdouyin.com,www.douyin.com"
}

func (r *vDouyinCom) Parse(uri string) (download.Downloader, error) {
	vid, err := r.getVideoID(uri)
	if err != nil {
		return nil, err
	}
	meta, err := r.getMeta(uri, vid)
	if err != nil {
		return nil, err
	}
	title := fmt.Sprintf("%s_%d", meta.ItemList[0].Desc, meta.ItemList[0].AuthorUserID)
	spec := []*download.Specification{{URL: meta.ItemList[0].Video.PlayAddr.URLList[0]}}
	return download.NewDownloadURL(title, title+".mp4", spec), nil
}

func (r *vDouyinCom) getMeta(originURL, vid string) (*vDouyinComMetaResp, error) {
	uri := fmt.Sprintf("https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/?item_ids=%s", vid)
	resp := new(vDouyinComMetaResp)
	err := config.ReqCli.New(http.MethodGet, uri).WithHeaders(prepareCommonHeader(originURL, nil)).Unmarshal(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusMsg != "" {
		return nil, fmt.Errorf(resp.StatusMsg)
	}
	return resp, nil
}

func (r *vDouyinCom) getVideoID(uri string) (string, error) {
	// v.douyin.com/adfasdfadf/
	// www.iesdouyin.com/share/video/1234
	// www.douyin.com/video/1234

	uriParsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	switch uriParsed.Host {
	case "v.douyin.com":
		location, err := config.ReqCli.New(http.MethodGet, uri).WithRedirect(false).ResponseHeaderByKey("Location")
		if err != nil {
			return "", err
		}
		locationParsed, err := url.Parse(location)
		if err != nil {
			return "", err
		}
		if locationParsed.Host != "www.iesdouyin.com" && locationParsed.Host != "www.douyin.com" {
			return "", fmt.Errorf("url %q unsupport", uri)
		}
		uriParsed = locationParsed
		fallthrough
	case "www.iesdouyin.com":
		return getMatchStringByRegs(uriParsed.Path, vDouyinComDouyinIdReg), nil
	case "www.douyin.com":
		return getMatchStringByRegs(uriParsed.Path, vDouyinComDouyinIdReg), nil
	default:
		return "", fmt.Errorf("url %q unsupport", uri)
	}
}

type vDouyinComMetaResp struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	ItemList   []struct {
		Desc    string `json:"desc"`
		GroupID int64  `json:"group_id"`
		AwemeID string `json:"aweme_id"`
		Video   struct {
			Height   int `json:"height"`
			Duration int `json:"duration"`
			PlayAddr struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"play_addr"`
			Cover struct {
				URLList []string `json:"url_list"`
				URI     string   `json:"uri"`
			} `json:"cover"`
			Width        int    `json:"width"`
			Ratio        string `json:"ratio"`
			HasWatermark bool   `json:"has_watermark"`
			Vid          string `json:"vid"`
		} `json:"video"`
		ShareURL string `json:"share_url"`
		Duration int    `json:"duration"`
		Music    struct {
			CoverMedium struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"cover_medium"`
			CoverThumb struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"cover_thumb"`
			PlayURL struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"play_url"`
			Duration int    `json:"duration"`
			Mid      string `json:"mid"`
			Title    string `json:"title"`
			Author   string `json:"author"`
			ID       int64  `json:"id"`
			CoverHd  struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"cover_hd"`
			CoverLarge struct {
				URLList []string `json:"url_list"`
				URI     string   `json:"uri"`
			} `json:"cover_large"`
		} `json:"music"`
		IsLiveReplay bool `json:"is_live_replay"`
		CreateTime   int  `json:"create_time"`
		ChaList      []struct {
			Desc    string `json:"desc"`
			Cid     string `json:"cid"`
			ChaName string `json:"cha_name"`
		} `json:"cha_list"`
		ShareInfo struct {
			ShareWeiboDesc string `json:"share_weibo_desc"`
			ShareDesc      string `json:"share_desc"`
			ShareTitle     string `json:"share_title"`
		} `json:"share_info"`
		AuthorUserID int64 `json:"author_user_id"`
		Author       struct {
			ShortID      string `json:"short_id"`
			Signature    string `json:"signature"`
			AvatarLarger struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"avatar_larger"`
			AvatarThumb struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"avatar_thumb"`
			UniqueID     string `json:"unique_id"`
			UID          string `json:"uid"`
			AvatarMedium struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"avatar_medium"`
			Nickname string `json:"nickname"`
		} `json:"author"`
	} `json:"item_list"`
}

var vDouyinComDouyinIdReg = []*regexp.Regexp{
	regexp.MustCompile(`video/(.*?)/$`),
	regexp.MustCompile(`video/(.*?)$`),
}

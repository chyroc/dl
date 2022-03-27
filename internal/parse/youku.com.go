package parse

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/resource"
)

func NewVYoukuCom() Parser {
	return &vYoukuCom{}
}

type vYoukuCom struct{}

func (r *vYoukuCom) Kind() string {
	return "v.youku.com"
}

func (r *vYoukuCom) Parse(uri string) (resource.Resource, error) {
	meta, err := r.getMeta(uri)
	if err != nil {
		return nil, err
	}

	title := meta.Data.Video.Title
	specs := []*resource.Specification{}
	for _, v := range meta.Data.Stream {
		specs = append(specs, &resource.Specification{
			Size:       v.Size,
			Definition: resource.MayConvertDefinition(v.StreamType),
			URL:        v.Segs[0].CdnURL,
		})
	}
	return resource.NewURLWithSpecification(title+".mp4", specs), nil
}

func (r *vYoukuCom) getVideoID(uri string) (string, error) {
	match := youkuComRegId.FindStringSubmatch(uri)
	if len(match) == 2 {
		return match[1], nil
	}
	return "", fmt.Errorf("parse %q, get video_id failed", uri)
}

func (r *vYoukuCom) getMeta(uri string) (*youkuComGetMetaResp, error) {
	videoID, err := r.getVideoID(uri)
	if err != nil {
		return nil, err
	}
	query := map[string]string{
		"vid":       videoID,
		"ccode":     "0532",
		"client_ip": "192.168.1.1",
		"utid":      getMmstatEtag(),
		"client_ts": strconv.FormatInt(time.Now().Unix(), 10),
	}
	header := map[string]string{
		"Host":            "ups.youku.com",
		"Referer":         uri,
		"User-Agent":      userAgent,
		"Accept-Language": "en-us,en;q=0.5",
	}
	resp := new(youkuComGetMetaResp)

	err = config.ReqCli.New(http.MethodGet, "https://ups.youku.com/ups/get.json").WithHeaders(header).WithQuerys(query).Unmarshal(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type youkuComGetMetaResp struct {
	Cost float64 `json:"cost"`
	Data struct {
		Stream []struct {
			AudioLang         string `json:"audio_lang"`
			DrmType           string `json:"drm_type"`
			Height            int    `json:"height"`
			Logo              string `json:"logo"`
			M3U8URL           string `json:"m3u8_url"`
			MediaType         string `json:"media_type"`
			MillisecondsAudio int    `json:"milliseconds_audio"`
			MillisecondsVideo int    `json:"milliseconds_video"`
			Segs              []struct {
				CdnURL                 string `json:"cdn_url"`
				Fileid                 string `json:"fileid"`
				Key                    string `json:"key"`
				Secret                 string `json:"secret"`
				Size                   int    `json:"size"`
				TotalMillisecondsAudio int    `json:"total_milliseconds_audio"`
				TotalMillisecondsVideo int    `json:"total_milliseconds_video"`
			} `json:"segs"`
			Size         int64  `json:"size"`
			StreamType   string `json:"stream_type"`
			SubtitleLang string `json:"subtitle_lang"`
			Width        int    `json:"width"`
		} `json:"stream"`
		Uploader struct {
			Avatar struct {
				Big    string `json:"big"`
				Large  string `json:"large"`
				Middle string `json:"middle"`
				Small  string `json:"small"`
				Xlarge string `json:"xlarge"`
			} `json:"avatar"`
			Homepage string `json:"homepage"`
			UID      string `json:"uid"`
			Username string `json:"username"`
		} `json:"uploader"`
		Ups struct {
			Psid           string `json:"psid"`
			Strf           bool   `json:"strf"`
			Stsp           bool   `json:"stsp"`
			UpsClientNetip string `json:"ups_client_netip"`
			UpsTs          string `json:"ups_ts"`
		} `json:"ups"`
		User struct {
			IP         string `json:"ip"`
			PartnerVip bool   `json:"partnerVip"`
			UID        string `json:"uid"`
			Ytid       string `json:"ytid"`
		} `json:"user"`
		Video struct {
			Encodeid    string  `json:"encodeid"`
			ID          int     `json:"id"`
			Limited     int     `json:"limited"`
			Logo        string  `json:"logo"`
			Seconds     float64 `json:"seconds"`
			Source      int     `json:"source"`
			StSorted    int     `json:"st_sorted"`
			StreamTypes struct {
				Default []string `json:"default"`
			} `json:"stream_types"`
			Subcategories []interface{} `json:"subcategories"`
			Tags          []string      `json:"tags"`
			Title         string        `json:"title"`
			Type          []string      `json:"type"`
			UID           int           `json:"uid"`
			Userid        int           `json:"userid"`
			Username      string        `json:"username"`
			VideoidPlay   int           `json:"videoid_play"`
			Weburl        string        `json:"weburl"`
		} `json:"video"`
	} `json:"data"`
}

func getMmstatEtag() string {
	resp, err := http.Get("https://log.mmstat.com/eg.js")
	if err != nil {
		return ""
	}
	etag := resp.Header.Get("etag")
	return etag[1 : len(etag)-1] // un quoto
}

var youkuComRegId = regexp.MustCompile(`v.youku.com/v_show/id_(.*?)=*`)

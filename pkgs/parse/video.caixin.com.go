package parse

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/chyroc/dl/pkgs/config"
	"github.com/chyroc/dl/pkgs/resource"
)

func NewVideoCaixinCom() Parser {
	return &videoCaixinCom{}
}

type videoCaixinCom struct{}

func (r *videoCaixinCom) Kind() string {
	return "video.caixin.com"
}

func (r *videoCaixinCom) ExampleURLs() []string {
	return []string{"https://video.caixin.com/2021-09-09/101770028.html"}
}

func (r *videoCaixinCom) Parse(uri string) (resource.Resourcer, error) {
	id, title, err := r.getVideoID(uri)
	if err != nil {
		return nil, err
	}
	token, err := r.getVideoToken(uri, id)
	if err != nil {
		return nil, err
	}
	meta, err := r.getVideoMeta(uri, token.GetPlayInfoToken)
	if err != nil {
		return nil, err
	}

	// 组装数据
	specs := []*resource.Specification{}
	for _, v := range meta.PlayInfoList {
		specs = append(specs, &resource.Specification{
			Size:       v.Size,
			Definition: resource.MayConvertDefinition(v.Definition),
			URL:        v.MainPlayURL,
		})
	}

	return resource.NewURLWithSpecification(title+".mp4", specs), nil
}

func (r *videoCaixinCom) getVideoID(uri string) (string, string, error) {
	header := prepareCommonHeader(uri, nil)
	text, err := config.ReqCli.New(http.MethodGet, uri).WithHeaders(header).Text()
	if err != nil {
		return "", "", err
	}
	id := getMatchString(text, videoCaixinComIdReg)
	topic := getMatchString(text, videoCaixinComTopicReg)

	return id, topic, nil
}

func (r *videoCaixinCom) getVideoToken(originURL string, videoID string) (*videoCaixinComGetVideoTokenBase64, error) {
	uri := "https://gateway.caixin.com/api/appservice/player/token"
	query := map[string]string{
		"vid": videoID,
	}
	header := prepareCommonHeader(originURL, nil)
	resp := new(videoCaixinComGetVideoTokenResp)

	err := config.ReqCli.New(http.MethodGet, uri).WithQuerys(query).WithHeaders(header).Unmarshal(resp)
	if err != nil {
		return nil, err
	} else if resp.Message != "" {
		return nil, fmt.Errorf(resp.Message)
	}
	bs, err := base64.StdEncoding.DecodeString(resp.Token)
	if err != nil {
		return nil, err
	}
	meta := new(videoCaixinComGetVideoTokenBase64)
	if err = json.Unmarshal(bs, meta); err != nil {
		return nil, err
	}

	return meta, nil
}

func (r *videoCaixinCom) getVideoMeta(originURL string, token string) (*videoCaixinComMeta, error) {
	uri := fmt.Sprintf("https://vod.bytedanceapi.com/?%s&ssl=true", token)
	header := prepareCommonHeader(originURL, nil)
	resp := new(videoCaixinComGetMetaResp)

	err := config.ReqCli.New(http.MethodGet, uri).WithHeaders(header).Unmarshal(resp)
	if err != nil {
		return nil, err
	} else if resp.ResponseMetadata.Error.Message != "" {
		return nil, fmt.Errorf(resp.ResponseMetadata.Error.Message)
	}

	return &resp.Result.Data, nil
}

type videoCaixinComGetVideoTokenResp struct {
	Message string `json:"Message"`
	Token   string `json:"token"`
}

type videoCaixinComGetVideoTokenBase64 struct {
	GetPlayInfoToken string `json:"GetPlayInfoToken"`
}

type videoCaixinComGetMetaResp struct {
	ResponseMetadata struct {
		RequestID string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Error     struct {
			CodeN   int    `json:"CodeN"`
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`
	} `json:"ResponseMetadata"`
	Result struct {
		EncryptKey string             `json:"EncryptKey"`
		CipherText string             `json:"CipherText"`
		Data       videoCaixinComMeta `json:"Data"`
	} `json:"Result"`
}

type videoCaixinComMeta struct {
	Status         int     `json:"Status"`
	VideoID        string  `json:"VideoID"`
	CoverURL       string  `json:"CoverUrl"`
	Duration       float64 `json:"Duration"`
	MediaType      string  `json:"MediaType"`
	EnableAdaptive bool    `json:"EnableAdaptive"`
	PlayInfoList   []struct {
		Bitrate          int     `json:"Bitrate"`
		FileHash         string  `json:"FileHash"`
		Size             int64   `json:"Size"`
		Height           int     `json:"Height"`
		Width            int     `json:"Width"`
		Format           string  `json:"Format"`
		Codec            string  `json:"Codec"`
		Logo             string  `json:"Logo"`
		Definition       string  `json:"Definition"`
		Quality          string  `json:"Quality"`
		Duration         float64 `json:"Duration"`
		EncryptionMethod string  `json:"EncryptionMethod"`
		PlayAuth         string  `json:"PlayAuth"`
		PlayAuthID       string  `json:"PlayAuthID"`
		MainPlayURL      string  `json:"MainPlayUrl"`
		BackupPlayURL    string  `json:"BackupPlayUrl"`
		URLExpire        int     `json:"UrlExpire"`
		FileID           string  `json:"FileID"`
		P2PVerifyURL     string  `json:"P2pVerifyURL"`
		PreloadInterval  int     `json:"PreloadInterval"`
		PreloadMaxStep   int     `json:"PreloadMaxStep"`
		PreloadMinStep   int     `json:"PreloadMinStep"`
		PreloadSize      int     `json:"PreloadSize"`
		CheckInfo        string  `json:"CheckInfo"`
	} `json:"PlayInfoList"`
	TotalCount int `json:"TotalCount"`
}

var (
	videoCaixinComIdReg    = regexp.MustCompile(`initPlayer\('(.*?)'`)
	videoCaixinComTopicReg = regexp.MustCompile(`topic="(.*?)"`)
)

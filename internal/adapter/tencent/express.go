package tencent

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	MusicExpressAPI = "http://base.music.qq.com/fcgi-bin/fcg_musicexpress.fcg"
)

type MusicExpress struct {
	Code    int      `json:"code"`
	SIP     []string `json:"sip"`
	ThirdIP []string `json:"thirdip"`
	Key     string   `json:"key"`
}

func createGuid() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Intn(10000000000-1000000000) + 1000000000)
}

func getVkey(guid string) (string, error) {
	query := map[string]string{
		"guid":   guid,
		"format": "json",
	}

	var m MusicExpress
	err := tencentRequest(query, MusicExpressAPI, &m)
	if err != nil {
		return "", err
	} else if m.Code != 0 {
		return "", fmt.Errorf("%d", m.Code)
	}

	return m.Key, nil
}

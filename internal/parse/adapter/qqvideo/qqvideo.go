package qqvideo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/chyroc/dl/internal/config"
)

// copy from: https://github.com/iawia002/lux
// license: MIT
type qqVideoInfo struct {
	Fl struct {
		Fi []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Cname string `json:"cname"`
			Fs    int    `json:"fs"`
		} `json:"fi"`
	} `json:"fl"`
	Vl struct {
		Vi []struct {
			Fn    string `json:"fn"`
			Ti    string `json:"ti"`
			Fvkey string `json:"fvkey"`
			Cl    struct {
				Fc int `json:"fc"`
				Ci []struct {
					Idx int `json:"idx"`
				} `json:"ci"`
			} `json:"cl"`
			Ul struct {
				UI []struct {
					URL string `json:"url"`
				} `json:"ui"`
			} `json:"ul"`
		} `json:"vi"`
	} `json:"vl"`
	Msg string `json:"msg"`
}

type qqKeyInfo struct {
	Key string `json:"key"`
}

const qqPlayerVersion string = "3.2.19.333"

func getVinfo(vid, defn, refer string) (qqVideoInfo, error) {
	url := fmt.Sprintf("http://vv.video.qq.com/getinfo?otype=json&platform=11&defnpayver=1&appver=%s&defn=%s&vid=%s", qqPlayerVersion, defn, vid)
	html, err := config.ReqCli.New(http.MethodGet, url).Text()
	if err != nil {
		return qqVideoInfo{}, err
	}
	jsonStrings := matchOneOf(html, `QZOutputJson=(.+);$`)
	if jsonStrings == nil || len(jsonStrings) < 2 {
		return qqVideoInfo{}, fmt.Errorf("get vinfo failed")
	}
	jsonString := jsonStrings[1]
	var data qqVideoInfo
	if err = json.Unmarshal([]byte(jsonString), &data); err != nil {
		return qqVideoInfo{}, err
	}
	return data, nil
}

type Stream struct {
	Quality string  `json:"quality"`
	Parts   []*Part `json:"parts"`
	Size    int64   `json:"size"`
}

// Part is the data structure for a single part of the video stream information.
type Part struct {
	URL  string `json:"url"`
	Size int64  `json:"size"`
	Ext  string `json:"ext"`
}

func genStreams(vid, cdn string, data qqVideoInfo) (map[string]*Stream, error) {
	streams := make(map[string]*Stream)
	var vkey string
	// number of fragments
	var clips int

	for _, fi := range data.Fl.Fi {
		var fmtIDPrefix string
		var fns []string
		if fi.Name == "shd" || fi.Name == "fhd" {
			fmtIDPrefix = "p"
			fmtIDName := fmt.Sprintf("%s%d", fmtIDPrefix, fi.ID%10000)
			fns = []string{strings.Split(data.Vl.Vi[0].Fn, ".")[0], fmtIDName, "mp4"}
			if len(fns) > 3 {
				// delete ID part
				// e0765r4mwcr.2.mp4 -> e0765r4mwcr.mp4
				fns = append(fns[:1], fns[2:]...)
			}
			clips = data.Vl.Vi[0].Cl.Fc
			if clips == 0 {
				clips = 1
			}
		} else {
			tmpData, err := getVinfo(vid, fi.Name, cdn)
			if err != nil {
				return nil, err
			}
			fns = strings.Split(tmpData.Vl.Vi[0].Fn, ".")
			if len(fns) >= 3 && matchOneOf(fns[1], `^p(\d{3})$`) != nil {
				fmtIDPrefix = "p"
			}
			clips = tmpData.Vl.Vi[0].Cl.Fc
			if clips == 0 {
				clips = 1
			}
		}

		var urls []*Part
		var totalSize int64
		var filename string
		for part := 1; part < clips+1; part++ {
			// Multiple fragments per streams
			if fmtIDPrefix == "p" {
				if len(fns) < 4 {
					// If the number of fragments > 0, the filename needs to add the number of fragments
					// n0687peq62x.p709.mp4 -> n0687peq62x.p709.1.mp4
					fns = append(fns[:2], append([]string{strconv.Itoa(part)}, fns[2:]...)...)
				} else {
					fns[2] = strconv.Itoa(part)
				}
			}
			filename = strings.Join(fns, ".")
			url := fmt.Sprintf("http://vv.video.qq.com/getkey?otype=json&platform=11&appver=%s&filename=%s&format=%d&vid=%s", qqPlayerVersion, filename, fi.ID, vid)
			html, err := config.ReqCli.New(http.MethodGet, url).Text()
			if err != nil {
				return nil, err
			}
			jsonStrings := matchOneOf(html, `QZOutputJson=(.+);$`)
			if jsonStrings == nil || len(jsonStrings) < 2 {
				return nil, fmt.Errorf("parse video fail")
			}
			jsonString := jsonStrings[1]

			var keyData qqKeyInfo
			if err = json.Unmarshal([]byte(jsonString), &keyData); err != nil {
				return nil, err
			}

			vkey = keyData.Key
			if vkey == "" {
				vkey = data.Vl.Vi[0].Fvkey
			}
			realURL := fmt.Sprintf("%s%s?vkey=%s", cdn, filename, vkey)
			readResp, err := config.ReqCli.New(http.MethodGet, realURL).Response()
			// size, err := request.Size(realURL, cdn)
			if err != nil {
				return nil, err
			}
			urlData := &Part{
				URL:  realURL,
				Size: readResp.ContentLength,
				Ext:  "mp4",
			}
			urls = append(urls, urlData)
			totalSize += readResp.ContentLength
		}
		streams[fi.Name] = &Stream{
			Parts:   urls,
			Size:    totalSize,
			Quality: fi.Cname,
		}
	}
	return streams, nil
}

type extractor struct{}

// New returns a qq extractor.
func New() *extractor {
	return &extractor{}
}

// Data is the main data structure for the whole video data.
type Data struct {
	Title   string             `json:"title"`
	Streams map[string]*Stream `json:"streams"`
}

// Extract is the main function to extract the data.
func (e *extractor) Extract(url string) (*Data, error) {
	vids := matchOneOf(url, `vid=(\w+)`, `/(\w+)\.html`)
	if vids == nil || len(vids) < 2 {
		return nil, fmt.Errorf("fail")
	}
	vid := vids[1]

	data, err := getVinfo(vid, "fhds", url)
	if err != nil {
		return nil, err
	}

	// API request error
	if data.Msg != "" {
		return nil, fmt.Errorf(data.Msg)
	}
	cdn := data.Vl.Vi[0].Ul.UI[0].URL
	streams, err := genStreams(vid, cdn, data)
	if err != nil {
		return nil, err
	}

	return &Data{
		Title:   data.Vl.Vi[0].Ti + ".mp4",
		Streams: streams,
	}, nil
}

func matchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		// (?flags): set flags within current group; non-capturing
		// s: let . match \n (default false)
		// https://github.com/google/re2/wiki/Syntax
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

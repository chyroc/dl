package parse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/chyroc/dl/internal/config"
	"github.com/chyroc/dl/internal/resource"
)

func NewWwwMissevanCom() Parser {
	return &wwwMissevanCom{}
}

type wwwMissevanCom struct{}

func (r *wwwMissevanCom) Kind() string {
	return "www.missevan.com"
}

func (r *wwwMissevanCom) ExampleURLs() []string {
	return []string{"https://www.missevan.com/sound/player?id=1303686"}
}

func (r *wwwMissevanCom) Parse(uri string) (resource.Resource, error) {
	videoID, err := r.getVideoID(uri)
	if err != nil {
		return nil, err
	}

	{
		resp, err := r.getVideoList(videoID)
		if err != nil {
			return nil, err
		} else if resp != nil {
			urls := make([]string, len(resp.Info.Episodes.Episode))
			wg := new(sync.WaitGroup)
			var finalErr error
			for i := range resp.Info.Episodes.Episode {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					if finalErr != nil {
						return
					}
					v := resp.Info.Episodes.Episode[i]
					resp, err := r.getVideoSingle(strconv.FormatInt(int64(v.SoundID), 10))
					if err != nil {
						finalErr = err
						return
					}
					urls[i] = resp.Info.Sound.Soundurl
				}(i)
			}
			wg.Wait()

			chapters := []resource.Resource{}
			for idx, v := range resp.Info.Episodes.Episode {
				// sid := strconv.FormatInt(v.SoundID, 10)
				chapters = append(chapters, resource.NewURL(fmt.Sprintf("%s_%d.mp3", v.Soundstr, v.SoundID), urls[idx]))
			}
			return resource.NewURLChapter(resp.Info.Drama.Name, chapters), nil
		}
	}

	resp, err := r.getVideoSingle(videoID)
	if err != nil {
		return nil, err
	}
	title := fmt.Sprintf("%s_%d", resp.Info.Sound.Soundstr, resp.Info.User.ID)
	return resource.NewURL(title+".mp3", resp.Info.Sound.Soundurl), nil
}

func (r *wwwMissevanCom) getVideoID(uri string) (string, error) {
	uriParsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	return uriParsed.Query().Get("id"), nil
}

func (r *wwwMissevanCom) getVideoList(videoID string) (*wwwMissevanComGetVideoListResp, error) {
	uri := fmt.Sprintf("https://www.missevan.com/dramaapi/getdramabysound?sound_id=%s", videoID)
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}
	errResp := new(wwwMissevanComErrResp)
	if err = json.Unmarshal([]byte(text), errResp); err == nil {
		return nil, nil
	}
	resp := new(wwwMissevanComGetVideoListResp)
	return resp, json.Unmarshal([]byte(text), resp)
}

func (r *wwwMissevanCom) getVideoSingle(videoID string) (*wwwMissevanComGetVideoSingleResp, error) {
	uri := fmt.Sprintf("https://www.missevan.com/sound/getsound?soundid=%s", videoID)
	text, err := config.ReqCli.New(http.MethodGet, uri).Text()
	if err != nil {
		return nil, err
	}
	errResp := new(wwwMissevanComErrResp)
	if err = json.Unmarshal([]byte(text), errResp); err == nil {
		return nil, nil
	}
	resp := new(wwwMissevanComGetVideoSingleResp)
	return resp, json.Unmarshal([]byte(text), resp)
}

type wwwMissevanComErrResp struct {
	Info string `json:"info"`
}

type wwwMissevanComGetVideoSingleResp struct {
	Info struct {
		Sound struct {
			ID                 int    `json:"id"`
			CatalogID          int    `json:"catalog_id"`
			CreateTime         int    `json:"create_time"`
			LastUpdateTime     int    `json:"last_update_time"`
			Duration           int    `json:"duration"`
			UserID             int    `json:"user_id"`
			Username           string `json:"username"`
			CoverImage         string `json:"cover_image"`
			Soundstr           string `json:"soundstr"`
			Intro              string `json:"intro"`
			Soundurl           string `json:"soundurl"`
			Soundurl128        string `json:"soundurl_128"`
			FrontCover         string `json:"front_cover"`
			AllComments        int    `json:"all_comments"`
			CommentsNum        int    `json:"comments_num"`
			ViewCountFormatted string `json:"view_count_formatted"`
			ForbiddenComment   int    `json:"forbidden_comment"`
			MosaicURL          string `json:"mosaic_url"`
			Breadcrumb         string `json:"breadcrumb"`
		} `json:"sound"`
		User struct {
			ID            int    `json:"id"`
			Intro         string `json:"intro"`
			Username      string `json:"username"`
			Icon          string `json:"icon"`
			Fansnum       string `json:"fansnum"`
			Soundnum      int    `json:"soundnum"`
			Authenticated int    `json:"authenticated"`
		} `json:"user"`
	} `json:"info"`
}

type wwwMissevanComGetVideoListResp struct {
	Info struct {
		Drama struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			Cover           string `json:"cover"`
			CoverColor      int    `json:"cover_color"`
			Abstract        string `json:"abstract"`
			UserID          int    `json:"user_id"`
			Username        string `json:"username"`
			Checked         int    `json:"checked"`
			ViewCount       int    `json:"view_count"`
			CommentCount    int    `json:"comment_count"`
			Catalog         int    `json:"catalog"`
			SubscriptionNum int    `json:"subscription_num"`
			ShowRevenue     int    `json:"show_revenue"`
			NameLetters     string `json:"name_letters"`
			IntegrityName   string `json:"integrity_name"`
			CatalogName     string `json:"catalog_name"`
		} `json:"drama"`
		Episodes struct {
			Episode []struct {
				ID           int    `json:"id"`
				Name         string `json:"name"`
				DramaID      int    `json:"drama_id"`
				SoundID      int64  `json:"sound_id"`
				Date         int    `json:"date"`
				Order        int    `json:"order"`
				CreateTime   int    `json:"create_time"`
				ModifiedTime int    `json:"modified_time"`
				Video        int    `json:"video"`
				NeedPay      int    `json:"need_pay"`
				Soundstr     string `json:"soundstr"`
				Duration     int    `json:"duration"`
				Checked      int    `json:"checked"`
			} `json:"episode"`
		} `json:"episodes"`
		Current struct {
			ID           int    `json:"id"`
			Name         string `json:"name"`
			DramaID      int    `json:"drama_id"`
			SoundID      int    `json:"sound_id"`
			Date         int    `json:"date"`
			Order        int    `json:"order"`
			CreateTime   int    `json:"create_time"`
			ModifiedTime int    `json:"modified_time"`
			Video        int    `json:"video"`
		} `json:"current"`
	} `json:"info"`
}

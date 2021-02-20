package utils

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Movies ST data
type Movies struct {
	ID         string
	Order      string
	Title      string
	MovieURL   string
	PictureURL string
	Date       string
}

func dateFormat(oriTime string) (newTime string) {
	timeParse, err := time.Parse("2006-01-02T15:04:05Z07:00", oriTime)
	if err != nil {
		log.Fatal(" [API - Format Time]: ", err)
	}
	newTime = timeParse.Format("2006-01-02 15:04:05")
	return
}

func urlFormat(oriURL string) (newURL string) {
	decodeURL, err := url.QueryUnescape(oriURL)
	if err != nil {
		log.Fatal(" [API - Format URL]: ", err)
	}
	rep1 := strings.Replace(decodeURL, "ulizasekailab", "https", -1)
	newURL = strings.Replace(rep1, "videoquery=", "", -1)
	return
}

func float2string(oriFloat float64) (newData string) {
	newData = strconv.FormatFloat(oriFloat, 'f', -1, 64)
	return
}

// STToken app token
func STToken() (token string) {
	userURL := "https://st-api.st-channel.jp/v1/users"
	headers := MiniHeaders{
		"User-Agent": UserAgent,
	}
	res := Minireq.Post(userURL, headers)
	data := res.RawJSON()
	token = data.(map[string]interface{})["api_token"].(string)
	return
}

// STInfo app main info
func STInfo(token string, sinceID string) (stdata []*Movies) {
	mainURL := "https://st-api.st-channel.jp/v1/movies?"

	headers := MiniHeaders{
		"User-Agent":    UserAgent,
		"authorization": "Bearer " + token,
	}

	var sinceOrder string
	if sinceID != "" {
		sinceOrder = "999"
	} else {
		sinceID = "0"
		sinceOrder = "0"
	}

	params := MiniParams{
		"since_id":    sinceID,
		"device_type": "2",
		"since_order": sinceOrder,
		"sort":        "order",
	}

	res := Minireq.Get(mainURL, headers, params)
	data := res.RawJSON()

	movies := data.(map[string]interface{})["movies"].([]interface{})
	for _, movie := range movies {
		mData := movie.(map[string]interface{})
		movidesData := new(Movies)
		movidesData.ID = float2string(mData["id"].(float64))
		movidesData.Order = float2string(mData["order"].(float64))
		movidesData.Title = strings.TrimSpace(mData["title"].(string))
		movidesData.MovieURL = urlFormat(mData["movie_url_everyone"].(string))
		movidesData.PictureURL = mData["thumbnail_path"].(string)
		movidesData.Date = dateFormat(mData["publish_start_date"].(string))
		stdata = append(stdata, movidesData)
	}
	return
}

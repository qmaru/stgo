package utils

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"
)

func init() {
    log.SetFlags(log.Lshortfile | log.LstdFlags)
}

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
    req, err := http.NewRequest("POST", userURL, nil)
    if err != nil {
        log.Fatal(" [API - Token Request Error]: ", err)
    }
    req.Header.Add("User-Agent", DefaultUA())
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(" [API - Token Response Error]: ", err)
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatal(" [API - Token Body Error]: ", err)
    }
    data := make(map[string]interface{})
    json.Unmarshal(body, &data)
    token = data["api_token"].(string)
    return
}

// STInfo app main info
func STInfo(token string, sinceID string) (stdata []*Movies) {
    mainURL := "https://st-api.st-channel.jp/v1/movies?"

    req, err := http.NewRequest("GET", mainURL, nil)
    if err != nil {
        log.Fatal(" [API - Info Request Error]: ", err)
    }
    req.Header.Add("User-Agent", DefaultUA())
    req.Header.Add("authorization", "Bearer "+token)

    var sinceOrder string
    if sinceID != "" {
        sinceOrder = "999"
    } else {
        sinceID = "0"
        sinceOrder = "0"
    }

    params := make(url.Values)
    params.Add("since_id", sinceID)
    params.Add("device_type", "2")
    params.Add("since_order", sinceOrder)
    params.Add("sort", "order")
    req.URL.RawQuery = params.Encode()

    res, err := client.Do(req)
    if err != nil || res.StatusCode == 401 {
        log.Println(" [API - Info Authorization Error]: ", err)
        return
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatal(" [API - Info Body Error]: ", err)
    }
    data := make(map[string]interface{})
    json.Unmarshal(body, &data)

    movies := data["movies"].([]interface{})
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

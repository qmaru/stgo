package utils

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path"
    "regexp"
    "strconv"
    "strings"
    "time"
)

func init() {
    log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func minireq(url string) (data []byte) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(" [HLS - Request Error]: ", err)
    }
    req.Header.Add("User-Agent", DefaultUA())
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(" [HLS - Response Error]: ", err)
    }
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatal(" [HLS - Body Error]: ", err)
    }
    data = body
    return
}

func videoHOST(url string) (host string) {
    urlPart := strings.Split(url, "/")
    host = "http://" + urlPart[2] + "/" + urlPart[3] + "/"
    return
}

func rProgress(i int, amp float64) {
    progress := float64(i) * amp
    num := int(progress / 10)
    if num < 1 {
        num = 1
    }
    pstyle := strings.Repeat(">", num)
    log.Printf("%s [%.2f %%]\r", pstyle, progress)
}

// HLSDownload get hls data
func HLSDownload(hlsURL string, savedir string) {
    filename := time.Now().Format("20060102150405") + ".ts"
    savePath := path.Join(savedir, filename)
    videoFile, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    if err != nil {
        log.Fatal(" [HLS - Download Open File Error]: ", err)
    }

    hlsmain := minireq(hlsURL)

    regBestRule := regexp.MustCompile(`(?m)^[\w\-\.\/\:\?\&\=\%\,\+]+`)
    bestURL := regBestRule.FindAllString(string(hlsmain), -1)[0]
    m3u8Raw := string(minireq(bestURL))

    regKeyRule := regexp.MustCompile(`"(.*?)"`)
    regURLRule := regexp.MustCompile(`(?m)^[\w\-\.\/\:\?\&\=\%\,\+]+`)
    keyRaw := regKeyRule.FindAllString(m3u8Raw, -1)
    keyurl := strings.Replace(keyRaw[0], "\"", "", -1)
    vurls := regURLRule.FindAllString(m3u8Raw, -1)

    vhost := videoHOST(bestURL)
    keyBytes := minireq(keyurl)

    total := float64(len(vurls))
    part, err := strconv.ParseFloat(fmt.Sprintf("%.5f", 100.0/total), 64)
    if err != nil {
        log.Fatal(" [HLS - Download Float to String Error]: ", err)
    }
    for i, url := range vurls {
        i = i + 1
        videoBytes := minireq(vhost + url)
        decrtVideo := AesDecrypt(videoBytes, keyBytes)
        offset, err := videoFile.Seek(0, os.SEEK_END)
        if err != nil {
            log.Fatal(" [HLS - Download File Seek Error]: ", err)
        }
        videoFile.WriteAt(decrtVideo, offset)
        rProgress(i, part)
    }
    defer videoFile.Close()
}

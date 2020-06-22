package main

import (
    "fmt"
    "stgo/utils"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path"
    "strings"

    assetfs "github.com/elazarl/go-bindata-assetfs"
    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/static"
    "github.com/gin-gonic/gin"
    "github.com/zserge/lorca"
)

// var staticDir = "./webui/build"
var savedir string = createDownload()
var movies *utils.Movies

type binaryFileSystem struct {
    fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
    return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
    if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
        if _, err := b.fs.Open(p); err != nil {
            return false
        }
        return true
    }
    return false
}

// BinaryFileSystem bindata2gin
func BinaryFileSystem(root string) *binaryFileSystem {
    fs := &assetfs.AssetFS{
        Asset: utils.Asset,
        AssetDir: utils.AssetDir,
        AssetInfo: utils.AssetInfo,
        Prefix: root,
    }
    return &binaryFileSystem{fs}
}

type stURL struct {
    URL string
}

func init() {
    log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func checkFolder(cpath string) bool {
    _, err := os.Stat(cpath)
    if err != nil && os.IsNotExist(err) {
        return false
    }
    return true
}

func createDownload() string {
    workPath, err := os.Getwd()
    if err != nil {
        log.Fatal(" [Main - Download Folder Get Workdir]: ", err)
    }
    downloadPath := path.Join(workPath, "download")
    if checkFolder(downloadPath) {
        return downloadPath
    }
    os.Mkdir(downloadPath, os.ModePerm)
    return downloadPath
}

func getToken(new bool) string {
    workPath, err := os.Getwd()
    if err != nil {
        log.Fatal(" [Main - Token Folder Get Workdir]: ", err)
    }
    tokenTmp := path.Join(workPath, ".token")
    if new {
        token := utils.STToken()
        ioutil.WriteFile(tokenTmp, []byte(token), 0666)
        return token
    }
    if checkFolder(tokenTmp) {
        token, err := ioutil.ReadFile(tokenTmp)
        if err != nil {
            log.Fatal(" [Main - Token Read File]: ", err)
        }
        return string(token)
    }
    token := utils.STToken()
    ioutil.WriteFile(tokenTmp, []byte(token), 0666)
    return token
}

func stList(c *gin.Context) {
    var status int
    token := getToken(false)
    sinceID := c.DefaultQuery("next_id", "")

    stinfo := utils.STInfo(token, sinceID)
    status = 1
    if utils.CheckNil(stinfo) {
        newToken := getToken(true)
        stinfo = utils.STInfo(newToken, sinceID)
        if utils.CheckNil(stinfo) {
            stinfo = []*utils.Movies{}
            status = 0
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  status,
        "data":    stinfo,
        "message": "STChannel Movie List",
    })
}

func stdown(c *gin.Context) {
    c.Header("Content-Type", "application/json")

    var sturl stURL
    c.BindJSON(&sturl)
    playlist := sturl.URL

    if playlist == "" {
        c.JSON(http.StatusOK, gin.H{
            "status":  0,
            "message": "No Playlist",
        })
    } else {
        utils.HLSDownload(playlist, savedir)
        c.JSON(http.StatusOK, gin.H{
            "status":  1,
            "message": "Finished",
        })
    }
}

func runUI() error {
    ui, err := lorca.New("", "", 800, 600)
    if err != nil {
        return err
    }

    ui.Load("http://localhost:61800")
    <-ui.Done()
    return nil
}

func runServer() {
    listenAddr := "localhost:61800"
    log.Println("Listen: " + listenAddr)
    gin.SetMode(gin.ReleaseMode)
    // gin.SetMode(gin.DebugMode)

    // file, _ := os.Create("access.log")
    // gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
    gin.DefaultWriter = os.Stdout

    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"*"}
    config.AllowMethods = []string{"POST", "OPTION"}
    router := gin.New()
    router.Use(cors.New(config))
    router.Use(gin.Recovery())
    router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s - %s %s %d %s %s\n",
            param.TimeStamp.Format("2006-01-02 15:04:05"),
            param.Method,
            param.Path,
            param.StatusCode,
            param.Latency,
            param.ErrorMessage,
        )
    }))
    // router.Use(static.Serve("/", static.LocalFile(staticDir, true)))
    router.Use(static.Serve("/", BinaryFileSystem("webui/build")))

    router.GET("/stlist", stList)
    router.POST("/stdown", stdown)
    router.Run(listenAddr)
}

func main() {
    go runServer()

    err := runUI()
    if err != nil {
        log.Fatal(err)
    }
}

package apis

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// StaticHand 静态文件
func StaticHand(static fs.FS) gin.HandlerFunc {
	return func(c *gin.Context) {
		upath := c.Request.URL.Path
		if !strings.HasPrefix(upath, "/stlist") {
			http.FileServer(http.FS(static)).ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

// Server 接口服务
func Server(addr string, static fs.FS) {
	log.Println("Listen: " + addr)
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	gin.DefaultWriter = os.Stdout

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "OPTION"}
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

	router.Use(StaticHand(static))
	router.GET("/stlist", STList)
	router.Run(addr)
}

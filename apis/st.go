package apis

import (
	"io/ioutil"
	"net/http"
	"os"
	"stgo/utils"

	"github.com/gin-gonic/gin"
)

var movies *utils.Movies
var tokenFile *os.File

func init() {
	token := utils.STToken()
	tokenFile, _ = ioutil.TempFile("", "st-token-")
	tokenFile.Write([]byte(token))
}

func checkFile(cpath string) bool {
	_, err := os.Stat(cpath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// STList 获取电影列表
func STList(c *gin.Context) {
	if checkFile(tokenFile.Name()) {
		tokenByte, _ := os.ReadFile(tokenFile.Name())
		sinceID := c.DefaultQuery("next_id", "")
		stinfo := utils.STInfo(string(tokenByte), sinceID)

		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"data":    stinfo,
			"message": "STChannel Movie List",
		})
	}
}

package main

import (
	"embed"
	"io/fs"
	"log"

	"stgo/apis"

	"github.com/zserge/lorca"
)

//go:embed webui/build
var build embed.FS

// SPAIndex 载入静态文件
func SPAIndex() fs.FS {
	fsys := fs.FS(build)
	buildStatic, _ := fs.Sub(fsys, "webui/build")
	return buildStatic
}

// RunUI 图形界面
func RunUI(addr string) error {
	ui, err := lorca.New("", "", 480, 800)
	if err != nil {
		return err
	}

	ui.Load("http://" + addr)
	<-ui.Done()
	return nil
}

func main() {
	listenAddr := "127.0.0.1:61800"

	go apis.Server(listenAddr, SPAIndex())

	err := RunUI(listenAddr)
	if err != nil {
		log.Fatal(err)
	}
}

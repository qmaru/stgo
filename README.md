# stgo

简单下载 [STchannel](https://play.google.com/store/apps/details?id=jp.co.shueisha.stchannel) 内的视频，所有下载的视频将保存在执行目录。  

+ Golang / [gin](https://github.com/gin-gonic/gin) 提供数据接口
+ React 构建用户界面
+ [lorca](https://github.com/zserge/lorca) 连接 Golang 和 React
+ [go-bindata](https://github.com/go-bindata/go-bindata) 和 [go-bindata-assetfs](https://github.com/elazarl/go-bindata-assetfs) 打包静态文件

#### 运行环境

+ Golang 1.13
+ Node.js 12 LTS
+ React 16
+ material-ui 4.6

可选
+ [upx](https://upx.github.io/)

#### 项目打包
```go
// webui 目录下打包静态文件
yarn build

// 编译静态文件
go-bindata -o utils/bindata.go -pkg=utils webui/build/...

// 生成单文件 隐藏执行窗口
go build -ldflags "-H windowsgui"

// 生成最小单文件
go build -ldflags "-s -w -H windowsgui"

// upx压缩文件（可选）
upx stgo
```

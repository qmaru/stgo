# stgo

获取 [STchannel](https://play.google.com/store/apps/details?id=jp.co.shueisha.stchannel) 内的视频列表。

+ Golang / [gin](https://github.com/gin-gonic/gin) 提供数据接口
+ React 构建用户界面
+ [lorca](https://github.com/zserge/lorca) 连接 Golang 和 React

#### 运行环境

+ Golang 1.16
+ Node.js 12 LTS
+ React 17
+ material-ui 4.11

可选
+ [upx](https://upx.github.io/)

#### 项目打包
```go
// webui 目录下打包静态文件
yarn build

// 生成单文件 隐藏执行窗口
go build -ldflags "-H windowsgui"

// 生成最小单文件
go build -ldflags "-s -w -H windowsgui"

// upx压缩文件（可选）
upx stgo
```

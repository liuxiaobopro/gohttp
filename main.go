package main

import (
	"gohttp/middleware"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化一个gin实例
	r := gin.Default()

	r.Use(middleware.Trace())
	r.Static("static", "front/static")
	r.LoadHTMLGlob("front/view/**/*")

	r.GET("/", ViewIndexIndex)

	r.GET("/log", ViewLogIndex)

	// 监听服务
	r.Run(":9000")
}

func ViewIndexIndex(c *gin.Context) {
	resData := make(gin.H, 0)
	resData["time"] = time.Now().Format("2006-01-02 15:04:05")
	c.HTML(http.StatusOK, "index/index.html", resData)
}

func ViewLogIndex(c *gin.Context) {
	resData := make(gin.H, 0)
	// 获取testlog目录下所有文件和文件夹
	files, err := ioutil.ReadDir("./testlog")
	if err != nil {
		logrus.Errorf("ReadDir failed, err:%v \n", err)
	}

	file1 := make([]string, 0)
	file2 := make([]string, 0)
	for _, v := range files {
		if v.IsDir() {
			file1 = append(file1, v.Name())
		} else {
			file2 = append(file2, v.Name())
		}
	}

	file1 = append(file1, file2...)

	resData["files"] = file1

	c.HTML(http.StatusOK, "log/index.html", resData)
}

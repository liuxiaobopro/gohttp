package main

import (
	"fmt"
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
	r.LoadHTMLGlob("view/**/*")

	r.GET("/", func(c *gin.Context) {
		logrus.Infof("Hello World! %s", time.Now())
		c.JSON(http.StatusOK, gin.H{
			"message": time.Now(),
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/log", func(c *gin.Context) {
		c.HTML(http.StatusOK, "log/index.html", reeaDir)
	})

	// 监听服务
	r.Run(":9000")
}

func reeaDir() *gin.H {
	// 获取testlog目录下所有文件和文件夹
	files, err := ioutil.ReadDir("./testlog")
	if err != nil {
		logrus.Errorf("ReadDir failed, err:%v \n", err)
	}

	for _, v := range files {
		if v.IsDir() {
			fmt.Println("dir:", v.Name())
		} else {
			fmt.Println("file:", v.Name())
		}
	}

	return &gin.H{
		"message": "pong",
	}
}

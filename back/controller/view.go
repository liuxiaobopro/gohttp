package controller

import (
	"gohttp/back/common/tools"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ViewIndexIndex(c *gin.Context) {
	resData := make(gin.H, 0)
	resData["time"] = time.Now().Format("2006-01-02 15:04:05")
	c.HTML(http.StatusOK, "index.html", resData)
}

func ViewLogIndex(c *gin.Context) {
	resData := make(gin.H, 0)

	//#region 读取文件
	lines, err := tools.ReadFileLine("D:\\1_liuxiaobo\\testlog\\log.txt", 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	//#endregion

	resData["code"] = lines

	c.HTML(http.StatusOK, "log/index.html", resData)
}

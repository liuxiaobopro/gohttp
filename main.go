package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"gohttp/back/middleware"
	modelRes "gohttp/back/model/res"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	r.Use(middleware.Trace())
	r.Static("static", "./front/static")
	r.LoadHTMLGlob("front/view/**/*")

	//#region 页面
	rg1 := r.Group("")
	{
		rg1.GET("/", ViewIndexIndex)
		rg1.GET("/log", ViewLogIndex)
	}
	//#endregion

	//#region api
	rg2 := r.Group("/api")
	{
		rg2.GET("/show_fold", ShowFolds)
	}
	//#endregion

	r.Run(":9000")
}

func ViewIndexIndex(c *gin.Context) {
	resData := make(gin.H, 0)
	resData["time"] = time.Now().Format("2006-01-02 15:04:05")
	c.HTML(http.StatusOK, "index/index.html", resData)
}

func ViewLogIndex(c *gin.Context) {
	resData := make(gin.H, 0)

	c.HTML(http.StatusOK, "log/index.html", resData)
}

func ShowFolds(c *gin.Context) {
	var (
		menuList      = make([]*modelRes.LogIndexShowFoldsRes, 0)
		menuChildList = make([]*modelRes.LogIndexShowFoldsRes, 0)
	)

	//#region 获取参数
	pathParam := c.Query("path")
	//#endregion

	if pathParam == "" {
		//#region 获取所有盘符
		drives := []string{}
		for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
			path := string(drive) + ":\\"
			_, err := os.Stat(path)
			if err == nil {
				drives = append(drives, path)
			}
		}
		//#endregion

		//#region 整理数据
		for k, drive := range drives {
			// 读取目录中的所有文件和子目录
			entries, err := os.ReadDir(drive)
			if err != nil {
				logrus.Errorf("read dir error: %v", err)
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			var (
				menuChildList1 = make([]*modelRes.LogIndexShowFoldsRes, 0)
				menuChildList2 = make([]*modelRes.LogIndexShowFoldsRes, 0)
			)

			// 打印每个子目录的名称
			for k1, entry := range entries {
				var (
					title string
				)

				if entry.IsDir() {
					title = "<span class=\"layui-badge layui-bg-orange\">F</span>" + entry.Name()
					menuChildList1 = append(menuChildList1, &modelRes.LogIndexShowFoldsRes{
						Title:  title,
						Id:     k1 + 1,
						PathId: fmt.Sprintf("%d-%d", k, k1),
					})
				} else {
					title = "<span class=\"layui-badge layui-bg-green\">D</span>" + entry.Name()
					menuChildList2 = append(menuChildList2, &modelRes.LogIndexShowFoldsRes{
						Title:  title,
						Id:     k1 + 1,
						PathId: fmt.Sprintf("%d-%d", k, k1),
					})
				}
			}

			menuChildList = append(menuChildList1, menuChildList2...)

			menuList = append(menuList, &modelRes.LogIndexShowFoldsRes{
				Title:    drive,
				Id:       k,
				PathId:   strconv.Itoa(k),
				Children: menuChildList,
			})
		}
		//#endregion
	}

	c.JSON(http.StatusOK, menuList)
}

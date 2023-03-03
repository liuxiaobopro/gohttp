package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"gohttp/back/middleware"
	modelReq "gohttp/back/model/req"
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
		rg2.GET("/open_fold", OpenFold)
		rg2.GET("/read_file", ReadFile)
	}
	//#endregion

	// fp := "D:\\1_project\\github\\PowerWechatTutorial\\log.txt"

	// // 写入10w行数据
	// f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	fmt.Println("open file err=", err)
	// 	return
	// }
	// defer f.Close()

	// // 循环写入
	// for i := 1; i <= 10000; i++ {
	// 	s := fmt.Sprintf("%010d ## \r\n", i)
	// 	f.WriteString(s)
	// }

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

func OpenFold(c *gin.Context) {
	//#region 获取参数
	pathParam := c.Query("filepath")
	if pathParam == "" {
		c.JSON(http.StatusBadRequest, "path is empty")
		return
	}
	//#endregion

	//#region 读取目录中的所有文件和子目录
	entries, err := os.ReadDir(pathParam)
	if err != nil {
		logrus.Errorf("read dir error: %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var (
		menuList      = make([]*modelRes.LogIndexShowFoldsRes, 0)
		menuChildList = make([]*modelRes.LogIndexShowFoldsRes, 0)
	)
	//#endregion

	//#region 整理数据
	for k, entry := range entries {
		var (
			title string
			id    = strconv.Itoa(k)
		)

		if entry.IsDir() {
			// title = "<span class=\"layui-badge layui-bg-orange\">F</span>" + entry.Name()
			// menuChildList = append(menuChildList, &modelRes.LogIndexShowFoldsRes{
			// 	Title:  title,
			// 	Id:     id,
			// 	PathId: fmt.Sprintf("%s-%s", pathParam, id),
			// })
		} else {
			title = "<span class=\"layui-badge layui-bg-green\">D</span> " + entry.Name()
			menuChildList = append(menuChildList, &modelRes.LogIndexShowFoldsRes{
				Title: title,
				Id:    id,
				Path:  pathParam,
				Name:  entry.Name(),
			})
		}
	}

	menuList = append(menuList, &modelRes.LogIndexShowFoldsRes{
		Title:    pathParam,
		Spread:   true,
		Children: menuChildList,
	})
	//#endregion

	c.JSON(http.StatusOK, menuList)
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
			id := strconv.Itoa(k)
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
					id1   = strconv.Itoa(k1)
				)

				if entry.IsDir() {
					title = "<span class=\"layui-badge layui-bg-orange\">F</span> " + entry.Name()
					menuChildList1 = append(menuChildList1, &modelRes.LogIndexShowFoldsRes{
						Title:  title,
						Id:     id1,
						PathId: fmt.Sprintf("%s-%s", id, id1),
					})
				} else {
					title = "<span class=\"layui-badge layui-bg-green\">D</span> " + entry.Name()
					menuChildList2 = append(menuChildList2, &modelRes.LogIndexShowFoldsRes{
						Title:  title,
						Id:     id1,
						PathId: fmt.Sprintf("%s-%s", id, id1),
					})
				}
			}

			menuChildList = append(menuChildList1, menuChildList2...)

			menuList = append(menuList, &modelRes.LogIndexShowFoldsRes{
				Title:    drive,
				Id:       id,
				PathId:   id,
				Children: menuChildList,
			})
		}
		//#endregion
	}

	c.JSON(http.StatusOK, menuList)
}

func ReadFile(c *gin.Context) {
	//#region 获取参数
	var r modelReq.LogIndexReadFileReq
	if err := c.ShouldBind(&r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	limit := 1000

	n := (r.Page - 1) * limit
	m := limit

	if r.Path == "" || r.Name == "" {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	//#endregion

	//#region 拼接文件路径
	var filepath string
	if runtime.GOOS == "windows" {
		filepath = fmt.Sprintf("%s\\%s", r.Path, r.Name)
	} else {
		filepath = fmt.Sprintf("%s/%s", r.Path, r.Name)
	}
	//#endregion

	logrus.Infof("filepath: %s", filepath)

	//#region 读取文件
	file, err := os.Open(filepath)
	if err != nil {
		logrus.Errorf("打开文件错误: %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0, m)

	// 跳过前n-1行
	for i := 1; i <= n && scanner.Scan(); i++ {
		// do nothing
	}

	// 读取n到n+m-1行
	for i := 1; i <= m && scanner.Scan(); i++ {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logrus.Errorf("读取文件错误: %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	//#endregion

	c.JSON(http.StatusOK, lines)
}

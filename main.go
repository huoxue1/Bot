package main

import (
	_ "Bot/plugins/All"
	_ "Bot/plugins/Hello"
	_ "Bot/plugins/Robbery"
	"Bot/plugins/daka"
	"Bot/plugins/fileSearch"
	_ "Bot/plugins/flash"
	_ "Bot/plugins/recall"
	_ "Bot/plugins/refresh"
	"fmt"
	Bot "github.com/3343780376/go-mybots"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func main() {

	hand := Bot.Hand()
	handHttp(hand)
	Bot.LoadFilter("./config.json")
	err := hand.Run("0.0.0.0:80")
	if err != nil {
		log.Println("端口错误")
	}
	go daka.Cr()
	log.Println("正在监听")
}

func handHttp(engine *gin.Engine) {
	engine.LoadHTMLFiles("./templete/fiction.html")
	engine.StaticFS("/log", http.Dir("./plugins/logs"))

	engine.POST("/hook", func(context *gin.Context) {
		if strings.Contains(context.Request.Header.Get("User-Agent"), "GitHub") {
			log.Println("开始执行pull ")
			command := exec.Command("/bin/sh", "-c", "git pull")
			err := command.Start()
			if err != nil {
				log.Println("命令执行失败")
			}

		}
	})

	engine.GET("/fiction", func(context *gin.Context) {
		context.HTML(http.StatusOK, "fiction.html", gin.H{"content": fileSearch.File})
	})
	engine.GET("/fiction/:filename", func(context *gin.Context) {
		param := fileSearch.File[context.Param("filename")]
		context.Writer.WriteHeader(http.StatusOK)
		context.Header("Content-Disposition", "attachment;filename="+param)
		context.Header("Content-Type", "application/text/plain")
		file, err := ioutil.ReadFile("./fiction/" + param)
		if err != nil {
			return
		}
		context.Header("Accept-Length", fmt.Sprintf("%d", len(file)))
		_, _ = context.Writer.Write([]byte(file))
	})
}

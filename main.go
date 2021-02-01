package main

import (
	_ "Bot/plugins/All"
	_ "Bot/plugins/Hello"
	_ "Bot/plugins/Robbery"
	"Bot/plugins/daka"
	"Bot/plugins/fileSearch"
	_ "Bot/plugins/flash"
	_ "Bot/plugins/refresh"
	"fmt"
	"github.com/3343780376/go-bot"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	go_bot.Run("127.0.0.1:8000")
	go daka.Cr()
	go_bot.LoadFilter("./config.json")
	handHttp()
}

func handHttp() {
	engine := gin.New()
	engine.LoadHTMLFiles("./templete/fiction.html")
	engine.StaticFS("/log", http.Dir("./plugins/logs"))
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
	err := engine.Run("0.0.0.0:80")
	if err != nil {
		log.Println("端口错误")
	}
}

package main

import (
	_ "Bot/plugins/All"
	_ "Bot/plugins/Hello"
	_ "Bot/plugins/Robbery"
	"Bot/plugins/daka"
	_ "Bot/plugins/refresh"
	"fmt"
	Bot "github.com/3343780376/go-mybots"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {

	hand := Bot.Hand()
	handHttp(hand)
	go daka.Cr()
	Bot.LoadFilter("./config.json")
	err := hand.Run("0.0.0.0:80")
	if err != nil {
		log.Println("端口错误")
	}
	log.Println("正在监听")
}

func handHttp(engine *gin.Engine) {
	getwd, _ := os.Getwd()
	fmt.Println(getwd)
	engine.StaticFS("/log", http.Dir("./plugins/logs"))
}

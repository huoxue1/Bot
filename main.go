package main

import (
	_ "Bot/plugins/All"
	_ "Bot/plugins/Hello"
	_ "Bot/plugins/Robbery"
	_ "Bot/plugins/refresh"
	Bot "github.com/3343780376/go-mybots"
	"log"
)

func main() {
	hand := Bot.Hand()
	err := hand.Run("127.0.0.1:8000")
	if err != nil {
		log.Println("端口错误")
	}
	log.Println("正在监听")
}

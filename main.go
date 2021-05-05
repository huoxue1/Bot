package main

import (
	"fengyeBot/models"
	_ "fengyeBot/plugin/curse"
	_ "fengyeBot/plugin/welcome"
	"github.com/3343780376/leafBot"
	"github.com/3343780376/leafBot/plugins"
	"os"
)

func init() {
	plugins.UseEchoHandle()
}

func main() {
	//go gui.InitWindow()
	models.Xorminit()
	path, _ := os.Getwd()
	leafBot.LoadConfig(path+"/config/config.json", leafBot.JSON)
	leafBot.InitBots()
}

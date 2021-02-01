package flash

import (
	"fmt"
	go_bot "github.com/3343780376/go-bot"
	"regexp"
)

var bot *go_bot.Bot

func botInit() {
	bot = go_bot.GetBot(3343780376)
}

func init() {
	go botInit()
	go_bot.ViewMessage = append(go_bot.ViewMessage, go_bot.ViewMessageApi{OnMessage: Flash,
		MessageType: "", SubType: ""})
}

func Flash(event go_bot.Event) {
	compile := regexp.MustCompile(`\[CQ:image,type=flash,file=(.*?)\]`)
	if compile.MatchString(event.Message) {
		for _, i2 := range compile.FindAllStringSubmatch(event.Message, -1) {
			if event.MessageType == "private" {
				_ = bot.SendPrivateMsg(3180808826,
					fmt.Sprintf("来自私聊消息%v[CQ:image,file=%v]", event.UserId, i2[1]), false)
			} else {
				_ = bot.SendPrivateMsg(3180808826,
					fmt.Sprintf("来自群消息%v的%v所发消息[CQ:image,file=%v]", event.GroupId, event.UserId, i2[1]),
					false)
			}

		}
	}
}

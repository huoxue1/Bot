package refresh

import (
	"Bot/model"
	"github.com/3343780376/go-bot"
)

func init() {
	go_bot.ViewMessage = append(go_bot.ViewMessage, go_bot.ViewMessageApi{OnMessage: Refresh,
		MessageType: go_bot.MessageTypeApi.Group, SubType: ""})
	Num = 0
	UserId = 0
	go botInit()
}

var bot *go_bot.Bot

func botInit() {
	bot = go_bot.GetBot(2177120078)
}

var (
	UserId int
	Num    int
)

func Refresh(event go_bot.Event) {
	if event.SelfId == 3343780376 {
		return
	}
	if event.UserId == UserId {
		Num += 1
	} else {
		Num = 0
	}
	connect := model.DbInit()
	defer connect.Close()
	UserId = event.UserId
	if Num >= 4 {
		connect.Update(-2, event)
		bot.SetGroupBan(event.GroupId, event.UserId, 10*60)
		bot.SendGroupMsg(event.GroupId, "你刷屏了"+bot.MessageAt(event.UserId).Message, false)
		Num = 0
	}
}

package refresh

import (
	"Bot/Integral"
	"github.com/3343780376/go-bot"
	"log"
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
	UserId = event.UserId
	if Num >= 4 {
		xlsx := Integral.Xlsx{Event: event, Sheet: ""}
		err := xlsx.XlsxInit()
		_, err = xlsx.Decrease(2)
		bot.SetGroupBan(event.GroupId, event.UserId, 10*60)
		_ = bot.SendGroupMsg(event.GroupId, "你刷屏了"+bot.MessageAt(event.UserId).Message, false)
		if err != nil {
			log.Println(err)
		}
		Num = 0
	}
}

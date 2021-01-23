package refresh

import (
	"Bot/Integral"
	"github.com/3343780376/go-mybots"
	"log"
)

func init() {
	go_mybots.ViewMessage = append(go_mybots.ViewMessage, go_mybots.ViewMessageApi{OnMessage: Refresh,
		MessageType: go_mybots.MessageTypeApi.Group, SubType: ""})
	Num = 0
	UserId = 0
}

var (
	bot    = go_mybots.Bots{Address: "127.0.0.1", Port: 5700, Admin: 3343780376}
	UserId int
	Num    int
)

func Refresh(event go_mybots.Event) {
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
		err = bot.SetGroupBan(event.GroupId, event.UserId, 10*60)
		_, err = bot.SendGroupMsg(event.GroupId, "你刷屏了"+go_mybots.MessageAt(event.UserId).Message, false)
		if err != nil {
			log.Println(err)
		}
		Num = 0
	}
}

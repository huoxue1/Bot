package flash

import (
	"fmt"
	go_mybots "github.com/3343780376/go-mybots"
	"regexp"
)

var bot = go_mybots.Bots{Address: "127.0.0.1", Port: 5700, Admin: 3343780376}

func init() {
	go_mybots.ViewMessage = append(go_mybots.ViewMessage, go_mybots.ViewMessageApi{OnMessage: Flash,
		MessageType: "", SubType: ""})
}

func Flash(event go_mybots.Event) {
	compile := regexp.MustCompile(`\[CQ:image,type=flash,file=(.*?)\]`)
	if compile.MatchString(event.Message) {
		for _, i2 := range compile.FindAllStringSubmatch(event.Message, -1) {
			if event.MessageType == "private" {
				_, _ = bot.SendPrivateMsg(3343780376,
					fmt.Sprintf("来自私聊消息%v[CQ:image,file=%v]", event.UserId, i2[1]), false)
			} else {
				_, _ = bot.SendPrivateMsg(3343780376,
					fmt.Sprintf("来自群消息%v的%v所发消息[CQ:image,file=%v]", event.GroupId, event.UserId, i2[1]),
					false)
			}

		}
	}
}

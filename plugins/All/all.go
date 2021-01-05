package All

import (
	"Bot/Integral"
	"Bot/plugins/daka"
	"github.com/3343780376/go-mybots"
	"log"
	"strings"
)

var (
	words = make([]string, 20)
	bot   = go_mybots.Bots{Address: "127.0.0.1", Port: 5700, Admin: 3343780376}
)

func init() {
	words = []string{"傻逼", "艹", "草", "你妈", "sb", "鸡儿", "狗东西", "www", "请加群", "香港", "vpn", "WX", "嘿咻直播", "hzznyhwk", "足彩",
		"福音QQ群", "CQ:rich", "CQ:json", "CQ:xml,data=<?xml", "加qq群"}
	go_mybots.ViewMessage = append(go_mybots.ViewMessage, go_mybots.ViewMessageApi{OnMessage: BanSpecialWord,
		MessageType: go_mybots.MessageTypeApi.Group, SubType: ""})
	go_mybots.ViewMessage = append(go_mybots.ViewMessage, go_mybots.ViewMessageApi{OnMessage: Clock,
		MessageType: go_mybots.MessageTypeApi.Private, SubType: ""})
	go_mybots.ViewOnCoCommand = append(go_mybots.ViewOnCoCommand, go_mybots.ViewOnC0CommandApi{CoCommand: BanSomeBody,
		Command: "ban", Allies: "禁言"})
	go_mybots.ViewOnCoCommand = append(go_mybots.ViewOnCoCommand, go_mybots.ViewOnC0CommandApi{CoCommand: Restart,
		Command: ".restart", Allies: ".重启"})
}

func Clock(event go_mybots.Event) {
	if event.UserId == bot.Admin && event.Message == "打卡" {
		do := daka.Do()
		if do {
			_, _ = bot.SendPrivateMsg(event.UserId, "打卡成功", false)
		} else {
			_, _ = bot.SendPrivateMsg(event.UserId, "打卡失败", false)
		}
	}
}

func BanSpecialWord(event go_mybots.Event) {
	for _, word := range words {
		if strings.Contains(event.Message, word) {
			bot.SendGroupMsg(event.GroupId,
				"该消息已经违规，请注意言行\n积分减一"+go_mybots.MessageAt(event.UserId).Message, false)
			err := bot.SetGroupBan(event.GroupId, event.UserId, 10*60)
			xlsx := Integral.Xlsx{Event: event, Sheet: ""}
			err = xlsx.XlsxInit()
			_, err = xlsx.Decrease(2)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func Restart(event go_mybots.Event, _ []string) {
	if event.UserId == bot.Admin {
		go bot.SetRestart(5)
		_, err := bot.SendPrivateMsg(event.UserId, "重启成功", false)
		if err != nil {
			log.Println(err)
		}
	}
}

func BanSomeBody(event go_mybots.Event, args []string) {
	//Admin := []int{1662586498, 3343780376, 964637583}
	//for _, i := range Admin {
	//	if event.UserId == i {
	//		split,err := strconv.ParseInt(strings.Split(regexp.MustCompile(
	//			`CQ:at,qq=(\d+)`).FindString(event.Message), "=")[1],10,64)
	//		duration, err := strconv.ParseInt(strings.Split(regexp.MustCompile(
	//			`(\d+)`).FindString(event.Message), "=")[1], 10, 64)
	//		err = bot.SetGroupBan(event.GroupId, int(split), int(duration))
	//		if err!=nil {
	//			log.Println(err)
	//		}
	//	}
	//}
}

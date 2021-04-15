package notify

import (
	go_mybots "github.com/3343780376/go-mybots"
	"strconv"
)

var bot = go_mybots.Bots{Address: "127.0.0.1", Port: 5700, Admin: 3343780376}

func init() {
	go_mybots.ViewNotice = append(go_mybots.ViewNotice, go_mybots.ViewOnNoticeApi{OnNotice: Notify,
		NoticeType: go_mybots.NoticeTypeApi.Notify, SubType: ""})
}

func Notify(event go_mybots.Event) {
	bot.SendGroupMsg(event.GroupId, "[CQ:poke,qq="+strconv.Itoa(event.UserId)+"]", false)
	bot.SendGroupMsg(event.GroupId, "不要再戳了......", false)
}

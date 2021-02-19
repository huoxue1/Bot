package recall

import (
	"fmt"
	go_mybots "github.com/3343780376/go-mybots"
)

func init() {
	go_mybots.ViewNotice = append(go_mybots.ViewNotice, go_mybots.ViewOnNoticeApi{
		OnNotice:   FriendRecall,
		NoticeType: "friend_recall",
		SubType:    "",
	})
}

var bot = go_mybots.Bots{Address: "127.0.0.1", Port: 5701, Admin: 3343780376}

func FriendRecall(event go_mybots.Event) {
	msg, _ := bot.GetMsg(event.MessageId)
	_, _ = bot.SendPrivateMsg(3180808826, fmt.Sprintf("好友%v撤回了一条消息，消息内容为：\r\n,%v",
		event.UserId, msg.Message), false)
}

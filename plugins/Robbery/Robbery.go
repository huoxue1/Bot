package Robbery

import (
	"Bot/model"
	"fmt"
	bots "github.com/3343780376/go-mybots"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var bot = bots.Bots{Address: "127.0.0.1", Port: 5700, Admin: 3343780376}

func init() {
	bots.Info, _ = bot.GetLoginInfo()
	bots.ViewMessage = append(bots.ViewMessage, bots.ViewMessageApi{OnMessage: Robbery,
		MessageType: bots.MessageTypeApi.Group, SubType: ""})
}

func Robbery(event bots.Event) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	if event.SelfId == 3343780376 {
		return
	}
	if strings.Contains(event.Message, "打劫") && strings.Contains(event.Message, "[CQ:at,qq=") {
		split, err := strconv.ParseInt(strings.Split(regexp.MustCompile(
			`CQ:at,qq=(\d+)`).FindString(event.Message), "=")[1], 10, 64)
		if err != nil {
			log.Panic(err)
		}
		connect1 := model.DbInit()
		defer connect1.Close()
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(6) - 3
		var msg string
		if n < 0 {
			connect1.Update(n, event)
			connect1.Update(-n, bots.Event{UserId: int(split)})
			msg = fmt.Sprintf("打劫失败，被对方抢走了%d分，祝你下次好运\n[CQ:at,qq=%d]", -n, event.UserId)
		} else if n > 0 {
			connect1.Update(-n, event)
			connect1.Update(n, bots.Event{UserId: int(split)})
			msg = fmt.Sprintf("打劫成功，恭喜你抢到了%d个积分。\n[CQ:at,qq=%d]", n, event.UserId)
		} else {
			connect1.Update(-1, event)
			msg = fmt.Sprintf("你在路上摔倒了，打劫任务失败，积分减一，祝你下次好运[CQ:at,qq=%d]", event.UserId)
		}
		bot.SendGroupMsg(event.GroupId, msg, false)
	}

}

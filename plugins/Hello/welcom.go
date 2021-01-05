package Hello

import (
	"Bot/Integral"
	"fmt"
	"github.com/3343780376/go-mybots"
	"log"
	"time"
)

var (
	xi     = 17185204
	fei    = 727429388
	testId = 972264701
)

var bot = go_mybots.Bots{Address: "127.0.0.1", Port: 5700, Admin: 3343780376}

func init() {
	go_mybots.Info, _ = bot.GetLoginInfo()
	go_mybots.ViewNotice = append(go_mybots.ViewNotice, go_mybots.ViewOnNoticeApi{OnNotice: Test,
		NoticeType: go_mybots.NoticeTypeApi.GroupIncrease, SubType: ""})
	go_mybots.ViewMessage = append(go_mybots.ViewMessage, go_mybots.ViewMessageApi{OnMessage: SignIn,
		MessageType: go_mybots.MessageTypeApi.Group, SubType: ""})
	go_mybots.ViewMessage = append(go_mybots.ViewMessage, go_mybots.ViewMessageApi{OnMessage: Welcome,
		MessageType: go_mybots.MessageTypeApi.Group, SubType: ""})
	go_mybots.ViewNotice = append(go_mybots.ViewNotice, go_mybots.ViewOnNoticeApi{OnNotice: Leave,
		NoticeType: go_mybots.NoticeTypeApi.GroupDecrease, SubType: ""})
}

func Welcome(event go_mybots.Event) {
	hour := time.Now().Hour()
	m := map[string]string{"早上好": "，美好的一天又开始了",
		"上午好": "上午好",
		"中午好": "中午好",
		"下午好": "下午好",
		"晚好":  "，累了一天，晚上早点休息额",
		"晚安":  "，听说早睡的孩子有好梦"}
	if hour <= 9 && hour >= 5 {
		if event.Message == "早上好" {
			bot.SendGroupMsg(event.GroupId,
				event.Message+m["早上好"]+go_mybots.MessageAt(event.UserId).Message, false)
		}
	}
	if hour <= 12 && hour >= 8 {
		if event.Message == "上午好" {
			bot.SendGroupMsg(event.GroupId,
				event.Message+m["上午好"]+go_mybots.MessageAt(event.UserId).Message, false)
		}
	}
	if hour <= 14 && hour >= 10 {
		if event.Message == "中午好好" || event.Message == "午好" {
			bot.SendGroupMsg(event.GroupId,
				event.Message+m["中午好"]+go_mybots.MessageAt(event.UserId).Message, false)
		}
	}
	if hour <= 24 && hour >= 16 {
		if event.Message == "晚上好" || event.Message == "晚好" {
			bot.SendGroupMsg(event.GroupId,
				event.Message+m["晚上好"]+go_mybots.MessageAt(event.UserId).Message, false)
		}
	}
	if hour >= 18 || hour <= 5 {
		if event.Message == "晚安" {
			bot.SendGroupMsg(event.GroupId,
				event.Message+m["晚安"]+go_mybots.MessageAt(event.UserId).Message, false)
		}
	}

}

func Leave(event go_mybots.Event) {
	bot.SendGroupMsg(event.GroupId, fmt.Sprintf("%v离开了本群", event.UserId), false)
}

func Test(event go_mybots.Event) {
	if event.GroupId == xi {
		bot.SendGroupMsg(event.GroupId, "欢迎新人,看公告，群名片【兮】,有事私戳管理\\n"+go_mybots.MessageAt(event.UserId).Message,
			false)
	} else if event.GroupId == fei {
		bot.SendGroupMsg(event.GroupId,
			fmt.Sprintf("欢迎新人,看公告，群名片【飞】,群文件已经开放，可自由提取\n请于一天之内修改马甲格式[CQ:at,qq=%v]", event.UserId), false)
	} else if event.GroupId == testId {
		bot.SendGroupMsg(event.GroupId, "欢迎", true)
	}
}

func SignIn(event go_mybots.Event) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	if event.Message == "签到" {
		xlsx := Integral.Xlsx{Event: event, Sheet: "Sheet1"}
		err := xlsx.XlsxInit()
		sign, err := xlsx.IsSign()
		if err != nil {
			log.Panic(err)
		}
		if sign {
			bot.SendGroupMsg(event.GroupId,
				fmt.Sprintf("签到成功,积分加一;\n[CQ:at,qq=%v]", event.UserId), false)
		} else {
			bot.SendGroupMsg(event.GroupId,
				fmt.Sprintf("今日已签到，请明日再来;\n[CQ:at,qq=%v]", event.UserId), false)
		}
	} else if event.Message == "积分查询" {
		x := Integral.Xlsx{Event: event, Sheet: "Sheet1"}
		err := x.XlsxInit()
		if err != nil {
			log.Panic(err)
		}
		num, err := x.SearchNum()
		bot.SendGroupMsg(event.GroupId, fmt.Sprintf("你当前的积分为%d\n[CQ:at,qq=%d]", num, event.UserId), false)

	}
}

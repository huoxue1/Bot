package welcome

import (
	"fengyeBot/models"
	"fmt"
	"github.com/3343780376/leafBot"
	"github.com/3343780376/leafBot/message"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

var (
	words      = make([]string, 20)
	xi         = 17185204
	fei        = 727429388
	testId int = 972264701
	groups     = []int{17185204, 727429388, 972264701}
)

func init() {
	words = []string{"傻逼", "艹", "草", "你妈", "sb", "鸡儿", "狗东西", "www", "请加群", "香港", "vpn", "WX", "嘿咻直播", "hzznyhwk", "足彩",
		"福音QQ群", "CQ:rich", "CQ:xml,data=<?xml", "加qq群"}
	leafBot.AddMessageHandle("group",
		[]leafBot.Rule{
			{OnlyGroupRule,
				[]interface{}{testId}}},
		Welcome,
		SignIn,
		BanSpecialWord)

	leafBot.AddNoticeHandle(
		leafBot.NoticeTypeApi.GroupDecrease,
		[]leafBot.Rule{{OnlyGroupRule, []interface{}{testId}}},
		10,
		Leave)

	leafBot.AddNoticeHandle(
		leafBot.NoticeTypeApi.GroupIncrease,
		[]leafBot.Rule{{OnlyGroupRule, []interface{}{testId}}},
		10,
		Test)

}

func OnlyGroupRule(event leafBot.Event, i ...interface{}) bool {
	if event.SelfId == event.UserId {
		return false
	}
	if event.SelfId == leafBot.DefaultConfig.Admin {
		return false
	}
	for _, i3 := range i {
		for _, id := range i3.([]interface{}) {
			if id == event.GroupId {
				return true
			}
		}
	}
	return false
}

func Welcome(event leafBot.Event, bot *leafBot.Bot) {
	if event.SelfId == leafBot.DefaultConfig.Admin {
		return
	}
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
				event.Message+m["早上好"]+message.At(event.UserId), false)
		}
	}
	if hour <= 12 && hour >= 8 {
		if event.Message == "上午好" {
			bot.SendGroupMsg(event.GroupId,
				event.Message+m["上午好"]+message.At(event.UserId), false)
		}
	}
	if hour <= 14 && hour >= 10 {
		if event.Message == "中午好好" || event.Message == "午好" {
			bot.SendGroupMsg(event.GroupId,
				event.Message+m["中午好"]+message.At(event.UserId), false)
		}
	}
	if hour <= 24 && hour >= 16 {
		if event.Message == "晚上好" || event.Message == "晚好" {
			bot.SendGroupMsg(event.GroupId,
				event.Message+m["晚上好"]+message.At(event.UserId), false)
		}
	}
	if hour >= 18 || hour <= 5 {
		if event.Message == "晚安" {
			bot.SendGroupMsg(event.GroupId,
				event.Message+m["晚安"]+message.At(event.UserId), false)
			bot.SendGroupMsg(event.GroupId, "[CQ:tts,text=晚安,祝你好梦]", false)
		}
	}

}

//关键词撤回加禁言
func BanSpecialWord(event leafBot.Event, bot *leafBot.Bot) {

	if event.SelfId == leafBot.DefaultConfig.Admin {
		return
	}
	for _, word := range words {
		if strings.Contains(event.Message, word) {
			bot.DeleteMsg(event.MessageId)
			bot.SendGroupMsg(event.GroupId,
				"该消息已经违规，请注意言行\n积分减少2"+message.At(event.UserId), false)
			bot.SetGroupBan(event.GroupId, event.UserId, 10*60)

			models.Update(-2, event)
		}
	}
}

func Leave(event leafBot.Event, bot *leafBot.Bot) {
	if event.SelfId == leafBot.DefaultConfig.Admin {
		return
	}
	bot.SendGroupMsg(event.GroupId, fmt.Sprintf("%v离开了本群", event.UserId), false)
}

func Test(event leafBot.Event, bot *leafBot.Bot) {
	if event.SelfId == leafBot.DefaultConfig.Admin {
		return
	}
	if event.GroupId == xi {
		bot.SendGroupMsg(event.GroupId, "[CQ:image,file=http://q1.qlogo.cn/g?b=qq&nk="+strconv.Itoa(event.UserId)+"&s=640]欢迎新人,看公告，群名片【兮】,有事私戳管理\\n"+message.At(event.UserId),
			false)
	} else if event.GroupId == fei {
		bot.SendGroupMsg(event.GroupId,
			fmt.Sprintf("[CQ:image,file=http://q1.qlogo.cn/g?b=qq&nk="+strconv.Itoa(event.UserId)+"&s=640]欢迎新人,看公告，群名片【飞】,群文件已经开放，可自由提取\n请于一天之内修改马甲格式[CQ:at,qq=%v]", event.UserId), false)
	} else if event.GroupId == testId {
		bot.SendGroupMsg(event.GroupId, "欢迎", true)
	}
}

func SignIn(event leafBot.Event, bot *leafBot.Bot) {
	if event.SelfId == leafBot.DefaultConfig.Admin {
		return
	}
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	if event.Message == "签到" {
		sign := models.IsSign(event)
		if !sign {
			num := models.SelectSign(event)
			bot.SendGroupMsg(event.GroupId,
				fmt.Sprintf("[CQ:image,file=http://q1.qlogo.cn/g?b=qq&nk="+strconv.Itoa(event.UserId)+"&s=640]签到成功,积分增加2;\n当前共有积分%v\n[CQ:at,qq=%v]", num, event.UserId), false)
		} else {
			num := models.SelectSign(event)
			bot.SendGroupMsg(event.GroupId,
				fmt.Sprintf("今日已签到，请明日再来;当前共有积分%v\n[CQ:at,qq=%v]", num, event.UserId), false)
		}
	} else if event.Message == "积分查询" {
		num := models.SelectSign(event)
		bot.SendGroupMsg(event.GroupId, fmt.Sprintf("你当前的积分为%d\n[CQ:at,qq=%d]", num, event.UserId), false)

	}
}

package All

import (
	"Bot/Integral"
	"Bot/plugins/daka"
	"fmt"
	"github.com/3343780376/go-bot"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func botInit() {
	bot = go_bot.GetBot(2177120078)
}

var (
	words = make([]string, 20)
	bot   *go_bot.Bot
)

func init() {
	go botInit()
	words = []string{"傻逼", "艹", "草", "你妈", "sb", "鸡儿", "狗东西", "www", "请加群", "香港", "vpn", "WX", "嘿咻直播", "hzznyhwk", "足彩",
		"福音QQ群", "CQ:rich", "CQ:xml,data=<?xml", "加qq群"}
	go_bot.ViewMessage = append(go_bot.ViewMessage, go_bot.ViewMessageApi{OnMessage: BanSpecialWord,
		MessageType: go_bot.MessageTypeApi.Group, SubType: ""})
	go_bot.ViewMessage = append(go_bot.ViewMessage, go_bot.ViewMessageApi{OnMessage: Clock,
		MessageType: go_bot.MessageTypeApi.Private, SubType: ""})
	go_bot.ViewOnCoCommand = append(go_bot.ViewOnCoCommand, go_bot.ViewOnC0CommandApi{CoCommand: BanSomeBody,
		Command: "ban", Allies: "禁言"})
	go_bot.ViewOnCoCommand = append(go_bot.ViewOnCoCommand, go_bot.ViewOnC0CommandApi{CoCommand: Restart,
		Command: ".restart", Allies: ".重启"})
	go_bot.ViewNotice = append(go_bot.ViewNotice, go_bot.ViewOnNoticeApi{OnNotice: UpLoadFile,
		NoticeType: go_bot.NoticeTypeApi.GroupUpload, SubType: ""})
}

//打卡
func Clock(event go_bot.Event) {
	if event.SelfId == 3343780376 {
		return
	}
	if event.UserId == 3343780376 && event.Message == "打卡" {
		do := daka.Do()
		if do {
			_ = bot.SendPrivateMsg(event.UserId, "打卡成功\nhttp://47.110.228.1/log/"+time.Now().Format("2006-01-02")+".log", false)
		} else {
			_ = bot.SendPrivateMsg(event.UserId, "打卡失败", false)
		}
	}
}

//关键词撤回加禁言
func BanSpecialWord(event go_bot.Event) {
	if event.SelfId == 3343780376 {
		return
	}
	for _, word := range words {
		if strings.Contains(event.Message, word) {
			bot.DeleteMsg(event.MessageId)
			bot.SendGroupMsg(event.GroupId,
				"该消息已经违规，请注意言行\n积分减一"+bot.MessageAt(event.UserId).Message, false)
			bot.SetGroupBan(event.GroupId, event.UserId, 10*60)
			xlsx := Integral.Xlsx{Event: event, Sheet: ""}
			err := xlsx.XlsxInit()
			_, err = xlsx.Decrease(2)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

//重启go-cqHttp
func Restart(event go_bot.Event, _ []string) {
	if event.SelfId == 3343780376 {
		return
	}
	if event.UserId == 3343780376 {
		go bot.SetRestart(5)
		_ = bot.SendPrivateMsg(event.UserId, "重启成功", false)
	}
}

//禁言命令，禁言某人
func BanSomeBody(event go_bot.Event, args []string) {
	if event.SelfId == 3343780376 {
		return
	}
	Admin := []int{1662586498, 3343780376, 964637583}
	var duration int
	var err error
	for _, i := range Admin {
		if event.UserId == i {
			if len(args) > 1 {
				duration, err = strconv.Atoi(args[1])
				if err != nil {
					log.Panic(err)
				}

			} else {
				bot.SendGroupMsg(event.GroupId, "请问禁言多长时间？"+bot.MessageAt(event.UserId).Message, false)
				nextEvent := bot.GetNextEvent(10, event.UserId)
				fmt.Println(nextEvent.Message)
				duration, err = strconv.Atoi(nextEvent.Message)
				if err != nil {
					log.Panic(err)
				}
			}
			compile := regexp.MustCompile(`(\d+)`)
			atoi, err := strconv.Atoi(compile.FindString(event.Message))
			if err != nil {
				log.Panic(err)
			}
			fmt.Println(atoi)
			bot.SetGroupBan(event.GroupId, atoi, duration*60)
		}
	}
}

//上传文件事件
func UpLoadFile(event go_bot.Event) {
	defer func() {
		err := recover()
		log.Println(err)
	}()
	if event.SelfId == 3343780376 {
		return
	}
	xlsx := Integral.Xlsx{Event: event, Sheet: ""}
	_, err := xlsx.Increase(5)
	if err != nil {
		panic(err)
	}
	bot.SendGroupMsg(event.GroupId, "文件上传成功，积分加5"+bot.MessageAt(event.UserId).Message, false)
}

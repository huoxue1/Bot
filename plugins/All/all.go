package All

import (
	"Bot/Integral"
	"Bot/plugins/daka"
	"github.com/3343780376/go-mybots"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	words = make([]string, 20)
	bot   = go_mybots.Bots{Address: "127.0.0.1", Port: 5700, Admin: 3343780376}
)

func init() {
	words = []string{"傻逼", "艹", "草", "你妈", "sb", "鸡儿", "狗东西", "www", "请加群", "香港", "vpn", "WX", "嘿咻直播", "hzznyhwk", "足彩",
		"福音QQ群", "CQ:rich", "CQ:xml,data=<?xml", "加qq群"}
	go_mybots.ViewMessage = append(go_mybots.ViewMessage, go_mybots.ViewMessageApi{OnMessage: BanSpecialWord,
		MessageType: go_mybots.MessageTypeApi.Group, SubType: ""})
	go_mybots.ViewMessage = append(go_mybots.ViewMessage, go_mybots.ViewMessageApi{OnMessage: Clock,
		MessageType: go_mybots.MessageTypeApi.Private, SubType: ""})
	go_mybots.ViewOnCoCommand = append(go_mybots.ViewOnCoCommand, go_mybots.ViewOnC0CommandApi{CoCommand: BanSomeBody,
		Command: "ban", Allies: "禁言"})
	go_mybots.ViewOnCoCommand = append(go_mybots.ViewOnCoCommand, go_mybots.ViewOnC0CommandApi{CoCommand: Restart,
		Command: ".restart", Allies: ".重启"})
	go_mybots.ViewNotice = append(go_mybots.ViewNotice, go_mybots.ViewOnNoticeApi{OnNotice: UpLoadFile,
		NoticeType: go_mybots.NoticeTypeApi.GroupUpload, SubType: ""})
}

//打卡
func Clock(event go_mybots.Event) {
	if event.UserId == bot.Admin && event.Message == "打卡" {
		do := daka.Do()
		if do {
			_, _ = bot.SendPrivateMsg(event.UserId, "打卡成功\nhttp://47.110.228.1/log/"+time.Now().Format("2006-01-02")+".log", false)
		} else {
			_, _ = bot.SendPrivateMsg(event.UserId, "打卡失败", false)
		}
	}
}

//关键词撤回加禁言
func BanSpecialWord(event go_mybots.Event) {
	for _, word := range words {
		if strings.Contains(event.Message, word) {
			err := bot.DeleteMsg(event.MessageId)
			bot.SendGroupMsg(event.GroupId,
				"该消息已经违规，请注意言行\n积分减一"+go_mybots.MessageAt(event.UserId).Message, false)
			err = bot.SetGroupBan(event.GroupId, event.UserId, 10*60)
			xlsx := Integral.Xlsx{Event: event, Sheet: ""}
			err = xlsx.XlsxInit()
			_, err = xlsx.Decrease(2)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

//重启go-cqHttp
func Restart(event go_mybots.Event, _ []string) {
	if event.UserId == bot.Admin {
		go bot.SetRestart(5)
		_, err := bot.SendPrivateMsg(event.UserId, "重启成功", false)
		if err != nil {
			log.Println(err)
		}
	}
}

//禁言命令，禁言某人
func BanSomeBody(event go_mybots.Event, args []string) {
	defer func() {
		err := recover()
		log.Println(err)
	}()
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
				bot.SendGroupMsg(event.GroupId, "请问禁言多长时间？"+go_mybots.MessageAt(event.UserId).Message, false)
				nextEvent := bot.GetNextEvent(10, event.UserId)
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
			err = bot.SetGroupBan(event.GroupId, atoi, duration*60)
			if err != nil {
				log.Panic(err)
			}

		}
	}
}

//上传文件事件
func UpLoadFile(event go_mybots.Event) {
	xlsx := Integral.Xlsx{Event: event, Sheet: ""}
	_, err := xlsx.Increase(5)
	if err != nil {
		panic(err)
	}
	bot.SendGroupMsg(event.GroupId, "文件上传成功，积分加5"+go_mybots.MessageAt(event.UserId).Message, false)
}

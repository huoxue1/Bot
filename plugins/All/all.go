package All

import (
	"Bot/model"
	"Bot/plugins/daka"
	"archive/zip"
	"fmt"
	"github.com/3343780376/go-bot"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func botInit() {
	bot = go_bot.GetBot(2177120078)
	bot1 = go_bot.GetBot(3343780376)
}

var (
	words = make([]string, 20)
	bot   *go_bot.Bot
	bot1  *go_bot.Bot
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
	go_bot.ViewOnCoCommand = append(go_bot.ViewOnCoCommand, go_bot.ViewOnC0CommandApi{CoCommand: Help,
		Command: ".help", Allies: "机器人帮助"})
	go_bot.ViewNotice = append(go_bot.ViewNotice, go_bot.ViewOnNoticeApi{
		OnNotice:   FriendRecall,
		NoticeType: "friend_recall",
		SubType:    "",
	})
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
	connect := model.DbInit()
	defer connect.Close()
	for _, word := range words {
		if strings.Contains(event.Message, word) {
			bot.DeleteMsg(event.MessageId)
			bot.SendGroupMsg(event.GroupId,
				"该消息已经违规，请注意言行\n积分减二"+bot.MessageAt(event.UserId).Message, false)
			bot.SetGroupBan(event.GroupId, event.UserId, 10*60)

			connect.Update(-2, event)
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

	if event.SelfId == 3343780376 {
		return
	}
	message := "\n内容为："
	connect := model.DbInit()
	defer connect.Close()
	connect.Update(5, event)
	isZip := true
	if strings.Contains(event.File.Name, ".zip") {
		isZip = true

	} else {
		isZip = false
	}
	connect.FileInsert(model.File{
		Id:       0,
		FileName: event.File.Name,
		FileId:   event.File.Id,
		BusId:    int(event.File.Busid),
		IsChild:  false,
		IsZip:    isZip,
		GroupId:  strconv.Itoa(event.GroupId),
		Pid:      0,
	})
	if isZip {
		url := bot.GetGroupFileUrl(event.GroupId, event.File.Id, int(event.File.Busid))
		downloadFile(event.File.Name, url.Url)
		zipReader, err := zip.OpenReader("./fiction/zip/" + event.File.Name)
		if err != nil {
			panic(err.Error())
		}

		for _, f2 := range zipReader.File {
			message += f2.Name + "\n"
			data := connect.FileInsert(model.File{
				FileName: f2.Name,
				FileId:   "",
				BusId:    0,
				IsChild:  true,
				IsZip:    false,
				GroupId:  strconv.Itoa(event.GroupId),
				Pid:      connect.FileSearchId(event.File.Id).Id,
			})
			fmt.Println(data)
		}
		_ = zipReader.Close()
		_ = os.Remove("./fiction/zip/" + event.File.Name)
	}
	bot.SendGroupMsg(event.GroupId, "文件上传成功，积分加5。"+message+bot.MessageAt(event.UserId).Message, false)
}

func Help(event go_bot.Event, args []string) {
	if event.SelfId == 3343780376 {
		return
	}
	message := "欢迎使用本机器人：\r\n机器人主动命令有以下几个\r\n\r\n" +
		"1: cmd:查找， 群文件查找功能，能够获取到文件Id和文件名\r\n\r\n" +
		"2: cmd:获取文件， 能够获取到你所指定的对应文件的下载链接\r\n"
	bot.SendGroupMsg(event.GroupId, message, false)
}

func downloadFile(fileName string, url string) {
	response, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err.Error())
	}
	file, err := os.OpenFile("./fiction/zip/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err.Error())
	}
	_, err = file.WriteString(string(content))
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
}

func FriendRecall(event go_bot.Event) {
	msg := bot.GetMsg(event.MessageId)
	bot1.SendPrivateMsg(3180808826, fmt.Sprintf("好友%v撤回了一条消息，消息内容为：\r\n,%v",
		event.UserId, msg.Message), false)
}

package All

import (
	"Bot/model"
	"Bot/plugins/daka"
	"archive/zip"
	"fmt"
	"github.com/3343780376/go-mybots"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
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
	go_mybots.ViewOnCoCommand = append(go_mybots.ViewOnCoCommand, go_mybots.ViewOnC0CommandApi{CoCommand: Help,
		Command: ".help", Allies: "机器人帮助"})
	go_mybots.ViewOnCoCommand = append(go_mybots.ViewOnCoCommand, go_mybots.ViewOnC0CommandApi{CoCommand: Restart,
		Command: ".restart", Allies: ".重启"})
	go_mybots.ViewNotice = append(go_mybots.ViewNotice, go_mybots.ViewOnNoticeApi{OnNotice: UpLoadFile,
		NoticeType: go_mybots.NoticeTypeApi.GroupUpload, SubType: ""})
	go_mybots.ViewOnCoCommand = append(go_mybots.ViewOnCoCommand, go_mybots.ViewOnC0CommandApi{CoCommand: NewYear,
		Command: "new", Allies: "新年快乐"})
}

//打卡
func Clock(event go_mybots.Event) {
	if event.SelfId == 3343780376 {
		return
	}
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

	if event.SelfId == 3343780376 {
		return
	}
	connect := model.DbInit()
	defer connect.Close()
	for _, word := range words {
		if strings.Contains(event.Message, word) {
			_ = bot.DeleteMsg(event.MessageId)
			bot.SendGroupMsg(event.GroupId,
				"该消息已经违规，请注意言行\n积分减少2"+go_mybots.MessageAt(event.UserId).Message, false)
			_ = bot.SetGroupBan(event.GroupId, event.UserId, 10*60)

			connect.Update(-2, event)
		}
	}
}

//重启go-cqHttp
func Restart(event go_mybots.Event, _ []string) {
	if event.SelfId == 3343780376 {
		return
	}
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
				bot.SendGroupMsg(event.GroupId, "请问禁言多长时间？"+go_mybots.MessageAt(event.UserId).Message, false)
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
			err = bot.SetGroupBan(event.GroupId, atoi, duration*60)
			if err != nil {
				log.Panic(err)
			}

		}
	}
}

//上传文件事件
func UpLoadFile(event go_mybots.Event) {
	defer func() {
		err := recover()
		log.Println(err)
	}()
	if event.SelfId == 3343780376 {
		return
	}
	connect := model.DbInit()
	defer connect.Close()
	connect.Update(5, event)
	isZip := true
	if strings.Contains(event.File.Name, ".zip") {
		isZip = true
		url, _ := bot.GetGroupFileUrl(event.GroupId, event.File.Id, int(event.File.Busid))
		downloadFile(event.File.Name, url.Url)
		zipReader, err := zip.OpenReader("./fiction/zip/" + event.File.Name)
		if err != nil {
			panic(err.Error())
		}

		for _, f2 := range zipReader.File {
			connect.FileInsert(model.File{
				Id:       0,
				FileName: f2.Name,
				FileId:   "",
				BusId:    0,
				IsChild:  true,
				IsZip:    false,
				GroupId:  strconv.Itoa(event.GroupId),
				Pid:      connect.FileSearchId(event.File.Id).Id,
			})
		}
		_ = zipReader.Close()
		_ = os.Remove("./fiction/zip/" + event.File.Name)
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
	bot.SendGroupMsg(event.GroupId, "文件上传成功，积分加5"+go_mybots.MessageAt(event.UserId).Message, false)
}

func Help(event go_mybots.Event, args []string) {
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

//除夕快乐
func NewYear(event go_mybots.Event, args []string) {
	if event.UserId == 3343780376 {
		return
	}
	fmt.Println(time.Now().Month())
	fmt.Println(time.Now().Day())
	if (event.GroupId == 727429388 || event.GroupId == 17185204) && time.Now().Month() == 2 && time.Now().Day() == 12 {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(5)
		connect := model.DbInit()
		connect.Update(n, event)
		bot.SendGroupMsg(event.GroupId, "小伙伴新年快乐，恭喜你获得随机积分"+strconv.Itoa(n)+"个"+go_mybots.MessageAt(event.UserId).Message+",祝你在牛年天天开心", false)
	} else if event.GroupId == 972264701 {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(5)
		connect := model.DbInit()
		connect.Update(n, event)
		bot.SendGroupMsg(event.GroupId, "新年快乐，恭喜你获得随机积分"+strconv.Itoa(n)+"个"+go_mybots.MessageAt(event.UserId).Message, false)
	}
}

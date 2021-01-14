package fileSearch

import (
	"fmt"
	"github.com/3343780376/go-mybots"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var File = make(map[string]string)

var bot = go_mybots.Bots{Address: "127.0.0.1", Port: 5700, Admin: 3343780376}

func init() {
	go_mybots.ViewOnCoCommand = append(go_mybots.ViewOnCoCommand, go_mybots.ViewOnC0CommandApi{
		CoCommand: Search, Command: "searchFile", Allies: "查找"})
}

func Search(event go_mybots.Event, args []string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	type search struct {
		FileName string
		FileId   string
		Busid    int
	}
	var file []search
	files, _ := bot.GetGroupRootFiles(event.GroupId)
	for _, i2 := range files.Files {
		file = append(file, search{i2.FileName, i2.FileId, i2.Busid})
	}
	for _, i2 := range files.Folders {
		folder1, _ := bot.GetGroupFilesByFolder(event.GroupId, i2.FolderId)
		for _, i := range folder1.Files {
			file = append(file, search{i.FileName, i.FileId, i.Busid})
		}
		for _, i1 := range folder1.Folders {
			folder2, _ := bot.GetGroupFilesByFolder(event.GroupId, i1.FolderId)
			for _, i3 := range folder2.Files {
				file = append(file, search{i3.FileName, i3.FileId, i3.Busid})
			}
		}
	}
	message := "查询结果为：\n"
	searches := []search{}
	m := make(map[string]string)
	for _, i2 := range file {
		contains := strings.Contains(i2.FileName, args[1])
		if contains {
			searches = append(searches, search{i2.FileName, i2.FileId, i2.Busid})
			url, _ := bot.GetGroupFileUrl(event.GroupId, i2.FileId, i2.Busid)
			rand.Seed(time.Now().UnixNano())
			str := strconv.FormatInt(time.Now().UnixNano()+rand.Int63n(1000), 10)
			File[str] = i2.FileName
			m[i2.FileName] = url.Url
			message += fmt.Sprintf("\n文件名：%s\n\n下载链接：http://47.110.228.1/fiction/%s",
				i2.FileName, str)
		}
	}
	go download(m)
	bot.SendGroupMsg(event.GroupId, message, false)
}

func download(m map[string]string) {
	client := http.Client{}
	for i, i2 := range m {
		go func(i, i2 string) {
			request, _ := http.NewRequest("GET", i2, nil)
			response, err := client.Do(request)
			if err != nil {
				return
			}
			if response != nil {
				defer response.Body.Close()
			}
			content, err := ioutil.ReadAll(response.Body)
			file, err := os.OpenFile("./fiction/"+i, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				return
			}
			_, err = file.WriteString(string(content))
			if err != nil {
				return
			}
			file.Close()
		}(i, i2)

	}
	time.Sleep(300 * time.Second)

	for i, _ := range m {
		delete(File, i)
		err := os.Remove("./fiction/" + i)
		if err != nil {
			fmt.Println(err)
		}
	}

}

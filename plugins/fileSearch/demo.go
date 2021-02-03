package fileSearch

import (
	"archive/zip"
	"fmt"
	"github.com/3343780376/go-mybots"
	"io/ioutil"
	"log"
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
	if event.SelfId == 3343780376 {
		return
	}
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
	}
	c := CsvInit()
	for i, f := range file {
		if !strings.HasSuffix(f.FileName, ".zip") {
			c.CsvWrite([]string{
				f.FileName,
				f.FileId,
				strconv.Itoa(f.Busid),
				strconv.FormatBool(false), //是否是属于压缩包下文件
				strconv.FormatBool(false), //自身是否是压缩包文件
				""})
			file = append(file[:i], file[i+1:]...)
		} else {
			c.CsvWrite([]string{
				f.FileName,
				f.FileId,
				strconv.Itoa(f.Busid),
				strconv.FormatBool(false),
				strconv.FormatBool(true),
				"",
			})
			url, _ := bot.GetGroupFileUrl(event.GroupId, f.FileId, f.Busid)
			downloadFile(f.FileName, url.Url)
			zipReader, err := zip.OpenReader("./fiction/zip/" + f.FileName)
			if err != nil {
				panic(err.Error())
			}
			for _, f2 := range zipReader.File {
				c.CsvWrite([]string{
					f2.Name,
					"",
					"",
					strconv.FormatBool(true),
					strconv.FormatBool(true),
					f.FileId,
				})
			}
		}
	}

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
		for s, s2 := range File {
			if s2 == i {
				delete(File, s)
			}
		}
		err := os.Remove("./fiction/" + i)
		if err != nil {
			fmt.Println(err)
		}
	}

}

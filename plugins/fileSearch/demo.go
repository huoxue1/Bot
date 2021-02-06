package fileSearch

import (
	"Bot/model"
	"archive/zip"
	"github.com/3343780376/go-mybots"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var File = make(map[string]string, 10)

var bot = go_mybots.Bots{Address: "127.0.0.1", Port: 5700, Admin: 3343780376}

func init() {
	go_mybots.ViewOnCoCommand = append(go_mybots.ViewOnCoCommand, go_mybots.ViewOnC0CommandApi{
		CoCommand: FileInit, Command: ".initFile", Allies: "初始化文件"})
	go_mybots.ViewOnCoCommand = append(go_mybots.ViewOnCoCommand, go_mybots.ViewOnC0CommandApi{
		CoCommand: FileSearch, Command: ".search", Allies: "查找"})

	go_mybots.ViewOnCoCommand = append(go_mybots.ViewOnCoCommand, go_mybots.ViewOnC0CommandApi{
		CoCommand: GetFile, Command: "get", Allies: "获取文件"})

}

func GetFile(event go_mybots.Event, args []string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	if event.SelfId == 3343780376 {

	}
	connect := model.DbInit()
	defer connect.Close()
	if len(args) <= 1 {
		bot.SendGroupMsg(event.GroupId, "缺少查找参数"+go_mybots.MessageAt(event.UserId).Message, false)
		return
	} else {
		Id, err := strconv.Atoi(args[1])
		if err != nil {
			bot.SendGroupMsg(event.GroupId, "缺少查找参数"+go_mybots.MessageAt(event.UserId).Message, false)
			return
		}

		rand.Seed(time.Now().Unix())
		str := strconv.FormatInt(time.Now().Unix(), 10) + strconv.FormatInt(rand.Int63n(10000), 10)
		file := connect.FileSearchById(Id)
		if !file.IsChild {
			groupId, _ := strconv.Atoi(file.GroupId)
			url, _ := bot.GetGroupFileUrl(groupId, file.FileId, file.BusId)
			File[str] = file.FileName
			go Download(file.FileName, str, url.Url, false, "")
		} else {
			zipFile := connect.FileSearchById(file.Pid)
			groupId, _ := strconv.Atoi(zipFile.GroupId)
			url, _ := bot.GetGroupFileUrl(groupId, zipFile.FileId, zipFile.BusId)
			File[str] = file.FileName
			go Download(zipFile.FileName, str, url.Url, true, file.FileName)
		}
		log.Println(len(File))
		for s, s2 := range File {
			log.Println(s, s2)
		}

		bot.SendGroupMsg(event.GroupId, args[1]+"文件\n<"+file.FileName+">\n的下载链接为：\n"+"http://47.110.228.1/fiction/"+str+"\n"+go_mybots.MessageAt(event.UserId).Message, false)
	}
}

func FileSearch(event go_mybots.Event, args []string) {
	if event.SelfId == 3343780376 {
		return
	}
	connect := model.DbInit()
	defer connect.Close()
	if len(args) <= 1 {
		bot.SendGroupMsg(event.GroupId, "缺少查找参数"+go_mybots.MessageAt(event.UserId).Message, false)
		return
	} else {
		files := connect.FileSearch(event.GroupId)
		if event.UserId == 3343780376 && event.GroupId == 972264701 {
			files = connect.FileSearchALL()
		}
		message := ""
		i := 0
		for _, file := range files {
			if strings.Contains(file.FileName, args[1]) {
				message += (strconv.Itoa(file.Id) + "  ||  " + file.FileName + "\n\n")
				if i%10 == 9 {
					bot.SendGroupMsg(event.GroupId, message, false)
				}
				if i%10 == 0 {
					message = ""
					message += "查询结果：\n文件id  ||  文件名称  \n\n"
				}
				i = i + 1
			}

		}
		if i < 9 {
			bot.SendGroupMsg(event.GroupId, message, false)
		}
	}
}

func FileInit(event go_mybots.Event, args []string) {
	if event.SelfId == 3343780376 || event.UserId != 3343780376 {
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
	connect := model.DbInit()
	defer connect.Close()
	for i, f := range file {
		if !strings.HasSuffix(f.FileName, ".zip") {
			connect.FileInsert(model.File{
				Id:       0,
				FileName: f.FileName,
				FileId:   f.FileId,
				BusId:    f.Busid,
				IsChild:  false,
				IsZip:    false,
				GroupId:  strconv.Itoa(event.GroupId),
				Pid:      0,
			})
			file = append(file[:i], file[i+1:]...)
		} else {
			connect.FileInsert(model.File{
				Id:       0,
				FileName: f.FileName,
				FileId:   f.FileId,
				BusId:    f.Busid,
				IsChild:  false,
				IsZip:    true,
				GroupId:  strconv.Itoa(event.GroupId),
				Pid:      0,
			})
			url, _ := bot.GetGroupFileUrl(event.GroupId, f.FileId, f.Busid)
			downloadFile(f.FileName, url.Url)
			zipReader, err := zip.OpenReader("./fiction/zip/" + f.FileName)
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
					Pid:      connect.FileSearchId(f.FileId).Id,
				})
			}
			_ = zipReader.Close()
			_ = os.Remove("./fiction/zip/" + f.FileName)
		}
	}
	bot.SendGroupMsg(event.GroupId, "文件初始化完成"+go_mybots.MessageAt(event.UserId).Message, false)

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

func Download(fileName, randNum string, url string, isZip bool, resultFileName string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	client := http.Client{}
	if isZip {
		response, err := client.Get(url)
		if err != nil {
			panic(err.Error())
		}
		defer response.Body.Close()
		content, _ := ioutil.ReadAll(response.Body)
		file, err := os.OpenFile("./fiction/zip/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err.Error())
		}
		_, err = file.WriteString(string(content))
		if err != nil {
			panic(err.Error())
		}
		file.Close()
		err = DeCompress("./fiction/zip/"+fileName, "./fiction/")
		if err != nil {
			panic(err.Error())
		}

		zipFile, err := zip.OpenReader("./fiction/zip/" + fileName)
		if err != nil {
			panic(err.Error())
		}
		for _, i2 := range zipFile.File {
			if i2.Name != resultFileName {
				_ = os.Remove("./fiction/" + i2.Name)
			}

		}
		err = zipFile.Close()
		if err != nil {
			panic(err)
		}
		_ = os.Remove("./fiction/zip/" + fileName)
		time.Sleep(time.Duration(120) * time.Second)
		err = os.Remove("./fiction/" + resultFileName)
		if err != nil {
			panic(err)
		}
		delete(File, randNum)
	} else {
		response, err := client.Get(url)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		file, err := os.OpenFile("./fiction/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err.Error())
		}
		_, err = file.WriteString(string(content))
		if err != nil {
			panic(err.Error())
		}
		file.Close()
		time.Sleep(time.Duration(120) * time.Second)
		delete(File, randNum)
		err = os.Remove("./fiction/" + fileName)
		if err != nil {
			panic(err)
		}
	}
}

func CopyFile(srcFileName string, dstFileName string) {
	defer func() {
		i := recover()
		if i != nil {
			log.Println("文件移动失败")
		}
	}()
	//打开源文件
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		log.Panic("源文件读取失败,原因是:%v\n", err)
	}
	defer func() {
		err = srcFile.Close()
		if err != nil {
			log.Panic("源文件关闭失败,原因是:%v\n", err)
		}
	}()

	//创建目标文件,稍后会向这个目标文件写入拷贝内容
	distFile, err := os.Create(dstFileName)
	if err != nil {
		log.Panic("目标文件创建失败,原因是:%v\n", err)
	}
	defer func() {
		err = distFile.Close()
		if err != nil {
			log.Panic("目标文件关闭失败,原因是:%v\n", err)
		}
	}()
	//定义指定长度的字节切片,每次最多读取指定长度
	var tmp = make([]byte, 1024*4)
	//循环读取并写入
	for {
		n, err := srcFile.Read(tmp)
		n, _ = distFile.Write(tmp[:n])
		if err != nil {
			if err == io.EOF { //读到了文件末尾,并且写入完毕,任务完成返回(关闭文件的操作由defer来完成)
				return
			} else {
				log.Panic("拷贝过程中发生错误,错误原因为:%v\n", err)
			}
		}
	}
}

func DeCompress(zipFile string, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(getDir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

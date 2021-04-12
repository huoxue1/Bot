package models

import (
	log "github.com/sirupsen/logrus"
	"strconv"
)

type File struct {
	Id       int    `xorm:"not null pk autoincr INT(11)"`
	Filename string `xorm:"not null VARCHAR(100)"`
	Fileid   string `xorm:"VARCHAR(100)"`
	Busid    int    `xorm:"INT(11)"`
	Ischild  int    `xorm:"not null default 0 comment('是否属于压缩包下文件') TINYINT(1)"`
	Iszip    int    `xorm:"not null default 1 comment('自身是否是压缩包文件') TINYINT(1)"`
	Groupid  string `xorm:"not null VARCHAR(100)"`
	Pid      int    `xorm:"comment('上级压缩包文件id') INT(11)"`
}

func FileInsert(file *File) bool {
	_, err := X.Insert(file)
	if err != nil {
		log.Debugln(err)
		return false
	}
	return true
}

func FileSearchByGroup(groupId int64) map[int]*File {
	files := make(map[int]*File)
	err := X.Where("groupId=?", strconv.FormatInt(groupId, 10)).Find(&files)
	if err != nil {
		log.Debugln(err)
		return nil
	}
	return files
}

func FileSearchALL() []*File {
	files := make(map[int]*File)
	err := X.Find(&files)
	if err != nil {
		log.Debugln(err)
		return nil
	}
	var f []*File
	for _, file := range files {
		f = append(f, file)
	}
	return f
}

func FileSearchById(id int) *File {
	file := new(File)
	_, err := X.Where("id=?", id).Get(file)
	if err != nil {
		log.Debugln(err)
		return nil
	}
	return file
}

func FileSearchId(fileId string) *File {
	file := new(File)
	_, err := X.Where("fileId=?", fileId).Get(file)
	if err != nil {
		log.Debugln(err)
		return nil
	}
	return file
}

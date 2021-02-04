package model

import (
	"strconv"
)

type File struct {
	Id       int    `json:"id" db:"Id"`
	FileName string `json:"file_name" db:"fileName"`
	FileId   string `json:"file_id" db:"fileId"`
	BusId    int    `json:"busid" db:"busid"`
	IsChild  bool   `json:"is_child" db:"isChild"`
	IsZip    bool   `json:"is_zip" db:"isZip"`
	GroupId  string `json:"group_id" db:"groupId"`
	Pid      int    `json:"pid" db:"pid"`
}

func (con *Connect) FileInsert(file File) bool {
	se := con.FileSearchId(file.FileId)
	if se.Id != 0 {
		_, _ = con.Db.Exec("update file set fileName=?,file.fileId=?,file.busid = ?, file.isChild = ?,file.isZip = ?,file.groupId=?,file.pid=? where Id = ?",
			file.FileName, file.FileId, file.BusId, file.IsChild, file.IsZip, file.GroupId, file.Pid, se.Id)
	} else {

		_, err := con.Db.Exec("insert into file (filename, fileid, busid, ischild, iszip, groupid, pid) VALUES (?,?,?,?,?,?,?)",
			file.FileName, file.FileId, file.BusId, file.IsChild, file.IsZip, file.GroupId, file.Pid)
		if err != nil {
			return false
		}
		return true
	}
	return true
}

func (con *Connect) FileSearch(groupId int) []File {
	var fileList []File
	err := con.Db.Select(&fileList, "select * from file where groupId = ? and isZip = 0", strconv.Itoa(groupId))
	if err != nil {
		return []File{}
	} else {
		return fileList
	}
}

func (con *Connect) FileSearchById(id int) File {
	var file File
	err := con.Db.Get(&file, "select * from file where Id = ?", id)
	if err != nil {
		return File{}
	} else {
		return file
	}
}

func (con *Connect) FileSearchId(fileId string) File {
	var file File
	err := con.Db.Get(&file, "select * from file where file.fileId = ?", fileId)
	if err != nil {
		return File{}
	} else {
		return file
	}
}

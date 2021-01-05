package Integral

import (
	"fmt"
	gomybots "github.com/3343780376/go-mybots"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"log"
	"strconv"
	"time"
)

var (
	file  *excelize.File
	rows  [][]string
	err   error
	sheet string
)

type Xlsx struct {
	Event gomybots.Event `json:"event"`
	Sheet string         `json:"sheet"`
}

type Do interface {
	XlsxInit() error
	InsertData() error
	IsExit() bool
	Increase(n int) (bool, error)
	Decrease(n int) (bool, error)
	IsSign() (bool, error)
}

func (x Xlsx) XlsxInit() error {
	if x.Event.GroupId == 727429388 {
		sheet = "fei"
	} else if x.Event.GroupId == 17185204 {
		sheet = "xi"
	} else {
		sheet = "other"
	}
	var err error
	file, err = excelize.OpenFile("../Bot/templete/group.xlsx")
	if err != nil {
		return err
	}
	rows, err = file.GetRows(sheet)
	return err
}

//插入积分数据
func (x Xlsx) InsertData() error {

	err = file.SetCellValue(sheet, fmt.Sprintf("A%d", len(rows)+1), x.Event.UserId)
	err = file.SetCellValue(sheet, fmt.Sprintf("B%d", len(rows)+1), 10)
	err = file.SetCellValue(sheet, fmt.Sprintf("C%d", len(rows)+1), x.Event.Sender.Card)
	err = file.SetCellValue(sheet, fmt.Sprintf("D%d", len(rows)+1), 0)
	err = file.Save()
	return err
}

//判断积分数据是否存在
func (x Xlsx) IsExit() bool {
	if err != nil {
		log.Println(err)
	}
	for _, i2 := range rows {
		if len(i2) == 0 {
			continue
		} else {
			if i2[0] == strconv.Itoa(x.Event.UserId) {
				return true
			}
		}
	}
	return false
}

//增加n个积分
func (x Xlsx) Increase(n int) (bool, error) {
	defer file.Save()
	var err error
	if !x.IsExit() {
		err = x.InsertData()
	}
	for i1, i2 := range rows {
		if len(i2) != 0 {
			if i2[0] == strconv.Itoa(x.Event.UserId) {
				oldData, err := strconv.ParseInt(i2[1], 10, 32)
				err = file.SetCellValue(sheet, fmt.Sprintf("B%v", i1+1), int(oldData)+n)
				return true, err
			}
		}

	}
	return false, err
}

//减少n个积分
func (x Xlsx) Decrease(n int) (bool, error) {
	defer file.Save()
	var err error
	if !x.IsExit() {
		err = x.InsertData()
	}
	for i1, i2 := range rows {
		if len(i2) != 0 {
			if i2[0] == strconv.Itoa(x.Event.UserId) {
				oldData, err := strconv.ParseInt(i2[1], 10, 32)
				err = file.SetCellValue(sheet, fmt.Sprintf("B%v", i1+1), int(oldData)-n)
				return true, err
			}
		}

	}
	return false, err
}

//判断是否签到
func (x Xlsx) IsSign() (bool, error) {
	day := time.Now().Day()
	defer file.Save()
	var err error
	if !x.IsExit() {
		err = x.InsertData()
	}
	for i1, i2 := range rows {
		if len(i2) != 0 {
			if i2[0] == strconv.Itoa(x.Event.UserId) {
				i, err := strconv.ParseInt(i2[3], 10, 32)
				if int(i) != day {
					_, err := x.Increase(2)
					err = file.SetCellValue(sheet, fmt.Sprintf("D%d", i1+1), day)
					return true, err
				} else {
					return false, err
				}
			}
		}

	}
	return false, err
}

func (x Xlsx) SearchNum() (int, error) {
	//day := time.Now().Day()
	defer file.Save()
	var err error
	if !x.IsExit() {
		err = x.InsertData()
	}
	for _, i2 := range rows {
		if len(i2) != 0 {
			if i2[0] == strconv.Itoa(x.Event.UserId) {
				i, err := strconv.ParseInt(i2[1], 10, 32)
				return int(i), err
			}
		}
	}
	return 0, err
}

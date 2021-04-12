package models

import (
	go_mybots "github.com/3343780376/go-mybots"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Sign struct {
	Userid string `xorm:"not null pk comment('qq') VARCHAR(15)"`
	Num    int    `xorm:"not null default 10 comment('积分数值') INT(11)"`
	Card   string `xorm:"VARCHAR(100)"`
	Day    int    `xorm:"not null default 10 comment('上次签到日期') INT(11)"`
}

func SelectSign(event go_mybots.Event) int {
	Exist(event)
	s := new(Sign)
	s.Userid = strconv.Itoa(event.UserId)
	_, err := X.Get(s)
	if err != nil {
		log.Debugln(err)
	}
	return s.Num
}

func Exist(event go_mybots.Event) {
	s := new(Sign)
	s.Userid = strconv.Itoa(event.UserId)
	exist, err := X.Exist(s)
	if !exist || err != nil {
		s.Card = event.Sender.Card
		s.Num = 10
		s.Day = time.Now().Day() - 1
		_, err := X.Insert(s)
		if err != nil {
			log.Debugln(err)
		}
	}
}

/*
	return:
		已经签到 ：true
		未签到： false
*/
func IsSign(event go_mybots.Event) bool {
	Exist(event)
	s := new(Sign)
	s.Userid = strconv.Itoa(event.UserId)
	_, err := X.Get(s)
	if err != nil {
		log.Debugln(err)
		return false
	}
	if s.Day == time.Now().Day() {
		return true
	} else {
		Update(2, event)

		_, err := X.Where("userId=?", strconv.Itoa(event.UserId)).Update(&Sign{Day: time.Now().Day()})
		if err != nil {
			log.Debugln(err)
		}
		return false
	}
}

func Update(n int, event go_mybots.Event) {
	Exist(event)
	sign := SelectSign(event)
	_, err := X.Where("userId=?", strconv.Itoa(event.UserId)).Update(&Sign{Num: sign + n})
	if err != nil {
		log.Debugln(err)
	}

}

package curse

import (
	"fmt"
	"github.com/3343780376/leafBot"
	"os"
	"strconv"
	"time"
)

var ()

func init() {
	leafBot.AddCommandHandle(Curse, ".curse", []string{"课程表"}, []leafBot.Rule{{func(event leafBot.Event, i ...interface{}) bool {
		if event.SelfId == leafBot.DefaultConfig.Admin {
			return true
		}
		return false
	}, nil}}, 10, false)
}

func Curse(event leafBot.Event, bot *leafBot.Bot, args []string) {
	path, _ := os.Getwd()

	if event.GroupId == 681168003 || event.UserId == leafBot.DefaultConfig.Admin {
		if len(args) <= 1 {
			if event.MessageType == "group" {
				_ = bot.SendMsg(event.MessageType, event.UserId, event.GroupId, fmt.Sprintf("[CQ:image,file=file:///%v]", path+"/templete/image/week"+strconv.Itoa(getWeek())+".jpg"), false)
			} else {
				_ = bot.SendMsg(event.MessageType, event.UserId, event.GroupId, fmt.Sprintf("[CQ:image,file=file:///%v]", path+"/templete/image/week"+strconv.Itoa(getWeek())+".jpg"), false)
			}
		}
	}
}

func getWeek() int {
	week := 2
	mouth, _ := strconv.Atoi(time.Now().Format("1"))
	day := time.Now().Day()

	if mouth == 3 && (day >= 8 && day <= 14) {
		week = 3
	} else if mouth == 3 && (day >= 15 && day <= 21) {
		week = 4
	} else if mouth == 3 && (day >= 22 && day <= 28) {
		week = 5
	} else if (mouth == 3 && (day >= 29 && day <= 31)) || (mouth == 4 && (day >= 1 && day <= 4)) {
		week = 6
	} else if mouth == 4 && (day >= 5 && day <= 11) {
		week = 7
	} else if mouth == 4 && (day >= 12 && day <= 18) {
		week = 8
	} else if mouth == 4 && (day >= 19 && day <= 25) {
		week = 9
	} else if (mouth == 4 && (day >= 26 && day <= 30)) || (mouth == 5 && (day >= 1 && day <= 2)) {
		week = 10
	} else if mouth == 5 && (day >= 3 && day <= 9) {
		week = 11
	} else if mouth == 5 && (day >= 10 && day <= 16) {
		week = 12
	} else if mouth == 5 && (day >= 17 && day <= 23) {
		week = 13
	} else if mouth == 5 && (day >= 24 && day <= 30) {
		week = 14
	} else if (mouth == 5 && day == 31) || (mouth == 6 && (day >= 1 && day <= 6)) {
		week = 15
	} else if mouth == 6 && (day >= 7 && day <= 13) {
		week = 16
	} else if mouth == 6 && (day >= 14 && day <= 20) {
		week = 17
	}
	return week
}

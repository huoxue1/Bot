package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	log "github.com/sirupsen/logrus"
)

var (
	X *xorm.Engine
)

func Xorminit() {
	var err error
	X, err = xorm.NewEngine("mysql", "user1:1743224847gou#@(rm-bp16ar289akg5k41hjo.mysql.rds.aliyuncs.com)/groupdb?charset=utf8")
	if err != nil {
		log.Fatal("数据库链接失败")
	}
}

package fileSearch

import (
	"encoding/csv"
	"log"
	"os"
	"sync"
)

type Csv struct {
	File  *os.File
	Mutex sync.Mutex
}

func CsvInit() *Csv {
	file, err := os.OpenFile("./templete/fileData.csv", os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return &Csv{
		File:  file,
		Mutex: sync.Mutex{},
	}
}

func (c *Csv) CsvWrite(data []string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("csv文件读写错误", err)
		}
	}()
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	w := csv.NewWriter(c.File)
	w.Comma = ','
	w.UseCRLF = true
	err := w.Write(data)
	if err != nil {
		log.Panic(err)
	}
	w.Flush()
}

func (c *Csv) CsvRead() [][]string {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("csv文件读取错误")
		}
	}()
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	w := csv.NewReader(c.File)
	data, err := w.ReadAll()
	if err != nil {
		panic(err.Error())
	}
	return data
}

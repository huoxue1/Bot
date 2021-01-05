package main

import (
	"log"
	"regexp"
	"strings"
)

func main() {
	s := "打劫[CQ:at,qq=3343780376]123"
	println(strings.Contains(s, "打劫"))
	compile := regexp.MustCompile(`(\d+)`)
	findString := compile.FindAllString(s,-1)
	log.Println(findString)


}

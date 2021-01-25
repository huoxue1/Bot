package main

import (
	"fmt"
	"strconv"
)

type name struct {
	Id int64
}

func main() {
	n := name{Id: 13121321321321}
	str := strconv.FormatInt(n.Id, 10)
	fmt.Println(str)
}

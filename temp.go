package main

import (
	"fmt"
	"strings"
)

func main() {
	split := strings.Split("禁言 10", " ")
	for _, s := range split {
		fmt.Println(s)
	}
}

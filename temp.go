package main

import (
	"math/rand"
	"time"
)

func main() {
	println(time.Now().UnixNano())
	println(time.Now().Unix())
	rand.Seed(time.Now().UnixNano())
	println(rand.Intn(1000))
}

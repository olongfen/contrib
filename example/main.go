package main

import "github.com/srlemon/go-log"

func main() {
	l := log.NewLogFile("demo.log")
	l1, _ := log.NewLog(l)
	l1.Println("aaaaaaaaaa")
}
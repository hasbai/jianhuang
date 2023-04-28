package main

import (
	"os"
	"os/signal"
	"syscall"
)

var url = "https://teicn.oss-cn-hongkong.aliyuncs.com/teicarmx64.7z"
var threads = 2

func main() {
	supervisor := NewSupervisor()
	go supervisor.Run()
	for i := 0; i < threads; i++ {
		supervisor.AddRunner()
	}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	os.Exit(0)
}

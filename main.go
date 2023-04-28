package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"math"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var url = "https://teicn.oss-cn-hongkong.aliyuncs.com/teicarmx64.7z"
var threads = int(math.Max(4, float64(runtime.NumCPU())))

func main() {
	app := &cli.App{
		Name:  "jian huang",
		Usage: "http benchmark tool",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "concurrency",
				Aliases: []string{"c"},
				Usage:   "specify concurrency",
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Args().First() != "" {
				url = ctx.Args().First()
			}
			if ctx.Int("concurrency") > 0 {
				threads = ctx.Int("concurrency")
			}
			run()
			return nil
		},
	}

	panic(app.Run(os.Args))
}

func run() {
	supervisor := NewSupervisor()
	go supervisor.Run()
	for i := 0; i < threads; i++ {
		supervisor.AddRunner()
	}
	fmt.Println("Benchmarking", url)
	fmt.Println("Concurrency:", threads)
	//fmt.Println("start api server...")
	//panic(http.ListenAndServe("localhost:8080", nil))

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	os.Exit(0)
}

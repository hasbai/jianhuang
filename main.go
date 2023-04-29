package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var url = ""
var threads = int(math.Max(4, float64(runtime.NumCPU())))

func main() {
	app := &cli.App{
		Name:  "jian huang",
		Usage: "http benchmark tool",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "concurrency",
				Aliases: []string{"c"},
				Usage:   "specify concurrency, default is cpu cores, minimum is 4",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "enable debug mode",
			},
		},
		Action: func(ctx *cli.Context) error {
			url = ctx.Args().First()
			// https://teicn.oss-cn-hongkong.aliyuncs.com/teicarmx64.7z
			if url == "" {
				return nil
			}
			if ctx.Int("concurrency") > 0 {
				threads = ctx.Int("concurrency")
			}
			if ctx.Bool("debug") {
				go profile()
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

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	os.Exit(0)
}

func profile() {
	fmt.Println("start api server...")
	panic(http.ListenAndServe("localhost:8080", nil))
}

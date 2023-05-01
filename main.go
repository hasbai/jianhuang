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

const UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"

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
			url := ctx.Args().First()
			// https://teicn.oss-cn-hongkong.aliyuncs.com/teicarmx64.7z
			if url == "" {
				return nil
			}
			threads := int(math.Max(4, float64(runtime.NumCPU())))
			if ctx.Int("concurrency") > 0 {
				threads = ctx.Int("concurrency")
			}
			if ctx.Bool("debug") {
				go profile()
			}
			run(url, threads)
			return nil
		},
	}

	panic(app.Run(os.Args))
}

func run(url string, threads int) {
	s := NewSupervisor(url)
	s.Run(threads)

	go showStat(url, threads)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
	os.Exit(0)
}

func profile() {
	fmt.Println("profile server started at http://localhost:8080/debug/pprof/")
	panic(http.ListenAndServe("localhost:8080", nil))
}

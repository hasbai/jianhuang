package main

import (
	"fmt"
	"github.com/apoorvam/goterminal"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	netstat "github.com/shirou/gopsutil/net"
	"math"
	"os"
	"time"
)

var TerminalWriter = goterminal.New(os.Stdout)

func showStat(url string, threads int) {
	initialNetCounter, _ := netstat.IOCounters(true)

	for {
		percent, _ := cpu.Percent(time.Second, false)
		memStat, _ := mem.VirtualMemory()
		netCounter, _ := netstat.IOCounters(true)
		loadStat, _ := load.Avg()

		fmt.Fprintf(TerminalWriter, "Benchmarking: %s\n", url)
		fmt.Fprintf(TerminalWriter, "Concurrency: %d\n", threads)

		fmt.Fprintf(TerminalWriter, "CPU: %.2f%% \n", percent)
		fmt.Fprintf(TerminalWriter, "Memory: %.2f%% \n", memStat.UsedPercent)
		fmt.Fprintf(TerminalWriter, "Load: %.2f %.2f %.2f\n", loadStat.Load1, loadStat.Load5, loadStat.Load15)
		for i := 0; i < len(netCounter); i++ {
			if netCounter[i].BytesRecv == 0 && netCounter[i].BytesSent == 0 {
				continue
			}
			receivedBytes := netCounter[i].BytesRecv - initialNetCounter[i].BytesRecv
			sentBytes := netCounter[i].BytesSent - initialNetCounter[i].BytesSent
			fmt.Fprintf(TerminalWriter, "Nic:%v,Recv %s(%s/s),Send %s(%s/s)\n", netCounter[i].Name,
				readableBytes(netCounter[i].BytesRecv),
				readableBytes(receivedBytes),
				readableBytes(netCounter[i].BytesSent),
				readableBytes(sentBytes))
		}
		initialNetCounter = netCounter
		TerminalWriter.Clear()
		_ = TerminalWriter.Print()
		time.Sleep(1 * time.Millisecond)
	}
}

var sizes = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

func readableBytes(bytes uint64) (expression string) {
	if bytes == 0 {
		return "0B"
	}
	i := math.Ilogb(float64(bytes)) / 10
	pow := 1 << (i * 10)
	return fmt.Sprintf("%.2f%s", float64(bytes)/float64(pow), sizes[i])
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GitHamburg/windows-agent/cron"
	"github.com/GitHamburg/windows-agent/funcs"
	"github.com/GitHamburg/windows-agent/g"
	"github.com/GitHamburg/windows-agent/http"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	check := flag.Bool("check", false, "check collector")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.InitLog()

	g.InitRootDir()
	g.InitLocalIps()
	g.InitRpcClients()

	funcs.BuildMappers()

	go cron.InitDataHistory()

	cron.ReportAgentStatus()
	cron.SyncMinePlugins()
	cron.SyncBuiltinMetrics()
	cron.SyncTrustableIps()
	cron.Collect()

	go http.Start()

	select {}
}

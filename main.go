package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"

	"os"

	"github.com/gin-gonic/gin"
	"xdt.com/hm-diag/api"
	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/devdis"
	"xdt.com/hm-diag/diag"
	"xdt.com/hm-diag/regist"
	"xdt.com/hm-diag/util"
)

// TODO: unit test
// TODO: log format -- TRACE,INFO,WARN,ERROR

//go:embed web/release/*
var webFS embed.FS

//go:embed api/swagger_ui/*
var swagFS embed.FS

type Opt struct {
	Port          int
	GitRepoUrl    string
	MinerUrl      string
	IntervalSec   uint
	GitRepoDir    string
	LanDevIntface string
	Verbose       bool
}

var opt Opt

func usage() {
	fmt.Fprintf(os.Stdout, "Helium Diagnostic\n")
	fmt.Fprintf(os.Stdout, "Usage: hm-diag [options] [get|server] \n\n")
	fmt.Fprintf(os.Stdout, "Subcommand:\n")
	fmt.Fprintf(os.Stdout, "  get\n    get info to stdout\n")
	fmt.Fprintf(os.Stdout, "  server\n    run http server, can omit it\n")
	fmt.Fprintf(os.Stdout, "Options:\n")
	flag.PrintDefaults()
}

func init() {
	flag.IntVar(&opt.Port, "p", 8090, "server listening port")
	flag.StringVar(&opt.LanDevIntface, "lan", "eth0", "lan device discovery net interface")
	flag.StringVar(&opt.MinerUrl, "m", "http://127.0.0.1:4467", "miner http url")
	flag.StringVar(&opt.GitRepoDir, "gitRepo",
		"/home/pi/hnt_iot", "program docker-compose working git dir")
	flag.StringVar(&opt.GitRepoUrl, "gitRepoUrl",
		"https://github.com/HummingbirdIot/hnt_iot_release.git", "hnt iot git url")
	flag.UintVar(&opt.IntervalSec, "i", 30, "data refresh interval in seconds")
	flag.BoolVar(&opt.Verbose, "v", false, "verbose log")
	flag.Usage = usage

	flag.Parse()
	// store config as global
	config.InitConf(config.GlobalConfig{
		ApiPort:       opt.Port,
		LanDevIntface: opt.LanDevIntface,
		MinerUrl:      opt.MinerUrl,
		GitRepoDir:    opt.GitRepoDir,
		GitRepoUrl:    opt.GitRepoUrl,
		IntervalSec:   opt.IntervalSec,
	})
}

func main() {
	diag.InitTask(*config.Config())
	diagTask := diag.TaskInstance()

	if flag.Arg(0) == "get" {
		if !opt.Verbose {
			log.SetOutput(io.Discard)
		}
		diagTask.Do()
		s, _ := json.MarshalIndent(diagTask.Data(), "", "  ")
		os.Stdout.WriteString(string(s))
		return
	} else if flag.Arg(0) == "server" || flag.Arg(0) == "" {
		optJson, _ := json.Marshal(opt)
		log.Println("options: ", string(optJson))

		// init job
		go diagTask.StartTaskJob(true)
		go regist.StartRegistJob()
		util.Sgo(devdis.Init, "init device discovery error")

		// http server
		r := gin.Default()
		r.Use(api.CORSMiddleware())
		api.Route(r, webFS, swagFS)
		r.Run(fmt.Sprintf(":%d", opt.Port))
	} else {
		flag.Usage()
	}
}

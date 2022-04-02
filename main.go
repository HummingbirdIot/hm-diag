package main

import (
	"context"
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
	"xdt.com/hm-diag/link"
	"xdt.com/hm-diag/regist"
	"xdt.com/hm-diag/util"
)

// set by build with -ldflags eg:
// go build -ldflags "-X main.Version=${version} -X main.Githash=`git rev-parse HEAD`"
var (
	Githash = ""
	Version = ""
)

// TODO: unit test
// TODO: log format -- TRACE,INFO,WARN,ERROR

var (
	//go:embed web/release/*
	webFS embed.FS

	//go:embed api/swagger_ui/*
	swagFS embed.FS
)

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
	fmt.Fprintf(os.Stdout, "  version\n    get version info to stdout\n")
	fmt.Fprintf(os.Stdout, "Options:\n")
	flag.PrintDefaults()
}

func init() {
	setVersion()

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

func setVersion() {
	config.Version = Version
	config.Githash = Githash
}

func main() {
	diag.InitTask(*config.Config())
	diagTask := diag.TaskInstance()
	rootCtx := context.Background()
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

		// link
		link.Start(rootCtx)

		// http server
		r := gin.Default()
		r.Use(api.CORSMiddleware()).Use(api.PrivateAccessMiddle())
		api.Route(r, webFS, swagFS)
		r.Run(fmt.Sprintf(":%d", opt.Port))
	} else if flag.Arg(0) == "version" {
		fmt.Printf("version: %s githash: %s\n", Version, Githash)
	} else {
		flag.Usage()
	}
}

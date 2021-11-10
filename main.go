package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"

	"net/http"
	"os"

	"xdt.com/hm-diag/diag"
	"xdt.com/hm-diag/regist"
)

type Opt struct {
	Port        int
	MinerUrl    string
	IntervalSec uint
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
	flag.StringVar(&opt.MinerUrl, "m", "http://127.0.0.1:4467", "miner http url")
	flag.UintVar(&opt.IntervalSec, "i", 30, "data refresh interval in seconds")
	flag.Usage = usage
}

var task *diag.Task

func main() {
	flag.Parse()
	task = &diag.Task{Config: diag.TaskConfig{MinerUrl: opt.MinerUrl, IntervalSec: opt.IntervalSec}}
	if flag.Arg(0) == "get" {
		log.SetOutput(io.Discard)
		task.Do()
		s, _ := json.MarshalIndent(task.Data(), "", "  ")
		os.Stdout.WriteString(string(s))
		return
	} else if flag.Arg(0) == "server" || flag.Arg(0) == "" {
		optJson, _ := json.Marshal(opt)
		log.Println("options: ", string(optJson))
		go task.StartTaskJob(true)
		register := &regist.Register{ApiPort: opt.Port, RegistryApiPort: 6753, ReistIntervalSec: 30}
		go register.StartRegistJob()

		route(task, register)

		log.Println("server listening on port " + strconv.Itoa(opt.Port))
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(opt.Port), nil))
	} else {
		flag.Usage()
	}
}

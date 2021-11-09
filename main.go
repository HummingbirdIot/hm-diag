package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"

	"net/http"
	"os"

	"xdt.com/hm-diag/diag"
)

type Opt struct {
	Port        string
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
	flag.StringVar(&opt.Port, "p", "8090", "server listening port")
	flag.StringVar(&opt.MinerUrl, "m", "http://127.0.0.1:4467", "miner http url")
	flag.UintVar(&opt.IntervalSec, "i", 30, "data refresh interval in seconds")
	flag.Usage = usage
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		return
	}
	w.Header().Add("Content-type", "application/json")
	d := task.GetData()
	d.Data["a-notice"] = `do not use this api path "/" to integrate, use api under path "api/"`
	j, _ := json.MarshalIndent(d, "", "  ")
	fmt.Fprint(w, string(j))
}

func minerInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	d := task.GetMinerInfo()
	j, _ := json.MarshalIndent(d, "", "  ")
	fmt.Fprint(w, string(j))
}

func hardwareInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	d := task.GetHardwareInfo()
	j, _ := json.MarshalIndent(d, "", "  ")
	fmt.Fprint(w, string(j))
}

var task diag.Task

func main() {
	flag.Parse()
	task = diag.Task{Config: diag.TaskConfig{MinerUrl: opt.MinerUrl, IntervalSec: opt.IntervalSec}}
	if flag.Arg(0) == "get" {
		log.SetOutput(io.Discard)
		task.DoTask()
		s, _ := json.MarshalIndent(task.GetData(), "", "  ")
		os.Stdout.WriteString(string(s))
		return
	} else if flag.Arg(0) == "server" || flag.Arg(0) == "" {
		optJson, _ := json.Marshal(opt)
		log.Println("options: ", string(optJson))
		go task.StartTask(true)
		http.HandleFunc("/", homeHandler)
		http.HandleFunc("/api/v1/hardware", hardwareInfoHandler)
		http.HandleFunc("/api/v1/miner", minerInfoHandler)
		log.Println("server listening on port " + opt.Port)
		log.Fatal(http.ListenAndServe(":"+opt.Port, nil))
	} else {
		flag.Usage()
	}
}

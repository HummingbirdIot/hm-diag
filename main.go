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
	fmt.Fprintf(os.Stdout, "Usage: [options] [get | server] \n\n")
	fmt.Fprintf(os.Stdout, "Subcommand:\n")
	fmt.Fprintf(os.Stdout, "  get\n	get info to stdout\n")
	fmt.Fprintf(os.Stdout, "  server\n	run http server, can omit it\n")
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
	w.Header().Add("Content-type", "application/json")
	d := task.GetData()
	j, _ := json.MarshalIndent(d, "", "  ")
	fmt.Fprint(w, string(j))

	// t1, err := template.ParseFiles("tmpl/index.html")
	// if err != nil {
	// 	panic(err)
	// }
	// t1.Execute(w, d)
}

var task diag.Task

func main() {
	flag.Parse()
	optJson, _ := json.Marshal(opt)
	log.Println("options: ", string(optJson))
	task = diag.Task{Config: diag.TaskConfig{MinerUrl: opt.MinerUrl, IntervalSec: opt.IntervalSec}}
	if flag.Arg(0) == "get" {
		log.SetOutput(io.Discard)
		task.DoTask()
		s, _ := json.MarshalIndent(task.GetData(), "", "  ")
		os.Stdout.WriteString(string(s))
		return
	}
	go task.StartTask(true)
	http.HandleFunc("/", homeHandler)
	log.Println("server listening on port " + opt.Port)
	log.Fatal(http.ListenAndServe(":"+opt.Port, nil))
}

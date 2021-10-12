package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"xdt.com/hm-diag/miner"
)

var (
	port     string
	minerUrl string
)

func usage() {
	fmt.Fprintf(os.Stderr, "Helium Diagnostic\n")
	fmt.Fprintf(os.Stderr, "Usage: [options] [get | server] \n\n")
	fmt.Fprintf(os.Stderr, "Subcommand:\n")
	fmt.Fprintf(os.Stderr, "  get\n	get info to stdout\n")
	fmt.Fprintf(os.Stderr, "  server\n	run http server, can omit it\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&port, "p", "8090", "server listening port")
	flag.StringVar(&minerUrl, "m", "http://127.0.0.1:4467", "miner http url")
	flag.Usage = usage
}

type CacheData struct {
	time time.Time
	data interface{}
}

var minerCacheData CacheData

func getMinerData() []byte {
	fetchTime := time.Now()
	res := miner.FetchData(minerUrl)
	res["fetch_time"] = fetchTime
	jsonBuf, _ := json.MarshalIndent(res, " ", " ")
	return jsonBuf
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	delay, _ := time.ParseDuration("30s")
	if minerCacheData.data == nil || minerCacheData.time.Add(delay).Before(time.Now()) {
		fetchTime := time.Now()
		jsonBuf := getMinerData()
		minerCacheData.time = fetchTime
		minerCacheData.data = string(jsonBuf)
	}
	w.Header().Add("Content-type", "application/json")
	fmt.Fprintf(w, "%s\n", minerCacheData.data)

	// t1, err := template.ParseFiles("tmpl/index.html")
	// if err != nil {
	// 	panic(err)
	// }
	// t1.Execute(w, res)
}

func main() {
	flag.Parse()
	if flag.Arg(0) == "get" {
		log.SetOutput(ioutil.Discard)
		jsonBuf := getMinerData()
		os.Stdout.WriteString(string(jsonBuf) + "\n")
		return
	}
	delay, _ := time.ParseDuration("-10d")
	minerCacheData = CacheData{time: time.Now().Add(delay)}
	http.HandleFunc("/", homeHandler)
	log.Println("server listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

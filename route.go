package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"xdt.com/hm-diag/ctrl"
	"xdt.com/hm-diag/diag"
	"xdt.com/hm-diag/regist"
)

var diagTask *diag.Task
var register *regist.Register

func route(_diagTask *diag.Task, _register *regist.Register) {
	diagTask = _diagTask
	register = _register
	RouteState()
	RouteCtrl()
}

func RouteCtrl() {
	http.HandleFunc("/api/v1/device/reboot", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		res := map[string]interface{}{
			"msg": "receive reboot request, to reboot",
		}
		j, _ := json.Marshal(res)
		w.Write(j)
		go ctrl.RebootDevice()
	})

	http.HandleFunc("/api/v1/miner/resync", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		err := ctrl.ResyncMiner()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := map[string]interface{}{
				"msg": "receive resync miner request, but got error:" + err.Error(),
			}
			j, _ := json.Marshal(res)
			w.Write(j)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
}

func RouteState() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/v1/device/state", deviceInfoHandler)
	http.HandleFunc("/api/v1/miner/state", minerInfoHandler)
	http.HandleFunc("/registInfo", registInfoHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		return
	}
	w.Header().Add("Content-type", "application/json")
	d := diagTask.Data()
	if d.Data != nil {
		d.Data["aNotice"] = `do not use this api path "/" to integrate, use api under path "api/"`
	}
	j, _ := json.MarshalIndent(d, "", "  ")
	fmt.Fprint(w, string(j))
}

func minerInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	var d diag.TaskData
	if c := r.URL.Query().Get("cache"); c == "true" {
		d = diagTask.MinerInfo()
	} else {
		d = diag.TaskData{Data: diagTask.FetchMinerInfo(), FetchTime: time.Now()}
	}
	j, _ := json.MarshalIndent(d, "", "  ")
	fmt.Fprint(w, string(j))
}

func deviceInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	var d diag.TaskData
	if c := r.URL.Query().Get("cache"); c == "true" {
		d = diagTask.DeviceInfo()
	} else {
		d = diag.TaskData{Data: diagTask.FetchDeviceInfo(), FetchTime: time.Now()}
	}
	j, _ := json.MarshalIndent(d, "", "  ")
	fmt.Fprint(w, string(j))
}

func registInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	d, err := register.GetRegistInfo()
	var str string
	if err != nil {
		str = "error: " + err.Error()
	} else {
		j, _ := json.MarshalIndent(d, "", "  ")
		str = string(j)
	}
	fmt.Fprint(w, str)
}

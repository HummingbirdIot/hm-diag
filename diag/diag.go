package diag

import (
	"encoding/json"
	"log"
	"time"

	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/diag/device"
	"xdt.com/hm-diag/diag/miner"
	"xdt.com/hm-diag/link"
	"xdt.com/hm-diag/link/message"
)

type TaskConfig struct {
	MinerUrl    string
	IntervalSec uint
}

type Task struct {
	Config     TaskConfig
	time       time.Time
	data       map[string]interface{}
	taskTicker *time.Ticker
	quitTask   chan struct{}
}

type TaskData struct {
	FetchTime time.Time              `json:"fetchTime"`
	Data      map[string]interface{} `json:"data"`
}

type AllStateInfo struct {
	Device device.DeviceInfo `json:"device"`
	Miner  miner.MinerInfo   `json:"miner"`
}

var taskSingleton *Task

func InitTask(conf config.GlobalConfig) {
	taskSingleton = &Task{Config: TaskConfig{MinerUrl: conf.MinerUrl, IntervalSec: conf.IntervalSec}}
}

func TaskInstance() *Task {
	return taskSingleton
}

func (t *Task) Data() TaskData {
	return TaskData{Data: t.data, FetchTime: t.time}
}

func (t *Task) DeviceInfo() TaskData {
	var data map[string]interface{}
	if t.data != nil {
		data = t.data["device"].(map[string]interface{})
	}
	return TaskData{Data: data, FetchTime: t.time}
}

func (t *Task) MinerInfo() TaskData {
	var data map[string]interface{}
	if t.data != nil {
		data = t.data["miner"].(map[string]interface{})
	}
	return TaskData{Data: data, FetchTime: t.time}
}

func (t *Task) StartTaskJob(runRightNow bool) {
	log.Printf("task job scheduler start")
	if t.quitTask != nil {
		close(t.quitTask)
	}
	t.taskTicker = time.NewTicker(time.Duration(t.Config.IntervalSec) * time.Second)
	quitTask := make(chan struct{})
	go func() {
		for {
			select {
			case <-t.taskTicker.C:
				t.Do()
			case <-quitTask:
				t.taskTicker.Stop()
			}
		}
	}()

	if runRightNow {
		t.Do()
	}
}

func (t *Task) Stop() {
	close(t.quitTask)
}

func (t *Task) Do() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("do task error", r)
		}
	}()
	log.Println("to do task...")
	resMap := make(map[string]interface{})

	m := miner.FetchData(t.Config.MinerUrl)
	resMap["miner"] = m

	resMap["device"] = t.FetchDeviceInfo()

	t.data = resMap
	t.time = time.Now()

	if link.SingleClientConnected() {
		v, _ := miner.PacketForwardVersion()
		config, _ := link.GetClientConfig()
		var res map[string]interface{}
		res = t.Data().Data
		res["time"] = t.Data().FetchTime
		res["version"] = v
		buf, err := json.Marshal(res)
		if err != nil {
			log.Println("Marshal task data error: ", err)
		}
		client := link.GetSingleClient()
		err = client.WriteMessage(message.OfHttpResponse(
			config.ID+"/hotspotInfoCache",
			message.HttpResponseData{
				Body: string(buf),
			}))
		if err != nil {
			log.Println("WriteMessage error: ", err)
		}
	} else {
		log.Println("ws client not Connected")
	}
	log.Println("task done")
}

func (t *Task) FetchMinerInfo() map[string]interface{} {
	return miner.FetchData(t.Config.MinerUrl)
}

func (t *Task) FetchDeviceInfo() map[string]interface{} {
	resMap := device.GetInfo()

	wifi, err := device.GetWifiInfo()
	if err != nil {
		log.Println("fetch wifi info error", err)
	}
	resMap["wifi"] = wifi

	cpuTemp, _ := device.GetCpuTemp()
	resMap["cpuTemp"] = cpuTemp

	cpuFreq, _ := device.GetCpuFreq()
	resMap["cpuFreq"] = cpuFreq

	return resMap
}

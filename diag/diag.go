package diag

import (
	"encoding/json"
	"time"

	"github.com/kpango/glg"
	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/diag/device"
	"xdt.com/hm-diag/diag/miner"
	"xdt.com/hm-diag/link"
	"xdt.com/hm-diag/link/message"
)

type TaskConfig struct {
	MinerUrl     string
	MinerGrpcUrl string
	IntervalSec  uint
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
	taskSingleton = &Task{Config: TaskConfig{MinerUrl: conf.MinerUrl, IntervalSec: conf.IntervalSec, MinerGrpcUrl: conf.MinerGrpcUrl}}
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
	glg.Debug("task job scheduler start")
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
			glg.Error("do task error", r)
		}
	}()
	glg.Debug("to do task...")
	resMap := make(map[string]interface{})

	m := miner.FetchData(t.Config.MinerUrl, t.Config.MinerGrpcUrl)
	resMap["miner"] = m

	resMap["device"] = t.FetchDeviceInfo()

	t.data = resMap
	t.time = time.Now()

	if link.Connected() {
		v, _ := miner.PacketForwardVersion()
		config, _ := link.GetClientConfig()
		var res map[string]interface{}
		res = t.Data().Data
		res["time"] = t.Data().FetchTime
		res["packetForwardVersion"] = v
		buf, err := json.Marshal(res)
		if err != nil {
			glg.Error("Marshal task data error: ", err)
		}
		err = link.ReportData(message.OfReportData(config.ID+"/hotspotInfoCache", string(buf)))
		if err != nil {
			glg.Error("ReportData error: ", err)
		}
	} else {
		glg.Info("ws client not Connected")
	}
	glg.Debug("task done")
}

func (t *Task) FetchMinerInfo() map[string]interface{} {
	return miner.FetchData(t.Config.MinerUrl, t.Config.MinerGrpcUrl)
}

func (t *Task) FetchDeviceInfo() map[string]interface{} {
	resMap := device.GetInfo()

	wifi, err := device.GetWifiInfo()
	if err != nil {
		glg.Error("fetch wifi info error", err)
	}
	resMap["wifi"] = wifi

	cpuTemp, _ := device.GetCpuTemp()
	resMap["cpuTemp"] = cpuTemp

	cpuFreq, _ := device.GetCpuFreq()
	resMap["cpuFreq"] = cpuFreq

	return resMap
}

package main

import (
	"log"
	"time"

	"xdt.com/hm-diag/hardware"
	"xdt.com/hm-diag/miner"
)

type TaskConfig struct {
	MinerUrl    string
	IntervalSec uint
}

type Task struct {
	Config     TaskConfig
	time       time.Time
	data       interface{}
	taskTicker *time.Ticker
	quitTask   chan struct{}
}

type TaskData struct {
	FetchTime time.Time   `json:"fetch_time"`
	Data      interface{} `json:"data"`
}

func (task *Task) GetData() TaskData {
	return TaskData{Data: task.data, FetchTime: task.time}
}

func (task *Task) StartTask(runRightNow bool) {
	log.Printf("task scheduler start")
	if task.quitTask != nil {
		close(task.quitTask)
	}
	task.taskTicker = time.NewTicker(time.Duration(task.Config.IntervalSec) * time.Second)
	quitTask := make(chan struct{})
	go func() {
		for {
			select {
			case <-task.taskTicker.C:
				task.DoTask()
			case <-quitTask:
				task.taskTicker.Stop()
			}
		}
	}()

	if runRightNow {
		defer func() {
			if r := recover(); r != nil {
				log.Println("do task error", r)
			}
		}()
		task.DoTask()
	}
}

func (task *Task) Stop() {
	close(task.quitTask)
}

func (task *Task) DoTask() {
	log.Println("to do task...  =======================>")
	resMap := make(map[string]interface{})

	m := miner.FetchData(task.Config.MinerUrl)
	resMap["miner"] = m

	resMap["hardware"] = GetHardwareInfo()

	task.data = resMap
	task.time = time.Now()

	log.Println("task done <=======================")
}

func GetHardwareInfo() map[string]interface{} {
	// resMap := make(map[string]interface{})

	resMap := hardware.GetInfo()

	wifi := hardware.GetWifiInfo()
	resMap["wifi"] = wifi

	cpuTemp, _ := hardware.GetCpuTemp()
	resMap["cpu_temp"] = cpuTemp

	cpuFreq, _ := hardware.GetCpuFreq()
	resMap["cpu_freq"] = cpuFreq

	return resMap
}

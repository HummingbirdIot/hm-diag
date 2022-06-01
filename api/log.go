package api

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"xdt.com/hm-diag/diag/device"
)

const (
	logPath = "/tmp/log/"
)

func packLogs() (string, error) {
	logFileName := logPath + "log" + strconv.FormatInt(time.Now().Unix(), 10) + ".tar.gz"
	err := os.MkdirAll(logPath, 0777)
	if err != nil {
		return "", err
	}

	//删除之前的日志文件
	dir, _ := ioutil.ReadDir(logPath)
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{logPath, d.Name()}...))
	}

	// 创建一个缓冲区用来保存压缩文件内容
	var buf bytes.Buffer
	// 创建一个压缩文档
	tw := tar.NewWriter(&buf)

	until := time.Now()
	since := time.Date(until.Year(), until.Month(), until.Day(), 0, 0, 0, 0, until.Location())
	pktfwdLog, err := device.QueryPktfwdLog(since, until, "")
	if err != nil {
		return "", err
	}

	minerLog, err := device.QueryMinerLog("", 5000)
	if err != nil {
		return "", err
	}

	dhcpcd, err := device.QueryDhcpcdLog()
	if err != nil {
		return "", err
	}

	diagLog, err := device.QueryDiagLog()
	if err != nil {
		return "", err
	}

	hiotLog, err := device.QueryHiotLog()
	if err != nil {
		return "", err
	}
	// 定义一堆文件
	// 将文件写入到压缩文档tw
	var files = []struct {
		Name, Body string
	}{
		{"hm-diag.txt", diagLog},
		{"hiot.txt", hiotLog},
		{"packet-forward.txt", pktfwdLog},
		{"miner.txt", minerLog},
		{"dhcpcd.txt", dhcpcd},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatal(err)
	}

	// 将压缩文档内容写入文件 file.tar.gz
	f, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	buf.WriteTo(f)
	return logFileName, nil
}

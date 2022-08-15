package diag

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/kpango/glg"
	"xdt.com/hm-diag/config"
)

const (
	logPath           = "/tmp/log/"
	DIAG_LOG_SCRIPT   = "journalctl -u hm-diag.service -S {since} -U {until} -n 100000 --no-hostname -o cat"
	HIOT_LOG_SCRIPT   = "journalctl -u hiot -S {since} -U {until} -n 100000 --no-hostname -o cat"
	DHCPCD_LOG_SCRIPT = "journalctl -u dhcpcd.service -n 5000 -o cat"
)

func PackLogs() (string, error) {
	var writer *gzip.Writer

	logFileName := logPath + "log" + time.Now().Format("2006-01-02 15:04:05") + ".tar.gz"
	err := os.MkdirAll(logPath, 0777)
	if err != nil {
		return "", err
	}

	//删除之前的日志文件
	dir, _ := ioutil.ReadDir(logPath)
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{logPath, d.Name()}...))
	}

	// 将压缩文档内容写入文件 file.tar.gz
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		glg.Fatal(err)
	}

	if writer, err = gzip.NewWriterLevel(f, gzip.BestCompression); err != nil {
		glg.Fatalln(err)
	}
	defer writer.Close()

	// 创建一个缓冲区用来保存压缩文件内容
	// 创建一个压缩文档
	tw := tar.NewWriter(writer)

	until := time.Now()
	since := until.Add(-time.Hour * 24)

	pktfwdLogChan := make(chan string)
	go func() {
		glg.Debug("srart query pktfwd log")
		pktfwd, err := QueryPktfwdLog(since, until, "")
		if err != nil {
			pktfwdLogChan <- ""
		}
		pktfwdLogChan <- pktfwd
	}()

	minerLogChan := make(chan string)
	go func() {
		glg.Debug("srart query miner log")
		miner, err := QueryMinerLog("", 5000)
		if err != nil {
			minerLogChan <- ""
		}
		minerLogChan <- miner
	}()

	dhcpcdLogChan := make(chan string)
	go func() {
		glg.Debug("srart query dhcpcd log")
		dhcpcd, err := QueryDhcpcdLog()
		if err != nil {
			dhcpcdLogChan <- ""
		}
		dhcpcdLogChan <- dhcpcd
	}()

	diagLogChan := make(chan string)
	go func() {
		glg.Debug("srart query diag log")
		diag, err := QueryDiagLog(since, until)
		if err != nil {
			diagLogChan <- ""
		}
		diagLogChan <- diag
	}()

	hiotLogChan := make(chan string)
	go func() {
		glg.Debug("srart query hiot log")
		hiot, err := QueryHiotLog(since, until)
		if err != nil {
			hiotLogChan <- ""
		}
		hiotLogChan <- hiot
	}()

	diagLog := <-diagLogChan
	hiotLog := <-hiotLogChan
	minerLog := <-minerLogChan
	dhcpcdLog := <-dhcpcdLogChan
	pktfwdLog := <-pktfwdLogChan

	// 定义一堆文件
	// 将文件写入到压缩文档tw
	var files = []struct {
		Name, Body string
	}{
		{"hm-diag.txt", diagLog},
		{"hiot.txt", hiotLog},
		{"packet-forward.txt", pktfwdLog},
		{"miner.txt", minerLog},
		{"dhcpcd.txt", dhcpcdLog},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			glg.Error(err)
			return "", err
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			glg.Error(err)
			return "", err
		}
	}
	if err := tw.Close(); err != nil {
		glg.Error(err)
		return "", err
	}

	return logFileName, nil
}

func QueryDiagLog(since, until time.Time) (string, error) {
	s := since.Format("'2006-01-02 15:04:05'")
	u := until.Format("'2006-01-02 15:04:05'")
	script := strings.ReplaceAll(DIAG_LOG_SCRIPT, "{since}", s)
	script = strings.ReplaceAll(script, "{until}", u)
	glg.Debug("exec cmd: ", script)
	cmd := exec.Command("bash", "-c", script)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func QueryHiotLog(since, until time.Time) (string, error) {
	s := since.Format("'2006-01-02 15:04:05'")
	u := until.Format("'2006-01-02 15:04:05'")
	script := strings.ReplaceAll(HIOT_LOG_SCRIPT, "{since}", s)
	script = strings.ReplaceAll(script, "{until}", u)
	glg.Debug("exec cmd: ", script)
	cmd := exec.Command("bash", "-c", script)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func QueryDhcpcdLog() (string, error) {
	cmd := exec.Command("bash", "-c", DHCPCD_LOG_SCRIPT)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func QueryPktfwdLog(since, until time.Time, filterTxt string) (string, error) {
	queryCmd := config.MAIN_SCRIPT + " pktfwdLog"
	cmdStr := fmt.Sprintf("%s %s %s %s",
		queryCmd,
		since.Format("'2006-01-02 15:04:05'"),
		until.Format("'2006-01-02 15:04:05'"),
		"'"+filterTxt+"'")
	glg.Debug("exec cmd:", cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = config.Config().GitRepoDir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func QueryMinerLog(filterTxt string, maxLines uint) (string, error) {
	queryCmd := config.MAIN_SCRIPT + " minerLog"
	cmdStr := fmt.Sprintf("%s %s %d",
		queryCmd,
		"'"+filterTxt+"'", maxLines)
	glg.Debug("exec cmd:", cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = config.Config().GitRepoDir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

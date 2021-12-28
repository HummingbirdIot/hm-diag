package ctrl

import (
	"bufio"
	"encoding/base64"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	expect "github.com/ThomasRooney/gexpect"
	"github.com/pkg/errors"
	"xdt.com/hm-diag/util"
)

const mainWorkDir = "/home/pi/hnt_iot/"

// const mainWorkDir = "./mock/"

const resyncMinerCmd = mainWorkDir + "trim_miner.sh"
const snapshotTakeCmd = mainWorkDir + "snapshot_take.sh"
const snapshotLoadCmd = mainWorkDir + "snapshot_load.sh"

type SnapshotStateRes struct {
	File  string    `json:"file"`
	State string    `json:"state"`
	Time  time.Time `json:"time"`
}

func ResyncMiner() error {
	log.Println("to resync miner")
	log.Println("exec cmd: bash", resyncMinerCmd)
	cmd := exec.Command("bash", resyncMinerCmd)
	data, err := cmd.Output()
	if err != nil {
		log.Println("[error] resync miner error:", err.Error(), string(data))
		return err
	}
	log.Println("resync miner output:", string(data))
	return nil
}

func SnapshotTake() {
	fn := func() error {
		log.Println("spawn cmd: bash", snapshotTakeCmd, "take")
		cmd := exec.Command("bash", snapshotTakeCmd, "take")
		p, err := cmd.StdoutPipe()
		if err != nil {
			return errors.WithStack(errors.WithMessage(err, "start snapshot cmd pipe error"))
		}
		r := bufio.NewReader(p)

		err = cmd.Start()
		if err != nil {
			return errors.WithStack(errors.WithMessage(err, "start snapshot cmd error"))
		}
		for err == nil {
			ln, _, errIn := r.ReadLine()
			err = errIn
			if err == nil {
				s := string(ln)
				log.Println("snapshot cmd output:", s)
			} else if err == io.EOF {
				break
			} else {
				log.Println("read snapshot cmd ouput error:", err.Error())
			}
		}

		err = cmd.Wait()
		if err != nil {
			return errors.WithStack(errors.WithMessage(err, "snapshot exit error"))
		}
		return nil
	}
	util.Sgo(fn, "snapshot take error")
}

func SnapshotState() (*SnapshotStateRes, error) {
	var result SnapshotStateRes
	resPrefix := ">>>state:"
	log.Println("spawn cmd: bash", snapshotTakeCmd)
	cmd := exec.Command("bash", snapshotTakeCmd, "state")
	p, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "start snapshot state cmd pipe error"))
	}
	r := bufio.NewReader(p)

	err = cmd.Start()
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "start snapshot state cmd error"))
	}
	for err == nil {
		ln, _, errIn := r.ReadLine()
		err = errIn
		if err == nil {
			s := string(ln)
			log.Println("snapshot state cmd output:", s)
			if strings.HasPrefix(s, resPrefix) {
				s = strings.TrimPrefix(s, resPrefix)
				result, err = parseSnapshotStateResult(s)
				if err != nil {
					return nil, err
				}
			}
		} else if err == io.EOF {
			break
		} else {
			log.Println("read snapshot state cmd ouput error:", err.Error())
		}
	}

	err = cmd.Wait()
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "snapshot exit error"))
	}
	return &result, nil
}

func parseSnapshotStateResult(s string) (SnapshotStateRes, error) {
	result := SnapshotStateRes{}
	arr := strings.Split(s, ",")
	for _, a := range arr {
		subArr := strings.Split(a, "=")
		field := subArr[0]
		value := subArr[1]
		if field == "time" && value != "" {
			u, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return result, errors.WithStack(
					errors.WithMessage(err, "snapshot state result format eroror:"+value))
			}
			result.Time = time.Unix(u, 0)
		}
		if field == "file" && value != "" {
			result.File = base64.StdEncoding.EncodeToString([]byte(value))
		}
		if field == "state" {
			result.State = value
		}
	}
	return result, nil
}

func SnapshotLoad(file string) {
	fn := func() error {
		cmd := "bash " + snapshotLoadCmd + " " + file
		log.Println("spawn cmd:", cmd)
		child, err := expect.Spawn(cmd)
		if err != nil {
			return err
		}
		for err == nil {
			ln, err := child.ReadLine()
			if err != nil {
				return errors.WithMessage(err, "load snapshort error when read output")
			}
			log.Println("load snapshot output:", ln)
		}
		err = child.Wait()
		if err != nil {
			err = errors.WithMessage(err, "load snapthot exit error")
			return err
		}
		return nil
	}
	util.Sgo(fn, "snapshot load error")
}

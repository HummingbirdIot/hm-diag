package ctrl

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/kpango/glg"
	"github.com/pkg/errors"
	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/diag/jsonrpc"
	"xdt.com/hm-diag/util"
)

const (
	MakerAddr        = "14DdSjvEkBQ46xQ24LAtHwQkAeoUUZHfGCosgJe33nRQ6rZwPG3"
	resyncMinerCmd   = "./trim_miner.sh"
	snapshotTakeCmd  = "./snapshot_take.sh take"
	snapshotStateCmd = "./snapshot_take.sh state"
	snapshotLoadCmd  = "./snapshot_load.sh"
	restartMinerCmd  = config.MAIN_SCRIPT + " restartMiner"
)

var (
	onboardingCmd = "docker exec hnt_iot_helium-miner_1 /opt/miner/lib/miner-0.1.0/priv/gateway_rs/helium_gateway -c /opt/miner/lib/miner-0.1.0/priv/gateway_rs/ add --owner {owner} --payer {payer} --mode light"
)

type SnapshotStateRes struct {
	File  string    `json:"file"`
	State string    `json:"state"`
	Time  time.Time `json:"time"`
}

type onboardingCliResp struct {
	Address     string `json:"address"`
	Fee         int64  `json:"fee"`
	Mode        string `json:"mode"`
	Owner       string `json:"owner"`
	Payer       string `json:"payer"`
	Staking_fee int64  `json:"staking fee"`
	Txn         string `json:"txn"`
}

func ResyncMiner() error {
	glg.Info("to resync miner")
	glg.Debug("exec cmd: bash", resyncMinerCmd)
	cmd := exec.Command("bash", resyncMinerCmd)
	cmd.Dir = config.Config().GitRepoDir
	data, err := cmd.Output()
	if err != nil {
		glg.Error("[error] resync miner error:", err.Error(), string(data))
		return err
	}
	glg.Info("resync miner output:", string(data))
	return nil
}

func RestartMiner() error {
	glg.Info("to restart miner")
	glg.Debug("exec cmd:", restartMinerCmd)
	cmd := exec.Command("bash", strings.Split(restartMinerCmd, " ")...)
	cmd.Dir = config.Config().GitRepoDir
	data, err := cmd.Output()
	if err != nil {
		glg.Error("[error] restart miner error:", err.Error(), string(data))
		return err
	}
	glg.Info("restart miner output:", string(data))
	return nil
}

func SnapshotTake() {
	fn := func() error {
		glg.Debug("spawn cmd: bash ", snapshotTakeCmd)
		cmd := exec.Command("bash", strings.Split(snapshotTakeCmd, " ")...)
		cmd.Dir = config.Config().GitRepoDir
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
				glg.Debug("snapshot cmd output:", s)
			} else if err == io.EOF {
				break
			} else {
				glg.Error("read snapshot cmd ouput error:", err.Error())
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
	glg.Debug("spawn cmd:", snapshotStateCmd)
	cmd := exec.Command("bash", strings.Split(snapshotStateCmd, " ")...)
	cmd.Dir = config.Config().GitRepoDir
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
			glg.Debug("snapshot state cmd output:", s)
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
			glg.Error("read snapshot state cmd ouput error:", err.Error())
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
		glg.Debug("exec cmd: bash " + snapshotLoadCmd + " " + file)
		cmd := exec.Command("bash", snapshotLoadCmd, file)
		cmd.Dir = config.Config().GitRepoDir
		p, err := cmd.StdoutPipe()
		if err != nil {
			return errors.WithStack(errors.WithMessage(err, "start snapshot load cmd pipe error"))
		}
		r := bufio.NewReader(p)

		err = cmd.Start()
		if err != nil {
			return errors.WithStack(errors.WithMessage(err, "start snapshot load cmd error"))
		}
		for err == nil {
			ln, _, errIn := r.ReadLine()
			err = errIn
			if err == nil {
				s := string(ln)
				glg.Info("snapshot load cmd output:", s)
			} else if err == io.EOF {
				break
			} else {
				glg.Error("read snapshot load cmd ouput error:", err.Error())
			}
		}

		err = cmd.Wait()
		if err != nil {
			return errors.WithStack(errors.WithMessage(err, "snapshot load exit error"))
		}
		return nil
	}
	util.Sgo(fn, "snapshot load error")
}

func GenOnboardingTxn(ownerAddr string, payerAddr string) (string, error) {
	if ownerAddr == "" || payerAddr == "" {
		return "", fmt.Errorf("ownerAddr and payerAddr must be provided")
	}
	glg.Info(ownerAddr, payerAddr)
	var out bytes.Buffer
	onboardingCmd = strings.Replace(onboardingCmd, "{owner}", ownerAddr, 1)
	onboardingCmd = strings.Replace(onboardingCmd, "{payer}", payerAddr, 1)
	cmd := exec.Command("bash", "-c", onboardingCmd)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	var resp onboardingCliResp
	err = json.NewDecoder(&out).Decode(&resp)
	if err != nil {
		return "", err
	}
	glg.Infof("%+v", resp)
	return resp.Txn, nil
	// jrClient := jsonrpc.Client{Url: config.Config().MinerUrl}

	// re, err := jrClient.Call("txn_add_gateway", map[string]string{
	// 	"owner": ownerAddr,
	// 	"payer": payerAddr,
	// })
	// if err != nil {
	// 	return "", err
	// }
	// result, ok := re.(map[string]interface{})
	// if !ok {
	// 	return "", fmt.Errorf("miner txn add gateway request result error: %#v", re)
	// }
	// txn, ok := result["result"].(string)
	// if !ok {
	// 	return "", fmt.Errorf("miner txn add gateway request result error: %#v", result)
	// }
	// return txn, nil
}

func GenAssertLocationTxn(ownerAddr, payerAddr, location string, nonce int) (string, error) {
	// TODO: is free ? payerAddr = ownerAddr ?
	if ownerAddr == "" || payerAddr == "" || location == "" {
		return "", fmt.Errorf("ownerAddr, location and payerAddr must be provided")
	}
	jrClient := jsonrpc.Client{Url: config.Config().MinerUrl}

	re, err := jrClient.Call("txn_assert_location", map[string]string{
		"owner": ownerAddr,
		"payer": payerAddr,
		"h3":    location,
	})
	if err != nil {
		return "", err
	}
	result, ok := re.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("miner txn assert location request result error: %#v", re)
	}
	txn, ok := result["result"].(string)
	if !ok {
		return "", fmt.Errorf("miner txn assert location request result error: %#v", result)
	}
	return txn, nil
}

type HeliumHotspot struct {
	SpeculativeNonce int     `json:"speculative_nonce"`
	Lng              float64 `json:"lng"`
	Lat              float64 `json:"lat"`
	TimestampAdded   string  `json:"timestamp_added"`
	// more field ...
}
type HeliumHotspotResp struct {
	Data HeliumHotspot `json:"data"`
}

type OnboardingRecord struct {
	OnboardingKey string           `json:"onboarding_key"`
	PublicAddress string           `json:"public_address"`
	Maker         *OnboardingMaker `json:"maker"`
}

type OnboardingMaker struct {
	id                 int    `json:"id"`
	Address            string `json:"address"`
	LocationNonceLimit int    `json:"locationNonceLimit"`
}

func isFreeAssertLocation(hotspotAddr string) (bool, error) {
	url := "https://api.helium.io/v1/hotspots/" + hotspotAddr
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	var hotspot HeliumHotspotResp
	err = json.NewDecoder(resp.Body).Decode(&hotspot)
	if err != nil {
		return false, err
	}

	url = "https://onboarding.dewi.org/api/v2/hotspots/" + hotspotAddr
	resp, err = http.Get(url)
	if err != nil {
		return false, err
	}

	var onboardingRecord OnboardingRecord
	err = json.NewDecoder(resp.Body).Decode(&onboardingRecord)
	if err != nil {
		return false, err
	}

	isFree := false
	makerLimit := 0
	if onboardingRecord.Maker != nil {
		makerLimit = onboardingRecord.Maker.LocationNonceLimit
	}
	isFree = hotspot.Data.SpeculativeNonce < makerLimit
	return isFree, nil
}

package onboarding

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/diag/miner"
)

const (
	defaultMakerName   = "hummingbird"
	ONBOARDING_API_URL = "https://onboarding.dewi.org"
)

type onboardingResp struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Success bool        `json:"success"`
}

func CheckOnboarding() bool {
	addr := miner.GetPeerAddr(config.Config().MinerUrl)
	pubkey := strings.Split(addr, "/")[2]
	api := ONBOARDING_API_URL + "/api/v2/hotspots/" + pubkey
	log.Println(api)
	resp, err := http.Get(api)
	if err != nil {
		log.Println("onboarding API error,", err)
		return false
	}
	if resp.StatusCode >= 400 {
		log.Printf("onboarding API error, status code: %d", resp.StatusCode)
		return false
	}
	var respData onboardingResp
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		log.Println("read onboarding API body error,", err)
		return false
	}

	if respData.Code != 200 || !respData.Success {
		log.Println("onboarding API code not 200 or success is false")
		return false
	}

	maker, ok := respData.Data.(map[string]interface{})["maker"].(map[string]interface{})
	if !ok {
		log.Println("get onboarding API maker error")
		return false
	}

	makerName, ok := maker["name"].(string)
	if !ok {
		log.Println("get onboarding API makerName error")
		return false
	}

	if makerName != defaultMakerName {
		log.Printf("makername value error, want :%s,have : %s", defaultMakerName, makerName)
		return false
	}

	return true
}

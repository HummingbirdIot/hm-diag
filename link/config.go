package link

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/diag/miner"
)

const (
	configFilePath = config.ETC_DIR + "/linkClient.json"
)

var (
	clientConfig *ClientConfig
)

type ClientConfig struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Server string `json:"server"`
}

func InitClientConfig() {
	c, err := loadClientConfig()
	log.Println("client config : ", c)
	if err != nil {
		log.Println("load client config error", err)
	} else {
		clientConfig = c
	}
}

func GetClientConfig() (*ClientConfig, error) {
	if clientConfig == nil {
		return nil, fmt.Errorf("client config not ready yet")
	}
	return clientConfig, nil
}

func SaveClientConfig(c ClientConfig) error {
	buf, err := json.Marshal(c)
	defaultErr := errors.WithMessage(err, "marshal client config")
	if err != nil {
		return defaultErr
	}
	err = os.WriteFile(configFilePath, buf, 0664)
	if err != nil {
		return defaultErr
	}
	var conf ClientConfig
	json.Unmarshal(buf, &conf)
	log.Println("set client config cache, content : ", conf)
	clientConfig = &conf
	return nil
}

func loadClientConfig() (*ClientConfig, error) {
	f, err := os.Open(configFilePath)
	defaultErr := errors.WithMessage(err, "load client config")
	if err != nil {
		return nil, defaultErr

	}
	var conf ClientConfig
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		errors.WithMessage(err, "client config content format")
	}

	if conf.ID == "" {
		log.Println("client config file not have ID")
		m := miner.FetchData(config.Config().MinerUrl)
		if addr, ok := m["peerAddr"].(string); ok {
			conf.ID = strings.Split(addr, "/")[2]
		} else {
			conf.ID = uuid.NewString()
		}
		log.Println("set client config file not have ID: ", conf.ID)
	}

	return &conf, nil
}

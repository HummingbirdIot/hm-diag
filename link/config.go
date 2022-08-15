package link

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kpango/glg"
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
	c, err := LoadClientConfig()
	if err != nil {
		glg.Error("load client config error", err)
	} else {
		glg.Info("load clieng config: ", c)
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
	glg.Info("set client config cache, content : ", conf)
	clientConfig = &conf
	return nil
}

func LoadClientConfig() (*ClientConfig, error) {
	f, err := os.Open(configFilePath)
	defaultErr := errors.WithMessage(err, "load client config")
	if err != nil {
		if os.IsNotExist(err) {
			//配置文件不存在
			glg.Warn("File config not exist")
		} else {
			return nil, defaultErr
		}
	}
	var conf ClientConfig
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		errors.WithMessage(err, "client config content format")
	}

	if conf.ID == "" {
		glg.Info("client config file not have ID")
		addr := miner.GetPeerAddr(config.Config().MinerUrl)
		if addr != "" {
			conf.ID = strings.Split(addr, "/")[2]
		} else {
			return nil, fmt.Errorf("hotspot peerAddr is empty")
		}
		glg.Info("set client config file ID: ", conf.ID)
	}

	return &conf, nil
}

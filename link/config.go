package link

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"xdt.com/hm-diag/config"
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

func init() {
	c, err := loadClientConfig()
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

	return &conf, nil
}

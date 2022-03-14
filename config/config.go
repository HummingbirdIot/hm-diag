package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type CONF_BOOL int

const (
	CONF_DEFAULT CONF_BOOL = 0
	CONF_ON      CONF_BOOL = 1
	CONF_OFF     CONF_BOOL = 2
)

type GlobalConfig struct {
	ApiPort       int
	LanDevIntface string
	MinerUrl      string
	IntervalSec   uint
	GitRepoDir    string
	GitRepoUrl    string
	PublicAccess  CONF_BOOL
}

type ConfiFileData struct {
	PublicAccess CONF_BOOL `json:"publicAccess"`
}

const MAIN_SCRIPT = "./hummingbird_iot.sh"
const PROXY_ETC_DIR = "/usr/local/etc/hm-diag"
const PROXY_ETC_REPO = PROXY_ETC_DIR + "/git-repo-proxy.json"
const PROXY_ETC_RELEASE = PROXY_ETC_DIR + "/git-release-proxy.json"
const CONF_ETC_FILE = PROXY_ETC_DIR + "/config.json"
const GITHUB_URL = "https://github.com/"

var conf *GlobalConfig

func InitConf(cf GlobalConfig) {
	conf = &cf
	if cf.PublicAccess == CONF_DEFAULT {
		confFile, err := ReadConfigFile()
		log.Printf("read config file content: %#v", confFile)
		if err != nil {
			log.Println(err)
			cf.PublicAccess = CONF_OFF
		} else {
			conf.PublicAccess = confFile.PublicAccess
		}
	}
}

func Config() *GlobalConfig {
	if conf == nil {
		panic("GlobalConfig is nil")
	}
	return conf
}

func ReadConfigFile() (*ConfiFileData, error) {
	f, err := os.Open(CONF_ETC_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("conf file is not exist")
			return &ConfiFileData{PublicAccess: CONF_DEFAULT}, nil
		}
		return nil, fmt.Errorf("can't open config file %s", CONF_ETC_FILE)
	}
	var conf ConfiFileData
	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		return nil, fmt.Errorf("can't parse config read from config file %s error: %s", CONF_ETC_FILE, err)
	}
	return &conf, nil
}

func SaveConfigFile(conf ConfiFileData) error {
	buf, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	err = os.WriteFile(CONF_ETC_FILE, buf, 0644)
	if err != nil {
		return err
	}
	// update live config
	Config().PublicAccess = conf.PublicAccess
	return nil
}

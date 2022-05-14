package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
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
	Password      string
}

type ConfiFileData struct {
	PublicAccess CONF_BOOL `json:"publicAccess"`
	Password     string    `json:"password"`
}

const MAIN_SCRIPT = "./hummingbird_iot.sh"
const ETC_DIR = "/usr/local/etc/hm-diag"
const PROXY_ETC_REPO = ETC_DIR + "/git-repo-proxy.json"
const PROXY_ETC_RELEASE = ETC_DIR + "/git-release-proxy.json"
const CONF_ETC_FILE = ETC_DIR + "/config.json"
const GITHUB_URL = "https://github.com/"

var conf *GlobalConfig

func InitConf(cf GlobalConfig) {
	conf = &cf
	loadConfigFile(conf)
}

func Config() *GlobalConfig {
	if conf == nil {
		panic("GlobalConfig is nil")
	}
	return conf
}

func loadConfigFile(cf *GlobalConfig) {
	if cf.PublicAccess == CONF_DEFAULT || cf.PublicAccess == CONF_OFF {
		confFile, err := ReadConfigFile()
		log.Printf("read config file content: %#v", confFile)
		if err != nil {
			log.Println(err)
			cf.PublicAccess = CONF_OFF
		} else {
			conf.PublicAccess = confFile.PublicAccess
		}

		if confFile == nil || confFile.Password == "" {
			inet, err := net.InterfaceByName("eth0")
			if err != nil {
				conf.Password = "Hiot@2022"
				fmt.Printf("fail to get net interfaces: %v", err)
			} else {
				conf.Password = inet.HardwareAddr.String()
			}
		} else {
			log.Println("set hotspot password from config.json , password: ", confFile.Password)
			conf.Password = confFile.Password
		}
		log.Println("hotspot password: ", conf.Password)
	}
}

func ReadConfigFile() (*ConfiFileData, error) {
	f, err := os.Open(CONF_ETC_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("conf file is not exist")
			return &ConfiFileData{PublicAccess: CONF_ON}, nil
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

func SaveConfigFile(c ConfiFileData) error {
	buf, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(CONF_ETC_FILE, buf, 0644)
	if err != nil {
		return err
	}
	// update live config
	Config().PublicAccess = c.PublicAccess
	Config().Password = c.Password
	loadConfigFile(conf)
	return nil
}

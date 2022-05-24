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
	ApiPort           int
	LanDevIntface     string
	MinerUrl          string
	IntervalSec       uint
	GitRepoDir        string
	GitRepoUrl        string
	PublicAccess      CONF_BOOL
	Password          string
	DashboardPassword bool
}

type ConfiFileData struct {
	PublicAccess      CONF_BOOL `json:"publicAccess"`
	DashboardPassword bool      `json:"dashboardPassword"`
	Password          string    `json:"password"`
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
	if cf.PublicAccess == CONF_DEFAULT {
		confFile, err := ReadConfigFile()
		log.Printf("read config file content: %#v", confFile)
		if err != nil {
			log.Println(err)
			cf.PublicAccess = CONF_OFF
			cf.DashboardPassword = false
		} else {
			conf.PublicAccess = confFile.PublicAccess
			conf.DashboardPassword = confFile.DashboardPassword
		}
		checkPassword(confFile)
	}
}

func checkPassword(confFile *ConfiFileData) {
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
	Config().Password = conf.Password
	Config().DashboardPassword = conf.DashboardPassword
	checkPassword(&conf)
	return nil
}

package config

import (
	"log"
	"os"
	"testing"
)

var defaultConf = GlobalConfig{
	ApiPort:       8090,
	LanDevIntface: "eth0",
	MinerUrl:      "http://127.0.0.1:4467",
	GitRepoDir:    "/home/pi/hnt_iot",
	GitRepoUrl:    "https://github.com/HummingbirdIot/hnt_iot_release.git",
	IntervalSec:   30,
}

var configFile = ConfiFileData{
	CONF_DEFAULT, false, "123456",
}

func setup() {
	InitConf(defaultConf)
	log.Println("init config")
}

func TestSaveConfigFile(t *testing.T) {
	err := SaveConfigFile(configFile)
	if err != nil {
		t.Error(err)
	}
}

func TestReadConfigFile(t *testing.T) {
	conf, err := ReadConfigFile()
	if err != nil {
		t.Error(err)
	}
	if conf.DashboardPassword != configFile.DashboardPassword || conf.Password != configFile.Password || conf.PublicAccess != configFile.PublicAccess {
		t.Error("invalid config")
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

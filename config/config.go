package config

type GlobalConfig struct {
	MinerUrl    string
	IntervalSec uint
	GitRepoDir  string
	GitRepoUrl  string
}

const PROXY_ETC_DIR = "/usr/local/etc/hm-diag"
const PROXY_ETC_REPO = PROXY_ETC_DIR + "/git-repo-proxy.json"
const PROXY_ETC_RELEASE = PROXY_ETC_DIR + "/git-release-proxy.json"
const GITHUB_URL = "https://github.com/"

var conf *GlobalConfig

func InitConf(cf GlobalConfig) {
	conf = &cf
}

func Config() *GlobalConfig {
	if conf == nil {
		panic("GlobalConfig is nil")
	}
	return conf
}

package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"strings"
)

const proxyConfFile = ".proxyconf"

type ProxyType string

const (
	MIRROR     ProxyType = "mirror"
	URL_PREFIX ProxyType = "urlPrefix"
)

type ProxyItem struct {
	Type  ProxyType `json:"type"`
	Value string    `json:"value"`
}
type ProxyConf struct {
	ReleaseFileProxy ProxyItem `json:"releaseFileProxy"`
}

func SetRepoMirrorProxy(repoDir string, proxy ProxyItem) error {
	// check params
	if proxy.Type != MIRROR {
		msg := "can't set proxy type " + string(proxy.Type) + " other than type 'mirror'"
		log.Println(msg)
		return fmt.Errorf(msg)
	}

	err := checkProxyUrl(proxy.Value)
	if err != nil {
		return err
	}

	// change dir to repo dir
	err = os.Chdir(repoDir)
	if err != nil {
		return err
	}

	// get old url insteadof
	str, err := gitRepoMirrorConf(repoDir)
	if err != nil {
		return err
	}
	log.Println("git config list url insteadof output:\n", str)

	if str != "" {
		strArr := strings.Split(str, "=")
		pre := strArr[0]

		// unset url insteadof
		script := "git config --unset " + pre
		cmd := exec.Command("bash", "-c", script)
		buf, err := cmd.Output()
		if err != nil {
			return err
		}
		str = string(buf)
		log.Println("git config unset url insteadof output:\n", str)
	}

	// set url insteadof
	if proxy.Value == "" {
		return nil
	}
	script := "git config url." + proxy.Value + ".insteadof https://github.com"
	cmd := exec.Command("bash", "-c", script)
	buf, err := cmd.Output()
	if err != nil {
		return err
	}
	str = string(buf)
	log.Println("git config set url insteadof output:\n", str)

	return nil
}

func SetReleaseFileProxy(repoDir string, proxy ProxyItem) error {
	err := checkProxyUrl(proxy.Value)
	if err != nil {
		return err
	}
	err = os.Chdir(repoDir)
	if err != nil {
		return err
	}
	proxyConf := ProxyConf{ReleaseFileProxy: proxy}
	proxyConfBuf, err := json.MarshalIndent(proxyConf, "", "  ")
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	confFilePath := path.Join(wd, proxyConfFile)
	err = os.WriteFile(confFilePath, proxyConfBuf, 0664)
	if err != nil {
		return err
	}
	log.Println("writen config into file:", confFilePath)
	return nil
}

func RepoProxy(repoDir string) (*ProxyItem, error) {
	s, err := RepoMirrorUrl(repoDir)
	if err != nil {
		return nil, err
	} else {
		if s == "" {
			return nil, nil
		}
		return &ProxyItem{Type: MIRROR, Value: s}, nil
	}
}

func RepoMirrorUrl(repoDir string) (string, error) {
	err := os.Chdir(repoDir)
	if err != nil {
		return "", err
	}
	s, err := gitRepoMirrorConf(repoDir)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`url\.(.*)\.insteadof.*`)
	f := re.FindStringSubmatch(s)
	return f[1], nil
}

func ReleaseFileProxy(repoDir string) (*ProxyItem, error) {
	err := os.Chdir(repoDir)
	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	confFilePath := path.Join(wd, proxyConfFile)
	if _, err := os.Stat(confFilePath); err == nil {
		log.Println("read file.........")
		buf, err := os.ReadFile(confFilePath)
		if err != nil {
			return nil, err
		}
		var item ProxyConf
		err = json.Unmarshal(buf, &item)
		log.Println("read file.........", string(buf))
		if err != nil {
			return nil, err
		}
		return &item.ReleaseFileProxy, nil
	} else {
		log.Println("read file.........err", err)
		return nil, nil
	}
}

func checkProxyUrl(urlStr string) error {
	if urlStr == "" {
		return nil
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	if url.Scheme != "http" && url.Scheme != "https" {
		return fmt.Errorf("proxy scheme should be https or http")
	}
	return nil
}

func gitRepoMirrorConf(repoDir string) (string, error) {
	err := os.Chdir(repoDir)
	if err != nil {
		return "", err
	}
	script := "git config --list | grep -e \"url\\..*insteadof=.*\""
	cmd := exec.Command("bash", "-c", script)
	buf, err := cmd.Output()
	str := string(buf)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			// ignore
		} else {
			log.Println("git config list url insteadof error:", reflect.TypeOf(err))
			return "", err
		}
	}
	return str, nil
}

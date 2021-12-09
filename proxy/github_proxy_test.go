package proxy

import (
	"testing"
)

const gitDir = "/home/pi/hnt_iot_release"

func TestSetRepoMirror(t *testing.T) {
	err := SetRepoMirrorProxy(
		gitDir,
		ProxyItem{Type: MIRROR, Value: "https://hub.fastgit.org"},
	)
	if err != nil {
		t.Error(err)
	}
}

func TestSetReleaseFileProxy(t *testing.T) {
	err := SetReleaseFileProxy(
		gitDir,
		ProxyItem{Type: URL_PREFIX, Value: "https://ghproxy.com/"},
	)
	if err != nil {
		t.Error(err)
	}
}

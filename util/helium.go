package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const HELIUM_API_URL = "https://api.helium.io/"

func HeliumApiProxy(path string) (string, error) {
	api := HELIUM_API_URL + path
	resp, err := http.Get(api)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("helium API error, status code: %d", resp.StatusCode)
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

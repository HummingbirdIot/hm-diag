package link

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/kpango/glg"
	"github.com/pkg/errors"
	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/link/message"
)

func requestLocal(mr *message.HttpRequest) (*message.HttpResponse, error) {
	rd := mr.Data
	// for security reasons, only support request http api on the server
	urlStr := LocalHost() + rd.URL
	glg.Infof("rpc to %v", urlStr)
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid url "+rd.URL)
	}
	h := http.Header{}
	for k, v := range rd.Header {
		h.Set(k, v)
	}
	r := &http.Request{
		Method: rd.Method,
		URL:    u,
		Header: h,
		Body:   io.NopCloser(strings.NewReader(rd.Body)),
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	resData := message.OfHttpResponse(
		mr.Meta.MsgId,
		message.HttpResponseData{
			StatusCode:    resp.StatusCode,
			Header:        resp.Header,
			ContentLength: resp.ContentLength,
			Body:          string(body),
		})

	return resData, err
}

func LocalHost() string {
	return fmt.Sprintf("http://127.0.0.1:%d", config.Config().ApiPort)
}

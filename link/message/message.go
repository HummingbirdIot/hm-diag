package message

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const DEFAULT_CONTENT_TYPE = "application/json"

// put the message type value in message header TYPE_HEADER
type TypeHeaderValue string

const TYPE_HEADER = "type"
const (
	TYPE_HTTP_REQUEST  = "HttpRequest"
	TYPE_HTTP_RESPONSE = "HttpResponse"
	TYPE_REPORT_DATA   = "ReportData"
)

// put the message type value in message header DATA_TYPE_HEADER for type TYPE_REPORT_DATA
type DataTypeHeaderValue string

const DATA_TYPE_HEADER = "DataType"
const (
	D_TYPE_STATE = "State"
)

type Meta struct {
	MsgId       string            `json:"msgId"`
	ContentType string            `json:"contentType"` // for future
	Header      map[string]string `json:"header"`
}

type ReportData struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

// type for http over ws

type HttpRequest struct {
	Meta Meta            `json:"meta"`
	Data HttpRequestData `json:"data"`
}
type HttpResponse struct {
	Meta Meta             `json:"meta"`
	Data HttpResponseData `json:"data"`
}

type HttpRequestData struct {
	Method        string            `json:"method"`
	URL           string            `json:"url"` // path and query parameters
	Header        map[string]string `json:"headers"`
	ContentLength int64             `json:"contentLength"`
	Body          string            `json:"body"`
}

type HttpResponseData struct {
	StatusCode    int         `json:"statusCode"`
	Header        http.Header `json:"header"`
	ContentLength int64       `json:"contentLength"`
	Body          string      `json:"body"`
}

func OfReportData(d any) *ReportData {
	return &ReportData{
		Meta: Meta{
			MsgId:       uuid.NewString(),
			ContentType: DEFAULT_CONTENT_TYPE,
			Header: map[string]string{
				TYPE_HEADER:      TYPE_REPORT_DATA,
				DATA_TYPE_HEADER: D_TYPE_STATE,
			},
		},
		Data: d,
	}
}

func OfHttpRequest(msgId string, d HttpRequestData) *HttpRequest {
	return &HttpRequest{
		Meta: Meta{
			MsgId:       msgId,
			ContentType: DEFAULT_CONTENT_TYPE,
			Header: map[string]string{
				TYPE_HEADER: TYPE_HTTP_REQUEST,
			},
		},
		Data: d,
	}
}

func OfHttpResponse(msgId string, d HttpResponseData) *HttpResponse {
	return &HttpResponse{
		Meta: Meta{
			MsgId:       msgId,
			ContentType: DEFAULT_CONTENT_TYPE,
			Header: map[string]string{
				TYPE_HEADER: TYPE_HTTP_REQUEST,
			},
		},
		Data: d,
	}
}

func Typeof(msg *map[string]interface{}) (t string, err error) {
	t, err = headerOf(msg, TYPE_HEADER)
	return
}

func DataTypeof(msg *map[string]interface{}) (t string, err error) {
	t, err = headerOf(msg, DATA_TYPE_HEADER)
	return
}

func headerOf(msg *map[string]interface{}, name string) (string, error) {
	m := (*msg)["meta"]
	if meta, ok := m.(map[string]interface{}); ok {
		if header, ok := meta["header"].(map[string]interface{}); ok {
			if t, ok := header[name].(string); ok {
				return t, nil
			}
		}
	}
	return "", fmt.Errorf("can't get type of message")
}

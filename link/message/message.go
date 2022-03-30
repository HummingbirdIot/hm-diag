package message

import (
	"net/http"

	"github.com/google/uuid"
)

type Meta struct {
	MsgId       string            `json:"msgId"`
	ContentType string            `json:"contentType"` // for future
	Header      map[string]string `json:"header"`
}

const DEFAULT_CONTENT_TYPE = "application/json"

type Msg[T any] struct {
	Meta Meta `json:"meta"`
	Data T    `json:"data"`
}

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

// type for http over ws

type HttpRequest Msg[HttpRequestData]
type HttpResponse Msg[HttpResponseData]

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

func OfReportData(d any) *Msg[any] {
	return &Msg[any]{
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

func OfHttpRequest(msgId string, d HttpRequestData) *Msg[HttpRequestData] {
	return &Msg[HttpRequestData]{
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

func OfHttpResponse(msgId string, d HttpResponseData) *Msg[HttpResponseData] {
	return &Msg[HttpResponseData]{
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

func Typeof(msg *Msg[any]) string {
	t := msg.Meta.Header[TYPE_HEADER]
	return t
}

func DataTypeof(msg *Msg[any]) string {
	t := msg.Meta.Header[DATA_TYPE_HEADER]
	return t
}

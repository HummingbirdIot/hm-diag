package api

import (
	"time"

	"xdt.com/hm-diag/ctrl"
	"xdt.com/hm-diag/devdis"
	"xdt.com/hm-diag/diag"
	"xdt.com/hm-diag/diag/device"
	"xdt.com/hm-diag/diag/miner"
)

// swagger:response AllState
type AllState struct {
	// in: body
	Body struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    diag.AllStateInfo `json:"data"`
	}
}

// swagger:response MinerInfo
type MinerInfo struct {
	// in: body
	Body struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    miner.MinerInfo `json:"data"`
	}
}

// swagger:response DeviceInfo
type DeviceInfo struct {
	// in: body
	Body struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Data    device.DeviceInfo `json:"data"`
	}
}

// swagger:parameters all-state miner-state device-state
type StateParams struct {
	// in: query
	// default:false
	Cache bool `json:"cache"`
}

// swagger:response DevDis
type DevDis struct {
	// in: body
	Body struct {
		Code    int        `json:"code"`
		Message string     `json:"message"`
		Dev     devdis.Dev `json:"data"`
	}
}

// swagger:response SnapshotStateRes
type SnapshotStateRes struct {
	// in: body
	Body struct {
		Code    int                   `json:"code"`
		Message string                `json:"message"`
		Data    ctrl.SnapshotStateRes `json:"data"`
	}
}

// swagger:response ProxyItem
type ProxyItem struct {
	// in: body
	Body struct {
		Code    int            `json:"code"`
		Message string         `json:"message"`
		Data    ctrl.ProxyItem `json:"data"`
	}
}

// swagger:parameters config-proxy-update
type ProxyItemUpdate struct {
	// in: query
	// required: true
	Item string `json:"item"`
	// in: body
	Data ctrl.ProxyItem `json:"data"`
}

// swagger:parameters config-proxy-get
type ProxyItemGet struct {
	// in: query
	// required: true
	Item string `json:"item"`
}

// swagger:parameters miner-log
type MinerLogParams struct {
	// format: yyyy-MM-ddTHH:mm:ss
	// in: query
	// required: true
	Since time.Time `json:"item"`
	// format: yyyy-MM-ddTHH:mm:ss
	// in: query
	// required: true
	Until time.Time `json:"until"`
	// in: query
	Filter string `json:"filter"`

	// in: query
	Limit uint `json:"limit"`
}

// swagger:parameters miner-snapshot-download
type MinerSnapshotDownloadParams struct {
	// in:path
	Name string `json:"name"`
}

// swagger:parameters proxy-heliumApi
type HeliumApiProxyParams struct {
	// Helim API path
	// in:query
	// required: true
	Path string `json:"path"`
}

// swagger:response EmptyBody
type EmptyBody struct {
	// in: body
	Body struct {
		Code    int      `json:"code"`
		Message string   `json:"message"`
		Data    struct{} `json:"data"`
	}
}

// swagger:response StringBody
type StringBody struct {
	// in: body
	Body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}
}

// swagger:response BoolBody
type BoolBody struct {
	// in: body
	Body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    bool   `json:"data"`
	} `json:"data"`
}

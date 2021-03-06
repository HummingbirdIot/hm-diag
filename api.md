


# Hm-diag API
API is for： 
- retrieving information about hotspots 
- genernal operation and maintenance

Unified return data structure:
```
{
  "code": 0, 
  "data": {}  
  "message": “OK”
}
```
Unless otherwise specified, the `code` value is the same as the http status.
  

## Informations

### Version

1.0.0

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json

### Produces
  * application/json

## All endpoints

###  public

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /api/v1/config/proxy | [config proxy get](#config-proxy-get) | Get proxy config |
| POST | /api/v1/config/proxy | [config proxy update](#config-proxy-update) | Set proxy config |
| POST | /api/v1/device/reboot | [device reboot](#device-reboot) | Reboot Device |
| GET | /api/v1/device/state | [device state](#device-state) | Get device info |
| GET | /api/v1/lan/hotspot | [lan hotspots](#lan-hotspots) | Get devices(hotspots) address in LAN |
| POST | /api/v1/miner/restart | [miner restart](#miner-restart) | Restart miner |
| POST | /api/v1/miner/resync | [miner resync](#miner-resync) | Resync miner |
| GET | /api/v1/miner/state | [miner state](#miner-state) | Get miner info |
| POST | /api/v1/workspace/update | [workspace update](#workspace-update) | Update worksapce (main git repo) |
| GET | /api/v1/workspace/update | [workspace update get](#workspace-update-get) | Check workspace update |
  


## Paths

### <span id="config-proxy-get"></span> Get proxy config (*config-proxy-get*)

```
GET /api/v1/config/proxy
```

Proxy config is about git repository or git release files
`item` query parameter shoulbe: "gitRelease" or "gitRepo"

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| item | `query` | string | `string` |  | ✓ |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#config-proxy-get-200) | OK |  |  | [schema](#config-proxy-get-200-schema) |

#### Responses


##### <span id="config-proxy-get-200"></span> 200
Status: OK

###### <span id="config-proxy-get-200-schema"></span> Schema
   
  

[ConfigProxyGetOKBody](#config-proxy-get-o-k-body)

###### Inlined models

**<span id="config-proxy-get-o-k-body"></span> ConfigProxyGetOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | |  |  |
| Message | string| `string` |  | |  |  |
| data | [ProxyItem](#proxy-item)| `models.ProxyItem` |  | |  |  |



### <span id="config-proxy-update"></span> Set proxy config (*config-proxy-update*)

```
POST /api/v1/config/proxy
```

roxy config is about git repository or git release files
`item` query parameter shoulbe: "gitRelease" or "gitRepo"

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| item | `query` | string | `string` |  | ✓ |  |  |
| data | `body` | [ProxyItem](#proxy-item) | `models.ProxyItem` | |  | |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#config-proxy-update-200) | OK |  |  | [schema](#config-proxy-update-200-schema) |

#### Responses


##### <span id="config-proxy-update-200"></span> 200
Status: OK

###### <span id="config-proxy-update-200-schema"></span> Schema
   
  

[ConfigProxyUpdateOKBody](#config-proxy-update-o-k-body)

###### Inlined models

**<span id="config-proxy-update-o-k-body"></span> ConfigProxyUpdateOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | |  |  |
| Data | [interface{}](#interface)| `interface{}` |  | |  |  |
| Message | string| `string` |  | |  |  |



### <span id="device-reboot"></span> Reboot Device (*device-reboot*)

```
POST /api/v1/device/reboot
```

API will return immediately, you can check

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#device-reboot-200) | OK |  |  | [schema](#device-reboot-200-schema) |

#### Responses


##### <span id="device-reboot-200"></span> 200
Status: OK

###### <span id="device-reboot-200-schema"></span> Schema
   
  

[DeviceRebootOKBody](#device-reboot-o-k-body)

###### Inlined models

**<span id="device-reboot-o-k-body"></span> DeviceRebootOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | |  |  |
| Data | [interface{}](#interface)| `interface{}` |  | |  |  |
| Message | string| `string` |  | |  |  |



### <span id="device-state"></span> Get device info (*device-state*)

```
GET /api/v1/device/state
```

this will show device state

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| cache | `query` | boolean | `bool` |  |  |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#device-state-200) | OK |  |  | [schema](#device-state-200-schema) |

#### Responses


##### <span id="device-state-200"></span> 200
Status: OK

###### <span id="device-state-200-schema"></span> Schema
   
  

[DeviceStateOKBody](#device-state-o-k-body)

###### Inlined models

**<span id="device-state-o-k-body"></span> DeviceStateOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | |  |  |
| Message | string| `string` |  | |  |  |
| data | [DeviceInfo](#device-info)| `models.DeviceInfo` |  | |  |  |



### <span id="lan-hotspots"></span> Get devices(hotspots) address in LAN (*lan-hotspots*)

```
GET /api/v1/lan/hotspot
```

Device discovery by net interface `eth0`

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#lan-hotspots-200) | OK |  |  | [schema](#lan-hotspots-200-schema) |

#### Responses


##### <span id="lan-hotspots-200"></span> 200
Status: OK

###### <span id="lan-hotspots-200-schema"></span> Schema
   
  

[LanHotspotsOKBody](#lan-hotspots-o-k-body)

###### Inlined models

**<span id="lan-hotspots-o-k-body"></span> LanHotspotsOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | |  |  |
| Message | string| `string` |  | |  |  |
| data | [Dev](#dev)| `models.Dev` |  | |  |  |



### <span id="miner-restart"></span> Restart miner (*miner-restart*)

```
POST /api/v1/miner/restart
```

Restart miner container

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#miner-restart-200) | OK |  |  | [schema](#miner-restart-200-schema) |

#### Responses


##### <span id="miner-restart-200"></span> 200
Status: OK

###### <span id="miner-restart-200-schema"></span> Schema
   
  

[MinerRestartOKBody](#miner-restart-o-k-body)

###### Inlined models

**<span id="miner-restart-o-k-body"></span> MinerRestartOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | |  |  |
| Data | [interface{}](#interface)| `interface{}` |  | |  |  |
| Message | string| `string` |  | |  |  |



### <span id="miner-resync"></span> Resync miner (*miner-resync*)

```
POST /api/v1/miner/resync
```

Clean miner data and restart miner, miner will resync data

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#miner-resync-200) | OK |  |  | [schema](#miner-resync-200-schema) |

#### Responses


##### <span id="miner-resync-200"></span> 200
Status: OK

###### <span id="miner-resync-200-schema"></span> Schema
   
  

[MinerResyncOKBody](#miner-resync-o-k-body)

###### Inlined models

**<span id="miner-resync-o-k-body"></span> MinerResyncOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | |  |  |
| Data | [interface{}](#interface)| `interface{}` |  | |  |  |
| Message | string| `string` |  | |  |  |



### <span id="miner-state"></span> Get miner info (*miner-state*)

```
GET /api/v1/miner/state
```

Get miner info and state

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| cache | `query` | boolean | `bool` |  |  |  |  |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#miner-state-200) | OK |  |  | [schema](#miner-state-200-schema) |

#### Responses


##### <span id="miner-state-200"></span> 200
Status: OK

###### <span id="miner-state-200-schema"></span> Schema
   
  

[MinerStateOKBody](#miner-state-o-k-body)

###### Inlined models

**<span id="miner-state-o-k-body"></span> MinerStateOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | |  |  |
| Message | string| `string` |  | |  |  |
| data | [MinerInfo](#miner-info)| `models.MinerInfo` |  | |  |  |



### <span id="workspace-update"></span> Update worksapce (main git repo) (*workspace-update*)

```
POST /api/v1/workspace/update
```

Trigger workspace update and return immediately

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#workspace-update-200) | OK |  |  | [schema](#workspace-update-200-schema) |

#### Responses


##### <span id="workspace-update-200"></span> 200
Status: OK

###### <span id="workspace-update-200-schema"></span> Schema
   
  

[WorkspaceUpdateOKBody](#workspace-update-o-k-body)

###### Inlined models

**<span id="workspace-update-o-k-body"></span> WorkspaceUpdateOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | |  |  |
| Data | [interface{}](#interface)| `interface{}` |  | |  |  |
| Message | string| `string` |  | |  |  |



### <span id="workspace-update-get"></span> Check workspace update (*workspace-update-get*)

```
GET /api/v1/workspace/update
```

Whether worksapce (main git repo) is update available

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#workspace-update-get-200) | OK |  |  | [schema](#workspace-update-get-200-schema) |

#### Responses


##### <span id="workspace-update-get-200"></span> 200
Status: OK

###### <span id="workspace-update-get-200-schema"></span> Schema
   
  

[WorkspaceUpdateGetOKBody](#workspace-update-get-o-k-body)

###### Inlined models

**<span id="workspace-update-get-o-k-body"></span> WorkspaceUpdateGetOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | |  |  |
| Data | boolean| `bool` |  | |  |  |
| Message | string| `string` |  | |  |  |



## Models

### <span id="all-state-info"></span> AllStateInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| device | [DeviceInfo](#device-info)| `DeviceInfo` |  | |  |  |
| miner | [MinerInfo](#miner-info)| `MinerInfo` |  | |  |  |



### <span id="dev"></span> Dev


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Address | string| `string` |  | |  |  |
| Host | string| `string` |  | |  |  |
| Name | string| `string` |  | |  |  |
| Port | uint16 (formatted integer)| `uint16` |  | |  |  |



### <span id="device-info"></span> DeviceInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| CpuFreq | uint64 (formatted integer)| `uint64` |  | |  |  |
| CpuPercent | []double (formatted number)| `[]float64` |  | |  |  |
| CpuTemp | string| `string` |  | |  |  |
| NetInterface | [][NetInterfaceInfo](#net-interface-info)| `[]*NetInterfaceInfo` |  | |  |  |
| disk | [DiskInfo](#disk-info)| `DiskInfo` |  | |  |  |
| host | [HostInfo](#host-info)| `HostInfo` |  | |  |  |
| mem | [MemInfo](#mem-info)| `MemInfo` |  | |  |  |
| wifi | [WifiInfo](#wifi-info)| `WifiInfo` |  | |  |  |



### <span id="disk-info"></span> DiskInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Free | uint64 (formatted integer)| `uint64` |  | |  |  |
| Fstype | string| `string` |  | |  |  |
| Path | string| `string` |  | |  |  |
| Total | uint64 (formatted integer)| `uint64` |  | |  |  |
| Used | uint64 (formatted integer)| `uint64` |  | |  |  |
| UsedPercent | double (formatted number)| `float64` |  | |  |  |



### <span id="host-info"></span> HostInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| BootTime | uint64 (formatted integer)| `uint64` |  | |  |  |
| Hostname | string| `string` |  | |  |  |
| Uptime | uint64 (formatted integer)| `uint64` |  | |  |  |



### <span id="info-p2p-status"></span> InfoP2pStatus


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Connected | string| `string` |  | |  |  |
| Dialable | string| `string` |  | |  |  |
| Height | uint64 (formatted integer)| `uint64` |  | |  |  |
| NatType | string| `string` |  | |  |  |



### <span id="info-summary"></span> InfoSummary


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| FirmwareVersion | string| `string` |  | |  |  |
| Height | string| `string` |  | |  |  |
| Name | string| `string` |  | |  |  |
| Version | string| `string` |  | |  |  |



### <span id="mem-info"></span> MemInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Available | uint64 (formatted integer)| `uint64` |  | |  |  |
| Buffers | uint64 (formatted integer)| `uint64` |  | |  |  |
| Cached | uint64 (formatted integer)| `uint64` |  | |  |  |
| Free | uint64 (formatted integer)| `uint64` |  | |  |  |
| Shared | uint64 (formatted integer)| `uint64` |  | |  |  |
| Total | uint64 (formatted integer)| `uint64` |  | |  |  |
| Used | uint64 (formatted integer)| `uint64` |  | |  |  |



### <span id="miner-info"></span> MinerInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| InfoHeight | uint64 (formatted integer)| `uint64` |  | |  |  |
| InfoRegion | string| `string` |  | |  |  |
| PeerAddr | string| `string` |  | |  |  |
| infoP2pStatus | [InfoP2pStatus](#info-p2p-status)| `InfoP2pStatus` |  | |  |  |
| infoSummary | [InfoSummary](#info-summary)| `InfoSummary` |  | |  |  |
| peerBook | [PeerBook](#peer-book)| `PeerBook` |  | |  |  |



### <span id="net-interface-info"></span> NetInterfaceInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Addrs | [][NetInterfaceInfoAddrsItems0](#net-interface-info-addrs-items0)| `[]*NetInterfaceInfoAddrsItems0` |  | |  |  |
| HardwareAddr | string| `string` |  | |  |  |
| Index | uint64 (formatted integer)| `uint64` |  | |  |  |
| Mtu | uint64 (formatted integer)| `uint64` |  | |  |  |
| Name | string| `string` |  | |  |  |



#### Inlined models

**<span id="net-interface-info-addrs-items0"></span> NetInterfaceInfoAddrsItems0**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Addr | string| `string` |  | |  |  |



### <span id="peer-book"></span> PeerBook


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Address | string| `string` |  | |  |  |
| ConnectionCount | uint64 (formatted integer)| `uint64` |  | |  |  |
| LastUpdated | uint64 (formatted integer)| `uint64` |  | |  |  |
| ListenAddrCount | uint64 (formatted integer)| `uint64` |  | |  |  |
| ListenAddresses | []string| `[]string` |  | |  |  |
| Name | string| `string` |  | |  |  |
| Nat | string| `string` |  | |  |  |
| sessions | [Sessions](#sessions)| `Sessions` |  | |  |  |



#### Inlined models

**<span id="sessions"></span> Sessions**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Local | string| `string` |  | |  |  |
| Name | string| `string` |  | |  |  |
| P2p | string| `string` |  | |  |  |
| Remote | string| `string` |  | |  |  |



### <span id="proxy-item"></span> ProxyItem


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Value | string| `string` |  | |  |  |
| type | [ProxyType](#proxy-type)| `ProxyType` |  | |  |  |



### <span id="proxy-type"></span> ProxyType


  

| Name | Type | Go type | Default | Description | Example |
|------|------|---------| ------- |-------------|---------|
| ProxyType | string| string | |  |  |



### <span id="snapshot-state-res"></span> SnapshotStateRes


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| File | string| `string` |  | |  |  |
| State | string| `string` |  | |  |  |
| Time | date-time (formatted string)| `strfmt.DateTime` |  | |  |  |



### <span id="wifi-info"></span> WifiInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Connected | boolean| `bool` |  | |  |  |
| Name | string| `string` |  | |  |  |
| Powered | boolean| `bool` |  | |  |  |



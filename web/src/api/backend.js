import http from "./http";

export async function stateGet() {
  return await http.get('/inner/state')
}

// item: gitRepo or gitRelease
export async function proxyConfigGet(item) {
  return http.get(`/api/v1/config/proxy?item=${item}`)
}

// item: gitRepo or gitRelease
// body: { type: string, value: string }
export async function proxyConfigSet(item, body) {
  return http.post(`/api/v1/config/proxy?item=${item}`, body)
}

// device
export async function deviceReboot() {
  return await http.post("/api/v1/device/reboot");
}

export async function lanHotspots() {
  return await http.get("/api/v1/lan/hotspot");
}

export async function logQuery(params) {
  const { logType, filter, fromTime, toTime, limitLine } = params;
  let url = `/inner/api/v1/log?type=${logType}&filter=${filter}`;
  if (logType === "pktfwdLog") {
    url += `&since=${fromTime}&until=${toTime}`;
  } else {
    url += `&limit=${limitLine}`;
  }
  return await http.get(url);
}

// miner
export async function minerRestart() {
  return await http.post("/api/v1/miner/restart");
}
export async function minerResync() {
  return await http.post("/api/v1/miner/resync");
}

// workspace
export async function workspaceUpdateCheck() {
  const api = "/api/v1/workspace/update";
  return await http.get(api);
}

export async function workspaceUpdate() {
  const api = "/api/v1/workspace/update";
  return await http.post(api);
}

export async function workspaceReset() {
  const api = "/inner/api/v1/workspace/reset";
  await http.post(api);
}

// snapshot

export async function snap() {
  const api = "/inner/api/v1/miner/snapshot";
  return await http.post(api);
}

export async function snapState() {
  const api = "/inner/api/v1/miner/snapshot/state";
  return await http.get(api);
}

export function snapDownload(fileName) {
  open(`/inner/api/v1/miner/snapshot/file/${fileName}`, "_blank");
}

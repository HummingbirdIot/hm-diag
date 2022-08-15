import http from "./http";

export async function login(sec) {
  return await http.post("/api/v1/login", sec);
}

export async function stateGet() {
  return await http.get("/inner/state");
}

export async function version() {
  return await http.get("/inner/api/v1/version");
}

// item: gitRepo or gitRelease
export async function proxyConfigGet(item) {
  return http.get(`/api/v1/config/proxy?item=${item}`);
}

// item: gitRepo or gitRelease
// body: { type: string, value: string }
export async function proxyConfigSet(item, body) {
  return http.post(`/api/v1/config/proxy?item=${item}`, body);
}

export async function configGet() {
  return http.get("/inner/api/v1/config/safe")
}

export async function configSet(conf) {
  return http.post("/inner/api/v1/config/safe", conf)
}

// check if remote ip for server is private ip (eg: localhost, LAN IP)
export async function isViaPrivate() {
  return http.get("/inner/api/v1/safe/isViaPrivate")
}

// device
export async function deviceReboot() {
  return await http.post("/api/v1/device/reboot");
}

export async function blinkDeviceLight(durSec) {
  const d = durSec ? durSec : "";
  return await http.post(`/inner/api/v1/device/light/blink?durSec=${d}`);
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

export async function networkTest() {
  return await http.get("/inner/api/v1/network/ping");
}

// miner
export async function minerRestart() {
  return await http.post("/api/v1/miner/restart");
}
export async function minerResync() {
  return await http.post("/api/v1/miner/resync");
}
export async function onboarding(ownerAddr) {
  return await http.post("/inner/api/v1/miner/txn/onboarding?owner=" + ownerAddr);
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

// export async function snap() {
//   const api = "/inner/api/v1/miner/snapshot";
//   return await http.post(api);
// }

// export async function snapState() {
//   const api = "/inner/api/v1/miner/snapshot/state";
//   return await http.get(api);
// }

// export function snapDownload(fileName) {
//   open(`/inner/api/v1/miner/snapshot/file/${fileName}`, "_blank");
// }


//onboarding

export async function checkOnboarding() {
  const api = "/inner/api/v1/onboarding";
  return await http.get(api);
}

//log
export async function downloadLog() {
  return await http.get("/inner/api/v1/log/download",{
    responseType: 'blob'
  });
}
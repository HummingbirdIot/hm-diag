import http from "./http";
import axios from "axios";
const client = axios.create({ baseURL: "http://gh.xdt.com/hapi/" });
const heliumClient = axios.create({ baseURL: "https://api.helium.io/" });

export async function fetchHeliumHeight() {
  const api = `/inner/api/v1/proxy/heliumApi?path=/v1/blocks/height&t=${Date.now()}`;
  const r = await http.get(api)
  return JSON.parse(r).data.height
}
export async function blockHeight() {
  const r = await client.get('/v1/blocks/height');
  return r.data.data.height
}

export async function blockHeightFromHelium() {
  const r = await heliumClient.get('/v1/blocks/height');
  return r.data.data.height
}
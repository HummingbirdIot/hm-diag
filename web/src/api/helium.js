import http from "./http";
import axios from "axios";
const client = axios.create({ baseURL: "http://gh.xdt.com/hapi/" });

export async function fetchHeliumHeight() {
  const api = `/inner/api/v1/proxy/heliumApi?path=/v1/blocks/height&t=${Date.now()}`;
  console.log(api);
  const r = await http.get(api)
  return JSON.parse(r).data.height
}
export async function blockHeight() {
  const r = await client.get('/v1/blocks/height');
  return r.data.data.height
}
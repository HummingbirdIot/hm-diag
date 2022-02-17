import http from "./http";

export async function fetchHeliumHeight() {
  const api = `/inner/api/v1/proxy/heliumApi?path=/v1/blocks/height&t=${Date.now()}`;
  console.log(api);
  const r = await http.get(api)
  return JSON.parse(r).data.height
}

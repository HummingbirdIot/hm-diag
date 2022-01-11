import http from "./http";

export async function checkWorkspaceUpdate() {
  const api = "/api/v1/workspace/update";
  const r = await http.get(api);
  return r.data
}

export async function workspaceUpdate() {
  const api = "/api/v1/workspace/update";
  const r = await http.post(api);
  return r.data
}


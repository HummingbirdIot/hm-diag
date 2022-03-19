import axios from "axios";
import { Notify } from "vant";

export function getBase() {
  const s = window.location.search;
  const p = window.location.pathname;
  let ip = "";
  if (s.length > 1 && s.indexOf("ip") > 0) {
    const f = s
      .substring(1)
      .split("&")
      .find((i) => i.indexOf("ip=") === 0);
    ip = f ? f.substring(3) : "";
    return "/proxy/" + ip;
  } else if (p.indexOf("/proxy/") === 0) {
    const arr = p.split("/");
    return arr.slice(0, 3).join("/");
  } else {
    return "/";
  }
}

function getToken() {
  const s = window.location.search;
  let tk = "";
  if (s.length > 1 && s.indexOf("tk") > 0) {
    const f = s
      .substring(1)
      .split("&")
      .find((i) => i.indexOf("tk=") === 0);
    tk = f ? f.substring(3) : "";
  }
  return tk;
}

const client = axios.create({ baseURL: getBase() });

client.interceptors.request.use(
  (req) => {
    Object.assign(req.headers, {
      Authorization: getToken(),
    });
    return req;
  },
  (err) => {
    throw err;
  }
);

client.interceptors.response.use(
  (res) => {
    if (res.data.code < 400) {
      return res.data.data;
    } else {
      console.log("api error: ", res);
      Notify(res.data.message);
      throw res.data.message;
    }
  },
  (err) => {
    console.log("api error: ", err);
    const msg = err.response?.data?.message ? err.response.data.message : err.message ? err.message : err;
    Notify(msg);
    throw msg;
  }
);

export default client;

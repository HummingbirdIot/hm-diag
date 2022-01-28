import axios from "axios";
import { Notify } from "vant";

const client = axios.create({ baseURL: "/" });

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

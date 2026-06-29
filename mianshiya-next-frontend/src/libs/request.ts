import axios from "axios";

// 创建 Axios 实例
// Server 端(SSR) 用 NEXT_PUBLIC_API_URL，Docker 中指向 backend 容器
// Client 端(浏览器) 用空 baseURL，通过 Nginx 代理到同源
const BASE_URL = typeof window === "undefined"
  ? (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8101")
  : "";
const myAxios = axios.create({
  baseURL: BASE_URL,
  timeout: 60000,
  withCredentials: true,
});

// 创建请求拦截器
myAxios.interceptors.request.use(
  function (config) {
    // 请求执行前执行
    if (typeof window !== "undefined") {
      const token = window.localStorage.getItem("token");
      if (token) {
        config.headers = config.headers ?? {};
        (config.headers as any).Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  function (error) {
    // 处理请求错误
    return Promise.reject(error);
  },
);

// 创建响应拦截器
myAxios.interceptors.response.use(
  // 2xx 响应触发
  function (response) {
    // 处理响应数据
    const { data } = response;
    // 未登录
    if (data.code === 40100) {
      if (typeof window !== "undefined") {
        window.localStorage.removeItem("token");
      }
      // 不是获取用户信息接口，或者不是登录页面，则跳转到登录页面
      const responseUrl = response.request?.responseURL ?? "";
      if (
        typeof window !== "undefined" &&
        !responseUrl.includes("user/get/login") &&
        !window.location.pathname.includes("/user/login")
      ) {
        window.location.href = `/user/login?redirect=${window.location.href}`;
      }
    } else if (data.code !== 0) {
      // 其他错误
      throw new Error(data.message ?? "服务器错误");
    }
    return data;
  },
  // 非 2xx 响应触发
  function (error) {
    // 处理响应错误
    return Promise.reject(error);
  },
);

export default myAxios;

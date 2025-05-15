import { fetchApi } from "../api";

const API_ENDPOINTS = {
  login: "/api/v1alpha1/login",
  loginCallback: "/api/v1alpha1/login/callback",
};

export const loginApi = {
  getLogin: () => fetchApi(API_ENDPOINTS.login, "GET", null),
};

export const LoginCallbacksApi = {
  getLoginCallback: (payload) =>
    fetchApi(API_ENDPOINTS.loginCallback, "POST", payload),
};

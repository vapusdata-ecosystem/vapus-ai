import { fetchApi } from "../api";

const API_ENDPOINTS = {
  ORGACCESSTOKEN: "/api/v1alpha1/auth/token",
};

export const OrgAccessTokenApi = {
  getOrgAccessToken: (payload) =>
    fetchApi(API_ENDPOINTS.ORGACCESSTOKEN, "POST", payload),
};

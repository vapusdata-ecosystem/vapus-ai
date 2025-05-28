import { fetchApi } from "../api";

const API_ENDPOINTS = {
  PLATFORM: "/api/v1alpha1/platform",
};

export const platformApi = {
  getPlatform: () => fetchApi(API_ENDPOINTS.PLATFORM, "GET", null),
};

import { fetchApi } from "../api";

const API_ENDPOINTS = {
  // AI Studio
  PLATFORM_DOMAIN: "/api/v1alpha1/domains",
  PLATFORM_CREATE: "",
};

export const platformDomainApi = {
  getplatformdomain: () => fetchApi(API_ENDPOINTS.PLATFORM_DOMAIN, "GET", null),
};

export const platformCreateApi = {
  getplatformCreate: (payload) =>
    fetchApi(API_ENDPOINTS.PLATFORM_CREATE, "POST", payload),
};

import { fetchApi } from "../api";

const API_ENDPOINTS = {
  // AI Studio
  PLATFORM_DOMAIN: "/api/v1alpha1/domains",
};

export const platformDomainApi = {
  getplatformdomain: () => fetchApi(API_ENDPOINTS.PLATFORM_DOMAIN, "GET", null),
};

export const platformDomainCreateApi = {
  getplatformDomainCreate: (payload) =>
    fetchApi(API_ENDPOINTS.PLATFORM_DOMAIN, "POST", payload),
};

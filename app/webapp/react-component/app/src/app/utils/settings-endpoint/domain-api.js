import { fetchApi } from "../api";

const API_ENDPOINTS = {
  DOMAIN: "/api/v1alpha1/domains/{domain_id}",
  ADD_USERS: "/api/v1alpha1/domains/{domain_id}/users",
};

export const domainApi = {
  getDomains: () => fetchApi(API_ENDPOINTS.DOMAIN, "GET", null),
};

export const addUsersApi = {
  getAddUsers: (payload) => fetchApi(API_ENDPOINTS.ADD_USERS, "PUT", payload),
};

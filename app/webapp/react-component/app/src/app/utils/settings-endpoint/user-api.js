import { fetchApi } from "../api";

const API_ENDPOINTS = {
  USER: "/api/v1alpha1/users/anand@vapusdata.com",
};

export const userApi = {
  getuser: () => fetchApi(API_ENDPOINTS.USER, "GET", null),
  getuserId: () => fetchApi(API_ENDPOINTS.USER, "GET", null),
};

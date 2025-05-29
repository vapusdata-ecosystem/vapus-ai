import { fetchApi } from "../api";

const API_ENDPOINTS = {
  USER: "/api/v1alpha1/users",
};

export const userApi = {
  getuser: (action) => {
    const url = action
      ? `${API_ENDPOINTS.USER}/?action=${action}`
      : API_ENDPOINTS.USER;
    return fetchApi(url, "GET", null);
  },
  getuserId: (user_id) =>
    fetchApi(`${API_ENDPOINTS.USER}/${user_id}`, "GET", null),
};

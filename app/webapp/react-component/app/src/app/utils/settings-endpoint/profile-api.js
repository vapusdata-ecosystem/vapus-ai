import { fetchApi } from "../api";

const API_ENDPOINTS = {
  USER: "/api/v1alpha1/users",
};

export const userProfileApi = {
  getuserProfile: (userId) =>
    fetchApi(`${API_ENDPOINTS.USER}/${userId}`, "GET", null),
};
export const updateProfile = {
  getUpdateProfile: (payload) => fetchApi(API_ENDPOINTS.USER, "POST", payload),
};

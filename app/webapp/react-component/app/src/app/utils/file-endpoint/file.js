import { fetchApi } from "../api";

const API_ENDPOINTS = {
  USER: "/api/v1alpha1/utility/upload",
};

export const UploadFileApi = {
  getUploadFile: (uploadPayload) =>
    fetchApi(API_ENDPOINTS.USER, "POST", uploadPayload),
};

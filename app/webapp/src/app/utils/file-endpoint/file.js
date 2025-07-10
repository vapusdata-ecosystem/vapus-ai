import { fetchApi } from "../api";

const API_ENDPOINTS = {
  UPLOAD: "/api/v1alpha1/utility/upload",
  DOWNLOAD: "/api/v1alpha1/utility/download",
};

export const UploadFileApi = {
  getUploadFile: (uploadPayload) =>
    fetchApi(API_ENDPOINTS.UPLOAD, "POST", uploadPayload),
};

export const DownloadFileApi = {
  getDownloadFile: (params) => {
    const queryParams = new URLSearchParams();
    if (params && params.path) {
      queryParams.append("path", params.path);
    }
    const urlWithParams = `${API_ENDPOINTS.DOWNLOAD}?${queryParams.toString()}`;
    return fetchApi(urlWithParams, "GET", null);
  },
};

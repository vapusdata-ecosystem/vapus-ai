import { fetchApi } from "../api";

const API_ENDPOINTS = {
  SECRET_STORE: "/api/v1alpha1/secrets",
};

export const secretStoreApi = {
  getSecretStore: () => fetchApi(API_ENDPOINTS.SECRET_STORE, "GET", null),
  getSecretStoreName: (name) =>
    fetchApi(`${API_ENDPOINTS.SECRET_STORE}/${name}`, "GET", null),
};

export const secretStoreFormApi = {
  getSecretStoreForm: (payload) =>
    fetchApi(API_ENDPOINTS.SECRET_STORE, "POST", payload),
};

export const secretStoreArchiveApi = {
  getSecretStoreArchive: (name) =>
    fetchApi(`${API_ENDPOINTS.SECRET_STORE}/${name}`, "DELETE", null),
};

import { fetchApi } from "../api";

const API_ENDPOINTS = {
  // AI Studio
  MODELS_REGISTRY: "/api/v1alpha1/aistudio/models",
};

export const modelsRegistryApi = {
  getModelsRegistry: () => fetchApi(API_ENDPOINTS.MODELS_REGISTRY, "GET", null),
  getModelsRegistryID: (ai_model_node_id) =>
    fetchApi(
      `${API_ENDPOINTS.MODELS_REGISTRY}/${ai_model_node_id}`,
      "GET",
      null
    ),
};

export const modelsRegistryFormApi = {
  getmodelsRegistryForm: (payload) =>
    fetchApi(API_ENDPOINTS.MODELS_REGISTRY, "POST", payload),
};

export const modelsRegistryUpdateFormApi = {
  getmodelsRegistryUpdateForm: (payload) =>
    fetchApi(API_ENDPOINTS.MODELS_REGISTRY, "PUT", payload),
};

export const modelsRegistryArchiveApi = {
  getModelsRegistryArchive: (ai_model_node_id) =>
    fetchApi(
      `${API_ENDPOINTS.MODELS_REGISTRY}/${ai_model_node_id}`,
      "DELETE",
      null
    ),
};

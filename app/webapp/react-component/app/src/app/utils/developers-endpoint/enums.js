import { fetchApi } from "../api";

const API_ENDPOINTS = {
  ENUMS: "/api/v1alpha1/resources-spec",
  BEDROCKGUARDRAILS: "/api/v1alpha1/aistudio/guardrails-plugins",
};

export const enumsApi = {
  getEnums: () => fetchApi(API_ENDPOINTS.ENUMS, "GET", null),
};

export const bedrockGuardrailsApi = {
  getBedrockGuardrailsApi: () =>
    fetchApi(API_ENDPOINTS.BEDROCKGUARDRAILS, "GET", null),
};

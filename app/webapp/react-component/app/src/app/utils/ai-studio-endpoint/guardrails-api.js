import { fetchApi } from "../api";

const API_ENDPOINTS = {
  GUARDRAILS: "/api/v1alpha1/aistudio/guardrails",
};

export const GuardrailApi = {
  getGuardrail: () => fetchApi(API_ENDPOINTS.GUARDRAILS, "GET", null),
  getGuardrailId: (guardrail_id) =>
    fetchApi(`${API_ENDPOINTS.GUARDRAILS}/${guardrail_id}`, "GET", null),
};

export const GuardrailFormApi = {
  getGuardrailForm: (payload) =>
    fetchApi(API_ENDPOINTS.GUARDRAILS, "POST", payload),
};
export const GuardrailUpdateFormApi = {
  getGuardrailUpdateForm: (payload) =>
    fetchApi(API_ENDPOINTS.GUARDRAILS, "PUT", payload),
};

export const GuardrailArchiveApi = {
  getGuardrailArchive: (guardrail_id) =>
    fetchApi(`${API_ENDPOINTS.GUARDRAILS}/${guardrail_id}`, "DELETE", null),
};

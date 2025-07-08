import { fetchApi } from "../api";

const API_ENDPOINTS = {
  PROMPTS: "/api/v1alpha1/aistudio/prompts",
};

export const PromptsApi = {
  getPrompts: () => fetchApi(API_ENDPOINTS.PROMPTS, "GET", null),
  getPromptsId: (prompt_id) =>
    fetchApi(`${API_ENDPOINTS.PROMPTS}/${prompt_id}`, "GET", null),
};

export const promptsFormApi = {
  getPromptsForm: (payload) => fetchApi(API_ENDPOINTS.PROMPTS, "POST", payload),
};

export const promptsUpdateFormApi = {
  getPromptsUpdteForm: (payload) =>
    fetchApi(API_ENDPOINTS.PROMPTS, "PUT", payload),
};

export const promptArchiveApi = {
  getPromptArchive: (prompt_id) =>
    fetchApi(`${API_ENDPOINTS.PROMPTS}/${prompt_id}`, "DELETE", null),
};

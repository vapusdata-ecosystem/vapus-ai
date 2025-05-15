import { fetchApi } from "../api";

const API_ENDPOINTS = {
  ENUMS: "/api/v1alpha1/resources-spec",
};

export const enumsApi = {
  getEnums: () => fetchApi(API_ENDPOINTS.ENUMS, "GET", null),
};

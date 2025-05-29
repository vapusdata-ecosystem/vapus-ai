import { fetchApi } from "../api";

const API_ENDPOINTS = {
  PLUGINS: "/api/v1alpha1/plugins",
};

export const pluginsApi = {
  getPlugins: () => fetchApi(API_ENDPOINTS.PLUGINS, "GET", null),
  getPluginsId: (plugin_id) =>
    fetchApi(`${API_ENDPOINTS.PLUGINS}/${plugin_id}`, "GET", null),
};

export const pluginsFormApi = {
  getpluginsForm: (payload) => fetchApi(API_ENDPOINTS.PLUGINS, "POST", payload),
};

export const pluginsUpdateFormApi = {
  getPluginsUpdateForm: (payload) =>
    fetchApi(API_ENDPOINTS.PLUGINS, "PUT", payload),
};

export const pluginsArchiveApi = {
  getPluginsArchive: (plugin_id) =>
    fetchApi(`${API_ENDPOINTS.PLUGINS}/${plugin_id}`, "DELETE", null),
};

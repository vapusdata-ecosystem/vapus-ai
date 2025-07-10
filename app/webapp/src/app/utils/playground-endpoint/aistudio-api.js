import { fetchApi } from "../api";


const API_ENDPOINTS = {
  CREATE_CHAT: "/api/v1alpha1/aistudio/chats",
  AI_CHAT:"/api/v1alpha1/aistudio/chat",
  AI_GATEWAY:"/gateway/v1/chat/completions"
};

export const AiChatApi = {
  getAiChat: (payload, options ) => fetchApi(API_ENDPOINTS.AI_CHAT, "POST", payload, options),
};

export const AiGatewayChatApi = {
  getAiGatewayChat: (payload, options) => fetchApi(API_ENDPOINTS.AI_GATEWAY, "POST", payload, options),
};

export const createNewChatApi = {
  getCreateChat: (payload) => fetchApi(API_ENDPOINTS.CREATE_CHAT, "POST", payload),
};

export const chatHistoryApi = {
  getChatHistory: () => fetchApi(API_ENDPOINTS.CREATE_CHAT, "GET", null ),
   getChatData: (chat_id) =>
      fetchApi(`${API_ENDPOINTS.CREATE_CHAT}?chat_id=${chat_id}`, "GET", null),
};


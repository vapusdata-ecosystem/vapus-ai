export let accessToken = "";
const BASE_URL = "http://127.0.0.1:9017";

// Function to get cookie value by name
const getCookie = (name) => {
  if (typeof document === "undefined") {
    return null;
  }

  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) {
    return parts.pop().split(";").shift();
  }
  return null;
};

// Function to initialize access token from cookies
const initializeAccessToken = () => {
  const tokenFromCookie = getCookie("access_token");
  if (tokenFromCookie && !accessToken) {
    accessToken = tokenFromCookie;
  }
};

export const setAccessToken = (token) => {
  accessToken = token;
};

// Function to get the current access token
export const getAccessToken = () => {
  if (!accessToken) {
    initializeAccessToken();
  }
  return accessToken;
};

if (typeof window !== "undefined") {
  initializeAccessToken();
}

export const fetchApi = async (endpoint, method, payload, options = {}) => {
  console.log("submit the data", payload);

  // Ensure we have the latest token from cookies
  const currentToken = getAccessToken();

  const fullUrl = `${BASE_URL}${endpoint}`;
  const defaultOptions = {
    method: method,
    headers: {
      authorization: `Bearer ${currentToken}`,
      "Content-Type": "application/json",
    },
    redirect: "follow",
    ...options,
  };

  // If options already has headers, merge them with our default headers
  if (options.headers) {
    defaultOptions.headers = {
      ...defaultOptions.headers,
      ...options.headers,
    };
  }

  // Add the payload
  if (payload) {
    defaultOptions.body = JSON.stringify(payload);
  }

  const response = await fetch(fullUrl, defaultOptions);
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }
  return response.json();
};

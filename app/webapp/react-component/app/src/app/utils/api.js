export let accessToken = "";
const BASE_URL = "http://127.0.0.1:9017";

// Function to update the access token from other components
export const setAccessToken = (token) => {
  accessToken = token;
};

// Function to get the current access token
export const getAccessToken = () => {
  return accessToken;
};

export const fetchApi = async (endpoint, method, payload, options = {}) => {
  console.log("submit the data", payload);

  const fullUrl = `${BASE_URL}${endpoint}`;

  const defaultOptions = {
    method: method,
    headers: {
      authorization: `Bearer ${accessToken}`,
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

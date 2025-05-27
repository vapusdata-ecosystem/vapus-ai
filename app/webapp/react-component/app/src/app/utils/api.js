export let accessToken =
  "eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjkwMTQiLCJzdWIiOiJWYXB1c0RhdGEgYWNjZXNzIHRva2VuIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjE6OTAxNCJdLCJleHAiOjE3NDg0Mjg5NDAsIm5iZiI6MTc0ODM0MjU0MSwiaWF0IjoxNzQ4MzQyNTQwLCJqdGkiOiJhMDA5YjRkZi00YjViLTQ2OTctODVhMi1lNjY0ODk0NzE5MjUiLCJzY29wZSI6eyJ1c2VySWQiOiJhbmFuZEB2YXB1c2RhdGEuY29tIiwiYWNjb3VudElkIjoiYWNjLTQzZGMwMGIzLTA0NzYtNDQ1Yy04MDU0LTE0Mzc2MDY4MTFjYyIsIm9yZ2FuaXphdGlvbklkIjoiMmMwNWJjOGYtMzNiMy00MzZhLWE5MDAtZGVjOTc0ZjdhZWEyIiwib3JnYW5pemF0aW9uUm9sZSI6IlNFUlZJQ0VfT1dORVIiLCJyb2xlU2NvcGUiOiJPcmdhbml6YXRpb25TY29wZSJ9fQ.AXLHcEZBN4fMPb86uBC1gFJ1f-UgyfI09AKrvoK1uVxDo9k_BVGedb2E-RP0rega_1O-iPFshm8Im4EKrzK4uYHDAGzUWPXTfZ_WygP_6c8VmpnfdrZjn1o1BT8naNxEWMJkn2pCphON86s-wcSpTWSpC-kzENBkd400ez638Mwul1zw";
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

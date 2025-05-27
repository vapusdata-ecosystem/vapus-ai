export let accessToken =
  "eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjkwMTQiLCJzdWIiOiJWYXB1c0RhdGEgYWNjZXNzIHRva2VuIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjE6OTAxNCJdLCJleHAiOjE3NDg0MDkzOTksIm5iZiI6MTc0ODMyMzAwMCwiaWF0IjoxNzQ4MzIyOTk5LCJqdGkiOiI0NTllMjFlOC0xNGVjLTQyMDMtYjY4My05YTYzNzQyMWU0NWQiLCJzY29wZSI6eyJ1c2VySWQiOiJhbmFuZEB2YXB1c2RhdGEuY29tIiwiYWNjb3VudElkIjoiYWNjLWIzZmYyOGMzLWY3YzYtNDdlMi1hYzc5LTFmM2IzZWViOGVmNCIsImRvbWFpbklkIjoiZG1uLWI3NzM0OTFkLWU5M2ItNGU4NC05NmZiLWRkYjgzN2IyM2EzMyIsImRvbWFpblJvbGUiOiJET01BSU5fT1dORVJTIiwicm9sZVNjb3BlIjoiRG9tYWluU2NvcGUiLCJwbGF0Zm9ybVJvbGUiOiJQTEFURk9STV9VU0VSUyJ9fQ.ANJ4NFS6w9RQgsPV3neb9JhM5ZqZdbFFgSqijtR3mD_v54J3FezX0txxGj8iDmZS6uJSy2RAslOplI91BE0yhpEDAUe53NzDQZE2l_WrI9kcaes1BGevrbsd9ZhabVCGICS3Etkxrms4v7fyTecwK0yhWQWRIjV5MCxJuWuNHqFe1xOc";
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

const accessToken =
  "eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjkwMTQiLCJzdWIiOiJWYXB1c0RhdGEgYWNjZXNzIHRva2VuIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjE6OTAxNCJdLCJleHAiOjE3NDc5Nzc3NjQsIm5iZiI6MTc0Nzg5MTM2NSwiaWF0IjoxNzQ3ODkxMzY0LCJqdGkiOiJiMTg4NmQxMC0yMDg3LTQ5ZjAtODUzMC03M2Y3Zjg2ODQxYTUiLCJzY29wZSI6eyJ1c2VySWQiOiJhbmFuZEB2YXB1c2RhdGEuY29tIiwiYWNjb3VudElkIjoiYWNjLWIzZmYyOGMzLWY3YzYtNDdlMi1hYzc5LTFmM2IzZWViOGVmNCIsImRvbWFpbklkIjoiZG1uLWI3NzM0OTFkLWU5M2ItNGU4NC05NmZiLWRkYjgzN2IyM2EzMyIsImRvbWFpblJvbGUiOiJET01BSU5fT1dORVJTIiwicm9sZVNjb3BlIjoiRG9tYWluU2NvcGUiLCJwbGF0Zm9ybVJvbGUiOiJQTEFURk9STV9VU0VSUyJ9fQ.AYrZ1pxUC2HvX0obZautNCPSdEgoDa3vEilapsDzY0qlyqBKlEcXJ2tWOVkAJprTjM1M91MUQgEVZjCcVZ--gFscAV2j3CBU1Iqed1obIWBSS8lHwLr0_Mc2Wu8WOP4GppF2IAxYzVOG9TreZn8shWzT57MY5a9WMbJZvY-0C8dolfou";
const BASE_URL = "http://127.0.0.1:9017";

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

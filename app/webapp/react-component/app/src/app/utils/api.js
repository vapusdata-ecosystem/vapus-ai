const accessToken =
  "eyJhbGciOiJFUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjkwMTQiLCJzdWIiOiJWYXB1c0RhdGEgYWNjZXNzIHRva2VuIiwiYXVkIjpbImh0dHA6Ly8xMjcuMC4wLjE6OTAxNCJdLCJleHAiOjE3NDc4OTE5MTQsIm5iZiI6MTc0NzgwNTUxNSwiaWF0IjoxNzQ3ODA1NTE0LCJqdGkiOiJiNzMwY2JlZi0yYjliLTQwZTktYTM5Yi00MTk5MDRhNWVkNWEiLCJzY29wZSI6eyJ1c2VySWQiOiJhbmFuZEB2YXB1c2RhdGEuY29tIiwiYWNjb3VudElkIjoiYWNjLWIzZmYyOGMzLWY3YzYtNDdlMi1hYzc5LTFmM2IzZWViOGVmNCIsImRvbWFpbklkIjoiZG1uLWI3NzM0OTFkLWU5M2ItNGU4NC05NmZiLWRkYjgzN2IyM2EzMyIsImRvbWFpblJvbGUiOiJET01BSU5fT1dORVJTIiwicm9sZVNjb3BlIjoiRG9tYWluU2NvcGUiLCJwbGF0Zm9ybVJvbGUiOiJQTEFURk9STV9VU0VSUyJ9fQ.AGEO1J3qF5fcgXBU7YzR7V7qPaRBmfnH59_hbS2zdy5khVBS4T_Pql6PnEvYm2IrbTMM90qf05Xv9ESkY4lhIlcdAUaBMNpfhvfEFkA6H0mk65ABbkuJHC3pjDxYE-fBarqVhb4vuyqBPdhkYKbggqfncc--S54IGm8QKOJvBzAyDyhH";
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

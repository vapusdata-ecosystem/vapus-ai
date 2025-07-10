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

// Helper function to parse multiple JSON objects from streaming response
const parseStreamingResponse = (responseText) => {
  const chunks = [];
  
  // Handle case where JSON objects are concatenated without proper separation
  try {
    const singleJson = JSON.parse(responseText);
    return [singleJson];
  } catch (error) {
    console.log("Parsing as streaming response...");
  }
  
  // Check if this is Server-Sent Events (SSE) format with "data:" prefix
  if (responseText.includes('data:')) {
    return parseSSEResponse(responseText);
  }
  
  // Handle regular concatenated JSON objects
  const parts = responseText.split(/(?<=})\s*(?={)/);
  
  for (const part of parts) {
    const trimmedPart = part.trim();
    if (trimmedPart) {
      try {
        const parsed = JSON.parse(trimmedPart);
        chunks.push(parsed);
      } catch (error) {
        console.error('Failed to parse chunk:', trimmedPart, error);
      }
    }
  }
  
  return chunks;
};

// Helper function to parse Server-Sent Events (SSE) format
const parseSSEResponse = (responseText) => {
  const chunks = [];
  const lines = responseText.split('\n');
  
  for (const line of lines) {
    const trimmedLine = line.trim();
    
    // Skip empty lines and [DONE] markers
    if (!trimmedLine || trimmedLine === 'data: [DONE]' || trimmedLine === '[DONE]') {
      continue;
    }
    
    // Handle SSE data format
    if (trimmedLine.startsWith('data: ')) {
      const jsonStr = trimmedLine.substring(6);
      
      try {
        const parsed = JSON.parse(jsonStr);
        chunks.push(parsed);
      } catch (error) {
        console.error('Failed to parse SSE chunk:', jsonStr, error);
      }
    }
  }
  
  return chunks;
};

// Helper function to combine streaming chunks into final response
const combineStreamingChunks = (chunks) => {
  if (!chunks || chunks.length === 0) {
    return null;
  }
  
  // If only one chunk, return it as is
  if (chunks.length === 1) {
    return chunks[0];
  }
  
  // Combine multiple chunks 
  let combinedContent = '';
  let finalChunk = null;
  
  for (const chunk of chunks) {
    // Handle original format
    if (chunk.result?.choices?.[0]?.messages?.content) {
      combinedContent += chunk.result.choices[0].messages.content;
      finalChunk = chunk;
    }
    // Handle  SSE format (delta.content)
    else if (chunk.choices?.[0]?.delta?.content) {
      combinedContent += chunk.choices[0].delta.content;
      finalChunk = chunk;
    }
  }
  
  // Create final response with combined content
  if (finalChunk && combinedContent) {
    // Handle original format
    if (finalChunk.result?.choices?.[0]?.messages) {
      return {
        ...finalChunk,
        result: {
          ...finalChunk.result,
          choices: [{
            ...finalChunk.result.choices[0],
            messages: {
              ...finalChunk.result.choices[0].messages,
              content: combinedContent
            }
          }]
        }
      };
    }
    // Handle  SSE format 
    else if (finalChunk.choices?.[0]?.delta) {
      return {
        result: {
          created: finalChunk.created,
          model: finalChunk.model,
          object: "chat.completion",
          choices: [{
            messages: {
              role: "assistant",
              content: combinedContent,
              structuredContent: [],
              toolCalls: []
            },
            finishReason: finalChunk.choices[0].finish_reason || "",
            index: finalChunk.choices[0].index || "0",
            delta: null
          }],
          usage: [],
          event: "data",
          id: finalChunk.id,
          serviceTier: "",
          systemFingerprint: finalChunk.system_fingerprint || ""
        }
      };
    }
  }
  
  return chunks[chunks.length - 1]; 
};

// Updated fetchApi function with streaming support
export const fetchApi = async (endpoint, method, payload, options = {}) => {
  console.log("submit the data", payload);

  // Ensure we have the latest token from cookies
  const currentToken = getAccessToken();

  const fullUrl = `${BASE_URL}${endpoint}`;
  
  // Create default headers
  const defaultHeaders = {
    authorization: `Bearer ${currentToken}`,
    "Content-Type": "application/json",
  };

  // Merge headers properly - options.headers should override defaults where there are conflicts
  const mergedHeaders = {
    ...defaultHeaders,
    ...(options.headers || {}),
  };

  // Create the final options object
  const defaultOptions = {
    method: method,
    headers: mergedHeaders,
    redirect: "follow",
    // Spread other options except headers (since we handled headers separately)
    ...Object.fromEntries(
      Object.entries(options).filter(([key]) => key !== 'headers')
    ),
  };

  // Add the payload
  if (payload) {
    defaultOptions.body = JSON.stringify(payload);
  }

  const response = await fetch(fullUrl, defaultOptions);
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }

  // Get response as text first
  const responseText = await response.text();
  
  // Check if response is streaming (multiple JSON objects)
  const chunks = parseStreamingResponse(responseText);
  
  // Return combined result
  return combineStreamingChunks(chunks);
};

"use client";
import Head from "next/head";
import { useEffect, useState, useCallback } from "react";
import Script from "next/script";
import { useRouter } from "next/navigation";
import Header from "@/app/components/platform/header";
import AIToolCallPopup from "../aitoolcallpoppup";
import NestedDropdown from "@/app/components/nestedDropdown";
import PromptDropdown from "@/app/components/promptDropdown";
import { toast } from "react-toastify";
import ContextModal from "../ai-context-popup";
import { AiChatApi, AiGatewayChatApi, chatHistoryApi, createNewChatApi } from "@/app/utils/playground-endpoint/aistudio-api";

export default function AIStudio({ response, globalContext, aiStudioChat }) {
  const router = useRouter();
  
  // State management
  const [isToolCallModalOpen, setIsToolCallModalOpen] = useState(false);
  const [isContextModalOpen, setIsContextModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isStreaming, setIsStreaming] = useState(false);
  const [currentChatId, setCurrentChatId] = useState("");
  const [apiHeaders, setApiHeaders] = useState({}); 
  const [chatMessages, setChatMessages] = useState([]); 
  const [chatHistory, setChatHistory] = useState([]);
  const [error, setError] = useState(null);
  const [formData, setFormData] = useState({
    spec: {
      Tools: [],
    },
  });

  // Enhanced function to format message content with colorful code blocks
  const formatMessageContent = (content) => {
    // Split content by code blocks
    const parts = content.split(/(```[\s\S]*?```)/g);
    
    return parts.map((part, index) => {
      if (part.startsWith('```') && part.endsWith('```')) {
        // Extract language and code
        const lines = part.split('\n');
        const firstLine = lines[0].replace('```', '');
        const language = firstLine.trim() || 'text';
        const code = lines.slice(1, -1).join('\n');
        
        return (
          <div key={index} className="my-4">
            <div className="bg-[#0d1117] rounded-lg overflow-hidden border border-gray-700">
              <div className="flex items-center justify-between bg-[#161b22] px-4 py-2 text-sm border-b border-gray-700">
                <span className="text-gray-300 font-medium">{language}</span>
                <button
                  onClick={() => copyToClipboard(code)}
                  className="flex items-center space-x-1 bg-pink-900 hover:bg-pink-700 cursor-pointer px-3 py-1 rounded text-white transition-all duration-200 text-xs font-medium"
                >
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 2z" />
                  </svg>
                  <span>Copy</span>
                </button>
              </div>
              <div className="syntax-highlight-container">
                <pre className="p-4 overflow-x-auto text-sm leading-relaxed">
                  <code 
                    className={`language-${language} hljs`} 
                    dangerouslySetInnerHTML={{
                      __html: highlightCodeWithColors(code, language)
                    }} 
                  />
                </pre>
              </div>
            </div>
          </div>
        );
      } else {
        // Regular text content
        return (
          <span key={index} className="whitespace-pre-wrap">
            {part}
          </span>
        );
      }
    });
  };

  // Enhanced function to highlight code with colors
  const highlightCodeWithColors = (code, language) => {
    if (typeof window !== 'undefined' && window.hljs) {
      try {
        let highlightedCode;
        if (language && window.hljs.getLanguage(language)) {
          highlightedCode = window.hljs.highlight(code, { language }).value;
        } else {
          highlightedCode = window.hljs.highlightAuto(code).value;
        }
        
        // Apply additional color styling if hljs classes aren't working
        return highlightedCode
          .replace(/<span class="hljs-keyword">/g, '<span class="hljs-keyword" style="color: #ff7b72; font-weight: 600;">')
          .replace(/<span class="hljs-string">/g, '<span class="hljs-string" style="color: #a5d6ff;">')
          .replace(/<span class="hljs-number">/g, '<span class="hljs-number" style="color: #79c0ff;">')
          .replace(/<span class="hljs-comment">/g, '<span class="hljs-comment" style="color: #8b949e; font-style: italic;">')
          .replace(/<span class="hljs-function">/g, '<span class="hljs-function" style="color: #d2a8ff;">')
          .replace(/<span class="hljs-variable">/g, '<span class="hljs-variable" style="color: #ffa657;">')
          .replace(/<span class="hljs-title">/g, '<span class="hljs-title" style="color: #7ee787;">')
          .replace(/<span class="hljs-attr">/g, '<span class="hljs-attr" style="color: #79c0ff;">')
          .replace(/<span class="hljs-tag">/g, '<span class="hljs-tag" style="color: #7ee787;">')
          .replace(/<span class="hljs-name">/g, '<span class="hljs-name" style="color: #7ee787;">')
          .replace(/<span class="hljs-built_in">/g, '<span class="hljs-built_in" style="color: #ffa657;">')
          .replace(/<span class="hljs-literal">/g, '<span class="hljs-literal" style="color: #79c0ff;">')
          .replace(/<span class="hljs-type">/g, '<span class="hljs-type" style="color: #ffa657;">')
          .replace(/<span class="hljs-symbol">/g, '<span class="hljs-symbol" style="color: #79c0ff;">')
          .replace(/<span class="hljs-class">/g, '<span class="hljs-class" style="color: #ffa657;">')
          .replace(/<span class="hljs-operator">/g, '<span class="hljs-operator" style="color: #ff7b72;">')
          .replace(/<span class="hljs-punctuation">/g, '<span class="hljs-punctuation" style="color: #e6edf3;">')
          .replace(/<span class="hljs-regexp">/g, '<span class="hljs-regexp" style="color: #7ee787;">')
          .replace(/<span class="hljs-selector-tag">/g, '<span class="hljs-selector-tag" style="color: #7ee787;">')
          .replace(/<span class="hljs-selector-class">/g, '<span class="hljs-selector-class" style="color: #ffa657;">')
          .replace(/<span class="hljs-selector-id">/g, '<span class="hljs-selector-id" style="color: #79c0ff;">')
          .replace(/<span class="hljs-property">/g, '<span class="hljs-property" style="color: #79c0ff;">')
          .replace(/<span class="hljs-value">/g, '<span class="hljs-value" style="color: #a5d6ff;">');
      } catch (e) {
        console.error('Highlighting error:', e);
        return code;
      }
    }
    return code;
  };


  // Function to copy code to clipboard
  const copyToClipboard = async (text) => {
    try {
      await navigator.clipboard.writeText(text);
      toast.success('Code copied to clipboard!');
    } catch (err) {
      console.error('Failed to copy:', err);
      toast.error('Failed to copy code');
    }
  };

  // Function to add message to chat
  const addMessageToChat = useCallback((message, isUser = false, isThinking = false) => {
    setChatMessages(prev => [...prev, {
      id: Date.now() + Math.random(),
      content: message,
      isUser,
      isThinking,
      timestamp: new Date().toLocaleTimeString()
    }]);
  }, []);

  // Function to update thinking message with actual response
  const updateThinkingMessage = useCallback((responseContent) => {
    setChatMessages(prev => 
      prev.map(msg => 
        msg.isThinking ? { ...msg, content: responseContent, isThinking: false } : msg
      )
    );
  }, []);

  // Function to set API headers
  const buildApiHeaders = useCallback((modelNodeId) => {
    const headers = {
      "Accept": "application/x-ndjson",
      "Content-Type": "application/x-ndjson",
    };

    if (modelNodeId) {
      headers["x-aimodelnode"] = modelNodeId;
    }
    
    return headers;
  }, []);
  
  // Function to update headers when model changes
  const updateApiHeaders = useCallback((modelNodeId) => {
    const newHeaders = buildApiHeaders(modelNodeId);
    setApiHeaders(newHeaders);
    console.log("API Headers updated:", newHeaders);
  }, [buildApiHeaders]);

  // Modal handlers
  const openToolCallModal = () => setIsToolCallModalOpen(true);
  const closeToolCallModal = () => setIsToolCallModalOpen(false);
  const toggleContextPopup = () => setIsContextModalOpen(!isContextModalOpen);

  // set toolcall
  const handleAddTool = (toolData) => {
    setFormData((prevData) => ({
      ...prevData,
      spec: {
        ...prevData.spec,
        Tools: [toolData],
      },
    }));
    toast.success("Tool added successfully");
  };

  // set modelNodeId
  const handleModelSelect = ({ modelNodeId, modelName }) => {
    console.log("Model received:", modelNodeId, modelName);
    updateApiHeaders(modelNodeId);
  };

  // set promptId
  const handlePromptSelection = (promptId) => {
    console.log('prompt received:', promptId);
  };

 // Initialize chat mode
useEffect(() => {
  if (!aiStudioChat) {
    const chatModeRadio = document.querySelector(
      'input[name="aiInteractionMode"][value="CHAT_MODE"]'
    );
    if (chatModeRadio) {
      chatModeRadio.checked = true;
    }
    handleInteractionModeChange("CHAT_MODE");
  }
}, []); 

  // Create payload  
    const createPayload = useCallback(() => {
    const inputArea = document.getElementById("input");
    const promptInput = document.getElementById("promptInput");
    const aiModel = document.getElementById("aiModel").value;
    const temperature = document.getElementById("temperatureValue").value;
    const maxTokens = document.getElementById("maxTokens").value;
    const topk = document.getElementById("topk").value;
    const topP = document.getElementById("topP").value;
    const promptId = document.getElementById("promptId").value;
    const currentChatIdFromDOM = document.getElementById("currentChatId").value;
    const sysMess = document.getElementById("systemMessage"); 

    // Get current interaction mode
    const interactionMode = document.querySelector(
      'input[name="aiInteractionMode"]:checked'
    )?.value || "P2P";

    // Get stream toggle state
    const streamButton = document.getElementById("toggleStreamButton");
    const isstream = streamButton?.getAttribute("aria-checked") === "true";

    // Determine input source
    const input = inputArea?.value || "";
    const promptInputValue = promptInput?.value || "";

    // Build payload with corrected structure
    let payload = {
      messages: [],
    };

    // Add promptId and promptInput if available
    if (promptId) {
      payload.promptId = promptId;
      if (promptInputValue !== "") {
        try {
          payload.promptInput = JSON.parse(promptInputValue);
        } catch (error) {
          payload.promptInput = promptInputValue;
        }
      }
    }

    // Add tools if available
    if (formData.spec.Tools && formData.spec.Tools.length > 0) {
      payload.tools = formData.spec.Tools;
    }

    // Add user input message
    if (input) {
      if (input !== "") {
        payload.messages.push({
          role: "user",
          content: input
        });
      } else {
        toast.error("Please provide a valid input for the AI model.");
        return null;
      }
    }
    
    // Add system message if available
    const systemMessageContent = sysMess?.value || "";
    payload.messages.push({
      role: "system",
      content: systemMessageContent
    });

    // Add temperature
    if (temperature) {
      payload.temperature = parseFloat(temperature);
    }

    // Add topP
    if (topP) {
      payload.topP = parseFloat(topP);
    }

    // Add maxTokens
    if (maxTokens) {
      payload.max_output_tokens = Number(maxTokens);
    }

    // Add model information
    if (!aiModel) {
      toast.error("Please select an AI model");
      return null;
    }

    const [modelNodeId, modelName] = aiModel.split("||");
    
    if (!modelName || !modelNodeId) {
      toast.error("Invalid model selection");
      return null;
    }

    if (modelName) {
      payload.model = modelName;
    }

    // Add modelNodeId
    if (modelNodeId) {
      payload.modelNodeId = modelNodeId;
      updateApiHeaders(modelNodeId);
    }

    // Set mode based on interaction mode
    payload.mode = interactionMode;

    // Add chat_id
    const chatId = currentChatIdFromDOM || currentChatId;
    if (chatId && chatId !== "") {
      payload.chat_id = chatId;
    }

    // Add stream setting
    payload.stream = isstream;

    // Add topK if available
    if (topk) {
      payload.topK = parseInt(topk);
    }

    return payload;
  }, [formData.spec.Tools, currentChatId, updateApiHeaders]);

  // Submit API call 
  const submitForm = async (payload) => {
    try {
      setIsLoading(true);
      setIsStreaming(true);

      // Get user input before clearing
      const inputArea = document.getElementById("input");
      const userMessage = inputArea?.value || "";

      // Add user message to chat
      if (userMessage) {
        addMessageToChat(userMessage, true);
      }

      // Add thinking indicator
      addMessageToChat("Thinking...", false, true);
      document.getElementById("input").value = "";
      document.getElementById("promptInput").value = "";
      document.getElementById("input").classList.remove("hidden");
      document.getElementById("promptInput").classList.add("hidden");

      let result;

      // Prepare headers for API calls
      const requestHeaders = {
        headers: apiHeaders
      };

      console.log("Sending API request with headers:", requestHeaders);

      // Conditional API calls based on mode
      if (payload.mode === 'CHAT_MODE') {
        result = await AiChatApi.getAiChat(payload, requestHeaders);
        console.log("Chat mode result:", result);
      } else {
        if (payload.stream) {
          result = await AiGatewayChatApi.getAiGatewayChat(payload, requestHeaders);
          console.log("Stream mode result:", result);
        } else {
          result = await AiGatewayChatApi.getAiGatewayChat(payload, requestHeaders);
          console.log("Api mode result:", result);
        }
      }

      console.log("API result:", result);

      // Process the response 
      if (result && result.result && result.result.choices && result.result.choices.length > 0) {
        const assistantMessage = result.result.choices[0].messages;
        if (assistantMessage && assistantMessage.content) {
          updateThinkingMessage(assistantMessage.content);
        } else {
          updateThinkingMessage("Sorry, I couldn't generate a response.");
          toast.error("No content in response");
        }
      } else {
        updateThinkingMessage("Sorry, I couldn't generate a response.");
        toast.error("No response received from API");
      }

      clearInputs();

    } catch (error) {
      console.error("Error sending message:", error);
      updateThinkingMessage("Sorry, there was an error processing your request.");
      toast.error(`Failed to send message: ${error.message}`);
    } finally {
      setIsLoading(false);
      setIsStreaming(false);
      
      const inputArea = document.getElementById("input");
      if (inputArea) {
        inputArea.disabled = false;
    }
  }
};

  // Clear input fields after submission
  const clearInputs = () => {
    const inputArea = document.getElementById("input");
    const promptInputArea = document.getElementById("promptInput");
    
    if (inputArea) {
      inputArea.value = "";
      inputArea.disabled = false;
    }
    
    if (promptInputArea) {
      promptInputArea.value = "";
    }
  };

  // submit function
const submitInput = async () => {
  if (isLoading) {
    toast.info("Please wait for the current request to complete");
    return;
  }

  const payload = createPayload();
  if (!payload) return;

  // Disable input during processing
  const inputArea = document.getElementById("input");
  if (inputArea) {
    inputArea.disabled = true;
  }

  try {
    // Check if we need to create a new chat first
    const currentChatIdFromDOM = document.getElementById("currentChatId")?.value;
    const activeChatId = currentChatIdFromDOM || currentChatId;
    
    // If no active chat ID exists, create a new chat first
    if (!activeChatId || activeChatId === "") {
      console.log("No active chat found, creating new chat first...");
      
      // Show loading state for chat creation
      toast.info("Creating new chat...");
      
      try {
        // Create new chat
        const newChatId = await createNewChatAndGetId();
        
        if (!newChatId) {
          throw new Error("Failed to create new chat");
        }
        
        // Update payload with new chat ID
        payload.chat_id = newChatId;
        
        console.log("New chat created with ID:", newChatId);
        
      } catch (chatCreationError) {
        console.error("Error creating new chat:", chatCreationError);
        toast.error("Failed to create new chat. Please try again.");
        
        // Re-enable input
        if (inputArea) {
          inputArea.disabled = false;
        }
        return;
      }
    }

    console.log("Payload being sent:", payload);
    console.log("Headers being sent:", apiHeaders);
    
    // Now submit the message
    await submitForm(payload);
    
  } catch (error) {
    console.error("Error in submitInput:", error);
    toast.error(`Failed to process request: ${error.message}`);
    
    // Re-enable input
    if (inputArea) {
      inputArea.disabled = false;
    }
  }
};

  // Handle interaction mode changes
  function handleInteractionModeChange(divId) {
    document.getElementsByName("aiInteractionMode").forEach((element) => {
      if (element.value === divId) {
        document.getElementById(divId)?.classList.remove("hidden");
      } else {
        document.getElementById(element.value)?.classList.add("hidden");
      }
    });
  }

  const createNewChatAndGetId = async () => {
  try {
    setIsLoading(true);
    const payload = { action: "CREATE" };
    
    const requestHeaders = {
      headers: apiHeaders
    };

    const createChat = await createNewChatApi.getCreateChat(payload, requestHeaders);

    console.log("Chat created:", createChat);

    if (createChat?.output && createChat.output.length > 0) {
      const chatInfo = createChat.output[0];
      const chatId = chatInfo.chatId;

      if (chatId) {
        // Update state and DOM
        setCurrentChatId(chatId);
        const chatIdInput = document.getElementById("currentChatId");
        if (chatIdInput) {
          chatIdInput.value = chatId;
        }

        // Clear chat messages for new chat
        setChatMessages([]);

        // Update URL
        const currentBaseUrl = globalContext?.CurrentUrl || window.location.pathname;
        const newUrl = `${currentBaseUrl}?aiChatId=${chatId}`;
        
        try {
          // Update URL without waiting for navigation
          window.history.pushState({ chatId }, '', newUrl);
        } catch (urlError) {
          console.warn("Failed to update URL:", urlError);
        }

        return chatId;
      } else {
        throw new Error("Failed to get chat ID from response");
      }
    } else {
      throw new Error("Invalid response format from chat creation API");
    }
  } catch (error) {
    console.error("Error creating new chat:", error);
    throw error; 
  } finally {
    setIsLoading(false);
  }
};

  // Create new chat 
  const createNewChat = async () => {
  try {
    const chatId = await createNewChatAndGetId();
    toast.success("New chat created successfully");
    
    // Navigate to new chat URL after a short delay
    setTimeout(() => {
      const currentBaseUrl = globalContext?.CurrentUrl || window.location.pathname;
      const newUrl = `${currentBaseUrl}?aiChatId=${chatId}`;
      
      try {
        router.push(newUrl);
      } catch {
        window.location.href = newUrl;
      }
    }, 1000);
    
  } catch (error) {
    console.error("Error creating new chat:", error);
    toast.error("Failed to create new chat. Please try again.");
  }
};

  // set the chat history
  useEffect(() => {
    const fetchChatHistory = async () => {
      try {
        setIsLoading(true);
        setError(null);
        
        // Fetch chat history using your API
        const data = await chatHistoryApi.getChatHistory();       
        // Extract the output array from the response
        const historyOutput = data?.output || [];
        
        if (historyOutput.length > 0) {
          setChatHistory(historyOutput);
        } else {
          setChatHistory([]);
          console.warn("No chat history found in response");
        }
      } catch (error) {
        console.error("Failed to fetch chat history:", error);
        setError(error.message);
        setChatHistory([]);
      } finally {
        setIsLoading(false);
      }
    };
  
    fetchChatHistory();
  }, []);

  // set tha chat history data on msg area 
  const loadChatData = async (chatId) => {
   try {
     setIsLoading(true);
     setError(null);
     
     // Call the API to get chat data
     const chatData = await chatHistoryApi.getChatData(chatId);
     
     if (chatData?.output && chatData.output.length > 0) {
       const chatInfo = chatData.output[0];
       const messages = chatInfo.messages || [];
       
       // Transform API messages to chat message format
       const transformedMessages = messages.map((msg, index) => ({
         id: Date.now() + index + Math.random(),
         content: msg.content,
         isUser: msg.role === 'user',
         isThinking: false,
         timestamp: new Date().toLocaleTimeString()
       }));
       
       setChatMessages(transformedMessages);
       setCurrentChatId(chatId);
       const chatIdInput = document.getElementById("currentChatId");
       if (chatIdInput) {
         chatIdInput.value = chatId;
       }
       
     } else {
       toast.error("No chat data found");
       setChatMessages([]);
     }
   } catch (error) {
     console.error("Error loading chat data:", error);
     setError(error.message);
     toast.error(`Failed to load chat: ${error.message}`);
     setChatMessages([]);
   } finally {
     setIsLoading(false);
   }
 };

  // function to handle chat history Onclick
  const handleChatHistoryClick = async (chatId, event) => {
    event.preventDefault();
  
    await loadChatData(chatId);
    
    // Optional: Update URL without full page redirect
    const currentBaseUrl = globalContext?.CurrentUrl || window.location.pathname;
    const newUrl = `${currentBaseUrl}?chat_id=${chatId}`;
    
    // Update browser history without reload
    window.history.pushState({ chatId }, '', newUrl);
  };

  // Load chat data if chat_id is in URL
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const chatIdFromUrl = urlParams.get('chat_id');
    
    if (chatIdFromUrl && chatIdFromUrl !== currentChatId) {
     
      loadChatData(chatIdFromUrl);
    }
  }, []); 
  
  // Close options bar
  function closeOptionsBar() {
    const rightbar = document.getElementById("rightbar");
    rightbar?.classList.add("hidden");
  }

  // Temperature controls
  function updateTemperatureValue(value) {
    const tempValue = document.getElementById("temperatureValue");
    if (tempValue) {
      tempValue.value = parseFloat(value).toFixed(1);
    }
  }

  function updateSliderValue(value) {
    const numericValue = parseFloat(value);
    if (numericValue >= 0.0 && numericValue <= 1.0) {
      const tempSlider = document.getElementById("temperatureSlider");
      if (tempSlider) {
        tempSlider.value = numericValue.toFixed(1);
      }
    }
  }

  // Handle Enter key press
  function EnterInput(event) {
    if (event.key === "Enter" && !event.shiftKey) {
      event.preventDefault();
      submitInput();
    }
  }

  // Toggle streaming
  function toggleStream(element) {
    const button = element.tagName === "BUTTON" ? element : element.closest("button");
    if (!button) return;
    
    const isOn = button.getAttribute("aria-checked") === "true";
    const newState = !isOn;
    
    button.classList.toggle("bg-orange-700", newState);
    button.classList.toggle("bg-gray-300", !newState);
    
    const indicator = button.querySelector("span");
    if (indicator) {
      indicator.style.transform = newState ? "translateX(32px)" : "translateX(0)";
    }
    
    const label = document.getElementById("toggleStateLabel");
    if (label) {
      label.textContent = newState ? "Stream ON" : "Stream OFF";
    }
    
    button.setAttribute("aria-checked", newState.toString());
  }

  // Cancel streaming
  const cancelStream = () => {
    setIsStreaming(false);
    setIsLoading(false);
    
    // Update any thinking messages
    setChatMessages(prev => 
      prev.map(msg => 
        msg.isThinking ? { ...msg, content: "Request cancelled", isThinking: false } : msg
      )
    );
    
    // Re-enable input
    const inputArea = document.getElementById("input");
    if (inputArea) {
      inputArea.disabled = false;
    }
    
    toast.info("Stream cancelled");
  };

  return (
    <>
      <Head>
        <meta charSet="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>AI Studio</title>

        {/* External Stylesheets */}
        <link
          href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism.min.css"
          rel="stylesheet"
        />
        <link
          href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-okaidia.min.css"
          rel="stylesheet"
        />
        <link
          href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/styles/github-dark.min.css"
          rel="stylesheet"
        />
        <link
          rel="stylesheet"
          href="https://cdnjs.cloudflare.com/ajax/libs/jsoneditor/9.10.0/jsoneditor.min.css"
        />

        <style jsx global>{`
        /* Global Modal Styles */
           .modal {
             display: none;
             position: fixed;
             top: 0;
             left: 0;
             width: 100%;
             height: 100%;
             background: rgba(0, 0, 0, 0.5);
             align-items: center;
             justify-content: center;
             z-index: 1000;
           }
           
           .modal.active {
             display: flex;
           }
           
           /* Sidebar Animation */
           .rightbar {
             transition: transform 0.3s ease-in-out;
           }
           
           .rightbar.hidden {
             transform: translateX(100%);
           }
           
           .rightbar.visible {
             transform: translateX(0);
           }
           
           /* Loading Animations */
           .loading-spinner {
             animation: spin 1s linear infinite;
           }
           
           @keyframes spin {
             from { transform: rotate(0deg); }
             to { transform: rotate(360deg); }
           }
           
           /* Pulse animation for thinking indicator */
           .thinking-pulse {
             animation: pulse 1.5s ease-in-out infinite;
           }
           
           @keyframes pulse {
             0%, 100% { opacity: 1; }
             50% { opacity: 0.5; }
           }
           
           /* Enhanced Syntax Highlighting Styles */
           .syntax-highlight-container .hljs {
             background: #0d1117 !important;
             color: #e6edf3 !important;
             font-family: 'SF Mono', 'Monaco', 'Inconsolata', 'Roboto Mono', 'Source Code Pro', monospace;
             font-size: 14px;
             line-height: 1.5;
             border-radius: 6px;
           }
           
           /* Keywords (if, else, function, const, let, var, etc.) */
           .syntax-highlight-container .hljs-keyword {
             color: #ff7b72 !important;
             font-weight: 600;
           }
           
           /* Strings */
           .syntax-highlight-container .hljs-string {
             color: #a5d6ff !important;
           }
           
           /* Numbers */
           .syntax-highlight-container .hljs-number {
             color: #79c0ff !important;
           }
           
           /* Comments */
           .syntax-highlight-container .hljs-comment {
             color: #8b949e !important;
             font-style: italic;
           }
           
           /* Function names */
           .syntax-highlight-container .hljs-function,
           .syntax-highlight-container .hljs-title.function_ {
             color: #d2a8ff !important;
           }
           
           /* Variables */
           .syntax-highlight-container .hljs-variable {
             color: #ffa657 !important;
           }
           
           /* Class names, HTML tags */
           .syntax-highlight-container .hljs-title,
           .syntax-highlight-container .hljs-tag,
           .syntax-highlight-container .hljs-name {
             color: #7ee787 !important;
           }
           
           /* HTML Attributes */
           .syntax-highlight-container .hljs-attr {
             color: #79c0ff !important;
           }
           
           /* Built-in functions and types */
           .syntax-highlight-container .hljs-built_in,
           .syntax-highlight-container .hljs-type {
             color: #ffa657 !important;
           }
           
           /* Literals (true, false, null, undefined) */
           .syntax-highlight-container .hljs-literal {
             color: #79c0ff !important;
           }
           
           /* Symbols and operators */
           .syntax-highlight-container .hljs-symbol,
           .syntax-highlight-container .hljs-operator {
             color: #ff7b72 !important;
           }
           
           /* Punctuation */
           .syntax-highlight-container .hljs-punctuation {
             color: #e6edf3 !important;
           }
           
           /* Regular expressions */
           .syntax-highlight-container .hljs-regexp {
             color: #7ee787 !important;
           }
           
           /* CSS Selectors */
           .syntax-highlight-container .hljs-selector-tag {
             color: #7ee787 !important;
           }
           
           .syntax-highlight-container .hljs-selector-class {
             color: #ffa657 !important;
           }
           
           .syntax-highlight-container .hljs-selector-id {
             color: #79c0ff !important;
           }
           
           /* CSS Properties and Values */
           .syntax-highlight-container .hljs-property {
             color: #79c0ff !important;
           }
           
           .syntax-highlight-container .hljs-value {
             color: #a5d6ff !important;
           }
           
           /* Additional language-specific highlighting */
           .syntax-highlight-container .hljs-meta {
             color: #8b949e !important;
           }
           
           .syntax-highlight-container .hljs-doctag {
             color: #8b949e !important;
           }
           
           .syntax-highlight-container .hljs-section {
             color: #1f6feb !important;
             font-weight: bold;
           }
           
           .syntax-highlight-container .hljs-addition {
             color: #aff5b4 !important;
             background-color: #033a16 !important;
           }
           
           .syntax-highlight-container .hljs-deletion {
             color: #ffdcd7 !important;
             background-color: #67060c !important;
           }
           
           .syntax-highlight-container .hljs-emphasis {
             font-style: italic;
           }
           
           .syntax-highlight-container .hljs-strong {
             font-weight: bold;
           }
           
           /* Chat Message Styles */
           .chat-message {
             margin-bottom: 16px;
             animation: slideIn 0.3s ease-out;
           }
           
           @keyframes slideIn {
             from {
               opacity: 0;
               transform: translateY(10px);
             }
             to {
               opacity: 1;
               transform: translateY(0);
             }
           }
           
           .chat-message.user {
             text-align: right;
           }
           
           .chat-message.assistant {
             text-align: left;
           }
           
           /* Input Field Styles */
           .input-disabled {
             opacity: 0.6;
             cursor: not-allowed;
           }
           
           /* Button Hover Effects */
           .button-hover-effect {
             transition: all 0.2s ease-in-out;
           }
           
           .button-hover-effect:hover {
             transform: translateY(-1px);
             box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
           }
           
           /* Stream Toggle Button */
           .stream-toggle {
             transition: all 0.3s ease;
           }
           
           .stream-toggle .indicator {
             transition: transform 0.3s ease;
           }
           
           /* Code Block Copy Button */
           .copy-button {
             opacity: 0.7;
             transition: opacity 0.2s ease;
           }
           
           .copy-button:hover {
             opacity: 1;
           }
           
           /* Chat History Styles */
           .chat-history-item {
             transition: background-color 0.2s ease;
             cursor: pointer;
             border-radius: 8px;
             padding: 8px 12px;
             margin-bottom: 4px;
           }
           
           .chat-history-item:hover {
             background-color: rgba(156, 163, 175, 0.1);
           }
           
           .chat-history-item.active {
             background-color: rgba(59, 130, 246, 0.1);
             border-left: 3px solid #3b82f6;
           }
           
           /* Toast Notification Styles (if using react-toastify) */
           .Toastify__toast {
             border-radius: 8px;
           }
           
           .Toastify__toast--success {
             background: linear-gradient(135deg, #10b981, #059669);
           }
           
           .Toastify__toast--error {
             background: linear-gradient(135deg, #ef4444, #dc2626);
           }
           
           .Toastify__toast--info {
             background: linear-gradient(135deg, #3b82f6, #2563eb);
           }
           
           /* Dropdown Styles */
           .dropdown-menu {
             box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
             border-radius: 8px;
             border: 1px solid rgba(229, 231, 235, 0.2);
           }
           
           .dropdown-item {
             transition: background-color 0.15s ease;
           }
           
           .dropdown-item:hover {
             background-color: rgba(156, 163, 175, 0.1);
           }
           
           /* Tool Call Modal Styles */
           .tool-modal {
             backdrop-filter: blur(4px);
           }
           
           .tool-modal .modal-content {
             background: white;
             border-radius: 12px;
             box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
             max-height: 90vh;
             overflow-y: auto;
           }
           
           /* Context Modal Styles */
           .context-modal .modal-content {
             background: linear-gradient(135deg, #f8fafc, #f1f5f9);
             border-radius: 12px;
             box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
           }
           
           /* Scrollbar Styles */
           .custom-scrollbar {
             scrollbar-width: thin;
             scrollbar-color: rgba(156, 163, 175, 0.5) transparent;
           }
           
           .custom-scrollbar::-webkit-scrollbar {
             width: 6px;
           }
           
           .custom-scrollbar::-webkit-scrollbar-track {
             background: transparent;
           }
           
           .custom-scrollbar::-webkit-scrollbar-thumb {
             background-color: rgba(156, 163, 175, 0.5);
             border-radius: 3px;
           }
           
           .custom-scrollbar::-webkit-scrollbar-thumb:hover {
             background-color: rgba(156, 163, 175, 0.7);
           }
           
           /* Temperature Slider */
           .temperature-slider {
             appearance: none;
             background: linear-gradient(to right, #3b82f6, #8b5cf6, #ef4444);
             border-radius: 5px;
             height: 8px;
             outline: none;
           }
           
           .temperature-slider::-webkit-slider-thumb {
             appearance: none;
             background: white;
             border: 2px solid #3b82f6;
             border-radius: 50%;
             cursor: pointer;
             height: 20px;
             width: 20px;
           }
           
           .temperature-slider::-moz-range-thumb {
             background: white;
             border: 2px solid #3b82f6;
             border-radius: 50%;
             cursor: pointer;
             height: 16px;
             width: 16px;
           }
           
           /* Responsive Design */
           @media (max-width: 768px) {
             .syntax-highlight-container .hljs {
               font-size: 12px;
             }
             
             .chat-message {
               margin-bottom: 12px;
             }
             
             .rightbar {
               width: 100% !important;
             }
           }
           
           @media (max-width: 640px) {
             .modal .modal-content {
               margin: 20px;
               max-height: calc(100vh - 40px);
             }
             
             .syntax-highlight-container pre {
               padding: 12px !important;
             }
           }
           
           /* Focus States */
           input:focus,
           textarea:focus,
           select:focus {
             outline: 2px solid #3b82f6;
             outline-offset: 2px;
           }
           
           /* Accessibility */
           .sr-only {
             position: absolute;
             width: 1px;
             height: 1px;
             padding: 0;
             margin: -1px;
             overflow: hidden;
             clip: rect(0, 0, 0, 0);
             white-space: nowrap;
             border: 0;
           }
           
           /* Dark mode support */
           @media (prefers-color-scheme: dark) {
             .tool-modal .modal-content {
               background: #1f2937;
               color: #f3f4f6;
             }
             
             .context-modal .modal-content {
               background: linear-gradient(135deg, #1f2937, #111827);
               color: #f3f4f6;
             }
           }
        
        `}</style>
        </Head>
         

      {/* External Script Tags */}
      <Script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/jsoneditor/9.10.0/jsoneditor.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/highlight.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/languages/go.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/prism.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-python.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-javascript.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-sql.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-go.min.js" />
      <Script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/plugins/autoloader/prism-autoloader.min.js" />

       <div className="bg-zinc-800 flex h-screen">
        <div className="h-screen w-full">
          <Header
            sectionHeader="AI Studio Interface"
            hideBackListingLink={true}
            backListingLink="./"
          />

          <button
            id="optionsIcon"
            className="fixed bottom-4 right-4 bg-orange-700 p-3 rounded-full shadow-lg focus:outline-none z-50 cursor-pointer hover:bg-orange-600 transition-colors"
            onClick={() => {
              const rightbar = document.getElementById("rightbar");
              const isHidden = rightbar?.classList.contains("hidden");
              rightbar?.classList.toggle("hidden", !isHidden);
              rightbar?.classList.toggle("visible", isHidden);
            }}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth="2"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </button>

          {/* Chat and Editor Section */}
          <div id="aistudiodiv" className="flex flex-1">
            <div
              className="flex bg-zinc-800 p-2 flex-col w-full overflow-y-auto overflow-x-auto scrollbar"
              style={{ height: "calc(100vh - 65px)" }}
            >
              {/* Chat Messages Display Area */}
              <div
                id="aiChatResult"
                className="ml-8 mr-8 flex-grow rounded-lg bg-gray-50 table-container overflow-y-auto scrollbar"
                style={{
                  height: "500px",
                  background:
                    "radial-gradient(circle, black -25%, transparent 10%) 0 0, radial-gradient(circle, black -25%, transparent 10%) 10px 10px",
                  backgroundSize: "20px 20px",
                }}
              >
                {chatMessages.length === 0 ? (
                  <div className="flex items-center justify-center h-full text-gray-500">
                    <p>Start a conversation...</p>
                  </div>
                ) : (
                  <div className="space-y-4 p-4">
                    {chatMessages.map((message) => (
                       <div
                         key={message.id}
                         className={`flex ${message.isUser ? 'justify-end' : 'justify-start'}`}
                       >
                         <div 
                           className={`inline-block max-w-[85%] px-4 py-2 rounded-3xl ${
                             message.isUser
                               ? 'bg-[#1b1b1b] text-white'
                               : message.isThinking
                               ? 'bg-gray-300 text-gray-700 thinking-animation'
                               : 'text-gray-100 bg-transparent'
                           }`}
                           style={{
                             wordWrap: 'break-word',
                             wordBreak: 'break-word',
                             overflowWrap: 'break-word'
                           }}
                         >
                           {message.isThinking ? (
                             <p className="text-sm whitespace-pre-wrap">{message.content}</p>
                           ) : (
                             <div className="text-sm">
                               {message.isUser ? (
                                 <p className="whitespace-pre-wrap">{message.content}</p>
                               ) : (
                                 <div className="prose prose-invert max-w-none">
                                   {formatMessageContent(message.content)}
                                 </div>
                               )}
                             </div>
                           )}
                         </div>
                       </div>
                     ))}
                  </div>
                )}
              </div>

              {/* Input Area with Dropdown */}
              <div
                id="userInputGround"
                className="ml-8 mr-8 mt-2 bg-[#1b1b1b] p-1 rounded-lg shadow-md border border-zinc-500">
                <textarea
                  id="input"
                  onKeyDown={EnterInput}
                  onInput={(e) => {
                    // Auto-resize textarea
                    e.target.style.height = 'auto';
                    e.target.style.height = Math.min(e.target.scrollHeight, 300) + 'px';
                  }}
                  disabled={isLoading}
                  className="w-full min-h-[64px] max-h-[300px] p-3 rounded-lg focus:outline-none text-gray-100 overflow-y-auto scrollbar disabled:opacity-50 disabled:cursor-not-allowed resize-none leading-5"
                  style={{
                    height: '64px',
                    lineHeight: '1.25rem'
                  }}
                  placeholder={isLoading ? "Sending message..." : "Type your message here..."}
                  rows="1"
                />
                <textarea
                  id="promptInput"
                  onKeyDown={EnterInput}
                  onInput={(e) => {
                    // Auto-resize textarea
                    e.target.style.height = 'auto';
                    e.target.style.height = Math.min(e.target.scrollHeight, 300) + 'px';
                  }}
                  className="w-full min-h-[64px] max-h-[300px] p-3 rounded-lg focus:outline-none text-gray-100 hidden resize-none leading-5"
                  style={{
                    height: '64px',
                    lineHeight: '1.25rem'
                  }}
                  placeholder="Type your message here..."
                  rows="1"
                />
                <div className="flex h-[40px] items-center justify-between mt-2 mb-2 px-3">
                  {/* toggleContextPopup Button*/}
                  <div className="flex gap-x-2">
                    <button
                      title="View selected dataproduct prompts"
                      aria-label="datafiles"
                      data-testid="datafiles"
                      className="flex h-10 w-10 items-center justify-center rounded-full transition-colors hover:opacity-70 focus-visible:outline-none focus-visible:outline-black bg-orange-700"
                      onClick={toggleContextPopup}
                      disabled={isLoading}
                    >
                      <svg viewBox="0 0 24 24" className="h-6 w-6 shrink-0">
                        <path
                          fill="currentColor"
                          d="M16.5 6v11.5a4 4 0 0 1-4 4a4 4 0 0 1-4-4V5A2.5 2.5 0 0 1 11 2.5A2.5 2.5 0 0 1 13.5 5v10.5a1 1 0 0 1-1 1a1 1 0 0 1-1-1V6H10v9.5a2.5 2.5 0 0 0 2.5 2.5a2.5 2.5 0 0 0 2.5-2.5V5a4 4 0 0 0-4-4a4 4 0 0 0-4 4v12.5a5.5 5.5 0 0 0 5.5 5.5a5.5 5.5 0 0 0 5.5-5.5V6z"
                        />
                      </svg>
                    </button>
                  </div>
              
                  <button
                    aria-label="startChat"
                    data-testid="startChat"
                    id="startChat"
                    className={`flex cursor-pointer h-10 w-10 items-center justify-center rounded-full transition-colors hover:opacity-70 focus-visible:outline-none focus-visible:outline-black disabled:text-[#f4f4f4] disabled:hover:opacity-100 disabled:dark:bg-token-text-quaternary bg-orange-700 disabled:bg-[#D7D7D7] ${
                      isLoading ? "hidden" : ""
                    }`}
                    onClick={submitInput}
                    disabled={isLoading}
                  >
                    <svg
                      viewBox="0 0 32 32"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-8 w-8"
                    >
                      <path
                        fillRule="evenodd"
                        clipRule="evenodd"
                        d="M15.1918 8.90615C15.6381 8.45983 16.3618 8.45983 16.8081 8.90615L21.9509 14.049C22.3972 14.4953 22.3972 15.2189 21.9509 15.6652C21.5046 16.1116 20.781 16.1116 20.3347 15.6652L17.1428 12.4734V22.2857C17.1428 22.9169 16.6311 23.4286 15.9999 23.4286C15.3688 23.4286 14.8571 22.9169 14.8571 22.2857V12.4734L11.6652 15.6652C11.2189 16.1116 10.4953 16.1116 10.049 15.6652C9.60265 15.2189 9.60265 14.4953 10.049 14.049L15.1918 8.90615Z"
                        fill="currentColor"
                      />
                    </svg>
                  </button>
              
                  <button
                    aria-label="inprogressChat"
                    data-testid="inprogressChat"
                    id="inprogressChat"
                    className={`flex h-10 w-10 items-center justify-center rounded-full transition-colors hover:opacity-70 focus-visible:outline-none focus-visible:outline-black disabled:text-[#f4f4f4] disabled:hover:opacity-100 disabled:dark:bg-token-text-quaternary bg-orange-700 disabled:bg-[#D7D7D7] ${
                      isLoading ? "" : "hidden"
                    }`}
                    onClick={cancelStream}
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="h-8 w-8"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    >
                      <rect
                        x="9"
                        y="9"
                        width="6"
                        height="6"
                        fill="currentColor"
                      />
                    </svg>
                  </button>
                </div>
              </div>
            </div>

            {/* right sidebar  */}
           <div className="w-[220px] min-w-[220px] bg-[#1b1b1b] p-2 text-gray-100 flex flex-col" style={{ height: "calc(100vh - 65px)" }}>
              <div className="flex-shrink-0 mb-2 space-y-2 text-gray-100">
                <label
                  htmlFor="aiInteractionMode"
                  className="block font-medium text-sm mb-1"
                >
                  Interaction Mode
                </label>
                <div className="flex items-center space-x-4 mt-2"></div>
                <label className="inline-flex items-center cursor-pointer">
                  <input
                    type="radio"
                    name="aiInteractionMode"
                    value="P2P"
                    className="peer hidden"
                    onClick={() => handleInteractionModeChange("P2P")}
                    defaultChecked
                  />
                  <div className=" w-4 h-4 rounded-full border-2 peer-checked:bg-orange-700 transition" />
                  <span className=" ml-2 mr-2 text-sm">Api</span>
                </label>
                <label className=" inline-flex items-center cursor-pointer">
                  <input
                    type="radio"
                    name="aiInteractionMode"
                    value="CHAT_MODE"
                    className="peer hidden"
                    onClick={() => handleInteractionModeChange("CHAT_MODE")}
                  />
                  <div className="w-4 h-4 rounded-full border-2  peer-checked:bg-orange-700 transition" />
                  <span className="ml-2 text-sm">Chat</span>
                </label>
              </div>
            
              <div className="flex-shrink-0 mb-2 mt-2 text-gray-100">
                <NestedDropdown onModelSelect={handleModelSelect} />
                <input type="text" id="aiModel" hidden />
              </div>
            
              {/* API Mode Section */}
              <div id="P2P" className="flex-shrink-0 text-gray-100">
                <div className="relative">
                  <label
                    htmlFor="toggleStateLabel"
                    id="toggleStateLabel"
                    className="block text-sm font-medium mb-1"
                  >
                    Stream OFF
                  </label>
                  <button
                    id="toggleStreamButton"
                    className="relative cursor-pointer inline-flex h-6 w-16 items-center rounded-full bg-gray-300 transition-colors focus:outline-none focus:ring-2 focus:ring-[#ff5f1f]"
                    role="switch"
                    aria-checked="false"
                    onClick={(e) => toggleStream(e.target)}
                  >
                    <span className="absolute left-1 w-6 h-6 bg-white rounded-full shadow-md transition-transform transform"></span>
                  </button>
                </div>
            
                <div className="mb-2 mt-2">
                  <PromptDropdown onPromptSelect={handlePromptSelection} />
                  <input type="text" id="promptId" hidden />
                </div>
            
                <div id="toolCall" className="flex mt-2 mb-2">
                  <button
                    id="addToolCallButton"
                    type="button"
                    onClick={openToolCallModal}
                    className="bg-orange-700 w-full px-2 py-2 text-sm rounded-lg focus:outline-none cursor-pointer hover:bg-orange-600 transition-colors"
                    disabled={isLoading}
                  >
                    Add Tool Call
                  </button>
                  <AIToolCallPopup
                    isOpen={isToolCallModalOpen}
                    onClose={closeToolCallModal}
                    onAddTool={handleAddTool}
                  />
                </div>
              </div>
            
              {/* Chat Mode Section */}
              <div id="CHAT_MODE" className="hidden flex-1 flex flex-col min-h-0 overflow-hidden">
                <div className="flex-shrink-0 mb-2 mt-2">
                  <a
                    className="flex p-2 bg-orange-700 hover:bg-zinc-900 items-center justify-between text-sm font-medium transition-colors rounded-lg cursor-pointer"
                    href="#"
                    onClick={createNewChat}
                  >
                    <span className="flex flex-row items-center justify-start">
                      <div slot="icon" className="w-4 text-neutral-white">
                        <svg
                          viewBox="0 0 24 24"
                          className="h-full w-6"
                          stroke="currentColor"
                        >
                          <path
                            fill="currentColor"
                            d="M5 3c-1.11 0-2 .89-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7h-2v7H5V5h7V3zm12.78 1a.7.7 0 0 0-.48.2l-1.22 1.21l2.5 2.5L19.8 6.7c.26-.26.26-.7 0-.95L18.25 4.2c-.13-.13-.3-.2-.47-.2m-2.41 2.12L8 13.5V16h2.5l7.37-7.38z"
                          />
                        </svg>
                      </div>
                      <span className="ml-2 pl-2 justify-center text-sm">
                        Start New Chat
                      </span>
                    </span>
                  </a>
                </div>
            
                {/* Chat History Section */}
                <div className="flex-1 flex flex-col min-h-0 overflow-hidden">
                  <label
                    htmlFor="chatHistory"
                    className="flex-shrink-0 block font-semibold text-sm mb-2"
                  >
                    Chat History
                  </label>
                  <div className="flex-1 overflow-hidden">
                    {isLoading ? (
                      <div className="p-2 text-gray-500">Loading chat history...</div>
                    ) : error ? (
                      <div className="p-2 text-red-500">Error loading chat history: {error}</div>
                    ) : (
                      <div className="h-full overflow-y-auto scrollbar">
                        <ul className="space-y-2 pr-1">
                          {chatHistory.length > 0 ? (
                            chatHistory.map((chat, index) =>
                              chat.currentLog && (
                                <li key={chat.chatId || index} className="border-b border-zinc-500">
                                  <a
                                    href="#"
                                    onClick={(e) => handleChatHistoryClick(chat.chatId, e)}
                                    className="block p-2 text-gray-100 hover:bg-zinc-900 hover:text-white cursor-pointer transition-colors rounded"
                                  >
                                    <p className="font-semibold text-xs break-words">
                                      {chat.currentLog.value.length > 35
                                        ? chat.currentLog.value.substring(0, 35) + "..."
                                        : chat.currentLog.value}
                                    </p>
                                    <p className="text-xs text-gray-400 mt-1">
                                      Started: {new Date(parseInt(chat.startedAt) * 1000).toLocaleString()}
                                    </p>
                                  </a>
                                </li>
                              )
                            )
                          ) : (
                            <li className="p-2 bg-zinc-900 shadow-sm rounded-md">
                              <p className="font-semibold text-xs break-words">No history available</p>
                            </li>
                          )}
                        </ul>
                      </div>
                    )}
                  </div>
                </div>                  
              </div>
            </div>

            {/* temperature menu */}
            <div
              id="rightbar"
              className="rightbar fixed top-20 right-0 w-64 h-full bg-[#1b1b1b] shadow-lg hidden rounded-lg text-gray-100"
            >
              <div className="p-4 bg-zinc-900 flex justify-between items-center rounded-md">
                <h2 className="text-lg font-bold">Sidebar</h2>
                <button
                  id="closeSidebar"
                  onClick={closeOptionsBar}
                  className="hover:text-gray-100 focus:outline-none cursor-pointer"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    strokeWidth="2"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                </button>
              </div>

              <div className="p-4">
                <label
                  htmlFor="temperatureSelect"
                  className="block text-sm font-medium mb-1"
                >
                  Select Temperature
                </label>
                <div
                  id="temperatureSelect"
                  className="flex items-center space-x-4 mt-2"
                >
                  <input
                    id="temperatureSlider"
                    type="range"
                    min="0.0"
                    max="1.0"
                    step="0.1"
                    defaultValue="0.7"
                    className="w-full h-2 bg-orange-700 rounded-lg appearance-none cursor-pointer"
                    onChange={(e) => updateTemperatureValue(e.target.value)}
                  />
                  <input
                    id="temperatureValue"
                    type="number"
                    step="0.1"
                    min="0.0"
                    max="1.0"
                    defaultValue="0.7"
                    className="w-16 px-2 py-1 border border-gray-300 rounded-lg text-center focus:outline-none focus:ring focus:ring-orange-700"
                    onChange={(e) => updateSliderValue(e.target.value)}
                  />
                </div>

                <label
                  htmlFor="maxTokens"
                  className="block text-sm font-medium mt-4 mb-1"
                >
                  Select Max Tokens
                </label>
                <input
                  type="text"
                  id="maxTokens"
                  placeholder="Max Tokens..."
                  className="w-full px-3 py-2 mb-1 text-gray-100 shadow-md rounded-lg focus:outline-none focus:ring text-sm focus:ring-orange-700"
                />

                <label
                  htmlFor="topk"
                  className="block text-sm font-medium mt-4 mb-1"
                >
                  Select Top K
                </label>
                <input
                  type="text"
                  id="topk"
                  placeholder="Top K value between 0 - 1000"
                  className="w-full px-3 py-2 text-gray-100 mb-1 shadow-md rounded-lg focus:outline-none focus:ring text-sm focus:ring-orange-700"
                />

                <label
                  htmlFor="topP"
                  className="block text-sm font-medium mt-4 mb-1"
                >
                  Select Top P
                </label>
                <input
                  type="text"
                  id="topP"
                  placeholder="Top P value between 0.0 - 1.0"
                  className="w-full px-3 py-2 text-gray-100 mb-1 shadow-md rounded-lg focus:outline-none focus:ring text-sm focus:ring-orange-700"
                />
              </div>
            </div>

            <input
              type="text"
              id="currentChatId"
              hidden
            />

            <div id="currentChat" className="hidden">
              {aiStudioChat ? JSON.stringify(aiStudioChat) : ""}
            </div>
          </div>
        </div>

        {/* Context Popup Component */}
        <ContextModal 
          isOpen={isContextModalOpen} 
          onClose={() => setIsContextModalOpen(false)} 
        />
      </div>
    </>
  );
}
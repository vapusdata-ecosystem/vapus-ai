// import { addUserInputBox } from "./chat-interface.js";
import { addUserInputBox } from "./chat-interface.js";
import { LocalDb } from "./indexdb.js";
let CurrentAbortController = null;

function escapeHTML(str) {
  str = str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
  return decodeURIComponent(originalString);
}

export async function aiInterfaceAction(
  apiUrl,
  chatAPIUrl,
  tokenKey,
  resultDiv,
  contextType,
  contextValue,
  modelNodeId,
  promptId,
  input,
  temperature,
  topP,
  modelName,
  maxTokens,
  topK
) {
  const promptInput = document.getElementById("promptInput").value;
  if (input === "" && promptInput === "") {
    showErrorMessage(
      "User Input Error",
      "Please provide a valid input for the AI model."
    );
    return;
  }
  const isStream = getToggleStreamValue() === "true";
  const reqId = generateUUID();
  const qReqId = `q-${reqId}`;
  const rReqId = `r-${reqId}`;
  const pRReqId = `pr-${reqId}`;
  let dInput = input;
  if (promptInput !== "") {
    dInput = dInput + " " + promptInput;
  }
  let responseElem = `
      <div class="chat-loader" id="loader-${pRReqId}">
        Replying<span>.</span><span>.</span><span>.</span>
      </div>
      <div class="mb-2 flex justify-start hidden" id=${rReqId}>
      <div class="text-left rounded-xl bg-zinc-900 text-gray-100" id=${pRReqId}>
      </div>
      </div>`;
  let currentChatId = "";
  const currentChatObj = document.getElementById("currentChatId");
  if (currentChatObj) {
    currentChatId = currentChatObj.value;
  }
  currentChatId = currentChatId.trim();
  let opBox = document.getElementById(resultDiv);
  addUserInputBox("promptInput", dInput, opBox);
  opBox.innerHTML += responseElem;
  opBox.scrollTop = opBox.scrollHeight;
  const sysMess = document.getElementById("systemMessage");
  const radios = Array.from(document.getElementsByName("aiInteractionMode"));

  // Find the checked radio button
  const modeSel = radios.find((radio) => radio.checked);
  const mode = modeSel.value;

  const myHeaders = new Headers();
  myHeaders.append("Accept", "application/x-ndjson");
  myHeaders.append("Content-Type", "application/x-ndjson");
  myHeaders.append("x-aimodelnode", modelNodeId);
  const apiToken = getCookie(tokenKey);
  myHeaders.append("Authorization", `Bearer ${apiToken}`);
  const fName = document.getElementById("functionName").value;
  const fDescription = document.getElementById("functionDescription").value;
  var toolCallVal = document.getElementById("toolCallSchema").value;
  const toolcall = [];
  if (toolCallVal) {
    const fSchema = JSON.parse(toolCallVal.value);
    if (fName !== "" || fSchema !== "" || fDescription !== "") {
      toolcall.push({
        type: "function",
        functionSchema: {
          name: fName,
          description: fDescription,
          parameters: JSON.stringify(fSchema),
        },
      });
      document.getElementById("functionName").value = "";
      document.getElementById("toolCallSchema").value = "";
      document.getElementById("functionDescription").value = "";
    }
  }
  document.getElementById("startChat").classList.add("hidden");
  document.getElementById("inprogressChat").classList.remove("hidden");

  let payload = {
    messages: [],
  };
  if (promptId) {
    payload.promptId = promptId;
    if (promptInput !== "") {
      payload.promptInput = JSON.parse(promptInput);
    }
  }
  if (toolcall.length > 0) {
    payload.tools = toolcall;
  }
  if (input) {
    if (input !== "") {
      payload.messages.push({
        role: "user",
        content: input,
      });
    } else {
      showAlert(
        AlertError,
        "User Input Error",
        "Please provide a valid input for the AI model."
      );
      return;
    }
  }
  if (temperature) {
    payload.temperature = parseFloat(temperature);
  }
  if (topP) {
    payload.topP = topP;
  }
  if (maxTokens) {
    payload.max_output_tokens = Number(maxTokens);
  }
  if (contextType && contextValue) {
    payload.contexts = [
      {
        key: contextType,
        value: contextValue,
      },
    ];
  }
  if (modelName) {
    payload.model = modelName;
  }
  if (sysMess) {
    payload.messages.push({
      role: "system",
      content: sysMess.value,
    });
  }
  if (modelNodeId) {
    payload.modelNodeId = modelNodeId;
  }
  payload.mode = mode;
  document.getElementById("input").value = "";
  document.getElementById("promptInput").value = "";
  document.getElementById("input").classList.remove("hidden");
  document.getElementById("promptInput").classList.add("hidden");
  try {
    if (mode === "CHAT_MODE") {
      return chatServer(
        payload,
        chatAPIUrl,
        myHeaders,
        resultDiv,
        rReqId,
        pRReqId,
        currentChatId
      );
    } else {
      if (isStream) {
        payload.stream = true;
        return streamServer(
          payload,
          apiUrl,
          myHeaders,
          resultDiv,
          rReqId,
          pRReqId
        );
      }
      return fetchServer(
        payload,
        apiUrl,
        myHeaders,
        resultDiv,
        rReqId,
        pRReqId
      );
    }
  } finally {
  }
}

async function chatServer(
  payload,
  apiUrl,
  myHeaders,
  resultDiv,
  messageDivId,
  pDivId,
  currentChatId
) {
  try {
    CurrentAbortController = new AbortController();
    const signal = CurrentAbortController.signal;
    payload.chat_id = currentChatId;
    const raw = JSON.stringify(payload);
    console.log("apiUrl:");
    console.log(apiUrl);
    console.log(payload);

    const response = await fetch(apiUrl, {
      method: "POST",
      headers: myHeaders,
      body: raw,
      redirect: "follow",
      signal,
    });

    if (!response.body) {
      console.error("ReadableStream not supported in this environment");
      return;
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder("utf-8");
    const streamedContentDiv = document.getElementById(resultDiv);
    const messageDiv = document.getElementById(messageDivId);
    messageDiv.classList.remove("hidden");
    const pDiv = document.getElementById(pDivId);
    let contentpre = document.createElement("div");
    contentpre.classList.add(
      "text-sm",
      "p-2",
      "rounded-lg",
      "mt-2",
      "break-words",
      "whitespace-pre-wrap"
    );
    pDiv.appendChild(contentpre);
    document.getElementById("loader-" + pDivId).style.display = "none";

    let done = false;
    while (!done) {
      const { value, done: readerDone } = await reader.read();
      if (value) {
        let decodedValue = decoder.decode(value);
        decodedValue = decodedValue.trim();
        decodedValue = decodedValue.replace(/}{/g, "},{");
        decodedValue = decodedValue.replace(/}\n{/g, "},{");
        try {
          let strval = "[" + decodedValue + "]";
          let objVal = JSON.parse(strval);
          for (let i = 0; i < objVal.length; i++) {
            objVal[i].result.choices.forEach((val) => {
              contentpre.innerHTML += decodeHTMLEntities(val.messages.content);
            });
            // contentpre.innerHTML += decodeHTMLEntities(objVal[i].result.output.content);
          }
          // contentpre.textContent += escapeHTML(objVal.result.output.content);
        } catch (e) {
          console.log(e);
        }
        // Append content directly to the div
        streamedContentDiv.scrollTop = streamedContentDiv.scrollHeight; // Auto-scroll to the bottom
      }
      done = readerDone;
    }
    contentpre.innerHTML = renderMarkdown(contentpre.innerHTML);
    pDiv.appendChild(contentpre);
    pDiv.innerHTML = formatHtmlContent(pDiv.id);
  } catch (error) {
    console.error("Error fetching data:", error);
  } finally {
    document.getElementById("startChat").classList.remove("hidden");
    document.getElementById("inprogressChat").classList.add("hidden");
    CurrentAbortController = null;
  }
}

async function streamServer(
  payload,
  apiUrl,
  myHeaders,
  resultDiv,
  messageDivId,
  pDivId
) {
  try {
    CurrentAbortController = new AbortController();
    const signal = CurrentAbortController.signal;
    const raw = JSON.stringify(payload);
    const response = await fetch(apiUrl, {
      method: "POST",
      headers: myHeaders,
      body: raw,
      redirect: "follow",
      signal,
    });

    if (!response.body) {
      console.error("ReadableStream not supported in this environment");
      return;
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder("utf-8");
    const streamedContentDiv = document.getElementById(resultDiv);
    const messageDiv = document.getElementById(messageDivId);
    messageDiv.classList.remove("hidden");
    const pDiv = document.getElementById(pDivId);
    let contentpre = document.createElement("div");
    contentpre.classList.add(
      "text-sm",
      "p-2",
      "rounded-lg",
      "mt-2",
      "break-words",
      "whitespace-pre-wrap"
    );
    pDiv.appendChild(contentpre);
    document.getElementById("loader-" + pDivId).style.display = "none";

    let done = false;
    while (!done) {
      const { value, done: readerDone } = await reader.read();
      if (value) {
        let decodedValue = decoder.decode(value);
        console.log("Task log data ++++++++++++++++++0000", decodedValue);
        if (decodedValue.includes("data")) {
          decodedValue = decodedValue.replace(/data:/g, "");
        } else if (decodedValue.includes("DONE")) {
          break;
        }
        decodedValue = decodedValue.trim();
        decodedValue = decodedValue.replace(/}{/g, "},{");
        decodedValue = decodedValue.replace(/}\n{/g, "},{");
        try {
          let strval = "[" + decodedValue + "]";
          let objVal = JSON.parse(strval);
          for (let i = 0; i < objVal.length; i++) {
            objVal[i].choices.forEach((val) => {
              contentpre.innerHTML += decodeHTMLEntities(val.delta.content);
            });
          }
        } catch (e) {
          console.log(e);
        }
        streamedContentDiv.scrollTop = streamedContentDiv.scrollHeight; // Auto-scroll to the bottom
      }
      done = readerDone;
    }
    contentpre.innerHTML = renderMarkdown(contentpre.innerHTML);
    pDiv.appendChild(contentpre);
    pDiv.innerHTML = formatHtmlContent(pDiv.id);
  } catch (error) {
    console.error("Error fetching data:", error);
  } finally {
    document.getElementById("startChat").classList.remove("hidden");
    document.getElementById("inprogressChat").classList.add("hidden");
    CurrentAbortController = null;
  }
}

async function fetchServer(
  payload,
  apiUrl,
  myHeaders,
  resultDiv,
  messageDivId,
  pDivId
) {
  const raw = JSON.stringify(payload);
  try {
    const response = await fetch(apiUrl, {
      method: "POST",
      headers: myHeaders,
      body: raw,
      redirect: "follow",
    });
    const jsonResponse = await response.json();
    const streamedContentDiv = document.getElementById(resultDiv);

    const messageDiv = document.getElementById(messageDivId);
    messageDiv.classList.remove("hidden");
    const pDiv = document.getElementById(pDivId);
    let contentpre = document.createElement("div");
    contentpre.classList.add(
      "text-gray-100",
      "text-sm",
      "p-2",
      "rounded-lg",
      "mt-2",
      "break-words",
      "whitespace-pre-wrap"
    );
    document.getElementById("loader-" + pDivId).style.display = "none";
    // Assuming the response has the content in a similar structure
    if (jsonResponse.choices.length > 0) {
      jsonResponse.choices.forEach((val) => {
        if (val.messages.content !== "") {
          console.log(
            val.messages.content,
            "------------------------------------->>>>>>>>>>>>>>>>>>>>>>"
          );
          contentpre.innerHTML += renderMarkdown(
            decodeHTMLEntities(val.messages.content)
          );
        } else if (val.messages.toolCalls.length > 0) {
          val.messages.toolCalls.forEach((toolCall) => {
            const tool = document.createElement("pre");
            const toolCode = document.createElement("code");
            toolCode.className = "language-json";
            let pm = JSON.parse(toolCall.functionSchema.parameters);
            toolCode.textContent = JSON.stringify(pm, null, 2);
            tool.appendChild(toolCode);
            contentpre.appendChild(tool);
          });
        } else {
          contentpre.textContent = "No content found in response!";
        }
      });
    } else {
      contentpre.textContent = "No content found in response!";
    }

    pDiv.appendChild(contentpre);
    pDiv.innerHTML = formatHtmlContent(pDiv.id);
    // Append content directly to the div
    // Scroll to the bottom if needed
    streamedContentDiv.scrollTop = streamedContentDiv.scrollHeight;
  } catch (error) {
    console.error("Error fetching data:", error);
  } finally {
    document.getElementById("startChat").classList.remove("hidden");
    document.getElementById("inprogressChat").classList.add("hidden");
  }
}

export async function crawlUrlWithContent(urlDiv, scrap = false) {
  let url = document.getElementById(urlDiv).value;
  if (scrap) {
    try {
      // If scraping flag is enabled, adjust headers or options as needed
      const response = await fetch(url, {
        method: "GET",
        headers: {
          "Content-Type": "text/html",
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch URL: ${response.statusText}`);
      }

      // Get the content as text
      const htmlContent = await response.text();

      const tempDiv = document.createElement("div");
      tempDiv.innerHTML = htmlContent;

      // Extract plain text from the parsed HTML
      const plainText = tempDiv.textContent || tempDiv.innerText || "";

      document.getElementById("contextValue").value = plainText;
      // If the content needs further parsing (e.g., extracting JSON), handle here
      return url; // Returns the raw HTML content of the URL
    } catch (error) {
      console.error("Error while crawling the URL:", error.message);
      showAlert(
        AlertError,
        "Url - " + url,
        "Failed to crawl the URL. Please check the console for details."
      );
      return null;
    }
  } else {
    document.getElementById("contextValue").value = url;
    document.getElementById("contextType").value = "URL Reference";
    return url;
  }
}

// Placeholder functions for upload options
export function uploadFromComputer() {
  const fileInput = document.createElement("input");
  fileInput.type = "file";
  fileInput.accept = "*/*"; // Accept any file type
  fileInput.style.display = "none";

  fileInput.addEventListener("change", (event) => {
    const files = event.target.files;
    if (files.length > 0) {
      const file = files[0];
      const reader = new FileReader();
      // File content read as text
      reader.onload = (e) => {
        const fileContent = e.target.result;
        try {
          // Attempt to parse as JSON
          const jsonData = JSON.parse(fileContent);
          document.getElementById("contextValue").value = JSON.stringify(
            jsonData,
            null,
            2
          );
        } catch (error) {
          // If parsing fails, treat as a plain string
          document.getElementById("contextValue").value = fileContent;
        }
      };

      // Read the file content as text
      reader.readAsText(file);
    }
  });

  // Trigger the file input
  fileInput.click();
}

const AIStudioDatabase = {
  dbName: "AIStudioDB",
  storeName: "StudioPrompts",
  pKey: "userStudioAuthId",
};

let localDbConn = new LocalDb(
  AIStudioDatabase.dbName,
  AIStudioDatabase.storeName,
  AIStudioDatabase.pKey
);
window.addEventListener("DOMContentLoaded", () => {
  console.log("Page loaded. Initialilizing IndexDB...");
  localDbConn.openIndexDB();
  console.log("IndexDB initialized successfully!");
});

export function addContextLocally(obj, isPrompt) {
  try {
    const payload = {
      userStudioAuthId: obj.userId + "_" + obj.domain,
      value: [
        {
          id: generateUUIDv4(),
          userId: obj.userId,
          domain: obj.domain,
          timestamp: Date.now(),
          userPrompt: isPrompt,
          content: obj.content,
        },
      ],
    };

    // localDbConn.patchData(payload, payload.userStudioAuthId);
  } catch (error) {
    console.error("Error storing data locally:", error);
  } finally {
    console.log("Data stored locally!");
  }
}

export function loadAIStudioChat(canvasId) {
  try {
    const dataDiv = document.getElementById("currentChat");
    const chatObj = JSON.parse(dataDiv.innerText);
    const chatId = chatObj.chatId;
    const canvas = document.getElementById(canvasId);
    if (chatObj.messages.length > 0) {
      chatObj.messages.forEach((msg) => {
        let chId = generateUUIDv4();
        if (msg.role.toLowerCase() === "user") {
          canvas.innerHTML += `<div class="mb-2 flex justify-end">
              <div class="text-right">
                  <div class="bg-[#2d2d2d] text-sm text-gray-100 p-2 rounded-xl mt-2 break-words">
                  ${msg.content}
                  </div>
              </div>
            </div>`;
        } else if (msg.role.toLowerCase() === "assistant") {
          let content = renderMarkdown(decodeHTMLEntities(msg.content));
          const repDiv = document.createElement("div");
          repDiv.classList.add("mb-2", "flex", "justify-start");
          const repContDiv = document.createElement("div");
          repContDiv.classList.add(
            "text-left",
            "rounded-xl",
            "bg-zinc-900 text-gray-100"
          );
          repContDiv.id = "r-" + chId;
          repContDiv.innerHTML = content;
          repDiv.appendChild(repContDiv);
          canvas.appendChild(repDiv);
          repContDiv.innerHTML = formatHtmlContent("r-" + chId);
        }
      });
    }
    canvas.scrollTop = canvas.scrollHeight;
  } catch (error) {
    console.error("Error loading data locally:", error);
  }
}

function cancelStream() {
  if (CurrentAbortController) {
    console.log("Aborting current fetch stream...");
    CurrentAbortController.abort();
    showInfoMessage(
      "Message Cancelled",
      "The current chat message has been cancelled."
    );
  }
}

async function selectPrompt(promptId, tokenKey, apiUrl) {
  promptId = promptId.trim();
  if (promptId === "") {
    showErrorMessage("Prompt Selection Error", "Please select a valid prompt.");
    return;
  }
  apiUrl = apiUrl + "/" + promptId;
  const response = await fetch(
    apiUrl,
    getRequestOptions(tokenKey, "GET", null)
  );

  // Process the response
  if (!response.ok) {
    showErrorMessage(
      "Prompt Selection Error",
      "Failed to select the prompt. Please try again or continue without selecting a prompt."
    );
    return;
  }

  const result = await response.json();
  console.log(result);
  console.log("====================");
  var variablesMap = {};
  if (result === null) {
    showErrorMessage(
      "Prompt Selection Error",
      "Failed to select the prompt. Please try again or continue without selecting a prompt."
    );
    return;
  }
  if (result.output.length === 0) {
    showErrorMessage(
      "Prompt Selection Error",
      "Failed to select the prompt. Please try again or continue without selecting a prompt."
    );
    return;
  }
  result.output.forEach((output) => {
    if (output.spec.variables.length > 0) {
      output.spec.variables.forEach((variable) => {
        variable = variable.trim();
        variable = variable.replace("{{", "");
        variable = variable.replace("}}", "");
        variablesMap[variable] = "";
      });
      document.getElementById("promptInput").value = JSON.stringify(
        variablesMap,
        null,
        2
      );
      document.getElementById("promptInput").classList.remove("hidden");
      document.getElementById("input").classList.add("hidden");
    }
  });
}

function retryInput(content, inputElementId) {
  document.getElementById(inputElementId).value = content;
  var event = new KeyboardEvent("keydown", {
    key: "Enter",
    code: "Enter",
    keyCode: 13,
    which: 13,
    bubbles: true,
  });

  input.dispatchEvent(event);
}

window.addContextLocally = addContextLocally;
window.crawlUrlWithContent = crawlUrlWithContent;
window.uploadFromComputer = uploadFromComputer;
window.aiInterfaceAction = aiInterfaceAction;
window.loadAIStudioChat = loadAIStudioChat;
window.cancelStream = cancelStream;
window.selectPrompt = selectPrompt;
window.retryInput = retryInput;
// window.addEventListener("DOMContentLoaded", () => {
//   promptEditor = new JSONEditor(document.getElementById("promptInput"), {
//     mode: "code",
//     mainMenuBar: false,
//     onError: function (err) {
//       console.error(err);
//     },
//   });
// });

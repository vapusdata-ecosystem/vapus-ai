import {
  addchatCodeBox, addchatDataBox, addchatDataErrBox, addchatDataSuccessBox, addChatDateTimeBox, addChatReasoningDataBox,
  addChatStateBox,
  addUserInputBox, buildAndDownloadFile
} from './chat-interface.js';
import { dataFileUploader } from './fileloader.js';
import { LocalDb } from './indexdb.js';

export let CurrentAbortController = null;

let taskLogSSEHandler = false;

var tablesRendered = false;

const DatafabricDatabase = {
  dbName: "DataNabhikDB",
  storeName: "VapusQueries",
  pKey: "queryId",
  batchSize: 200,
  waitTime: 50
};

const NeedAgentRun = false;
export let currentChatId = '';
export let currChatObj;
export let selectedMessageId = '';
export let attachDatasets = [];
let localDbConn = new LocalDb(DatafabricDatabase.dbName, DatafabricDatabase.storeName, DatafabricDatabase.pKey);
window.addEventListener("DOMContentLoaded", () => {
  localDbConn.openIndexDB();
});

export function toggleInputType() {
  const inputType = document.getElementById('inputFormat').value;
  const textInput = document.getElementById('textInput');

  if (inputType === 'query') {
    textInput.classList.remove('hidden');
  } else if (inputType === 'api') {
    textInput.classList.add('hidden');
  }
}

async function uploadDataset(api, tokenKey, resource, fileCanvasId, chatId) {
  if (chatId !== "") {
    currentChatId = chatId;
    console.log("currentChatId ---- ", currentChatId);
    console.log("attachDatasets ---- ", attachDatasets);
    await dataFileUploader(api, tokenKey, resource, chatId, fileCanvasId,attachDatasets);
  } else {
    if (currChatObj) {
      if (currChatObj !== null && currChatObj !== undefined) {
        currentChatId = currChatObj.chatId;
        await dataFileUploader(api, tokenKey, resource, currentChatId, fileCanvasId,attachDatasets);
      } else {
        showAlert(AlertError, "Error", "Please start a chat to upload a dataset");
      }
    } else {
      showAlert(AlertError, "Error", "Please start a chat to upload a dataset");
    }
  }

}

async function createNewChat(manageUrl, tokenKey) {
  const myHeaders = new Headers();
  myHeaders.append("Accept", "application/x-ndjson");
  myHeaders.append("Content-Type", "application/x-ndjson");
  const apiToken = getCookie(tokenKey);
  myHeaders.append(
    "Authorization",
    `Bearer ${apiToken}`);
  const payload = JSON.stringify({
    action: "CREATE"
  });
  try {
    const response = await fetch(
      manageUrl,
      {
        method: "POST",
        headers: myHeaders,
        body: payload,
        redirect: "follow",
      }
    );
    if (!response.ok) {
      showErrorMessage("Error", "Error while creating new chat, please try again");
      return;
    } else {
      const result = await response.json();
      if (result.output !== undefined && result.output.length === 1) {
        return result.output[0];
      } else {
        showErrorMessage("Error", "Error while creating new chat, please try again");
      }
    }
  } catch (error) {
    console.error("Error creating new chat:", error);
    showErrorMessage("Error", "Error while creating new chat, please try again");
    return null;
  }
}
export async function dataNabhikAction(tokenKey, apiUrl, manageUrl, username,downloadFileUrl, chId,taskLogStreamUrl) {
  const input = document.getElementById('input').value;
  document.getElementById('input').disabled = true;
  const dpId = document.getElementById('dataProduct').value;
  if (input === "") {
    showErrorMessage("Missing Input", "Please enter a valid input to proceed");
    return;
  }
  const genSug = document.getElementById('nabhik-suggetion-generic');
  genSug.classList.add("hidden");
  if (chId !== "") {
    currentChatId = chId;
  } else {
    if (currChatObj !== null && currChatObj !== undefined) {
      currentChatId = currChatObj.chatId;
    } else {
      currentChatId = '';
      currChatObj = await createNewChat(manageUrl, tokenKey);
      if (currChatObj === null || currChatObj === undefined) {
        showErrorMessage("Error", "Error while creating new chat, please try again");
        return;
      } else {
        currentChatId = currChatObj.chatId;
        updateBrowserUrl("chatId", currentChatId);

      }
    }
  }
  
  document.getElementById("ask").classList.add("hidden");
  document.getElementById("inprogressChat").classList.remove("hidden");
  var queryParams = {
    chatId: currentChatId,
    input: input,
  };
  if (dpId === "") {
    queryParams.dataproducts = [];
  } else {
    queryParams.dataproducts = [dpId];
  }
  if (selectedMessageId === "") {
    queryParams.messageId = [];
  } else {
    queryParams.messageId = [selectedMessageId];
  }
  console.log("attachDatasets ---- ", attachDatasets);
  if (attachDatasets.length > 0) {
    queryParams.fileData = [];
    attachDatasets.forEach(dataset => {
      queryParams.fileData.push({
        name: dataset
      });
    });
  }
  selectedMessageId = '';
  const myHeaders = new Headers();
  myHeaders.append("Accept", "application/x-ndjson");
  myHeaders.append("Content-Type", "application/x-ndjson");
  const apiToken = getCookie(tokenKey);
  myHeaders.append(
    "Authorization",
    `Bearer ${apiToken}`);
  const payload = JSON.stringify(queryParams);
  // requestOptions.mode = "no-cors";
  // if (streamChat === 'true') {
  return submitQueryForStream(apiUrl, myHeaders, payload, username,downloadFileUrl,tokenKey,taskLogStreamUrl);
}
async function submitQueryForStream(url, myHeaders, payload, username,downloadFileUrl,tokenKey,taskLogStreamUrl) {
  const startTime = getEpochTime();
  try {
    CurrentAbortController = new AbortController();
    const signal = CurrentAbortController.signal;
    let fabricCanvas = document.getElementById("nabhikCanvasBoard");
    fabricCanvas.classList.remove("hidden");
    const messageCanvasElem = document.createElement("div");
    messageCanvasElem.classList.add("my-2","overflow-y-auto","scrollbar");
    const tempMessid = generateUUID();
    const suggestionDiv = document.getElementById('suggestionDiv');
    messageCanvasElem.id = tempMessid;
    fabricCanvas.appendChild(messageCanvasElem);
    const messageCanvas = document.getElementById(tempMessid);
    messageCanvas.classList.add("relative", "p-4","text-sm","rounded-lg");
    const reasoningCanvas = document.createElement("details");
    reasoningCanvas.classList.add("p-2", "rounded-lg", "my-2", "text-xs", "text-gray-100", "chat-loader", "font-semibold","cursor-pointer","border","border-zinc-500");
    reasoningCanvas.id = "reasoning" + tempMessid;
    reasoningCanvas.innerHTML = `<summary class="text-xs" id="reasoning-summary-${tempMessid}">Thinking...</summary>`;
    reasoningCanvas.open = false;

    // messageCanvas.classList.add("border", "border-gray-400", "[clip-path:polygon(15px_0,calc(100%_-_15px)_0,100%_15px,100%_calc(100%_-_15px),calc(100%_-_15px)_100%,15px_100%,0_calc(100%_-_15px),0_15px)]");
    addUserInputBox("input",document.getElementById('input').value, messageCanvas);
    const nabhikCanvas = document.createElement("div");
    nabhikCanvas.classList.add("relative","border","hover:border-zinc-500","bg-[#1b1b1b]", "border-[#1b1b1b]", "p-4","rounded-lg");
    messageCanvas.appendChild(nabhikCanvas);
    nabhikCanvas.appendChild(reasoningCanvas);
    messageCanvas.scrollTop = messageCanvas.scrollHeight;
    document.getElementById('input').value = "";
    var tableMap = new Map();
    var tableHeaderMap = new Map();
    const response = await fetch(
      url,
      {
        method: "POST",
        headers: myHeaders,
        body: payload,
        redirect: "follow",
        signal,
      }
    );

    if (!response.body) {
      console.error("ReadableStream not supported in this environment");
      return;
    }
    const reader = response.body.getReader();
    const decoder = new TextDecoder("utf-8");
    let done = false;

    let dataFields = [];
    let messCahe;
    let counter = 0;
    let responseId = "";
    let contentTables = [];
    let tableBodyMap = new Map();

    while (!done) {
      messageCanvas.scrollTop = messageCanvas.scrollHeight;
      fabricCanvas.scrollTop = fabricCanvas.scrollHeight;
      nabhikCanvas.scrollTop = nabhikCanvas.scrollHeight;
      const { value, done: readerDone } = await reader.read();
      if (value) {
        let decodedValue = decoder.decode(value);
        decodedValue = decodedValue.trim();
        try {
          decodedValue = decodedValue.replace(/}{/g, "},{");
          decodedValue = decodedValue.replace(/}\n{/g, "},{");
          let strval;
          let objVal;
          try {
            strval = "[" + decodedValue + "]";
            objVal = JSON.parse(strval);
            if (objVal.length < 0) {
              continue;
            }
            messCahe = null;
          } catch (error) {
            if (messCahe !== null) {
              messCahe = messCahe + decodedValue;
            } else {
              messCahe = decodedValue;
            }
            try {
              strval = "[" + messCahe + "]";
              objVal = JSON.parse(strval);
            } catch (error) {
              continue;
            }
            console.error("Error parsing JSON: Will concartenate with prvious messages", error);
          }

          objVal.forEach(async responseJson => {
            if (responseJson !== null && responseJson.message !== undefined) {
              addchatDataErrBox(responseJson.message, nabhikCanvas);
              return;
            }
            if (responseJson === null || responseJson.result === null || responseJson.result.output === null || responseJson.result.output === undefined) {
              addchatDataErrBox("Error while querying data product either there is no data for this query or there is some internal server error, please try again or contact the data product owner",
                nabhikCanvas);
              return;
            }
            switch (responseJson.result.output.event) {
              case "CHAT_OVERFLOWEN":
                addchatDataBox(null,responseJson.result.output.data.content, nabhikCanvas, true);
                endChatOptions(null,"EXPIRED");
                break;
              case "TASK_CREATED":
                console.log("Task Created", responseJson.result);
                const taskCanvas = document.createElement("div");
                taskCanvas.classList.add("relative","border","hover:border-zinc-500","bg-[#1b1b1b]", "border-[#1b1b1b]", "p-2","rounded-lg","m-2","max-h-[600px]","overflow-y-auto","scrollbar");
                taskCanvas.id = responseJson.result.taskId;
                
                const taskLink = document.createElement("a");
                taskLink.classList.add("text-gray-100", "text-sm", "rounded-lg", "px-2", "py-1", "border-[2px]", "border-zinc-500", "cursor-pointer", "hover:bg-orange-700");
                taskLink.href = "/ui/ai/manage/nabhiktasks/" + responseJson.result.taskId;
                taskLink.innerHTML = "Click here to view task details";
                taskLink.target = "_blank";
                nabhikCanvas.appendChild(taskLink);
                nabhikCanvas.appendChild(taskCanvas);
                loadTaskCurrentLogs(responseJson.result.taskId,"", tokenKey, taskLogStreamUrl,downloadFileUrl, document.getElementById("nabhikCanvasBoard"),"");
                break;
              case "END":
                if (responseJson.result.output.data !== undefined && responseJson.result.output.data.final !== undefined) {
                  if (responseJson.result.output.data.final.reason === "SUCCESSFULL") {
                    addchatDataSuccessBox(responseJson.result.output.data.content, nabhikCanvas);
                  } else if (responseJson.result.output.data.final.reason === "DONE") {
                    break;
                  } else {
                    addchatDataErrBox(JSON.stringify({
                      key: responseJson.result.output.data.final.reason,
                      value: responseJson.result.output.data.final.metadata
                    }), nabhikCanvas);
                  }
                  break;
                } 
              case "ABORTED":
                addchatDataErrBox(JSON.stringify({
                  key: responseJson.result.output.data.final.reason,
                  value: responseJson.result.output.data.final.metadata
                }), nabhikCanvas);
                break;
              case "START":
                setChatId(responseJson.result.chatId);
                break;
              case "DATETIME":
                if (responseJson.result.output.data.content !== "") {
                  addChatDateTimeBox("", responseJson.result.output.data.content, nabhikCanvas);
                }
                break;
              case "STATE":
                addChatStateBox(responseJson.result.output.data.content,responseJson.result.output.data.contentType, reasoningCanvas);
                break;
              case "REASONINGS":
                if (responseJson.result.output.data.content !== undefined && responseJson.result.output.data.content.length > 0) {
                  addChatReasoningDataBox(responseJson.result.output.data.content, nabhikCanvas);
                }
              case "FILE_DATA":
                if (responseJson.result.output.files) {
                  if (responseJson.result.output.files.name.length > 0) {
                    // await renderNabhikFiles(responseJson.result.output.files, messageCanvas, responseId);
                    await handleNabhikFile(nabhikCanvas, responseJson.result.output.files, tokenKey, downloadFileUrl,true);
                  }
                }
                break;
              case "RESPONSE_ID":
                responseId = responseJson.result.responseId;
                nabhikCanvas.id = responseId;
                break;
              case "SUGGESTIONS":
                addstreamSuggestionBox(responseJson.result.output.data.content, suggestionDiv,submitInput);
                break;
              default:
                if (responseJson.result.output.data.contentType === "JSON") {
                  let content = JSON.parse(responseJson.result.output.data.content);
                  addchatDataBox(content.key, content.value, nabhikCanvas, false);
                // } else if (responseJson.result.output.data.contentType === "MARKUP") {
                //   addstreamMarkdownDataBox(responseJson.result.output.data.content,responseJson.result.output.data.contentId, messageCanvas);
                } else if (responseJson.result.output.data.contentType === "CLICKSET") {
                  addstreamClickableBox(responseJson.result.output.data.content, nabhikCanvas,submitInput);
                } else if (responseJson.result.output.data.contentType === "CODE") {
                  addchatCodeBox(responseJson.result.output.data.content, nabhikCanvas);
                } else {
                  addchatDataBox(null,responseJson.result.output.data.content, nabhikCanvas, true);
                }
                break;
            }
          });

        } catch (error) {
          console.error(error);
          addstreamDataErrBox("Error while processing your request, please try again", nabhikCanvas);
        } finally {
          console.log("Done reading data");
        }
      }
      done = readerDone;
    }
    try {
      reasoningCanvas.classList.remove("chat-loader");
      let timeDiff = getEpochTime() - startTime;
      addChatStateBox("Done: Time - " + timeDiff + "sec", "PLAIN_TEXT", reasoningCanvas);
      reasoningCanvas.open = false;
      // hideStreamDiv(responseId);
      // hideStreamDiv(currentChatId);
      const divSperator = document.createElement("div");
      divSperator.classList.add("flex", "mt-2", "cursor-pointer", "text-xs", "text-gray-400", "hover:text-gray-100");
      divSperator.innerHTML = "Select Message";
      const innerDiv = document.createElement("div");
      innerDiv.innerHTML = copySvg;
      divSperator.addEventListener("click", function () {
        selectedMessageId = nabhikCanvas.id;
        const toast = document.getElementById('toast');
        toast.textContent = "message selected";
        toast.classList.add('show');
      });
      const seperator = document.createElement("hr");
      seperator.classList.add("border-gray-400", "border-2", "my-2");
      divSperator.appendChild(innerDiv);
      messageCanvas.appendChild(divSperator);
      messageCanvas.appendChild(seperator);
      fabricCanvas.appendChild(messageCanvas);
      try {
        document.getElementById("reasoning-summary-" + tempMessid).innerHTML = "Thoughts";
      } catch {
        console.error(error);
      }
    } catch (error) {
      console.error(error);
    }
  } catch (error) {
    console.error("Error:", error);
  } finally {
    document.getElementById('input').disabled = false;
    document.getElementById("ask").classList.remove("hidden");
    document.getElementById("inprogressChat").classList.add("hidden");
    CurrentAbortController = null;
  }
}

function setChatId(op) {
  if (op !== "") {
    currentChatId = op;
    return;
  }
}

function submitInput(content) {
  document.getElementById('input').value = content;
  document.getElementById('ask').click();
}

function handleDataset(tableMap,dataset, dataTableBody, dataFields, responseId) {
  let tablesRendered = false;
  let result = parseAndreturnJSONArrayElem(dataset);
  let datasetList = tableMap.get(responseId);
  if (datasetList === undefined) {
    datasetList = result[0];
  } else {
    datasetList = datasetList.concat(result[0]);
  }
  tableMap.set(responseId, datasetList);
  // setTimeout(() => {
  //   addTableRow(dataTableBody, dataFields, result[0],result[1]);
  //   // addDataSetInIndexDb(responseId, dbBatch);
  //   tablesRendered = true;
  // }, 10);
  return 0;
}

function addDataSetInIndexDb(iden, dataSet) {
  try {
    const payload = {
      queryId: iden,
      value: dataSet
    }
    localDbConn.patchData(payload, iden);
  } catch (error) {
    console.error("Error storing dataset locally:", error);
  }
}

export async function exportResultSet(responseId) {
  const dataset = await localDbConn.retrieveData(responseId);
  const format = document.getElementById('exportFormat').value;
  // let parsedData = JSON.parse(dataset);
  // console.log(parsedData);
  if (format === 'JSON') {
    dataToJSON(dataset.value, '');
  } else if (format === 'CSV') {
    dataToCSV(dataset.value, '');
  } else {
    console.error('Invalid export format');
  }
}

export function showPrompts() {
  try {
    const dp = document.getElementById("dataProduct").value;
    if (dp === "") {
      document.getElementById("dp-prompts").classList.add("hidden");
      showErrorMessage("Invalid param", "Please select a data product");
      return;
    } else {
      document.getElementById("prompts-" + dp).classList.remove("hidden");
    }
  } catch (error) {
    showErrorMessage("Invalid param", "Please select a data product");
  }
}

export function closePrompts(id) {
  document.getElementById(id).classList.add("hidden");
}

function retryInput(content,inputElementId) {
  document.getElementById(inputElementId).value = content;
  var event = new KeyboardEvent('keydown', {
    key: 'Enter',
    code: 'Enter',
    keyCode: 13,
    which: 13,
    bubbles: true,
  });

  input.dispatchEvent(event);
}

function selectPromptIntoInput(prompt, inputDiv, promptId) {
  try {
    const input = document.getElementById(inputDiv);
    input.value = prompt;
    var event = new KeyboardEvent('keydown', {
      key: 'Enter',
      code: 'Enter',
      keyCode: 13,
      which: 13,
      bubbles: true,
    });

    input.dispatchEvent(event);
    closePrompts(promptId)

  } catch (error) {
    console.error("Error selecting prompt into input:", error);
  }
}

async function loadNabhikChat(username, accessTokenKey, downloadUrl, canvasId, currChatObject,taskLogStreamUrl) {
  try {
    if (currChatObject !== null) {
      const canvas = document.getElementById(canvasId);
      canvas.classList.remove("hidden");
      if (currChatObject.logs !== undefined && currChatObject.logs.length > 0) {
        currChatObject.logs.forEach(log => {
          const logCanvas = document.createElement("div");
          logCanvas.id = generateUUID();
          logCanvas.classList.add("my-2","w-full","rounded-lg","p-2");
          addUserInputBox("input",log.input, logCanvas);
          if (log !== undefined) {
              const messageCanvas = document.createElement("div");
              messageCanvas.classList.add("relative","border","hover:border-zinc-500","bg-[#1b1b1b]", "border-[#1b1b1b]", "p-4","rounded-lg");
              messageCanvas.id = log.messageId;
              var timelog = document.createElement("div");
              timelog.classList.add("absolute", "-top-3", "-left-15", "text-xs", "font-semibold", "text-gray-100","rounded-lg","bg-[#1b1b1b]","hover:display","hidden");

              if (log.startTime) {
                addChatDateTimeBox("", log.startTime, timelog);
                messageCanvas.appendChild(timelog);
              }
              if (log.instructions) {
                addchatDataBox(null, log.instructions, messageCanvas, false);
              }
              if (log.reasonings !== undefined && log.reasonings.length > 0) {
                addChatReasoningDataBox(log.reasonings, messageCanvas);
              }
              if (log.output !== undefined && log.output !== "") {
                addchatDataBox(null, log.output, messageCanvas, false);
                if (log.messageType.toLowerCase() === "task") {
                    const taskLink = document.createElement("a");
                    taskLink.classList.add("text-gray-100","text-sm","rounded-lg","px-2","py-1","border-[2px]","border-zinc-500","cursor-pointer","hover:bg-orange-700");
                    taskLink.href = "/ui/ai/manage/nabhiktasks/"+log.taskId;
                    taskLink.innerHTML = "Click here to view task details";
                    taskLink.target = "_blank";
                    messageCanvas.appendChild(taskLink);
                }
              }
              if (log.outputAssets !== undefined && log.outputAssets.length > 0) {
                log.outputAssets.forEach(async asset => {
                  if (asset !== undefined && asset.file !== undefined && asset.file.name !== "") {
                    await handleNabhikFile(messageCanvas, asset.file, accessTokenKey, downloadUrl, true);
                  }
                });
              }
              if (log.messageType.toLowerCase() === "task") {
                const taskCanvas = document.createElement("div");
                taskCanvas.classList.add("relative","border","hover:border-zinc-500","bg-[#1b1b1b]", "border-[#1b1b1b]", "p-2","rounded-lg","m-2","max-h-[600px]","overflow-y-auto","scrollbar");
                taskCanvas.id = log.taskId;
                messageCanvas.appendChild(taskCanvas);
              }
              const divSperator = document.createElement("div");
              divSperator.classList.add("flex", "mt-2", "cursor-pointer", "text-xs", "text-gray-400", "hover:text-gray-100");
              divSperator.innerHTML = "Select Message";
              const innerDiv = document.createElement("div");
              innerDiv.innerHTML = copySvg;
              divSperator.addEventListener("click", function () {
                selectedMessageId = log.messageId;
                const toast = document.getElementById('toast');
                toast.textContent = "message selected";
                toast.classList.add('show');
              });
              divSperator.appendChild(innerDiv);
              messageCanvas.appendChild(divSperator);
              logCanvas.appendChild(messageCanvas);
          }
          canvas.appendChild(logCanvas);
        });
        endChatOptions(currChatObject.base);
      }
      canvas.scrollTop = canvas.scrollHeight;
      await loadTaskCurrentLogs("",currChatObject.chatId, accessTokenKey, taskLogStreamUrl,downloadUrl,canvas,"COMPLETED");
    }
  } catch (error) {
    console.error("Error loading fabric chat:", error);
  }
}

async function endChatOptions(base,status){
  if (base !== null && base !== undefined) {
    if(base.status){
      status = base.status;
    }
  }
  if (status === "END" || status === "EXPIRED") {
    document.getElementById("userInputGround").classList.add("hidden");
    document.getElementById("endedChatMessage").classList.remove("hidden");
  }
}

async function loadTaskCurrentLogs(taskId,nabhikChatId, accessTokenKey, taskLogStreamUrl,downloadFileUrl, mainCanvas,status) {
  try {
    const myHeaders = new Headers();
    myHeaders.append("Accept", "application/x-ndjson");
    myHeaders.append("Content-Type", "application/x-ndjson");
    const apiToken = getCookie(accessTokenKey);
    myHeaders.append(
      "Authorization",
      `Bearer ${apiToken}`);
    let payload = {
      nabhikChatId: nabhikChatId,
      status:status,
      taskIds: [],
    };
    if (taskId !== "") {
      payload.taskIds.push(taskId);
    }
    const raw = JSON.stringify(payload);
    const response = await fetch(
      taskLogStreamUrl,
      {
        method: "POST",
        headers: myHeaders,
        body: raw,
        redirect: "follow",
        // signal,
      }
    );

    if (!response.body) {
      console.error("ReadableStream not supported in this environment");
      return;
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder("utf-8");
    let done = false;
    while (!done) {
      const { value, done: readerDone } = await reader.read();
      if (value) {
        let decodedValue = decoder.decode(value);
        if (decodedValue.includes("data")) {
          decodedValue = decodedValue.replace(/data:/g, "");
        } else {
          break;
        }
        decodedValue = decodedValue.trim();
        decodedValue = decodedValue.replace(/}\n\n {/g, "},{");
        if (decodedValue.includes("[DONE]")) {
          decodedValue = decodedValue.replace(/\[DONE\]/g, "");
          decodedValue = decodedValue.trim();
        }
        try {
          console.log("Decoded Value--------------------------------------------", decodedValue);
          let strval = "[" + decodedValue + "]";
          let logChunks = JSON.parse(strval);
          for (let i = 0; i < logChunks.length; i++) {
            const logChunk = logChunks[i];
            const taskCanvas = document.getElementById(logChunk.taskId);
            switch (logChunk.event) {
              case "END":
                if (logChunk.data.final.reason === "SUCCESSFULL") {
                  addchatDataSuccessBox(logChunk.data.content, taskCanvas);
                } else if (logChunk.data.final.reason === "DONE" || logChunk.data.final.reason === "END" || logChunk.data.final.reason === "ABORTED") {
                  break;
                } else {
                  addchatDataErrBox(JSON.stringify({
                    key: logChunk.data.final.reason,
                    value: logChunk.data.final.metadata
                  }), taskCanvas);
                }
                break;
              case "ABORTED":
                addchatDataErrBox(JSON.stringify({
                  key: logChunk.data.final.reason,
                  value: logChunk.data.final.metadata
                }), taskCanvas);
                break;
              case "DATETIME":
                if (logChunk.data.content !== "") {
                  addChatDateTimeBox("", logChunk.data.content, taskCanvas);
                }
                break;
              case "STATE":
                addChatStateBox(logChunk.data.content, logChunk.data.contentType, reasoningCanvas);
                break;

              case "REASONINGS":
                if (logChunk.data.content !== undefined && logChunk.data.content.length > 0) {
                  addChatReasoningDataBox(logChunk.data.content, taskCanvas);
                }
              case "FILE_DATA":
                if (logChunk.files) {
                  if (logChunk.files.name.length > 0) {
                    // await renderNabhikFiles(responseJson.result.output.files, messageCanvas, responseId);
                    await handleNabhikFile(taskCanvas, logChunk.files, accessTokenKey, downloadFileUrl, true);
                  }
                }
                break;
              case "SUGGESTIONS":
                addstreamSuggestionBox(logChunk.data.content, suggestionDiv, submitInput);
                break;
              default:
                if (logChunk.data.contentType) {
                  addchatDataBox(null, logChunk.data.content, taskCanvas, true);
                } else {
                  if (logChunk.data.contentType === "JSON") {
                    let content = JSON.parse(logChunk.data.content);
                    addchatDataBox(content.key, content.value, taskCanvas, false);
                  } else if (logChunk.data.contentType === "CODE") {
                    addchatCodeBox(logChunk.data.content, taskCanvas);
                  } else {
                    addchatDataBox(null, logChunk.data.content, taskCanvas, true);
                  }
                }
                break;
            }
          }
        } catch (e) {
          console.log(e);
        }
        mainCanvas.scrollTop = mainCanvas.scrollHeight;
      }
      done = readerDone;
    }
  } catch (error) {
    console.error("Error loading task logs:", error);
  } finally {
    taskLogSSEHandler = false;
  }
}

async function handleNabhikFile(canvas, file, accessTokenKey, downloadUrl,showContent=false) {
  const mainBox = document.createElement("div");
  mainBox.classList.add("m-2");
  const nBox = document.createElement("div");
  nBox.classList.add("flex","items-center", "text-xs","cursor-pointer", "text-gray-100");
  var svgFiles = document.createElement("span");
  svgFiles.classList.add("w-6", "h-6");
  svgFiles.innerHTML = `
  <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
    <path fill-rule="evenodd" d="M6 2a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6H6zm7 1.5V8h4.5L13 3.5zm-3 8.75a2.25 2.25 0 1 1 0 4.5 2.25 2.25 0 0 1 0-4.5zm-3 5.25a1.75 1.75 0 1 1 0-3.5 1.75 1.75 0 0 1 0 3.5zm8-5.25a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3z" clip-rule="evenodd"/>
  </svg>
  `;
  nBox.appendChild(svgFiles);
  var assetContext = document.createElement("span");
  downloadNabhikFile(file.path, renderMDContent(file.description),accessTokenKey, canvas,downloadUrl,showContent, false);
  assetContext.classList.add("ml-2");
  assetContext.textContent = file.description;
  nBox.appendChild(assetContext);
  const downloadLink = document.createElement("a");
  downloadLink.classList.add("text-gray-100","text-sm","rounded-lg","px-2","py-1","border-[2px]","border-zinc-500","cursor-pointer"); 
  downloadLink.href = "#";
  // downloadLink.textContent = text;
  downloadLink.onclick = () => downloadNabhikFile(file.path,renderMDContent(file.description), accessTokenKey, canvas,downloadUrl,false, true);
  downloadLink.appendChild(svgFiles);
  nBox.appendChild(downloadLink);
  mainBox.appendChild(nBox);
  canvas.appendChild(mainBox);
  canvas.scrollTop = canvas.scrollHeight;
}

async function downloadNabhikFile(fileName, fileDescription,accessTokenKey, canvas,downloadUrl, showContent=false, download=false) {
  const myHeaders = new Headers();
  myHeaders.append("Accept", "application/x-ndjson");
  myHeaders.append("Content-Type", "application/json");
  const apiToken = getCookie(accessTokenKey);
  myHeaders.append(
    "Authorization",
    `Bearer ${apiToken}`);
  const payload = JSON.stringify({ fileName: fileName });
  downloadUrl = downloadUrl + "?fileName=" + fileName;
  try {
    const response = await fetch(downloadUrl, {
      method: "GET",
      headers: myHeaders,
      redirect: "follow",
    });
    if (!response.ok) {
      showErrorMessage("Error", "Error downloading file, please try again");
    }

    const result = await response.json();
    result.output.forEach(file => {
      if (file.name !== undefined && file.name !== "") {
          buildAndDownloadFile(file,fileDescription,canvas,showContent,download);
      }

    });
  } catch (error) {
    console.error("Error downloading file:", error);
    showErrorMessage("Error", "Error downloading file, please try again");
  }
}

function clearAttachDataset(id) {
  attachDatasets = [];
}


function cancelStream() {
  if (CurrentAbortController) {
    CurrentAbortController.abort();
    showInfoMessage("Message Cancelled", "The current chat message has been cancelled.");
  }
}

window.toggleInputType = toggleInputType;
window.exportResultSet = exportResultSet;
window.showPrompts = showPrompts;
window.closePrompts = closePrompts;
window.loadNabhikChat = loadNabhikChat;
window.currChatObj = currChatObj;
window.uploadDataset = uploadDataset;
window.cancelStream = cancelStream;
window.retryInput = retryInput;
window.dataNabhikAction = dataNabhikAction;
window.clearAttachDataset = clearAttachDataset;
window.selectPromptIntoInput = selectPromptIntoInput;
window.loadTaskCurrentLogs = loadTaskCurrentLogs;

import {
  addTable,
  addTableBody,
  addUserInputBox, addchatCodeBox, addchatDataBox,
  addchatDataErrBox, addchatDataSuccessBox
} from './chat-interface.js';
import { dataFileUploader } from './fileloader.js';
import { LocalDb } from './indexdb.js';
let CurrentAbortController = null;
const DataDatabaseDatabase = {
  dbName: "VapusData",
  storeName: "dataServerQueries",
  pKey: "queryId",
  batchSize: 200,
  waitTime: 50
};
let localDbConn = new LocalDb(DataDatabaseDatabase.dbName, DataDatabaseDatabase.storeName, DataDatabaseDatabase.pKey);
window.addEventListener("DOMContentLoaded", () => {
  localDbConn.openIndexDB();
});
const NeedAgentRun = false;
let attachDatasets = [];
function toggleInputType() {
  const inputType = document.getElementById('inputFormat').value;
  const textInput = document.getElementById('textInput');

  if (inputType === 'query') {
    textInput.classList.remove('hidden');
  } else if (inputType === 'api') {
    textInput.classList.add('hidden');
  }
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

async function uploadDataset(api, tokenKey, resource, fileCanvasId) {
  let resourceId = generateUUID();
  console.log("Resource ID: ", resourceId, "Resource: ", resource, "File Canvas ID: ", fileCanvasId,"API: ", api,"TokenKey: ", tokenKey);
  await dataFileUploader(api, tokenKey, resource, resourceId, fileCanvasId, attachDatasets);
}

async function queryDataserverAction(tokenKey, apiUrl, username,downloadFileUrl) {
  const input = document.getElementById('input').value;
  document.getElementById('input').disabled = true;
  const dpId = document.getElementById('dataProduct').value;
  if (input === "") {
    showErrorMessage("Missing Input", "Please enter a valid input to proceed");
    return;
  }

  document.getElementById("ask").classList.add("hidden");
  document.getElementById("inprogressChat").classList.remove("hidden");
  var queryParams = {
    textQuery: input,
  };
  if (dpId !== "") {
    queryParams.dataproducts = [dpId];
  }
  
  if (attachDatasets.length > 0) {
    queryParams.fileData = [];
    attachDatasets.forEach(dataset => {
      queryParams.fileData.push({
        name: dataset
      });
    });
  }
  const myHeaders = new Headers();
  myHeaders.append("Accept", "application/x-ndjson");
  myHeaders.append("Content-Type", "application/json");
  const apiToken = getCookie(tokenKey);
  myHeaders.append(
    "Authorization",
    `Bearer ${apiToken}`);
  const payload = JSON.stringify(queryParams);
  return submitQueryForStream(apiUrl, myHeaders, payload, username,downloadFileUrl,tokenKey);
}
async function submitQueryForStream(url, myHeaders, payload, downloadFileUrl,tokenKey) {
  const startTime = getEpochTime();
  try {
    CurrentAbortController = new AbortController();
  const signal = CurrentAbortController.signal;
    let fabricCanvas = document.getElementById("dataServerOutput");
    fabricCanvas.classList.remove("hidden");
    const messageCanvasElem = document.createElement("div");
    messageCanvasElem.classList.add("my-2","overflow-y-auto","scrollbar","bg-[#1b1b1b]", "text-gray-100");
    const tempMessid = generateUUID();
    messageCanvasElem.id = tempMessid;
    fabricCanvas.appendChild(messageCanvasElem);
    const messageCanvas = document.getElementById(tempMessid);
    messageCanvas.classList.add("relative", "border", "border-zinc-500", "p-4","text-sm","rounded-lg");

    // messageCanvas.classList.add("border", "border-gray-400", "[clip-path:polygon(15px_0,calc(100%_-_15px)_0,100%_15px,100%_calc(100%_-_15px),calc(100%_-_15px)_100%,15px_100%,0_calc(100%_-_15px),0_15px)]");
    addUserInputBox("input",document.getElementById('input').value, messageCanvas);
    
    messageCanvas.scrollTop = messageCanvas.scrollHeight;
    var responseCanvas = document.createElement("div");
    responseCanvas.classList.add("my-2","overflow-y-auto","scrollbar", "rounded-lg");
    responseCanvas.id = generateUUID();
    messageCanvas.appendChild(responseCanvas);
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
    let responseId = "";
    let contentTables = [];
    let tableBodyMap = new Map();

    while (!done) {
      messageCanvas.scrollTop = messageCanvas.scrollHeight;
      fabricCanvas.scrollTop = fabricCanvas.scrollHeight;
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
            console.log(objVal,"+++++++++++++++++++++++++++++++++++++++++++++");
            messCahe = null;
          } catch (error) {
            console.log(error,"++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++000000000000errorerrorerror");
            if (messCahe !== null) {
              messCahe = messCahe + decodedValue;
            } else {
              messCahe = decodedValue;
            }
            console.log(messCahe,"++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++000000000000");
            try {
              strval = "[" + messCahe + "]";
              objVal = JSON.parse(strval);
            } catch (error) {
              continue;
            }
            console.error("Error parsing JSON: Will concartenate with previous messages", error);
          }

          objVal.forEach(responseJson => {
            console.log(responseJson,"-------------------------",responseJson.message);
            if (responseJson !== null && responseJson.message !== undefined) {
              console.log("Message:>>>>>>>>>>>>>>>>>>>>>>>");
              addchatDataErrBox(responseJson.message, responseCanvas);
              return;
            }
            console.log("==============================================");
            console.log(responseJson.result);
            console.log(responseJson.result.output);
            if (responseJson === null || responseJson.result === null || responseJson.result.output === null || responseJson.result.output === undefined) {
              addchatDataErrBox("Error while querying data product either there is no data for this query or there is some internal server error, please try again or contact the data product owner",
                responseCanvas);
              return;
            }
            console.log(responseJson.result.output.event);
            switch (responseJson.result.output.event) {
              case "END":
                if (responseJson.result.output.data.final.reason === "SUCCESSFULL") {
                  addchatDataSuccessBox(responseJson.result.output.data.content, responseCanvas);
                } else if (responseJson.result.output.data.final.reason === "DONE") {
                  break;
                } else {
                  addchatDataErrBox(JSON.stringify({
                    key: responseJson.result.output.data.final.reason,
                    value: responseJson.result.output.data.final.metadata
                  }), responseCanvas);
                }
                break;
              case "ABORTED":
                addchatDataErrBox(JSON.stringify({
                  key: responseJson.result.output.data.final.reason,
                  value: responseJson.result.output.data.final.metadata
                }), responseCanvas);
                break;
              case "DATASET_END":
                const datadiv = document.createElement("div");
                datadiv.classList.add("overflow-x-auto","scrollbar" ,"border-gray-800","p-2");
                const dataTableElem = document.getElementById(responseJson.result.output.data.contentId);
                datadiv.appendChild(dataTableElem);
                responseCanvas.appendChild(datadiv);
                console.log("tableMap:>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", tableMap);
                let dataset = tableMap.get(responseJson.result.output.data.contentId);
                let csvColumns = tableHeaderMap.get(responseJson.result.output.data.contentId);
                console.log("CSV Columns:>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", csvColumns);
                console.log("Dataset:>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", dataset);
                renderDatatable(responseJson.result.output.data.contentId, null, csvColumns, dataset)
                break;
              case "DATASET_START":
                let table = addTable(responseCanvas, responseJson.result.output.data.contentId);
                table.classList.add("stripe","row-border","compact");
                contentTables.push(responseJson.result.output.data.contentId);
                let dataTableBody = addTableBody(table);
                tableBodyMap.set(responseJson.result.output.data.contentId, dataTableBody);
                dataFields = JSON.parse(responseJson.result.output.data.content);
                tableHeaderMap.set(responseJson.result.output.data.contentId, dataFields);
                console.log(document.getElementById(table.id),"Table ID:>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", table.id);
                break;
              case "DATASET":
                handleDataset(tableMap,responseJson.result.output.data.content,  responseJson.result.output.data.contentId);
                break;
              case "FILE_DATA":
                if (responseJson.result.output.files) {
                  if (responseJson.result.output.files.name.length > 0) {
                    renderNabhikFiles(responseJson.result.output.files, responseCanvas, responseId);
                    handleNabhikFile(responseCanvas, responseJson.result.output.files, responseJson.result.output.files.description, tokenKey, downloadFileUrl,false);
                  }
                }
                break;
              default:
                if (responseJson.result.output.data.contentType === "JSON") {
                  let content = JSON.parse(responseJson.result.output.data.content);
                  addchatDataBox(content.key, content.value, responseCanvas, false);
                } else if (responseJson.result.output.data.contentType === "CODE") {
                  addchatCodeBox(responseJson.result.output.data.content, responseCanvas);
                } else {
                  addchatDataBox(null,responseJson.result.output.data.content, responseCanvas, false);
                }
                break;
            }
          });

        } catch (error) {
          console.error(error);
          addchatDataErrBox("Error while processing your request, please try again", responseCanvas);
        } finally {
          messageCanvas.appendChild(responseCanvas);
          fabricCanvas.appendChild(messageCanvas);
          console.log("Done reading data");
        }
      }
      done = readerDone;
    }
    try {
      let timeDiff = getEpochTime() - startTime;
      console.log(responseCanvas,"++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++000000000000");
      messageCanvas.appendChild(responseCanvas);
      fabricCanvas.appendChild(messageCanvas);
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

function handleDataset(tableMap,dataset, responseId) {
  let result = parseAndreturnJSONArrayElem(dataset);
  let datasetList = tableMap.get(responseId);
  if (datasetList === undefined) {
    datasetList = result[0];
  } else {
    datasetList = datasetList.concat(result[0]);
  }
  tableMap.set(responseId, datasetList);
  addDataSetInIndexDb(responseId, dataset)
  return 0;
}

function addDataSetInIndexDb(iden, dataSet) {
  try {
    const payload = {
      queryId: iden,
      value: dataSet
    }
    localDbConn.putData(payload);
  } catch (error) {
    console.error("Error storing dataset locally:", error);
  }
}

async function exportResultSet(responseId) {
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

function showPrompts() {
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

function closePrompts(id) {
  document.getElementById(id).classList.add("hidden");
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


function handleNabhikFile(canvas, file, description, accessTokenKey, downloadUrl,showContent=false) {
  const mainBox = document.createElement("div");
  mainBox.classList.add("m-2");
  const nBox = document.createElement("div");
  nBox.classList.add("flex","items-center", "text-xs","cursor-pointer", "text-gray-700");
  var svgFiles = document.createElement("span");
  svgFiles.classList.add("w-6", "h-6");
  svgFiles.innerHTML = `
  <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
    <path fill-rule="evenodd" d="M6 2a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6H6zm7 1.5V8h4.5L13 3.5zm-3 8.75a2.25 2.25 0 1 1 0 4.5 2.25 2.25 0 0 1 0-4.5zm-3 5.25a1.75 1.75 0 1 1 0-3.5 1.75 1.75 0 0 1 0 3.5zm8-5.25a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3z" clip-rule="evenodd"/>
  </svg>
  `;
  // nBox.appendChild(svgFiles);
  console.log("File:", file);
  console.log("Description:", description);
  console.log("Download URL:", downloadUrl);
  console.log("Access Token Key:", accessTokenKey);
  console.log("Canvas:", canvas);
  console.log("Show Content:", showContent);
  var assetContext = document.createElement("span");
  downloadNabhikFile(file.path, accessTokenKey, canvas,downloadUrl,showContent, false);
  assetContext.classList.add("ml-2");
  assetContext.textContent = description;
  nBox.appendChild(assetContext);
  const downloadLink = document.createElement("a");
  downloadLink.classList.add("text-gray-700","text-sm","rounded-lg","px-2","py-1","border-[2px]","border-gray-300","cursor-pointer"); 
  downloadLink.href = "#";
  // downloadLink.textContent = text;
  downloadLink.onclick = () => downloadNabhikFile(file.path, accessTokenKey, canvas,downloadUrl,false, true);
  downloadLink.appendChild(svgFiles);
  nBox.appendChild(downloadLink);
  mainBox.appendChild(nBox);
  canvas.appendChild(mainBox);
  canvas.scrollTop = canvas.scrollHeight;
}

async function downloadNabhikFile(fileName, accessTokenKey, canvas,downloadUrl, showContent=false, download=false) {
  const myHeaders = new Headers();
  myHeaders.append("Accept", "application/x-ndjson");
  myHeaders.append("Content-Type", "application/x-ndjson");
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
          buildAndDownloadFile(file.name, file.format, file.data,canvas,showContent,download);
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
    showInfoMessage("Message Cancelled", "The current query has been cancelled.");
  }
}

window.toggleInputType = toggleInputType;
window.exportResultSet = exportResultSet;
window.showPrompts = showPrompts;
window.closePrompts = closePrompts;
window.selectPromptIntoInput = selectPromptIntoInput;
window.uploadDataset = uploadDataset;
window.cancelStream = cancelStream;
window.queryDataserverAction = queryDataserverAction;
window.clearAttachDataset = clearAttachDataset;
window.selectPromptIntoInput = selectPromptIntoInput;
window.retryInput = retryInput;

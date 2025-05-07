import { dataFileUploader } from './fileloader.js';

var agentStudioDataset = [];

async function uploadAgentDataset(api, tokenKey, resource, fileCanvasId, resourceId) {
    if (resourceId === "") {
        resourceId = generateUUID();
    }
    await dataFileUploader(api, tokenKey, resource, resourceId, fileCanvasId, agentStudioDataset);
}

async function validateAgent(tokenKey, apiUrl, payloadDiv,canvasId) {
    alert("Validating agent");
    const payloadStr = document.getElementById(payloadDiv).innerHTML;
    const payload = JSON.stringify(JSON.parse(payloadStr));
    const startTime = getEpochTime();
    var runId = generateUUID();
    const myHeaders = new Headers();
    myHeaders.append("Accept", "application/x-ndjson");
    myHeaders.append("Content-Type", "application/x-ndjson");
    const apiToken = getCookie(tokenKey);
    myHeaders.append("Authorization", `Bearer ${apiToken}`);
    let fabricCanvas = document.getElementById(canvasId);
    fabricCanvas.innerHTML = "";
    const reasoningCanvas = document.createElement("details");
    reasoningCanvas.classList.add("bg-white", "p-2", "rounded-lg", "my-2", "text-xs", "text-blue-700", "chat-loader", "font-semibold","cursor-pointer");
    reasoningCanvas.id = "reasoning" + runId;
    reasoningCanvas.innerHTML = `<summary class="text-xs text-gray-900 font-semibold" id="reasoning-summary-${runId}">Thinking...</summary>`;
    reasoningCanvas.open = true;
    const messageCanvasElem = document.createElement("div");
    messageCanvasElem.id = runId;
    messageCanvasElem.classList.add("my-2","overflow-y-auto","w-[100%]");
    messageCanvasElem.appendChild(reasoningCanvas);
    fabricCanvas.appendChild(messageCanvasElem);
    const messageCanvas = document.getElementById(runId);
    const response = await fetch(apiUrl, {
        method: "POST",
        headers: myHeaders,
        redirect: "follow",
        body: payload,
    });
    if (!response.ok) {
        addstreamDataErrBox("Error while querying data product either there is no data for this query or there is some internal server error, please try again or contact the data product owner",
            messageCanvas);
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

          objVal.forEach(responseJson => {
            console.log("Response JSON: ", responseJson);
            if (responseJson !== null && responseJson.message !== "") {
              addstreamDataErrBox(responseJson.message, messageCanvas);
              return;
            }
            if (responseJson === null || responseJson.result === null || responseJson.result.output === null || responseJson.result.output === undefined) {
              addstreamDataErrBox("Error while querying data product either there is no data for this query or there is some internal server error, please try again or contact the data product owner",
                messageCanvas);
              return;
            }
            switch (responseJson.result.output.event) {
              case "END":
                if (responseJson.result.output.data.final.reason === "SUCCESSFULL") {
                  addstreamDataSuccessBox(responseJson.result.output.data.content, messageCanvas);
                } else if (responseJson.result.output.data.final.reason === "DONE") {
                  break;
                } else {
                  addstreamDataErrBox(JSON.stringify({
                    key: responseJson.result.output.data.final.reason,
                    value: responseJson.result.output.data.final.metadata
                  }), messageCanvas);
                }
                break;
              case "ABORTED":
                addstreamDataErrBox(JSON.stringify({
                  key: responseJson.result.output.data.final.reason,
                  value: responseJson.result.output.data.final.metadata
                }), messageCanvas);
                break;
              case "START":
                setChatId(responseJson.result.chatId);
                break;
              case "CONTENT_HEAD":
                  addstreamHeadingBox(responseJson.result.output.data.content,responseJson.result.output.data.content.contentId, messageCanvas);
                  break;
              case "TEXT_CONTENT":
                console.log("Text content type ========>>>>>>>>>>>>>>>>>>>>", responseJson.result.output.data.contentType);
                if (responseJson.result.output.data.contentType === "JSON") {
                  const [content, isArray] = isJSONorJSONArray(responseJson.result.output.data.content);
                  if (isArray) {
                    content.forEach(element => {
                      addstreamDataBox("", element, messageCanvas, false);
                    });
                  } else {
                    addstreamDataBox("", content.value, messageCanvas, false);
                  }
                  // addstreamDataBox(content.key, content.value, messageCanvas, false);
                } else if (responseJson.result.output.data.contentType === "FORMATTED_CONTENT") {
                  addstreamMarkdownDataBox(responseJson.result.output.data.content,responseJson.result.output.data.contentId, messageCanvas);
                } else {
                  addstreamDataBox("Info", responseJson.result.output.data.content, messageCanvas, false);
                }
                break;
              case "TITLE":
                addstreamTitleBox(responseJson.result.output.data.content, messageCanvas);
                break;
              case "CONTENT_SUBHEAD":
                addstreamSubHeadingBox(responseJson.result.output.data.content, responseJson.result.output.data.content.contentId, messageCanvas);
                break;
              case "CONTENT_DATA":
                addstreamHeadingContentBox(responseJson.result.output.data.content, responseJson.result.output.data.content.contentId, messageCanvas);
                break;
              case "DATASET_END":
                const datadiv = document.createElement("div");
                datadiv.classList.add("overflow-x-auto", "w-auto", "border-gray-800");
                let dataTable = document.getElementById(responseJson.result.output.data.contentId);
                datadiv.appendChild(dataTable);
                messageCanvas.appendChild(datadiv);
                addTableFooter(responseJson.result.output.data.contentId);
                // addstreamDataBox("Download", "You can select the file format for the report and then click on the download button below", messageCanvas, false);
                addDatasetDownloadButton(responseJson.result.output.data.contentId, messageCanvas);
                break;
              case "DATASET_START":
                // addDataSetHead(messageCanvas, "DataSet Id: ", responseId);
                let table = addTable("dataServerOutput", responseJson.result.output.data.contentId);
                contentTables.push(responseJson.result.output.data.contentId);
                let dataTableBody = addTableBody(table);
                tableBodyMap.set(responseJson.result.output.data.contentId, dataTableBody);
                dataFields = JSON.parse(responseJson.result.output.data.content);
                addTableHeader(table, dataFields);
                break;
              case "STATE":
                addstreamStateBox(responseJson.result.output.data.content,responseJson.result.output.data.contentType, reasoningCanvas);
                break;
              case "FILE_DATA":
                renderFiles(responseJson.result.output.files, messageCanvas, responseId);
                break;
              case "RESPONSE_ID":
                responseId = responseJson.result.responseId;
                messageCanvas.id = responseId;
                break;
              default:
                if (responseJson.result.output.data.contentType === "JSON") {
                  let content = JSON.parse(responseJson.result.output.data.content);
                  addstreamDataBox(content.key, content.value, messageCanvas, false);
                } else if (responseJson.result.output.data.contentType === "DATASET") {
                  counter = handleDataset(responseJson.result.output.data.dataset, tableBodyMap.get(responseJson.result.output.data.contentId), dataFields, responseJson.result.output.data.contentId, counter);
                  counter++;
                } else if (responseJson.result.output.data.contentType === "CLICKSET") {
                  addstreamClickableBox(responseJson.result.output.data.content, messageCanvas,submitInput);
                } else {
                  addstreamDataBox("Info", responseJson.result.output.data.content, messageCanvas, false);
                }
                break;
            }
          });

        } catch (error) {
          console.error(error);
          addstreamDataErrBox("Error while processing your request, please try again", messageCanvas);
        } finally {
          console.log("Done reading data");
        }
      }
      done = readerDone;
    }
    try {
      reasoningCanvas.classList.remove("chat-loader");
      let timeDiff = getEpochTime() - startTime;
      addstreamStateBox("Done: Time - " + timeDiff + "sec", "PLAIN_TEXT", reasoningCanvas);
      reasoningCanvas.open = false;
    //   fabricCanvas.appendChild(messageCanvas);
      try {
        document.getElementById("reasoning-summary-" + runId).innerHTML = "Thoughts";
      } catch {
        console.error("Error setting reasoning summary");
      }
      if (contentTables.length > 0) {
        contentTables.forEach(tableId => {
          setTimeout(() => {
            addTablePagination({
              tableId: tableId,
              rowsPerPage: 10,
              prevPageBtn: document.getElementById(tableId + "-prevPage"),
              nextPageBtn: document.getElementById(tableId + "-nextPage"),
              currentPageSizeSpan: document.getElementById(tableId + "-currentPageSize"),
              totalElementsSpan: document.getElementById(tableId + "-totalElements"),
            });
            // addDataSetInIndexDb(responseId, dbBatch);
          }, 2 * DatafabricDatabase.waitTime);
        });
      }
    } catch (error) {
      console.error(error);
  } finally {
    
  }
}

window.uploadAgentDataset = uploadAgentDataset;
window.validateAgent = validateAgent;
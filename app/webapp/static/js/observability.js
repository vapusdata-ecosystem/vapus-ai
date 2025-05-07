function appendUserInput(user, content, canvas) {
  let seperator = document.createElement("div");
  seperator.style.margin = "4px 0";
  seperator.classList.add("flex", "justify-end");
  seperator.innerHTML = `<div class="bg-gray-200 font-semibold break-words px-4 py-2 rounded-lg mt-3 mb-2 w-full">
                            <p class="text-xs text-blue-700 mt-1">${user}:</p>
                            <p class="text-sm text-gray-700-500"> ${content}</p>
                        </div>`;
  canvas.appendChild(seperator);
}

export async function dataproductObservabilityAction(tokenKey, apiUrl, username) {
  const input = document.getElementById('input').value;
  document.getElementById('input').disabled = true;
  const dpId = document.getElementById('dataProduct').value;
  if (input === "") {
    showErrorMessage(" Please enter a valid input");
    return;
  }
  apiUrl = apiUrl + "?q=" + input;
  if (dpId !== "") {
    apiUrl = apiUrl + "&dataProducts=" + dpId;
  }

  const myHeaders = new Headers();
  myHeaders.append("Accept", "application/x-ndjson");
  myHeaders.append("Content-Type", "application/x-ndjson");
  const apiToken = getCookie(tokenKey);
  myHeaders.append(
    "Authorization",
    `Bearer ${apiToken}`);
  return queryObservabilityAPI(apiUrl, myHeaders, username);
}

async function queryObservabilityAPI(url, myHeaders, username) {
  try {
    let fabricCanvas = document.getElementById("metricOutput");
    fabricCanvas.classList.remove("hidden");
    const messageCanvasElem = document.createElement("div");
    messageCanvasElem.classList.add("my-2");
    const tempMessid = generateUUID();
    messageCanvasElem.id = tempMessid;
    fabricCanvas.appendChild(messageCanvasElem);
    const messageCanvas = document.getElementById(tempMessid);
    appendUserInput(username, document.getElementById('input').value, messageCanvas);
    document.getElementById('input').value = "";
    const response = await fetch(
      url,
      {
        method: "GET",
        headers: myHeaders,
      }
    );

    if (!response.ok) {
      showErrorMessage("Error while querying data product either there is no data for this query or there is some internal server error, please try again or contact the data product owner");
      return;
    }
    const jsonResponse = await response.json();
    if (response.message !== null && response.message !== undefined) {
      addstreamDataErrBox(response.message, messageCanvas);
      return;
    }
    if (
      jsonResponse.output === null || jsonResponse.output === undefined
    ) {
      addstreamDataErrBox("Error while querying data product either there is no data for this query or there is some internal server error, please try again or contact the data product owner",
        messageCanvas);
      return;
    }

    if (jsonResponse.output.metrics.length < 1) {
      addstreamDataErrBox("Error while querying data product either there is no data for this query or there is some internal server error, please try again or contact the data product owner",
        messageCanvas);
      return;
    }
    const responseId = jsonResponse.output.responseId;
    addstreamDataBox("Metrics", "The following metrics are available for the data product", messageCanvas, false);
    addstreamDataBox("Total Records", jsonResponse.output.totalPulls, messageCanvas, false);
    addstreamDataBox("Last Pull At", jsonResponse.output.lastPullAt, messageCanvas, false);
    const dataTable = addTable("metricOutput", responseId);
    const dataTableBody = addTableBody(dataTable);
    const columns = [
      "dataProducts",
      "status",
      "error",
      "resultLength",
      "queryer",
      "queryerOrganization",
      "queriedAt",
      "query",
      "dataStoreQuery",
    ];
    addTableHeader(dataTable, columns);
    let rowbatch = [];
    jsonResponse.output.metrics.forEach(metric => {
      rowbatch.push({
        dataProducts: metric.dataProducts,
        status: metric.status,
        error: metric.error,
        resultLength: metric.resultLength,
        queryer: metric.queryer,
        queryerOrganization: metric.queryerOrganization,
        queriedAt: metric.queriedAt,
        query: metric.query,
        dataStoreQuery: metric.dataStoreQuery,
      });
      if (rowbatch.length > 100) {
        addTableRow(dataTableBody, columns, rowbatch);
        rowbatch = [];
      }
    });

    addTableFooter(dataTable, responseId);
    addTablePagination({
      tableId: dataTable.id,
      rowsPerPage: 10,
      prevPageBtn: document.getElementById(responseId + "-prevPage"),
      nextPageBtn: document.getElementById(responseId + "-nextPage"),
      currentPageSizeSpan: document.getElementById(responseId + "-currentPageSize"),
      totalElementsSpan: document.getElementById(responseId + "-totalElements"),
    });
  } catch (error) {
    console.error("Error: ", error);
  }
}


window.dataproductObservabilityAction = dataproductObservabilityAction;
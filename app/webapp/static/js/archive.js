function getEndpoint(dpId, streamChat) {
  let conso = document.getElementById('actionUris').innerText;
  let actionUris = JSON.parse(document.getElementById('actionUris').innerText);
  let containerEP = document.getElementById('address-' + dpId).innerText;
  let result = URL.parse(containerEP)
  if (streamChat === 'true') {
    result.pathname = actionUris.QueryStreamServerUrl;
  } else {
    result.pathname = actionUris.QueryServerUrl;
  }
  return result.toString();
}

async function displayResponseTable(response) {
    try {
      const hiddenQueryResultSet = document.getElementById('hiddenQueryResultSet');
      const exportAction = document.getElementById('exportAction');
      // Clear previous table rows
      var resultDataSet = [];
      const data = JSON.parse(response);
      if (data === null || data.body === null || data.body.length === 0) {
        handleNoResponseOrErr(data.metadata);
      }
      if (data.body.length > 0) {
        columns = Object.keys(data.body[0]);
        headerRow = document.createElement('tr')
        columns.forEach(column => {
          headerCell = document.createElement('th');
          headerCell.innerHTML = column;
          headerCell.classList.add('px-2', 'py-2', 'text-left', 'text-sm', 'font-medium', 'text-gray-700', 'border-b');
          headerRow.appendChild(headerCell);
        });
        responseTableSetHeader.appendChild(headerRow);
      }
      data.body.forEach(vals => {
        const row = document.createElement('tr');
        row.classList.add('hover:bg-gray-100');
        columns.forEach(column => {
          cell = document.createElement('td');
          try {
            cell.innerHTML = JSON.parse(vals[column]);
          } catch (error) {
            cell.innerHTML = vals[column];
          }
          cell.classList.add('px-2', 'py-2', 'text-sm', 'text-gray-600', 'border-b');
          row.appendChild(cell);
        });
        responseTableBody.appendChild(row);
      });
      outputMetadata.innerHTML = `<h2 class="text-medium font-semibold mb-2">Metadata</h2>`
      data.metadata.forEach(val => {
        const metaDiv = document.createElement('div');
        key = toTitleCase(val.key);
        metaDiv.innerHTML = `
            <span class="px-4 py-2 'text-sm' text-white "><strong>${key}</strong>: </span>
              <span class="px-4 py-2 'text-sm' text-white break-words">${val.value}</span>
        `;
        outputMetadata.appendChild(metaDiv);
      });
      outputMetadata.classList.remove('hidden');
      responseTableContainer.classList.remove('hidden');
      exportAction.classList.remove('hidden');
      document.getElementById('loadingOverlay').style.display = 'none';
    } catch (error) {
      handleNoResponseOrErr(error);
    }
  }
  
  function showDatastores(id) {
    const models = document.getElementById(`datastores-` + id);
    modelElems = document.getElementsByClassName("datastores");
    for (let i = 0; i < modelElems.length; i++) {
      modelElems[i].classList.add("hidden");
    }
    models.classList.remove("hidden");
  }

  

function handleNoResponseOrErr(errorMD) {
  let message = '';
  try {
    if ((errorMD === null) || (errorMD === undefined) || (errorMD.length === 0)) {
      message = 'Error While Processing the request, please try again';
    } else {
      message = message + errorMD[0].key + ': ' + errorMD[0].value;
    }
  } catch (error) {
    message = 'Error While Processing the request, please try again';
  }
  const responseTableContainer = document.getElementById('responseTableContainer');
  const responseTableBody = document.getElementById('responseTableBody');
  // Clear previous table rows
  responseTableBody.innerHTML = '';
  const row = document.createElement('tr');

  // row.innerHTML = `
  //       <td class="px-4 py-2 text-gray-700 bg-gray-100 border-b">${message}</td>
  //   `;
  row.innerHTML = `
    <td class="px-4 py-2 text-gray-700 bg-gray-100 border-b">Error while querying data product either there is no data for this query or there is some internal server error, please try again or contact the data product owner</td>
`;
  responseTableBody.appendChild(row);
  // Show the table container
  responseTableContainer.classList.remove('hidden');
}

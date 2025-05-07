export async function fileUploader(apiUrl, accessTokenKey, resource, resourceId) {
    return new Promise((resolve, reject) => {
        const fileInput = document.createElement("input");
        fileInput.type = "file";
        fileInput.accept = ".csv,application/json,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"; // Accept CSV, JSON, and XLSX files
        fileInput.style.display = "none";

        fileInput.addEventListener("change", async (event) => {
            const files = event.target.files;
            if (files.length > 0) {
                const file = files[0];
                const reader = new FileReader();

                // Read the file content
                reader.onload = async (e) => {
                    const fileContent = e.target.result;
                    try {
                        // Process file data
                        let data = strToUniArray(fileContent);
                        data = uint8ArrayToBase64(data);
                        const fileExt = getFileExtension(file.name);

                        const payload = {
                            resourceId: resourceId,
                            resource: resource,
                            objects: [
                                {
                                    name: file.name,
                                    data: data,
                                    format: fileExt.toUpperCase(),
                                },
                            ],
                        };

                        const myHeaders = new Headers();
                        myHeaders.append("Accept", "application/json");
                        myHeaders.append("Content-Type", "application/json");
                        const apiToken = getCookie(accessTokenKey);
                        myHeaders.append("Authorization", `Bearer ${apiToken}`);

                        const params = JSON.stringify(payload);

                        // Send the file upload request
                        const response = await fetch(apiUrl, {
                            method: "POST",
                            headers: myHeaders,
                            redirect: "follow",
                            body: params,
                        });

                        if (!response.ok) {
                            showErrorMessage("Upload Failed", "File upload failed, please try after some time");
                            reject("File upload failed");
                            return;
                        }

                        const resultBody = await response.json();
                        let fileResults = [];
                        if (resultBody.output.length > 0) {
                            resultBody.output.forEach((res) => {
                                fileResults.push({
                                    name: res.object.name,
                                    format: res.object.format,
                                    responsePath: res.responsePath,
                                    fid: res.fid
                                });
                            });
                        }
                        resolve(fileResults); // Resolve the promise with the result
                    } catch (error) {
                        console.error("Error uploading file:", error);
                        showAlert(AlertError, "Upload Failed", "File upload failed, please try after some time");
                        reject([]); // Reject the promise in case of errors
                    }
                };

                // Read the file as text
                reader.readAsText(file);
            } else {
                reject("No file selected"); // Reject if no file is selected
            }
        });

        // Trigger the file input
        fileInput.click();
    });
}

export async function dataFileUploader(apiUrl, tokenKey, resource, resourceId, fileCanvasId, datasetsList) {
  const fileCanvas = document.getElementById(fileCanvasId);
  const result = await fileUploader(apiUrl, tokenKey, resource, resourceId);
  console.log("File upload result: ", result);
  if (result.length > 0) {
      result.forEach(file => {
          if (file.name !== undefined && file.name !== "") {
              addDatasetFile(file, fileCanvas, datasetsList);
          }
      });
  }
}


export function addDatasetFile(file, fileCanvas, datasetsList) {
  const liItem = document.createElement("li");
  liItem.classList.add("flex", "p-1", "items-center", "justify-between");
  const filenamediv = document.createElement("div");
  filenamediv.textContent = file.name + "attach";
  filenamediv.id = file.fid + "-name";
  filenamediv.classList.add("p-1", "m-1",  "border-1", "text-gray-100", "border-orange-800",
      "text-xs", "shadow-sm", "rounded-md", "overflow-hidden", "break-words", "w-auto");
  const attachDiv = document.createElement("div");
  attachDiv.innerHTML = `<svg viewBox="0 0 24 24" width="1.2em" height="1.2em" class="size-3"><path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6z"></path></svg>`;
  attachDiv.classList.add("hidden","p-1", "m-1", "bg-orange-800", "border-1", "text-gray-100", "border-orange-800",  "cursor-pointer");
 
  const fileInfo = document.createElement("input");
  fileInfo.type = "hidden";
  fileInfo.value = file.responsePath;
  fileInfo.id = file.fid;
  attachDiv.onclick = function () {
      attachDataset(this, file.fid, datasetsList);
  };
  console.log(file);
  const existIndex = datasetsList.indexOf(file.responsePath);
  if (existIndex < 0) {
      datasetsList.push(file.responsePath);
      console.log("attachDataset ---- ", datasetsList);
      showInfoMessage("File Attached", "File attached successfully");
  }
  console.log("attachDataset ---- ", attachDiv, file.fid, datasetsList);
  const detachdiv = document.createElement("div");
  detachdiv.innerHTML = `<svg viewBox="0 0 24 24" width="1.2em" height="1.2em" class="size-3">
  <path fill="currentColor" d="M5 11h14v2H5z"></path>
</svg>`;
  detachdiv.classList.add( "p-1", "m-1", "text-gray-100", "border-1",  "border-orange-800", "cursor-pointer");
  detachdiv.onclick = function () {
      dettachDataset(this, file.fid, datasetsList);
  };
  detachdiv.id = file.fid + "-detach";
  attachDiv.id = file.fid + "-attach";

  liItem.appendChild(fileInfo);
  liItem.appendChild(filenamediv);
  liItem.appendChild(attachDiv);
  liItem.appendChild(detachdiv);
  fileCanvas.appendChild(liItem);
}

function attachDataset(el, filesetId, datasetsList) {
  const fileDiv = document.getElementById(filesetId + "-name");
  if (fileDiv === null) {
      return;
  }
  // fileDiv.classList.add("fileuploader-list-selected");
  const filePath = document.getElementById(filesetId);
  if (filePath === null) {
      return;
  }
  console.log("attachDataset ---- ", filePath.value, datasetsList);
  const existIndex = datasetsList.indexOf(filePath.value);
  if (existIndex < 0) {
      datasetsList.push(filePath.value);
      console.log("attachDataset ---- ", datasetsList);
      showInfoMessage("File Attached", "File attached successfully");
  }
  el.classList.add("hidden");
  const detachdiv = document.getElementById(filesetId + "-detach");
  detachdiv.classList.remove("hidden");
}

export function dettachDataset(el, filesetId, datasetsList) {
  const fileDiv = document.getElementById(filesetId + "-name");
  if (fileDiv === null) {
      return;
  }
  // fileDiv.classList.remove("fileuploader-list-selected");
  const filePath = document.getElementById(filesetId);
  if (filePath === null) {
      return;
  }
  const existIndex = datasetsList.indexOf(filePath.value);
  if (existIndex > -1) {
      datasetsList.pop(filePath.value);
      showInfoMessage("File Dettached", "File dettached successfully");
  }
  el.classList.add("hidden");
  const attachDiv = document.getElementById(filesetId + "-attach");
  attachDiv.classList.remove("hidden");
}

function CleanDatasetAttachment(el, filecanvas, datasetsList) {
  const fileCanvas = document.getElementById(filecanvas);
  if (!fileCanvas) {
      console.error("File canvas element not found");
      return;
  }

  // Clear datasets list
  datasetsList.length = 0;
  
  // Get all li items in the fileCanvas
  const liItems = fileCanvas.querySelectorAll('li');
  
  // Iterate through each li item
  liItems.forEach(liItem => {
      // Find the file ID input element
      const fileInput = liItem.querySelector('input[type="hidden"]');
      if (fileInput) {
          const fileId = fileInput.id;
          
          // Reset the UI state for each item
          const nameElement = document.getElementById(fileId + "-name");
          if (nameElement) {
              nameElement.classList.remove("fileuploader-list-selected");
          }
          
          const attachDiv = document.getElementById(fileId + "-attach");
          if (attachDiv) {
              attachDiv.classList.remove("hidden");
          }
          
          const detachDiv = document.getElementById(fileId + "-detach");
          if (detachDiv) {
              detachDiv.classList.add("hidden");
          }
      }
  });
  
  const fileDiv = document.getElementById(filesetId + "-name");
  if (fileDiv === null) {
      return;
  }
  fileDiv.classList.remove("fileuploader-list-selected");
  const filePath = document.getElementById(filesetId);
  if (filePath === null) {
      return;
  }
  const existIndex = datasetsList.indexOf(filePath.value);
  if (existIndex > -1) {
      datasetsList.pop(filePath.value);
      showInfoMessage("File Attached", "File attached successfully");
  }
  el.classList.add("hidden");
  const attachDiv = document.getElementById(filesetId + "-attach");
  attachDiv.classList.remove("hidden");
}


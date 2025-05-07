let editor;

// Initialize CodeMirror
function initCodeMirror() {
    editor = CodeMirror(document.getElementById("yamlEditor"), {
        mode: "yaml",
        lineNumbers: true,
        lineWrapping: true,
        theme: "dracula",
        viewportMargin: Infinity,
        lint: true,
        gutters: ["CodeMirror-linenumbers", "CodeMirror-lint-markers"],
        indentUnit: 2, // YAML standard is 2 spaces
        tabSize: 2,
        autoCloseBrackets: true,
        matchBrackets: true,
        styleActiveLine: true,
        highlightSelectionMatches: {showToken: true, annotateScrollbar: true}
    });
    // editor.setSize("50%", "50%");
}

// Show loading overlay
function showLoading() {
    document.getElementById("loading-overlay").classList.remove("hidden");
}

// Hide loading overlay
function hideLoading() {
    document.getElementById("loading-overlay").classList.add("hidden");
}


async function VapusResourceActionHandler(specId, action, method, title, tokenKey, apiUrl,isRedirect) {
    if (action.toLowerCase() === "archive" || action.toLowerCase() === "delete") {
        ShowConfirm("VapusData Action",
                "Are you sure you want to delete this resource?",
                () => {
                console.log("Confirmed!");
                resourceGetCall(method, tokenKey, apiUrl)
                return;
                });
    } else if (action.toLowerCase() === "validate" || action.toLowerCase() === "publish" || action.toLowerCase() === "unpublish") {
        var jsonSpec = document.getElementById(specId).innerHTML;
        var jsPayload = JSON.parse(jsonSpec);
        resourcePostCall(method, tokenKey, apiUrl,jsonSpec);
        return;
    } else if (action.toLowerCase() === "sync") {
        resourceGetCall(method, tokenKey, apiUrl,null);
        return;
    }
    if (action.toLowerCase() === "create" || action.toLowerCase() === "update" || action.toLowerCase() === "deploy" || action.toLowerCase() === "add_users") {
        document.getElementById('actionTitle').innerHTML = "";
        document.getElementById('actionTitle').innerHTML = action.toUpperCase();
        document.getElementById('yamlSpecTitle').innerHTML = "";
        document.getElementById('yamlSpecTitle').innerHTML = title;
        openYAMLedModal(apiUrl,
            tokenKey,
            specId, method);
    }
}

async function CreateResourceViaYaml(specId,method, title, tokenKey, apiUrl) {
    document.getElementById('yamlSpecTitle').innerHTML = "";
    document.getElementById('yamlSpecTitle').innerHTML = title;
    console.log(apiUrl, tokenKey, specId, method,"================================");
    renderYamlSection(apiUrl,
        tokenKey,
        specId, method);
}

function renderYamlSection(apiUrl, tokenKey, contentDiv) {
    console.log(apiUrl);
    // params = JSON.stringify({ apiUrl: apiUrl, tokenKey: tokenKey });
    // document.getElementById("yamlParams").innerText = params;
    // console.log(document.getElementById("yamlParams").innerText);

    if (!editor) {
        initCodeMirror();
    }
    if (contentDiv != null) {
        val = document.getElementById(contentDiv).innerText;
        if (val) {
            // Trim whitespace before setting the value in the editor
            editor.setValue(val.trim());
        }
    }
}

async function resourceGetCall( method ,tokenKey, apiUrl) {
    try {
        const apiToken = getCookie(tokenKey);
        const myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", `Bearer ${apiToken}`);
        const requestOptions = {
            method: method,
            headers: myHeaders,
            redirect: "follow"
        };
        const response = await fetch(apiUrl, requestOptions);
        // Process the response
        if (!response.ok) {
            const errorData = await response.text();
            // If the response is JSON, you can parse it (optional)
            message = handleInappResponseError(errorData);
            showErrorMessage("VapusData Action", message);
        } else {
            const result = await response.json();
            hideLoading();  // Hide the loading overlay
            showInfoMessage("VapusData Action","Resource action executed successfully");
            location.reload(true);
        }

    } catch (error) {
        console.log(error);
        showAlert(AlertError, "VapusData Action", "Execution failed due to an error");
    } finally {
        hideLoading();
    }
}

async function resourcePostCall(method ,tokenKey, apiUrl,payload) {
    try {
        console.log('Payload:', payload);
        const apiToken = getCookie(tokenKey);
        const myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", `Bearer ${apiToken}`);
        const requestOptions = {
            method: method,
            headers: myHeaders,
            redirect: "follow",
        };
        if (payload  !== null || payload !== undefined) {
            requestOptions.body = payload;
        }
        const response = await fetch(apiUrl, requestOptions);
        // Process the response
        if (!response.ok) {
            const errorData = await response.text();
            // If the response is JSON, you can parse it (optional)
            message = handleInappResponseError(errorData);
            showErrorMessage("VapusData Action", message);
        } else {
            const result = await response.json();
            hideLoading();  // Hide the loading overlay
            showInfoMessage("VapusData Action","Resource action executed successfully");
            location.reload(true);
        }

    } catch (error) {
        console.log(error);
        showAlert(AlertError, "VapusData Action", "Execution failed due to an error");
    } finally {
        hideLoading();
    }
}


// Function to handle API call and show loading spinner
async function VapusDataAct() {
    const mParams = document.getElementById("modalParams").innerText;
    const params = JSON.parse(mParams);
    showLoading();  // Show loading overlay
    yamlContent = getYAMLEditorVal();  // Get the YAML content from the editor
    let jsonPayload;
    try {
        jsonPayload = jsyaml.load(yamlContent);
        console.log('JSON:', jsonPayload);
    } catch (error) {
        showAlert(AlertError, "VapusData Action", "Invalid YAML content");
        hideLoading();
        return;
    }
    console.log('JSON:', jsonPayload);
    let method = document.getElementById("apiMethod").value;
    if (method === "") {
        method = "GET";
    }
    try {
        // Make the API call
        const apiToken = getCookie(params.tokenKey);
        const myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", `Bearer ${apiToken}`);

        const payload = JSON.stringify(jsonPayload);
        const requestOptions = {
            method: method,
            headers: myHeaders,
            body: payload,
            redirect: "follow"
        };
        const response = await fetch(params.apiUrl, requestOptions);
        // Process the response
        if (!response.ok) {
            const errorData = await response.text();
            // If the response is JSON, you can parse it (optional)
            message = handleInappResponseError(errorData);
            showAlert(AlertError, "VapusData Action", message);
            // throw new Error('Action failed');
        } else {
            const result = await response.json();
            hideLoading();  // Hide the loading overlay
            showAlert(AlertInfo, "VapusData Action", result.dmResp.message);
            location.reload(true);
        }

    } catch (error) {
        console.log(error);
        showAlert(AlertError, "VapusData Action", "Execution failed due to an error");
    } finally {
        hideLoading();
    }
}

// Function to handle API call and show loading spinner
async function YamlCreateResourceAction(params) {
    // const mParams = document.getElementById("modalParams").innerText;
    // const params = JSON.parse(mParams);
    showLoading();  // Show loading overlay
    yamlContent = getYAMLEditorVal();  // Get the YAML content from the editor
    let jsonPayload;
    try {
        jsonPayload = jsyaml.load(yamlContent);
        console.log('JSON:', jsonPayload);
    } catch (error) {
        showAlert(AlertError, "VapusData Action", "Invalid YAML content");
        hideLoading();
        return;
    }
    console.log('JSON:', jsonPayload);
    if (params.method === "") {
        method = "GET";
    }
    try {
        // Make the API call
        const apiToken = getCookie(params.tokenKey);
        const myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");
        myHeaders.append("Authorization", `Bearer ${apiToken}`);

        const payload = JSON.stringify(jsonPayload);
        const requestOptions = {
            method: params.method,
            headers: myHeaders,
            body: payload,
            redirect: "follow"
        };
        const response = await fetch(params.apiUrl, requestOptions);
        // Process the response
        if (!response.ok) {
            const errorData = await response.text();
            // If the response is JSON, you can parse it (optional)
            message = handleInappResponseError(errorData);
            showAlert(AlertError, "VapusData Action", message);
            // throw new Error('Action failed');
        } else {
            const result = await response.json();
            hideLoading();  // Hide the loading overlay
            showAlert(AlertInfo, "VapusData Action", result.dmResp.message);
            location.reload(true);
        }

    } catch (error) {
        console.log(error);
        showAlert(AlertError, "VapusData Action", "Execution failed due to an error");
    } finally {
        hideLoading();
    }
}

function getYAMLEditorVal() {
    return editor.getValue();
}

// Load YAML file content into CodeMirror
function loadFileContent(event) {
    const file = event.target.files[0];
    if (file) {
        const reader = new FileReader();
        reader.onload = function (e) {
            editor.setValue(e.target.result);
        };
        reader.readAsText(file);
    }
}

// Initialize CodeMirror when the modal opens
function openYAMLedModal(apiUrl, tokenKey, contentDiv, apiMethod) {
    console.log(apiUrl);
    document.getElementById("yamlModal").classList.remove("hidden");
    params = JSON.stringify({ apiUrl: apiUrl, tokenKey: tokenKey });
    document.getElementById("modalParams").innerText = params;
    console.log(document.getElementById("modalParams").innerText);
    if (apiMethod !== null || apiMethod !== undefined) {
        document.getElementById("apiMethod").value = apiMethod;
    }

    if (!editor) {
        initCodeMirror();
    }
    if (contentDiv != null) {
        val = document.getElementById(contentDiv).innerText;
        if (val) {
            editor.setValue(val);
        }
    }
}

function openYAMLedModalWithText(apiUrl, tokenKey, text) {
    console.log(apiUrl);
    document.getElementById("yamlModal").classList.remove("hidden");
    params = JSON.stringify({ apiUrl: apiUrl, tokenKey: tokenKey });
    document.getElementById("modalParams").innerText = params;
    console.log(document.getElementById("modalParams").innerText);

    if (!editor) {
        initCodeMirror();
    }
    editor.setValue(text);
}

function closeYAMLedModal() {
    document.getElementById("yamlModal").classList.add("hidden");
    if (editor) {
        editor.setValue(""); // Clear the editor content
    }
}
// Adjust the textarea height based on content
function adjustYamlFieldHeight(textarea) {
    textarea.style.height = "auto";  // Reset the height
    textarea.style.height = (textarea.scrollHeight) + "px";  // Set the new height based on content
}
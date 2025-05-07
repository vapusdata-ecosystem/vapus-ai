function convertArrayKeys(obj) {
    if (Array.isArray(obj)) {
        return obj.map(convertArrayKeys);
    }
    else if (obj !== null && typeof obj === 'object') {
        const newObj = {};
        for (const key in obj) {
            if (Object.prototype.hasOwnProperty.call(obj, key)) {
                const match = key.match(/^(.*)\[(\d+)\]$/);
                if (match) {
                    const baseKey = match[1];
                    const index = parseInt(match[2], 10);
                    if (!newObj.hasOwnProperty(baseKey)) {
                        newObj[baseKey] = [];
                    }
                    newObj[baseKey][index] = convertArrayKeys(obj[key]);
                } else {
                    newObj[key] = convertArrayKeys(obj[key]);
                }
            }
        }
        return newObj;
    }
    return obj;
}

function getFormData(formData) {
    const data = {};

    for (const [key, value] of formData.entries()) {
        if (key.startsWith("spec.")){
            const keys = key.split('.');
            let current = data;
            keys.forEach((k, index) => {
                if (index === keys.length - 1) {
                    current[k] = value;
                } else {
                    if (!current[k]) {
                        current[k] = {};
                    }
                    current = current[k];
                }
            });
        }
    }
    const finalData = convertArrayKeys(data);
    return finalData.spec;
}

function convertNumericStrings(input) {
    if (typeof input === "string") {
        const commaSeparatedRegex = /^\s*-?\d+(\.\d+)?\s*,\s*(-?\d+(\.\d+)?\s*,\s*)*-?\d+(\.\d+)?\s*$/;
        console.log("input", input);
        if (commaSeparatedRegex.test(input)) {
            console.log("comma separated");
            return input.split(",").map(s => Number(s.trim()));
        }
        if (/^-?\d+(\.\d+)?$/.test(input)) {
            return Number(input);
        }
        return input;
    }

    if (Array.isArray(input)) {
        return input.map(convertNumericStrings);
    }

    if (input !== null && typeof input === "object") {
        const newObj = {};
        for (const key in input) {
            if (Object.prototype.hasOwnProperty.call(input, key)) {
                newObj[key] = convertNumericStrings(input[key]);
            }
        }
        return newObj;
    }
    return input;
}

function scanNetworkParams(formData) {
    console.log("Network Params", formData);
    if (formData.netParams) {
        console.log("Found netParams", formData.netParams.port);
        formData.netParams.port = formData.netParams.port[0];
    } else {
        formData.netParams = {};
        return formData;
    }
    // if (formData.netParams.databases !== "") {
    //     formData.netParams.databases = formData.netParams.databases.split(",");
    // } else {
    //     formData.netParams.databases = [];
    // }
    // if (formData.netParams.databasePrefixes !== "") {
    //     formData.netParams.databasePrefixes = formData.netParams.databasePrefixes.split(",");
    // } else {
    //     formData.netParams.databasePrefixes = [];
    // }
    if (formData.netParams.dsCreds.length > 0) {
        formData.netParams.dsCreds.forEach((dsCred) => {
            if (dsCred.priority) {
                dsCred.priority = dsCred.priority[0];
            }
        });
    }
    return formData;
}

function scanVapusSchedule(formData) {
    if (formData.metaSyncSchedule === undefined || formData.metaSyncSchedule === null) {
        formData.metaSyncSchedule = {};
        return formData;
    }
    if (formData.metaSyncSchedule.limit) {
        console.log("Found vapus schedule", formData.metaSyncSchedule.limit);
        formData.metaSyncSchedule.limit = formData.metaSyncSchedule.limit[0];
    }
    if (formData.metaSyncSchedule.cronTab) {
        if (formData.metaSyncSchedule.cronTab.frequencyTab) {
            if (formData.metaSyncSchedule.cronTab.frequencyTab.length > 0) {
                formData.metaSyncSchedule.cronTab.frequencyTab.forEach((frequencyTab) => {
                    if (frequencyTab.frequency) {
                        frequencyTab.frequencyInterval = frequencyTab.frequencyInterval[0];
                    }
                }
                );
            }
        }
        if (formData.metaSyncSchedule.runAt ==="") {
            formData.metaSyncSchedule.runAt = null;
        }
    }

    return formData;
}

 // Tagging Code 
function setupTagInput(inputId, tagsContainerId, existingTags = []) {
    const tagsContainer = document.getElementById(tagsContainerId);
    const input = document.getElementById(inputId);
    let tagContents = [];
    const baseInputName = input.getAttribute("name");
    
    // Populate existing tags
    existingTags.forEach(tag => createTag(tag, tagsContainer, tagContents, baseInputName));
    
    input.addEventListener('keydown', function (event) {
        if (event.key === 'Enter') {
            event.preventDefault();
            const tagContent = input.value.trim();
            createTag(tagContent, tagsContainer, tagContents, baseInputName);
            input.value = '';
        }
    });
    
    tagsContainer.addEventListener('click', function (event) {
        deleteTag(event, tagContents);
    });
} 

// Tag creation
function createTag(tagContent, tagsContainer, tagContents, baseInputName) {
    if (tagContents.length < 4 && tagContent !== '') {
        const tag = document.createElement('li');
        tag.className = `bg-gray-200 text-gray-800 px-3 py-1 rounded-full flex items-center ${baseInputName}`;
        
        tag.innerHTML = `
            ${tagContent}
            <button class="ml-2 text-red-500 hover:text-red-700 delete-button">&times;</button>
        `;

        tagsContainer.appendChild(tag);
        tagContents.push(tagContent);
    } else {
        console.log("Opps.......!");
        showErrorMessage("Tag", "Field cannot be greater than 4");
    }
}

// Tag deletion
function deleteTag(event, tagContents) {
    if (event.target.classList.contains('delete-button')) {
        const tagElement = event.target.parentNode;
        const tagText = tagElement.firstChild.textContent.trim();

        // Remove the tag from tagContents array
        const index = tagContents.indexOf(tagText);
        if (index !== -1) {
            tagContents.splice(index, 1);
        }
        
        tagElement.remove(); // Remove the tag from UI
        console.log(tagContents);
    }
}

// converting multiple tags into array of string
function getTagsByClass(className) {
    const tags = document.getElementsByClassName(className);
    return Array.from(tags).map(tag => {
        let value = tag.getAttribute("data-value");
        if (!value) {
            value = tag.firstChild.textContent.trim();
        }
        return value;
    });
}

async function submitCreateForm(dataObj, tokenKey, apiurl,listingUrl) {
    try {
        showLoadingButton();
        const payload = {
            spec: dataObj
        };
        try {
            const apiToken = getCookie(tokenKey);
            const myHeaders = new Headers();
            myHeaders.append("Content-Type", "application/json");
            myHeaders.append("Authorization", `Bearer ${apiToken}`);
            console.log("Header: ", myHeaders)
            console.log("payload: ", payload)
            console.log("apiToken: ", apiToken)
            const response = await fetch(apiurl, {
                method: 'POST',
                headers: myHeaders,
                body: JSON.stringify(payload)
            });

            if (response.ok) {
                showInfoMessage("Resource Created", "Resource created successfully.");
                const output = await response.json();
                console.log("Resource created:", output);
                var resourceInfo = output.result;
                if (resourceInfo) {
                    showInfoMessage("Resource created", "Resource creation initiation successfully.");
                    console.log("Listing URL: ", listingUrl);
                    redirectToLink(listingUrl);
                } else {
                    showInfoMessage("Resource Created", "Resource created successfully.");
                }
            } else {
                showErrorMessage("Resource Creation Failed", `Resource creation failed with status ${response.status}.`);
            }
        } catch (error) {
            console.error("Error sending API request:", error);
            showErrorMessage("Resource Creation Failed", `Resource creation failed`);
        }
    } catch (error) {
        console.error("Error parsing form data:", error);
        showErrorMessage("Source Connector Error", "Error parsing form data. Please check the console for more details.");
    } finally {
        hideLoadingButton();
        return false;
    }
}

async function submitUpdateForm(dataObj, tokenKey, apiurl, listingUrl) {
    try {
        showLoadingButton();
        const payload = {
            spec: dataObj
        };
        try {
            console.log("I am in UpdateForm")
            const apiToken = getCookie(tokenKey);
            const myHeaders = new Headers();
            myHeaders.append("Content-Type", "application/json");
            myHeaders.append("Authorization", `Bearer ${apiToken}`);
            console.log("Header: ", myHeaders)
            console.log("payload: ", payload)
            console.log("apiToken: ", apiToken)
            const response = await fetch(apiurl, {
                method: 'PUT',
                headers: myHeaders,
                body: JSON.stringify(payload)
            });

            if (response.ok) {
                showInfoMessage("Resource Created", "Resource created successfully.");
                const output = await response.json();
                console.log("Resource created:", output);
                // Redirecting
                let currLocation = window.location.href;
                let redirectAPI = currLocation.substring(0, currLocation.lastIndexOf("/"));
                console.log("redirectAPI:>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", redirectAPI);
                window.location.replace(redirectAPI)
            } else {
                showErrorMessage("Resource Creation Failed", `Resource creation failed with status ${response.status}.`);
            }
        } catch (error) {
            console.error("Error sending API request:", error);
            showErrorMessage("Resource Creation Failed", `Resource creation failed`);
        }
    } catch (error) {
        console.error("Error parsing form data:", error);
        showErrorMessage("Source Connector Error", "Error parsing form data. Please check the console for more details.");
    } finally {
        hideLoadingButton();
        return false;
    }
}

function showLoadingButton() {
    document.getElementById("submit").disabled = true;
    document.getElementById("submit").classList.add("hidden");
    document.getElementById("loading-overlay").classList.remove("hidden");
}

function hideLoadingButton() {
    document.getElementById("submit").disabled = false;
    document.getElementById("submit").classList.remove("hidden");
    document.getElementById("loading-overlay").classList.add("hidden");
}
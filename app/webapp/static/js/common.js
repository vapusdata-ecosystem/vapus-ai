const AlertSuccess = "success";
const AlertError = "error";
const AlertInfo = "info";

const currentUri = encodeURIComponent(window.location.href);

const copySvg = `
<svg viewBox="0 0 24 24" width="1.2em" height="1.2em" class="h-4 w-4 rounded-sm p-0.5 hover:bg-neutral-10.5">
    <path fill="currentColor" d="M19 21H8V7h11m0-2H8a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h11a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2m-3-4H4a2 2 0 0 0-2 2v14h2V3h12z">
    </path>
</svg>
`;

const itemSvg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="green" width="36px" height="36px">
                    <path d="M9 16.2l-3.5-3.5-1.4 1.4L9 19 20 8l-1.4-1.4z"></path>
                </svg>`

const suggestionSvg = `
<svg
  xmlns="http://www.w3.org/2000/svg"
  fill="none"
  stroke="currentColor"
  stroke-width="2"
  viewBox="0 0 24 24"
  stroke-linecap="round"
  stroke-linejoin="round"
  width="1.5em" height="1.5em"
>
  <!-- Sparkles icon -->
  <path d="M9 7c0-1.075.574-2.303 1.12-3.184.232-.38.347-.57.567-.746.193-.156.414-.28.652-.367.222-.173.438-.203.66-.203s.438.03.66.203c.238.087.46.211.652.367.22.177.335.367.567.746C13.426 4.697 14 5.925 14 7s-.574 2.303-1.12 3.184c-.232.38-.347.57-.567.746-.193.156-.414.28-.652.367-.222.173-.438.203-.66.203s-.438-.03-.66-.203c-.238-.087-.46-.211-.652-.367-.22-.177-.335-.367-.567-.746C9.574 9.303 9 8.075 9 7zm7.5 9.5c0-.713.38-1.528.744-2.112.155-.25.233-.375.383-.49.132-.1.285-.178.447-.232.207-.076.43-.116.675-.116s.468.04.675.116c.162.054.315.132.447.232.15.115.228.24.383.49.364.584.744 1.399.744 2.112 0 .713-.38 1.528-.744 2.112-.155.25-.233.375-.383.49-.132.1-.285.178-.447.232-.207.076-.43.116-.675.116s-.468-.04-.675-.116c-.162-.054-.315-.132-.447-.232-.15-.115-.228-.24-.383-.49-.364-.584-.744-1.399-.744-2.112zM3.5 13.5c0-.57.304-1.221.594-1.688.119-.19.179-.285.294-.372.101-.08.218-.144.342-.187.158-.058.328-.088.52-.088s.361.03.52.088c.124.043.241.107.342.187.115.087.175.182.294.372.29.467.594 1.118.594 1.688 0 .57-.304 1.221-.594 1.688-.119.19-.179.285-.294.372-.101.08-.218.144-.342.187-.158.058-.328.088-.52.088s-.361-.03-.52-.088c-.124-.043-.241-.107-.342-.187-.115-.087-.175-.182-.294-.372-.29-.467-.594-1.118-.594-1.688z" />
</svg>
`
const mimeTypes = {
    "png": "image/png",
    "csv": "text/csv",
    "json": "application/json",
    "yaml": "application/x-yaml",
    "xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    "pdf": "application/pdf",
    "html": "text/html",
    "txt": "text/plain",
};

const eligibleRenderedFileFormats = ["png", "pdf", "html","jpg","jpeg"];

function getEpochTime() {
    return Math.floor(new Date().getTime() / 1000);
}

function toggleSidebar() {
    alert("here");
    const sidebar = document.getElementById("sidebar");
    sidebar.classList.toggle("collapsed");
    alert(sidebar.classList)
}

function copyToClipboard(text) {
    // const el = document.createElement('textarea');
    // el.value = text;
    // document.body.appendChild(el);
    // el.select();
    navigator.clipboard.writeText(text);
    const toast = document.getElementById('toast');
    toast.textContent = "copied to clipboard";
    toast.classList.add('show');

    // Hide the toast after 2 seconds
    setTimeout(() => {
        toast.classList.remove('show');
    }, 1000);
    // document.body.removeChild(el);
}

function copyDivToClipboard(divId) {
    const text = document.getElementById(divId).innerHTML;
    navigator.clipboard.writeText(text);
    const toast = document.getElementById('toast');
    toast.textContent = "copied to clipboard";
    toast.classList.add('show');

    // Hide the toast after 2 seconds
    setTimeout(() => {
        toast.classList.remove('show');
    }, 1000);
    // document.body.removeChild(el);
}

function showErrorMessage(header, text) {
    const toast = document.getElementById('errorMessage');
    toast.textContent = header + ": " + text;
    toast.classList.add('show');

    // Hide the toast after 2 seconds
    setTimeout(() => {
        toast.classList.remove('show');
    }, 1000);
    // document.body.removeChild(el);
}

function showInfoMessage(header, text) {
    const toast = document.getElementById('infoMessage');
    toast.textContent = header + ": " + text;
    toast.classList.add('show');

    // Hide the toast after 2 seconds
    setTimeout(() => {
        toast.classList.remove('show');
    }, 1000);
    // document.body.removeChild(el);
}




function copyToClipboardUsingElement(el) {
    text = document.getElementById(el).innerHTML;
    navigator.clipboard.writeText(text);
}

function getRandomColor() {
    const letters = '0123456789ABCDEF';
    let color = '#';
    for (let i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
}

function downloadElementIntoYAML(id, name) {
    // Get the text content you want to download
    const text = document.getElementById(id).innerText;
    // Convert text to a Blob object
    const blob = new Blob([text], { type: "text/yaml" });
    // Create a link element for download
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${name}.yaml`; // Set the filename for the download

    // Trigger the download
    link.click();

    // Clean up by revoking the object URL
    URL.revokeObjectURL(link.href);
    document.body.removeChild(link);
}

function dataToYAML(data, name) {
    // Convert text to a Blob object
    const blob = new Blob([data], { type: "text/yaml" });

    // Create a link element for download
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${name}.yaml`; // Set the filename for the download

    // Trigger the download
    link.click();

    // Clean up by revoking the object URL
    URL.revokeObjectURL(link.href);
    document.body.removeChild(link);
}

function dataToJSON(data, name) {
    const jsonString = JSON.stringify(data, null, 2);
    // Convert text to a Blob object
    const blob = new Blob([jsonString], { type: "application/json" });

    // Create a link element for download
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${name}.json`; // Set the filename for the download

    // Trigger the download
    link.click();

    // Clean up by revoking the object URL
    URL.revokeObjectURL(link.href);
    document.body.removeChild(link);
}

function dataToCSV(data, name) {
    // Convert text to a Blob object
    if (name == null) {
        name = Date.now();
    }
    const content = convertToCSV(data);
    const blob = new Blob([content], { type: "text/csv" });

    // Create a link element for download
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${name}.csv`; // Set the filename for the download

    // Trigger the download
    link.click();

    // Clean up by revoking the object URL
    URL.revokeObjectURL(link.href);
    document.removeChild(link);
}

function convertToCSV(data) {
    const csvRows = [];
    const headers = Object.keys(data[0]);
    csvRows.push(headers.join(","));

    for (const row of data) {
        const values = headers.map(header => {
            const escaped = ("" + row[header]).replace(/"/g, '\\"');
            return `"${escaped}"`;
        });
        csvRows.push(values.join(","));
    }

    return csvRows.join("\n");
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

// Function to show the custom alert
function showAlert(type, title, message) {
    // Remove existing alert if it exists
    const existingAlert = document.getElementById("custom-alert");
    if (existingAlert) {
        existingAlert.remove();
    }

    // Create the alert container
    const alertContainer = document.createElement("div");
    alertContainer.id = "custom-alert";
    alertContainer.className = "fixed inset-0 z-50 flex items-center justify-center bg-zinc-600/90";

    // Create the alert box
    const alertBox = document.createElement("div");
    alertBox.id = "alert-box";
    alertBox.className = "bg-white w-auto p-6 rounded-lg shadow-lg border-l-4";

    // Set the border color based on the type
    switch (type) {
        case AlertSuccess:
            alertBox.classList.add("border-green-500");
            break;
        case AlertError:
            alertBox.classList.add("border-red-500");
            break;
        case AlertInfo:
        default:
            alertBox.classList.add("border-blue-500");
            break;
    }

    // Create the alert header with title
    const header = document.createElement("div");
    header.className = "flex justify-between items-center mb-4";

    const alertTitle = document.createElement("h2");
    alertTitle.id = "alert-title";
    alertTitle.className = "text-xl font-semibold";
    alertTitle.textContent = title;

    // Close button
    const closeButton = document.createElement("button");
    closeButton.className = "text-gray-500 hover:text-gray-700";
    closeButton.onclick = closeAlert;

    const closeIcon = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    closeIcon.className = "h-6 w-6";
    closeIcon.setAttribute("fill", "none");
    closeIcon.setAttribute("viewBox", "0 0 24 24");
    closeIcon.setAttribute("stroke", "currentColor");

    const closePath = document.createElementNS("http://www.w3.org/2000/svg", "path");
    closePath.setAttribute("stroke-linecap", "round");
    closePath.setAttribute("stroke-linejoin", "round");
    closePath.setAttribute("stroke-width", "2");
    closePath.setAttribute("d", "M6 18L18 6M6 6l12 12");
    closeIcon.appendChild(closePath);

    closeButton.appendChild(closeIcon);

    // Append title and close button to the header
    header.appendChild(alertTitle);
    header.appendChild(closeButton);

    // Create the message paragraph
    const alertMessage = document.createElement("p");
    alertMessage.id = "alert-message";
    alertMessage.className = "text-gray-700 mb-4";
    alertMessage.textContent = message;

    // Create action button
    const actionButton = document.createElement("button");
    actionButton.className = "px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 cursor-pointer";
    actionButton.textContent = "OK";
    actionButton.onclick = closeAlert;

    // Append elements to alert box
    alertBox.appendChild(header);
    alertBox.appendChild(alertMessage);
    alertBox.appendChild(actionButton);

    // Append alert box to container
    alertContainer.appendChild(alertBox);

    // Append the alert container to the body
    document.body.appendChild(alertContainer);
}

// Function to close the custom alert
function closeAlert() {
    const alertBox = document.getElementById("custom-alert");
    if (alertBox) {
        alertBox.remove();
    }
}

function ShowConfirm(title, message, onConfirm) {
    // Remove existing confirm modal if it exists
    const existingConfirm = document.getElementById("custom-confirm");
    if (existingConfirm) {
        existingConfirm.remove();
    }

    // Create the confirm modal container
    const confirmContainer = document.createElement("div");
    confirmContainer.id = "custom-confirm";
    confirmContainer.className = "fixed inset-0 z-50 flex items-center justify-center bg-zinc-600/90";

    // Create the confirm box
    const confirmBox = document.createElement("div");
    confirmBox.className = "bg-white w-80 p-6 rounded-lg shadow-lg border-l-4 border-blue-500";

    // Create the confirm header with title
    const header = document.createElement("div");
    header.className = "flex justify-between items-center mb-4";

    const confirmTitle = document.createElement("h2");
    confirmTitle.className = "text-xl font-semibold";
    confirmTitle.textContent = title;

    // Close button (optional)
    const closeButton = document.createElement("button");
    closeButton.className = "text-gray-500 hover:text-gray-700";
    closeButton.onclick = () => confirmContainer.remove();

    const closeIcon = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    closeIcon.className = "h-6 w-6";
    closeIcon.setAttribute("fill", "none");
    closeIcon.setAttribute("viewBox", "0 0 24 24");
    closeIcon.setAttribute("stroke", "currentColor");

    const closePath = document.createElementNS("http://www.w3.org/2000/svg", "path");
    closePath.setAttribute("stroke-linecap", "round");
    closePath.setAttribute("stroke-linejoin", "round");
    closePath.setAttribute("stroke-width", "2");
    closePath.setAttribute("d", "M6 18L18 6M6 6l12 12");
    closeIcon.appendChild(closePath);

    closeButton.appendChild(closeIcon);

    // Append title and close button to the header
    header.appendChild(confirmTitle);
    header.appendChild(closeButton);

    // Create the message paragraph
    const confirmMessage = document.createElement("p");
    confirmMessage.className = "text-gray-700 mb-4";
    confirmMessage.textContent = message;

    // Create Yes button
    const yesButton = document.createElement("button");
    yesButton.className = "px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 mr-2 cursor-pointer";
    yesButton.textContent = "Yes";
    yesButton.onclick = () => {
        onConfirm();
        confirmContainer.remove();
    };

    // Create No button
    const noButton = document.createElement("button");
    noButton.className = "px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 cursor-pointer";
    noButton.textContent = "No";
    noButton.onclick = () => confirmContainer.remove();

    // Append elements to confirm box
    confirmBox.appendChild(header);
    confirmBox.appendChild(confirmMessage);
    confirmBox.appendChild(yesButton);
    confirmBox.appendChild(noButton);

    // Append confirm box to container
    confirmContainer.appendChild(confirmBox);

    // Append the confirm container to the body
    document.body.appendChild(confirmContainer);
}


// Close confirm function
function closeConfirm() {
    const confirmContainer = document.getElementById("custom-confirm");
    if (confirmContainer) {
        confirmContainer.remove();
    }
}

function getRequestOptions(tokenKey, method, jsPayload) {
    const apiToken = getCookie(tokenKey);
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    myHeaders.append("Authorization", `Bearer ${apiToken}`);

    const requestOptions = {
        method: method,
        headers: myHeaders,
        redirect: "follow"
    };
    if (jsPayload != null) {
        requestOptions.body = JSON.stringify(jsPayload);
    }
    return requestOptions;
}

function toggleActionDropdownMenu() {
    const dropdown = document.getElementById("actionDropdownMenu");
    dropdown.classList.toggle("hidden");
}

// Show loading overlay
function showLoading() {
    document.getElementById("loading-overlay").classList.remove("hidden");
}

// Hide loading overlay
function hideLoading() {
    document.getElementById("loading-overlay").classList.add("hidden");
}

function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        const r = (Math.random() * 16) | 0; // Generate a random number between 0 and 15
        const v = c === 'x' ? r : (r & 0x3) | 0x8; // Generate the correct digit for 'x' or 'y'
        return v.toString(16); // Convert to hexadecimal
    });
}

function generateUUIDv4() {
    return ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}

function generateColors(count, opacity = 1) {
    const colors = [];
    for (let i = 0; i < count; i++) {
        const r = Math.floor(Math.random() * 256);
        const g = Math.floor(Math.random() * 256);
        const b = Math.floor(Math.random() * 256);
        colors.push(`rgba(${r}, ${g}, ${b}, ${opacity})`);
    }
    return colors;
}

function generateColorsWithBorder(count, bgOp, borderOp) {
    const bg = [];
    const borders = [];
    for (let i = 0; i < count; i++) {
        const r = Math.floor(Math.random() * 256);
        const g = Math.floor(Math.random() * 256);
        const b = Math.floor(Math.random() * 256);
        bg.push(`rgba(${r}, ${g}, ${b}, ${bgOp})`);
        borders.push(`rgba(${r}, ${g}, ${b}, ${borderOp})`);
    }
    return {
        backgroundColors: bg,
        borderColors: borders
    };
}

function toTitleCase(str) {
    if (!str) {
        return ""
    }
    return str.toLowerCase().replace(/\b\w/g, s => s.toUpperCase());
}

function handleInappResponseError(errorData) {
    // If the response is JSON, you can parse it (optional)
    try {
        const parsedError = JSON.parse(errorData);
        return parsedError.message;
    } catch (e) {
        return errorData;
    }
}

function setHttp(link) {
    if (link.search(/^http[s]?\:\/\//) == -1) {
        link = 'http://' + link;
    }
    return link;
}


// Toggle the popup visibility
function toggleContextPopup() {
    const popup = document.getElementById('topKPopup');
    popup.classList.toggle('hidden');
}

function connectGoogleDrive() {
    showAlert(AlertInfo, "Connect to Google Drive", "Note: This feature is not yet implemented.");
}

function connectOneDrive() {
    showAlert(AlertInfo, "Connect to One Drive", "Note: This feature is not yet implemented.");
}

function toggleURLCrawlFlag() {
    const urlCrawlDiv = document.getElementById('urlCrawlDiv');
    urlCrawlDiv.classList.toggle('hidden'); // Toggle visibility
}

function loadmultiSelect(elem) {
    return new Choices(elem, {
        removeItemButton: true, // Adds a remove button to selected items
        placeholderValue: "Select options",
        searchPlaceholderValue: "Search...",
        shouldSortItems: true, // Sorts dropdown items alphabetically
        noChoicesText: 'No Agents to choose from',
        classNames: {
            placeholder: 'choices__placeholders'
        },
        maxItemText: (maxItemCount) => {
            return `Only ${maxItemCount} values can be added`;
        },
    });
}

const MultiSelectCustomClasses = {
    placeholder: 'choices__placeholders',
    containerOuter: ['choices'],
    containerInner: ['choices__inner'],
    input: ['choices__input'],
    inputCloned: ['choices__input--cloned'],
    list: ['choices__list'],
    listItems: ['choices__list--multiple'],
    listSingle: ['choices__list--single'],
    listDropdown: ['choices__list--dropdown'],
    item: ['choices__item'],
    itemSelectable: ['choices__item--selectable'],
    itemDisabled: ['choices__item--disabled'],
    itemChoice: ['choices__item--choice'],
    description: ['choices__description'],
    group: ['choices__group'],
    groupHeading: ['choices__heading'],
    button: ['choices__button'],
    activeState: ['is-active'],
    focusState: ['is-focused'],
    openState: ['is-open'],
    disabledState: ['is-disabled'],
    highlightedState: ['is-highlighted'],
    selectedState: ['is-selected'],
    flippedState: ['is-flipped'],
    loadingState: ['is-loading'],
    notice: ['choices__notice'],
    addChoice: ['choices__item--selectable', 'add-choice'],
    noResults: ['has-no-results'],
    noChoices: ['has-no-choices'],
};

function toggleSubmenu(submenuId) {
    const submenu = document.getElementById(submenuId);
    const arrow = document.getElementById(submenuId + '-Arrow');
    if (submenu.classList.contains('hidden')) {
        submenu.classList.remove('hidden'); // Show the submenu
        arrow.classList.add('rotate-180'); // Rotate the arrow
    } else {
        submenu.classList.add('hidden'); // Hide the submenu
        arrow.classList.remove('rotate-180'); // Reset arrow rotation
    }
}

function fetchAttributeFromDiv(divId, attribute) {
    const div = document.getElementById(divId);
    if (div) {
        return div.getAttribute(attribute);
    }
    return null;
}


// Function to add a new box dynamically
function addstreamDataBox(title, content, canvas, finalSuccess) {
    const box = document.createElement('div');
    box.classList.add('box', "flex", "mb-1", "h-auto", "items-center",  "p-1", "break-words", "text-gray-700");
    let fContent = "";
    try {
        var json = JSON.parse(content);
        var contentVals = Object.values(json)
        fContent = contentVals.join(", ");
    } catch (error) {
        fContent = content;
    }

    if (finalSuccess) {
        box.classList.remove("border-l-2");
        box.innerHTML = `
            <span>
                  <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" class="h-6 w-6"
                    viewBox="0 0 24 24" version="1.1">
                    <g id="surface1">
                      <path
                        style=" stroke:none;fill-rule:nonzero;fill:rgb(36.862746%,86.666667%,37.64706%);fill-opacity:1;"
                        d="M 11.664062 23.78125 C 11.574219 23.78125 11.496094 23.722656 11.472656 23.636719 C 11.445312 23.550781 11.480469 23.457031 11.554688 23.410156 L 13.886719 21.941406 L 10.851562 23.691406 C 10.757812 23.746094 10.640625 23.714844 10.585938 23.621094 C 10.527344 23.527344 10.554688 23.410156 10.648438 23.347656 L 13.378906 21.640625 L 10.023438 23.574219 C 9.929688 23.625 9.8125 23.59375 9.757812 23.503906 C 9.699219 23.410156 9.726562 23.292969 9.816406 23.234375 L 14.355469 20.375 L 9.148438 23.378906 C 9.054688 23.429688 8.933594 23.398438 8.878906 23.304688 C 8.824219 23.214844 8.851562 23.09375 8.941406 23.035156 L 11.171875 21.664062 L 8.515625 23.195312 C 8.421875 23.246094 8.304688 23.214844 8.25 23.121094 C 8.195312 23.03125 8.222656 22.910156 8.3125 22.851562 L 8.433594 22.78125 L 8.015625 23.019531 C 7.921875 23.078125 7.800781 23.046875 7.742188 22.953125 C 7.683594 22.859375 7.714844 22.734375 7.808594 22.675781 L 15.578125 17.738281 L 7.101562 22.628906 C 7.007812 22.679688 6.890625 22.648438 6.835938 22.554688 C 6.78125 22.464844 6.808594 22.34375 6.898438 22.285156 L 7.316406 22.03125 L 6.664062 22.410156 C 6.570312 22.464844 6.449219 22.433594 6.390625 22.339844 C 6.335938 22.246094 6.363281 22.121094 6.460938 22.066406 L 12.46875 18.34375 L 6.042969 22.054688 C 5.949219 22.105469 5.832031 22.074219 5.773438 21.980469 C 5.71875 21.886719 5.75 21.769531 5.839844 21.710938 L 11.023438 18.523438 L 5.511719 21.707031 C 5.417969 21.761719 5.296875 21.730469 5.238281 21.636719 C 5.183594 21.539062 5.214844 21.417969 5.308594 21.359375 L 9.148438 19.019531 L 5.0625 21.378906 C 4.96875 21.429688 4.851562 21.398438 4.792969 21.304688 C 4.738281 21.214844 4.765625 21.09375 4.859375 21.035156 L 11.308594 17.0625 L 4.554688 20.960938 C 4.457031 21.015625 4.339844 20.984375 4.28125 20.890625 C 4.226562 20.796875 4.253906 20.675781 4.347656 20.621094 L 12.839844 15.335938 L 3.988281 20.445312 C 3.894531 20.496094 3.773438 20.464844 3.71875 20.371094 C 3.664062 20.277344 3.691406 20.160156 3.785156 20.101562 L 4.117188 19.902344 L 3.695312 20.144531 C 3.601562 20.195312 3.484375 20.160156 3.429688 20.070312 C 3.375 19.976562 3.402344 19.859375 3.492188 19.800781 L 5.703125 18.464844 L 3.382812 19.804688 C 3.289062 19.859375 3.167969 19.828125 3.109375 19.734375 C 3.054688 19.640625 3.082031 19.515625 3.179688 19.460938 L 11.691406 14.183594 L 2.925781 19.246094 C 2.832031 19.296875 2.710938 19.265625 2.65625 19.175781 C 2.601562 19.082031 2.628906 18.960938 2.71875 18.902344 L 11.320312 13.578125 L 2.503906 18.667969 C 2.40625 18.71875 2.289062 18.683594 2.234375 18.59375 C 2.179688 18.5 2.207031 18.382812 2.296875 18.324219 L 7.289062 15.285156 L 2.210938 18.21875 C 2.117188 18.269531 2 18.234375 1.945312 18.144531 C 1.890625 18.050781 1.917969 17.929688 2.007812 17.875 L 4.804688 16.183594 L 1.976562 17.8125 C 1.882812 17.863281 1.765625 17.832031 1.710938 17.738281 C 1.65625 17.648438 1.683594 17.53125 1.773438 17.472656 L 7.503906 13.976562 L 1.714844 17.316406 C 1.621094 17.371094 1.5 17.339844 1.445312 17.246094 C 1.390625 17.152344 1.417969 17.027344 1.511719 16.972656 L 5.636719 14.46875 L 1.5 16.855469 C 1.40625 16.914062 1.28125 16.878906 1.226562 16.785156 C 1.171875 16.691406 1.203125 16.566406 1.296875 16.511719 L 1.359375 16.472656 L 1.34375 16.484375 C 1.25 16.535156 1.132812 16.5 1.078125 16.410156 C 1.023438 16.316406 1.050781 16.199219 1.140625 16.140625 L 2.445312 15.351562 L 1.1875 16.078125 C 1.09375 16.132812 0.972656 16.101562 0.917969 16.007812 C 0.863281 15.914062 0.890625 15.792969 0.984375 15.738281 L 8.382812 11.175781 L 0.984375 15.445312 C 0.890625 15.5 0.773438 15.46875 0.714844 15.375 C 0.660156 15.28125 0.6875 15.164062 0.777344 15.105469 L 8.929688 10.050781 L 0.800781 14.742188 C 0.707031 14.789062 0.59375 14.753906 0.539062 14.664062 C 0.488281 14.574219 0.511719 14.460938 0.597656 14.398438 L 4.855469 11.804688 L 0.695312 14.207031 C 0.597656 14.257812 0.484375 14.226562 0.429688 14.132812 C 0.371094 14.042969 0.398438 13.925781 0.488281 13.867188 L 8.15625 9.105469 L 0.585938 13.476562 C 0.492188 13.527344 0.371094 13.496094 0.316406 13.402344 C 0.261719 13.308594 0.289062 13.191406 0.382812 13.132812 L 6.800781 9.167969 L 0.527344 12.789062 C 0.433594 12.839844 0.316406 12.804688 0.261719 12.714844 C 0.207031 12.621094 0.234375 12.5 0.324219 12.445312 L 2.675781 11.015625 L 0.515625 12.261719 C 0.421875 12.308594 0.308594 12.277344 0.253906 12.1875 C 0.199219 12.097656 0.222656 11.980469 0.308594 11.921875 L 3.578125 9.921875 L 0.523438 11.683594 C 0.429688 11.738281 0.308594 11.707031 0.253906 11.613281 C 0.195312 11.519531 0.226562 11.398438 0.320312 11.339844 L 1.09375 10.871094 L 0.554688 11.179688 C 0.460938 11.226562 0.347656 11.195312 0.296875 11.105469 C 0.242188 11.015625 0.265625 10.902344 0.351562 10.839844 L 4.765625 8.109375 L 0.632812 10.496094 C 0.539062 10.546875 0.421875 10.515625 0.367188 10.421875 C 0.3125 10.332031 0.339844 10.214844 0.425781 10.152344 L 7.234375 5.832031 L 0.8125 9.539062 C 0.71875 9.589844 0.601562 9.558594 0.546875 9.46875 C 0.488281 9.375 0.515625 9.253906 0.605469 9.195312 L 6.96875 5.117188 L 1.109375 8.5 C 1.015625 8.554688 0.894531 8.523438 0.839844 8.429688 C 0.78125 8.335938 0.808594 8.21875 0.902344 8.160156 L 4.882812 5.636719 L 1.4375 7.628906 C 1.34375 7.683594 1.222656 7.652344 1.167969 7.558594 C 1.109375 7.46875 1.140625 7.347656 1.230469 7.289062 L 4.070312 5.488281 L 1.847656 6.773438 C 1.753906 6.828125 1.632812 6.796875 1.578125 6.703125 C 1.523438 6.609375 1.546875 6.492188 1.640625 6.429688 L 4.390625 4.65625 L 2.429688 5.789062 C 2.335938 5.839844 2.21875 5.808594 2.164062 5.71875 C 2.105469 5.628906 2.128906 5.511719 2.21875 5.449219 L 5.261719 3.347656 L 3.582031 4.316406 C 3.492188 4.367188 3.371094 4.339844 3.316406 4.25 C 3.257812 4.160156 3.28125 4.039062 3.371094 3.976562 L 4.824219 2.964844 C 4.71875 2.8125 4.75 2.691406 4.847656 2.632812 L 7.484375 1.113281 C 7.578125 1.058594 7.695312 1.089844 7.753906 1.179688 C 7.808594 1.269531 7.785156 1.386719 7.699219 1.449219 L 7.453125 1.621094 L 9.355469 0.523438 C 9.449219 0.472656 9.5625 0.503906 9.621094 0.589844 C 9.675781 0.679688 9.65625 0.796875 9.570312 0.859375 L 6.523438 2.964844 L 11.203125 0.261719 C 11.296875 0.207031 11.417969 0.238281 11.472656 0.332031 C 11.527344 0.421875 11.5 0.542969 11.410156 0.601562 L 8.660156 2.375 L 12.347656 0.25 C 12.441406 0.191406 12.566406 0.222656 12.621094 0.316406 C 12.679688 0.410156 12.648438 0.535156 12.554688 0.589844 L 9.714844 2.390625 L 13.292969 0.324219 C 13.386719 0.273438 13.503906 0.304688 13.558594 0.394531 C 13.613281 0.488281 13.585938 0.605469 13.5 0.667969 L 9.515625 3.183594 L 14.210938 0.472656 C 14.304688 0.421875 14.425781 0.453125 14.480469 0.542969 C 14.535156 0.636719 14.507812 0.753906 14.421875 0.8125 L 8.054688 4.894531 L 15.257812 0.734375 C 15.351562 0.683594 15.46875 0.714844 15.527344 0.804688 C 15.582031 0.898438 15.554688 1.015625 15.464844 1.074219 L 8.65625 5.402344 L 16.179688 1.058594 C 16.273438 1.003906 16.394531 1.035156 16.449219 1.128906 C 16.507812 1.222656 16.476562 1.34375 16.386719 1.402344 L 11.976562 4.125 L 16.816406 1.332031 C 16.910156 1.277344 17.03125 1.308594 17.085938 1.402344 C 17.140625 1.496094 17.113281 1.617188 17.019531 1.675781 L 16.246094 2.144531 L 17.265625 1.558594 C 17.359375 1.503906 17.480469 1.535156 17.535156 1.628906 C 17.589844 1.722656 17.5625 1.84375 17.46875 1.902344 L 14.191406 3.90625 L 17.773438 1.835938 C 17.867188 1.789062 17.984375 1.820312 18.039062 1.914062 C 18.089844 2.003906 18.066406 2.121094 17.976562 2.179688 L 15.640625 3.601562 L 18.21875 2.113281 C 18.3125 2.058594 18.4375 2.089844 18.492188 2.183594 C 18.550781 2.277344 18.519531 2.402344 18.421875 2.457031 L 12 6.425781 L 18.785156 2.507812 C 18.878906 2.457031 18.996094 2.488281 19.054688 2.578125 C 19.109375 2.671875 19.082031 2.792969 18.992188 2.851562 L 11.320312 7.613281 L 19.367188 2.96875 C 19.460938 2.917969 19.578125 2.949219 19.632812 3.042969 C 19.6875 3.132812 19.660156 3.253906 19.570312 3.3125 L 15.324219 5.898438 L 19.773438 3.328125 C 19.867188 3.273438 19.988281 3.304688 20.042969 3.398438 C 20.101562 3.492188 20.070312 3.613281 19.980469 3.671875 L 11.828125 8.726562 L 20.289062 3.839844 C 20.386719 3.789062 20.503906 3.820312 20.558594 3.914062 C 20.613281 4.003906 20.585938 4.125 20.496094 4.183594 L 13.089844 8.746094 L 20.734375 4.332031 C 20.828125 4.277344 20.949219 4.3125 21.003906 4.40625 C 21.058594 4.5 21.03125 4.617188 20.9375 4.675781 L 19.628906 5.464844 L 21.007812 4.667969 C 21.105469 4.617188 21.222656 4.648438 21.277344 4.742188 C 21.332031 4.835938 21.304688 4.957031 21.210938 5.011719 L 21.15625 5.046875 L 21.253906 4.988281 C 21.347656 4.9375 21.46875 4.96875 21.527344 5.0625 C 21.582031 5.15625 21.550781 5.277344 21.457031 5.335938 L 17.332031 7.839844 L 21.546875 5.40625 C 21.640625 5.355469 21.757812 5.386719 21.8125 5.480469 C 21.871094 5.574219 21.839844 5.691406 21.75 5.75 L 16.027344 9.242188 L 21.847656 5.882812 C 21.941406 5.832031 22.058594 5.863281 22.113281 5.957031 C 22.167969 6.050781 22.140625 6.167969 22.050781 6.226562 L 19.246094 7.921875 L 22.078125 6.285156 C 22.175781 6.234375 22.292969 6.269531 22.347656 6.359375 C 22.402344 6.453125 22.375 6.570312 22.285156 6.628906 L 17.289062 9.667969 L 22.324219 6.761719 C 22.417969 6.707031 22.539062 6.738281 22.59375 6.832031 C 22.648438 6.925781 22.621094 7.046875 22.53125 7.101562 L 13.933594 12.429688 L 22.617188 7.414062 C 22.710938 7.363281 22.828125 7.394531 22.882812 7.488281 C 22.941406 7.582031 22.914062 7.699219 22.820312 7.757812 L 14.308594 13.03125 L 22.871094 8.089844 C 22.964844 8.039062 23.082031 8.074219 23.136719 8.164062 C 23.191406 8.257812 23.164062 8.378906 23.074219 8.4375 L 20.859375 9.773438 L 23.007812 8.53125 C 23.105469 8.480469 23.222656 8.511719 23.277344 8.605469 C 23.335938 8.699219 23.304688 8.820312 23.210938 8.875 L 22.878906 9.074219 L 23.125 8.933594 C 23.21875 8.882812 23.335938 8.914062 23.390625 9.007812 C 23.449219 9.101562 23.421875 9.21875 23.328125 9.277344 L 14.835938 14.5625 L 23.289062 9.683594 C 23.382812 9.632812 23.5 9.667969 23.554688 9.757812 C 23.609375 9.851562 23.582031 9.96875 23.492188 10.027344 L 17.035156 14.007812 L 23.390625 10.335938 C 23.484375 10.285156 23.601562 10.320312 23.65625 10.410156 C 23.710938 10.503906 23.683594 10.621094 23.59375 10.679688 L 19.765625 13.015625 L 23.445312 10.890625 C 23.539062 10.839844 23.65625 10.871094 23.710938 10.964844 C 23.769531 11.054688 23.738281 11.175781 23.652344 11.234375 L 18.460938 14.421875 L 23.480469 11.523438 C 23.574219 11.472656 23.691406 11.503906 23.746094 11.597656 C 23.804688 11.691406 23.777344 11.808594 23.6875 11.867188 L 17.675781 15.589844 L 23.480469 12.238281 C 23.574219 12.1875 23.691406 12.222656 23.746094 12.3125 C 23.800781 12.40625 23.773438 12.527344 23.683594 12.585938 L 23.273438 12.832031 L 23.453125 12.730469 C 23.546875 12.671875 23.667969 12.703125 23.726562 12.796875 C 23.785156 12.890625 23.753906 13.015625 23.660156 13.070312 L 15.890625 18.011719 L 23.332031 13.714844 C 23.425781 13.664062 23.542969 13.699219 23.597656 13.789062 C 23.652344 13.882812 23.625 14 23.535156 14.058594 L 23.269531 14.21875 C 23.367188 14.1875 23.472656 14.230469 23.515625 14.324219 C 23.558594 14.417969 23.527344 14.527344 23.4375 14.578125 L 21.207031 15.953125 L 23.082031 14.871094 C 23.175781 14.820312 23.292969 14.851562 23.351562 14.945312 C 23.40625 15.035156 23.378906 15.15625 23.289062 15.214844 L 18.753906 18.074219 L 22.8125 15.730469 C 22.90625 15.671875 23.027344 15.703125 23.085938 15.800781 C 23.144531 15.894531 23.113281 16.015625 23.019531 16.074219 L 20.28125 17.785156 L 22.5 16.507812 C 22.59375 16.457031 22.714844 16.488281 22.769531 16.578125 C 22.824219 16.671875 22.796875 16.789062 22.707031 16.847656 L 20.378906 18.320312 L 22.105469 17.320312 C 22.199219 17.269531 22.316406 17.300781 22.371094 17.390625 C 22.425781 17.484375 22.402344 17.601562 22.3125 17.660156 L 19.929688 19.191406 L 21.53125 18.269531 C 21.625 18.21875 21.742188 18.25 21.796875 18.339844 C 21.851562 18.429688 21.828125 18.546875 21.738281 18.609375 L 20.855469 19.179688 C 20.949219 19.128906 21.066406 19.15625 21.125 19.25 C 21.179688 19.339844 21.15625 19.460938 21.066406 19.519531 L 19.820312 20.363281 C 19.910156 20.371094 19.980469 20.433594 20 20.519531 C 20.019531 20.605469 19.980469 20.695312 19.90625 20.738281 L 15.605469 23.21875 C 15.511719 23.273438 15.390625 23.246094 15.335938 23.152344 C 15.277344 23.0625 15.304688 22.941406 15.394531 22.882812 L 16.238281 22.308594 L 13.996094 23.605469 C 13.902344 23.65625 13.78125 23.625 13.726562 23.535156 C 13.671875 23.441406 13.699219 23.324219 13.785156 23.261719 L 14.671875 22.691406 L 12.871094 23.734375 C 12.777344 23.785156 12.65625 23.753906 12.601562 23.664062 C 12.542969 23.570312 12.570312 23.449219 12.664062 23.390625 L 15.046875 21.859375 L 11.765625 23.753906 C 11.734375 23.773438 11.699219 23.78125 11.664062 23.78125 Z M 11.664062 23.78125 " />
                      <path style=" stroke:none;fill-rule:nonzero;fill:rgb(100%,100%,100%);fill-opacity:1;"
                        d="M 10.796875 17.375 C 10.714844 17.375 10.640625 17.324219 10.609375 17.25 C 10.582031 17.171875 10.597656 17.085938 10.660156 17.03125 L 10.863281 16.835938 L 10.167969 17.234375 C 10.078125 17.285156 9.964844 17.261719 9.90625 17.175781 C 9.847656 17.089844 9.863281 16.972656 9.941406 16.90625 L 11.09375 15.949219 L 9.621094 16.800781 C 9.527344 16.855469 9.410156 16.824219 9.351562 16.734375 C 9.292969 16.644531 9.316406 16.527344 9.402344 16.464844 L 9.949219 16.078125 L 9.277344 16.464844 C 9.183594 16.519531 9.066406 16.488281 9.007812 16.398438 C 8.953125 16.308594 8.976562 16.191406 9.0625 16.128906 L 10.476562 15.136719 L 8.867188 16.066406 C 8.773438 16.117188 8.65625 16.089844 8.597656 16 C 8.542969 15.910156 8.5625 15.792969 8.652344 15.730469 L 11.601562 13.636719 L 8.316406 15.53125 C 8.234375 15.578125 8.125 15.558594 8.0625 15.480469 C 8 15.40625 8.003906 15.292969 8.070312 15.222656 L 8.886719 14.335938 L 7.761719 14.984375 C 7.671875 15.035156 7.558594 15.011719 7.5 14.929688 C 7.4375 14.847656 7.449219 14.734375 7.523438 14.664062 L 8.058594 14.167969 L 7.347656 14.582031 C 7.261719 14.628906 7.152344 14.605469 7.09375 14.527344 C 7.03125 14.449219 7.039062 14.335938 7.105469 14.265625 L 7.8125 13.558594 L 6.9375 14.0625 C 6.851562 14.109375 6.742188 14.089844 6.679688 14.011719 C 6.617188 13.9375 6.621094 13.824219 6.6875 13.753906 L 7.558594 12.820312 C 7.5 12.78125 7.460938 12.714844 7.460938 12.644531 C 7.464844 12.574219 7.5 12.511719 7.5625 12.472656 L 8.269531 12.066406 C 8.355469 12.015625 8.464844 12.039062 8.527344 12.113281 C 8.585938 12.191406 8.582031 12.300781 8.515625 12.375 L 7.890625 13.050781 L 8.910156 12.460938 C 8.996094 12.410156 9.109375 12.429688 9.167969 12.511719 C 9.230469 12.589844 9.222656 12.703125 9.152344 12.773438 L 8.445312 13.484375 L 9.398438 12.933594 C 9.488281 12.882812 9.601562 12.90625 9.660156 12.988281 C 9.722656 13.070312 9.710938 13.183594 9.636719 13.253906 L 9.101562 13.75 L 9.816406 13.335938 C 9.902344 13.289062 10.011719 13.308594 10.074219 13.386719 C 10.132812 13.460938 10.128906 13.574219 10.0625 13.644531 L 9.25 14.53125 L 15.074219 11.167969 C 15.164062 11.117188 15.28125 11.144531 15.339844 11.234375 C 15.398438 11.324219 15.375 11.441406 15.289062 11.503906 L 12.339844 13.601562 L 14.167969 12.542969 C 14.261719 12.492188 14.378906 12.519531 14.433594 12.609375 C 14.492188 12.699219 14.46875 12.820312 14.382812 12.878906 L 12.972656 13.875 L 13.492188 13.574219 C 13.582031 13.519531 13.699219 13.550781 13.757812 13.640625 C 13.8125 13.726562 13.792969 13.847656 13.707031 13.910156 L 13.007812 14.40625 C 13.078125 14.398438 13.148438 14.429688 13.191406 14.488281 C 13.25 14.574219 13.234375 14.691406 13.152344 14.757812 L 12.003906 15.714844 L 12.128906 15.640625 C 12.21875 15.585938 12.332031 15.609375 12.394531 15.691406 C 12.457031 15.777344 12.445312 15.890625 12.367188 15.960938 L 11.777344 16.519531 C 11.882812 16.675781 11.851562 16.796875 11.753906 16.851562 L 10.894531 17.347656 C 10.867188 17.367188 10.832031 17.375 10.796875 17.375 Z M 10.796875 17.375 " />
                      <path style=" stroke:none;fill-rule:nonzero;fill:rgb(100%,100%,100%);fill-opacity:1;"
                        d="M 11.023438 13.277344 C 10.9375 13.277344 10.859375 13.21875 10.832031 13.136719 C 10.808594 13.050781 10.839844 12.960938 10.914062 12.914062 L 11.914062 12.230469 L 11.726562 12.339844 C 11.632812 12.394531 11.511719 12.367188 11.457031 12.273438 C 11.398438 12.183594 11.421875 12.0625 11.511719 12.003906 L 12.898438 11.046875 L 12.378906 11.34375 C 12.285156 11.394531 12.171875 11.367188 12.113281 11.277344 C 12.058594 11.1875 12.082031 11.070312 12.164062 11.007812 L 12.9375 10.476562 C 12.84375 10.515625 12.734375 10.480469 12.683594 10.390625 C 12.632812 10.300781 12.660156 10.1875 12.742188 10.128906 L 13.480469 9.597656 C 13.402344 9.609375 13.328125 9.578125 13.28125 9.519531 C 13.21875 9.433594 13.234375 9.316406 13.316406 9.25 L 14.308594 8.410156 C 14.222656 8.457031 14.113281 8.4375 14.054688 8.359375 C 13.992188 8.28125 13.996094 8.171875 14.0625 8.097656 L 14.789062 7.3125 C 14.691406 7.164062 14.722656 7.042969 14.816406 6.988281 L 15.417969 6.640625 C 15.503906 6.59375 15.613281 6.613281 15.671875 6.691406 C 15.734375 6.765625 15.730469 6.878906 15.664062 6.949219 L 15.242188 7.410156 L 16.050781 6.941406 C 16.140625 6.890625 16.253906 6.917969 16.3125 7 C 16.371094 7.085938 16.359375 7.199219 16.28125 7.269531 L 15.289062 8.109375 L 16.652344 7.320312 C 16.746094 7.265625 16.863281 7.292969 16.921875 7.382812 C 16.980469 7.472656 16.957031 7.59375 16.871094 7.65625 L 16.226562 8.117188 L 17.023438 7.65625 C 17.117188 7.605469 17.230469 7.636719 17.289062 7.726562 C 17.34375 7.816406 17.320312 7.933594 17.238281 7.996094 L 16.476562 8.519531 L 17.015625 8.207031 C 17.105469 8.160156 17.222656 8.1875 17.277344 8.277344 C 17.335938 8.367188 17.3125 8.484375 17.226562 8.546875 L 15.84375 9.5 L 16.371094 9.195312 C 16.464844 9.144531 16.582031 9.171875 16.636719 9.265625 C 16.695312 9.355469 16.671875 9.472656 16.582031 9.535156 L 15.582031 10.21875 L 15.773438 10.105469 C 15.867188 10.058594 15.984375 10.089844 16.039062 10.183594 C 16.09375 10.277344 16.0625 10.394531 15.972656 10.453125 L 11.125 13.25 C 11.09375 13.269531 11.058594 13.277344 11.023438 13.277344 Z M 11.023438 13.277344 " />
                      <path
                        style=" stroke:none;fill-rule:nonzero;fill:rgb(14.901961%,14.901961%,14.901961%);fill-opacity:1;"
                        d="M 16.152344 6.8125 C 14.273438 5.785156 11.734375 12.292969 10.394531 13.597656 C 9.605469 12.902344 8.664062 11.378906 7.554688 12.460938 C 7.078125 12.90625 6.113281 13.578125 6.71875 14.285156 C 7.589844 15.199219 8.550781 16.039062 9.4375 16.941406 C 12.273438 19.035156 14.699219 11.535156 16.4375 9.863281 C 17.417969 8.5625 18.09375 7.601562 16.152344 6.8125 Z M 16.828125 8.429688 C 15.0625 11.132812 13.265625 13.835938 11.492188 16.539062 C 10.011719 18.0625 8.449219 14.820312 7.253906 14.058594 C 6.648438 13.621094 8.179688 12.035156 8.613281 12.613281 C 11.277344 15.648438 10.695312 13.269531 12.8125 10.753906 C 13.628906 10.03125 14.921875 6.136719 16.054688 7.328125 C 16.449219 7.570312 17.125 7.847656 16.828125 8.429688 Z M 16.828125 8.429688 " />
                      <path
                        style=" stroke:none;fill-rule:nonzero;fill:rgb(14.901961%,14.901961%,14.901961%);fill-opacity:1;"
                        d="M 23.890625 11.964844 C 22.34375 -8.003906 -7.183594 0.230469 1.792969 18.101562 C 7.699219 28.1875 24.058594 23.664062 23.890625 11.964844 Z M 15.964844 22.734375 C -2.664062 27.988281 -4.980469 -1.703125 14.296875 0.828125 C 25.191406 3.023438 26.441406 18.9375 15.964844 22.734375 Z M 15.964844 22.734375 " />
                    </g>
                  </svg>
            </span>
            <span class="p-1 text-xs w-4/5 text-gray-800  break-words">
                ${fContent}
            </span>`;
    } else {
        if (title !== "" || title.toLowerCase() !== "info") {
        box.innerHTML = `
                <span class="p-1 text-sm w-auto font-bold text-gray-800">
                    ${title}
                </span>
                <span class="p-1 ml-2 text-xs w-auto  text-gray-800  break-words">
                    ${fContent}
                </span>`;
        } else {
            box.innerHTML = `
                <span class="p-1 text-xs w-4/5 text-gray-800  break-words">
                    ${fContent}
                </span>`;
        }
    }
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
    
}

// Function to add a new box dynamically
function addstreamCodeBox(content, canvas) {
    var fContent = "";
    try {
        var json = JSON.parse(content);
        var contentVals = Object.values(json)
        fContent = contentVals.join(", ");
    } catch (error) {
        fContent = content;
    }

    const box = document.createElement('div');
    box.classList.add('box', "flex", "mb-1", "h-auto", "items-center",  "p-1", "break-words", "text-gray-700");
    
    var code = document.createElement('code');
    code.classList.add('text-xs', 'text-gray-800', 'font-semibold', 'break-words');
    if (typeof content === 'object') {
        code.classList.add('language-json');
    }
    var pre = document.createElement('pre');
    code.innerHTML = content;
    pre.appendChild(code);
    box.appendChild(pre);
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
}

function formatUnixTime(unixTimestamp) {
    const date = new Date(unixTimestamp * 1000); // convert seconds to milliseconds
    return date.toLocaleString(); // formatted according to local timezone
  }

function addstreamTitleBox(content, canvas) {
    const box = document.createElement('div');
    box.classList.add('box', "flex", "mb-1", "h-auto", "b-4/6", "items-center",  "p-1", "break-words", "border-l-2", "text-gray-700");
    let fContent = "";
    try {
        var json = JSON.parse(content);
        var contentVals = Object.values(json)
        fContent = contentVals.join(", ");
    } catch (error) {
        fContent = content;
    }


    box.innerHTML = `
            <span class="p-1 text-md w-4/5 font-bold text-gray-800  break-words underlined">
                ${fContent}
            </span>`;
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
    
}

function addstreamMarkdownDataBox(fContent,contentId, canvas) {
    const box = document.createElement('div');
    fContent = decodeHTMLEntities(fContent);
    box.id = contentId;
    box.classList.add('box', "mb-1", "h-auto", "b-4/6","text-xs", "items-center", "p-1",  "border-l-2", "text-gray-700","w-4/5");
    fContent = renderMarkdownBase(fContent);
    const rawHTML = marked.parse(fContent);
    // Sanitize HTML content
    const cleanHTML = DOMPurify.sanitize(rawHTML);
    // box.innerHTML = `
    //         <span class="p-1 break-words">
    //             <pre>`+fContent+`</pre>
    //         </span>`;
            box.innerHTML = fContent;
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
}

function addstreamHeadingBox(content, contentId, canvas) {
    const box = document.createElement('div');
    let fContent = "";
    try {
        var json = JSON.parse(content);
        var contentVals = Object.values(json)
        fContent = contentVals.join(", ");
    } catch (error) {
        fContent = content;
    }

    box.id = contentId;
    box.classList.add('box', "mb-1", "h-auto", "b-4/6", "items-center", "break-words", "border-l-2","w-4/5");
    cont = document.createElement('div');
    cont.classList.add('p-1', 'text-md', 'font-bold', 'text-gray-900', 'break-words');
    cont.innerHTML = fContent;
    box.appendChild(cont);
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
}

function addstreamSubHeadingBox(content, contentId, canvas) {
    const box = document.getElementById(contentId);
    let fContent = "";
    try {
        var json = JSON.parse(content);
        var contentVals = Object.values(json)
        fContent = contentVals.join(", ");
    } catch (error) {
        fContent = content;
    }

    box.id = contentId;
    box.classList.add('box', "mb-1", "h-auto", "b-4/6", "items-center", "break-words", "border-l-2","w-4/5");
    cont = document.createElement('div');
    cont.classList.add('p-1', 'text-sm', 'font-bold', 'text-gray-900', 'break-words');
    cont.innerHTML = fContent;
    box.appendChild(cont);
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
}

function addstreamHeadingContentBox(content, contentId, canvas) {
    const box = document.getElementById(contentId);
    let fContent = "";
    try {
        var json = JSON.parse(content);
        var contentVals = Object.values(json)
        fContent = contentVals.join(", ");
    } catch (error) {
        fContent = content;
    }

    box.id = contentId;
    box.classList.add('box', "mb-1", "h-auto", "b-4/6", "items-center", "break-words", "border-l-2","w-4/5");
    cont = document.createElement('div');
    cont.classList.add('p-1', 'text-xs', 'font-bold', 'text-gray-900', 'break-words');
    pCont = document.createElement('p');
    pCont.innerHTML = fContent;
    cont.appendChild(pCont);
    box.appendChild(cont);
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
}



// Function to add a new box dynamically
function addstreamClickableBox(content, canvas, clickfunc) {
    let nonArray = false;
    const box = document.createElement('div');

    try {
        let fCont = JSON.parse(content);
        if (!Array.isArray(fCont)) {
            nonArray = true;
        } else {
            box.classList.add("grid", "grid-cols-1", "md:grid-cols-3", "gap-6", "mt-2");
            fCont.forEach((element) => {
                const subBox = document.createElement('span');
                subBox.classList.add('bg-orange-700', 'shadow-md', 'rounded-lg', 'shadow-orange-700', 'text-white', 'cursor-pointer', 'mb-2', 'h-auto', 'items-center', 'p-1', 'text-xs', 'break-words', 'hover:bg-pink-900',"w-4/5");
                subBox.innerText = element;
                subBox.onclick = () => {
                    clickfunc(element);
                }
                box.appendChild(subBox);
            });
        }
    } catch (error) {
        nonArray = true;
    }
    if (nonArray) {
        box.classList.add('bg-orange-700', 'shadow-md', 'rounded-lg', 'shadow-orange-700', 'text-white', 'cursor-pointer', 'flex', 'mb-2', 'h-auto', 'items-center', 'p-1', 'break-words', 'hover:bg-pink-900');
        box.innerHTML = `
        <span class="p-1 text-xs break-words">
            ${content}
        </span>`;
        box.onclick = () => {
            clickfunc(content);
        }
    }


    canvas.appendChild(box);

    canvas.scrollTop = canvas.scrollHeight;
}

function addstreamSuggestionBox(content, canvas, clickfunc) {
    let nonArray = false;
    const box = document.createElement('div');
    canvas.innerHTML = "";
    try {
        let fCont = JSON.parse(content);
        if (!Array.isArray(fCont)) {
            nonArray = true;
        } else {
            box.classList.add("flex", "gap-2","pr-6","pl-6","w-4/5");
            fCont.forEach((element) => {
                const subBox = document.createElement('div');
                const item = document.createElement('span');
                subBox.classList.add('flex','bg-blue-100','border-blue-800','z-50', 'shadow-md', 'rounded-lg', 'shadow-blue-800', 'text-gray-100', 'cursor-pointer', 'h-auto', 'items-center', 'p-1', 'text-xs', 'break-words');
                item.innerText = element;
                subBox.onclick = () => {
                    clickfunc(element);
                }
                subBox.innerHTML = suggestionSvg;
                subBox.appendChild(item);
                box.appendChild(subBox);
            });
        }
    } catch (error) {
        nonArray = true;
    }
    canvas.appendChild(box);

    canvas.scrollTop = canvas.scrollHeight;
}

function addstreamDataSuccessBox(cont, canvas) {
    let obj = {};
    try {
        obj = JSON.parse(cont);
    } catch (error) {
        console.log(error);
        obj = {
            key: "Status",
            value: cont
        }
    }
    const box = document.createElement('div');
    box.classList.add('min-w-0', "max-w-full");
    content = `
        <div class="mt-2 w-full p-1 shadow-md rounded-lg bg-white">
            <div class="text-primary-black flex min-h-16 justify-start">
              <div class="w-1 self-stretch border px-1.5 bg-green-800 border-green-800">
              </div>
              <div class="flex grow items-left justify-between overflow-hidden p-1 bg-white">
                <div class="flex w-full gap-4">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="h-6 w-6" fill="green" stroke="currentColor">
                        <path d="M9 16.2l-3.5-3.5-1.4 1.4L9 19 20 8l-1.4-1.4z"></path>
                    </svg>
                  <div class="flex w-full flex-col text-green-800">
                      <p class="text-sm font-bold items-left">
                          ${obj.key}
                      </p>
                      <p class="text-xs font-medium">
                          ${obj.value}
                      </p>
                  </div>
                </div>
              <div class="ml-2">
              </div>
            </div>
          </div>
        </div>
    `
    box.innerHTML = content;
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
}

function parseAndreturnJSONArrayElem(stringifiedData) {
    let parsedData;
    let isJsonArray = false;
    try {
        parsedData = JSON.parse(stringifiedData);
    } catch (e) {
        console.error('Invalid JSON:', e);
        return;
    }
    if (!Array.isArray(parsedData)) {
        console.error('Parsed data is not an array.');
        return;
    }

    // Check only the first item to determine if it's a JSON array
    if (parsedData.length > 0) {
        const item = parsedData[0];
        if (Array.isArray(item)) {
            isJsonArray = false;
        } else if (typeof item === 'object' && item !== null) {
            isJsonArray = true;
        } else {
            isJsonArray = false;
        }
    }
    
    return [parsedData, isJsonArray]; // Return as array to properly return multiple values
}


function addDataSetHead(canvas, title, text) {
    const box = document.createElement('div');
    box.innerHTML = `
    <div class="flex justify-start mb-1 justify-between">
    <button class="flex items-center px-4 py-2 bg-gray-600 text-gray-100 rounded-lg hover:bg-pink-900" 
    title="Copy to clipboard" onclick="copyToClipboard('${text}')">
    ${title}: ${copySvg}
  </button>
  </div>`;
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
}

function addTableFooter(reqId) {
    const table = document.getElementById(reqId);
    try {
        let tFooter = document.createElement('tfoot');
        tFooter.classList.add("bg-white");
        tFooter.innerHTML = `
        <tr>
      <td colspan="6" class="px-1 py-1">
        <div class="flex items-center justify-between">
          <!-- Left Side: Pagination Controls -->
          <div class="flex items-center space-x-2">
            <button id="${reqId}-prevPage" class="px-1 py-1 bg-black text-xs text-white rounded-lg">Previous</button>
            <button id="${reqId}-nextPage" class="px-1 py-1 bg-black text-xs text-white rounded-lg">Next</button>
            <div class="text-xs text-gray-800">
              Showing <span id="${reqId}-currentPageSize">0</span> of <span id="${reqId}-totalElements">0</span> entries
            </div>
          </div>
          <!-- Right Side: Pagination Info -->
        </div>
      </td>
    </tr>
        `;
        table.appendChild(tFooter);
    } catch (error) {
        console.error(error);
    }
}

function addTableRow(tableBody, columns, datalist, isObject) {
    try {
        datalist.forEach(data => {
            let row = document.createElement('tr');
            if (isObject) {
                columns.forEach(column => {
                    let cell = document.createElement('td');
                    try {
                        cell.innerHTML = JSON.parse(data[column])
                    } catch (error) {
                        cell.innerHTML = data[column];
                    }
                    cell.classList.add('px-1', 'py-1', 'text-xs', 'text-gray-600', 'border-b', 'border-l', 'border-r');
                    row.appendChild(cell);
                });
            }
            else {
                columns.forEach((value, index) => {
                    let cell = document.createElement('td');
                    try {
                        cell.innerHTML = JSON.parse(data[index])
                    } catch (error) {
                        cell.innerHTML = data[index];
                    }
                    cell.classList.add('px-1', 'py-1', 'text-xs', 'text-gray-600', 'border-b', 'border-l', 'border-r');
                    row.appendChild(cell);
                });
                tableBody.appendChild(row);


            }
            tableBody.appendChild(row);
        });
    } catch (error) {
        console.error(error);
    }
}

const stepSvg = `    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 m-1" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <!-- Outer Chevron -->
      <polyline points="9 18 15 12 9 6" 
                stroke="url(#gradient1)" 
                stroke-width="2" 
                stroke-linecap="round" 
                stroke-linejoin="round" 
                fill="none"/>
      
      <!-- Inner Chevron for Depth -->
      <polyline points="9 17 14 12 9 7" 
                stroke="url(#gradient2)" 
                stroke-width="2" 
                stroke-linecap="round" 
                stroke-linejoin="round" 
                fill="none"/>
      
      <!-- Decorative Dots -->
      <circle cx="5" cy="12" r="2" fill="url(#gradient3)" />
      <circle cx="19" cy="12" r="2" fill="url(#gradient3)" />
      
      <!-- Gradient Definitions -->
      <defs>
        <linearGradient id="gradient1" x1="9" y1="18" x2="15" y2="6" gradientUnits="userSpaceOnUse">
          <stop offset="0%" stop-color="white"/>
          <stop offset="100%" stop-color="white"/>
        </linearGradient>
        <linearGradient id="gradient2" x1="9" y1="17" x2="14" y2="7" gradientUnits="userSpaceOnUse">
          <stop offset="0%" stop-color="white"/>
          <stop offset="100%" stop-color="white" stop-opacity="0.5"/>
        </linearGradient>
        <radialGradient id="gradient3" cx="12" cy="12" r="10" fx="12" fy="12">
          <stop offset="0%" stop-color="white"/>
          <stop offset="100%" stop-color="white"/>
        </radialGradient>
      </defs>
    </svg>`;

function uploadFromComputer() {
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = '.csv,application/json,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'; // Accept CSV, JSON, and XLSX files
    fileInput.style.display = 'none';

    fileInput.addEventListener('change', (event) => {
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
                    document.getElementById('contextValue').value = JSON.stringify(jsonData, null, 2);
                } catch (error) {
                    // If parsing fails, treat as a plain string
                    document.getElementById('contextValue').value = fileContent;
                }
            };

            // Read the file content as text
            reader.readAsText(file);
        }
    });

    // Trigger the file input
    fileInput.click();
}

function strToUniArray(str) {
    const encoder = new TextEncoder(); // Defaults to UTF-8
    return encoder.encode(str);
}

function uint8ArrayToBase64(uint8Array) {
    let binary = '';
    uint8Array.forEach((byte) => {
        binary += String.fromCharCode(byte);
    });
    return btoa(binary);
}

function showTab(tabId) {
    const tabs = document.querySelectorAll('.tab-content');
    tabs.forEach(tab => tab.classList.add('hidden'));

    document.getElementById(tabId).classList.remove('hidden');

    const buttons = document.querySelectorAll('.tab-button');
    buttons.forEach(button => button.classList.remove('active-tab'));

    document.querySelector(`[onclick="showTab('${tabId}')"]`).classList.add('active-tab');
}

function getLoaderTextDiv(text, reqId) {
    return `<div class="chat-loader" id="loader-${reqId}">
      ${text}<span>.</span><span>.</span><span>.</span>
    </div>`;
}

function hideStreamDiv(reqId) {
    try {
        el = document.getElementById(`loader-${reqId}`);
        el.remove();
    } catch (error) {
        console.error(error);
    }

}

function addDatasetDownloadButton(tableId, canvas) {
    const box = document.createElement('div');
    box.innerHTML = `
    <button onclick="exportResultSet('${tableId}')"
    class="flex items-center mb-1 px-2 py-2 text-gray-100 text-xs rounded-lg bg-black">
    <svg viewBox="0 0 24 24" fill="black" class="h-4 w-4" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M5 20h14v2H5v-2zm7-2c-.28 0-.53-.11-.71-.29L8 13.41l1.41-1.41L11 14.17V4h2v10.17l1.59-1.59L16 13.41l-3.29 3.29c-.18.18-.43.29-.71.29z"
        fill="#FFFFFF" />
    </svg>
  </button>
    `
    canvas.appendChild(box);
    canvas.scrollTop = canvas.scrollHeight;
}

async function renderNabhikFiles(response,fileDescription, canvas, responseId) {
    if (response.format) {
        if (eligibleRenderedFileFormats.includes(response.format.toLowerCase())) {
            renderFiles(response,fileDescription, canvas, responseId);
        }
    }
}

function closeChartModal(modalId) {
    const modal = document.getElementById(modalId);
    modal.classList.add("hidden");
}

function showChartModal(contentId) {
    const myObject = document.getElementById(contentId);
    const src = myObject.getAttribute("data");
    // Set it on the modal's object
    const modalObject = document.getElementById("modalobject-" + contentId);
    modalObject.setAttribute("data", src);
    // Show the modal
    const myModal = document.getElementById("view-" + contentId);
    myModal.classList.remove("hidden");
}

function renderFiles(response,fileDescription, canvas, contentid) {
    let objectId = "file-" + contentid + "-" + generateUUIDv4();
    // let objectId = response.name;
    const box = document.createElement("div");
    const description = document.createElement("div");
    description.classList.add("text-xs", "text-gray-100", "font-semibold", "break-words");
    if (fileDescription !== "") {
        description.innerHTML = fileDescription+"<br>";
    } else if (response.description !== "") {
        description.innerHTML = response.description+"<br>";
    } else {
        description.innerHTML = response.name+"<br>";
    }
    box.appendChild(description);
    const object = document.createElement("object");
    object.id = objectId;
    let fileType = mimeTypes[response.format.toLowerCase()];
    object.type = fileType;
    const fileUrl = `data:${fileType};base64,${response.data}`;
    object.data = fileUrl;
    object.classList.add("w-[100%]", "h-[400px]","flex","justify-center","scrollbar");
    box.classList.add("mt-2", "p-1","mr-[80px]","ml-[80px]","h-[500px]", "shadow-md", "rounded-lg", "bg-[#1b1b1b]");
    box.appendChild(object);
    canvas.appendChild(box);

    const modelBox = document.createElement("div");
    modelBox.classList.add("chartmodal", "hidden");
    modelBox.style.display = "hidden";
    modelBox.id = "view-" + objectId;
    modelBox.innerHTML = `<div class="chartmodal-content">
      <span id="closemodal-${objectId}" class="chartclose" onclick="closeChartModal('view-${objectId}')">
      &times;
      </span>
      <object id="modalobject-${objectId}"
              data=""
              type="${fileType}"
              style="width: 100%; height: 100%; border: 1px solid #ccc;">
      </object>
      
    </div>`;
    box.appendChild(modelBox);

    const nBox = document.createElement("div");
    nBox.classList.add("mt-4","flex");
    const downloadLink = document.createElement("a");
    downloadLink.classList.add("items-center", "px-4", "py-2", "text-orange-700", "rounded-lg", "hover:text-pink-900", "mr-2");
    downloadLink.href = fileUrl;
    var svgFiles = document.createElement("span");
    svgFiles.classList.add("w-6", "h-6");
    svgFiles.innerHTML = `
    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
      <path fill-rule="evenodd" d="M6 2a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6H6zm7 1.5V8h4.5L13 3.5zm-3 8.75a2.25 2.25 0 1 1 0 4.5 2.25 2.25 0 0 1 0-4.5zm-3 5.25a1.75 1.75 0 1 1 0-3.5 1.75 1.75 0 0 1 0 3.5zm8-5.25a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3z" clip-rule="evenodd"/>
    </svg>
    `;
    downloadLink.appendChild(svgFiles);
    downloadLink.download = response.name;
    nBox.appendChild(downloadLink);

    const viewLink = document.createElement("button");
    viewLink.classList.add("items-center", "px-4", "py-2", "text-orange-700", "rounded-lg", "hover:text-pink-900", "mr-2", "cursor-pointer");
    viewLink.textContent = `View`;
    viewLink.id = "viewmodal-" + objectId;
    // viewLink.onclick = showChartModal(contentId);
    viewLink.addEventListener("click", function () {
        showChartModal(objectId);
    });
    nBox.appendChild(viewLink);

    canvas.appendChild(nBox);
}

function renderstreamImage(response, canvas) {
    const box = document.createElement("div");
    const imageDataUrl = `data:image/${response.format.toLowerCase()};base64,${response.data}`;
    const img = document.createElement("img");
    img.src = imageDataUrl;
    img.alt = response.name;
    img.classList.add("w-2/3", "h-96");
    box.classList.add("mt-2", "p-1", "shadow-md", "rounded-lg", "bg-white");
    box.appendChild(img);
    canvas.appendChild(box);

    const nBox = document.createElement("div");
    nBox.classList.add("mt-4");
    const downloadLink = document.createElement("a");
    downloadLink.classList.add("items-center", "px-4", "py-2", "bg-gray-600", "text-gray-100", "rounded-lg", "hover:bg-pink-900", "mr-2");
    downloadLink.href = imageDataUrl;
    var svgFiles = document.createElement("span");
    svgFiles.classList.add("w-6", "h-6");
    svgFiles.innerHTML = `
    <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
      <path fill-rule="evenodd" d="M6 2a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8l-6-6H6zm7 1.5V8h4.5L13 3.5zm-3 8.75a2.25 2.25 0 1 1 0 4.5 2.25 2.25 0 0 1 0-4.5zm-3 5.25a1.75 1.75 0 1 1 0-3.5 1.75 1.75 0 0 1 0 3.5zm8-5.25a1.5 1.5 0 1 1 0 3 1.5 1.5 0 0 1 0-3z" clip-rule="evenodd"/>
    </svg>
    `;
    downloadLink.innerHTML = svgFiles;
    downloadLink.download = response.name;
    nBox.appendChild(downloadLink);
    canvas.appendChild(nBox);
}


function isJSONorJSONArray(str) {
    try {
        const parsed = JSON.parse(str);

        if (typeof parsed === "object" && parsed !== null) {
            if (Array.isArray(parsed)) {
                // const allObjects = parsed.every(item => typeof item === "object" && item !== null);
                return [parsed, true];
            } else {
                return [parsed, false];
            }
        }
        return [undefined, false];
    } catch (e) {
        // If parsing fails, it's not valid JSON
        return [undefined, false];
    }
}

function getFileExtension(fileName) {
    // Find the last dot in the file name
    const lastDotIndex = fileName.lastIndexOf('.');
    // If there's no dot or it's the first character, return an empty string
    if (lastDotIndex === -1 || lastDotIndex === 0) {
        return '';
    }
    // Return the substring after the last dot
    return fileName.slice(lastDotIndex + 1);
}

async function fileUploader(apiUrl, accessTokenKey, resource, resourceId) {
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
                        myHeaders.append("Accept", "application/x-ndjson");
                        myHeaders.append("Content-Type", "application/x-ndjson");
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
                            showAlert(AlertError, "Upload Failed", "File upload failed, please try after some time");
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

function updateBrowserUrl(key, value) {
    let currentUrl = window.location.href;

    // Check if the URL already contains query parameters
    let newUrl = new URL(currentUrl);

    // Add a new query parameter (or update if exists)
    newUrl.searchParams.set(key, value);

    // Update the URL in the browser without reloading
    window.history.pushState({}, '', newUrl);
}

function eliminateSubStrings(input, substrings) {
    substrings.forEach(sub => {
        input = input.replace(sub, "");
    });
    return input;
}


function renderMarkdown(markdownText) {
    // Use Marked to parse the markdown into HTML
    const rawHtml = marked.parse(markdownText, {
        highlight: function (code, lang) {
            // Use Prism.js for syntax highlighting
            if (Prism.languages[lang]) {
                return Prism.highlight(code, Prism.languages[lang], lang);
            } else {
                // Return plain code if the language isn't supported
                return code;
            }
        },
        gfm: true, // Enable GitHub Flavored Markdown
        breaks: false, // Disable auto <br> conversion
        smartypants: true,
        xhtml: false,
    });
    return rawHtml;
}

function renderMarkdownBase(markdownText) {
    // Use Marked to parse the markdown into HTML without syntax highlighting
    const rawHtml = marked.parse(markdownText, {
        gfm: true, // Enable GitHub Flavored Markdown
        breaks: true, // Convert line breaks to <br> tags
    });
    return rawHtml;
}

function decodeHTMLEntities(text) {
    const textarea = document.createElement("textarea");
    textarea.innerHTML = text;
    const decodedText = textarea.value;
    textarea.remove(); // Optional, ensures cleanup
    return decodedText;
}

function formatHtmlContent(resId) {
    const container = document.getElementById(resId);
    // Find all code blocks and wrap them with headers
    container.querySelectorAll("pre > code").forEach((codeBlock) => {
        const language = codeBlock.className.replace("language-", "") || "plaintext";

        // Create a wrapper div for the code block
        const wrapper = document.createElement("div");
        wrapper.classList.add("code-wrapper", "bg-black");
        wrapper.style.position = "relative";
        wrapper.style.border = "1px solid #ccc";
        wrapper.style.borderRadius = "8px";
        wrapper.style.overflow = "hidden";
        wrapper.style.wordBreak = "break-word";

        // Create the header
        const header = document.createElement("div");
        header.className = "code-header";
        header.style.display = "flex";
        header.style.justifyContent = "space-between";
        header.style.alignItems = "center";
        header.style.backgroundColor = "#2d2d2d";
        header.style.color = "white";
        header.style.padding = "8px 16px";
        header.style.fontSize = "14px";
        header.style.fontWeight = "bold";
        header.style.marginBottom = "4";

        // Add the language label
        const langLabel = document.createElement("span");
        langLabel.textContent = language;
        header.appendChild(langLabel);

        // Add the copy button
        const copyButton = document.createElement("button");
        copyButton.textContent = "Copy";
        copyButton.style.background = "#ff5f1f";
        copyButton.style.color = "white";
        copyButton.style.border = "none";
        copyButton.style.borderRadius = "4px";
        copyButton.style.padding = "4px 8px";
        copyButton.style.cursor = "pointer";
        copyButton.addEventListener('click', function () {
            copyToClipboard(codeBlock.textContent);
        });
        codeBlock.classList.add("p-1");
        const codeLines = codeBlock.innerHTML.split("\n").map((line) => {
            return `<span class="code-line p-2">${line}</span>`;
        }).join("\n");
        codeBlock.innerHTML = codeLines;
        header.appendChild(copyButton);

        // Wrap the code block
        const preElement = codeBlock.parentElement; // <pre> element
        wrapper.appendChild(header);
        wrapper.appendChild(preElement.cloneNode(true)); // Clone the original <pre> block

        // Replace the original <pre> block with the wrapper
        preElement.parentElement.replaceChild(wrapper, preElement);
    });
    return container.innerHTML;
}

function redirectToLink(link) {
    window.location.replace(link);
}

function renderMDContent(text) {
  // Simple Markdown detection (looks for common Markdown patterns)
  const isMarkdown = /^#|^[-*] |^[0-9]+\. |\[.*\]\(.*\)|`{1,3}[^`]|^\> /.test(text);
  if (isMarkdown) {
    try {
      marked.setOptions({
        gfm: true,
        breaks: true,
        headerIds: true,  // Add IDs to headings for anchor links
        headerPrefix: '', // No prefix for header IDs
      });
      const html = marked.parse(text);
      return html
        .replace(/<h1>/g, '<h1 class="markdown-h1">')
        .replace(/<h2>/g, '<h2 class="markdown-h2">')
        .replace(/<h3>/g, '<h3 class="markdown-h3">');
    } catch (e) {
      console.error("Markdown parsing failed:", e);
      return `<div class="text-content">${escapeHtml(text)}</div>`;
    }
  }
  
  // Plain text fallback
  return `<div class="text-content">${escapeHtml(text)}</div>`;
}

// Basic HTML escaping for plain text
function escapeHtml(unsafe) {
  return unsafe
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}

function downloadFile(file) {
    const binaryData = atob(file.data);
    const byteNumbers = new Uint8Array(binaryData.length);
    for (let i = 0; i < binaryData.length; i++) {
        byteNumbers[i] = binaryData.charCodeAt(i);
    }

    // Create a Blob with the decoded data and the correct MIME type
    const blob = new Blob([byteNumbers], { type: mimeTypes[file.format] });
    // Create a URL for the Blob
    const url = URL.createObjectURL(blob);

    // Create a temporary anchor element to trigger the download
    const a = document.createElement("a");
    a.href = url;
    a.download = `${file.name}`;
    document.body.appendChild(a);

    // Trigger the download
    a.click();

    // Cleanup: Remove the anchor element and revoke the Blob URL
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}
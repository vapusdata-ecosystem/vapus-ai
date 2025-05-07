function showOverloadinglay(id) {
  document.getElementById(id).style.display = 'flex';
}

function hideOverloadinglay(id) {
  document.getElementById(id).style.display = 'none';
}


async function ManageDataWorkerDeployment(action, depId, url, tokenKey) {
  // Show the loading overlay within the specific section
  showOverloadinglay('loadingOverlay');
  try {
    // Make the API call
    var queryParams = {
      action: action,
      spec: {
        deploymentId: depId,
      }
    };
    const response = await fetch(url, getRequestOptions(tokenKey, "POST", queryParams));

    // Process the response
    if (!response.ok) {
      throw new Error('Data Worker Deployment Action failed');
    }

    const result = await response.json();
    showAlert(AlertInfo, "Data Worker Deployment", "Action completed successfully");
    location.reload(true);

  } catch (error) {
    showAlert(AlertError, "Data Worker Deployment", "Deployment action failed: ");
  } finally {
    // Hide the loading overlay after the call completes
    hideOverloadinglay('loadingOverlay');
  }
}

async function ManageVDCDeployment(action, depId, url, tokenKey) {
  // Show the loading overlay within the specific section
  showOverloadinglay('loadingOverlay');

  try {
    // Make the API call
    var queryParams = {
      action: action,
      spec: {
        deploymentId: depId,
      }
    };
    const response = await fetch(url, getRequestOptions(tokenKey, "POST", queryParams));

    // Process the response
    if (!response.ok) {
      throw new Error('VDC Deployment Action failed');
    }

    const result = await response.json();
    showAlert(AlertInfo, "VDC Deployment", "Action completed successfully");
    location.reload(true);

  } catch (error) {
    showAlert(AlertError, "VDC Deployment", "Deployment action failed: ");
  } finally {
    // Hide the loading overlay after the call completes
    hideOverloadinglay('loadingOverlay');
  }
}

async function upgradeBaseOs(apiUrl, tokenKey, action) {
  // Show the loading overlay within the specific section
  showOverloadinglay('loadingOverlay');
  console.log(action, apiUrl, tokenKey);
  try {
    // Make the API call
    var queryParams = {
      actions: action
    };

    const response = await fetch(apiUrl, getRequestOptions(tokenKey, "POST", queryParams));

    // Process the response
    if (!response.ok) {
      throw new Error('Upgrade failed');
    }

    const result = await response.json();
    showAlert(AlertSuccess, 'Organization Settings', 'Base Os Upgraded successfully');

  } catch (error) {
    console.log(error);
    showAlert(AlertError, 'Organization Settings', 'Base Os Upgrade failed');
  } finally {
    // Hide the loading overlay after the call completes
    hideOverloadinglay('loadingOverlay');
  }
}

async function manageAccessRequests(url, dpId, tokenKey, commenttextElement, reqId, statusElem) {
  showLoading();
  const comment = document.getElementById(commenttextElement).value;
  const status = document.getElementById(statusElem).value;
  const payload = {
    spec: {
      status: status,
      requestId: reqId,
      dataProductId: dpId,
      comments: [
        {
          comment: comment
        }
      ]
    }
  };
  requestOptions = getRequestOptions(tokenKey, "POST", payload);
  try {
    const response = await fetch(url, requestOptions);
    if (!response.ok) {
      showAlert(AlertError, "Request Access", "Data Product Access Request Failed");
    }

    const result = await response.json();
    hideLoading();  // Hide the loading overlay
    showAlert(AlertInfo, "Request Access", "Access Request is successfully submitted");
    location.reload(true);

  } catch (error) {
    console.log(error);
    showAlert(AlertError, "Request Access", "Data Product Access Request Failed");
  } finally {
    hideLoading();
  }
}

function getVdcEndpoint(url, uri) {
  console.log(url, uri);
  url = setHttp(url);
  console.log(url);
  parsedurl = new URL(url);
  parsedurl.pathname = uri;
  return parsedurl.toString();
}

async function checkVdcStatus(host, uri, tokenKey, statuselem) {
  document.getElementById(statuselem).innerText = "checking...";

  try {
    // url = getVdcEndpoint("127.0.0.1:9073", uri);
    url = getVdcEndpoint(host, uri);
    console.log(url);
    requestOptions = getRequestOptions(tokenKey, "GET", null);
    const response = await fetch(url, requestOptions);
    if (!response.ok) {
      showAlert(AlertError, "Check VDC Status", "Error checking VDC status");
    }

    const result = await response.json();
    document.getElementById(statuselem).parentElement.classList.remove("hidden");
    document.getElementById(statuselem).innerText = result.status;
    // hideOverloadinglay('loadingOverlay');
    showAlert(AlertInfo, "Check VDC Status", "VDC status is successfully synced");
  } catch (error) {
    document.getElementById(statuselem).innerText = "";
    showAlert(AlertError, "Check VDC Status", "Error checking VDC status");
  } finally {
    // hideOverloadinglay('loadingOverlay');
  }
}

async function manageVdc(host, uri, tokenKey, statuselem, action) {
  showLoading();
  try {
    // url = getVdcEndpoint("127.0.0.1:9073", uri);
    url = getVdcEndpoint(host, uri);
    console.log(url);
    console.log(action);
    const payload = {
      action: action,
    };
    requestOptions = getRequestOptions(tokenKey, "POST", payload);
    const response = await fetch(url, requestOptions);
    if (!response.ok) {
      showAlert(AlertError, "VDC Action", "VDC Action Failed");
    }
    const result = await response.json();
    console.log(result);
    hideLoading();  // Hide the loading overlay
    showAlert(AlertInfo, "VDC Action", "VDC Action is successfully completed");
    document.getElementById(statuselem).parentElement.classList.remove("hidden");
    document.getElementById(statuselem).innerText = result.status;

  } catch (error) {
    console.log(error);
    showAlert(AlertError, "VDC Action", "VDC Action Failed");
  } finally {
    hideLoading();
  }
}

function toggleDropdown() {
    const dropdownMenu = document.getElementById('dropdown-menu');
    dropdownMenu.classList.toggle('hidden');
}

document.addEventListener('DOMContentLoaded', function () {
    const sidebar = document.getElementById('sidebar-content');
    const buttons = sidebar.getElementsByTagName('a');

    for (let button of buttons) {
        button.addEventListener('click', function () {
            // Remove 'selected' class from all buttons
            for (let btn of buttons) {
                btn.classList.remove('selected');
            }
            // Add 'selected' class to the clicked button
            this.classList.add('selected');
        });
    }

    // Set the selected class based on GlobalContext.CurrentSideBar
    const currentSideBar = '{{ .GlobalContext.CurrentSideBar }}';
    console.log(currentSideBar);
    if (currentSideBar) {
        const currentElement = document.getElementById(currentSideBar);
        console.log(currentSideBar);
        if (currentElement) {
            console.log(currentElement);
            currentElement.classList.add('selected');
        }
    }

    // Close dropdown if clicked outside
    // document.addEventListener('click', function(event) {
    //     const dropdownMenu = document.getElementById('dropdown-menu');
    //     const menuButton = document.getElementById('menu-button');
    //     if (!menuButton.contains(event.target) && !dropdownMenu.contains(event.target)) {
    //         dropdownMenu.classList.add('hidden');
    //     }
    // });

    // // Toggle dropdown on button click
    // const menuButton = document.getElementById('menu-button');
    // menuButton.addEventListener('click', function(event) {
    //     event.stopPropagation(); // Prevent the click event from propagating to the document
    //     toggleDropdown();
    // });

    // // Handle dropdown item selection
    // const dropdownItems = document.querySelectorAll('#dropdown-menu a');
    // dropdownItems.forEach(item => {
    //     item.addEventListener('click', function() {
    //         dropdownItems.forEach(itm => itm.classList.remove('selected'));
    //         this.classList.add('selected');
    //     });
    // });
});
// document.addEventListener('DOMContentLoaded', function() {
//     const dropdownMenu = document.getElementById('dropdown-menu');
//     const dropdownItems = dropdownMenu.getElementsByTagName('a');

//     for (let item of dropdownItems) {
//         item.addEventListener('click', function() {
//             // Remove 'selected' class from all items
//             for (let itm of dropdownItems) {
//                 itm.classList.remove('selected');
//             }
//             // Add 'selected' class to the clicked item
//             this.classList.add('selected');
//         });
//     }
// });
// // Close dropdown if clicked outside
// document.addEventListener('click', function(event) {
// const dropdownMenu = document.getElementById('dropdown-menu');
// const menuButton = document.getElementById('menu-button');
// if (!menuButton.contains(event.target) && !dropdownMenu.contains(event.target)) {
//     dropdownMenu.classList.add('hidden');
// }
// });


async function requestAccess(url, dpId, tokenKey, commenttextElement) {
    console.log(url, dpId, tokenKey, commenttextElement);
    showLoading();
    const comment = document.getElementById(commenttextElement).value;
    const payload = {
        spec: {
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

async function manageAccessRequests(url, dpId, tokenKey, commenttextElement, reqId, actionElem, currentStatus) {
    console.log(url, dpId, tokenKey, commenttextElement);
    showLoading();
    const comment = document.getElementById(commenttextElement).value;
    let status = document.getElementById(actionElem).value;
    if (status == "") {
        status = currentStatus;
    }
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
            showAlert(AlertError, "Request Access", "Data Product Access Request Message Failed");
        }

        const result = await response.json();
        hideLoading();  // Hide the loading overlay
        showAlert(AlertInfo, "Request Access", "Access Request message is successfully submitted");
        location.reload(true);

    } catch (error) {
        console.log(error);
        showAlert(AlertError, "Request Access", "Data Product Access Request Message Failed");
    } finally {
        hideLoading();
    }
}
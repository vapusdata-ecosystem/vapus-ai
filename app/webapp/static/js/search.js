function handleSearchResponseOrErr(error,resultFId) {
    const resultsContainer = document.getElementById(resultFId);
    console.log(error);
    // Clear previous table rows
    const row = document.createElement('div');
    
    row.innerHTML = 'No results found.';
    row.className = 'col-span-4 sm:col-span-4 text-gray-100  m-2 text-left text-lg';
        resultsContainer.prepend(row);
    // Show the table container
    }
      

// Function to perform search and display results as cards
function discoverMarketplace(apiUrl, tokenKey,bAUrl, queryFId,resultFId) {
    console.log(apiUrl, tokenKey,queryFId,resultFId);
    const query = document.getElementById(queryFId).value;
    showLoading();
    // Filter or process search results (using mock data here)
    fetch(`${apiUrl}?q=${query}`, getRequestOptions(tokenKey, "GET",null))
    .then((response) => response.text())
    .then((result) => displaySearchResponseTable(query,result,resultFId,bAUrl))
    .catch((error) => handleSearchResponseOrErr(error,resultFId))
    .finally(() => hideLoading());
}

function exploreDatasources(apiUrl, tokenKey,bAUrl, queryFId,category,resultFId) {
    console.log(apiUrl, tokenKey,queryFId,category,resultFId);
    const query = document.getElementById(queryFId).value;
    
    // Filter or process search results (using mock data here)
    fetch(`${apiUrl}?q=${query}&search_type=${category}`, getRequestOptions(tokenKey, "GET",null))
    .then((response) => response.text())
    .then((result) => displaySearchResponseTable(query,result,resultFId,bAUrl))
    .catch((error) => handleSearchResponseOrErr(error,resultFId));
}

function displaySearchResponseTable(q,results,resultFId,bAUrl) {
    const resultsContainer = document.getElementById(resultFId);
    // resultsContainer.innerHTML = '';
    // Generate HTML for each result as a card
    const seperator = document.createElement('div');
    seperator.innerHTML = `Query - ${q}`
    seperator.className = "col-span-4 sm:col-span-4 bg-zinc-800 border-gray-700 p-2  m-2 text-left my-2 text-xl text-orange-700 rounded-lg";
    resultsContainer.prepend(seperator);

    console.log(results);
    const data = JSON.parse(results);
    const host = window.location.host;
    data.results.forEach(obj => {
        const card = document.createElement('div');
        card.className = "p-4 bg-[#1b1b1b] text-gray-100 border border-zinc-500 rounded-lg shadow-lg";
        let bUrl = "";
        if (obj.result.url === ""){
            bUrl = "#";
        } else {
        bUrl = bAUrl+obj.result.resourceId;
        }
        console.log(bUrl);
        card.innerHTML = `
            <a href="${bUrl}" class="text-gray-100 hover:underline" target="_blank">
                <h3 class="font-semibold text-lg">${obj.result.name}</h3>
                <p class="text-gray-100 mt-2 text-xs">ID: ${obj.result.resourceId}</p>
            </a>
        `;
        resultsContainer.prepend(card);
    });

    if (data.results.length === 0) {
        const noResultsContainer = document.createElement('div');
        noResultsContainer.innerHTML = 'No results found.';
        noResultsContainer.className = 'col-span-4 sm:col-span-4 text-gray-100 m-2 text-left text-lg';
        resultsContainer.prepend(noResultsContainer);
    }

}

// Async function to call the API and handle streaming data
async function performStreamSearch(apiUrl, tokenKey, queryFId, categoryFId, resultFId) {
    const query = document.getElementById(queryFId).value;
    const category = document.getElementById(categoryFId).value;
    const resultsContainer = document.getElementById(resultFId);
    const messagesContainer = document.getElementById('messages-container');

    // Clear previous results
    messagesContainer.innerHTML = '';
    resultsContainer.innerHTML = '';
    showLoading();
    try {
        
        const response = await fetch(`${apiUrl}?q=${query}&search_type=${category}`,getRequestOptions(tokenKey, "GET",null));
        
        if (!response.body) {
            throw new Error('Readable stream not supported');
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder('utf-8');
        let done = false;

        while (!done) {
            const { value, done: streamDone } = await reader.read();
            done = streamDone;
            
            if (value) {
                // Decode the streamed data chunk
                const chunk = decoder.decode(value, { stream: true });
                // Parse the JSON chunk (assuming each chunk is a complete JSON object)
                const { messages, results } = JSON.parse(chunk);

                // Display messages above the cards as they come in
                messages.forEach(message => {
                    const messageElement = document.createElement('div');
                    messageElement.textContent = message;
                    messageElement.className = 'text-gray-600 mb-2';
                    messagesContainer.appendChild(messageElement);
                });

                // Display each result as a card below the messages
                results.forEach(result => {
                    const card = document.createElement('div');
                    card.className = "p-4 bg-white border border-gray-300 rounded-lg shadow-lg";
                    card.innerHTML = `
                        <h3 class="font-semibold text-lg">${result.title}</h3>
                        <p class="text-gray-600 mt-2">${result.description}</p>
                    `;
                    resultsContainer.appendChild(card);
                });
            }
        }
    } catch (error) {
        console.error('Error fetching data:', error);
        messagesContainer.innerHTML = `<p class="text-red-500">An error occurred while fetching data.</p>`;
    } finally {
        hideLoading();
    }
}

// // Event listener for the search button
// document.getElementById('search-button').addEventListener('click', performSearch);

// // Event listener for pressing Enter in the search input
// document.getElementById('search-input').addEventListener('keydown', function(event) {
//     if (event.key === 'Enter') {
//         performSearch();
//     }
// });


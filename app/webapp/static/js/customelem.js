function setupDropdown({
  dropdownToggleSelector,
  dropdownMenuSelector,
  inputFieldSelector = null, // Optional, for updating a hidden input or text field
  parentClass = 'parent',
  childClass = 'child',
  displayAttribute,
  valueTransform, // Default transform for display text
}) {
  const dropdownToggle = document.querySelector(dropdownToggleSelector);
  const dropdownMenu = document.querySelector(dropdownMenuSelector);
  const inputField = inputFieldSelector ? document.querySelector(inputFieldSelector) : null;

  // Toggle dropdown visibility
  dropdownToggle.addEventListener('click', () => {
    dropdownMenu.classList.toggle('show');
  });

  // Close dropdown when clicking outside
  document.addEventListener('click', (event) => {
    if (!dropdownToggle.contains(event.target) && !dropdownMenu.contains(event.target)) {
      dropdownMenu.classList.remove('show');
    }
  });

  // Handle selection of dropdown items
  dropdownMenu.addEventListener('click', (event) => {
    const item = event.target;
    // Ignore clicks on parent items
    if (item.classList.contains(parentClass)) {
      return;
    }
    // Handle selection of child items
    if (item.classList.contains(childClass)) {
      let displayValue = item.getAttribute(displayAttribute);
      // Update the toggle text
      if (valueTransform !== undefined) {
        displayValue = valueTransform(displayValue);
      }
      dropdownToggle.textContent = displayValue;
      // Update the input field, if provided
      if (inputField) {
        const selectedValue = item.dataset.value;
        inputField.value = selectedValue;
      }
      // Close the dropdown
      dropdownMenu.classList.remove('show');
    }
  });
}



function addTablePagination({
tableId,
rowsPerPage,
prevPageBtn,
nextPageBtn,
currentPageSizeSpan,
totalElementsSpan,
}) {
const table = document.getElementById(tableId);
const rows = table.querySelectorAll("tbody tr");
const totalRows = rows.length;
const totalPages = Math.ceil(totalRows / rowsPerPage);
let currentPage = 1;
// Function to display rows for the current page
function displayPage(page) {
  const start = (page - 1) * rowsPerPage;
  const end = Math.min(page * rowsPerPage, totalRows);
  rows.forEach((row, index) => {
    row.style.display = index >= start && index < end ? "" : "none";
  });

  // Update pagination info
  currentPageSizeSpan.textContent = end - start;
  totalElementsSpan.textContent = totalRows;

  // Enable/Disable navigation buttons
  prevPageBtn.disabled = page === 1;
  nextPageBtn.disabled = page === totalPages;
}

// Add event listeners for navigation
prevPageBtn.addEventListener("click", () => {
  if (currentPage > 1) {
    currentPage--;
    displayPage(currentPage);
  }
});

nextPageBtn.addEventListener("click", () => {
  if (currentPage < totalPages) {
    currentPage++;
    displayPage(currentPage);
  }
});

// Initialize table display
displayPage(currentPage);
}

function renderDatatable(tableDivId, csvRows, headers, arrayRows,pageLength,yAxisLength) {
  if ($.fn.DataTable.isDataTable('#' + tableDivId)) {
    $('#' + tableDivId).DataTable().clear().destroy();
  }
  var data = [];
  if (arrayRows) {
    data = arrayRows.map(row => {
      return headers.reduce((obj, key, index) => {
        obj[key] = row[index];
        return obj;
      }, {});
    });
  } else {
    data = csvRows;
  }
  var csvColumns = headers.map(col => ({
    title: col,
    data: col
  }));
  if (!yAxisLength) {
    yAxisLength = true;
  }
  if (!pageLength) {
    pageLength = 10;
  }
  $.fn.dataTable.ext.errMode = 'none';
  const tableSet = $('#' + tableDivId).DataTable({
    data: data,
    columns: csvColumns,
    pageLength: pageLength,
    searching: true,
    ordering: true,
    order: [], 
    select: true,
    scrollX: true,    // Enable horizontal scrolling
    scrollY: yAxisLength, // Enable vertical scrolling with fixed height
    scrollCollapse: true,
    responsive: true,
    initComplete: function () {
      const api = this.api();
      $(this.api().table().node())
            .addClass('scrollbar')
      api.columns().every(function () {
        const column = this;
        // const select = $('<select><option value="">All</option></select>')
        //   .appendTo($(column.header()).empty())
        //   .on('change', function () {
        //     const val = $.fn.dataTable.util.escapeRegex($(this).val());
        //     column.search(val ? `^${val}$` : '', true, false).draw();
        //   });

        // column.data().unique().sort().each(function (d) {
        //   select.append(`<option value="${d}">${d}</option>`);
        // });
      });
    }
  });
  return tableSet;
}
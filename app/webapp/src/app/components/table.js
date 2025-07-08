"use client";

import React, { useState, useEffect, useRef } from "react";
import LoadingOverlay from "./loading/loading";

const DataTable = ({
  id,
  data = [],
  columns = [],
  pageLength = 10,
  searching = true,
  ordering = true,
  select = true,
  responsive = true,
  loading = false,
  loadingText = "Loading data...", 
  className = "",
  filteredColumns = [],
}) => {
  const tableRef = useRef(null);
  const dataTableRef = useRef(null);
  const [scriptsReady, setScriptsReady] = useState(false);

  useEffect(() => {
    const checkScriptsLoaded = () => {
      return window.jQuery && window.jQuery.fn && window.jQuery.fn.DataTable;
    };

    const loadScripts = async () => {
      if (checkScriptsLoaded()) {
        setScriptsReady(true);
        return;
      }

      // Add DataTables CSS
      if (
        !document.querySelector('link[href*="dataTables.dataTables.min.css"]')
      ) {
        const link = document.createElement("link");
        link.rel = "stylesheet";
        link.href =
          "https://cdn.datatables.net/2.2.2/css/dataTables.dataTables.min.css";
        document.head.appendChild(link);
      }

      // Add jQuery
      if (!window.jQuery) {
        const jQueryScript = document.createElement("script");
        jQueryScript.src = "https://code.jquery.com/jquery-3.7.1.min.js";
        jQueryScript.async = true;
        document.body.appendChild(jQueryScript);

        await new Promise((resolve) => {
          jQueryScript.onload = resolve;
        });
      }

      // Add DataTables JS
      if (!window.jQuery?.fn?.DataTable) {
        const dataTablesScript = document.createElement("script");
        dataTablesScript.src =
          "https://cdn.datatables.net/2.2.2/js/dataTables.min.js";
        dataTablesScript.async = true;
        document.body.appendChild(dataTablesScript);

        await new Promise((resolve) => {
          dataTablesScript.onload = resolve;
        });
      }

      // Ensure everything is loaded
      let attempts = 0;
      const maxAttempts = 10;

      while (!checkScriptsLoaded() && attempts < maxAttempts) {
        await new Promise((resolve) => setTimeout(resolve, 200));
        attempts++;
      }

      if (checkScriptsLoaded()) {
        setScriptsReady(true);
      } else {
        console.error(
          "Failed to load DataTables scripts after multiple attempts"
        );
      }
    };

    loadScripts();

    return () => {
      if (dataTableRef.current) {
        dataTableRef.current.destroy();
        dataTableRef.current = null;
      }
    };
  }, []);

  // Add custom CSS for header filters
  useEffect(() => {
    if (!scriptsReady) return;
    const style = document.createElement("style");
    style.textContent = `
    .dt-header-filter {
      position: relative;
      cursor: pointer;
    }
    .dt-header-filter:after {
      content: '▼';
      font-size: 10px;
      margin-left: 5px;
      opacity: 0.6;
    }
    .dt-filter-dropdown {
      position: absolute;
      top: 100%;
      left: 0;
      z-index: 1000;
      display: none;
      min-width: 160px;
      max-height: 200px;
      overflow-y: auto;
  
      background-color: #333;
      border: 1px solid #444;
      border-radius: 4px;
      box-shadow: 0 4px 8px rgba(0,0,0,0.3);
    }
    .dt-filter-dropdown.show {
      display: block;
    }
    .dt-filter-option {
      padding: 6px 12px;
      cursor: pointer;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .dt-filter-option:hover {
      background-color: #444;
    }
    .dt-filter-option.active {
      background-color: #555;
    }
    .dt-filter-option.all {
      border-bottom: 1px solid #444;
      font-weight: bold;
    }
  
    /* Custom Scrollbar */
    .dt-filter-dropdown::-webkit-scrollbar {
      width: 6px;
      height: 8px;
      cursor: pointer;
    }
  
    .dt-filter-dropdown::-webkit-scrollbar-track {
      border-radius: 10vh;
      background: oklch(0.37 0.013 285.805);
    }
  
    .dt-filter-dropdown::-webkit-scrollbar-thumb {
      background: #dd290a;
      border-radius: 10vh;
    }
  
    .dt-filter-dropdown::-webkit-scrollbar-thumb:hover {
      background: #dd290a;
    }
  `;

    document.head.appendChild(style);

    return () => {
      document.head.removeChild(style);
    };
  }, [scriptsReady]);

  useEffect(() => {
    if (!scriptsReady || !tableRef.current || loading) return;

    const timer = setTimeout(() => {
      const $ = window.jQuery;
      if (dataTableRef.current) {
        dataTableRef.current.destroy();
        dataTableRef.current = null;
      }

      // Format columns
      const dtColumns = columns.map((col) => {
        if (typeof col === "string") {
          return { title: col, data: col };
        }
        return col;
      });

      try {
        // Initialize DataTable
        dataTableRef.current = $(tableRef.current).DataTable({
          data,
          columns: dtColumns,
          pageLength,
          searching,
          ordering,
          select,
          responsive,
          destroy: true,
          drawCallback: function () {
            console.log(
              "DataTable drawn successfully =================================>"
            );
          },
          initComplete: function () {
            const api = this.api();

            // Add custom header filter functionality
            this.api()
              .columns()
              .every(function (index) {
                const column = this;
                const columnInfo = dtColumns[index];
                const columnName = columnInfo.title || `Column ${index + 1}`;

                // Get the column key/name to match against filteredColumns
                const columnKey =
                  typeof columnInfo.data === "string"
                    ? columnInfo.data
                    : columnInfo.name || columnName;

                // Create the header filter container
                const headerCell = $(column.header());

                // Check if this column should have a filter
                const shouldHaveFilter =
                  filteredColumns.includes(columnKey) ||
                  filteredColumns.includes(index) ||
                  filteredColumns.includes(columnName);

                if (shouldHaveFilter) {
                  // Create filterable header
                  headerCell.html(
                    `<div class="dt-header-filter">${columnName}</div>`
                  );

                  // Create dropdown container
                  const dropdownContainer = $(
                    '<div class="dt-filter-dropdown"></div>'
                  );
                  const allOption = $(
                    '<div class="dt-filter-option all" data-value="">All</div>'
                  );
                  dropdownContainer.append(allOption);

                  // Add unique values to the dropdown
                  const uniqueValues = new Set();
                  column.data().each(function (value) {
                    if (
                      value !== null &&
                      value !== undefined &&
                      value.toString().trim() !== ""
                    ) {
                      uniqueValues.add(value.toString());
                    }
                  });

                  // Sort and add values to dropdown
                  Array.from(uniqueValues)
                    .sort()
                    .forEach((value) => {
                      const option = $(
                        `<div class="dt-filter-option" data-value="${value}">${value}</div>`
                      );
                      dropdownContainer.append(option);
                    });

                  // Append dropdown to header
                  headerCell.append(dropdownContainer);
                  headerCell
                    .find(".dt-header-filter")
                    .on("click", function (e) {
                      e.stopPropagation();

                      // Close all other dropdowns
                      $(".dt-filter-dropdown")
                        .not(dropdownContainer)
                        .removeClass("show");
                      dropdownContainer.toggleClass("show");
                    });

                  // Handle option selection
                  dropdownContainer.on(
                    "click",
                    ".dt-filter-option",
                    function () {
                      const value = $(this).data("value");

                      // Update active state
                      dropdownContainer
                        .find(".dt-filter-option")
                        .removeClass("active");
                      $(this).addClass("active");

                      // Apply filter
                      const val = $.fn.dataTable.util.escapeRegex(value);
                      column.search(val ? `^${val}$` : "", true, false).draw();
                      dropdownContainer.removeClass("show");
                    }
                  );
                } else {
                  // Just set the plain header text without filter functionality
                  headerCell.html(columnName);
                }
              });
            $(document).on("click", function () {
              $(".dt-filter-dropdown").removeClass("show");
            });
          },
        });
      } catch (error) {
        console.error("Error initializing DataTable:", error);
      }
    }, 300);

    return () => {
      clearTimeout(timer);
      //  click handler
      if (window.jQuery) {
        window.jQuery(document).off("click.dt-filter");
      }
    };
  }, [
    scriptsReady,
    data,
    columns,
    loading,
    pageLength,
    searching,
    ordering,
    select,
    responsive,
    filteredColumns,
  ]);

  return (
    <div className="datatable-container relative">
      {/* LoadingOverlay Component */}
      <LoadingOverlay 
        isLoading={loading}
        text={loadingText}
        size="default"
        showProgressBar={true}
        isOverlay={true}
        className="rounded-lg mt-80"
      />
      
      <div className="overflow-x-auto scrollbar rounded-lg shadow-md">
        <table
          ref={tableRef}
          id={id}
          className={`min-w-full divide-y divide-zinc-500 text-gray-100 border-2 border-zinc-500 text-xs ${className}`}
        >
          <thead className="bg-zinc-900 text-sm text-gray-500 uppercase tracking-wider">
            {/* Headers will be populated by DataTables with filter dropdowns */}
          </thead>
          <tbody className="bg-zinc-800 divide-y divide-zinc-500 break-words">
            {/* Body will be populated by DataTables */}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default DataTable;
"use client";

import React, { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import Sidebar from "@/app/components/platform/main-sidebar";
import Header from "@/app/components/platform/header";
import CreateNewButton from "@/app/components/add-new-button";
import { GuardrailApi } from "@/app/utils/ai-studio-endpoint/guardrails-api";

//  DataTable component with dynamic import
const DataTable = dynamic(() => import("@/app/components/table"), {
  ssr: false,
});

const GuardrailsTableClient = () => {
  const [guardrails, setGuardrails] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [error, setError] = useState(null);

  // Function to fetch the data

  const fetchGuardrailsData = async () => {
    try {
      const data = await GuardrailApi.getGuardrail(); // Use your API service
      return data.output || [];
    } catch (error) {
      console.error("Error fetching model nodes data:", error);
      return [];
    }
  };

  const transformGuardrailData = (guardrailItems) => {
    return guardrailItems.map((item) => ({
      Name: item.name || "N/A",
      "Guardrail Id": item.guardrailId || "N/A",
      Status: item.resourceBase?.status || "N/A",
      "View Details": `<a href="/ai-center/guardrails/${item.guardrailId}" target="_blank" title="Open in new tab">
                       <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" class="h-5 w-5" viewBox="0 0 40 40">
                        <path fill="rgb(251 146 60 / var(--tw-bg-opacity, 1))" d="M1.5 5.5H34.5V38.5H1.5z"></path>
                        <path fill="rgb(154 52 18 / var(--tw-bg-opacity, 1))" d="M34,6v32H2V6H34 M35,5H1v34h34V5L35,5z"></path>
                        <path fill="rgb(251 146 60 / var(--tw-bg-opacity, 1))" d="M30.611 13.611H37.055V15.944H30.611z" transform="rotate(90 33.833 14.778)"></path>
                        <path fill="rgb(251 146 60 / var(--tw-bg-opacity, 1))" d="M22 5H28.444V7.333H22z"></path>
                        <g>
                            <path fill="rgb(251 146 60 / var(--tw-bg-opacity, 1))" d="M18.707 16L28.707 6 24.207 1.5 38.5 1.5 38.5 15.793 34 11.293 24 21.293z"></path>
                            <path fill="rgb(154 52 18 / var(--tw-bg-opacity, 1))" d="M38,2v12.586l-3.293-3.293L34,10.586l-0.707,0.707L24,20.586L19.414,16l9.293-9.293L29.414,6 l-0.707-0.707L25.414,2H38 M39,1H23l5,5L18,16l6,6l10-10l5,5V1L39,1z"></path>
                        </g>
                      </svg>
                      </a>`,
    }));
  };

  // For loading data
  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      const guardrailsData = await fetchGuardrailsData();
      console.log("Guardrails data fetched:", guardrailsData);
      const transformedData = transformGuardrailData(guardrailsData);
      console.log("Transformed data:", transformedData);
      setGuardrails(transformedData);
      setIsLoading(false);
    };

    loadData();
  }, []);

  // Define columns for the DataTable
  const columns = ["Name", "Guardrail Id", "Status", "View Details"];
  const filteredColumns = ["Name"];

  return (
    <div className="bg-zinc-800 flex h-screen">
      {/* Sidebar Component */}
      <Sidebar />

      <div className="overflow-y-auto h-screen w-full">
        <Header
          sectionHeader="AI Guardrails"
          hideBackListingLink={true}
          backListingLink="./"
        />

        {isModalOpen && (
          <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
            <div className="bg-zinc-800 p-6 rounded-lg w-3/4 max-h-3/4 overflow-y-auto">
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl text-white">YAML Editor</h2>
                <button
                  onClick={() => setIsModalOpen(false)}
                  className="text-gray-400 hover:text-white"
                >
                  Close
                </button>
              </div>
            </div>
          </div>
        )}

        <div className="flex-grow p-2 w-full">
          <section id="tables" className="space-y-6">
            <div className="flex justify-between mb-2 items-center p-2">
              <CreateNewButton href="./guardrails/create" label="Add New" />
            </div>

            {error ? (
              <div className="flex justify-center items-center h-64 text-red-400">
                <div className="text-xl">Error: {error}</div>
              </div>
            ) : guardrails.length === 0 && !isLoading ? (
              <div className="overflow-x-auto scrollbar rounded-lg p-4 shadow-md text-gray-100">
                <table className="min-w-full divide-y divide-zinc-500 text-sm text-gray-100 border-2 border-zinc-500">
                  <thead className="bg-zinc-900 text-sm font-medium text-gray-500 uppercase tracking-wider"></thead>
                  <tbody className="bg-zinc-800 divide-y divide-zinc-500 break-words text-sm">
                    <tr>
                      <td
                        colSpan="4"
                        className="px-3 py-3 whitespace-nowrap text-center"
                      >
                        No AI Guardrails found
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="p-4">
                <DataTable
                  id="guardrailsDataTable"
                  data={guardrails}
                  columns={columns}
                  loading={isLoading}
                  filteredColumns={filteredColumns}
                />
              </div>
            )}
          </section>
        </div>
        <div hidden id="createTemplate"></div>
      </div>
    </div>
  );
};

// Dynamic import with no SSR for the component
const GuardrailsTable = dynamic(() => Promise.resolve(GuardrailsTableClient), {
  ssr: false,
});

export default GuardrailsTable;

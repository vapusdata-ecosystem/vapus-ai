"use client";

import React, { useState, useEffect } from "react";
import dynamic from "next/dynamic";
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
      const data = await GuardrailApi.getGuardrail(); 
      return data.output || [];
    } catch (error) {
      console.error("Error fetching  Guardrail data:", error);
        setError(error.message || "Frror fetching Guardrail data");
      return [];
    }
  };

  const transformGuardrailData = (guardrailItems) => {
    return guardrailItems.map((item) => ({
      Name: item.name || "N/A",
      "Guardrail Id": item.guardrailId || "N/A",
      Status: item.resourceBase?.status || "N/A",

      "View Details": `<a href="/ai-center/guardrails/${item.guardrailId}" target="_blank" class="relative group inline-flex items-center justify-center">
                          <!-- The Icon -->
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 200 200">
                            <circle cx="100" cy="100" r="90" stroke="rgb(207, 86, 46)" stroke-width="10" fill="none" />
                            <g transform="rotate(315, 100, 100)">
                              <line x1="60" y1="100" x2="140" y2="100" stroke="rgb(207, 86, 46)" stroke-width="10" stroke-linecap="round" />
                              <path d="M120,80 L140,100 L120,120" stroke="rgb(207, 86, 46)" stroke-width="10" stroke-linecap="round" stroke-linejoin="round" fill="none" />
                            </g>
                          </svg>

                          <!-- Tooltip: adjust position and spacing -->
                          <div class="absolute -top-6 left-20 -translate-x-1/2 
                                      hidden group-hover:block bg-gray-700 text-gray-100 text-xs rounded px-2 py-1 z-50 whitespace-nowrap">
                            View guardrail detail
                          </div>
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

        <div className="flex-grow w-full">
          <section id="tables" className="space-y-6">
           
              <CreateNewButton href="./guardrails/create" label="Add New" />
          

            {error ? (
              <div className="flex justify-center items-center h-64 text-red-400">
                <div className="text-xl">Error: {error}</div>
              </div>
            ) : (
              <div className="p-4">
                <DataTable
                  id="guardrailsDataTable"
                  data={guardrails}
                  columns={columns}
                  loading={isLoading}
                  filteredColumns={filteredColumns}
                  loadingText="Loading Guardrails..."
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
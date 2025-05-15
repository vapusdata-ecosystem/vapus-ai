"use client";

import React, { useEffect, useState } from "react";
import dynamic from "next/dynamic";
import Header from "@/app/components/platform/header";
import CreateNewButton from "@/app/components/add-new-button";
import Sidebar from "@/app/components/platform/main-sidebar";
import { modelsRegistryApi } from "@/app/utils/ai-studio-endpoint/models-registry-api";

//  DataTable component with dynamic import
const DataTable = dynamic(() => import("@/app/components/table"), {
  ssr: false,
});

const ModelNodesTable = () => {
  const [modelNodes, setModelNodes] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);

  // Function to fetch the data

  const fetchModelNodesData = async () => {
    try {
      const data = await modelsRegistryApi.getModelsRegistry();
      return data.output.aiModelNodes || [];
    } catch (error) {
      console.error("Error fetching model nodes data:", error);
      return [];
    }
  };

  const transformModelNodeData = (nodes) => {
    return nodes.map((node) => ({
      Name: node.name || "N/A",
      "Service Provider": node.attributes.serviceProvider || "N/A",
      "Node Id": node.modelNodeId || "N/A",
      Owner: node.resourceBase.domain || "N/A",
      Status: node.resourceBase.status || "N/A",
      "View Details": `<a href="/ai-center/models-registry/${node.modelNodeId}" target="_blank" title="Open in new tab">
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
      const modelNodesData = await fetchModelNodesData();
      const transformedData = transformModelNodeData(modelNodesData);
      setModelNodes(transformedData);
      setIsLoading(false);
    };

    loadData();
  }, []);

  // Define columns for the DataTable
  const columns = [
    "Name",
    "Service Provider",
    "Node Id",
    "Owner",
    "Status",
    "View Details",
  ];
  const filteredColumns = ["Name", "Service Provider"];

  return (
    <div className="bg-zinc-800 flex h-screen">
      {/* Sidebar Component */}
      <Sidebar />

      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="AI Models"
          hideBackListingLink={true}
          backListingLink="./"
        />

        {isModalOpen && (
          <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
            <div className="bg-zinc-800 p-6 rounded-lg w-3/4 max-h-3/4 overflow-y-auto scrollbar">
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
              <CreateNewButton
                href="./models-registry/create"
                label="Add New"
              />
            </div>

            <div className="p-4">
              <DataTable
                id="modelNodesDataTable"
                data={modelNodes}
                columns={columns}
                loading={isLoading}
                filteredColumns={filteredColumns}
                showCsvExport={true}
              />
            </div>
          </section>
        </div>
        <div hidden id="createTemplate"></div>
      </div>
    </div>
  );
};

// Dynamic import with no SSR for the component that uses jQuery/DataTables
const ModelNodesPage = dynamic(() => Promise.resolve(ModelNodesTable), {
  ssr: false,
});

export default ModelNodesPage;

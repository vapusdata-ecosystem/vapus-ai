"use client";
import React, { useEffect, useState } from "react";
import dynamic from "next/dynamic";
import Header from "@/app/components/platform/header";
import CreateNewButton from "@/app/components/add-new-button";
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
      "View Details": `<a href="/ai-center/models-registry/${node.modelNodeId}" target="_blank" class="relative group">
                          <!-- The Icon -->
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 200 200">
                            <circle cx="100" cy="100" r="90" stroke="rgb(207, 86, 46)" stroke-width="10" fill="none" />
                            <g transform="rotate(315, 100, 100)">
                              <line x1="60" y1="100" x2="140" y2="100" stroke="rgb(207, 86, 46)" stroke-width="10" stroke-linecap="round" />
                              <path d="M120,80 L140,100 L120,120" stroke="rgb(207, 86, 46)" stroke-width="10" stroke-linecap="round" stroke-linejoin="round" fill="none" />
                            </g>
                          </svg>

                          <!-- Custom Tooltip -->
                          <div class="absolute bottom-full left-1/2 transform -translate-x-1/2 ml-2
                                      hidden group-hover:block bg-gray-700 text-gray-100 text-xs rounded px-2 py-1 z-50 whitespace-nowrap">
                            View model detail
                          </div>
                      </a> `,
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

        <div className="flex-grow  w-full">
          <section id="tables" className="space-y-6">
   
              <CreateNewButton
                href="./models-registry/create"
                label="Add New"
              />
     

            <div className="p-4">
              <DataTable
                id="modelNodesDataTable"
                data={modelNodes}
                columns={columns}
                loading={isLoading}
                loadingText="Loading AI Models..."
                filteredColumns={filteredColumns}
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
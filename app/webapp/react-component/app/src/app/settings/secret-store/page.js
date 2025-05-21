"use client";

import React, { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import Sidebar from "@/app/components/platform/main-sidebar";
import Header from "@/app/components/platform/header";
import CreateNewButton from "@/app/components/add-new-button";
import { secretStoreApi } from "@/app/utils/settings-endpoint/secret-store-api";

// DataTable component with dynamic import
const DataTable = dynamic(() => import("@/app/components/table"), {
  ssr: false,
});

const SecretStoreClient = () => {
  const [secrets, setSecrets] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [createTemplate, setCreateTemplate] = useState("");

  // Function to fetch the secrets data

  const fetchSecretsData = async () => {
    try {
      const data = await secretStoreApi.getSecretStore(); // Use your API service
      return data.output || [];
    } catch (error) {
      console.error("Error fetching model nodes data:", error);
      setError(error.message || "Failed to fetch plugins data");
      return [];
    }
  };

  // Convert epoch timestamp to readable date format
  const epochConverterFull = (epoch) => {
    if (!epoch || epoch === "0") return "N/A";
    const date = new Date(parseInt(epoch) * 1000);
    return date.toLocaleString();
  };

  // Limit text length with ellipsis
  const limitLetters = (text, limit) => {
    if (!text) return "N/A";
    return text.length > limit ? text.substring(0, limit) + "..." : text;
  };

  const transformSecretsData = (secretItems) => {
    return secretItems.map((item) => {
      return {
        Name: item.name || "N/A",
        "Secret Type":
          item.secretType === "CUSTOM_SECRET"
            ? "Custom Secret"
            : "Vapus Credential",
        "Created At": epochConverterFull(item.resourceBase?.createdAt),
        Description: limitLetters(item.description, 30),
        // Use a string with HTML instead of a React component for DataTable

        "View Details": `<a href="secret-store/${item.name}" target="_blank" class="relative group inline-flex items-center justify-center">
                          <!-- The Icon -->
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 200 200">
                            <circle cx="100" cy="100" r="90" stroke="rgb(207, 86, 46)" stroke-width="10" fill="none" />
                            <g transform="rotate(315, 100, 100)">
                              <line x1="60" y1="100" x2="140" y2="100" stroke="rgb(207, 86, 46)" stroke-width="10" stroke-linecap="round" />
                              <path d="M120,80 L140,100 L120,120" stroke="rgb(207, 86, 46)" stroke-width="10" stroke-linecap="round" stroke-linejoin="round" fill="none" />
                            </g>
                          </svg>

                          <!-- Tooltip: adjust position and spacing -->
                          <div class="absolute -top-6 left-17 -translate-x-1/2 
                                      hidden group-hover:block bg-gray-700 text-gray-100 text-xs rounded px-2 py-1 z-50 whitespace-nowrap">
                            View user detail
                          </div>
                      </a>`,
      };
    });
  };

  // For loading data
  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      const secretsData = await fetchSecretsData();
      console.log("Secrets data fetched:", secretsData);
      const transformedData = transformSecretsData(secretsData);
      console.log("Transformed data:", transformedData);
      setSecrets(transformedData);
      setIsLoading(false);
    };

    loadData();
  }, []);

  // Define columns for the DataTable
  const columns = [
    "Name",
    "Secret Type",
    "Created At",
    "Description",
    "View Details",
  ];
  const filteredColumns = ["Name", "Secret Type"];

  const createNewResource = () => {
    console.log("Create new resource clicked");
    // Add your implementation here
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      {/* Sidebar Component */}
      <Sidebar />

      <div className="overflow-y-auto h-screen w-full">
        <Header
          sectionHeader="Secret Service"
          hideBackListingLink={true}
          backListingLink="./"
        />

        <div className="flex-grow p-2 w-full">
          <section id="tables" className="space-y-6">
            <div className="flex justify-between mb-2 items-center p-2">
              <CreateNewButton href="./secret-store/create" label="Add New" />
            </div>
            {error ? (
              <div className="flex justify-center items-center h-64 text-red-400">
                <div className="text-xl">Error: {error}</div>
              </div>
            ) : secrets.length === 0 && !isLoading ? (
              <div className="overflow-x-auto rounded-lg p-4 shadow-md text-gray-100">
                <table
                  className="min-w-full divide-y divide-zinc-500 text-gray-100 border-2 border-zinc-500"
                  id="secretStoreTable"
                >
                  <thead className="bg-zinc-900 divide-y divide-zinc-500 break-words   font-medium text-gray-500 uppercase">
                    <tr></tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200  ">
                    <tr>
                      <td
                        colSpan="5"
                        className="px-3 py-3 whitespace-nowrap text-center"
                      >
                        No Secrets found
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="overflow-x-auto rounded-lg p-4 shadow-md text-gray-100">
                <DataTable
                  id="secretStoreTable"
                  data={secrets}
                  columns={columns}
                  loading={isLoading}
                  filteredColumns={filteredColumns}
                  dangerouslySetInnerHTML={true}
                />
              </div>
            )}
          </section>
        </div>
        <div hidden id="createTemplate">
          {createTemplate}
        </div>
      </div>
    </div>
  );
};

// Dynamic import with no SSR for the component
const SecretStore = dynamic(() => Promise.resolve(SecretStoreClient), {
  ssr: false,
});

export default SecretStore;

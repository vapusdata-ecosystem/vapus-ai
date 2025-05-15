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
        Update: `<a href="secret-store/${item.name}" 
                      class="cursor-pointer mt-2 mb-2 rounded-full bg-orange-700 px-5 py-2 font-medium text-white transition-colors duration-200 hover:bg-pink-900 active:bg-orange-900" role="button">
                      Update
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
    "Update",
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

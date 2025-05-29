"use client";

import React, { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import Header from "@/app/components/platform/header";
import { pluginsApi } from "@/app/utils/settings-endpoint/plugins-api";
import CreateNewButton from "@/app/components/add-new-button";

const DataTable = dynamic(() => import("@/app/components/table"), {
  ssr: false,
});

const PluginSettingsClient = () => {
  const [plugins, setPlugins] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [createTemplate, setCreateTemplate] = useState("");

  //  fetch the plugins data
  const fetchPluginsData = async () => {
    try {
      const data = await pluginsApi.getPlugins();
      return data.output || [];
    } catch (error) {
      console.error("Error fetching model nodes data:", error);
      setError(error.message || "Failed to fetch plugins data");
      return [];
    }
  };

  const transformPluginsData = (pluginItems) => {
    return pluginItems.map((item) => {
      const status = item.resourceBase?.status || "N/A";
      return {
        "PLUGIN SERVICE": item.pluginService || "N/A",
        "PLUGIN TYPE": item.pluginType || "N/A",
        NAME: item.name || "N/A",
        SCOPE: item.scope || "N/A",
        Status: `<span class="px-3 py-1 text-sm font-medium ${
          status === "ACTIVE"
            ? "text-green-800 bg-green-100"
            : "text-red-800 bg-red-100"
        } rounded-full">${status}</span>`,

        "View Details": `<a href="/settings/plugins/${item.pluginId}" target="_blank" class="relative group">
                          <!-- The Icon -->
                          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 200 200">
                            <circle cx="100" cy="100" r="90" stroke="rgb(207, 86, 46)" stroke-width="10" fill="none" />
                            <g transform="rotate(315, 100, 100)">
                              <line x1="60" y1="100" x2="140" y2="100" stroke="rgb(207, 86, 46)" stroke-width="10" stroke-linecap="round" />
                              <path d="M120,80 L140,100 L120,120" stroke="rgb(207, 86, 46)" stroke-width="10" stroke-linecap="round" stroke-linejoin="round" fill="none" />
                            </g>
                          </svg>

                           <!-- Tooltip: adjust position and spacing -->
                          <div class="absolute -top-6 left-18 -translate-x-1/2 
                                      hidden group-hover:block bg-gray-700 text-gray-100 text-xs rounded px-2 py-1 z-50 whitespace-nowrap">
                            View plugins detail
                          </div>
                      </a> `,
      };
    });
  };

  // loading data
  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      const pluginsData = await fetchPluginsData();
      console.log("Plugins data fetched:", pluginsData);
      const transformedData = transformPluginsData(pluginsData);
      console.log("Transformed data:", transformedData);
      setPlugins(transformedData);
      setIsLoading(false);
    };

    loadData();
  }, []);

  // Define columns for the DataTable
  const columns = [
    "PLUGIN SERVICE",
    "PLUGIN TYPE",
    "NAME",
    "SCOPE",
    "Status",
    "View Details",
  ];
  const filteredColumns = ["PLUGIN SERVICE", "PLUGIN TYPE", "NAME"];

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Plugins"
          hideBackListingLink={true}
          backListingLink="./"
        />

        <div className="flex-grow p-2 w-full">
          <div className="flex justify-between mb-2 items-center p-2">
            <CreateNewButton href="./plugins/create" label="Add New" />
          </div>
          <section id="tables" className="space-y-6">
            {error ? (
              <div className="flex justify-center items-center h-64 text-red-400">
                <div className="text-xl">Error: {error}</div>
              </div>
            ) : plugins.length === 0 && !isLoading ? (
              <div className="overflow-x-auto scrollbar rounded-lg p-4 shadow-md text-gray-100">
                <table
                  className="min-w-full divide-y divide-zinc-500 text-xs text-gray-100 border-2 border-zinc-500"
                  id="pluginSettingsTable"
                >
                  <thead className="bg-zinc-900 divide-y divide-zinc-500 break-words text-sm font-medium text-gray-500 uppercase">
                    <tr></tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200 text-sm">
                    <tr>
                      <td
                        colSpan="6"
                        className="px-3 py-3 whitespace-nowrap text-center"
                      >
                        No Plugins found
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="overflow-x-auto scrollbar rounded-lg p-4  text-gray-100">
                <DataTable
                  id="pluginSettingsTable"
                  data={plugins}
                  columns={columns}
                  loading={isLoading}
                  filteredColumns={filteredColumns}
                  dangerouslySetInnerHTML={true}
                />
              </div>
            )}
          </section>
        </div>
      </div>
      <div hidden id="createTemplate">
        {createTemplate}
      </div>
    </div>
  );
};

// Dynamic import with no SSR for the component
const PluginSettings = dynamic(() => Promise.resolve(PluginSettingsClient), {
  ssr: false,
});

export default PluginSettings;

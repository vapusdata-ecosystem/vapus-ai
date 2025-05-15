"use client";

import React, { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import Sidebar from "@/app/components/platform/main-sidebar";
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
        "View Details": `<a href="/settings/plugins/${item.pluginId}" target="_blank" title="Open in new tab">
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
      {/* Sidebar Component */}
      <Sidebar />

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

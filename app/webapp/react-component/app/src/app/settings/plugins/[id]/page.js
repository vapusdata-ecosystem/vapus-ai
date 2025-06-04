"use client";

import React, { useState, useEffect } from "react";
import { use } from "react";
import Header from "@/app/components/platform/header";
import {
  pluginsApi,
  pluginsArchiveApi,
} from "@/app/utils/settings-endpoint/plugins-api";
import SectionHeaders from "@/app/components/section-headers";
import LoadingOverlay from "@/app/components/loading/loading";

export default function PluginDetailsPage({ params }) {
  const unwrappedParams = use(params);
  const pluginId = unwrappedParams?.id ? String(unwrappedParams.id).trim() : "";
  const [activeTab, setActiveTab] = useState("basic-info");
  const [pluginData, setPluginData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Function to handle tab switching
  const showTab = (tabId) => {
    setActiveTab(tabId);
  };

  useEffect(() => {
    const fetchPluginDetails = async () => {
      if (!pluginId) {
        console.error("No plugin ID available");
        setError("No plugin ID provided");
        setLoading(false);
        return;
      }

      try {
        const response = await pluginsApi.getPluginsId(pluginId);
        console.log("API Response:", response);

        if (!response) {
          console.error("No response received from server");
          setError("No response received from server");
          setLoading(false);
          return;
        }

        // Better response handling with more detailed logging
        if (response.output) {
          if (Array.isArray(response.output) && response.output.length > 0) {
            setPluginData(response.output[0]);
          } else if (typeof response.output === "object") {
            console.log("Setting plugin data from object:", response.output);
            setPluginData(response.output);
          } else {
            console.error("Unexpected output format:", response.output);
            setError("Unexpected data format received");
          }
        } else {
          console.log(
            "Response has no output property, checking if response itself is the data"
          );
          if (response.pluginId) {
            console.log("Using response as plugin data:", response);
            setPluginData(response);
          } else {
            console.error("Response does not contain plugin data:", response);
            setError("Invalid response format");
          }
        }

        setLoading(false);
      } catch (err) {
        console.error("Error fetching plugin details:", err);
        setError(err.message || "Failed to fetch plugin details");
        setLoading(false);
      }
    };

    if (pluginId) {
      fetchPluginDetails();
    } else {
      setError("No plugin ID provided");
      setLoading(false);
    }
  }, [pluginId]);

  if (error || (!loading && !pluginData)) {
    return (
      <div className="bg-zinc-800 flex h-screen">
        <div className="overflow-y-auto scrollbar h-screen w-full">
          <Header
            sectionHeader="Plugin Details"
            backListingLink="/settings/plugins"
          />
          <div className="flex justify-center items-center h-64 text-red-400">
            <div className="text-xl">Error: {error || "Plugin not found"}</div>
          </div>
        </div>
      </div>
    );
  }

  // Create a data object for SectionHeaders
  const apiServices = {
    plugin: {
      archive: pluginsArchiveApi.getPluginsArchive,
      delete: pluginsArchiveApi.getPluginsArchive,
    },
  };
  
  const headerResourceData = pluginData ? {
    id: pluginData.pluginId,
    name: pluginData.name || "Unnamed Plugin",
    createdAt: pluginData.resourceBase?.createdAt
      ? parseInt(pluginData.resourceBase.createdAt) * 1000
      : null,
    createdBy: pluginData.resourceBase?.createdBy,
    status: pluginData.resourceBase?.status || pluginData.status,
    resourceBase: pluginData.resourceBase,
    resourceType: "plugin",
    // Create action params for update functionality
    createActionParams: pluginData.createActionParams || {
      weblink: `/settings/plugins/${pluginData.pluginId}/update`,
    },
  } : null;

  return (
    <div className="bg-zinc-800 flex h-screen relative">
      <LoadingOverlay 
        isLoading={loading} 
        text="Loading plugin details"
        size="default"
        isOverlay={true}
        className="absolute bg-zinc-800 inset-0 z-10"
      />
      
      <div className="overflow-y-auto h-screen w-full">
        <Header
          sectionHeader="Plugin Details"
          hideBackListingLink={false}
          backListingLink="/settings/plugins"
        />

        <div className="flex-grow p-2 w-full text-gray-100">
          {pluginData && pluginData.pluginId && headerResourceData && (
            <SectionHeaders
              resourceId={pluginData.pluginId}
              resourceData={headerResourceData}
              resourceType="plugin"
              apiServices={apiServices}
            />
          )}

          {/* Tabs */}
          <div className="overflow-x-auto text-gray-100 bg-zinc-800 rounded-lg p-8 shadow-md">
            <div className="flex border-b border-zinc-500 font-semibold text-gray-50">
              <button
                onClick={() => showTab("basic-info")}
                className={`px-4 py-2 focus:outline-none ${
                  activeTab === "basic-info"
                    ? "bg-[oklch(0.205_0_0)] text-white rounded-t-[10px]"
                    : ""
                }`}
              >
                Basic Info
              </button>
            </div>

            {/* Tab Content */}
            <div
              id="basic-info"
              className={`mt-2 bg-[#1b1b1b] p-4 ${
                activeTab !== "basic-info" ? "hidden" : ""
              }`}
            >
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                {/* Display Name */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Display Name:
                  </p>
                  <p className="p-2">{pluginData?.name || "N/A"}</p>
                </div>

                {/* Plugin id */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Plugin id:
                  </p>
                  <p className="p-2">{pluginData?.pluginId || "N/A"}</p>
                </div>

                {/* Organization */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Organization:
                  </p>
                  <p className="p-2">{pluginData?.organization || "N/A"}</p>
                </div>

                {/* Plugin Type */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Plugin Type:
                  </p>
                  <p className="p-2">{pluginData?.pluginType || "N/A"}</p>
                </div>

                {/* Plugin Service */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Plugin Service:
                  </p>
                  <p className="p-2">{pluginData?.pluginService || "N/A"}</p>
                </div>

                {/* Status */}
                <div className="lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Status:
                  </p>
                  <p className="p-2">
                    <span
                      className={`px-3 py-1 text-sm font-medium ${
                        pluginData?.status === "ACTIVE" ||
                        pluginData?.resourceBase?.status === "ACTIVE"
                          ? "text-green-800 bg-green-100"
                          : "text-red-800 bg-red-100"
                      } rounded-full`}
                    >
                      {pluginData?.status ||
                        pluginData?.resourceBase?.status ||
                        "N/A"}
                    </span>
                  </p>
                </div>
              </div>

              {/* Parameters Section */}
              <div className="mt-6">
                <p className="text-xl mb-4 text-[1.25rem] font-bold text-[#f4d1c2] underline">
                  Parameters:
                </p>

                {Array.isArray(pluginData?.dynamicParams) &&
                pluginData.dynamicParams.length > 0 ? (
                  <div className="w-full bg-[#1b1b1b] rounded-lg">
                    <table className="min-w-full divide-y divide-zinc-500 text-gray-100 border-2 border-zinc-500 text-xs">
                      <thead className="bg-zinc-900 divide-y divide-zinc-500 break-words text-sm font-medium text-gray-500 uppercase">
                        <tr>
                          <th className="px-3 py-3 text-left tracking-wider">
                            Key
                          </th>
                          <th className="px-3 py-3 text-left tracking-wider">
                            Value
                          </th>
                        </tr>
                      </thead>
                      <tbody className="bg-zinc-800 divide-y divide-zinc-500">
                        {pluginData.dynamicParams.some(
                          (param) => param && param.key && param.key !== ""
                        ) ? (
                          pluginData.dynamicParams.map((param, index) =>
                            param && param.key && param.key !== "" ? (
                              <tr key={index}>
                                <td className="px-3 py-3 whitespace-nowrap">
                                  {param.key}
                                </td>
                                <td className="px-3 py-3 whitespace-nowrap">
                                  {param.value || ""}
                                </td>
                              </tr>
                            ) : null
                          )
                        ) : (
                          <tr>
                            <td colSpan="2" className="px-3 py-3 text-center">
                              No parameters available for this plugin.
                            </td>
                          </tr>
                        )}
                      </tbody>
                    </table>
                  </div>
                ) : (
                  <p className="p-2">
                    No parameters available for this plugin.
                  </p>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
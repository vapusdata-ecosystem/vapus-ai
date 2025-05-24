"use client";

import React, { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import Header from "@/app/components/platform/header";
import { platformDomainApi } from "@/app/utils/settings-endpoint/platform-domain";
import CreateNewButton from "@/app/components/add-new-button";

const DataTable = dynamic(() => import("@/app/components/table"), {
  ssr: false,
});

const PlatformDomains = () => {
  const [domains, setDomains] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [createTemplate, setCreateTemplate] = useState("");

  // Function to fetch the domains data
  const fetchDomainsData = async () => {
    try {
      const data = await platformDomainApi.getplatformdomain();
      return data.output?.domains || [];
    } catch (error) {
      console.error("Error fetching model nodes data:", error);
      setError(error.message || "Failed to fetch plugins data");
      return [];
    }
  };

  const transformDomainsData = (domainItems) => {
    return domainItems.map((item) => {
      const domainId = item.domainId;
      console.log("Checking domain ID:", domainId);
      let hasAccess = true;

      console.log("Has access:", hasAccess);

      const tickSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="inline"><polyline points="20 6 9 17 4 12"></polyline></svg>`;

      return {
        "Display Name": item.displayName || "N/A",
        Status: item.status || "N/A",
        Id: domainId || "N/A",
        "Domain Type": item.domainType || "N/A",
        "Has Access": `<span class="px-3 py-1 text-sm font-medium ${
          hasAccess ? "text-green-800 bg-green-100" : "text-red-800 bg-red-100"
        } rounded-full">${hasAccess ? tickSvg : "No"}</span>`,
      };
    });
  };

  // For loading data
  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      const domainsData = await fetchDomainsData();
      console.log("Domains data fetched:", domainsData);
      const transformedData = transformDomainsData(domainsData);
      console.log("Transformed data:", transformedData);
      setDomains(transformedData);
      setIsLoading(false);
    };

    loadData();
  }, []);

  // Define columns for the DataTable
  const columns = ["Display Name", "Status", "Id", "Domain Type", "Has Access"];
  const filteredColumns = ["Display Name", "Status", "Domain Type"];

  // Function to create new domain
  const createNewDomain = () => {
    document.getElementById("actionTitle").innerHTML = "Create New Domain";
    document.getElementById("yamlSpecTitle").innerHTML = "Enter your spec here";
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Domains"
          hideBackListingLink={true}
          backListingLink="./"
        />

        <div className="flex-grow p-2 w-full">
          <div className="flex justify-end mb-2 items-center p-2">
            {/* Add your create new resource button here */}
            <CreateNewButton href="./platform-domain/create" label="Add New" />
          </div>

          <section id="tables" className="space-y-6">
            {error ? (
              <div className="flex justify-center items-center h-64 text-red-400">
                <div className="text-xl">Error: {error}</div>
              </div>
            ) : domains.length === 0 && !isLoading ? (
              <div className="overflow-x-auto scrollbar rounded-lg p-4 shadow-md text-gray-100">
                <table
                  className="min-w-full divide-y divide-zinc-500 text-xs text-gray-100 border-2 border-zinc-500"
                  id="platformDomainsTable"
                >
                  <thead className="bg-zinc-900 divide-y divide-zinc-500 break-words text-sm font-medium text-gray-500 uppercase">
                    <tr></tr>
                  </thead>
                  <tbody className="bg-zinc-800 divide-y divide-zinc-500">
                    <tr>
                      <td
                        colSpan="5"
                        className="px-3 py-3 whitespace-nowrap text-center"
                      >
                        No Domains found
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="overflow-x-auto scrollbar rounded-lg p-4 shadow-md text-gray-100">
                <DataTable
                  id="platformDomainsTable"
                  data={domains}
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
const PlatformDomainsPage = dynamic(() => Promise.resolve(PlatformDomains), {
  ssr: false,
});

export default PlatformDomainsPage;

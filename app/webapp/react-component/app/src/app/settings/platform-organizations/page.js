"use client";

import React, { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import Header from "@/app/components/platform/header";
import { platformDomainApi } from "@/app/utils/settings-endpoint/platform-organization";
import CreateNewButton from "@/app/components/add-new-button";
import { userGlobalData } from "@/context/GlobalContext";
import { userProfileApi } from "@/app/utils/settings-endpoint/profile-api";

const DataTable = dynamic(() => import("@/app/components/table"), {
  ssr: false,
});

const PlatformDomains = () => {
  const [Organization, setDomains] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [createTemplate, setCreateTemplate] = useState("");
  const [domainMapData, setDomainMapData] = useState({});

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Get user data from global context
        const globalContext = await userGlobalData();

        // Check if userId exists
        if (globalContext?.userInfo?.userId) {
          const userId = globalContext.userInfo.userId;
          // Make API call to get user profile with userId
          const data = await userProfileApi.getuserProfile(userId);
          console.log("Anand", data);

          if (data.output?.users && data.output.users.length > 0) {
            // Store the organizationMap for access checking
            if (data.organizationMap) {
              setDomainMapData(data.organizationMap);
            }
          } else {
            console.error("No users found in API response");
          }
        } else {
          console.error("User ID not found in global context");
        }
      } catch (error) {
        console.error("Error fetching user data:", error);
        setError(error.message || "Failed to fetch user data");
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  const fetchDomainsData = async () => {
    try {
      const data = await platformDomainApi.getplatformdomain();
      console.log("Organization", data);
      return data.output?.organizations || [];
    } catch (error) {
      console.error("Error fetching Organization data:", error);
      setError(error.message || "Failed to fetch Organization data");
      return [];
    }
  };

  const transformDomainsData = (domainItems) => {
    return domainItems.map((item) => {
      const organizationId = item.organizationId;
      const hasAccess = domainMapData.hasOwnProperty(organizationId);

      const accessGrantedSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" class="inline">
        <defs>
          <linearGradient id="greenGrad" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" style="stop-color:#10b981;stop-opacity:1" />
            <stop offset="100%" style="stop-color:#059669;stop-opacity:1" />
          </linearGradient>
        </defs>
        <circle cx="12" cy="12" r="10" fill="url(#greenGrad)" stroke="#047857" stroke-width="1"/>
        <path d="M8 12l3 3 5-6" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
        <circle cx="12" cy="12" r="10" fill="none" stroke="rgba(16, 185, 129, 0.3)" stroke-width="2" opacity="0.8"/>
      </svg>`;

      const accessDeniedSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" class="inline">
        <defs>
          <linearGradient id="redGrad" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" style="stop-color:#ef4444;stop-opacity:1" />
            <stop offset="100%" style="stop-color:#dc2626;stop-opacity:1" />
          </linearGradient>
        </defs>
        <circle cx="12" cy="12" r="10" fill="url(#redGrad)" stroke="#b91c1c" stroke-width="1"/>
        <path d="M8 8l8 8M16 8l-8 8" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
        <circle cx="12" cy="12" r="10" fill="none" stroke="rgba(239, 68, 68, 0.3)" stroke-width="2" opacity="0.8"/>
      </svg>`;

      return {
        "Display Name": item.displayName || "N/A",
        Status: item.status || "N/A",
        Id: item.organizationId || "N/A",
        "Domain Type": item.organizationType || "N/A",
        "Has Access": hasAccess
          ? `<div class="ml-8">${accessGrantedSvg}</div>`
          : `<div class="ml-8">${accessDeniedSvg}</div>`,
      };
    });
  };

  // For loading Organization data - only run after organizationMap is loaded
  useEffect(() => {
    if (Object.keys(domainMapData).length >= 0) {
      const loadData = async () => {
        setIsLoading(true);
        const domainsData = await fetchDomainsData();
        const transformedData = transformDomainsData(domainsData);
        setDomains(transformedData);
        setIsLoading(false);
      };

      loadData();
    }
  }, [domainMapData]);

  // Define columns for the DataTable
  const columns = ["Display Name", "Status", "Id", "Domain Type", "Has Access"];
  const filteredColumns = ["Display Name", "Status", "Domain Type"];

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Organization"
          hideBackListingLink={true}
          backListingLink="./"
        />

        <div className="flex-grow p-2 w-full">
          <div className="flex justify-end mb-2 items-center p-2">
            {/* Add your create new resource button here */}
            <CreateNewButton
              href="./platform-organizations/create"
              label="Add New"
            />
          </div>

          <section id="tables" className="space-y-6">
            {error ? (
              <div className="flex justify-center items-center h-64 text-red-400">
                <div className="text-xl">Error: {error}</div>
              </div>
            ) : Organization.length === 0 && !isLoading ? (
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
                        No Organization found
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="overflow-x-auto scrollbar rounded-lg p-4 shadow-md text-gray-100">
                <DataTable
                  id="platformDomainsTable"
                  data={Organization}
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

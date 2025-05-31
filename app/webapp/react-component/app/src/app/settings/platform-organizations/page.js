"use client";

import React, { useState, useEffect, useCallback, useMemo } from "react";
import dynamic from "next/dynamic";
import Header from "@/app/components/platform/header";
import { platformDomainApi } from "@/app/utils/settings-endpoint/platform-organization";
import CreateNewButton from "@/app/components/add-new-button";
import { userGlobalData } from "@/context/GlobalContext";
import { userProfileApi } from "@/app/utils/settings-endpoint/profile-api";
import { strTitle } from "@/app/components/JS/common";

const DataTable = dynamic(() => import("@/app/components/table"), {
  ssr: false,
});

const PlatformDomains = () => {
  const [Organization, setOrganization] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [createTemplate, setCreateTemplate] = useState("");
  const [organizationMapData, setOrganizationMapData] = useState({});
  const [dataVersion, setDataVersion] = useState(0);

  // Fetch user profile data
  useEffect(() => {
    let isMounted = true;

    const fetchData = async () => {
      try {
        setIsLoading(true);
        setError(null);

        // Get user data from global context
        const globalContext = await userGlobalData();

        if (!isMounted) return;

        // Check if userId exists
        if (globalContext?.userInfo?.userId) {
          const userId = globalContext.userInfo.userId;
          // Make API call to get user profile with userId
          const data = await userProfileApi.getuserProfile(userId);
          console.log("User Profile Data:", data);

          if (!isMounted) return;

          if (data.output?.users && data.output.users.length > 0) {
            // Store the organizationMap for access checking
            if (data.organizationMap) {
              setOrganizationMapData(data.organizationMap);
            }
          } else {
            console.error("No users found in API response");
            setError("No users found in API response");
          }
        } else {
          console.error("User ID not found in global context");
          setError("User ID not found in global context");
        }
      } catch (error) {
        console.error("Error fetching user data:", error);
        if (isMounted) {
          setError(error.message || "Failed to fetch user data");
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    };

    fetchData();

    return () => {
      isMounted = false;
    };
  }, []);

  //  fetch organization data
  const fetchOrganizationData = useCallback(async () => {
    try {
      const data = await platformDomainApi.getplatformdomain();
      console.log("Organization Data:", data);
      return data.output?.organizations || [];
    } catch (error) {
      console.error("Error fetching Organization data:", error);
      setError(error.message || "Failed to fetch Organization data");
      return [];
    }
  }, []);

  //  transform organization data
  const transformOrganizationData = useCallback(
    (organizationItems) => {
      return organizationItems.map((item, index) => {
        const organizationId = item.organizationId;
        const hasAccess = organizationMapData.hasOwnProperty(organizationId);

        const greenGradId = `greenGrad-${index}`;
        const redGradId = `redGrad-${index}`;

        const accessGrantedSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" class="inline">
        <defs>
          <linearGradient id="${greenGradId}" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" style="stop-color:#10b981;stop-opacity:1" />
            <stop offset="100%" style="stop-color:#059669;stop-opacity:1" />
          </linearGradient>
        </defs>
        <circle cx="12" cy="12" r="10" fill="url(#${greenGradId})" stroke="#047857" stroke-width="1"/>
        <path d="M8 12l3 3 5-6" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
        <circle cx="12" cy="12" r="10" fill="none" stroke="rgba(16, 185, 129, 0.3)" stroke-width="2" opacity="0.8"/>
      </svg>`;

        const accessDeniedSvg = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" class="inline">
        <defs>
          <linearGradient id="${redGradId}" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" style="stop-color:#ef4444;stop-opacity:1" />
            <stop offset="100%" style="stop-color:#dc2626;stop-opacity:1" />
          </linearGradient>
        </defs>
        <circle cx="12" cy="12" r="10" fill="url(#${redGradId})" stroke="#b91c1c" stroke-width="1"/>
        <path d="M8 8l8 8M16 8l-8 8" stroke="white" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
        <circle cx="12" cy="12" r="10" fill="none" stroke="rgba(239, 68, 68, 0.3)" stroke-width="2" opacity="0.8"/>
      </svg>`;

        return {
          "Display Name": item.displayName || "N/A",
          Status: item.status || "N/A",
          Id: item.organizationId || "N/A",
          "Organization Type": strTitle(item.organizationType || "N/A"),
          "Has Access": hasAccess
            ? `<div class="ml-8">${accessGrantedSvg}</div>`
            : `<div class="ml-8">${accessDeniedSvg}</div>`,
        };
      });
    },
    [organizationMapData]
  );

  // Load organization data when organizationMapData is available
  useEffect(() => {
    let isMounted = true;

    const loadData = async () => {
      // Only proceed if we have organizationMapData
      if (typeof organizationMapData === "object") {
        try {
          setIsLoading(true);
          setError(null);

          const organizationData = await fetchOrganizationData();

          if (!isMounted) return;

          const transformedData = transformOrganizationData(organizationData);
          setOrganization(transformedData);
          setDataVersion((prev) => prev + 1);
        } catch (error) {
          console.error("Error loading organization data:", error);
          if (isMounted) {
            setError(error.message || "Failed to load organization data");
          }
        } finally {
          if (isMounted) {
            setIsLoading(false);
          }
        }
      }
    };

    loadData();

    return () => {
      isMounted = false;
    };
  }, [organizationMapData, fetchOrganizationData, transformOrganizationData]);

  // Define columns for the DataTable
  const columns = useMemo(
    () => ["Display Name", "Status", "Id", "Organization Type", "Has Access"],
    []
  );

  const filteredColumns = useMemo(
    () => ["Display Name", "Status", "Organization Type"],
    []
  );

  // Empty state component
  const EmptyState = () => (
    <div className="overflow-x-auto scrollbar rounded-lg p-4 shadow-md text-gray-100">
      <table
        className="min-w-full divide-y divide-zinc-500 text-xs text-gray-100 border-2 border-zinc-500"
        id="platformDomainsTable"
      >
        <thead className="bg-zinc-900 divide-y divide-zinc-500 break-words text-sm font-medium text-gray-500 uppercase">
          <tr>
            {columns.map((column, index) => (
              <th key={index} className="px-3 py-3 text-left">
                {column}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="bg-zinc-800 divide-y divide-zinc-500">
          <tr>
            <td
              colSpan={columns.length}
              className="px-3 py-3 whitespace-nowrap text-center"
            >
              No Organization found
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );

  // Error state component
  const ErrorState = () => (
    <div className="flex justify-center items-center h-64 text-red-400">
      <div className="text-xl">Error: {error}</div>
    </div>
  );

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
            <CreateNewButton
              href="./platform-organizations/create"
              label="Add New"
            />
          </div>

          <section id="tables" className="space-y-6">
            {error ? (
              <ErrorState />
            ) : Organization.length === 0 && !isLoading ? (
              <EmptyState />
            ) : (
              <div className="overflow-x-auto scrollbar p-4  text-gray-100">
                <DataTable
                  key={`platform-Organization-${dataVersion}`}
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

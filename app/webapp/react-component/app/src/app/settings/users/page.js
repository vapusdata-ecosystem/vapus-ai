"use client";

import React, { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import Header from "@/app/components/platform/header";
import { userApi } from "@/app/utils/settings-endpoint/user-api";
import { getGlobalData } from "@/context/GlobalContext";

// DataTable component with dynamic import
const DataTable = dynamic(() => import("@/app/components/table"), {
  ssr: false,
});

const UsersTable = () => {
  const [users, setUsers] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  const [globalData, setGlobalData] = useState({
    currentDomain: null,
  });

  // take data from globleContext
  useEffect(() => {
    const fetchGlobalData = async () => {
      const data = await getGlobalData();
      setGlobalData(data);
      console.log("Global data loaded:", data);
      console.log("current DomainType", data.currentDomain.domainType);
    };

    fetchGlobalData();
  }, []);

  // Function to convert epoch time to readable date
  const epochConverter = (epochTime) => {
    if (!epochTime) return "N/A";
    const date = new Date(Number(epochTime) * 1000);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");

    return `${year}-${month}-${day}`;
  };

  // Function to fetch the users data with conditional API call
  const fetchUsersData = async () => {
    try {
      let data;

      // Check domainType and call API endpoint
      if (globalData.currentDomain?.domainType === "SERVICE_DOMAIN") {
        data = await userApi.getuser("LIST_PLATFORM_USERS");
      } else {
        data = await userApi.getuser("LIST_USERS");
      }

      return data.output?.users || [];
    } catch (error) {
      console.error("Error fetching model nodes data:", error);
      return [];
    }
  };

  const transformUsersData = (usersItems) => {
    return usersItems.map((item) => ({
      "User Id": item.userId || "N/A",
      "Invited On": epochConverter(item.invitedOn) || "N/A",
      "Display Name":
        item.displayName || item.firstName + " " + item.lastName || "N/A",
      Status: item.status || item.resourceBase?.status || "N/A",

      "View Details": `<a href="/settings/users/${item.userId}" target="_blank" class="relative group inline-flex items-center justify-center">
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
    }));
  };

  // For loading data - now depends on globalData being loaded
  useEffect(() => {
    const loadData = async () => {
      // Only load data if globalData is available
      if (!globalData.currentDomain) return;

      setIsLoading(true);
      const usersData = await fetchUsersData();
      console.log("Users data fetched:", usersData);
      const transformedData = transformUsersData(usersData);
      console.log("Transformed data:", transformedData);
      setUsers(transformedData);
      setIsLoading(false);
    };

    loadData();
  }, [globalData]); // Add globalData as dependency

  // Define columns for the DataTable
  const columns = [
    "User Id",
    "Invited On",
    "Display Name",
    "Status",
    "View Details",
  ];
  const filteredColumns = ["User Id", "Display Name"];

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto h-screen w-full">
        <Header
          sectionHeader="Domain Users"
          hideBackListingLink={true}
          backListingLink="./"
        />

        <div className="flex-grow p-2 w-full">
          <section id="tables" className="space-y-6">
            <div className="flex justify-between mb-2 items-center p-2"></div>

            {error ? (
              <div className="flex justify-center items-center h-64 text-red-400">
                <div className="text-xl">Error: {error}</div>
              </div>
            ) : users.length === 0 && !isLoading ? (
              <div className="overflow-x-auto scrollbar rounded-lg p-4 shadow-md text-gray-100">
                <table className="min-w-full divide-y divide-zinc-500 text-sm text-gray-100 border-2 border-zinc-500">
                  <thead className="bg-zinc-900 text-sm font-medium text-gray-500 uppercase tracking-wider"></thead>
                  <tbody className="bg-zinc-800 divide-y divide-zinc-500 break-words text-sm">
                    <tr>
                      <td
                        colSpan="5"
                        className="px-3 py-3 whitespace-nowrap text-center"
                      >
                        No Users found
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="p-4">
                <DataTable
                  id="usersDataTable"
                  data={users}
                  columns={columns}
                  loading={isLoading}
                  filteredColumns={filteredColumns}
                />
              </div>
            )}
          </section>
        </div>
      </div>
    </div>
  );
};

// Dynamic import with no SSR for the component
const UsersPage = dynamic(() => Promise.resolve(UsersTable), {
  ssr: false,
});

export default UsersPage;

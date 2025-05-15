"use client";

import React, { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import Sidebar from "@/app/components/platform/main-sidebar";
import Header from "@/app/components/platform/header";
import { userApi } from "@/app/utils/settings-endpoint/user-api";

// DataTable component with dynamic import
const DataTable = dynamic(() => import("@/app/components/table"), {
  ssr: false,
});

const UsersTable = () => {
  const [users, setUsers] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  // Function to convert epoch time to readable date
  const epochConverter = (epochTime) => {
    if (!epochTime) return "N/A";
    const date = new Date(Number(epochTime) * 1000);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0");
    const day = String(date.getDate()).padStart(2, "0");

    return `${year}-${month}-${day}`;
  };
  // Function to fetch the users data
  const fetchUsersData = async () => {
    try {
      const data = await userApi.getuser();
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
      "View Details": `<a href="users/${item.userId}" target="_blank" title="Open in new tab">
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
    }));
  };

  // For loading data
  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      const usersData = await fetchUsersData();
      console.log("Users data fetched:", usersData);
      const transformedData = transformUsersData(usersData);
      console.log("Transformed data:", transformedData);
      setUsers(transformedData);
      setIsLoading(false);
    };

    loadData();
  }, []);

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
      {/* Sidebar Component */}
      <Sidebar />

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

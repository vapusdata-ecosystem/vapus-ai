// ResourcePlatform.jsx
"use client";
import { useState, useEffect } from "react";
import Link from "next/link";
import Sidebar from "@/app/components/platform/main-sidebar";

export default function ResourcePlatform() {
  const [resources, setResources] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchResources() {
      try {
        const response = await fetch("/resources.json");
        const data = await response.json();

        // Assuming the JSON structure has a Response.SpecMap similar to the template data
        if (data.Response && data.Response.SpecMap) {
          const resourceEntries = Object.entries(data.Response.SpecMap).map(
            ([resource, spec]) => ({ resource, spec })
          );
          setResources(resourceEntries);
        } else {
          setResources([]);
          console.error("Invalid data format in resources.json");
        }
      } catch (error) {
        console.error("Error fetching resources:", error);
        setResources([]);
      } finally {
        setLoading(false);
      }
    }

    fetchResources();
  }, []);

  const totalPages = Math.ceil(resources.length / rowsPerPage);
  const startIndex = (currentPage - 1) * rowsPerPage;
  const endIndex = startIndex + rowsPerPage;
  const currentResources = resources.slice(startIndex, endIndex);

  const handlePrevPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  const handleNextPage = () => {
    if (currentPage < totalPages) {
      setCurrentPage(currentPage + 1);
    }
  };

  // Function to download YAML
  const downloadYaml = (spec, resourceName) => {
    const blob = new Blob([spec], { type: "text/yaml" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${resourceName}.yaml`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      {/* Assuming main-sidebar component would be imported separately */}
      <Sidebar />

      <div className="overflow-y-auto scrollbar h-screen w-full">
        {/* Header would be imported separately */}
        <header className="bg-zinc-900 p-4 text-white">
          <h1 className="text-xl">Platform Resources</h1>
        </header>

        <div className="flex-grow p-2 w-full">
          <div className="overflow-x-auto scrollbar rounded-lg p-4 shadow-md">
            {loading ? (
              <div className="text-center text-gray-100 py-4">
                Loading resources...
              </div>
            ) : (
              <table
                className="min-w-full divide-y divide-zinc-500 text-gray-100 border-2 border-zinc-500"
                id="platformResourcesTable"
              >
                <thead className="bg-zinc-900 font-medium text-gray-500 uppercase tracking-wider">
                  <tr>
                    <th className="px-3 py-3 text-left text-xs">
                      Resource Name
                    </th>
                    <th className="px-3 py-3 text-left text-xs">
                      Generate Yaml Spec
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200 text-sm">
                  {currentResources.map(({ resource, spec }) => (
                    <tr key={resource}>
                      <td className="px-3 py-3 whitespace-nowrap">
                        {resource}
                      </td>
                      <td className="px-1 py-1 whitespace-nowrap">
                        <button
                          onClick={() => downloadYaml(spec, resource)}
                          className="flex items-center px-2 py-2 hover:text-pink-900 text-orange-700 rounded-lg"
                        >
                          <svg
                            className="w-6 h-6 mr-2"
                            viewBox="0 0 24 24"
                            fill="currentColor"
                            stroke="currentColor"
                            xmlns="http://www.w3.org/2000/svg"
                          >
                            <path
                              d="M5 20h14v2H5v-2zm7-2c-.28 0-.53-.11-.71-.29L8 13.41l1.41-1.41L11 14.17V4h2v10.17l1.59-1.59L16 13.41l-3.29 3.29c-.18.18-.43.29-.71.29z"
                              fill="#000"
                            />
                          </svg>
                          Download YAML
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
                <tfoot>
                  <tr>
                    <td colSpan="2" className="px-3 py-2">
                      <div className="flex items-center justify-between">
                        <div className="text-sm text-gray-400">
                          Showing {startIndex + 1} to{" "}
                          {Math.min(endIndex, resources.length)} of{" "}
                          {resources.length} resources
                        </div>
                        <div className="flex items-center space-x-2">
                          <button
                            onClick={handlePrevPage}
                            disabled={currentPage === 1}
                            className={`px-3 py-1 rounded-md ${
                              currentPage === 1
                                ? "bg-zinc-700 text-gray-500"
                                : "bg-zinc-600 text-white hover:bg-zinc-500"
                            }`}
                          >
                            Previous
                          </button>
                          <span className="text-gray-300">
                            Page {currentPage} of {totalPages}
                          </span>
                          <button
                            onClick={handleNextPage}
                            disabled={currentPage === totalPages}
                            className={`px-3 py-1 rounded-md ${
                              currentPage === totalPages
                                ? "bg-zinc-700 text-gray-500"
                                : "bg-zinc-600 text-white hover:bg-zinc-500"
                            }`}
                          >
                            Next
                          </button>
                        </div>
                      </div>
                    </td>
                  </tr>
                </tfoot>
              </table>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

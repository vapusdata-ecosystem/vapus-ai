"use client";

import { useState, useEffect } from "react";
import { format } from "date-fns";
import ActionDropdown from "./action-dropdown";

const SectionHeaders = ({
  resourceId,
  resourceData = null,
  resourceType = "default",
  fetchUrl = null,
  isLoading = false,
  customActions = null,
  customButton = null,
  apiServices = {},
}) => {
  const [resource, setResource] = useState(null);
  const [loading, setLoading] = useState(isLoading || !resourceData);
  const [error, setError] = useState(null);
  const [accessToken, setAccessToken] = useState(null);

  // Get access token from localStorage when component mounts
  useEffect(() => {
    if (typeof window !== "undefined") {
      const token = localStorage.getItem("AccessTokenKey") || "";
      setAccessToken(token);
    }
  }, []);

  // Fetch resource if not provided directly
  useEffect(() => {
    const fetchResource = async () => {
      if (resourceData) {
        setResource({
          ...resourceData,
          actionRules: resourceData.actionRules || [],
          createActionParams: resourceData.createActionParams || null,
          yamlSpec: resourceData.yamlSpec || null,
          resourceType: resourceData.resourceType || resourceType,
        });
        setLoading(false);
        return;
      }

      if (!fetchUrl) {
        setError("No resource data or fetch URL provided");
        setLoading(false);
        return;
      }

      try {
        const response = await fetch(fetchUrl);
        if (!response.ok) {
          throw new Error(`Failed to fetch ${resourceType}`);
        }

        const data = await response.json();

        // Handle different data structures based on resourceType
        let foundResource = null;

        switch (resourceType) {
          case "prompt":
            if (data && data.output && Array.isArray(data.output)) {
              foundResource = data.output.find(
                (item) => item.promptId === resourceId
              );

              if (foundResource && foundResource.resourceBase) {
                setResource({
                  id: resourceId,
                  name: foundResource.name,
                  createdAt:
                    parseInt(foundResource.resourceBase.createdAt) * 1000,
                  createdBy: foundResource.resourceBase.createdBy,
                  status: foundResource.resourceBase.status,
                  resourceBase: foundResource.resourceBase,
                  // Add action-related fields
                  actionRules: foundResource.actionRules || [],
                  createActionParams: foundResource.createActionParams || null,
                  yamlSpec: foundResource.yamlSpec || null,
                  resourceType: "prompt",
                });
              }
            }
            break;

          case "model":
            if (data && data.output) {
              if (data.output.modelNodeId === resourceId) {
                foundResource = data.output;
              } else {
                const potentialArrays = [
                  data.output.generativeModels,
                  data.output.embeddingModels,
                ].filter((arr) => Array.isArray(arr));

                for (const arr of potentialArrays) {
                  foundResource = arr?.find(
                    (item) => item.modelNodeId === resourceId
                  );
                  if (foundResource) break;
                }
              }

              if (foundResource) {
                setResource({
                  id: resourceId,
                  name: foundResource.name,
                  createdAt: foundResource.resourceBase?.createdAt
                    ? parseInt(foundResource.resourceBase.createdAt) * 1000
                    : null,
                  createdBy: foundResource.resourceBase?.createdBy,
                  status: foundResource.resourceBase?.status,
                  resourceBase: foundResource.resourceBase,
                  actionRules: foundResource.actionRules || [],
                  createActionParams: foundResource.createActionParams || null,
                  yamlSpec: foundResource.yamlSpec || null,
                  resourceType: "model",
                });
              }
            }
            break;

          default:
            if (data) {
              const extractResource = (obj) => {
                if (!obj) return null;
                if (obj.id === resourceId || obj.resourceId === resourceId)
                  return obj;

                // Check if the object has arrays that might contain our resource
                for (const key in obj) {
                  if (Array.isArray(obj[key])) {
                    const found = obj[key].find(
                      (item) =>
                        item.id === resourceId ||
                        item.resourceId === resourceId ||
                        item.promptId === resourceId ||
                        item.modelNodeId === resourceId
                    );
                    if (found) return found;
                  }
                }

                return null;
              };

              foundResource =
                extractResource(data) || extractResource(data.output);

              if (foundResource) {
                setResource({
                  id: resourceId,
                  name: foundResource.name,
                  createdAt: foundResource.resourceBase?.createdAt
                    ? parseInt(foundResource.resourceBase.createdAt) * 1000
                    : foundResource.createdAt
                    ? parseInt(foundResource.createdAt) * 1000
                    : null,
                  createdBy:
                    foundResource.resourceBase?.createdBy ||
                    foundResource.createdBy,
                  status:
                    foundResource.resourceBase?.status || foundResource.status,
                  resourceBase: foundResource.resourceBase,
                  // Add action-related fields
                  actionRules: foundResource.actionRules || [],
                  createActionParams: foundResource.createActionParams || null,
                  yamlSpec: foundResource.yamlSpec || null,
                  resourceType: resourceType,
                });
              }
            }
        }

        if (!foundResource) {
          setError(`${resourceType} not found for ID: ${resourceId}`);
        }

        setLoading(false);
      } catch (error) {
        console.error(`Error fetching ${resourceType}:`, error);
        setError(error.message);
        setLoading(false);
      }
    };

    if (resourceId && !resourceData) {
      fetchResource();
    } else if (resourceData) {
      setResource({
        ...resourceData,
        actionRules: resourceData.actionRules || [],
        createActionParams: resourceData.createActionParams || null,
        yamlSpec: resourceData.yamlSpec || null,
        resourceType: resourceData.resourceType || resourceType,
      });
      setLoading(false);
    }
  }, [resourceId, resourceData, resourceType, fetchUrl]);

  // Format date helper function
  const formatDate = (timestamp) => {
    if (!timestamp) return "N/A";
    try {
      return format(new Date(timestamp), "yyyy-MM-dd");
    } catch (error) {
      console.error("Date formatting error:", error);
      return "Invalid date";
    }
  };

  if (loading) {
    return (
      <div className="text-gray-400 text-sm">Loading resource data...</div>
    );
  }

  if (error) {
    return <div className="text-red-400 text-sm">Error: {error}</div>;
  }

  if (!resource) {
    return (
      <div className="text-gray-400 text-sm">No resource data available</div>
    );
  }

  return (
    <div className="flex justify-between mt-2 items-center bg-[#1b1b1b] text-gray-100 rounded-lg shadow-md p-2">
      <div className="items-center">
        <h1 className="text-xl text-gray-100 font-semibold mb-2">
          {resource.name || "Unnamed Resource"}
        </h1>
        <div className="flex text-gray-100 mb-4 pr-4">
          <svg
            className="w-6 h-6 mr-2 ml-2"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M8 7V3m8 4V3m-9 4h10a2 2 0 012 2v11a2 2 0 01-2 2H7a2 2 0 01-2-2V9a2 2 0 012-2zm3 4h4m-4 4h4"
            />
          </svg>
          <span className="text-sm">
            <span className="border-b-2 text-gray-100 font-semibold text-sm border-gray-800">
              {formatDate(resource.createdAt)}
            </span>
          </span>
          <svg
            className="w-6 h-6 mr-2 ml-2"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"
            />
          </svg>
          <span className="text-sm">
            <span className="border-b-2 text-gray-100 text-sm border-gray-800">
              {resource.createdBy || "N/A"}
            </span>
          </span>

          {resource.status && (
            <p>
              <span
                className={`px-3 py-2 ml-2 text-sm font-medium ${
                  resource.status === "ACTIVE"
                    ? "text-green-800 bg-green-100"
                    : "text-red-800 bg-red-100"
                } rounded-full`}
              >
                {resource.status}
              </span>
            </p>
          )}
        </div>
      </div>

      <div className="relative inline-block text-left">
        {customActions || (
          <ActionDropdown
            response={{
              ...resource,
              resourceId: resource.id,
              resourceType: resource.resourceType || resourceType,
            }}
            globalContext={{ AccessTokenKey: accessToken }}
            customButton={customButton}
            apiServices={apiServices}
            resourceType={resource.resourceType || resourceType}
          />
        )}
      </div>
    </div>
  );
};

export default SectionHeaders;

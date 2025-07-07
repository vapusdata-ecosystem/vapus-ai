"use client";
import { useState, useEffect, use } from "react";
import Header from "@/app/components/platform/header";
import SectionHeaders from "@/app/components/section-headers";
import {
  modelsRegistryApi,
  modelsRegistryArchiveApi,
} from "@/app/utils/ai-studio-endpoint/models-registry-api";
import LoadingOverlay from "@/app/components/loading/loading";

export default function AIModelDetailsPage({ params }) {
  const resolvedParams = use(params);
  const modelNodeId = resolvedParams.modelNodeId;

  const [modelDetails, setModelDetails] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchModelDetails = async () => {
      try {
        const response = await modelsRegistryApi.getModelsRegistry(modelNodeId);
        if (!response.output) {
          console.error("Data does not contain output property:", response);
          setError("Unexpected data format received from server");
          setLoading(false);
          return;
        }

        if (response.output.modelNodeId === modelNodeId) {
          setModelDetails(response.output);
        } else {
          const modelArrays = Object.entries(response.output)
            .filter(([key, value]) => Array.isArray(value))
            .map(([key, value]) => ({ key, value }));

          if (modelArrays.length > 0) {
            let foundModel = null;

            for (const { key, value } of modelArrays) {
              foundModel = value.find(
                (item) => item.modelNodeId === modelNodeId
              );
              if (foundModel) break;
            }

            if (foundModel) {
              setModelDetails(foundModel);
            } else {
              setError(`AI Model not found for the given ID: ${modelNodeId}`);
              setModelDetails(null);
            }
          } else {
            setError(`AI Model not found for the given ID: ${modelNodeId}`);
            setModelDetails(null);
          }
        }

        setLoading(false);
      } catch (err) {
        console.error("Error fetching model details:", err);
        setError(err.message);
        setLoading(false);
      }
    };

    if (modelNodeId) {
      fetchModelDetails();
    } else {
      setError("No model node ID provided");
      setLoading(false);
    }
  }, [modelNodeId]);

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);

    // Create a simple toast notification
    const toast = document.createElement("div");
    toast.textContent = "Copied to clipboard";
    toast.style.position = "fixed";
    toast.style.bottom = "20px";
    toast.style.left = "50%";
    toast.style.transform = "translateX(-50%)";
    toast.style.backgroundColor = "rgba(0, 0, 0, 0.7)";
    toast.style.color = "white";
    toast.style.padding = "10px 20px";
    toast.style.borderRadius = "5px";
    toast.style.zIndex = "1000";

    document.body.appendChild(toast);

    setTimeout(() => {
      document.body.removeChild(toast);
    }, 1000);
  };

  // function to check if a string is empty or undefined
  const stringCheck = (str) => {
    return str ? str : "N/A";
  };

  if (loading) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center relative">
         <LoadingOverlay 
                        isLoading={loading} 
                        text="Loading plugin details"
                        size="default"
                        isOverlay={true}
                        className="absolute inset-0 z-10 bg-zinc-800"
                      />
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-red-500 text-xl">Error: {error}</div>
      </div>
    );
  }

  if (!modelDetails) {
    return (
      <div className="bg-zinc-800 flex h-screen justify-center items-center">
        <div className="text-white text-xl">AI Model not found</div>
      </div>
    );
  }
  // Prepare data for SectionHeaders component

  const apiServices = {
    model: {
      archive: modelsRegistryArchiveApi.getModelsRegistryArchive,
      delete: modelsRegistryArchiveApi.getModelsRegistryArchive,
    },
  };

  const headerResourceData = {
    id: modelDetails.modelNodeId,
    name: modelDetails.name || modelDetails.displayName,
    createdAt: modelDetails.resourceBase?.createdAt
      ? parseInt(modelDetails.resourceBase.createdAt) * 1000
      : null,
    createdBy: modelDetails.resourceBase?.createdBy,
    status: modelDetails.resourceBase?.status,
    resourceBase: modelDetails.resourceBase,
    resourceType: "model",
    actionRules: modelDetails.actionRules || [],
    createActionParams: modelDetails.createActionParams || {
      weblink: `./${modelDetails.modelNodeId}/update`,
    },
    // Add YAML spec for download button
    yamlSpec: modelDetails.yamlSpec || JSON.stringify(modelDetails, null, 2),
  };

  return (
    <div className="bg-zinc-800 flex h-screen ">
      <div className="overflow-y-auto scrollbar h-screen w-full ">
        <Header
          sectionHeader="AI Model Details"
          hideBackListingLink={false}
          backListingLink="./"
        />

        <div className="flex-grow p-2 w-full text-gray-100">
          <SectionHeaders
            resourceId={modelDetails.modelNodeId}
            resourceData={headerResourceData}
            resourceType="model"
            apiServices={apiServices}
          />

          <div className="overflow-x-auto scrollbar shadow-md p-8">
            <div
              id="spec"
              className="mt-2 bg-[#1b1b1b] p-4 rounded-lg shadow-md"
            >
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Name:
                  </p>
                  <p className="break-words p-2">
                    {stringCheck(modelDetails.name)}
                  </p>
                </div>

                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Id:
                  </p>
                  <p
                    className="break-words p-2 cursor-pointer"
                    onClick={() => copyToClipboard(modelDetails.modelNodeId)}
                  >
                    {stringCheck(modelDetails.modelNodeId)}
                  </p>
                </div>

                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Scope:
                  </p>
                  <p className="break-words p-2">
                    {stringCheck(modelDetails.attributes?.scope)}
                  </p>
                </div>

                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Endpoint:
                  </p>
                  <p className="break-words p-2">
                    {stringCheck(modelDetails.attributes?.networkParams?.url)}
                  </p>
                </div>

                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    API Version:
                  </p>
                  <p className="break-words p-2">
                    {stringCheck(
                      modelDetails.attributes?.networkParams?.apiVersion
                    )}
                  </p>
                </div>

                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Model Path:
                  </p>
                  <p className="break-words p-2">
                    {stringCheck(
                      modelDetails.attributes?.networkParams?.localPath
                    )}
                  </p>
                </div>

                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Model Discovery Enabled:
                  </p>
                  <p className="break-words p-2">
                    {modelDetails.attributes?.discoverModels?.toString() ||
                      "N/A"}
                  </p>
                </div>

                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Service Provider:
                  </p>
                  <p className="break-words p-2">
                    {stringCheck(modelDetails.attributes?.serviceProvider)}
                  </p>
                </div>

                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Hosting:
                  </p>
                  <p className="break-words p-2">
                    {stringCheck(modelDetails.attributes?.hosting)}
                  </p>
                </div>

                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Approved Domains:
                  </p>

                  <ul className="list-disc ml-5 break-words p-2">
                    {modelDetails.attributes?.approvedDomains &&
                    modelDetails.attributes.approvedDomains.length > 0 ? (
                      modelDetails.attributes.approvedDomains.map(
                        (domain, index) => <li key={index}>{domain}</li>
                      )
                    ) : (
                      <li>ALL</li>
                    )}
                  </ul>
                </div>

                <div className=" lg:flex items-center">
                  <p className="text-base font-extralight text-[#f4d1c2] block">
                    Node Owners:
                  </p>

                  <ul className="list-disc ml-5 break-words p-2">
                    {modelDetails.nodeOwners &&
                      modelDetails.nodeOwners.map((owner, index) => (
                        <li key={index}>{owner}</li>
                      ))}
                  </ul>
                </div>
              </div>

              {/* Security Section */}
              {modelDetails.securityGuardrails && (
                <>
                  <h3 className="text-xl font-bold mb-4 mt-6 text-[#f4d1c2] underline">
                    Security
                  </h3>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                      <p className="text-base font-extralight text-[#f4d1c2] block">
                        Guardrails:
                      </p>

                      <ul className="list-disc ml-5 break-words p-2">
                        {modelDetails.securityGuardrails.guardrails &&
                        modelDetails.securityGuardrails.guardrails.length >
                          0 ? (
                          modelDetails.securityGuardrails.guardrails.map(
                            (guardrail, index) => (
                              <li key={index}>{guardrail}</li>
                            )
                          )
                        ) : (
                          <li>ALL</li>
                        )}
                      </ul>
                    </div>
                  </div>
                </>
              )}

              {/* Supported Models Section */}
              <h3 className="text-xl font-bold mb-4 mt-6 text-[#f4d1c2] underline">
                Supported Models
              </h3>
              <div
                className="bg-zinc-800  rounded-lg shadow-md overflow-y-auto scrollbar text-sm"
                style={{ maxHeight: "300px" }}
              >
                <table className="min-w-full divide-y divide-zinc-500 ">
                  <thead className="bg-zinc-900 sticky top-0 ">
                    <tr>
                      <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                        Model Name
                      </th>
                      <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                        Model Id
                      </th>
                      <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                        Model Type
                      </th>
                      <th className="px-3 py-3 text-left text-xs font-medium uppercase tracking-wider">
                        Native Model Owner
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-zinc-800 divide-y divide-zinc-500 border-l border-r border-zinc-500">
                    {/* Generative Models */}
                    {modelDetails.attributes?.generativeModels &&
                      modelDetails.attributes.generativeModels.map(
                        (model, index) => (
                          <tr
                            key={`gen-${index}`}
                            className="cursor-pointer hover:bg-zinc-600"
                            onClick={() => copyToClipboard(model.modelName)}
                          >
                            <td className="px-3 py-3 whitespace-nowrap">
                              {model.modelName}
                            </td>
                            <td className="px-3 py-3 whitespace-nowrap">
                              {model.modelId}
                            </td>
                            <td className="px-3 py-3 whitespace-nowrap">
                              {model.modelType}
                            </td>
                            <td className="px-3 py-3 whitespace-nowrap">
                              {stringCheck(model.ownedBy)}
                            </td>
                          </tr>
                        )
                      )}

                    {/* Embedding Models */}
                    {modelDetails.attributes?.embeddingModels &&
                      modelDetails.attributes.embeddingModels.map(
                        (model, index) => (
                          <tr
                            key={`emb-${index}`}
                            className="cursor-pointer hover:bg-zinc-600"
                            onClick={() => copyToClipboard(model.modelName)}
                          >
                            <td className="px-3 py-3 whitespace-nowrap">
                              {model.modelName}
                            </td>
                            <td className="px-3 py-3 whitespace-nowrap">
                              {model.modelId}
                            </td>
                            <td className="px-3 py-3 whitespace-nowrap">
                              {model.modelType}
                            </td>
                            <td className="px-3 py-3 whitespace-nowrap">
                              {stringCheck(model.ownedBy)}
                            </td>
                          </tr>
                        )
                      )}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

"use client";
import React, { useState, useEffect } from "react";
import { use } from "react";
import Header from "@/app/components/platform/header";
import { modelsRegistryUpdateFormApi } from "@/app/utils/ai-studio-endpoint/models-registry-api";
import { modelsRegistryApi } from "@/app/utils/ai-studio-endpoint/models-registry-api";
import NetworkCredentials from "@/app/components/formcomponets/NetworkCredentials";
import ToastContainerMessage from "@/app/components/notification/customToast";
import { toast } from "react-toastify";
import LoadingOverlay from "@/app/components/loading/loading";
import YamlEditorClient from "@/app/components/formcomponets/ymal";
import { useRouter } from "next/navigation";
import { enumsApi } from "@/app/utils/developers-endpoint/enums";
import { strTitle, wordLimit } from "@/app/components/JS/common";
import { GuardrailApi } from "@/app/utils/ai-studio-endpoint/guardrails-api";

const CreateModelsNodesUpdate = ({ params }) => {
  const router = useRouter();

  const unwrappedParams = use(params);
  const ai_model_node_id = unwrappedParams?.modelNodeId
    ? String(unwrappedParams.modelNodeId).trim()
    : "";
  const [formMode, setFormMode] = useState("form"); // 'form' or 'yaml'
  const [showDropdownMenu, setShowDropdownMenu] = useState(false);
  const [showGuardrailsDropdown, setShowGuardrailsDropdown] = useState(false);
  const [guardrails, setGuardrails] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [aiModelNode, setAiModelNode] = useState(null);
  const [initialNetworkParams, setInitialNetworkParams] = useState(null);
  const [networkParams, setNetworkParams] = useState({});
  const [serviceProviderLogos, setServiceProviderLogos] = useState({});
  const [enums, setEnums] = useState({
    SvcProvider: [],
    ResourceScope: [],
  });
  const [formData, setFormData] = useState({
    name: "",
    node_owners: "",
    attributes: {
      service_provider: "INVALID_PROVIDER",
      scope: "DOMAIN_SCOPE",
      approved_domains: [],
      discover_models: true,
    },
    securityGuardrails: {
      guardrails: [],
    },
  });

  const getSortedProviders = (providers, selectedProvider) => {
    return [
      ...providers.filter((provider) => provider === selectedProvider),
      ...providers.filter((provider) => provider !== selectedProvider),
    ];
  };

  // Fetch enums data
  useEffect(() => {
    const fetchEnumsData = async () => {
      try {
        setIsLoading(true);
        const response = await enumsApi.getEnums();
        const guardrailData = await GuardrailApi.getGuardrail();

        if (guardrailData && guardrailData.output) {
          setGuardrails(guardrailData.output);
          console.log("Guardrails loaded:", guardrailData.output);
        }

        // Process enum responses for resource scopes
        const newEnums = { SvcProvider: [], ResourceScope: [] };
        const enumResponses = response.enumResponse || [];

        enumResponses.forEach((enumData) => {
          if (enumData.name === "ResourceScope") {
            newEnums.ResourceScope = enumData.value || [];
          }
        });
        const serviceProviders = response.serviceProviderLogoMap || [];
        newEnums.SvcProvider = serviceProviders.map(
          (item) => item.serviceProvider
        );

        setEnums(newEnums);
        console.log("Resource scopes loaded:", newEnums.ResourceScope);
        console.log("Service providers loaded:", newEnums.SvcProvider);

        const logoMapObject = {};
        serviceProviders.forEach((item) => {
          logoMapObject[item.serviceProvider] = item.url;
        });

        setServiceProviderLogos(logoMapObject);
        console.log("Service provider logos loaded:", logoMapObject);
      } catch (error) {
        console.error("Failed to fetch enum data:", error);
        toast.error("Failed to fetch configuration data");
      } finally {
        setIsLoading(false);
      }
    };

    fetchEnumsData();
  }, []);

  // Fetch initial data on component mount
  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);
      try {
        console.log("Fetching models registry data", ai_model_node_id);
        const data = await modelsRegistryApi.getModelsRegistryID(
          ai_model_node_id
        );
        console.log("Data Fetched:", data);
        setAiModelNode(data);
        initializeFormData(data);
      } catch (error) {
        console.error("Error fetching model data:", error);
      } finally {
        setIsLoading(false);
      }
    };
    if (ai_model_node_id) {
      fetchData();
    } else {
      console.error("No model ID provided");
      toast.error("No model ID provided for update");
    }
  }, [ai_model_node_id]);

  // Initialize form data with fetched model data
  const initializeFormData = (data) => {
    if (
      data &&
      data.output &&
      data.output.aiModelNodes &&
      data.output.aiModelNodes.length > 0
    ) {
      const modelNode = data.output.aiModelNodes[0];
      const networkParamsData = modelNode.attributes.networkParams || {};
      setInitialNetworkParams(networkParamsData);
      setNetworkParams(networkParamsData);
      const nodeOwnersStr = modelNode.nodeOwners
        ? modelNode.nodeOwners.filter((owner) => owner.trim() !== "").join(",")
        : "";

      // Extract guardrails if present
      const existingGuardrails = modelNode.securityGuardrails?.guardrails || [];

      // Set form data
      setFormData({
        name: modelNode.name || "",
        node_owners: nodeOwnersStr,
        attributes: {
          service_provider:
            modelNode.attributes.serviceProvider || "INVALID_PROVIDER",
          scope: modelNode.attributes.scope || "DOMAIN_SCOPE",
          approved_domains: modelNode.attributes.approvedDomains || [],
          discover_models: true,
        },
        securityGuardrails: {
          guardrails: existingGuardrails.map((id) => ({
            id,
            name: getGuardrailNameById(id) || "Unknown Guardrail",
          })),
        },
      });
    }
  };

  //  function to get guardrail name from id
  const getGuardrailNameById = (id) => {
    const found = guardrails.find((g) => g.guardrailId === id);
    return found ? found.name : null;
  };

  // dummy domain data
  const userDomainRoles = [
    { DomainId: "domain1", name: "ANAND.VAPUSDATA.COM domain" },
  ];

  const domainMap = {
    domain1: "ANAND.VAPUSDATA.COM domain",
  };

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    const nameParts = name.split(".");

    if (nameParts[0] === "spec") {
      nameParts.shift();
    }

    setFormData((prevData) => {
      const newData = { ...prevData };
      let current = newData;

      // Navigate to the correct nested object
      for (let i = 0; i < nameParts.length - 1; i++) {
        if (!current[nameParts[i]]) {
          current[nameParts[i]] = {};
        }
        current = current[nameParts[i]];
      }

      // Set the value
      const lastKey = nameParts[nameParts.length - 1];
      if (type === "checkbox") {
        current[lastKey] = checked;
      } else {
        current[lastKey] = value;
      }

      return newData;
    });
  };

  const handleDomainCheckboxChange = (domainId) => {
    setFormData((prevData) => {
      const domains = [...prevData.attributes.approved_domains];
      const index = domains.indexOf(domainId);

      if (index === -1) {
        domains.push(domainId);
      } else {
        domains.splice(index, 1);
      }

      return {
        ...prevData,
        attributes: {
          ...prevData.attributes,
          approved_domains: domains,
        },
      };
    });
  };

  const handleGuardrailCheckboxChange = (guardrailId, guardrailName) => {
    setFormData((prevData) => {
      const selectedGuardrails = [
        ...(prevData.securityGuardrails.guardrails || []),
      ];
      const index = selectedGuardrails.findIndex((g) => g.id === guardrailId);

      if (index === -1) {
        selectedGuardrails.push({
          id: guardrailId,
          name: guardrailName,
        });
      } else {
        selectedGuardrails.splice(index, 1);
      }
      return {
        ...prevData,
        securityGuardrails: {
          ...prevData.securityGuardrails,
          guardrails: selectedGuardrails,
        },
      };
    });
  };

  const handleNetworkParamsChange = (params) => {
    if (JSON.stringify(networkParams) !== JSON.stringify(params)) {
      setNetworkParams(params);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Combine formData with networkParams
    let dataObj = {
      ...formData,
      attributes: {
        ...formData.attributes,
        networkParams: networkParams,
      },
    };

    if (dataObj.node_owners !== "") {
      dataObj.node_owners = dataObj.node_owners.split(",");
    } else {
      dataObj.node_owners = [];
    }

    // Extract guardrail IDs for API payload
    const guardrailIds = dataObj.securityGuardrails.guardrails.map((g) => g.id);

    // For API payload
    const payload = {
      spec: {
        name: dataObj.name,
        // model ID
        ...(aiModelNode &&
          aiModelNode.output &&
          aiModelNode.output.aiModelNodes &&
          aiModelNode.output.aiModelNodes[0] && {
            modelNodeId: aiModelNode.output.aiModelNodes[0].modelNodeId,
          }),
        node_owners: dataObj.node_owners,
        attributes: {
          service_provider: dataObj.attributes.service_provider,
          scope: dataObj.attributes.scope,
          approved_domains: dataObj.attributes.approved_domains,
          discover_models: dataObj.attributes.discover_models,
          networkParams: dataObj.attributes.networkParams,
        },
        securityGuardrails: {
          guardrails: guardrailIds,
        },
      },
    };

    // Submit form
    await submitForm(payload);
  };

  const submitForm = async (payload) => {
    try {
      setIsLoading(true);
      const output =
        await modelsRegistryUpdateFormApi.getmodelsRegistryUpdateForm(payload);
      let modelNodeId;
      if (
        aiModelNode &&
        aiModelNode.output &&
        aiModelNode.output.aiModelNodes &&
        aiModelNode.output.aiModelNodes[0]
      ) {
        modelNodeId = aiModelNode.output.aiModelNodes[0].modelNodeId;
      } else if (payload.spec.modelNodeId) {
        modelNodeId = payload.spec.modelNodeId;
      }

      toast.success("Resource updated successfully");

      setTimeout(() => {
        if (modelNodeId) {
          router.push(`/ai-center/models-registry/${modelNodeId}`);
        }
      }, 1000);
    } catch (error) {
      console.error("Error sending API request:", error);
      toast.error("Resource Update Failed");
    } finally {
      setIsLoading(false);
    }
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      const dropdown = document.getElementById("dropdownMenu");
      const button = document.getElementById("dropdownButton");

      if (
        dropdown &&
        button &&
        !dropdown.contains(event.target) &&
        !button.contains(event.target)
      ) {
        setShowDropdownMenu(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  // Close guardrails dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      const dropdown = document.getElementById("guardrailsDropdownMenu");
      const button = document.getElementById("guardrailsDropdownButton");

      if (
        dropdown &&
        button &&
        !dropdown.contains(event.target) &&
        !button.contains(event.target)
      ) {
        setShowGuardrailsDropdown(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Update AI Model Node"
          hideBackListingLink={false}
          backListingLink="./"
        />

        <ToastContainerMessage />
        <LoadingOverlay isLoading={isLoading} />

        <div className="flex-grow p-4 overflow-y-auto w-full">
          <section className="space-y-2">
            <div className="max-w-6xl mx-auto bg-[#1b1b1b] shadow rounded-lg p-2">
              {/* TOGGLE BUTTON BETWEEN YAML AND FORM */}
              <div className="text-gray-100 mb-2 flex justify-center">
                <button
                  id="formSpecB"
                  className="whitespace-nowrap border-b-2 py-2 px-2 text-md font-medium focus:outline-none ml-4 border-orange-700 text-orange-700 font-semibold"
                  onClick={() => setFormMode("form")}
                  suppressHydrationWarning
                >
                  Form
                </button>
              </div>

              {/* YAML FORM  */}
              {formMode === "yaml" ? (
                <div id="yamlSpec">{/* <YamlEditorClient /> */}</div>
              ) : (
                <div id="formSpec">
                  <form
                    id="dataSourceSpec"
                    className="space-y-2 border border-zinc-500 rounded-md text-gray-100 p-2"
                    onSubmit={handleSubmit}
                  >
                    <fieldset className="p-4 rounded">
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                        {/* Name */}
                        <div>
                          <label className="labels">Name</label>
                          <input
                            id="spec_name"
                            name="name"
                            type="text"
                            placeholder="Enter name"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                            value={formData.name}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                        {/* Node Owner */}
                        <div>
                          <label className="labels">Node Owners</label>
                          <input
                            id="spec_nodeOwners"
                            name="node_owners"
                            type="text"
                            placeholder="Enter node owners (comma separated)"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                            value={formData.node_owners}
                            onChange={handleInputChange}
                            suppressHydrationWarning
                          />
                        </div>
                      </div>

                      <fieldset className="rounded mb-4">
                        <legend className="text-lg font-semibold">
                          Attributes
                        </legend>
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                          {/* Scope */}
                          <div>
                            <label className="labels">Scope</label>
                            <select
                              name="attributes.scope"
                              className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                              value={formData.attributes.scope}
                              onChange={handleInputChange}
                              suppressHydrationWarning
                            >
                              <option value=" ">Select Scope</option>
                              {enums.ResourceScope?.map((scope) => (
                                <option key={scope} value={scope}>
                                  {strTitle(scope)}
                                </option>
                              ))}
                            </select>
                          </div>

                          {/* Select Domains - Only show when scope is "platform" */}
                          {formData.attributes.scope === "PLATFORM_SCOPE" && (
                            <div className="relative inline-block text-left w-full">
                              <div>
                                <label className="labels">Domains</label>
                                <button
                                  type="button"
                                  id="dropdownButton"
                                  className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                  onClick={() =>
                                    setShowDropdownMenu(!showDropdownMenu)
                                  }
                                  suppressHydrationWarning
                                >
                                  Select Domains
                                  <svg
                                    className="w-5 h-5 text-gray-500"
                                    xmlns="http://www.w3.org/2000/svg"
                                    viewBox="0 0 20 20"
                                    fill="currentColor"
                                  >
                                    <path
                                      fillRule="evenodd"
                                      d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                                      clipRule="evenodd"
                                    />
                                  </svg>
                                </button>
                              </div>

                              <div
                                id="dropdownMenu"
                                className={`absolute z-10 mt-2 w-full bg-zinc-800 rounded-md shadow-lg ${
                                  showDropdownMenu ? "" : "hidden"
                                }`}
                              >
                                <div className="p-2 space-y-2 max-h-60 overflow-y-auto">
                                  {userDomainRoles.map((domain) => (
                                    <div
                                      key={domain.DomainId}
                                      className="flex items-center"
                                    >
                                      <input
                                        type="checkbox"
                                        id={`domain-${domain.DomainId}`}
                                        value={domain.DomainId}
                                        checked={formData.attributes.approved_domains.includes(
                                          domain.DomainId
                                        )}
                                        onChange={() =>
                                          handleDomainCheckboxChange(
                                            domain.DomainId
                                          )
                                        }
                                        className="h-4 w-4 text-orange-700 border-gray-300 rounded focus:ring-orange-700 accent-orange-700"
                                      />
                                      <label
                                        htmlFor={`domain-${domain.DomainId}`}
                                        className="ml-2 text-sm"
                                      >
                                        {domainMap[domain.DomainId]}
                                      </label>
                                    </div>
                                  ))}
                                </div>
                              </div>
                            </div>
                          )}

                          {/* Select Guardrails */}
                          <div className="relative inline-block text-left w-full">
                            <div>
                              <label className="labels">Guardrails</label>
                              <button
                                type="button"
                                id="guardrailsDropdownButton"
                                className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                                onClick={() =>
                                  setShowGuardrailsDropdown(
                                    !showGuardrailsDropdown
                                  )
                                }
                                suppressHydrationWarning
                              >
                                Select Guardrails
                                <svg
                                  className="w-5 h-5 text-gray-500"
                                  xmlns="http://www.w3.org/2000/svg"
                                  viewBox="0 0 20 20"
                                  fill="currentColor"
                                >
                                  <path
                                    fillRule="evenodd"
                                    d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                                    clipRule="evenodd"
                                  />
                                </svg>
                              </button>
                            </div>

                            <div
                              id="guardrailsDropdownMenu"
                              className={`absolute z-10 mt-2 w-full bg-zinc-800 rounded-md shadow-lg ${
                                showGuardrailsDropdown ? "" : "hidden"
                              }`}
                            >
                              <div className="p-2 space-y-2 max-h-60 overflow-y-auto">
                                {guardrails && guardrails.length > 0 ? (
                                  guardrails.map((guardrail) => (
                                    <div
                                      key={guardrail.guardrailId}
                                      className="flex items-center"
                                    >
                                      <input
                                        type="checkbox"
                                        id={`guardrail-${guardrail.guardrailId}`}
                                        value={guardrail.guardrailId}
                                        checked={
                                          formData.securityGuardrails &&
                                          formData.securityGuardrails
                                            .guardrails &&
                                          formData.securityGuardrails.guardrails.some(
                                            (g) =>
                                              g.id === guardrail.guardrailId
                                          )
                                        }
                                        onChange={() =>
                                          handleGuardrailCheckboxChange(
                                            guardrail.guardrailId,
                                            guardrail.name
                                          )
                                        }
                                        className="h-4 w-4 text-orange-700 border-gray-300 rounded focus:ring-orange-700 accent-orange-700"
                                      />
                                      <label
                                        htmlFor={`guardrail-${guardrail.guardrailId}`}
                                        className="ml-2 text-sm"
                                      >
                                        {guardrail.name}
                                      </label>
                                    </div>
                                  ))
                                ) : (
                                  <div className="text-sm text-gray-400 p-2">
                                    No guardrails available
                                  </div>
                                )}
                              </div>
                            </div>
                          </div>
                        </div>
                        {/* >Service Provider */}
                        <div className="m-4">
                          <label className="labels">Service Provider</label>
                          <div className="mt-2">
                            <div className="rounded mb-4 max-h-36 scrollbar overflow-y-auto">
                              <div className="grid grid-cols-3 sm:grid-cols-5 md:grid-cols-8 lg:grid-cols-10 gap-4">
                                {getSortedProviders(
                                  enums.SvcProvider,
                                  formData.attributes.service_provider
                                ).map((provider) => (
                                  <div
                                    key={provider}
                                    className={`border rounded-md p-3 cursor-pointer transition-all duration-200 flex flex-col items-center ${
                                      formData.attributes.service_provider ===
                                      provider
                                        ? "border-orange-700 bg-zinc-700"
                                        : "border-zinc-600 hover:border-orange-700 hover:bg-zinc-700"
                                    }`}
                                    onClick={() => {
                                      const event = {
                                        target: {
                                          name: "attributes.service_provider",
                                          value: provider,
                                          type: "select",
                                        },
                                      };
                                      handleInputChange(event);
                                    }}
                                  >
                                    {/* Display image with text */}
                                    <div className="flex items-center gap-2">
                                      <div className="h-6 w-6 flex items-center justify-center">
                                        {serviceProviderLogos[provider] ? (
                                          <img
                                            src={serviceProviderLogos[provider]}
                                            alt={provider}
                                            className="max-h-full max-w-full object-contain"
                                          />
                                        ) : (
                                          <div className="h-8 w-8 bg-zinc-600 rounded-full flex items-center justify-center text-xs">
                                            {provider.substring(0, 2)}
                                          </div>
                                        )}
                                      </div>
                                      <div
                                        className="tooltip text-[10px] font-medium "
                                        tooltip={strTitle(provider)}
                                      >
                                        {wordLimit(strTitle(provider), 5)}
                                      </div>
                                    </div>
                                  </div>
                                ))}
                              </div>
                            </div>

                            {/* Hidden input for service provider value */}
                            <input
                              type="hidden"
                              name="attributes.service_provider"
                              value={formData.attributes.service_provider}
                            />
                          </div>
                        </div>
                      </fieldset>

                      {/* NetworkParams Credentials Component */}
                      <details
                        className="border border-zinc-500 p-4 rounded mb-4"
                        open
                      >
                        <summary className="text-lg font-semibold cursor-pointer">
                          Network Credentials
                        </summary>
                        <NetworkCredentials
                          onParamsChange={handleNetworkParamsChange}
                          initialParams={initialNetworkParams}
                        />
                      </details>
                    </fieldset>

                    {/* Submit Section */}
                    <div className="mt-4 flex justify-end space-x-2">
                      <button
                        type="submit"
                        className="px-6 py-2 bg-orange-700 text-white rounded-md shadow hover:bg-pink-900 cursor-pointer"
                        suppressHydrationWarning
                        disabled={isLoading}
                      >
                        {isLoading ? (
                          <svg
                            className="animate-spin h-6 w-6 text-white"
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                          >
                            <circle
                              className="opacity-25"
                              cx="12"
                              cy="12"
                              r="10"
                              stroke="currentColor"
                              strokeWidth="4"
                            ></circle>
                            <path
                              className="opacity-75"
                              fill="currentColor"
                              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                            ></path>
                          </svg>
                        ) : (
                          "Update"
                        )}
                      </button>
                    </div>
                  </form>
                </div>
              )}
            </div>
          </section>
        </div>
      </div>
    </div>
  );
};

export default CreateModelsNodesUpdate;

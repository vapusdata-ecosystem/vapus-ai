"use client";
import React, { useState, useEffect } from "react";
import { use } from "react";
import Header from "@/app/components/platform/header";
import NetworkCredentials from "@/app/components/formcomponets/NetworkCredentials";
import ToastContainerMessage from "@/app/components/notification/customToast";
import { strTitle, wordLimit } from "@/app/components/JS/common";
import LoadingOverlay from "@/app/components/loading/loading";
import {
  pluginsUpdateFormApi,
  pluginsApi,
} from "@/app/utils/settings-endpoint/plugins-api";
import { toast } from "react-toastify";
import { enumsApi } from "@/app/utils/developers-endpoint/enums";
import AddButton from "@/app/components/buttons/addButton";
import RemoveButton from "@/app/components/buttons/removeButton";
import { useRouter } from "next/navigation";

const UpdatePluginsForm = ({ params }) => {
  const router = useRouter();
  const unwrappedParams = use(params);
  const pluginId = unwrappedParams?.id ? String(unwrappedParams.id).trim() : "";

  const [formMode, setFormMode] = useState("form"); // 'form' or 'yaml'
  const [isLoading, setIsLoading] = useState(true);
  const [networkParams, setNetworkParams] = useState({});
  const [initialNetworkParams, setInitialNetworkParams] = useState(null);
  const [pluginTypesOptions, setPluginTypesOptions] = useState([]);
  const [resourceScopeOptions, setResourceScopeOptions] = useState([]);
  const [serviceProvidersOptions, setServiceProvidersOptions] = useState({});
  const [pluginTypeMap, setPluginTypeMap] = useState([]);
  const [dynamicParams, setDynamicParams] = useState([{ key: "", value: "" }]);
  const [formData, setFormData] = useState({
    name: "",
    pluginType: "",
    pluginService: "",
    scope: "",
    editable: true,
  });
  const [error, setError] = useState(null);
  const readonlyMode = true;

  useEffect(() => {
    const fetchPluginData = async () => {
      try {
        setIsLoading(true);
        const enumsData = await enumsApi.getEnums();
        const enumResponses = enumsData.enumResponse || [];
        const resourceScopeEnum = enumResponses.find(
          (item) => item.name === "ResourceScope"
        );
        if (resourceScopeEnum && resourceScopeEnum.value) {
          setResourceScopeOptions(resourceScopeEnum.value);
        }
        const pluginTypeMapData = enumsData.pluginTypeMap || [];
        setPluginTypeMap(pluginTypeMapData);

        const pluginTypes = pluginTypeMapData.map((item) => item.pluginTypes);
        setPluginTypesOptions(pluginTypes);

        // Check  pluginId
        if (pluginId) {
          const response = await pluginsApi.getPluginsId(pluginId);
          const pluginData = response.output?.[0] || {};
          if (pluginData) {
            const updatedFormData = {
              name: pluginData.name || "",
              pluginType: pluginData.pluginType || "",
              pluginService: pluginData.pluginService || "",
              scope: pluginData.scope || "DOMAIN_SCOPE",
              editable: pluginData.editable !== false,
            };

            setFormData(updatedFormData);

            // Set network params and initial network params
            if (pluginData.networkParams) {
              setNetworkParams(pluginData.networkParams);
              setInitialNetworkParams(pluginData.networkParams);
            }

            // Set dynamic params
            if (
              pluginData.dynamicParams &&
              pluginData.dynamicParams.length > 0
            ) {
              setDynamicParams(pluginData.dynamicParams);
            }
            if (pluginData.pluginType) {
              const selectedPluginType = pluginTypeMapData.find(
                (item) => item.pluginTypes === pluginData.pluginType
              );
              if (selectedPluginType && selectedPluginType.services) {
                setServiceProvidersOptions(selectedPluginType.services);
              }
            }
          }
        } else {
          console.error("No plugin ID provided");
          toast.error("No plugin ID provided for update");
        }
      } catch (error) {
        console.error("Failed to fetch plugin data:", error);
        setError(error.message);
        toast.error("Failed to fetch plugin configuration data");
      } finally {
        setIsLoading(false);
      }
    };

    fetchPluginData();
  }, [pluginId]);

  // Updated service provider options when plugin type changes
  useEffect(() => {
    if (formData.pluginType !== "") {
      const selectedPluginType = pluginTypeMap.find(
        (item) => item.pluginTypes === formData.pluginType
      );

      if (selectedPluginType && selectedPluginType.services) {
        setServiceProvidersOptions(selectedPluginType.services);

        const serviceProviderExists = Object.keys(
          selectedPluginType.services
        ).includes(formData.pluginService);

        if (!serviceProviderExists) {
          setFormData((prevData) => ({
            ...prevData,
            pluginService: "",
          }));
        }
      } else {
        setServiceProvidersOptions({});
      }
    } else {
      setServiceProvidersOptions({});
    }
  }, [formData.pluginType, pluginTypeMap]);

  const handleInputChange = (e) => {
    const { name, value, type, checked } = e.target;
    const nameParts = name.split(".");

    if (nameParts[0] === "spec") {
      nameParts.shift();
    }

    setFormData((prevData) => {
      const newData = { ...prevData };
      let current = newData;

      if (nameParts.length === 1) {
        if (type === "checkbox") {
          current[nameParts[0]] = checked;
        } else {
          current[nameParts[0]] = value;
        }
        return newData;
      }

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

  const handleNetworkParamsChange = (params) => {
    if (JSON.stringify(networkParams) !== JSON.stringify(params)) {
      setNetworkParams(params);
    }
  };

  const addDynamicParams = () => {
    setDynamicParams([...dynamicParams, { key: "", value: "" }]);
  };

  const removeDynamicParams = (index) => {
    const newParams = [...dynamicParams];
    newParams.splice(index, 1);
    setDynamicParams(newParams);
  };

  const handleDynamicParamChange = (index, field, value) => {
    const updatedParams = [...dynamicParams];
    updatedParams[index][field] = value;
    setDynamicParams(updatedParams);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Filter out any dynamic params with empty keys
    const filteredDynamicParams = dynamicParams.filter(
      (param) => param.key.trim() !== ""
    );

    // Combine formData with networkParams and dynamicParams
    let combinedFormData = {
      ...formData,
      networkParams: networkParams,
      dynamicParams: filteredDynamicParams,
    };

    if (
      combinedFormData.networkParams?.credentials &&
      Object.keys(combinedFormData.networkParams.credentials).length > 0
    ) {
      const genericCreds = combinedFormData.networkParams.credentials || {};
      const plugin_secrets = JSON.stringify(genericCreds);
      combinedFormData.data = plugin_secrets;
    }

    const pluginService = combinedFormData.pluginService || "";

    // payload
    const payload = {
      spec: {
        pluginId: pluginId,
        name: combinedFormData.name,
        scope: combinedFormData.scope,
        pluginType: combinedFormData.pluginType,
        pluginService: pluginService,
        editable: combinedFormData.editable,
        networkParams: combinedFormData.networkParams,
        dynamicParams: combinedFormData.dynamicParams,
      },
    };
    await submitForm(payload);
  };

  const submitForm = async (payload) => {
    try {
      setIsLoading(true);
      const output = await pluginsUpdateFormApi.getPluginsUpdateForm(
        payload,
        pluginId
      );

      const resourceInfo = output.result;
      if (resourceInfo) {
        toast.success(
          "Plugin Updated",
          `${resourceInfo.resource} Plugin updated successfully.`
        );
        setTimeout(() => {
          router.push(`/settings/plugins/${pluginId}`);
        }, 1000);
      } else {
        toast.success("Plugin Updated", "Plugin updated successfully.");
        setTimeout(() => {
          router.push(`/settings/plugins/${pluginId}`);
        }, 1000);
      }
    } catch (error) {
      console.error("Error sending API request:", error);
      toast.error("Plugin Update Failed");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Update Plugin"
          hideBackListingLink={false}
          backListingLink="./"
        />

        <ToastContainerMessage />
        <LoadingOverlay isLoading={isLoading} isOverlay={true}/>

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

              {/* FORM */}
              <div id="formSpec">
                <form
                  id="pluginForm"
                  className="space-y-2 border border-zinc-500 rounded-md text-gray-100 p-2"
                  onSubmit={handleSubmit}
                >
                  <fieldset className="p-4 rounded">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                      {/* Name */}
                      <div>
                        <label className="labels">Name</label>
                        <input
                          id="name"
                          name="name"
                          type="text"
                          placeholder="Enter plugin name"
                          className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                          value={formData.name}
                          onChange={handleInputChange}
                          required
                          suppressHydrationWarning
                        />
                      </div>

                      {/* Scope */}
                      <div>
                        <label className="labels">Scope</label>
                        <select
                          name="scope"
                          className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                          value={formData.scope}
                          onChange={handleInputChange}
                          suppressHydrationWarning
                        >
                          <option value="">Select Scope</option>
                          {resourceScopeOptions?.map((scope) => (
                            <option key={scope} value={scope}>
                              {strTitle(scope)}
                            </option>
                          ))}
                        </select>
                      </div>

                      {/* Plugin Type */}
                      <div>
                        <label className="labels">Plugin Type</label>
                        <select
                          name="pluginType"
                          className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                          value={formData.pluginType}
                          onChange={handleInputChange}
                          suppressHydrationWarning
                          disabled
                        >
                          <option value="">Select Plugin Type</option>
                          {pluginTypesOptions?.map((type) => (
                            <option key={type} value={type}>
                              {strTitle(type)}
                            </option>
                          ))}
                        </select>
                      </div>
                    </div>
                    {/* Service Provider */}
                    <div className="mb-4">
                      <label className="labels">Service Provider</label>
                      <div className="mt-2">
                        {!formData.pluginType ? (
                          <div className="text-gray-400 italic">
                            Please select a Plugin Type first to see available
                            Service Providers
                          </div>
                        ) : (
                          <>
                            <div className="rounded mb-4 max-h-36 scrollbar overflow-y-auto">
                              <div className="grid grid-cols-4 sm:grid-cols-4 md:grid-cols-6 lg:grid-cols-6 gap-4">
                                {Object.keys(serviceProvidersOptions).map(
                                  (provider) => (
                                    <div
                                      key={provider}
                                      className={`border rounded-md p-3 cursor-pointer transition-all duration-200 flex flex-col items-center ${
                                        formData.pluginService === provider
                                          ? "border-orange-700 bg-zinc-700"
                                          : "border-zinc-600 hover:border-orange-700 hover:bg-zinc-700"
                                      }`}
                                      onClick={() => {
                                        if (readonlyMode) return;
                                        const event = {
                                          target: {
                                            name: "pluginService",
                                            value: provider,
                                            type: "select",
                                          },
                                        };
                                        handleInputChange(event);
                                      }}
                                    >
                                      {/* display the image with Text  */}
                                      <div className="flex items-center gap-2">
                                        <div className="h-8 w-8 flex items-center justify-center">
                                          <img
                                            src={
                                              serviceProvidersOptions[provider]
                                            }
                                            alt={provider}
                                            className="max-h-full max-w-full object-contain"
                                          />
                                        </div>
                                        <div
                                          className="tooltip text-sm font-medium"
                                          tooltip={strTitle(provider)}
                                        >
                                          {wordLimit(strTitle(provider), 7)}
                                        </div>
                                      </div>
                                    </div>
                                  )
                                )}
                              </div>
                            </div>

                            {/* input for send the Service Provider while submiting the form  */}
                            <input
                              type="hidden"
                              name="pluginService"
                              value={formData.pluginService}
                            />
                          </>
                        )}
                      </div>
                    </div>

                    {/*  networkParams Credentials Component */}
                    <details className="border border-zinc-500 p-4 rounded mb-4">
                      <summary className="text-lg font-semibold cursor-pointer">
                        Network Credentials
                      </summary>
                      <NetworkCredentials
                        onParamsChange={handleNetworkParamsChange}
                        initialParams={initialNetworkParams}
                      />
                    </details>

                    {/* Dynamic  dynamicParams */}
                    <details className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                      <summary className="text-lg font-semibold text-gray-100 cursor-pointer">
                        Dynamic Parameters
                      </summary>
                      <fieldset className="rounded mb-4">
                        <div id="dynamicParamsContainer">
                          {dynamicParams.map((param, index) => (
                            <div
                              key={`param-${index}`}
                              className="border p-3 rounded mb-3"
                            >
                              <div className="grid grid-cols-2 md:grid-cols-2 gap-4 mb-4">
                                <div>
                                  <label className="labels">Key</label>
                                  <input
                                    type="text"
                                    placeholder="Enter Key"
                                    className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                    value={param.key}
                                    onChange={(e) =>
                                      handleDynamicParamChange(
                                        index,
                                        "key",
                                        e.target.value
                                      )
                                    }
                                    suppressHydrationWarning
                                  />
                                </div>
                                <div>
                                  <label className="labels">Value</label>
                                  <input
                                    type="text"
                                    placeholder="Enter Value"
                                    className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                                    value={param.value}
                                    onChange={(e) =>
                                      handleDynamicParamChange(
                                        index,
                                        "value",
                                        e.target.value
                                      )
                                    }
                                    suppressHydrationWarning
                                  />
                                </div>
                              </div>
                              {index > 0 && (
                                <RemoveButton
                                  onClick={() => removeDynamicParams(index)}
                                />
                              )}
                            </div>
                          ))}
                        </div>
                        <AddButton
                          name="+ Add Parameter"
                          onClick={addDynamicParams}
                        />
                      </fieldset>
                    </details>
                  </fieldset>

                  {/* Submit Section */}
                  <div className="mt-4 flex justify-end space-x-2">
                    <button
                      type="submit"
                      className="px-6 py-2 bg-orange-700 text-white rounded-md shadow hover:bg-pink-900"
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
                        "Submit"
                      )}
                    </button>
                  </div>
                </form>
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
};

export default UpdatePluginsForm;

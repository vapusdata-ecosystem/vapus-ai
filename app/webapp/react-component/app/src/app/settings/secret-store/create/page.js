"use client";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Sidebar from "@/app/components/platform/main-sidebar";
import Header from "@/app/components/platform/header";
import CredentialsForm from "@/app/components/generic-credentials";
import { secretStoreFormApi } from "@/app/utils/settings-endpoint/secret-store-api";
import ToastContainerMessage from "@/app/components/notification/customToast";
import LoadingOverlay from "@/app/components/loading/loading";
import { toast } from "react-toastify";
import { enumsApi } from "@/app/utils/developers-endpoint/enums";
import { strTitle } from "@/app/components/JS/common";

export default function SecretServiceForm({
  createActionParams,
  backListingLink,
}) {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("form");
  const [formData, setFormData] = useState({
    name: "",
    expiresAt: "",
    secretType: "",
    data: "",
    description: "",
  });
  const [isLoading, setIsLoading] = useState(false);
  const [credentialsData, setCredentialsData] = useState({});
  const [secretTypes, setSecretTypes] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Set current date as default for expiresAt
    const now = new Date();
    const localDateTime = now.toISOString().slice(0, 16);
    setFormData((prev) => ({
      ...prev,
      expiresAt: localDateTime,
    }));

    // Fetch enum data
    const fetchEnumsData = async () => {
      try {
        setIsLoading(true);
        const response = await enumsApi.getEnums();
        const enumResponses = response.enumResponse || [];

        enumResponses.forEach((enumData) => {
          if (enumData.name === "VapusSecretType") {
            setSecretTypes(enumData.value || []);
          }
        });
      } catch (error) {
        console.error("Failed to fetch enum data:", error);
        setError(error.message);
        toast.error("Failed to fetch configuration data");
      } finally {
        setIsLoading(false);
      }
    };

    fetchEnumsData();
  }, []);

  const handleTabChange = (tab) => {
    setActiveTab(tab);
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    // Strip "spec." from the name if present
    const fieldName = name.replace("spec.", "");
    setFormData((prev) => ({
      ...prev,
      [fieldName]: value,
    }));
  };

  const handleSecretTypeChange = (e) => {
    const value = e.target.value;
    setFormData((prev) => ({
      ...prev,
      secretType: value,
    }));
  };

  const handleCredentialsChange = (data) => {
    setCredentialsData(data);
  };

  const strToUniArray = (str) => {
    const encoder = new TextEncoder();
    return encoder.encode(str);
  };

  const uint8ArrayToBase64 = (uint8Array) => {
    let binary = "";
    uint8Array.forEach((byte) => {
      binary += String.fromCharCode(byte);
    });
    return btoa(binary);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Combine formData with credentialsData
    let combinedFormData = {
      ...formData,
      credentials: credentialsData,
    };

    console.log("Form Data:", combinedFormData);

    let dataObj = { ...combinedFormData };

    // Handle expiration date
    if (dataObj.expiresAt) {
      const runDate = new Date(dataObj.expiresAt);
      dataObj.expiresAt = Math.floor(runDate.getTime() / 1000);
    }

    if (dataObj.data === "" && dataObj.secretType === "VAPUS_CREDENTIAL") {
      const genericCreds = dataObj.credentials?.credentials || {};
      const vapus_secrets = JSON.stringify(genericCreds);
      const vapus_secrets_binary = strToUniArray(vapus_secrets);
      const vapus_secrets_base64 = uint8ArrayToBase64(vapus_secrets_binary);
      dataObj.data = vapus_secrets_base64;
    } else if (dataObj.data !== "") {
      const custom_secrets_binary = strToUniArray(dataObj.data);
      const custom_secrets_base64 = uint8ArrayToBase64(custom_secrets_binary);
      dataObj.data = custom_secrets_base64;
    }

    console.log("Transformed Data Object:", dataObj);
    console.log("Raw Credentials Data:", dataObj.credentials);

    // For API payload, you may want to separate them again
    const payload = {
      spec: {
        name: dataObj.name,
        expiresAt: dataObj.expiresAt,
        secretType: dataObj.secretType,
        data: dataObj.data,
        description: dataObj.description,
      },
    };

    console.log("Full API Payload:", payload);

    // Submit form
    await submitCreateForm(dataObj);
  };

  const submitCreateForm = async (dataObj) => {
    try {
      setIsLoading(true);

      const payload = {
        spec: dataObj,
      };
      const output = await secretStoreFormApi.getSecretStoreForm(payload);

      console.log("Resource created:", output);
      const resourceInfo = output.result;
      if (resourceInfo) {
        toast.success(
          "Resource Created",
          `${resourceInfo.resource} Resource created successfully.`
        );
        router.push(
          `${backListingLink || "/settings/secret-store"}/${
            resourceInfo.resourceId
          }`
        );
      }
    } catch (error) {
      console.error("Error sending API request:", error);
      toast.error("Resource Creation Failed", "Resource creation failed");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-zinc-800 flex h-screen ">
      <Sidebar />
      <div className="overflow-y-auto h-screen w-full overflow-y-auto scrollbar">
        <Header
          sectionHeader="Create Secret"
          hideBackListingLink={false}
          backListingLink="./"
        />
        <ToastContainerMessage />
        <LoadingOverlay isLoading={isLoading} />
        <div className="flex-grow p-4 overflow-y-auto w-full">
          <section id="grids" className="space-y-2">
            <div className="max-w-6xl mx-auto bg-[#1b1b1b] shadow rounded-lg p-2">
              <div className="text-gray-100 mb-2 flex justify-center">
                {/* <button
                  id="yamlSpecB"
                  className={`whitespace-nowrap border-b-2 ${
                    activeTab === "yaml"
                      ? "border-orange-700 text-orange-700 font-semibold"
                      : "border-black"
                  } py-2 px-2 text-md font-medium focus:outline-none`}
                  onClick={() => handleTabChange("yaml")}
                >
                  YAML
                </button> */}
                <button
                  id="formSpecB"
                  className={`whitespace-nowrap border-b-2 ${
                    activeTab === "form"
                      ? "border-orange-700 text-orange-700 font-semibold"
                      : "border-black"
                  } py-2 px-2 text-md font-medium focus:outline-none ml-4`}
                  onClick={() => handleTabChange("form")}
                  suppressHydrationWarning
                >
                  Form
                </button>
              </div>

              {activeTab === "yaml" && (
                <div id="yamlSpec">
                  {/* YAML editor would be imported here */}
                  <p className="text-gray-100">YAML Editor Placeholder</p>
                </div>
              )}

              {activeTab === "form" && (
                <div id="formSpec">
                  <form
                    id="SecretServiceSpec"
                    className="space-y-2 border border-zinc-500 rounded-md text-gray-100 p-2"
                    onSubmit={handleSubmit}
                  >
                    <fieldset className="p-4 rounded">
                      <div className="grid grid-cols-2 md:grid-cols-2 gap-4 mb-4">
                        {/* Name */}
                        <div>
                          <label className="labels">Name</label>
                          <input
                            id="name"
                            name="spec.name"
                            type="text"
                            placeholder="Enter name"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                            value={formData.name}
                            onChange={handleChange}
                            required
                            suppressHydrationWarning
                          />
                        </div>

                        {/* ExpiresAt */}
                        <div>
                          <label className="labels">Expires At</label>
                          <input
                            name="spec.expiresAt"
                            type="datetime-local"
                            id="dateInput"
                            placeholder="Select date and time"
                            max="2050-12-31T23:59"
                            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                            value={formData.expiresAt}
                            onChange={handleChange}
                            required
                          />
                        </div>

                        {/* Secret Type */}
                        <div className="flex-1 min-w-[200px]">
                          <div className="relative">
                            <label className="labels">Secret Type</label>
                            <select
                              id="modelNode"
                              name="spec.secretType"
                              className="mt-1 w-full flex justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
                              value={formData.secretType}
                              onChange={handleSecretTypeChange}
                              required
                              suppressHydrationWarning
                            >
                              <option value="">Select Secret Type</option>
                              {secretTypes.map((type, index) => (
                                <option
                                  key={index}
                                  value={type}
                                  className="text-sm "
                                >
                                  {strTitle(type)}
                                </option>
                              ))}
                            </select>
                          </div>
                        </div>
                      </div>

                      <fieldset className="rounded mb-4 border border-zinc-500 rounded-md shadow-sm p-4">
                        <legend className="text-lg font-semibold"></legend>

                        {/* Data */}
                        <div
                          id="custom-data"
                          className={`grid grid-cols-1 md:grid-cols-1 gap-4 mb-2 ${
                            formData.secretType === "CUSTOM_SECRET"
                              ? ""
                              : "hidden"
                          }`}
                        >
                          <div className="col-span-2">
                            <label className="block text-sm font-medium">
                              Data
                            </label>
                            <textarea
                              name="spec.data"
                              rows="3"
                              type="text"
                              placeholder="Enter Data"
                              className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full "
                              value={formData.data}
                              onChange={handleChange}
                            ></textarea>
                          </div>
                        </div>

                        {/* Generic Credentials */}
                        <div
                          id="generic-cred"
                          className={`grid grid-cols-1 md:grid-cols-1 gap-4 mb-2 ${
                            formData.secretType === "VAPUS_CREDENTIAL"
                              ? ""
                              : "hidden"
                          }`}
                        >
                          <CredentialsForm onChange={handleCredentialsChange} />
                        </div>

                        {/* Description */}
                        <div className="grid grid-cols-1 mt-2 md:grid-cols-1 gap-4">
                          <div className="col-span-2">
                            <label className="block text-sm font-medium">
                              Description
                            </label>
                            <textarea
                              name="spec.description"
                              rows="4"
                              type="text"
                              placeholder="Enter Description"
                              className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                              value={formData.description}
                              onChange={handleChange}
                            ></textarea>
                          </div>
                        </div>
                      </fieldset>
                    </fieldset>

                    {/* Submit Button */}
                    <div className="mt-4 flex justify-end space-x-2">
                      {!isLoading ? (
                        <button
                          type="submit"
                          id="submit"
                          className="px-6 py-2 bg-orange-700 text-white rounded-md shadow hover:bg-pink-900"
                          suppressHydrationWarning
                        >
                          Submit
                        </button>
                      ) : (
                        <button
                          type="button"
                          id="loading"
                          className="px-6 py-2 bg-orange-700 text-white rounded-md shadow hover:bg-pink-900"
                        >
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
                              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                            ></path>
                          </svg>
                        </button>
                      )}
                    </div>
                  </form>
                </div>
              )}
            </div>
          </section>
        </div>
        <div hidden id="createTemplate">
          {createActionParams?.yamlSpec}
        </div>
      </div>
    </div>
  );
}

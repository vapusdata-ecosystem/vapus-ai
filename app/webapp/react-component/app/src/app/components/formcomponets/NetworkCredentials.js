"use client";
import { useState, useEffect, useCallback } from "react";
import CredentialsData from "./credentials";

const NetworkCredentials = ({ onParamsChange, initialParams = {} }) => {
  // Consolidate related form fields into a single state object
  const [formData, setFormData] = useState({
    url: initialParams?.url || "",
    port: initialParams?.port || "",
    version: initialParams?.version || "",
    secretName: initialParams?.secretName || "",
    isAlreadyInSecret: initialParams?.isAlreadyInSecretBs || false,
    credentials: {
      username: initialParams?.credentials?.username || "",
      password: initialParams?.credentials?.password || "",
      apiToken: initialParams?.credentials?.apiToken || "",
      apiTokenType: initialParams?.credentials?.apiTokenType || "BASIC",
      aws_creds: initialParams?.credentials?.aws_creds || null,
      gcp_creds: initialParams?.credentials?.gcp_creds || null,
      azure_creds: initialParams?.credentials?.azure_creds || null,
    },
  });

  // Handle initialParams changes with a single effect
  useEffect(() => {
    if (initialParams && Object.keys(initialParams).length > 0) {
      setFormData({
        url: initialParams.url || "",
        port: initialParams.port || "",
        version: initialParams.version || "",
        secretName: initialParams.secretName || "",
        isAlreadyInSecret: initialParams.isAlreadyInSecretBs || false,
        credentials: {
          username: initialParams?.credentials?.username || "",
          password: initialParams?.credentials?.password || "",
          apiToken: initialParams?.credentials?.apiToken || "",
          apiTokenType: initialParams?.credentials?.apiTokenType || "BASIC",
          aws_creds: initialParams?.credentials?.aws_creds || null,
          gcp_creds: initialParams?.credentials?.gcp_creds || null,
          azure_creds: initialParams?.credentials?.azure_creds || null,
        },
      });
    }
  }, [initialParams]);

  // Generic handler for form field changes
  const handleFormChange = (field, value) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  // Handle CredentialsData component changes
  const handleCredentialsChange = (data) => {
    setFormData((prev) => ({
      ...prev,
      credentials: data.credentials,
    }));
  };

  // Format network parameters for parent component
  const getNetworkParams = useCallback(() => {
    const { url, port, version, secretName, isAlreadyInSecret, credentials } =
      formData;

    const netParams = {
      url,
      port: isNaN(parseInt(port)) ? "0" : String(parseInt(port)),
      version,
      secretName: isAlreadyInSecret ? secretName : "",
      isAlreadyInSecretBs: isAlreadyInSecret,
      credentials: {
        username: credentials?.username || "",
        password: credentials?.password || "",
        apiToken: credentials?.apiToken || "",
        apiTokenType: credentials?.apiTokenType || "BASIC",
      },
    };

    // Add aws  credentials if they exist
    if (
      credentials?.aws_creds &&
      Object.values(credentials.aws_creds).some(Boolean)
    ) {
      netParams.credentials.aws_creds = credentials.aws_creds;
    }

    // Always include GCP credentials with the specified format
    netParams.credentials.gcpCreds = {
      serviceAccountKey: credentials?.gcp_creds?.service_account_key || "",
      base64Encoded: credentials?.gcp_creds?.base64_encoded || false,
      projectId: credentials?.gcp_creds?.project_id || "",
      region: credentials?.gcp_creds?.region || "",
      zone: credentials?.gcp_creds?.zone || "",
    };

    // Add Azure credentials if they exist
    if (
      credentials?.azure_creds &&
      Object.values(credentials.azure_creds).some(Boolean)
    ) {
      netParams.credentials.azure_creds = credentials.azure_creds;
    }

    return netParams;
  }, [formData]);

  // Update parent component when data changes
  useEffect(() => {
    if (onParamsChange) {
      onParamsChange(getNetworkParams());
    }
  }, [getNetworkParams, onParamsChange]);

  const { url, port, version, secretName, isAlreadyInSecret } = formData;

  return (
    <div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label htmlFor="net_address" className="labels">
            Address
          </label>
          <input
            id="net_address"
            name="url"
            type="text"
            placeholder="Enter url"
            value={url}
            onChange={(e) => handleFormChange("url", e.target.value)}
            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
            suppressHydrationWarning
          />
        </div>
        <div>
          <label htmlFor="net_port" className="labels">
            Port
          </label>
          <input
            id="net_port"
            name="port"
            type="number"
            placeholder="Enter port"
            value={port}
            onChange={(e) => handleFormChange("port", e.target.value)}
            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
            suppressHydrationWarning
          />
        </div>
        <div>
          <label htmlFor="net_version" className="labels">
            Version
          </label>
          <input
            id="net_version"
            type="text"
            name="version"
            placeholder="Enter version"
            value={version}
            onChange={(e) => handleFormChange("version", e.target.value)}
            className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
            suppressHydrationWarning
          />
        </div>
      </div>

      <div className="mt-4">
        <div className="grid grid-cols-3 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <div className="flex items-center m-4">
            <input
              id="isAlreadyInSecretBs"
              name="isAlreadyInSecretBs"
              type="checkbox"
              className="h-4 w-4 text-indigo-600 border-gray-300 rounded accent-orange-700"
              checked={isAlreadyInSecret}
              onChange={(e) =>
                handleFormChange("isAlreadyInSecret", e.target.checked)
              }
            />
            <label
              htmlFor="isAlreadyInSecretBs"
              className="ml-2 block text-sm font-medium text-gray-100"
            >
              Already in Secret Store?
            </label>
          </div>

          {isAlreadyInSecret && (
            <div>
              <label htmlFor="secretName" className="labels">
                Secret Name
              </label>
              <input
                id="secretName"
                type="text"
                name="secretName"
                placeholder="Enter secret name"
                value={secretName}
                onChange={(e) => handleFormChange("secretName", e.target.value)}
                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                required={isAlreadyInSecret}
              />
            </div>
          )}
        </div>
      </div>

      <div className="mt-4">
        <CredentialsData
          onChange={handleCredentialsChange}
          isAlreadyInSecret={isAlreadyInSecret}
          value={{ credentials: formData.credentials }}
        />
      </div>
    </div>
  );
};

export default NetworkCredentials;

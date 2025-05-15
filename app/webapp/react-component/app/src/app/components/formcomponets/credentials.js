"use client";
import { enumsApi } from "@/app/utils/developers-endpoint/enums";
import { useState, useEffect, useRef } from "react";
import { strTitle } from "../JS/common";

/**
 * @typedef {Object} AwsCreds
 * @property {string} access_key_id
 * @property {string} secret_access_key
 * @property {string} region
 * @property {string} session_token
 * @property {string} role_arn
 */

/**
 * @typedef {Object} GcpCreds
 * @property {string} service_account_key
 * @property {boolean} base64_encoded
 * @property {string} project_id
 * @property {string} region
 * @property {string} zone
 */

/**
 * @typedef {Object} AzureCreds
 * @property {string} tenant_id
 * @property {string} client_id
 * @property {string} client_secret
 */

/**
 * @typedef {Object} Credentials
 * @property {string} username
 * @property {string} password
 * @property {string} apiToken
 * @property {string} apiTokenType
 * @property {AwsCreds} [aws_creds]
 * @property {GcpCreds} [gcp_creds]
 * @property {AzureCreds} [azure_creds]
 */

/**
 * @typedef {Object} CredentialsFormData
 * @property {Credentials} credentials
 */

/**
 * @typedef {Object} CredentialsFormProps
 * @property {function(CredentialsFormData): void} [onChange]
 * @property {boolean} isAlreadyInSecret - Flag indicating if credentials are already in secret store
 * @property {Partial<CredentialsFormData>} [value]
 * @property {Object} [enums]
 * @property {string[]} [enums.ApiTokenType]
 */

/**
 * @param {any} obj1
 * @param {any} obj2
 * @returns {boolean}
 */
function isEqual(obj1, obj2) {
  if (obj1 === obj2) return true;
  if (
    typeof obj1 !== "object" ||
    obj1 === null ||
    typeof obj2 !== "object" ||
    obj2 === null
  )
    return false;

  const keys1 = Object.keys(obj1);
  const keys2 = Object.keys(obj2);

  if (keys1.length !== keys2.length) return false;

  for (const key of keys1) {
    if (!keys2.includes(key)) return false;
    if (!isEqual(obj1[key], obj2[key])) return false;
  }

  return true;
}

/**
 * Credentials form component for handling various types of credentials
 * @param {CredentialsFormProps} props
 * @returns {JSX.Element}
 */
export default function CredentialsData({
  onChange,
  isAlreadyInSecret = false,
  value = {},
  enums = {},
}) {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [apiTokenTypes, setApiTokenTypes] = useState([]);

  const prevValueRef = useRef();

  // Fetch API token types from enum API
  useEffect(() => {
    const fetchApiTokenType = async () => {
      try {
        setIsLoading(true);
        const data = await enumsApi.getEnums();
        const enumResponses = data.enumResponse || [];
        const apiTokenTypeEnum = enumResponses.find(
          (item) => item.name === "ApiTokenType"
        );
        if (apiTokenTypeEnum && Array.isArray(apiTokenTypeEnum.value)) {
          setApiTokenTypes(apiTokenTypeEnum.value);
        } else {
          setApiTokenTypes([]);
          console.warn("ApiTokenType enum not found in response");
        }
      } catch (error) {
        console.error("Failed to fetch API Token Types:", error);
        setError(error.message);
        setApiTokenTypes([]);
      } finally {
        setIsLoading(false);
      }
    };

    fetchApiTokenType();
  }, []);

  // Use our fetched data or any enums provided via props
  const mergedEnums = {
    ApiTokenType: enums.ApiTokenType || apiTokenTypes,
  };

  //  initializing and copying objects
  const deepClone = (obj) => {
    if (obj === null || typeof obj !== "object") return obj;
    const result = Array.isArray(obj) ? [] : {};
    for (const key in obj) {
      if (Object.prototype.hasOwnProperty.call(obj, key)) {
        result[key] = deepClone(obj[key]);
      }
    }
    return result;
  };

  const [formData, setFormData] = useState(() => {
    const initialData = {
      credentials: {
        username: "",
        password: "",
        apiToken: "",
        apiTokenType: "",
      },
    };

    // Merge values from props if available
    if (value.credentials) {
      Object.keys(value.credentials).forEach((key) => {
        if (
          key === "aws_creds" ||
          key === "gcp_creds" ||
          key === "azure_creds"
        ) {
          if (value.credentials[key]) {
            initialData.credentials[key] = deepClone(value.credentials[key]);
          }
        } else {
          initialData.credentials[key] = value.credentials[key];
        }
      });
    }

    return initialData;
  });

  useEffect(() => {
    if (!prevValueRef.current) {
      prevValueRef.current = deepClone(value);
      return;
    }
    if (isEqual(prevValueRef.current, value)) {
      return;
    }

    prevValueRef.current = deepClone(value);
    const newFormData = deepClone(formData);

    if (value.credentials) {
      Object.keys(value.credentials).forEach((key) => {
        if (
          key === "aws_creds" ||
          key === "gcp_creds" ||
          key === "azure_creds"
        ) {
          if (value.credentials[key]) {
            newFormData.credentials[key] = deepClone(value.credentials[key]);
          }
        } else {
          newFormData.credentials[key] = value.credentials[key];
        }
      });
    }

    setFormData(newFormData);
  }, [value]);

  //  removing empty cloud credential sections
  const cleanupData = (data) => {
    const cleanData = deepClone(data);

    const hasValues = (obj) => {
      if (!obj) return false;

      return Object.values(obj).some((val) => {
        if (typeof val === "boolean") return true;
        if (typeof val === "object") return hasValues(val);
        return val !== null && val !== undefined && val !== "";
      });
    };

    if (
      cleanData.credentials?.aws_creds &&
      !hasValues(cleanData.credentials.aws_creds)
    ) {
      delete cleanData.credentials.aws_creds;
    }

    if (
      cleanData.credentials?.gcp_creds &&
      !hasValues(cleanData.credentials.gcp_creds)
    ) {
      delete cleanData.credentials.gcp_creds;
    }

    if (
      cleanData.credentials?.azure_creds &&
      !hasValues(cleanData.credentials.azure_creds)
    ) {
      delete cleanData.credentials.azure_creds;
    }

    return cleanData;
  };

  // handleInputChange to update local state and notify parent
  const handleInputChange = (e) => {
    const { name, value: inputValue, type } = e.target;
    const isCheckbox = type === "checkbox";
    const fieldValue = isCheckbox ? e.target.checked : inputValue;

    setFormData((prev) => {
      let newData = deepClone(prev);
      if (name.includes(".")) {
        const parts = name.split(".");
        let current = newData;
        for (let i = 0; i < parts.length - 1; i++) {
          if (!current[parts[i]]) {
            current[parts[i]] = {};
          }
          current = current[parts[i]];
        }
        current[parts[parts.length - 1]] = fieldValue;
      } else {
        newData[name] = fieldValue;
      }

      return newData;
    });

    // Call onChange outside the state update to prevent render cycles
    if (onChange) {
      let newData = deepClone(formData);

      if (name.includes(".")) {
        const parts = name.split(".");
        let current = newData;

        for (let i = 0; i < parts.length - 1; i++) {
          if (!current[parts[i]]) {
            current[parts[i]] = {};
          }
          current = current[parts[i]];
        }

        current[parts[parts.length - 1]] = fieldValue;
      } else {
        newData[name] = fieldValue;
      }

      const cleanData = cleanupData(newData);
      onChange(cleanData);
    }
  };

  return (
    <div id="genericCredentials">
      {/* Show credential fields only when not using secret store */}
      {!isAlreadyInSecret && (
        <>
          <div className="grid grid-cols-3 md:grid-cols-3 lg:grid-cols-4 gap-4 mb-4 mt-2">
            <div>
              <label htmlFor="username" className="labels">
                Username
              </label>
              <input
                id="username"
                name="credentials.username"
                type="text"
                placeholder="Enter username"
                value={formData.credentials.username || ""}
                onChange={handleInputChange}
                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                suppressHydrationWarning
              />
            </div>

            <div>
              <label htmlFor="password" className="labels">
                Password
              </label>
              <input
                id="password"
                name="credentials.password"
                type="password"
                placeholder="Enter password"
                value={formData.credentials.password || ""}
                onChange={handleInputChange}
                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                suppressHydrationWarning
              />
            </div>

            <div>
              <label htmlFor="apiToken" className="labels">
                API Token
              </label>
              <input
                id="apiToken"
                name="credentials.apiToken"
                type="text"
                placeholder="Enter API token"
                value={formData.credentials.apiToken || ""}
                onChange={handleInputChange}
                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                suppressHydrationWarning
              />
            </div>

            <div>
              <label htmlFor="apiTokenType" className="labels">
                API Token Type
              </label>
              <select
                id="apiTokenType"
                name="credentials.apiTokenType"
                value={formData.credentials.apiTokenType || ""}
                onChange={handleInputChange}
                className="mt-1 w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-zinc-800 text-sm"
                suppressHydrationWarning
                disabled={isLoading}
              >
                <option value="">Select Token Type</option>
                {mergedEnums.ApiTokenType.map((val) => (
                  <option key={val} value={val}>
                    {strTitle(val)}
                  </option>
                ))}
              </select>
              {isLoading && (
                <div className="text-xs mt-1 text-gray-300">Loading...</div>
              )}
              {error && (
                <div className="text-xs mt-1 text-red-400">
                  Error loading token types
                </div>
              )}
            </div>
          </div>
        </>
      )}

      {/* Show cloud credentials sections only when checkbox is NOT checked */}
      {!isAlreadyInSecret && (
        <>
          {/* AWS Creds */}
          <details className="border border-zinc-500 p-4 rounded mb-4 mt-2">
            <summary className="text-lg font-semibold cursor-pointer">
              AWS Creds
            </summary>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
              <div>
                <label htmlFor="access_key_id" className="labels">
                  Access Key ID
                </label>
                <input
                  id="access_key_id"
                  type="text"
                  name="credentials.aws_creds.access_key_id"
                  placeholder="Access Key ID"
                  value={formData.credentials.aws_creds?.access_key_id || ""}
                  onChange={handleInputChange}
                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                  suppressHydrationWarning
                />
              </div>
              <div>
                <label htmlFor="secret_access_key" className="labels">
                  Secret Access Key
                </label>
                <input
                  id="secret_access_key"
                  type="text"
                  name="credentials.aws_creds.secret_access_key"
                  placeholder="Secret Access Key"
                  value={
                    formData.credentials.aws_creds?.secret_access_key || ""
                  }
                  onChange={handleInputChange}
                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                  suppressHydrationWarning
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
              <div>
                <label htmlFor="region" className="labels">
                  Region
                </label>
                <input
                  id="region"
                  type="text"
                  name="credentials.aws_creds.region"
                  placeholder="Region"
                  value={formData.credentials.aws_creds?.region || ""}
                  onChange={handleInputChange}
                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                  suppressHydrationWarning
                />
              </div>
              <div>
                <label htmlFor="session_token" className="labels">
                  Session Token
                </label>
                <input
                  id="session_token"
                  type="text"
                  name="credentials.aws_creds.session_token"
                  placeholder="Session Token"
                  value={formData.credentials.aws_creds?.session_token || ""}
                  onChange={handleInputChange}
                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                  suppressHydrationWarning
                />
              </div>
            </div>

            <div className="mt-4">
              <label htmlFor="role_arn" className="labels">
                Role Arn
              </label>
              <input
                id="role_arn"
                type="text"
                name="credentials.aws_creds.role_arn"
                placeholder="Role Arn"
                value={formData.credentials.aws_creds?.role_arn || ""}
                onChange={handleInputChange}
                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                suppressHydrationWarning
              />
            </div>
          </details>

          {/* GCP Creds */}
          <details className="border border-zinc-500 p-4 rounded mb-4">
            <summary className="text-lg font-semibold cursor-pointer">
              GCP Creds
            </summary>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
              <div className="md:col-span-3">
                <label htmlFor="service_account_key" className="labels">
                  Service Account Key
                </label>
                <textarea
                  id="service_account_key"
                  name="credentials.gcp_creds.service_account_key"
                  placeholder="Service Account Key"
                  value={
                    formData.credentials.gcp_creds?.service_account_key || ""
                  }
                  onChange={handleInputChange}
                  className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                  suppressHydrationWarning
                ></textarea>
              </div>

              <div className="flex items-center mt-2">
                <input
                  id="base64_encoded"
                  type="checkbox"
                  name="credentials.gcp_creds.base64_encoded"
                  checked={
                    formData.credentials.gcp_creds?.base64_encoded || false
                  }
                  onChange={handleInputChange}
                  className="mr-2 h-4 w-4 text-orange-700 border-gray-300 rounded accent-orange-700"
                  suppressHydrationWarning
                />
                <label
                  htmlFor="base64_encoded"
                  className="labels cursor-pointer"
                >
                  Is Base64 Encoded?
                </label>
              </div>

              <div>
                <label htmlFor="gcp_project_id" className="labels">
                  Project ID
                </label>
                <input
                  id="gcp_project_id"
                  type="text"
                  name="credentials.gcp_creds.project_id"
                  placeholder="Project Id"
                  value={formData.credentials.gcp_creds?.project_id || ""}
                  onChange={handleInputChange}
                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                  suppressHydrationWarning
                />
              </div>

              <div>
                <label htmlFor="gcp_region" className="labels">
                  Region
                </label>
                <input
                  id="gcp_region"
                  type="text"
                  name="credentials.gcp_creds.region"
                  placeholder="Region"
                  value={formData.credentials.gcp_creds?.region || ""}
                  onChange={handleInputChange}
                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                  suppressHydrationWarning
                />
              </div>

              <div>
                <label htmlFor="gcp_zone" className="labels">
                  Zone
                </label>
                <input
                  id="gcp_zone"
                  type="text"
                  name="credentials.gcp_creds.zone"
                  placeholder="Zone"
                  value={formData.credentials.gcp_creds?.zone || ""}
                  onChange={handleInputChange}
                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                  suppressHydrationWarning
                />
              </div>
            </div>
          </details>

          {/* AZURE CREDS */}
          <details className="border border-zinc-500 p-4 rounded mb-4">
            <summary className="text-lg font-semibold cursor-pointer">
              Azure Creds
            </summary>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
              <div>
                <label htmlFor="azure_tenant_id" className="labels">
                  Tenant ID
                </label>
                <input
                  id="azure_tenant_id"
                  type="text"
                  name="credentials.azure_creds.tenant_id"
                  placeholder="Tenant Id"
                  value={formData.credentials.azure_creds?.tenant_id || ""}
                  onChange={handleInputChange}
                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                  suppressHydrationWarning
                />
              </div>

              <div>
                <label htmlFor="azure_client_id" className="labels">
                  Client ID
                </label>
                <input
                  id="azure_client_id"
                  type="text"
                  name="credentials.azure_creds.client_id"
                  placeholder="Client Id"
                  value={formData.credentials.azure_creds?.client_id || ""}
                  onChange={handleInputChange}
                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                  suppressHydrationWarning
                />
              </div>

              <div>
                <label htmlFor="azure_client_secret" className="labels">
                  Client Secret
                </label>
                <input
                  id="azure_client_secret"
                  type="text"
                  name="credentials.azure_creds.client_secret"
                  placeholder="Client Secret"
                  value={formData.credentials.azure_creds?.client_secret || ""}
                  onChange={handleInputChange}
                  className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                  suppressHydrationWarning
                />
              </div>
            </div>
          </details>
        </>
      )}
    </div>
  );
}

// Export the credential form data type for use in other files
export const CredentialsFormData = undefined; // Just for documentation

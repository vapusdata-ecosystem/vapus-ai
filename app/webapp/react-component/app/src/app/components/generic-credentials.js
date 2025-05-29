"use client";
import { useState, useEffect, useRef } from "react";

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
 * @property {string} secretName
 * @property {string} accessScope
 * @property {string} db
 * @property {Credentials} credentials
 */

/**
 * @typedef {Object} CredentialsFormProps
 * @property {function(CredentialsFormData): void} [onChange]
 * @property {Partial<CredentialsFormData>} [value]
 * @property {Object} [enums]
 * @property {string[]} [enums.DataSourceAccessScope]
 * @property {string[]} [enums.ApiTokenType]
 */

// Hardcoded enum values
const DATA_SOURCE_ACCESS_SCOPE = [
  "ANTHROPIC",
  "AWS",
  "AZURE_OPENAI",
  "AZURE_PHI",
  "BEDROCK",
  "BITBUCKET",
  "GCP",
  "GEMINI",
  "GENERIC",
  "GITHUB",
  "GITLAB",
  "GROQ",
  "HUGGING_FACE",
  "META",
  "MICROSOFT",
  "MISTRAL",
  "MONGO_ORG",
  "OLLAMA",
  "OPENAI",
  "REDHAT",
  "REDIS_ORG",
  "SELF_HOSTED",
  "TOGETHER",
  "VAPUS",
];

const API_TOKEN_TYPE = ["APIKEY", "BASIC", "BEARER"];

/**
 * Helper function for deep comparison of objects
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
export default function CredentialsForm({ onChange, value = {}, enums = {} }) {
  // Use our hardcoded enums if not provided
  const mergedEnums = {
    DataSourceAccessScope:
      enums.DataSourceAccessScope || DATA_SOURCE_ACCESS_SCOPE,
    ApiTokenType: enums.ApiTokenType || API_TOKEN_TYPE,
  };

  // Add state for tracking checkbox status
  const [isChecked, setIsChecked] = useState(false);

  // Use a ref to track previous value for comparison
  const prevValueRef = useRef();

  // Deep clone helper for initializing and copying objects
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

  // Initialize state for form data with a function to ensure it runs only once
  const [formData, setFormData] = useState(() => {
    // Fix: Create a proper structure with top-level fields and credentials object
    const initialData = {
      secretName: value.secretName || "",
      accessScope: value.accessScope || "",
      db: value.db || "",
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

  // Update local state when parent value changes, with deep comparison
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

    // Fix: Update top-level fields properly
    if (value.secretName !== undefined)
      newFormData.secretName = value.secretName;
    if (value.accessScope !== undefined)
      newFormData.accessScope = value.accessScope;
    if (value.db !== undefined) newFormData.db = value.db;

    // Handle credentials
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

  // Clean up data by removing empty cloud credential sections
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
  const handleCheckboxChange = (e) => {
    setIsChecked(e.target.checked);
  };

  // Helper function to format enum values for display
  const formatEnumValue = (prefix, value) => {
    return value
      .split("_")
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(" ");
  };

  return (
    <div id="genericCredentials">
      <div className="grid grid-cols-3 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <div className="flex items-center m-4">
          <input
            id="isAlreadyInSecretBs"
            name="dsCreds[0].isAlreadyInSecretBs"
            type="checkbox"
            className="h-4 w-4 text-indigo-600 border-gray-300 rounded accent-orange-700"
            checked={isChecked}
            onChange={handleCheckboxChange}
          />
          <label
            htmlFor="isAlreadyInSecretBs"
            className="ml-2 block text-sm font-medium text-gray-100"
          >
            Already in Secret Store?
          </label>
        </div>

        {/* Conditionally render Secret Name field when checkbox is checked */}
        {isChecked && (
          <div>
            <label htmlFor="secretName" className="labels">
              Secret Name
            </label>
            <input
              id="secretName"
              type="text"
              name="secretName"
              placeholder="Enter secret name"
              value={formData.secretName}
              onChange={handleInputChange}
              className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
            />
          </div>
        )}

        {/* Show remaining fields only when checkbox is NOT checked */}
        {!isChecked && (
          <>
            <div>
              <label htmlFor="accessScope" className="labels">
                Access Scope
              </label>
              <select
                id="accessScope"
                name="accessScope"
                value={formData.accessScope}
                onChange={handleInputChange}
                className="mt-1 w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-zinc-800 text-sm overflow-y-auto scrollbar"
                suppressHydrationWarning
              >
                <option value="ALL">Select Credential Access Scope</option>
                {mergedEnums.DataSourceAccessScope.map((val) => (
                  <option key={val} value={val}>
                    {formatEnumValue("DSCAS", val)}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label htmlFor="db" className="labels">
                Datastore
              </label>
              <input
                id="db"
                type="text"
                name="db"
                placeholder="Enter datastore"
                value={formData.db}
                onChange={handleInputChange}
                className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                suppressHydrationWarning
              />
            </div>

            <div>
              <label htmlFor="username" className="labels">
                Username
              </label>
              <input
                id="username"
                name="credentials.username"
                type="text"
                placeholder="Enter username"
                value={formData.credentials.username}
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
                value={formData.credentials.password}
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
                value={formData.credentials.apiToken}
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
                value={formData.credentials.apiTokenType}
                onChange={handleInputChange}
                className="mt-1 w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-zinc-800 text-sm"
                suppressHydrationWarning
              >
                <option value="">Select Token Type</option>
                {mergedEnums.ApiTokenType.map((val) => (
                  <option key={val} value={val}>
                    {formatEnumValue("", val)}
                  </option>
                ))}
              </select>
            </div>
          </>
        )}
      </div>

      {/* Show cloud credentials sections only when checkbox is NOT checked */}
      {!isChecked && (
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

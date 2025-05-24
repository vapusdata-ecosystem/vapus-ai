"use client";
import React from "react";

const AzureCredentials = ({ formData, handleInputChange }) => {
  return (
    <details className="border border-zinc-500 p-4 rounded mb-4">
      <summary className="text-lg font-semibold cursor-pointer">
        Azure Creds
      </summary>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label htmlFor="azure_tenant_id" className="labels">
            Tenant ID
          </label>
          <input
            id="azure_tenant_id"
            type="text"
            name="spec.attributes.network_params.credentials.azure_creds.tenant_id"
            placeholder="Tenant Id"
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            value={
              formData.attributes.network_params.credentials.azure_creds
                .tenant_id
            }
            onChange={handleInputChange}
          />
        </div>
        <div>
          <label htmlFor="azure_client_id" className="labels">
            Client ID
          </label>
          <input
            id="azure_client_id"
            type="text"
            name="spec.attributes.network_params.credentials.azure_creds.client_id"
            placeholder="Client Id"
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            value={
              formData.attributes.network_params.credentials.azure_creds
                .client_id
            }
            onChange={handleInputChange}
          />
        </div>
        <div>
          <label htmlFor="azure_client_secret" className="labels">
            Client Secret
          </label>
          <input
            id="azure_client_secret"
            type="text"
            name="spec.attributes.network_params.credentials.azure_creds.client_secret"
            placeholder="Client Secret"
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            value={
              formData.attributes.network_params.credentials.azure_creds
                .client_secret
            }
            onChange={handleInputChange}
          />
        </div>
      </div>
    </details>
  );
};

export default AzureCredentials;

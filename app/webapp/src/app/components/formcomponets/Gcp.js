"use client";
import React from "react";

const GCPCredentials = ({ formData, handleInputChange }) => {
  return (
    <details className="border border-zinc-500 p-4 rounded mb-4">
      <summary className="text-lg font-semibold cursor-pointer">
        GCP Creds
      </summary>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="md:col-span-3">
          <label htmlFor="gcp_service_account_key" className="labels">
            Service Account Key
          </label>
          <textarea
            id="gcp_service_account_key"
            name="spec.attributes.network_params.credentials.gcp_creds.service_account_key"
            placeholder="Service Account Key"
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            value={
              formData.attributes.network_params.credentials.gcp_creds
                .service_account_key
            }
            onChange={handleInputChange}
          ></textarea>
        </div>
        <div>
          <label className="flex items-center">
            <input
              id="gcp_base64_encoded"
              type="checkbox"
              name="spec.attributes.network_params.credentials.gcp_creds.base64_encoded"
              className="mr-2 h-4 w-4 text-orange-700 border-gray-300 rounded accent-orange-700"
              checked={
                formData.attributes.network_params.credentials.gcp_creds
                  .base64_encoded
              }
              onChange={handleInputChange}
            />
            <span className="text-sm font-medium">Is Base64 Encoded?</span>
          </label>
        </div>
        <div>
          <label htmlFor="gcp_project_id" className="labels">
            Project ID
          </label>
          <input
            id="gcp_project_id"
            type="text"
            name="spec.attributes.network_params.credentials.gcp_creds.project_id"
            placeholder="Project Id"
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            value={
              formData.attributes.network_params.credentials.gcp_creds
                .project_id
            }
            onChange={handleInputChange}
          />
        </div>
        <div>
          <label htmlFor="gcp_region" className="labels">
            Region
          </label>
          <input
            id="gcp_region"
            type="text"
            name="spec.attributes.network_params.credentials.gcp_creds.region"
            placeholder="Region"
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            value={
              formData.attributes.network_params.credentials.gcp_creds.region
            }
            onChange={handleInputChange}
          />
        </div>
        <div>
          <label htmlFor="gcp_zone" className="labels">
            Zone
          </label>
          <input
            id="gcp_zone"
            type="text"
            name="spec.attributes.network_params.credentials.gcp_creds.zone"
            placeholder="Zone"
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            value={
              formData.attributes.network_params.credentials.gcp_creds.zone
            }
            onChange={handleInputChange}
          />
        </div>
      </div>
    </details>
  );
};

export default GCPCredentials;

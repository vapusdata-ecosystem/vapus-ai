"use client";
import React from "react";

const AWSCredentials = ({ formData, handleInputChange }) => {
  return (
    <details className="border border-zinc-500 p-4 rounded mb-4 mt-2">
      <summary className="text-lg font-semibold cursor-pointer">
        AWS Creds
      </summary>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label htmlFor="aws_access_key_id" className="labels">
            Access Key ID
          </label>
          <input
            id="aws_access_key_id"
            type="text"
            name="spec.attributes.network_params.credentials.aws_creds.access_key_id"
            placeholder="Access Key ID"
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            value={
              formData.attributes.network_params.credentials.aws_creds
                .access_key_id
            }
            onChange={handleInputChange}
          />
        </div>
        <div>
          <label htmlFor="aws_secret_access_key" className="labels">
            Secret Access Key
          </label>
          <input
            id="aws_secret_access_key"
            type="text"
            name="spec.attributes.network_params.credentials.aws_creds.secret_access_key"
            placeholder="Secret Access Key"
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            value={
              formData.attributes.network_params.credentials.aws_creds
                .secret_access_key
            }
            onChange={handleInputChange}
          />
        </div>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
        <div>
          <label htmlFor="aws_region" className="labels">
            Region
          </label>
          <input
            id="aws_region"
            type="text"
            name="spec.attributes.network_params.credentials.aws_creds.region"
            placeholder="Region"
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            value={
              formData.attributes.network_params.credentials.aws_creds.region
            }
            onChange={handleInputChange}
          />
        </div>
        <div>
          <label htmlFor="aws_session_token" className="labels">
            Session Token
          </label>
          <input
            id="aws_session_token"
            type="text"
            name="spec.attributes.network_params.credentials.aws_creds.session_token"
            placeholder="Session Token"
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            value={
              formData.attributes.network_params.credentials.aws_creds
                .session_token
            }
            onChange={handleInputChange}
          />
        </div>
      </div>
      <div className="mt-4">
        <label htmlFor="aws_role_arn" className="labels">
          Role ARN
        </label>
        <input
          id="aws_role_arn"
          type="text"
          name="spec.attributes.network_params.credentials.aws_creds.role_arn"
          placeholder="Role Arn"
          className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
          value={
            formData.attributes.network_params.credentials.aws_creds.role_arn
          }
          onChange={handleInputChange}
        />
      </div>
    </details>
  );
};

export default AWSCredentials;

"use client";
import { useState, useEffect, useCallback } from "react";
import CredentialsForm from "./generic-credentials";

const NetworkCredentialsForm = ({ onParamsChange }) => {
  const [address, setAddress] = useState("");
  const [port, setPort] = useState(0);
  const [databases, setDatabases] = useState("");
  const [databasePrefixes, setDatabasePrefixes] = useState("");
  const [version, setVersion] = useState("");
  const [credentialsData, setCredentialsData] = useState({});

  const handleCredentialsChange = (data) => {
    setCredentialsData(data);
  };

  const getNetworkParams = useCallback(() => {
    let dbsArray = [];
    let prefixesArray = [];

    if (databases.trim() !== "") {
      dbsArray = databases.trim().split(",");
    }

    if (databasePrefixes.trim() !== "") {
      prefixesArray = databasePrefixes.trim().split(",");
    }

    const netParams = {
      address: address,
      port: isNaN(parseInt(port)) ? 0 : parseInt(port),
      databases: dbsArray,
      databasePrefixes: prefixesArray,
      version: version,
      dsCreds: [],
    };

    const creds = getGenericCredentialsData();
    creds.isAlreadyInSecretBs = true;
    netParams.dsCreds.push(creds);

    return netParams;
  }, [address, port, databases, databasePrefixes, version, credentialsData]);

  const getGenericCredentialsData = useCallback(() => {
    // Extract and return the credentials data
    if (!credentialsData || !credentialsData.credentials) {
      return {};
    }

    // Return the credentials object which contains all relevant credential fields
    return credentialsData.credentials;
  }, [credentialsData]);

  // Update parent component when data changes
  useEffect(() => {
    if (onParamsChange) {
      onParamsChange(getNetworkParams());
    }
  }, [getNetworkParams, onParamsChange]);

  return (
    <div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <input
            id="net_address"
            name="address"
            type="text"
            placeholder="Enter address"
            value={address}
            onChange={(e) => setAddress(e.target.value)}
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            suppressHydrationWarning
          />
        </div>
        <div>
          <input
            id="net_port"
            name="port"
            type="number"
            placeholder="Enter port"
            value={port}
            onChange={(e) => setPort(e.target.value)}
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            suppressHydrationWarning
          />
        </div>
        <div>
          <input
            type="text"
            name="databases"
            placeholder="Enter datasources (comma separated)"
            value={databases}
            onChange={(e) => setDatabases(e.target.value)}
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            suppressHydrationWarning
          />
        </div>

        <div>
          <input
            type="text"
            name="databasePrefixes"
            placeholder="Prefix1,Prefix2"
            value={databasePrefixes}
            onChange={(e) => setDatabasePrefixes(e.target.value)}
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            suppressHydrationWarning
          />
        </div>

        <div>
          <input
            type="text"
            name="version"
            placeholder="Enter version"
            value={version}
            onChange={(e) => setVersion(e.target.value)}
            className="w-full p-2 base-input-field placeholder-gray-300 placeholder:text-sm rounded-sm bg-[#3f3f46]"
            suppressHydrationWarning
          />
        </div>
      </div>

      <div className="mt-4">
        <CredentialsForm onChange={handleCredentialsChange} />
      </div>
    </div>
  );
};

export default NetworkCredentialsForm;

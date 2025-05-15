"use client";

import React, { useEffect, useState, useRef } from "react";

const ModelSelectionTable = () => {
  const [models, setModels] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedModels, setSelectedModels] = useState([]);
  const originalModelsRef = useRef([]);

  // Function to fetch the data
  const fetchModelsData = async () => {
    try {
      const response = await fetch("/data.json");
      const data = await response.json();

      //Checks if the expected nested structure exists in the JSON.
      if (data && data.output && data.output.aiModelNodes) {
        const allModels = []; //Initializes an empty array to store the extracted model data.

        data.output.aiModelNodes.forEach((node) => {
          if (node.attributes && node.attributes.generativeModels) {
            const serviceProvider = node.attributes.serviceProvider || "N/A";

            node.attributes.generativeModels.forEach((model) => {
              allModels.push({
                modelName: model.modelName,
                serviceProvider: serviceProvider,
              });
            });
          }
        });

        return allModels;
      }
      return [];
    } catch (error) {
      console.error("Error fetching models data:", error);
      return [];
    }
  };

  // For loading data
  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      const modelsData = await fetchModelsData();
      originalModelsRef.current = modelsData;

      // Initialize selectedModels with all models checked
      const modelIds = Array.from({ length: modelsData.length }, (_, i) => i);
      setSelectedModels(modelIds);

      setModels(modelsData);
      setIsLoading(false);
    };

    loadData();
  }, []);

  // Function to handle checkbox change
  const handleCheckboxChange = (modelId) => {
    if (selectedModels.includes(modelId)) {
      setSelectedModels(selectedModels.filter((id) => id !== modelId));
    } else {
      setSelectedModels([...selectedModels, modelId]);
    }
  };

  // Function to handle update button click
  const handleUpdate = () => {
    const selectedModelData = originalModelsRef.current
      .filter((_, index) => selectedModels.includes(index))
      .map((model) => ({
        modelName: model.modelName,
        serviceProvider: model.serviceProvider,
      }));

    console.log("Selected models to update:", selectedModelData);

    alert(`Updating ${selectedModelData.length} models`);
  };

  return (
    <details className="border border-zinc-500 p-4 rounded mb-4 mt-2">
      <summary className="text-lg font-semibold cursor-pointer">
        Models Selection
      </summary>
      {isLoading ? (
        <div className="flex justify-center items-center h-64 text-white">
          <div className="text-xl">Loading models data...</div>
        </div>
      ) : (
        <div className="overflow-x-auto rounded-lg shadow-md text-gray-100">
          {/* Table container with fixed height for scrolling */}
          <p className="my-2 text-xs">
            <span className="text-orange-700">{selectedModels.length}</span> of{" "}
            <span className="text-orange-700">
              {originalModelsRef.current.length}
            </span>{" "}
            models selected
          </p>
          <div className="max-h-120 overflow-y-auto border-2 border-zinc-500 scrollbar">
            <table className="min-w-full divide-y divide-zinc-500 text-gray-100 text-sm">
              <thead className="bg-zinc-900 sticky top-0 ">
                <tr>
                  <th className="px-6 py-3 text-left text-sm text-gray-500 uppercase tracking-wider">
                    Select
                  </th>
                  <th className="px-6 py-3 text-left text-sm text-gray-500 uppercase tracking-wider">
                    Model Name
                  </th>
                  <th className="px-6 py-3 text-left text-sm text-gray-500 uppercase tracking-wider">
                    Service Provider
                  </th>
                </tr>
              </thead>
              <tbody className="bg-zinc-800 divide-y divide-zinc-500 break-words">
                {models.map((model, index) => (
                  <tr key={index} className="hover:bg-zinc-700">
                    <td className="px-6 py-2 whitespace-nowrap">
                      <input
                        type="checkbox"
                        className="cursor-pointer w-4 h-4 accent-orange-700"
                        checked={selectedModels.includes(index)}
                        onChange={() => handleCheckboxChange(index)}
                      />
                    </td>
                    <td className="px-6 py-2 whitespace-nowrap text-xs">
                      {model.modelName || "N/A"}
                    </td>
                    <td className="px-6 py-2 whitespace-nowrap text-xs">
                      {model.serviceProvider || "N/A"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="flex justify-end items-center  py-2 px-2 text-white">
            <button
              onClick={handleUpdate}
              className="bg-orange-700 hover:bg-orange-600 text-white py-2 px-4 rounded cursor-pointer"
            >
              Update Models
            </button>
          </div>
        </div>
      )}
    </details>
  );
};

export default ModelSelectionTable;

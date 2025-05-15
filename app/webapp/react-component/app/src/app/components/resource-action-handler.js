"use client";

import { toast } from "react-toastify";
import ToastContainerMessage from "./notification/customToast";
import { useState } from "react";
import AlertPopup from "./notification/alertPopPup";

const ActionDropdownMenu = ({
  response = {},
  globalContext = {},
  isVisible = false,
  customButton = null,
  apiServices = {},
}) => {
  const {
    resourceId,
    createActionParams,
    actionRules,
    yamlSpec,
    resourceType,
  } = response;
  const { AccessTokenKey } = globalContext || {};

  const [alertOpen, setAlertOpen] = useState(false);
  const [alertMessage, setAlertMessage] = useState("");
  const [pendingAction, setPendingAction] = useState(null);

  const slugToTitle = (slug) => {
    if (!slug) return "";
    return slug
      .split("_")
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(" ");
  };

  const handleResourceAction = async (
    specId,
    action,
    method,
    title,
    tokenKey,
    apiUrl
  ) => {
    if (
      action.toLowerCase() === "archive" ||
      action.toLowerCase() === "delete"
    ) {
      setAlertMessage("Are you sure you want to delete this resource?");
      setPendingAction(() => async () => {
        console.log("Confirmed deletion!");

        try {
          const apiService =
            apiServices[resourceType]?.[action.toLowerCase()] ||
            apiServices[resourceType]?.["delete"] ||
            apiServices[resourceType]?.["archive"];

          if (apiService && typeof apiService === "function") {
            // Use the provided API service
            const result = await apiService(resourceId);
            console.log("Resource deleted successfully:", result);
            toast.success("Resource deleted successfully");
            window.location.href = "./";
          } else {
            console.error(
              "No API service available for this resource type and action"
            );
            toast.error(
              "Unable to delete resource. No API service configured."
            );
          }
        } catch (error) {
          console.error("Error deleting resource:", error);
          toast.error("Failed to delete resource. Please try again.");
        }
      });
      setAlertOpen(true);
      return;
    } else if (
      ["validate", "publish", "unpublish"].includes(action.toLowerCase())
    ) {
      console.log(`Handling ${action} action for ${resourceId}`);
      alert(`The ${action} functionality needs to be implemented.`);
      return;
    } else if (action.toLowerCase() === "sync") {
      console.log(`Handling sync action for ${resourceId}`);
      alert("The sync functionality needs to be implemented.");
      return;
    }

    if (
      ["create", "update", "deploy", "add_users"].includes(action.toLowerCase())
    ) {
      document.getElementById("actionTitle").innerHTML = action.toUpperCase();
      document.getElementById("yamlSpecTitle").innerHTML = title;
      openYAMLedModal(apiUrl, tokenKey, specId, method);
    }
  };

  // Function to handle alert confirmation
  const handleConfirm = () => {
    if (pendingAction) {
      pendingAction();
    }
    setAlertOpen(false);
  };

  // Function to handle alert cancellation
  const handleCancel = () => {
    setAlertOpen(false);
  };

  // Function to handle custom button click
  const handleCustomButtonClick = () => {
    if (customButton && customButton.onClick) {
      customButton.onClick();
    }
  };

  const downloadYaml = (id, name) => {
    // Get the text content you want to download
    const text = document.getElementById(id).innerText;
    // Convert text to a Blob object
    const blob = new Blob([text], { type: "text/yaml" });
    // Create a link element for download
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `${name}.yaml`;
    link.click();
    URL.revokeObjectURL(link.href);
  };

  const openYAMLedModal = (apiUrl, tokenKey, specId, method) => {
    console.log("Open YAML editor:", apiUrl, tokenKey, specId, method);
  };

  // Check if we should show the default archive action
  const shouldShowDefaultArchiveAction =
    (!actionRules || actionRules.length === 0) &&
    resourceType &&
    apiServices[resourceType] &&
    (apiServices[resourceType].archive || apiServices[resourceType].delete);

  return (
    <>
      <AlertPopup
        isOpen={alertOpen}
        message={alertMessage}
        onConfirm={handleConfirm}
        onCancel={handleCancel}
      />
      <div
        id="actionDropdownMenu"
        className={`absolute right-0 mt-2 w-60 text-gray-100 bg-zinc-800 border border-zinc-500 rounded-md shadow-lg z-10 ${
          isVisible ? "" : "hidden"
        }`}
      >
        <ToastContainerMessage />
        <ul className="py-1">
          {createActionParams && (
            <li className="text-sm p-1 rounded-md">
              <a
                href={createActionParams.weblink}
                className="flex items-center px-2 py-2 hover:bg-pink-900 hover:text-white rounded-md"
              >
                <svg
                  className="w-5 h-5 mr-2"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path d="M12 2L2 7v2c0 5.25 3.25 10.17 10 15 6.75-4.83 10-9.75 10-15V7L12 2zm0 3.84l7 3.89v.95c0 3.98-2.45 8.19-7 11.57-4.55-3.38-7-7.59-7-11.57v-.95l7-3.89zm-1 3.16v5.25l4.5 2.67.75-1.23-3.75-2.22V9H11z" />
                </svg>
                Update
              </a>
            </li>
          )}

          {/* Custom Button - Only renders when customButton prop is provided */}
          {customButton && customButton.show && (
            <li className="text-sm p-1 rounded-md">
              <a
                href="#"
                onClick={handleCustomButtonClick}
                className="flex items-center px-2 py-2 hover:bg-pink-900 hover:text-white rounded-md"
              >
                <svg
                  className="w-5 h-5 mr-2"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  {customButton.icon || (
                    <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
                  )}
                </svg>
                {customButton.title || "Custom Action"}
              </a>
            </li>
          )}

          {actionRules && actionRules.length > 0
            ? actionRules.map((rule, index) => (
                <li key={index} className="text-sm p-1 rounded-md">
                  <a
                    href="#"
                    onClick={() =>
                      handleResourceAction(
                        `${resourceId}-${rule.action}`,
                        rule.action,
                        rule.method,
                        rule.title,
                        AccessTokenKey,
                        rule.api,
                        rule.isRedirect
                      )
                    }
                    className="flex items-center px-2 py-2 hover:bg-pink-900 hover:text-white rounded-md"
                  >
                    {rule.action === "Upgrade" || rule.action === "Deploy" ? (
                      <svg
                        className="w-5 h-5 mr-2"
                        fill="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path d="M12 2L2 7v2c0 5.25 3.25 10.17 10 15 6.75-4.83 10-9.75 10-15V7L12 2zm0 3.84l7 3.89v.95c0 3.98-2.45 8.19-7 11.57-4.55-3.38-7-7.59-7-11.57v-.95l7-3.89zm-1 3.16v5.25l4.5 2.67.75-1.23-3.75-2.22V9H11z" />
                      </svg>
                    ) : (
                      <svg
                        className="w-5 h-5 mr-2"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path d="M15.232 4.232a1 1 0 0 1 1.415 0l3.122 3.122a1 1 0 0 1 0 1.415l-11.3 11.3a1 1 0 0 1-.707.293H5a1 1 0 0 1-1-1v-2.829a1 1 0 0 1 .293-.707l11.3-11.3zM16.646 2.818L17.788 3.96l-10.464 10.464-1.142-1.142L16.646 2.818zM3 21h18v2H3v-2z" />
                      </svg>
                    )}
                    {slugToTitle(rule.action)}
                  </a>
                  <div hidden id={`${resourceId}-${rule.action}`}>
                    {rule.yamlSpec}
                  </div>
                </li>
              ))
            : null}

          {/* Default archive action if no actions are defined but API services exist for this resource type */}
          {shouldShowDefaultArchiveAction && (
            <li className="text-sm p-1 rounded-md">
              <a
                href="#"
                onClick={() =>
                  handleResourceAction(
                    `${resourceId}-archive`,
                    "archive",
                    "DELETE",
                    "Archive Resource",
                    AccessTokenKey,
                    "",
                    false
                  )
                }
                className="flex items-center px-2 py-2 hover:bg-pink-900 hover:text-white rounded-md"
              >
                <svg
                  className="w-5 h-5 mr-2"
                  xmlns="http://www.w3.org/2000/svg"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z" />
                </svg>
                Archive
              </a>
            </li>
          )}

          {yamlSpec && (
            <li className="text-sm p-1 rounded-md">
              <div hidden id="yamlSpec">
                {yamlSpec}
              </div>
              <a
                href="#"
                onClick={() => downloadYaml("yamlSpec", resourceId)}
                className="flex items-center px-2 py-2 hover:bg-pink-900 hover:text-white rounded-lg"
              >
                <svg
                  className="w-6 h-6 mr-2"
                  viewBox="0 0 24 24"
                  fill="white"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    d="M5 20h14v2H5v-2zm7-2c-.28 0-.53-.11-.71-.29L8 13.41l1.41-1.41L11 14.17V4h2v10.17l1.59-1.59L16 13.41l-3.29 3.29c-.18.18-.43.29-.71.29z"
                    fill="#fff"
                  />
                </svg>
                Download Yaml Spec
              </a>
            </li>
          )}
        </ul>
      </div>
    </>
  );
};

export default ActionDropdownMenu;

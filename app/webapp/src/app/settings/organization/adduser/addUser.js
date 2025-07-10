"use client";
import { useState, useEffect } from "react";
import RemoveButton from "@/app/components/buttons/removeButton";
import AddButton from "@/app/components/buttons/addButton";
import { toast } from "react-toastify";
import { addUsersApi } from "@/app/utils/settings-endpoint/organization-api";
import ToastContainerMessage from "@/app/components/notification/customToast";
import MultiSelectDropdown from "@/app/components/multiSelectDropdown";
import { enumsApi } from "@/app/utils/developers-endpoint/enums";

// SVG Icon Components
const XIcon = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="24"
    height="24"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
  >
    <path d="M18 6L6 18"></path>
    <path d="M6 6L18 18"></path>
  </svg>
);

// Add organizationId as a prop
export default function AddUserModal({ isOpen, onClose, organizationId }) {
  const [users, setUsers] = useState([
    { userId: "", roles: [], inviteIfNotFound: true }, // Changed default to true
  ]);
  const [isLoading, setIsLoading] = useState(false);
  const [roleOptions, setRoleOptions] = useState([]);
  const [loadingRoles, setLoadingRoles] = useState(false);

  // Fetch the OrgRoles data
  useEffect(() => {
    const fetchEnumsData = async () => {
      try {
        setLoadingRoles(true);
        const response = await enumsApi.getEnums();
        console.log("Full enums response:", response);

        const OrgRoles = response.enumResponse?.find(
          (item) => item.name === "OrgRoles"
        );

        console.log("Found OrgRoles:", OrgRoles);

        const roleOptions =
          OrgRoles?.value?.map((role) => ({
            id: role,
            name: role,
          })) || [];

        setRoleOptions(roleOptions);
        console.log("Processed role options:", roleOptions);
      } catch (error) {
        console.error("Failed to fetch User Roles data:", error);
        toast.error("Failed to fetch configuration data");
      } finally {
        setLoadingRoles(false);
      }
    };

    if (isOpen) {
      fetchEnumsData();
    }
  }, [isOpen]);

  const addUserParams = () => {
    setUsers([...users, { userId: "", roles: [], inviteIfNotFound: true }]);
  };

  const removeUserParams = (index) => {
    const updatedUsers = [...users];
    updatedUsers.splice(index, 1);
    setUsers(updatedUsers);
  };

  const handleUserChange = (index, field, value) => {
    const updatedUsers = [...users];
    updatedUsers[index][field] = value;
    setUsers(updatedUsers);
  };

  const handleRoleSelectionChange = (index, selectedRoles) => {
    console.log("Selected roles received:", selectedRoles);
    const updatedUsers = [...users];

    // Convert selected role objects back to string array for API
    const roleStrings = selectedRoles.map((role) =>
      typeof role === "string" ? role : role.id || role.name
    );

    updatedUsers[index].roles = roleStrings;
    setUsers(updatedUsers);
    console.log("Updated user roles:", updatedUsers[index].roles);
  };

  const handleSubmit = async () => {
    // Add validation for organizationId
    if (!organizationId) {
      toast.error("Organization ID is required");
      return;
    }

    // Validate users data before submission
    const validUsers = users.filter(
      (user) => user.userId.trim() !== "" && user.roles.length > 0
    );

    if (validUsers.length === 0) {
      toast.error("Please enter valid user information with at least one role");
      return;
    }

    try {
      setIsLoading(true);
      const payload = {
        organizationId: organizationId,
        users: validUsers.map((user) => ({
          userId: user.userId,
          validTill: "0", // Added required field
          role: user.roles, // Array of role strings
          added: false, // Added required field
          inviteIfNotFound: user.inviteIfNotFound,
        })),
      };

      console.log("Full API Payload:", payload);
      const output = await addUsersApi.getAddUsers(payload);
      const resourceInfo = output?.result || {
        resource: "User",
        count: validUsers.length,
      };

      // Show toast with resource info
      toast.success(
        `${resourceInfo.resource} Added`,
        `${resourceInfo.count} ${resourceInfo.resource} User(s) added successfully.`
      );

      // Don't close modal automatically to allow success toast to be visible
      // User can manually close the modal after seeing the success message
    } catch (error) {
      console.error("Error adding users:", error);
      toast.error("Failed to add users");
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-zinc-600/80 bg-opacity-50 flex items-center justify-center p-4 z-50">
      <ToastContainerMessage />
      <div className="bg-[#1b1b1b] rounded-lg shadow-lg w-3/4 h-1/2.5 md:w-1/2 overflow-y-auto scrollbar text-gray-100">
        <div className="bg-zinc-900 px-4 py-2 flex justify-between items-center">
          <h2 className="text-xl font-semibold text-white">Add User</h2>
          <button
            onClick={onClose}
            className="text-zinc-400 hover:text-white cursor-pointer"
          >
            <XIcon />
          </button>
        </div>

        <div className="p-4 max-h-96 overflow-y-auto scrollbar">
          {users.map((user, index) => (
            <div
              key={index}
              className="mb-6 p-4 border border-zinc-700 rounded-md shadow-sm "
            >
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {/* User Id */}
                <div className="mb-4">
                  <label className="block mb-2 text-sm font-medium">
                    User Id
                  </label>
                  <input
                    type="text"
                    className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm p-1 mt-2 w-full"
                    placeholder="Enter userId"
                    value={user.userId}
                    onChange={(e) =>
                      handleUserChange(index, "userId", e.target.value)
                    }
                  />
                </div>

                {/* MultiSelectDropdown */}
                <div className="mb-4">
                  <label className="block mb-2 text-sm font-medium">Role</label>
                  <MultiSelectDropdown
                    options={roleOptions}
                    placeholder="Select roles"
                    onSelectionChange={(selectedRoles) =>
                      handleRoleSelectionChange(index, selectedRoles)
                    }
                    initialSelected={user.roles.map((role) => ({
                      id: role,
                      name: role,
                    }))}
                    isLoading={loadingRoles}
                    dropdownId={`dropdown-${index}`}
                  />
                </div>
              </div>

              <div className="mb-4 flex items-center">
                <input
                  type="checkbox"
                  id={`invite-${index}`}
                  className="w-4 h-4 mr-2 accent-orange-700 cursor-pointer"
                  checked={user.inviteIfNotFound}
                  onChange={(e) =>
                    handleUserChange(
                      index,
                      "inviteIfNotFound",
                      e.target.checked
                    )
                  }
                />
                <label htmlFor={`invite-${index}`} className="text-sm">
                  Invite if not registered
                </label>
              </div>

              {users.length > 1 && (
                <RemoveButton onClick={() => removeUserParams(index)} />
              )}
            </div>
          ))}

          <AddButton name="+ Add Another User" onClick={addUserParams} />
        </div>

        <div className="flex justify-end gap-3 p-4 border-t border-zinc-700">
          <button
            onClick={onClose}
            className="px-4 py-2 text-white bg-zinc-600 rounded-md cursor-pointer"
            disabled={isLoading}
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            className="px-6 py-2 bg-orange-700 text-white rounded-md shadow hover:bg-pink-900 cursor-pointer"
            disabled={isLoading}
          >
            {isLoading ? (
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
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
            ) : (
              "Submit"
            )}
          </button>
        </div>
      </div>
    </div>
  );
}

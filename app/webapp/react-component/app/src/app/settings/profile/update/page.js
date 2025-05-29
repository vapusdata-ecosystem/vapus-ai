"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import Header from "@/app/components/platform/header";
import YamlEditorClient from "@/app/components/formcomponets/ymal";
import ToastContainerMessage from "@/app/components/notification/customToast";
import LoadingOverlay from "@/app/components/loading/loading";
import {
  updateProfile,
  userProfileApi,
} from "@/app/utils/settings-endpoint/profile-api";
import AddButton from "@/app/components/buttons/addButton";
import RemoveButton from "@/app/components/buttons/removeButton";
import { userGlobalData } from "@/context/GlobalContext";
import { UploadFileApi } from "@/app/utils/file-endpoint/file";

export default function UpdateProfile() {
  const router = useRouter();
  const [activeTab, setActiveTab] = useState("form");
  const [isLoading, setIsLoading] = useState(true);
  const [userData, setUserData] = useState(null);
  const [domainMap, setDomainMap] = useState({});
  const [contextData, setContextData] = useState(null);
  const [formData, setFormData] = useState({
    displayName: "",
    avatar: "",
    addresses: [
      {
        street_address1: "",
        street_address2: "",
        city: "",
        state: "",
        zip_code: "",
        country: "",
        others: "",
      },
    ],
  });
  const [addressEntryCount, setAddressEntryCount] = useState(1);
  const [avatarFile, setAvatarFile] = useState(null);
  const [avatarPreview, setAvatarPreview] = useState("");
  const [uploadedAvatarPath, setUploadedAvatarPath] = useState("");

  // Fetch initial data and populate form
  useEffect(() => {
    const fetchData = async () => {
      try {
        setIsLoading(true);

        // Get user data from global context
        const globalContext = await userGlobalData();
        setContextData(globalContext);
        console.log("my data", globalContext);

        // Check if userId exists
        if (globalContext?.userInfo?.userId) {
          const userId = globalContext.userInfo.userId;
          console.log("User ID:", userId);

          // Make API call to get user profile with userId
          const data = await userProfileApi.getuserProfile(userId);
          console.log("data", data);

          const userInfo = data.output.users[0];
          setUserData(userInfo);
          setDomainMap(data.domainMap);

          // Populate form data with fetched user data
          populateFormData(userInfo);
        } else {
          console.error("User ID not found in global context");
        }
      } catch (error) {
        console.error("Error fetching user data:", error);
        toast.error("Failed to load user profile data");
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  // Function to populate form data from API response
  const populateFormData = (userInfo) => {
    // Set basic form data
    const newFormData = {
      displayName: userInfo.displayName || "",
      avatar: userInfo.profile?.avatar || "",
      addresses: [],
    };

    // Handle addresses
    if (userInfo.profile?.addresses && userInfo.profile.addresses.length > 0) {
      // If user has existing addresses, use them
      newFormData.addresses = userInfo.profile.addresses.map((addr) => ({
        street_address1: addr.street_address1 || "",
        street_address2: addr.street_address2 || "",
        city: addr.city || "",
        state: addr.state || "",
        zip_code: addr.zip_code || "",
        country: addr.country || "",
        others: addr.others || "",
      }));
      setAddressEntryCount(userInfo.profile.addresses.length);
    } else {
      // If no addresses exist, keep the default empty address
      newFormData.addresses = [
        {
          street_address1: "",
          street_address2: "",
          city: "",
          state: "",
          zip_code: "",
          country: "",
          others: "",
        },
      ];
      setAddressEntryCount(1);
    }

    // Set form data
    setFormData(newFormData);

    // Set avatar preview and uploaded path if avatar exists
    if (userInfo.profile?.avatar) {
      setAvatarPreview(userInfo.profile.avatar);
      setUploadedAvatarPath(userInfo.profile.avatar); // Store existing avatar path
    }
  };

  // Handle input changes
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    const fieldName = name.replace("spec.", "");

    setFormData((prevData) => ({
      ...prevData,
      [fieldName]: value,
    }));
  };

  //  avatar file upload
  const handleAvatarChange = async (e) => {
    const file = e.target.files[0];
    if (file) {
      // Validate file type
      const validTypes = ["image/jpeg", "image/jpg", "image/png", "image/gif"];
      if (!validTypes.includes(file.type)) {
        toast.error("Please select a valid image file (JPEG, PNG, GIF)");
        return;
      }

      // Validate file size (5MB limit)
      const maxSize = 5 * 1024 * 1024; // 5MB
      if (file.size > maxSize) {
        toast.error("File size must be less than 5MB");
        return;
      }

      setAvatarFile(file);

      // Create preview URL
      const reader = new FileReader();
      reader.onload = async (e) => {
        const base64Data = e.target.result;
        setAvatarPreview(base64Data);

        // Call upload API when image is selected
        await uploadAvatarImage(file, base64Data);
      };
      reader.readAsDataURL(file);
    }
  };

  // Function to upload avatar image
  const uploadAvatarImage = async (file, base64Data) => {
    try {
      setIsLoading(true);

      // Get file extension
      const fileExtension = file.name.split(".").pop().toLowerCase();
      const base64String = base64Data.split(",")[1];
      const userEmail = userData?.email || contextData?.userInfo?.email || "";

      if (!userEmail) {
        toast.error("User email not found. Cannot upload image.");
        return;
      }

      const uploadPayload = {
        resourceId: userEmail,
        resource: "Profile",
        objects: [
          {
            name: file.name,
            data: base64String,
            format: fileExtension,
          },
        ],
      };

      const uploadResponse = await UploadFileApi.getUploadFile(uploadPayload);
      console.log("Avatar upload response:", uploadResponse);

      // Extract responsePath from the upload response
      if (uploadResponse.output && uploadResponse.output.length > 0) {
        const responsePath = uploadResponse.output[0].responsePath;

        // Store the uploaded image path
        setUploadedAvatarPath(responsePath);

        // Update form data with the response path instead of base64
        setFormData((prevData) => ({
          ...prevData,
          avatar: responsePath, // Use responsePath instead of base64Data
        }));

        toast.success("Avatar uploaded successfully!");
      } else {
        throw new Error("Invalid upload response format");
      }
    } catch (error) {
      console.error("Error uploading avatar:", error);
      toast.error(error.message || "Failed to upload avatar");

      // Reset avatar on error
      setAvatarPreview("");
      setAvatarFile(null);
      setUploadedAvatarPath(""); // Reset the uploaded path
      setFormData((prevData) => ({
        ...prevData,
        avatar: "",
      }));
      // Clear the file input
      document.getElementById("avatar").value = "";
    } finally {
      setIsLoading(false);
    }
  };

  // Handle address changes
  const handleAddressChange = (index, field, value) => {
    const newAddresses = [...formData.addresses];
    newAddresses[index] = {
      ...newAddresses[index],
      [field]: value,
    };
    setFormData((prevData) => ({
      ...prevData,
      addresses: newAddresses,
    }));
  };

  // Add new address entry
  const addAddressEntry = () => {
    setFormData((prevData) => ({
      ...prevData,
      addresses: [
        ...prevData.addresses,
        {
          street_address1: "",
          street_address2: "",
          city: "",
          state: "",
          zip_code: "",
          country: "",
          others: "",
        },
      ],
    }));
    setAddressEntryCount((prev) => prev + 1);
  };

  // Remove address entry
  const removeAddressEntry = (index) => {
    if (formData.addresses.length > 1) {
      const newAddresses = formData.addresses.filter((_, i) => i !== index);
      setFormData((prevData) => ({
        ...prevData,
        addresses: newAddresses,
      }));
      setAddressEntryCount((prev) => prev - 1);
    }
  };

  // Handle remove avatar
  const handleRemoveAvatar = () => {
    setAvatarPreview("");
    setAvatarFile(null);
    setUploadedAvatarPath("");
    setFormData((prevData) => ({
      ...prevData,
      avatar: "",
    }));
    // Clear the file input
    document.getElementById("avatar").value = "";
  };

  // Handle form submission
  const submitCreateForm = async (e) => {
    e.preventDefault();

    if (!formData.displayName.trim()) {
      toast.error("Please fill in the display name");
      return;
    }

    // Filter out empty addresses (check if at least one field is filled)
    const filteredAddresses = formData.addresses.filter((addr) =>
      Object.values(addr).some((value) => value.trim() !== "")
    );

    try {
      setIsLoading(true);
      const payload = {
        action: "PATCH_USER",
        spec: {
          displayName: formData.displayName,
          userId: userData?.userId,
          profile: {
            avatar: formData.avatar || uploadedAvatarPath || "",
            addresses: filteredAddresses,
            description: userData?.profile?.description || "",
          },
        },
      };

      console.log("Submitting profile update data:", payload);
      // Call API
      const response = await updateProfile.getUpdateProfile(payload);

      console.log("Profile updated:", response);
      toast.success("Profile updated successfully!");
      router.push("./");
    } catch (error) {
      console.error("Error updating profile:", error);
      toast.error(error.message || "Failed to update profile");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="bg-zinc-800 flex h-screen">
      <div className="overflow-y-auto scrollbar h-screen w-full">
        <Header
          sectionHeader="Update Profile"
          hideBackListingLink={false}
          backListingLink="./"
        />
        <ToastContainerMessage />

        <LoadingOverlay isLoading={isLoading} />
        <div className="flex-grow p-4 overflow-y-auto w-full">
          <section className="space-y-2">
            <div className="max-w-6xl mx-auto bg-[#1b1b1b] shadow rounded-lg p-2">
              <div className="text-gray-100 mb-2 flex justify-center">
                <button
                  onClick={() => setActiveTab("form")}
                  className={`whitespace-nowrap border-b-2 py-2 px-2 text-md font-medium focus:outline-none ml-4 ${
                    activeTab === "form"
                      ? "border-orange-700 text-orange-700 font-semibold"
                      : ""
                  }`}
                  suppressHydrationWarning
                >
                  Form
                </button>
              </div>

              {activeTab === "yaml" && (
                <div id="yamlSpec">
                  <YamlEditorClient />
                </div>
              )}

              {activeTab === "form" && (
                <div id="formSpec">
                  <form
                    id="dataSourceSpec"
                    className="space-y-4 border border-zinc-500 rounded-md text-gray-100 p-4"
                    onSubmit={submitCreateForm}
                  >
                    <fieldset className="p-4 rounded">
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
                        {/* Avatar Upload */}
                        <div>
                          <label className="labels block text-sm font-medium mb-4">
                            Profile Avatar
                          </label>
                          <div className="flex flex-col items-center space-y-3">
                            {/* Avatar Preview */}
                            <div className="relative group">
                              <div className="w-32 h-32 rounded-full border-4 border-gray-600 overflow-hidden bg-zinc-700 flex items-center justify-center">
                                {avatarPreview ? (
                                  <img
                                    src={avatarPreview}
                                    alt="Avatar Preview"
                                    className="w-full h-full object-cover"
                                  />
                                ) : (
                                  <div className="text-gray-400 text-center">
                                    <svg
                                      className="w-12 h-12 mx-auto mb-2"
                                      fill="currentColor"
                                      viewBox="0 0 20 20"
                                    >
                                      <path
                                        fillRule="evenodd"
                                        d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                                        clipRule="evenodd"
                                      />
                                    </svg>
                                    <span className="text-xs">No Image</span>
                                  </div>
                                )}
                              </div>

                              {/* Edit Button Overlay */}
                              <div className="absolute inset-0 rounded-full bg-black bg-opacity-50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200 cursor-pointer">
                                <label
                                  htmlFor="avatar"
                                  className="cursor-pointer text-white text-sm font-medium px-3 py-1 bg-orange-700 rounded-full hover:bg-orange-600 transition-colors"
                                >
                                  Edit
                                </label>
                              </div>
                            </div>

                            {/* Hidden File Input */}
                            <input
                              id="avatar"
                              name="avatar"
                              type="file"
                              accept="image/*"
                              onChange={handleAvatarChange}
                              className="hidden"
                            />

                            {/* Upload Instructions */}
                            <p className="text-xs text-gray-400 text-center max-w-xs">
                              Click the circle to upload a new avatar.
                              <br />
                              Supported formats: JPEG, PNG, GIF. Max size: 5MB
                            </p>

                            {/* Remove Avatar Button (only show if avatar exists) */}
                            {avatarPreview && (
                              <button
                                type="button"
                                onClick={handleRemoveAvatar}
                                className="text-red-400 hover:text-red-300 text-sm underline cursor-pointer"
                              >
                                Remove Avatar
                              </button>
                            )}
                          </div>
                        </div>

                        {/* Display Name */}
                        <div>
                          <label
                            htmlFor="spec_displayName"
                            className="labels block text-sm font-medium mb-2"
                          >
                            Display Name *
                          </label>
                          <input
                            id="spec_displayName"
                            name="spec.displayName"
                            type="text"
                            placeholder="Enter Display name"
                            className="form-input-field p-3 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm w-full"
                            value={formData.displayName}
                            onChange={handleInputChange}
                            required
                            suppressHydrationWarning
                          />
                        </div>
                      </div>

                      {/* Addresses Section */}
                      <details className="shadow-sm border border-zinc-500 rounded-md shadow-sm p-4 mb-4">
                        <summary className="text-lg font-semibold text-gray-100 cursor-pointer">
                          Addresses
                        </summary>
                        <fieldset className="rounded mb-4">
                          <div id="addressesContainer">
                            {formData.addresses.map((address, index) => (
                              <div
                                key={`address-${index}`}
                                className="address-entry border border-zinc-600 p-4 rounded mb-4"
                              >
                                <h4 className="text-md font-medium text-gray-200 mb-3">
                                  Address {index + 1}
                                </h4>

                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                  {/* Street Address 1 */}
                                  <div>
                                    <label className="labels block text-sm font-medium mb-1">
                                      Street Address 1
                                    </label>
                                    <input
                                      type="text"
                                      name={`spec.addresses[${index}].street_address1`}
                                      placeholder="Enter street address 1"
                                      className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm mt-1 w-full"
                                      value={address.street_address1 || ""}
                                      onChange={(e) =>
                                        handleAddressChange(
                                          index,
                                          "street_address1",
                                          e.target.value
                                        )
                                      }
                                      suppressHydrationWarning
                                    />
                                  </div>

                                  {/* Street Address 2 */}
                                  <div>
                                    <label className="labels block text-sm font-medium mb-1">
                                      Street Address 2
                                    </label>
                                    <input
                                      type="text"
                                      name={`spec.addresses[${index}].street_address2`}
                                      placeholder="Enter street address 2 (optional)"
                                      className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm mt-1 w-full"
                                      value={address.street_address2 || ""}
                                      onChange={(e) =>
                                        handleAddressChange(
                                          index,
                                          "street_address2",
                                          e.target.value
                                        )
                                      }
                                      suppressHydrationWarning
                                    />
                                  </div>

                                  {/* City */}
                                  <div>
                                    <label className="labels block text-sm font-medium mb-1">
                                      City
                                    </label>
                                    <input
                                      type="text"
                                      name={`spec.addresses[${index}].city`}
                                      placeholder="Enter city"
                                      className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm mt-1 w-full"
                                      value={address.city || ""}
                                      onChange={(e) =>
                                        handleAddressChange(
                                          index,
                                          "city",
                                          e.target.value
                                        )
                                      }
                                      suppressHydrationWarning
                                    />
                                  </div>

                                  {/* State */}
                                  <div>
                                    <label className="labels block text-sm font-medium mb-1">
                                      State
                                    </label>
                                    <input
                                      type="text"
                                      name={`spec.addresses[${index}].state`}
                                      placeholder="Enter state"
                                      className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm mt-1 w-full"
                                      value={address.state || ""}
                                      onChange={(e) =>
                                        handleAddressChange(
                                          index,
                                          "state",
                                          e.target.value
                                        )
                                      }
                                      suppressHydrationWarning
                                    />
                                  </div>

                                  {/* Zip Code */}
                                  <div>
                                    <label className="labels block text-sm font-medium mb-1">
                                      Zip Code
                                    </label>
                                    <input
                                      type="text"
                                      name={`spec.addresses[${index}].zip_code`}
                                      placeholder="Enter zip code"
                                      className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm mt-1 w-full"
                                      value={address.zip_code || ""}
                                      onChange={(e) =>
                                        handleAddressChange(
                                          index,
                                          "zip_code",
                                          e.target.value
                                        )
                                      }
                                      suppressHydrationWarning
                                    />
                                  </div>

                                  {/* Country */}
                                  <div>
                                    <label className="labels block text-sm font-medium mb-1">
                                      Country
                                    </label>
                                    <input
                                      type="text"
                                      name={`spec.addresses[${index}].country`}
                                      placeholder="Enter country"
                                      className="form-input-field p-2 rounded-md bg-zinc-700 placeholder-gray-300 shadow-sm outline-none border-none text-sm mt-1 w-full"
                                      value={address.country || ""}
                                      onChange={(e) =>
                                        handleAddressChange(
                                          index,
                                          "country",
                                          e.target.value
                                        )
                                      }
                                      suppressHydrationWarning
                                    />
                                  </div>
                                </div>

                                {/* Others - Full Width */}
                                <div className="mt-4">
                                  <label className="labels block text-sm font-medium mb-1">
                                    Others
                                  </label>
                                  <textarea
                                    name={`spec.addresses[${index}].others`}
                                    placeholder="Enter additional address information"
                                    rows="3"
                                    className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
                                    value={address.others || ""}
                                    onChange={(e) =>
                                      handleAddressChange(
                                        index,
                                        "others",
                                        e.target.value
                                      )
                                    }
                                    suppressHydrationWarning
                                  />
                                </div>

                                {/* Remove Button */}
                                {formData.addresses.length > 1 && (
                                  <RemoveButton
                                    onClick={() => removeAddressEntry(index)}
                                  />
                                )}
                              </div>
                            ))}
                          </div>

                          {/* Add Address Button */}
                          <AddButton
                            name="+ Add Address"
                            onClick={addAddressEntry}
                          />
                        </fieldset>
                      </details>

                      {/* Submit Button */}
                      <div className="mt-8 flex justify-end">
                        <button
                          type="submit"
                          disabled={isLoading}
                          className={`px-8 py-3 bg-orange-700 text-white rounded-md shadow-md font-medium ${
                            isLoading
                              ? "opacity-50 cursor-not-allowed"
                              : "hover:bg-orange-600 transition-colors"
                          }`}
                          suppressHydrationWarning
                        >
                          {isLoading ? "Updating..." : "Update Profile"}
                        </button>
                      </div>
                    </fieldset>
                  </form>
                </div>
              )}
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}

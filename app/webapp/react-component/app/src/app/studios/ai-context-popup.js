"use client";
import { useState, useEffect } from "react";
import { toast } from "react-toastify";
import { UploadFileApi } from "@/app/utils/file-endpoint/file";
import { userGlobalData } from "@/context/GlobalContext";
import {
  userProfileApi,
} from "@/app/utils/settings-endpoint/profile-api";

export default function ContextModal({ isOpen, onClose }) {
  const [uploadedFiles, setUploadedFiles] = useState([]);
  const [isUploading, setIsUploading] = useState(false);
  const [userData, setUserData] = useState(null);
  const [contextData, setContextData] = useState(null);

  // Fetch user data when modal opens
  useEffect(() => {
    const fetchData = async () => {
      if (isOpen) {
        try {
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
          } else {
            console.error("User ID not found in global context");
          }
        } catch (error) {
          console.error("Error fetching user data:", error);
          toast.error("Failed to load user profile data");
        }
      }
    };

    fetchData();
  }, [isOpen]);

    const userEmail = userData?.email || contextData?.userInfo?.email || "";


  const uploadFromComputer = () => {

    const fileInput = document.createElement("input");
    fileInput.type = "file";
    fileInput.multiple = true;
    fileInput.accept = "*/*"; // Accept all file types

    fileInput.onchange = async (event) => {
      const files = event.target.files;
      if (files && files.length > 0) {
        console.log("Selected files:", files);
        
        // Process each file
        for (let file of files) {
          await handleFileUpload(file);
        }
      }
    };

    fileInput.click();
  };

  const handleFileUpload = async (file) => {
    // Validate file size (50MB limit - increased for various file types)
    const maxSize = 50 * 1024 * 1024;
    if (file.size > maxSize) {
      toast.error("File size must be less than 50MB");
      return;
    }

    // Get file type and create appropriate preview
    const isImage = file.type.startsWith('image/');
    
    if (isImage) {
      // Handle image files with preview
      const reader = new FileReader();
      reader.onload = async (e) => {
        const base64Data = e.target.result;
        
        const filePreview = {
          id: Date.now() + Math.random(), 
          name: file.name,
          type: file.type,
          size: file.size,
          preview: base64Data,
          file: file,
          uploaded: false,
          uploadPath: null,
          isImage: true
        };
        
        setUploadedFiles(prev => [...prev, filePreview]);
        await uploadFileToAPI(file, base64Data, filePreview.id);
      };
      reader.readAsDataURL(file);
    } else {
      // Handle non-image files
      const reader = new FileReader();
      reader.onload = async (e) => {
        const base64Data = e.target.result;
        
        const filePreview = {
          id: Date.now() + Math.random(), 
          name: file.name,
          type: file.type,
          size: file.size,
          file: file,
          uploaded: false,
          uploadPath: null,
          isImage: false
        };
        
        setUploadedFiles(prev => [...prev, filePreview]);
        await uploadFileToAPI(file, base64Data, filePreview.id);
      };
      reader.readAsDataURL(file);
    }
  };

  const uploadFileToAPI = async (file, base64Data, fileId) => {
    try {
      setIsUploading(true);
      console.log("Starting upload API call for:", file.name);

      // Check if user email is available
      if (!userEmail) {
        throw new Error("User email not found. Cannot upload file.");
      }

      // Get file extension
      const fileExtension = file.name.split(".").pop().toLowerCase();
      const base64String = base64Data.split(",")[1];

      const uploadPayload = {
        resourceId: userEmail,
        resource: "ContextModal",
        objects: [
          {
            name: file.name,
            data: base64String,
            format: fileExtension,
          },
        ],
      };

      console.log("Upload API payload:", uploadPayload);
      
      const uploadResponse = await UploadFileApi.getUploadFile(uploadPayload);
      console.log("Upload API response:", uploadResponse);

      // Update the file status
      if (uploadResponse.output && uploadResponse.output.length > 0) {
        const responsePath = uploadResponse.output[0].responsePath;
        
        setUploadedFiles(prev => 
          prev.map(file => 
            file.id === fileId 
              ? { ...file, uploaded: true, uploadPath: responsePath }
              : file
          )
        );

        console.log("File uploaded successfully:", responsePath);
        toast.success(`${file.name} uploaded successfully!`);
      } else {
        throw new Error("Invalid upload response format");
      }
    } catch (error) {
      console.error("Error uploading file:", error);
      toast.error(`Failed to upload ${file.name}: ${error.message}`);
      
      // Mark as failed
      setUploadedFiles(prev => 
        prev.map(file => 
          file.id === fileId 
            ? { ...file, uploaded: false, error: true }
            : file
        )
      );
    } finally {
      setIsUploading(false);
    }
  };

  const removeFile = (fileId) => {
    setUploadedFiles(prev => prev.filter(file => file.id !== fileId));
  };

  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const getFileIcon = (fileType) => {
    if (fileType.startsWith('image/')) {
      return '🖼️';
    } else if (fileType.includes('pdf')) {
      return '📄';
    } else if (fileType.includes('word') || fileType.includes('document')) {
      return '📝';
    } else if (fileType.includes('excel') || fileType.includes('spreadsheet')) {
      return '📊';
    } else if (fileType.includes('powerpoint') || fileType.includes('presentation')) {
      return '📽️';
    } else if (fileType.startsWith('video/')) {
      return '🎥';
    } else if (fileType.startsWith('audio/')) {
      return '🎵';
    } else if (fileType.includes('zip') || fileType.includes('rar') || fileType.includes('tar')) {
      return '🗜️';
    } else if (fileType.includes('text/')) {
      return '📄';
    } else {
      return '📎';
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-zinc-600/90 flex items-center justify-center z-50">
      <div className="bg-zinc-800 rounded-lg shadow-lg p-6 w-5/6 h-5/6 overflow-auto scrollbar text-gray-100">
        <h3 className="text-lg font-semibold mb-4">Additional Params</h3>
        {/* Options */}
        <ul className="divide-y divide-gray-200 mb-4">
          {/* Upload from Computer */}
          <li
            className={`flex items-center px-4 py-2 cursor-pointer ${
              userEmail 
                ? 'hover:bg-zinc-900' 
                : 'opacity-50 cursor-not-allowed'
            }`}
            onClick={uploadFromComputer}
          >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="w-6 h-6 mr-2"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path d="M3 4a2 2 0 012-2h10a2 2 0 012 2v8a2 2 0 01-2 2h-4v2h2a1 1 0 110 2H7a1 1 0 110-2h2v-2H5a2 2 0 01-2-2V4zm2 0v8h10V4H5z"/>
                <path d="M10 6a1 1 0 011 1v1.586l.293-.293a1 1 0 011.414 1.414l-2 2a1 1 0 01-1.414 0l-2-2a1 1 0 111.414-1.414L9 8.586V7a1 1 0 011-1z"/>
              </svg>
              <span className="text-sm">
                Upload from Computer (All File Types)
                {!userEmail && " - Loading user data..."}
              </span>
            {isUploading && (
              <div className="ml-2">
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-orange-500"></div>
              </div>
            )}
          </li>
        </ul>

        {/* Uploaded Files Display */}
        {uploadedFiles.length > 0 && (
          <div className="mb-4">
            <h4 className="text-md font-medium text-gray-200 mb-3">Uploaded Files:</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {uploadedFiles.map((file) => (
                <div key={file.id} className="relative group">
                  <div className="rounded-lg overflow-hidden bg-zinc-700 border-2 border-zinc-600 p-4">
                    {file.isImage ? (
                      // Image preview
                      <div className="aspect-square mb-2">
                        <img
                          src={file.preview}
                          alt={file.name}
                          className="w-full h-full object-cover rounded"
                        />
                      </div>
                    ) : (
                      // File icon for non-images
                      <div className="aspect-square mb-2 flex items-center justify-center text-6xl">
                        {getFileIcon(file.type)}
                      </div>
                    )}
                    
                    {/* File Info */}
                    <div className="text-white">
                      <div className="text-sm font-medium truncate" title={file.name}>
                        {file.name}
                      </div>
                      <div className="text-xs text-gray-300 mt-1">
                        {formatFileSize(file.size)}
                      </div>
                      <div className="flex items-center justify-between mt-2">
                        <span className={`text-xs ${
                          file.uploaded 
                            ? 'text-green-400' 
                            : file.error 
                              ? 'text-red-400' 
                              : 'text-yellow-400'
                        }`}>
                          {file.uploaded 
                            ? '✓ Uploaded' 
                            : file.error 
                              ? '✗ Failed' 
                              : '⏳ Uploading...'}
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* Remove Button */}
                  <button
                    onClick={() => removeFile(file.id)}
                    className="absolute cursor-pointer top-2 right-2 bg-red-500 hover:bg-red-600 text-white rounded-full w-6 h-6 flex items-center justify-center text-xs opacity-0 group-hover:opacity-100 transition-opacity"
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* System Message */}
        <div className="mt-4 p-4 rounded-lg">
          <label
            htmlFor="systemMessage"
            className="block text-sm font-medium"
          >
            System Message
          </label>
          <input
            type="text"
            id="systemMessage"
            placeholder="Enter your system message here..."
            className="w-full px-3 py-2 border border-zinc-500 rounded-lg focus:outline-none focus:ring focus:ring-[#ff5f1f] mb-4 bg-zinc-700 text-gray-100"
          />
        </div>

        {/* Context Type */}
        <div className="mb-4">
          <label htmlFor="contextType" className="block text-sm font-medium">
            Context Type
          </label>
          <input
            type="text"
            id="contextType"
            className="w-full p-2 border border-zinc-500 rounded-lg focus:outline-none focus:ring focus:ring-[#ff5f1f] bg-zinc-700 text-gray-100"
            placeholder="Enter RAG context type here..."
          />
        </div>

        {/* Context Value (Always Visible) */}
        <div className="mb-4">
          <label htmlFor="contextValue" className="block text-sm font-medium">
            Context Value
          </label>
          <textarea
            id="contextValue"
            className="form-textarea overflow-y-auto scrollbar text-sm rounded-md bg-zinc-700 placeholder-gray-300 p-2 shadow-sm w-full"
            rows="13"
            placeholder="Enter RAG context here..."
          ></textarea>
        </div>

        {/* Close Button */}
        <div className="flex justify-end mt-4">
          <button
            className="bg-red-500 text-gray-100 px-4 py-2 rounded-lg hover:bg-red-600 focus:outline-none focus:ring focus:ring-red-300"
            onClick={onClose}
          >
            Done
          </button>
        </div>
      </div>
    </div>
  );
}
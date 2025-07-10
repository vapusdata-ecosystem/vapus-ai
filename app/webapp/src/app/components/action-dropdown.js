"use client";

import { useState, useEffect, useRef } from "react";
import ActionDropdownMenu from "./resource-action-handler";

export default function ActionDropdown({
  response,
  globalContext,
  customButton,
  apiServices = {},
}) {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);

  const toggleDropdown = () => {
    setIsOpen(!isOpen);
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        onClick={toggleDropdown}
        className="flex items-center px-1 py-1 bg-orange-700 text-white rounded-md focus:outline-none hover:bg-pink-900 cursor-pointer"
      >
        Actions
        <svg
          className={`w-4 h-4 ml-2 transform ${isOpen ? "rotate-180" : ""}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </button>

      <ActionDropdownMenu
        response={response || {}}
        globalContext={globalContext || {}}
        isVisible={isOpen}
        customButton={customButton}
        apiServices={apiServices}
      />
    </div>
  );
}

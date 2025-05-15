"use client";
import { useState, useEffect, useRef } from "react";
import { strTitle } from "./JS/common";

const MultiSelectDropdown = ({
  options,
  placeholder = "Select Options",
  onSelectionChange,
  initialSelected = [],
  isLoading = false,
  dropdownId = "dropdown", // default fallback
}) => {
  const [showDropdownMenu, setShowDropdownMenu] = useState(false);
  const [selectedItems, setSelectedItems] = useState(initialSelected);
  const initialRender = useRef(true);

  useEffect(() => {
    if (initialRender.current) {
      initialRender.current = false;
      return;
    }
    if (onSelectionChange) {
      onSelectionChange(selectedItems);
    }
  }, [selectedItems]);

  useEffect(() => {
    const handleClickOutside = (event) => {
      const dropdown = document.getElementById(`${dropdownId}-menu`);
      const button = document.getElementById(`${dropdownId}-button`);
      if (
        dropdown &&
        button &&
        !dropdown.contains(event.target) &&
        !button.contains(event.target)
      ) {
        setShowDropdownMenu(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [dropdownId]);

  const handleItemCheckboxChange = (itemId) => {
    setSelectedItems((prev) =>
      prev.includes(itemId)
        ? prev.filter((id) => id !== itemId)
        : [...prev, itemId]
    );
  };

  const selectedNames = options
    .filter((item) => selectedItems.includes(item.id))
    .map((item) => strTitle(item.name));

  const buttonText =
    selectedNames.length > 0 ? selectedNames.join(", ") : placeholder;

  return (
    <div className="relative inline-block text-left w-full">
      <div>
        <button
          type="button"
          id={`${dropdownId}-button`}
          className=" w-full flex text-gray-300 justify-between items-center bg-zinc-800 px-4 py-2 text-sm font-medium focus:ring-1 focus:ring-orange-700"
          onClick={() => setShowDropdownMenu(!showDropdownMenu)}
        >
          {buttonText}
          <svg
            className="w-5 h-5 text-gray-500"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fillRule="evenodd"
              d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
              clipRule="evenodd"
            />
          </svg>
        </button>
      </div>

      <div
        id={`${dropdownId}-menu`}
        className={`absolute z-10 mt-2 w-full bg-zinc-800 rounded-md shadow-lg ${
          showDropdownMenu ? "" : "hidden"
        }`}
      >
        <div className="p-2 space-y-2 max-h-60 overflow-y-auto scrollbar">
          {isLoading ? (
            <div className="text-sm text-gray-400 py-2 text-center">
              Loading...
            </div>
          ) : options.length === 0 ? (
            <div className="text-sm text-gray-400 py-2 text-center">
              No options available
            </div>
          ) : (
            options.map((item) => (
              <div key={item.id} className="flex items-center">
                <input
                  type="checkbox"
                  id={`item-${dropdownId}-${item.id}`}
                  value={item.id}
                  checked={selectedItems.includes(item.id)}
                  onChange={() => handleItemCheckboxChange(item.id)}
                  className="h-4 w-4 cursor-pointer mr-2 accent-orange-700 rounded"
                />
                <label
                  htmlFor={`item-${dropdownId}-${item.id}`}
                  className="ml-2 text-sm cursor-pointer"
                >
                  {strTitle(item.name)}
                </label>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
};

export default MultiSelectDropdown;

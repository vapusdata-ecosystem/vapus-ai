import React from "react";

export default function AddButton({ name, onClick }) {
  return (
    <button
      type="button"
      onClick={onClick}
      className="mt-2 px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 cursor-pointer"
      suppressHydrationWarning
    >
      {name}
    </button>
  );
}

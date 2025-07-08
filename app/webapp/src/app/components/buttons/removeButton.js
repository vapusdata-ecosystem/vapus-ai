import React from "react";

export default function RemoveButton({ onClick, name = "Remove" }) {
  return (
    <button
      type="button"
      onClick={onClick}
      className="mt-2 px-2 py-1 bg-red-500 text-white rounded text-xs hover:bg-red-600 cursor-pointer"
    >
      {name}
    </button>
  );
}

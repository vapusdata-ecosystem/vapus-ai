// PromptCard.jsx
"use client";

import { useState } from "react";
import ViewDetailsSvg from "./view-details-icon";

export default function Card({ prompt, backListingLink }) {
  const [upvotes, setUpvotes] = useState(0);
  const [downvotes, setDownvotes] = useState(0);

  const handleUpvote = () => {
    setUpvotes(upvotes + 1);
  };

  const handleDownvote = () => {
    setDownvotes(downvotes + 1);
  };

  return (
    <div className="relative p-4 bg-[#1b1b1b] rounded-lg text-xs shadow-lg border border-zinc-600">
      <a
        href={`/ai-center/prompts/${prompt.PromptId}`}
        className="absolute top-2 right-2 text-gray-100 hover:text-blue-600"
        target="_blank"
        title="Open in new tab"
        rel="noreferrer"
      >
        <ViewDetailsSvg />
      </a>

      <h3 className="text-md font-semibold text-gray-100 mb-3">
        {prompt.Name}
      </h3>

      <div className="grid grid-cols-1 gap-2 text-gray-100">
        <div className="flex items-center justify-left">
          <span className="font-semibold pr-2">Owner:</span>
          <span>{prompt.PromptOwner}</span>
        </div>

        <div className="flex items-start justify-left">
          <span className="font-semibold pr-2">Prompt Types:</span>
          <div className="flex flex-wrap gap-1">
            {prompt.PromptTypes &&
              prompt.PromptTypes.map((tag, index) => (
                <span
                  key={index}
                  className="px-2 py-1 text-xs font-medium rounded-full text-purple-800 bg-purple-100"
                >
                  {tag}
                </span>
              ))}
          </div>
        </div>

        <div className="flex items-start justify-left">
          <span className="font-semibold pr-2">Labels:</span>
          <div className="flex flex-wrap gap-1">
            {prompt.Labels &&
              prompt.Labels.map((tag, index) => (
                <span
                  key={index}
                  className="px-2 py-1 text-xs font-medium rounded-full text-yellow-800 bg-yellow-100"
                >
                  {tag}
                </span>
              ))}
          </div>
        </div>

        <div className="flex items-center justify-left">
          <span className="font-semibold pr-2">Status:</span>
          <span
            className={`px-3 py-1 font-medium rounded-full ${
              prompt.ResourceBase?.Status === "ACTIVE"
                ? "text-green-800 bg-green-100"
                : "text-red-800 bg-red-100"
            }`}
          >
            {prompt.ResourceBase?.Status}
          </span>
        </div>
      </div>

      <div className="absolute bottom-2 right-2 flex items-center space-x-4">
        <button
          className="flex items-center text-gray-100 hover:text-blue-600"
          onClick={handleUpvote}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="currentColor"
            className="w-6 h-6"
            viewBox="0 0 24 24"
          >
            <path d="M12 2L4 14h16L12 2zm0 18v-6h-4v6h4zm-5 0h4v2H7v-2zm10 0v2h-4v-2h4z"></path>
          </svg>
          <span className="ml-1">{upvotes}</span>
        </button>

        <button
          className="flex items-center text-gray-100 hover:text-blue-600"
          onClick={handleDownvote}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="currentColor"
            className="w-6 h-6"
            viewBox="0 0 24 24"
          >
            <path d="M12 22L4 10h16L12 22zm0-18v6h-4V4h4zm-5 0H7v2h4V4zm10 0v2h-4V4h4z"></path>
          </svg>
          <span className="ml-1">{downvotes}</span>
        </button>
      </div>
    </div>
  );
}

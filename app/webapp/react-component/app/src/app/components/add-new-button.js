"use client";
import Link from "next/link";

export default function CreateNewButton({
  label = "Add New",
  href = "#",
  className = "text-white text-sm px-1 py-1 rounded-lg bg-orange-700 hover:bg-pink-900 text-lg flex items-center",
}) {
  return (
    <Link href={href} className={className}>
      <svg
        className="w-4 h-4 mr-1"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          fill="#fff"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M12 4v16m8-8H4"
        ></path>
      </svg>
      {label}
    </Link>
  );
}

export default function LoadingOverlay({ isLoading }) {
  return (
    <div
      className={`absolute inset-0 z-10 flex items-center h-full w-full justify-center bg-zinc-600/90 ${
        isLoading ? "" : "hidden"
      }`}
    >
      <svg
        className="animate-spin h-10 w-10 text-black m-4"
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
        />
        <path
          className="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
        />
      </svg>
      <span>Wait a moment...</span>
    </div>
  );
}

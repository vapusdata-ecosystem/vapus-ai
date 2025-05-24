export default function ViewDetailsSvg() {
  return (
    <div className="group inline-block">
      <svg
        viewBox="0 0 200 200"
        xmlns="http://www.w3.org/2000/svg"
        className="h-8 w-8"
      >
        {/* Circle */}
        <circle
          cx="100"
          cy="100"
          r="90"
          className="stroke-orange-700 group-hover:stroke-orange-500 transition duration-200"
          strokeWidth="10"
          fill="none"
        />

        {/* Arrow */}
        <g transform="rotate(315, 100, 100)">
          <line
            x1="60"
            y1="100"
            x2="140"
            y2="100"
            className="stroke-orange-700 group-hover:stroke-orange-500 transition duration-200"
            strokeWidth="10"
            strokeLinecap="round"
          />
          <path
            d="M120,80 L140,100 L120,120"
            className="stroke-orange-700 group-hover:stroke-orange-500 transition duration-200"
            strokeWidth="10"
            strokeLinecap="round"
            strokeLinejoin="round"
            fill="none"
          />
        </g>
      </svg>
      {/* <!-- Custom Tooltip --> */}
      <div
        className="absolute bottom-full right-1/2 transform-translate-y-1/2  
              hidden group-hover:block bg-gray-700 text-gray-100 text-xs rounded px-2 py-1 z-50 whitespace-nowrap"
      >
        View prompts detail
      </div>
    </div>
  );
}

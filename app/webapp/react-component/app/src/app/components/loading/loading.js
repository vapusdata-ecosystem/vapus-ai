import React from 'react';

export default function LoadingOverlay({ 
  isLoading,
  size = 'default', 
  text = 'Loading', 
  showProgressBar = true,
  className = '',
  isOverlay = false 
}) {
  // Size variants
  const sizeClasses = {
    small: {
      spinner: 'w-8 h-8',
      dot: 'w-1 h-1',
      progressBar: 'w-32 h-0.5'
    },
    default: {
      spinner: 'w-16 h-16',
      dot: 'w-2 h-2',
      progressBar: 'w-48 h-1'
    },
    large: {
      spinner: 'w-24 h-24',
      dot: 'w-3 h-3',
      progressBar: 'w-64 h-1.5'
    }
  };

  const currentSize = sizeClasses[size] || sizeClasses.default;

  // Don't render if not loading
  if (!isLoading) {
    return null;
  }

  // Base container classes
  const containerClasses = isOverlay 
    ? `absolute inset-0 z-10 flex items-center h-full w-full justify-center bg-zinc-600/90 ${className}`
    : `flex items-center justify-center min-h-screen bg-[#1b1b1b] ${className}`;

  return (
    <div className={containerClasses}>
      <div className="flex flex-col items-center space-y-6">
        {/* Animated Spinner */}
        <div className="relative">
          <div className={`${currentSize.spinner} border-4 border-slate-200 dark:border-slate-700 rounded-full animate-pulse`}></div>
          <div className={`absolute top-0 left-0 ${currentSize.spinner} border-4 border-transparent border-t-orange-700 border-r-orange-700 rounded-full animate-spin`}></div>
          <div className={`absolute top-1/2 left-1/2 ${currentSize.dot} bg-orange-700 rounded-full transform -translate-x-1/2 -translate-y-1/2 animate-ping`}></div>
        </div>

        {/* Loading text with animated dots */}
        <div className="flex items-center space-x-1">
          <span className="text-slate-600 dark:text-slate-400 font-medium">
            {text}
          </span>
          <div className="flex space-x-1">
            <div
              className="w-1 h-1 bg-orange-700 rounded-full animate-bounce"
              style={{ animationDelay: "0ms" }}
            ></div>
            <div
              className="w-1 h-1 bg-orange-700 rounded-full animate-bounce"
              style={{ animationDelay: "150ms" }}
            ></div>
            <div
              className="w-1 h-1 bg-orange-700 rounded-full animate-bounce"
              style={{ animationDelay: "300ms" }}
            ></div>
          </div>
        </div>

        {/* Progress bar */}
        {showProgressBar && (
          <div className={`${currentSize.progressBar} bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden`}>
            <div
              className="h-full bg-gradient-to-r from-orange-700 to-orange-500 rounded-full animate-pulse"
              style={{
                width: "60%",
                animation: "loading-progress 2s ease-in-out infinite",
              }}
            ></div>
          </div>
        )}
      </div>

      <style jsx>{`
        @keyframes loading-progress {
          0% {
            width: 0%;
          }
          50% {
            width: 70%;
          }
          100% {
            width: 100%;
          }
        }
      `}</style>
    </div>
  );
};
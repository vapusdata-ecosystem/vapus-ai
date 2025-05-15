import React from "react";

const AuthExpiredModal = ({ isOpen, onLogin, onStayHere }) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-zinc-600/80 flex items-center justify-center z-50">
      <div className="bg-zinc-800 rounded-lg shadow-lg w-3/4 md:w-1/2 lg:w-1/3 text-gray-100">
        <div className="bg-zinc-900 px-4 py-2 flex justify-between items-center rounded-t-lg">
          <h3 className="text-lg font-medium">Session Expired</h3>
          <button
            onClick={onStayHere}
            className="text-gray-400 hover:text-pink-900 cursor-pointer"
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path
                fillRule="evenodd"
                d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                clipRule="evenodd"
              ></path>
            </svg>
          </button>
        </div>

        <div className="p-6">
          <p className="mb-6">
            Your access token has expired. Would you like to log in again?
          </p>

          <div className="flex justify-end space-x-3">
            <button
              onClick={onStayHere}
              className="px-4 py-2 bg-zinc-700 hover:bg-zinc-600 rounded-md cursor-pointer"
            >
              Stay Here
            </button>
            <button
              onClick={onLogin}
              className="px-4 py-2 bg-orange-700 hover:bg-pink-900 rounded-md cursor-pointer"
            >
              Go to Login
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AuthExpiredModal;

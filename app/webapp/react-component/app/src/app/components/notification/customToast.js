"use client";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { createGlobalStyle } from "styled-components";

export default function ToastContainerMessage({
  position = "top-right",
  autoClose = 3000,
  hideProgressBar = false,
  newestOnTop = true,
  closeOnClick = true,
  rtl = false,
  pauseOnFocusLoss = true,
  draggable = true,
  pauseOnHover = true,
  theme = "dark",
}) {
  // Dynamic styles based on the theme
  const ToastStyles = createGlobalStyle`
    /* Toast container base styles */
    .Toastify__toast {
      width: 400px;
      font-size: 14px;
      font-weight: 500;
    }
    
    /* Error toast */
    .Toastify__toast--error {
      background-color: #dc7575 !important;
      color: #660000 !important;
    }
    
    .Toastify__progress-bar--error {
      background: #660000 !important;
    }
    
    /* Success toast */
    .Toastify__toast--success {
      background-color: #75dc9c !important;
      color: #003312 !important;
    }
    
    .Toastify__progress-bar--success {
      background: #003312 !important;
    }
    
    /* Info toast */
    .Toastify__toast--info {
      background-color: #75aadc !important;
      color: #002b66 !important;
    }
    
    .Toastify__progress-bar--info {
      background: #002b66 !important;
    }
  `;

  return (
    <>
      <ToastStyles />
      <ToastContainer
        position={position}
        autoClose={autoClose}
        hideProgressBar={hideProgressBar}
        newestOnTop={newestOnTop}
        closeOnClick={closeOnClick}
        rtl={rtl}
        pauseOnFocusLoss={pauseOnFocusLoss}
        draggable={draggable}
        pauseOnHover={pauseOnHover}
        theme={theme}
      />
    </>
  );
}

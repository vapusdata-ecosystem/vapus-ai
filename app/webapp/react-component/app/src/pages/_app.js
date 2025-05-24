"use client";
import React from "react";
import "../app/globals.css";
import { Auth0Provider } from "@auth0/nextjs-auth0";

export default function App({ Component, pageProps }) {
  return (
    <Auth0Provider>
      <Component {...pageProps} />
    </Auth0Provider>
  );
}

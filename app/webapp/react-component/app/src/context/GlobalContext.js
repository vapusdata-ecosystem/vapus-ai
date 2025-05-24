import { accessToken } from "@/app/utils/api";
import { domainApi } from "@/app/utils/settings-endpoint/domain-api";

// function to decode accessToken and get the data
function parseJwt(token) {
  try {
    const base64Url = token.split(".")[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map(function (c) {
          return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
        })
        .join("")
    );

    return JSON.parse(jsonPayload);
  } catch (error) {
    return null;
  }
}

// data from accessToken
export const userGlobalData = async () => {
  const token = accessToken;
  const userInfo = parseJwt(token)?.scope || null;
  return {
    userInfo,
  };
};

// user data
export const getGlobalData = async () => {
  try {
    const domains = await domainApi.getDomains();
    const currentDomain = domains?.output?.domains?.[0] || null;

    return {
      currentDomain,
    };
  } catch (error) {
    console.error("Error fetching global data:", error);
    return {
      domains: [],
      currentDomain: null,
      error,
    };
  }
};

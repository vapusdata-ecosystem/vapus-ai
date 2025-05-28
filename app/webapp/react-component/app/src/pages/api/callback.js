import { LoginCallbacksApi } from "@/app/utils/auth-endpoint/auth";
const AUTH_CONFIG = {
  loginPath: "/login",
  loginRedirectPath: "/login",
  callbackPath: "/api/callback",
  logoutPath: "/logout",
  homePath: "/dashboard",
  accessTokenCookieName: "access_token",
  idTokenCookieName: "id_token",
  cookiePath: "/",
};

export default async function handler(req, res) {
  console.log("Callback received with query params:", req.query);
  console.log("Code:", req.query.code);
  if (req.method !== "GET") {
    return res.status(405).json({ message: "Method Not Allowed" });
  }

  const { code } = req.query;

  if (!code) {
    return res.redirect(AUTH_CONFIG.loginPath);
  }

  try {
    const payload = {
      code,
      host: `http://127.0.0.1:9014/${AUTH_CONFIG.callbackPath}`,
    };

    const result = await LoginCallbacksApi.getLoginCallback(payload);

    if (!result || !result.token) {
      return res.redirect(`${AUTH_CONFIG.loginPath}?error=callback_failed`);
    }

    // Set cookies
    const cookieOptions = {
      path: AUTH_CONFIG.cookiePath,
      expires: new Date(result.token.validTill * 1000),
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "lax",
    };

    // Set ID token and access token cookies
    res.setHeader("Set-Cookie", [
      `${AUTH_CONFIG.idTokenCookieName}=${result.token.idToken}; Path=${
        cookieOptions.path
      }; Expires=${cookieOptions.expires.toUTCString()}; HttpOnly=${
        cookieOptions.httpOnly
      }; SameSite=${cookieOptions.sameSite}${
        cookieOptions.secure ? "; Secure" : ""
      }`,
      `${AUTH_CONFIG.accessTokenCookieName}=${result.token.accessToken}; Path=${
        cookieOptions.path
      }; Expires=${cookieOptions.expires.toUTCString()}; HttpOnly=${
        cookieOptions.httpOnly
      }; SameSite=${cookieOptions.sameSite}${
        cookieOptions.secure ? "; Secure" : ""
      }`,
    ]);

    // Check if registration is needed
    if (!result.token.accessToken) {
      return res.redirect(`${AUTH_CONFIG.loginPath}?register=true`);
    }

    // Redirect to home or stored redirect URL
    const redirectTo = req.cookies.loginRedirectUrl || AUTH_CONFIG.homePath;

    // Clear the redirect URL cookie if it exists
    if (req.cookies.loginRedirectUrl) {
      res.setHeader("Set-Cookie", [
        `loginRedirectUrl=; Path=${AUTH_CONFIG.cookiePath}; Expires=Thu, 01 Jan 1970 00:00:00 GMT`,
      ]);
    }

    return res.redirect(redirectTo);
  } catch (error) {
    console.error("Login callback error:", error);
    return res.redirect(`${AUTH_CONFIG.loginPath}?error=callback_error`);
  }
}

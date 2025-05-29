import { NextResponse } from "next/server";
import { LoginCallbacksApi } from "@/app/utils/auth-endpoint/auth";
import { cookies as getServerCookies } from "next/headers";

const AUTH_CONFIG = {
  loginPath: "/login",
  loginRedirectPath: "/login",
  callbackPath: "/auth/callback",
  logoutPath: "/logout",
  homePath: "/dashboard",
  accessTokenCookieName: "access_token",
  idTokenCookieName: "id_token",
  cookiePath: "/",
};

export async function GET(request) {
  const { searchParams } = new URL(request.url);
  const code = searchParams.get("code");

  const serverCookies = getServerCookies();

  if (!code) {
    return NextResponse.redirect(new URL(AUTH_CONFIG.loginPath, request.url));
  }

  try {
    const payload = {
      code,
      host: `http://127.0.0.1:9014${AUTH_CONFIG.callbackPath}`,
    };

    const result = await LoginCallbacksApi.getLoginCallback(payload);

    if (!result?.token?.accessToken) {
      return NextResponse.redirect(
        new URL(`${AUTH_CONFIG.loginPath}?register=true`, request.url)
      );
    }

    const redirectTo =
      serverCookies.get("loginRedirectUrl")?.value || AUTH_CONFIG.homePath;

    // Calculate expiry time for cookies in JS (e.g., in seconds)
    const expirySeconds =
      result.token.validTill - Math.floor(Date.now() / 1000);

    // Create an HTML page that sets the cookies  and then redirects
    const html = `
      <!DOCTYPE html>
      <html>
      <head><title>Logging in...</title></head>
      <body>
        <script>
          // Set the cookies with expiry, path and SameSite=Lax
          const expires = new Date(Date.now() + ${expirySeconds} * 1000).toUTCString();
          document.cookie = "${AUTH_CONFIG.idTokenCookieName}=${result.token.idToken}; path=${AUTH_CONFIG.cookiePath}; expires=" + expires + "; SameSite=Lax";
          document.cookie = "${AUTH_CONFIG.accessTokenCookieName}=${result.token.accessToken}; path=${AUTH_CONFIG.cookiePath}; expires=" + expires + "; SameSite=Lax";

          // Clear loginRedirectUrl cookie by setting expiration in the past
          document.cookie = "loginRedirectUrl=; path=${AUTH_CONFIG.cookiePath}; expires=Thu, 01 Jan 1970 00:00:00 GMT; SameSite=Lax";

          // Redirect to the original page
          window.location.href = "${redirectTo}";
        </script>
        <p>Logging in, please wait...</p>
      </body>
      </html>
    `;

    return new NextResponse(html, {
      status: 200,
      headers: {
        "Content-Type": "text/html",
      },
    });
  } catch (error) {
    console.error("Login callback error:", error);
    return NextResponse.redirect(
      new URL(`${AUTH_CONFIG.loginPath}?error=callback_error`, request.url)
    );
  }
}

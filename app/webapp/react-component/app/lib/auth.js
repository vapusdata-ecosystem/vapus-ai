"use client";
import { jwtDecode } from "jwt-decode";
import { loginApi } from "@/app/utils/auth-endpoint/auth";
import { useState, useEffect } from "react";
import { createRoot } from "react-dom/client";
import AuthExpiredModal from "@/app/components/notification/authExpiredPopPup";

const AUTH_CONFIG = {
  loginPath: "/login",
  loginRedirectPath: "/login",
  callbackPath: "/api/callback",
  logoutPath: "/login",
  homePath: "/settings/domain",
  accessTokenCookieName: "access_token",
  idTokenCookieName: "id_token",
  cookiePath: "/",
  tokenExpiryCheckInterval: 60000,
};

export class AuthService {
  constructor() {
    this.routes = {
      login: AUTH_CONFIG.loginPath,
      loginRedirect: AUTH_CONFIG.loginRedirectPath,
      loginCallback: AUTH_CONFIG.callbackPath,
      logout: AUTH_CONFIG.logoutPath,
      home: AUTH_CONFIG.homePath,
    };

    this.tokenCheckInterval = null;
    this.showingModal = false;

    // Only setup token check in browser environment and only once
    if (typeof window !== "undefined" && !this.tokenCheckInterval) {
      this.setupTokenExpiryCheck();
    }
  }

  // check for token expiration
  setupTokenExpiryCheck() {
    this.checkTokenExpiration();
    this.tokenCheckInterval = setInterval(() => {
      this.checkTokenExpiration();
    }, AUTH_CONFIG.tokenExpiryCheckInterval);
  }

  checkTokenExpiration() {
    // Skip check on login page
    if (
      typeof window !== "undefined" &&
      window.location.pathname === this.routes.login
    ) {
      return;
    }

    // If not authenticated or token expired, show modal instead of redirecting
    if (!this.isAuthenticated() && !this.showingModal) {
      this.showAuthExpiredModal();
    }
  }

  // Method to show the auth expired modal
  showAuthExpiredModal() {
    if (typeof window !== "undefined" && !this.showingModal) {
      this.showingModal = true;

      // Create a div for our modal
      const modalDiv = document.createElement("div");
      modalDiv.id = "auth-expired-modal-container";
      document.body.appendChild(modalDiv);
      const modalRoot = createRoot(modalDiv);

      // Render our modal with cleanup functions properly wrapped
      modalRoot.render(
        <AuthExpiredModal
          isOpen={true}
          onLogin={() => {
            this.cleanupModal(modalDiv, modalRoot);
            this.redirectToLogin();
          }}
          onStayHere={() => {
            this.cleanupModal(modalDiv, modalRoot);
          }}
        />
      );
    }
  }

  // Helper to clean up the modal
  cleanupModal(modalDiv, root) {
    if (root) {
      root.unmount();
    }
    if (modalDiv && modalDiv.parentNode) {
      modalDiv.parentNode.removeChild(modalDiv);
    }
    this.showingModal = false;
  }

  // Keep the existing redirectToLogin method for when the user chooses to log in
  redirectToLogin() {
    if (typeof window !== "undefined") {
      window.location.href = this.routes.login;
    }
  }

  async login(landingPage = this.routes.home) {
    if (typeof window !== "undefined") {
      try {
        // Store the landing page URL in a cookie for retrieval after authentication
        document.cookie = `loginRedirectUrl=${landingPage}; path=${AUTH_CONFIG.cookiePath}`;
        localStorage.setItem("loginRedirectUrl", landingPage);

        const data = await loginApi.getLogin();
        if (data && data.loginUrl) {
          window.location.href = data.loginUrl;
          return { success: true };
        } else {
          // Redirect with error parameter for invalid response
          const errorMsg = "Invalid response from server";
          window.location.href = `${
            this.routes.login
          }?error=${encodeURIComponent(errorMsg)}`;
          return { success: false, error: errorMsg };
        }
      } catch (error) {
        console.error("Error getting login URL:", error);

        // Redirect with error parameter when login fails
        const errorMsg = error.message || "Login failed. Please try again.";
        window.location.href = `${this.routes.login}?error=${encodeURIComponent(
          errorMsg
        )}`;

        return {
          success: false,
          error: errorMsg,
        };
      }
    }
    return { success: false, error: "Not in browser environment" };
  }

  logout() {
    this.clearAuthCookies();
    if (typeof window !== "undefined") {
      window.location.href = this.routes.logout;
    }
    return null;
  }

  //Clear authentication cookies
  clearAuthCookies() {
    this.setCookie(AUTH_CONFIG.accessTokenCookieName, "", -1);
    this.setCookie(AUTH_CONFIG.idTokenCookieName, "", -1);
    this.setCookie(`${AUTH_CONFIG.accessTokenCookieName}_expiry`, "", -1);
  }

  //Check if user is authenticated by looking for the access token and verifying it's not expired
  isAuthenticated() {
    if (typeof window === "undefined") return false;

    const accessToken = this.getCookie(AUTH_CONFIG.accessTokenCookieName);
    if (!accessToken) {
      return false;
    }
    return !this.isTokenExpired(accessToken);
  }

  /**
   * Check if a token is expired
   * @param {string} token
   * @returns {boolean}
   */
  isTokenExpired(token) {
    try {
      const decodedToken = jwtDecode(token);
      // Check if the token has an expiration time
      if (!decodedToken.exp) return true;

      // Compare expiration timestamp with current time
      const expiryDate = new Date(decodedToken.exp * 1000);
      return expiryDate <= new Date();
    } catch (error) {
      console.error("Error checking token expiry:", error);
      return true; // Consider token expired if we can't decode it
    }
  }

  /**
   * Get token expiration date
   * @returns {Date|null}
   */
  getTokenExpiry() {
    const accessToken = this.getCookie(AUTH_CONFIG.accessTokenCookieName);

    if (!accessToken) return null;

    try {
      const decodedToken = jwtDecode(accessToken);
      if (decodedToken.exp) {
        return new Date(decodedToken.exp * 1000);
      }
    } catch (error) {
      console.error("Error decoding token for expiry:", error);
    }

    // Fallback to cookie expiry if JWT decode fails
    if (typeof document !== "undefined") {
      const cookies = document.cookie.split(";");
      for (let i = 0; i < cookies.length; i++) {
        const cookie = cookies[i].trim();
        if (cookie.startsWith(`${AUTH_CONFIG.accessTokenCookieName}_expiry=`)) {
          const expiryTimestamp = parseInt(
            cookie.substring(
              `${AUTH_CONFIG.accessTokenCookieName}_expiry=`.length
            ),
            10
          );
          return new Date(expiryTimestamp);
        }
      }
    }

    return null;
  }

  /**
   * Get cookie by name
   * @param {string} name
   * @returns {string|null}
   */
  getCookie(name) {
    if (typeof document === "undefined") return null;

    const cookies = document.cookie.split(";");
    for (let i = 0; i < cookies.length; i++) {
      const cookie = cookies[i].trim();
      if (cookie.startsWith(name + "=")) {
        return cookie.substring(name.length + 1);
      }
    }
    return null;
  }

  /**
   * Set cookie with expiry
   * @param {string} name
   * @param {string} value
   * @param {number} daysToExpire
   */
  setCookie(name, value, daysToExpire) {
    if (typeof document === "undefined") return;

    let expires = "";
    if (daysToExpire) {
      const date = new Date();
      date.setTime(date.getTime() + daysToExpire * 24 * 60 * 60 * 1000);
      expires = `; expires=${date.toUTCString()}`;
    }

    document.cookie = `${name}=${value}${expires}; path=${AUTH_CONFIG.cookiePath}`;
  }

  /**
   * Set authentication tokens with expiry
   * @param {string} accessToken
   * @param {string} idToken
   * @param {number} validTill
   */
  setAuthTokens(accessToken, idToken, validTill) {
    if (!accessToken || !idToken || !validTill) {
      return;
    }

    const expiryDate = new Date(validTill * 1000);
    const now = new Date();
    const daysToExpiry = (expiryDate - now) / (1000 * 60 * 60 * 24);

    // Set the tokens
    this.setCookie(
      AUTH_CONFIG.accessTokenCookieName,
      accessToken,
      daysToExpiry
    );
    this.setCookie(AUTH_CONFIG.idTokenCookieName, idToken, daysToExpiry);

    // Set expiry info in a separate cookie
    this.setCookie(
      `${AUTH_CONFIG.accessTokenCookieName}_expiry`,
      expiryDate.getTime(),
      daysToExpiry
    );
  }

  /**
   * Get user info from the access token by decoding the JWT
   * @returns {Object|null}
   */
  getUserInfo() {
    const accessToken = this.getCookie(AUTH_CONFIG.accessTokenCookieName);
    const idToken = this.getCookie(AUTH_CONFIG.idTokenCookieName);

    if (!accessToken) {
      return null;
    }

    try {
      const tokenToDecode = idToken || accessToken;
      const decodedToken = jwtDecode(tokenToDecode);

      const scope = decodedToken.scope || {};

      const userInfo = {
        isAuthenticated: true,
        accessToken,
        idToken,
        sub: decodedToken.sub,
        userId: scope.userId,
        accountId: scope.accountId,
        domainId: scope.domainId,
        domainRole: scope.domainRole,
        roleScope: scope.roleScope,
        platformRole: scope.platformRole,
        decodedToken,
        exp: decodedToken.exp,
        expiresAt: new Date(decodedToken.exp * 1000),
      };

      return userInfo;
    } catch (error) {
      console.error("Error decoding JWT token:", error);
      return {
        isAuthenticated: true,
        accessToken,
        idToken,
        error: "Failed to decode token",
      };
    }
  }

  cleanup() {
    if (this.tokenCheckInterval) {
      clearInterval(this.tokenCheckInterval);
      this.tokenCheckInterval = null;
    }
  }
}

// Singleton instance
let authServiceInstance;

// Factory function to get auth service instance
export function getAuthService() {
  if (typeof window !== "undefined") {
    if (!authServiceInstance) {
      authServiceInstance = new AuthService();
    }
    return authServiceInstance;
  }

  // Return a minimal implementation for server-side rendering
  return {
    isAuthenticated: () => false,
    getUserInfo: () => null,
    login: () => ({ success: false, error: "Server-side call" }),
  };
}

/**
 * React component that requires authentication
 * @param {Component} Component
 */
export function withAuth(Component) {
  const WithAuthComponent = (props) => {
    const [loading, setLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
      const authService = getAuthService();
      if (!authService.isAuthenticated()) {
        router.replace(AUTH_CONFIG.loginPath);
      } else {
        setLoading(false);
      }
    }, [router]);

    // Show loading state or protected component
    if (loading) {
      return <div>Loading...</div>;
    }

    return <Component {...props} />;
  };

  // Add display name
  WithAuthComponent.displayName = `withAuth(${
    Component.displayName || Component.name || "Component"
  })`;

  return WithAuthComponent;
}

// Hook for accessing authentication in functional components
export function useAuth() {
  const [auth, setAuth] = useState({
    isAuthenticated: false,
    user: null,
    loading: true,
  });

  useEffect(() => {
    const authService = getAuthService();
    const isAuthenticated = authService.isAuthenticated();
    const user = authService.getUserInfo();

    setAuth({
      isAuthenticated,
      user,
      loading: false,
    });

    return () => {};
  }, []);

  return auth;
}

//  loginCallback handling for client side
export function handleCallbackResponse(tokenData) {
  if (
    !tokenData ||
    !tokenData.accessToken ||
    !tokenData.idToken ||
    !tokenData.validTill
  ) {
    return false;
  }

  // Set authentication tokens with expiry
  const authService = getAuthService();
  authService.setAuthTokens(
    tokenData.accessToken,
    tokenData.idToken,
    tokenData.validTill
  );
  return true;
}

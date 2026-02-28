import * as AuthSession from "expo-auth-session";
import * as SecureStore from "expo-secure-store";

const TOKEN_KEY = "auth0_access_token";
const REFRESH_TOKEN_KEY = "auth0_refresh_token";
const TOKEN_EXPIRY_KEY = "auth0_token_expiry";

// Promise to track ongoing refresh operation
let refreshPromise: Promise<string | null> | null = null;

const auth0Domain = process.env.EXPO_PUBLIC_AUTH0_DOMAIN;
const auth0ClientId = process.env.EXPO_PUBLIC_AUTH0_CLIENT_ID;
const auth0Audience = process.env.EXPO_PUBLIC_AUTH0_AUDIENCE;

if (!auth0Domain || !auth0ClientId || !auth0Audience) {
  console.warn("Auth0 env vars are missing");
}

const discovery = {
  authorizationEndpoint: `https://${auth0Domain}/authorize`,
  tokenEndpoint: `https://${auth0Domain}/oauth/token`,
  revocationEndpoint: `https://${auth0Domain}/oauth/revoke`,
  userInfoEndpoint: `https://${auth0Domain}/userinfo`,
};

export async function loginWithAuth0(): Promise<string | null> {
  // Use custom scheme for stable redirect URI
  const redirectUri = AuthSession.makeRedirectUri({
    scheme: "nqfrontend",
    path: "auth",
  });

  const request = new AuthSession.AuthRequest({
    clientId: auth0ClientId,
    redirectUri,
    responseType: AuthSession.ResponseType.Code,
    scopes: ["openid", "profile", "email", "offline_access"], // Added offline_access for refresh token
    extraParams: {
      audience: auth0Audience
    },
    usePKCE: true,
  });

  const result = await request.promptAsync(discovery);
  if (result.type !== "success" || !result.params.code) {
    return null;
  }

  const tokenResponse = await AuthSession.exchangeCodeAsync(
    {
      clientId: auth0ClientId,
      code: result.params.code,
      redirectUri,
      extraParams: {
        code_verifier: request.codeVerifier ?? "",
      },
    },
    discovery,
  );

  // Ensure we got an access token before proceeding
  if (!tokenResponse.accessToken) {
    return null;
  }

  // Store access token, refresh token, and expiry time
  await SecureStore.setItemAsync(TOKEN_KEY, tokenResponse.accessToken);
  
  if (tokenResponse.refreshToken) {
    await SecureStore.setItemAsync(REFRESH_TOKEN_KEY, tokenResponse.refreshToken);
  }
  
  // Calculate and store expiry time
  if (tokenResponse.expiresIn) {
    const expiryTime = Date.now() + tokenResponse.expiresIn * 1000;
    await SecureStore.setItemAsync(TOKEN_EXPIRY_KEY, expiryTime.toString());
  }
  
  return tokenResponse.accessToken;
}

export async function getAccessToken(): Promise<string | null> {
  const token = await SecureStore.getItemAsync(TOKEN_KEY);
  if (!token) return null;
  
  // Check if token is expired or will expire in the next 5 minutes
  const expiryStr = await SecureStore.getItemAsync(TOKEN_EXPIRY_KEY);
  if (expiryStr) {
    const expiry = parseInt(expiryStr, 10);
    const fiveMinutes = 5 * 60 * 1000;
    
    if (Date.now() + fiveMinutes >= expiry) {
      // Token is expired or about to expire, try to refresh
      // If a refresh is already in progress, wait for it
      if (refreshPromise) {
        return refreshPromise;
      }
      
      // Start a new refresh operation
      refreshPromise = refreshAccessToken().finally(() => {
        // Clear the promise when done
        refreshPromise = null;
      });
      
      const refreshed = await refreshPromise;
      if (refreshed) {
        return refreshed;
      }
      // Refresh failed, clear tokens and return null
      await logout();
      return null;
    }
  }
  
  return token;
}

async function refreshAccessToken(): Promise<string | null> {
  const refreshToken = await SecureStore.getItemAsync(REFRESH_TOKEN_KEY);
  if (!refreshToken) return null;
  
  try {
    const tokenResponse = await AuthSession.refreshAsync(
      {
        clientId: auth0ClientId ?? "",
        refreshToken,
      },
      discovery,
    );
    
    if (!tokenResponse.accessToken) {
      return null;
    }
    
    // Store new tokens
    await SecureStore.setItemAsync(TOKEN_KEY, tokenResponse.accessToken);
    
    if (tokenResponse.refreshToken) {
      await SecureStore.setItemAsync(REFRESH_TOKEN_KEY, tokenResponse.refreshToken);
    } else {
      console.warn(
        "Auth0 refresh response did not include a new refresh token. " +
        "Ensure your Auth0 refresh token rotation settings match this assumption.",
      );
    }
    
    if (tokenResponse.expiresIn) {
      const expiryTime = Date.now() + tokenResponse.expiresIn * 1000;
      await SecureStore.setItemAsync(TOKEN_EXPIRY_KEY, expiryTime.toString());
    }
    
    return tokenResponse.accessToken;
  } catch (error) {
    console.error("Failed to refresh token:", error);
    return null;
  }
}

export async function logout(): Promise<void> {
  await SecureStore.deleteItemAsync(TOKEN_KEY);
  await SecureStore.deleteItemAsync(REFRESH_TOKEN_KEY);
  await SecureStore.deleteItemAsync(TOKEN_EXPIRY_KEY);
}

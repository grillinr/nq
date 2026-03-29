import * as AuthSession from 'expo-auth-session';
import { logError, logInfo } from './logger';
import {
  getCurrentAccountId,
  getAccountTokens,
  saveAccountTokens,
  removeAccount,
  type AccountTokens,
} from './accountStorage';

// Promise to track ongoing refresh operation
let refreshPromise: Promise<string | null> | null = null;

const auth0Domain = process.env.EXPO_PUBLIC_AUTH0_DOMAIN;
const auth0ClientId = process.env.EXPO_PUBLIC_AUTH0_CLIENT_ID;
const auth0Audience = process.env.EXPO_PUBLIC_AUTH0_AUDIENCE;

if (!auth0Domain || !auth0ClientId || !auth0Audience) {
  logInfo('Auth0 env vars are missing');
}

const discovery = {
  authorizationEndpoint: `https://${auth0Domain}/authorize`,
  tokenEndpoint: `https://${auth0Domain}/oauth/token`,
  revocationEndpoint: `https://${auth0Domain}/oauth/revoke`,
  userInfoEndpoint: `https://${auth0Domain}/userinfo`,
};

export interface LoginResult {
  accessToken: string;
  refreshToken?: string;
  expiryTime?: number;
}

export async function loginWithAuth0(forceLogin: boolean = false): Promise<LoginResult | null> {
  logInfo('[Auth0] Starting login flow', forceLogin ? '(forced new login)' : '');

  // makeRedirectUri auto-detects the correct URI for the environment:
  // - In Expo Go: exp://<ip>:<port>/--/
  // - In a standalone/dev build: nqfrontend://
  // Both must be registered in the Auth0 dashboard Allowed Callback URLs.
  const redirectUri = AuthSession.makeRedirectUri({
    native: 'nqfrontend://',
  });

  // Log the redirect URI so you can register it in Auth0 if login fails (dev only)
  if (__DEV__) {
    logInfo('[Auth0] redirect_uri:', redirectUri);
  }

  const extraParams: Record<string, string> = {
    audience: auth0Audience ?? '',
  };

  // Force fresh login if requested (for adding new accounts)
  if (forceLogin) {
    extraParams.prompt = 'login';
    logInfo('[Auth0] Forcing fresh login prompt');
  }

  const request = new AuthSession.AuthRequest({
    clientId: auth0ClientId ?? '',
    redirectUri,
    responseType: AuthSession.ResponseType.Code,
    scopes: ['openid', 'profile', 'email', 'offline_access'], // Added offline_access for refresh token
    extraParams,
    usePKCE: true,
  });

  logInfo('[Auth0] Prompting for authentication');
  const result = await request.promptAsync(discovery);

  if (result.type !== 'success' || !result.params.code) {
    logError('[Auth0] Login failed or cancelled:', result);
    return null;
  }

  logInfo('[Auth0] Got authorization code, exchanging for tokens');
  const tokenResponse = await AuthSession.exchangeCodeAsync(
    {
      clientId: auth0ClientId ?? '',
      code: result.params.code,
      redirectUri,
      extraParams: {
        code_verifier: request.codeVerifier ?? '',
      },
    },
    discovery
  );

  // Ensure we got an access token before proceeding
  if (!tokenResponse.accessToken) {
    logError('[Auth0] Token exchange failed - no access token received');
    return null;
  }

  logInfo('[Auth0] Successfully received tokens');

  // Calculate expiry time
  const expiryTime = tokenResponse.expiresIn
    ? Date.now() + tokenResponse.expiresIn * 1000
    : undefined;

  return {
    accessToken: tokenResponse.accessToken,
    refreshToken: tokenResponse.refreshToken ?? undefined,
    expiryTime,
  };
}

export async function getAccessToken(): Promise<string | null> {
  const accountId = await getCurrentAccountId();
  if (!accountId) return null;

  const tokens = await getAccountTokens(accountId);
  if (!tokens) return null;

  // Check if token is expired or will expire in the next 5 minutes
  if (tokens.expiryTime) {
    const fiveMinutes = 5 * 60 * 1000;

    if (Date.now() + fiveMinutes >= tokens.expiryTime) {
      // Token is expired or about to expire, try to refresh
      // If a refresh is already in progress, wait for it
      if (refreshPromise) {
        return refreshPromise;
      }

      // Start a new refresh operation
      refreshPromise = refreshAccessToken(accountId).finally(() => {
        // Clear the promise when done
        refreshPromise = null;
      });

      const refreshed = await refreshPromise;
      if (refreshed) {
        return refreshed;
      }
      // Refresh failed, clear tokens and return null
      await removeAccount(accountId);
      return null;
    }
  }

  return tokens.accessToken;
}

async function refreshAccessToken(accountId: string): Promise<string | null> {
  const tokens = await getAccountTokens(accountId);
  if (!tokens?.refreshToken) return null;

  try {
    const tokenResponse = await AuthSession.refreshAsync(
      {
        clientId: auth0ClientId ?? '',
        refreshToken: tokens.refreshToken,
      },
      discovery
    );

    if (!tokenResponse.accessToken) {
      return null;
    }

    // Calculate new expiry time
    const expiryTime = tokenResponse.expiresIn
      ? Date.now() + tokenResponse.expiresIn * 1000
      : undefined;

    // Store new tokens
    const newTokens: AccountTokens = {
      accessToken: tokenResponse.accessToken,
      refreshToken: tokenResponse.refreshToken ?? tokens.refreshToken,
      expiryTime,
    };

    if (!tokenResponse.refreshToken) {
      logInfo(
        'Auth0 refresh response did not include a new refresh token. ' +
          'Ensure your Auth0 refresh token rotation settings match this assumption.'
      );
    }

    await saveAccountTokens(accountId, newTokens);

    return tokenResponse.accessToken;
  } catch (error) {
    logError('Failed to refresh token:', error);
    return null;
  }
}

/**
 * Logout and remove account
 * @param accountId - Optional account ID to remove. If not provided, removes current account.
 */
export async function logout(accountId?: string): Promise<void> {
  const idToRemove = accountId ?? (await getCurrentAccountId());
  if (!idToRemove) return;

  await removeAccount(idToRemove);
}

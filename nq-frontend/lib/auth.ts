import * as AuthSession from "expo-auth-session";
import * as SecureStore from "expo-secure-store";

const TOKEN_KEY = "auth0_access_token";

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
  const redirectUri = AuthSession.makeRedirectUri({
    path: "auth",
  });

  const request = new AuthSession.AuthRequest({
    clientId: auth0ClientId ?? "",
    redirectUri,
    responseType: AuthSession.ResponseType.Code,
    scopes: ["openid", "profile", "email"],
    extraParams: {
      audience: auth0Audience ?? "",
    },
    usePKCE: true,
  });

  const result = await request.promptAsync(discovery);
  if (result.type !== "success" || !result.params.code) {
    return null;
  }

  const tokenResponse = await AuthSession.exchangeCodeAsync(
    {
      clientId: auth0ClientId ?? "",
      code: result.params.code,
      redirectUri,
      extraParams: {
        code_verifier: request.codeVerifier ?? "",
      },
    },
    discovery,
  );

  if (!tokenResponse.accessToken) {
    return null;
  }

  await SecureStore.setItemAsync(TOKEN_KEY, tokenResponse.accessToken);
  return tokenResponse.accessToken;
}

export async function getAccessToken(): Promise<string | null> {
  return SecureStore.getItemAsync(TOKEN_KEY);
}

export async function logout(): Promise<void> {
  await SecureStore.deleteItemAsync(TOKEN_KEY);
}

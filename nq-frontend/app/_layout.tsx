import { Stack } from "expo-router";

export default function RootLayout() {
  // Debug: log environment and expo router globals to help trace boolean/string mismatch
  try {
    // eslint-disable-next-line no-console
    console.log("RootLayout render", {
      NODE_ENV: process.env.NODE_ENV,
      EXPO_ROUTER: (global as any)?.EXPO_ROUTER,
      EXPO_ROUTER_OPTIONS: (global as any)?.EXPO_ROUTER_OPTIONS,
    });
  } catch (e) {
    // eslint-disable-next-line no-console
    console.log("RootLayout logging failed", e);
  }

  return <Stack />;
}

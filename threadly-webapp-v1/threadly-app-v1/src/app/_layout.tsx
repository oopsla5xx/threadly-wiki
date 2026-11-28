import "../../global.css"

import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { PortalHost } from "@rn-primitives/portal"
import { Stack } from "expo-router"
import { ThemeProvider } from "expo-router/react-navigation"
import * as SplashScreen from "expo-splash-screen"
import { StatusBar } from "expo-status-bar"
import { GestureHandlerRootView } from "react-native-gesture-handler"
import { NAV_THEME } from "@/lib/theme"

SplashScreen.setOptions({
  duration: 500,
  fade: true,
})

const queryClient = new QueryClient()

export default function RootLayout() {
  return (
    <GestureHandlerRootView style={{ flex: 1 }}>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider value={NAV_THEME.light}>
          <StatusBar style="auto" />
          <Stack />
          <PortalHost />
        </ThemeProvider>
      </QueryClientProvider>
    </GestureHandlerRootView>
  )
}

const fs = require("node:fs")
const path = require("node:path")

const root = path.resolve(__dirname, "..")
const src = path.join(root, "src")
const backup = path.join(root, "project-backup")
const app = path.join(src, "app")
const features = path.join(src, "features")

if (fs.existsSync(backup)) {
  console.error("project-backup already exists. Move it before resetting again.")
  process.exitCode = 1
} else if (!fs.existsSync(app) || !fs.existsSync(features)) {
  console.error("Expected src/app and src/features to exist; nothing was changed.")
  process.exitCode = 1
} else {
  fs.mkdirSync(backup)
  fs.renameSync(app, path.join(backup, "app"))
  fs.renameSync(features, path.join(backup, "features"))
  fs.mkdirSync(app)
  fs.writeFileSync(
    path.join(app, "_layout.tsx"),
    `import "../../global.css"

import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { PortalHost } from "@rn-primitives/portal"
import { Stack } from "expo-router"
import { ThemeProvider } from "expo-router/react-navigation"
import { StatusBar } from "expo-status-bar"
import { GestureHandlerRootView } from "react-native-gesture-handler"
import { NAV_THEME } from "@/lib/theme"

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
`,
  )
  fs.writeFileSync(
    path.join(app, "index.tsx"),
    `import { Text, View } from "react-native"

export default function Index() {
  return (
    <View>
      <Text>Start building your app.</Text>
    </View>
  )
}
`,
  )
  fs.mkdirSync(features)
  console.log("Previous routes and features moved to project-backup/.")
}

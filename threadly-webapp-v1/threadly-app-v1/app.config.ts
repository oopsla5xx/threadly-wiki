import type { ConfigContext, ExpoConfig } from "expo/config"

type AppVariant = "development" | "staging" | "production"

const variants: Record<
  AppVariant,
  { name: string; scheme: string; iosBundleIdentifier: string; androidPackage: string }
> = {
  development: {
    name: "Threadly Dev",
    scheme: "threadly-dev",
    iosBundleIdentifier: "com.quindarts.threadly-app-v1.dev",
    androidPackage: "com.quindarts.threadlyappv1.dev",
  },
  staging: {
    name: "Threadly Staging",
    scheme: "threadly-staging",
    iosBundleIdentifier: "com.quindarts.threadly-app-v1.staging",
    androidPackage: "com.quindarts.threadlyappv1.staging",
  },
  production: {
    name: "Threadly",
    scheme: "threadly",
    iosBundleIdentifier: "com.quindarts.threadly-app-v1",
    androidPackage: "com.quindarts.threadlyappv1",
  },
}

function getVariant(): AppVariant {
  const value = process.env.APP_VARIANT
  if (value === undefined) return "development"
  if (value === "development" || value === "staging" || value === "production") return value
  throw new Error(`Invalid APP_VARIANT "${value}": expected development, staging, or production.`)
}

export default ({ config }: ConfigContext): ExpoConfig => {
  const variant = getVariant()
  const current = variants[variant]

  return {
    ...config,
    name: current.name,
    slug: config.slug ?? "threadly-app-v1",
    scheme: current.scheme,
    ios: { ...config.ios, bundleIdentifier: current.iosBundleIdentifier },
    android: { ...config.android, package: current.androidPackage },
    extra: { ...config.extra, appVariant: variant },
  }
}

import { spawnSync } from "node:child_process"
import { readFileSync } from "node:fs"
import { join } from "node:path"
import { describe, expect, it } from "vitest"

const cli = join(process.cwd(), "node_modules", "expo", "bin", "cli")
const eas = JSON.parse(readFileSync(join(process.cwd(), "eas.json"), "utf8")) as {
  build: Record<string, { environment: string; env: { APP_VARIANT: string } }>
}

function configFor(variant?: string) {
  const env: NodeJS.ProcessEnv = { ...process.env, EXPO_NO_DOTENV: "1" }
  if (variant === undefined) delete env.APP_VARIANT
  else env.APP_VARIANT = variant
  return spawnSync(process.execPath, [cli, "config", "--json"], { env, encoding: "utf8" })
}

describe("app variants", () => {
  it.each([
    ["development", "development", "Threadly Dev", "threadly-dev", ".dev"],
    ["staging", "preview", "Threadly Staging", "threadly-staging", ".staging"],
    ["production", "production", "Threadly", "threadly", ""],
  ])("resolves the %s build profile", (variant, environment, name, scheme, suffix) => {
    const profile = eas.build[variant]
    expect(profile.environment).toBe(environment)
    expect(profile.env.APP_VARIANT).toBe(variant)
    const result = configFor(profile.env.APP_VARIANT)
    expect(result.status, result.stderr).toBe(0)
    const config = JSON.parse(result.stdout)
    expect(config.name).toBe(name)
    expect(config.scheme).toBe(scheme)
    expect(config.ios.bundleIdentifier).toBe(`com.quindarts.threadly-app-v1${suffix}`)
    expect(config.android.package).toBe(`com.quindarts.threadlyappv1${suffix}`)
    expect(config.extra.appVariant).toBe(variant)
    expect(config.slug).toBe("threadly-app-v1")
    expect(config.ios.supportsTablet).toBe(false)
  })

  it("defaults to development and rejects an unknown variant", () => {
    const defaultConfig = configFor()
    expect(defaultConfig.status, defaultConfig.stderr).toBe(0)
    expect(JSON.parse(defaultConfig.stdout).extra.appVariant).toBe("development")
    const invalid = configFor("preview")
    expect(invalid.status).not.toBe(0)
    expect(invalid.stderr).toContain("Invalid APP_VARIANT")
  })
})

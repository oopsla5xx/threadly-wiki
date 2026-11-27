import { execFileSync, spawnSync } from "node:child_process"
import {
  cpSync,
  existsSync,
  mkdtempSync,
  mkdirSync,
  readFileSync,
  rmSync,
  writeFileSync,
} from "node:fs"
import { tmpdir } from "node:os"
import { join } from "node:path"
import { afterEach, describe, expect, it } from "vitest"

let fixture: string | undefined

afterEach(() => {
  if (fixture) rmSync(fixture, { recursive: true, force: true })
  fixture = undefined
})

function createFixture() {
  fixture = mkdtempSync(join(tmpdir(), "threadly-reset-"))
  mkdirSync(join(fixture, "scripts"))
  mkdirSync(join(fixture, "src", "app"), { recursive: true })
  mkdirSync(join(fixture, "src", "features"))
  cpSync(join(__dirname, "reset-project.js"), join(fixture, "scripts", "reset-project.js"))
  writeFileSync(join(fixture, "src", "app", "old.tsx"), "old route")
  writeFileSync(join(fixture, "src", "features", "old.ts"), "old feature")
  return fixture
}

describe("reset-project", () => {
  it("backs up existing routes and features and creates a fresh app", () => {
    const root = createFixture()
    execFileSync(process.execPath, [join(root, "scripts", "reset-project.js")])
    expect(readFileSync(join(root, "project-backup", "app", "old.tsx"), "utf8")).toBe("old route")
    expect(readFileSync(join(root, "project-backup", "features", "old.ts"), "utf8")).toBe(
      "old feature",
    )
    expect(existsSync(join(root, "src", "app", "index.tsx"))).toBe(true)
    expect(readFileSync(join(root, "src", "app", "_layout.tsx"), "utf8")).toContain("PortalHost")
    execFileSync(process.execPath, [
      join(process.cwd(), "node_modules", "prettier", "bin", "prettier.cjs"),
      "--config",
      join(process.cwd(), ".prettierrc"),
      "--check",
      join(root, "src", "app", "_layout.tsx"),
      join(root, "src", "app", "index.tsx"),
    ])
  })

  it("refuses to overwrite a previous backup", () => {
    const root = createFixture()
    mkdirSync(join(root, "project-backup"))
    const result = spawnSync(process.execPath, [join(root, "scripts", "reset-project.js")])
    expect(result.status).toBe(1)
    expect(result.stderr.toString()).toContain("already exists")
    expect(readFileSync(join(root, "src", "app", "old.tsx"), "utf8")).toBe("old route")
  })

  it("leaves routes untouched if features are missing", () => {
    const root = createFixture()
    rmSync(join(root, "src", "features"), { recursive: true })
    const result = spawnSync(process.execPath, [join(root, "scripts", "reset-project.js")])
    expect(result.status).toBe(1)
    expect(existsSync(join(root, "project-backup"))).toBe(false)
    expect(readFileSync(join(root, "src", "app", "old.tsx"), "utf8")).toBe("old route")
  })
})

# Threadly web

Next.js application using the App Router and TypeScript.

## Getting started

Use Node.js 22 or later and pnpm 11.21:

```sh
pnpm install --frozen-lockfile
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000). Routes are in `src/app`.

Run `pnpm lint` to check the code and `pnpm build` to create a production
build. Run `pnpm start` after building to serve it.

Playwright is configured for Chromium. Install its browser once with
`pnpm exec playwright install chromium`, then run `pnpm test:e2e`.

Tailwind CSS and shadcn/ui are configured in `src/app/globals.css` and
`components.json`. A Button component is available at
`src/components/ui/button.tsx`; add more components with
`pnpm exec shadcn add <component>`. Zod is installed for schema validation.
Framer Motion and Lucide React are available for animations and icons.

The root layout self-hosts Inter with `next/font`. The robots policy in
`src/app/robots.txt` points to `https://precedent.dev/sitemap.xml`; replace it
with the application's sitemap URL before deployment.

The application is a starter and is not yet connected to `threadly-api-v1`.

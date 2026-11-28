# Threadly web

Next.js application using the App Router and TypeScript.

## Getting started

Use Node.js 20.9 or later and Yarn 1:

```sh
yarn install --frozen-lockfile
yarn dev
```

Open [http://localhost:3000](http://localhost:3000). Routes are in `src/app`.

Run `yarn lint` to check the code and `yarn build` to create a production
build. Run `yarn start` after building to serve it.

Tailwind CSS and shadcn/ui are configured in `src/app/globals.css` and
`components.json`. A Button component is available at
`src/components/ui/button.tsx`; add more components with
`yarn shadcn add <component>`. Zod is installed for schema validation.

The application is a starter and is not yet connected to `threadly-api-v1`.

# Threadly app

Expo/React Native application using a feature-based structure.

## Getting started

```sh
yarn install --frozen-lockfile
yarn start
```

Use `yarn ios`, `yarn android`, or `yarn web` to select a platform.
The iOS and Android commands build and run native development apps locally.
Native modules not included in Expo Go require a development build; see
the [Expo development build guide](https://docs.expo.dev/develop/development-builds/introduction/).

## Structure

- `src/app`: Expo Router routes and the root providers.
- `src/features`: feature screens and their feature-specific code.

NativeWind is configured via `global.css`, `tailwind.config.js`, Babel, and Metro.
TanStack Query is available through the root provider. The remaining libraries
are installed but service integrations such as Clerk and PostHog
still need project credentials and application-specific configuration.

React Native Reusables is configured for NativeWind. Add components with
`yarn cli add <component>`.

Run `yarn typecheck`, `yarn lint`, `yarn format:check`, and `yarn test`
to validate changes. Husky installs a repository-wide pre-commit hook that runs
these checks for this app. Set `HUSKY=0` when installing in CI without hooks.

`yarn reset-project` moves `src/app` and `src/features` to `project-backup`
and creates an empty app. It refuses to overwrite an existing backup; move or
remove that backup yourself before resetting again.

## App environments

`APP_VARIANT` selects `development` (default), `staging`, or `production`.
Each variant has its own app name, deep-link scheme, iOS bundle ID, and Android
package. For example, run `APP_VARIANT=staging yarn start` locally. An invalid
value stops config evaluation rather than silently selecting another variant.

EAS build profiles map `development`, `staging`, and `production` to the EAS
environments `development`, `preview`, and `production` respectively. Install
EAS CLI and verify the Expo project link before running
`eas build --profile <profile>`. Configure environment-specific values in EAS;
do not commit secrets. Values exposed to client code are public.

Channels are reserved in the build profiles, but over-the-air updates are not
configured yet. Set up EAS Update for the linked project before using them.
When switching local native build variants, regenerate ignored native folders
with `APP_VARIANT=<variant> yarn expo prebuild --clean` before `yarn ios` or
`yarn android`; review any local native changes first because `--clean` deletes
those folders.

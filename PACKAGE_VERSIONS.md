# Package Versions - Latest Stable (February 2026)

This document lists all package versions used by the CLI tool. All versions are the **latest stable releases** as of February 25, 2026.

## Core Frameworks

| Package | Version | Notes |
|---------|---------|-------|
| `react` | ^19.2.4 | React 19 stable release |
| `react-dom` | ^19.2.4 | React DOM 19 stable release |
| `vue` | ^3.5.29 | Vue 3.5 stable |
| `svelte` | ^5.53.3 | Svelte 5 stable release |
| `solid-js` | ^1.9.11 | Solid 1.9 stable |
| `@angular/core` | ^21.1.5 | Angular 21 stable |

## Meta-Frameworks

Meta-frameworks shell out to upstream CLIs for scaffolding:

| Framework | CLI Command | Post-scaffold deps |
|-----------|------------|-------------------|
| Next.js | `npx create-next-app@latest` | ESLint, state, data fetching |
| Astro | `npm create astro@latest` | Tailwind via @tailwindcss/vite, ESLint |
| SvelteKit | `npx sv create` + `npx sv add` | ESLint, state, data fetching |

## Build Tools

| Package | Version | Notes |
|---------|---------|-------|
| `vite` | ^7.3.1 | Latest Vite 7 release |
| `@vitejs/plugin-react` | ^5.1.4 | Latest Vite React plugin |
| `@vitejs/plugin-vue` | ^6.0.4 | Latest Vite Vue plugin |
| `@sveltejs/vite-plugin-svelte` | ^6.2.4 | Latest Vite Svelte plugin for Svelte 5 |
| `@analogjs/vite-plugin-angular` | ^2.2.3 | AnalogJS Vite plugin for Angular 21 |

## TypeScript

| Package | Version | Notes |
|---------|---------|-------|
| `typescript` | ^5.9.3 | Latest TypeScript 5.9 stable |
| `typescript-eslint` | ^8.56.1 | Latest TypeScript ESLint integration |
| `@types/react` | ^19.2.14 | React 19 type definitions |
| `@types/react-dom` | ^19.2.3 | React DOM 19 type definitions |

## Routing

| Package | Version | Notes |
|---------|---------|-------|
| `react-router` | ^7.13.1 | React Router 7 (replaces react-router-dom) |
| `@tanstack/react-router` | ^1.163.2 | TanStack Router latest |
| `vue-router` | ^4.4.5 | Vue Router 4.x stable (5.0 available but new) |

## Styling

| Package | Version | Notes |
|---------|---------|-------|
| `tailwindcss` | ^4.2.1 | Tailwind CSS 4.x stable |
| `@tailwindcss/vite` | ^4.2.1 | Tailwind Vite plugin |
| `sass` | ^1.97.3 | Dart Sass latest |
| `styled-components` | ^6.3.11 | Styled Components 6.x |
| `bootstrap` | ^5.3.8 | Bootstrap 5.3 stable |

## State Management

| Package | Version | Notes |
|---------|---------|-------|
| `zustand` | ^5.0.11 | Zustand 5.x stable |
| `@reduxjs/toolkit` | ^2.11.2 | Redux Toolkit latest |
| `react-redux` | ^9.2.0 | React Redux 9 |
| `pinia` | ^3.0.4 | Pinia 3.x |

## Data Fetching

| Package | Version | Notes |
|---------|---------|-------|
| `@tanstack/react-query` | ^5.90.21 | TanStack Query v5 latest |
| `@tanstack/react-query-devtools` | ^5.91.3 | TanStack Query DevTools |
| `@tanstack/svelte-query` | ^6.0.18 | TanStack Svelte Query v6 (Svelte 5 runes) |
| `axios` | ^1.13.5 | Axios latest stable |
| `swr` | ^2.4.0 | SWR 2.x latest |

## Testing

| Package | Version | Notes |
|---------|---------|-------|
| `vitest` | ^4.0.18 | Vitest 4.0 with stable browser mode |
| `jest` | ^30.2.0 | Jest 30 (requires Node 18+) |
| `@testing-library/react` | ^16.3.2 | React Testing Library |
| `@testing-library/svelte` | ^5.3.1 | Svelte Testing Library |
| `@testing-library/jest-dom` | ^6.9.1 | Jest DOM matchers |
| `jsdom` | ^28.1.0 | JSDOM for Vitest |

## Linting

| Package | Version | Notes |
|---------|---------|-------|
| `eslint` | ^9.39.1 | ESLint 9.x stable (10.x just released, not adopted yet) |
| `@eslint/js` | ^9.39.1 | ESLint JavaScript config |
| `globals` | ^15.15.0 | Global variables definitions |
| `eslint-plugin-react-hooks` | ^7.0.1 | React hooks rules + React Compiler rules |
| `eslint-plugin-react-refresh` | ^0.5.2 | React Fast Refresh validation |

## UI Component Libraries

| Package | Version | Notes |
|---------|---------|-------|
| `@mui/material` | ^7.3.8 | Material UI v7 latest |
| `@chakra-ui/react` | ^3.33.0 | Chakra UI v3 latest |
| `antd` | ^6.0.0 | Ant Design v6 stable |
| `@angular/material` | ^21.1.5 | Angular Material 21 |
| `ng-zorro-antd` | ^21.1.0 | NG-ZORRO for Angular 21 |

## Form Management & Validation

| Package | Version | Notes |
|---------|---------|-------|
| `react-hook-form` | ^7.71.2 | React Hook Form latest |
| `@hookform/resolvers` | ^5.2.2 | Hookform resolvers latest |
| `@tanstack/react-form` | ^1.28.3 | TanStack Form v1 stable |
| `formik` | ^2.4.9 | Formik latest |
| `vee-validate` | ^4.15.1 | VeeValidate for Vue |
| `zod` | ^4.3.6 | Zod v4 schema validation |
| `yup` | ^1.7.1 | Yup v1 schema validation |

## Animation

| Package | Version | Notes |
|---------|---------|-------|
| `motion` | ^12.34.3 | Motion (formerly Framer Motion) latest |
| `gsap` | ^3.14.2 | GSAP v3 latest |
| `@react-spring/web` | ^10.0.3 | React Spring v10 |

## Data Visualization

| Package | Version | Notes |
|---------|---------|-------|
| `recharts` | ^3.7.0 | Recharts v3 latest |
| `echarts` | ^6.0.0 | Apache ECharts v6 |

## Internationalization

| Package | Version | Notes |
|---------|---------|-------|
| `react-i18next` | ^16.5.4 | React i18next v16 latest |
| `i18next` | ^25.8.13 | i18next v25 latest |
| `vue-i18n` | ^11.2.8 | Vue i18n v11 latest |

## Breaking Changes to Note

### Angular 21
- Jumped from 19.x to 21.x (via 20.x)
- Requires @analogjs/vite-plugin-angular 2.x for Vite integration

### @tanstack/svelte-query v6
- Migrated to Svelte 5 runes syntax
- No longer compatible with Svelte 4

### @hookform/resolvers v5
- Major API changes from v3

### @tanstack/react-form v1
- Stable release (was 0.x pre-release)

### vue-i18n v11
- Major upgrade from v10

### jsdom v28
- Major upgrade from v25

## Last Updated

February 25, 2026

## Verification Commands

To verify package versions in your generated project:

```bash
# Check installed versions
npm list <package-name>

# Check latest available versions
npm outdated

# Update to latest within range
npm update

# Check for security vulnerabilities
npm audit
```

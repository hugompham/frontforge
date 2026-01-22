# Package Versions - Latest Stable (December 2025)

This document lists all package versions used by the CLI tool. All versions are the **latest stable releases** as of December 10, 2025.

## Core Frameworks

| Package | Version | Notes |
|---------|---------|-------|
| `react` | ^19.2.1 | React 19 stable release |
| `react-dom` | ^19.2.1 | React DOM 19 stable release |
| `vue` | ^3.5.13 | Vue 3.5 stable, Vue 3.6 is in alpha |
| `svelte` | ^5.30.0 | Svelte 5 stable release |

## Build Tools

| Package | Version | Notes |
|---------|---------|-------|
| `vite` | ^7.2.7 | Latest Vite 7 release |
| `@vitejs/plugin-react` | ^5.1.2 | Latest Vite React plugin |
| `@vitejs/plugin-vue` | ^6.0.2 | Latest Vite Vue plugin |
| `@sveltejs/vite-plugin-svelte` | ^6.2.1 | Latest Vite Svelte plugin for Svelte 5 |

## TypeScript

| Package | Version | Notes |
|---------|---------|-------|
| `typescript` | ^5.9.3 | Latest TypeScript 5.9 stable |
| `typescript-eslint` | ^8.49.0 | Latest TypeScript ESLint integration |
| `@types/react` | ^19.2.7 | React 19 type definitions |
| `@types/react-dom` | ^19.2.3 | React DOM 19 type definitions |

## Routing

| Package | Version | Notes |
|---------|---------|-------|
| `react-router` | ^7.10.1 | React Router 7 (replaces react-router-dom) |
| `@tanstack/react-router` | ^1.140.2 | TanStack Router latest |
| `vue-router` | ^4.6.3 | Vue Router 4.6 stable |

## Styling

| Package | Version | Notes |
|---------|---------|-------|
| `tailwindcss` | ^4.1.17 | Tailwind CSS 4.x stable |
| `autoprefixer` | ^10.4.22 | Latest autoprefixer |
| `postcss` | ^8.5.6 | PostCSS 8.5 stable |
| `sass` | ^1.95.0 | Dart Sass latest (node-sass is deprecated) |
| `styled-components` | ^6.1.15 | Styled Components 6.x |

## State Management

| Package | Version | Notes |
|---------|---------|-------|
| `zustand` | ^5.0.9 | Zustand 5.x stable |
| `@reduxjs/toolkit` | ^2.11.1 | Redux Toolkit latest |
| `react-redux` | ^9.2.0 | React Redux 9 (requires React 18+) |
| `pinia` | ^3.0.4 | Pinia 3.x (dropped Vue 2 support) |

## Data Fetching

| Package | Version | Notes |
|---------|---------|-------|
| `@tanstack/react-query` | ^5.90.12 | TanStack Query v5 latest |
| `@tanstack/react-query-devtools` | ^5.91.1 | TanStack Query DevTools |
| `axios` | ^1.7.9 | Axios latest stable |
| `swr` | ^2.3.2 | SWR 2.x latest |

## Testing

| Package | Version | Notes |
|---------|---------|-------|
| `vitest` | ^4.0.15 | Vitest 4.0 with stable browser mode |
| `jest` | ^30.2.0 | Jest 30 (requires Node 18+) |
| `@testing-library/react` | ^16.3.0 | React Testing Library |
| `@testing-library/jest-dom` | ^6.9.1 | Jest DOM matchers |
| `jsdom` | ^25.0.1 | JSDOM for Vitest |

## Linting

| Package | Version | Notes |
|---------|---------|-------|
| `eslint` | ^9.39.1 | ESLint 9.x stable |
| `@eslint/js` | ^9.39.1 | ESLint JavaScript config |
| `globals` | ^15.15.0 | Global variables definitions |
| `eslint-plugin-react-hooks` | ^7.0.1 | React hooks rules + React Compiler rules |
| `eslint-plugin-react-refresh` | ^0.4.24 | React Fast Refresh validation |

## UI Component Libraries

| Package | Version | Notes |
|---------|---------|-------|
| `@mui/material` | ^7.4.0 | Material UI v7 latest |
| `@chakra-ui/react` | ^3.4.0 | Chakra UI v3 latest |
| `antd` | ^6.0.0 | Ant Design v6 stable |
| `tailwind-merge` | ^3.4.0 | Tailwind class merging utility |
| `primeng` | ^21.0.1 | PrimeNG for Angular latest |

## Form Management & Validation

| Package | Version | Notes |
|---------|---------|-------|
| `react-hook-form` | ^7.68.0 | React Hook Form latest |
| `formik` | ^2.4.9 | Formik latest |
| `vee-validate` | ^4.15.1 | VeeValidate for Vue latest |
| `zod` | ^4.1.13 | Zod v4 schema validation |
| `yup` | ^1.7.1 | Yup v1 schema validation |

## Animation

| Package | Version | Notes |
|---------|---------|-------|
| `framer-motion` | ^12.23.26 | Framer Motion latest |
| `gsap` | ^3.14.1 | GSAP v3 latest |

## Data Visualization

| Package | Version | Notes |
|---------|---------|-------|
| `recharts` | ^3.5.1 | Recharts v3 latest |
| `echarts` | ^6.0.0 | Apache ECharts v6 latest |

## Internationalization

| Package | Version | Notes |
|---------|---------|-------|
| `react-i18next` | ^16.4.1 | React i18next v16 latest |
| `i18next` | ^25.7.2 | i18next v25 latest |
| `vue-i18n` | ^10.0.8 | Vue i18n v10 latest |

## Version Selection Criteria

All package versions follow these criteria:

1. **Latest Stable** - Only production-ready stable releases
2. **Major Version Consistency** - Compatible with framework major versions
3. **Ecosystem Compatibility** - Tested combinations that work together
4. **Long-term Support** - Actively maintained with security updates

## Breaking Changes to Note

### React 19
- New features: Actions API, Server Components
- Some packages may not yet have full React 19 support

### Vite 7
- Requires Node.js 20.19+ or 22.12+
- Default browser target changed to 'baseline-widely-available'

### TypeScript 5.9
- Improved inference and type checking
- Better support for modern JavaScript features

### Tailwind CSS 4
- CSS-first configuration (no more JavaScript config)
- Up to 5x faster builds
- Requires Safari 16.4+, Chrome 111+, Firefox 128+

### Vitest 4
- Browser mode is now stable
- Built-in visual regression testing
- Requires Node 20+

### Jest 30
- Dropped Node 14, 16, 19, 21 support
- Minimum Node version: 18.x
- Better performance and memory management

### Pinia 3
- Dropped Vue 2 support completely
- Fully focused on Vue 3

### Zustand 5
- Dropped React <18 support
- Uses native useSyncExternalStore
- Smaller bundle size

## Last Updated

December 10, 2025

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

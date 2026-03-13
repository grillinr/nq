// https://docs.expo.dev/guides/using-eslint/
const path = require('path');
const { defineConfig } = require('eslint/config');
const { FlatCompat } = require('@eslint/eslintrc');
const expoConfig = require('eslint-config-expo/flat');
const prettierPlugin = require('eslint-plugin-prettier');
const prettierConfig = require('eslint-config-prettier');
const reactNativePlugin = require('eslint-plugin-react-native');

const compat = new FlatCompat({ baseDirectory: path.resolve(__dirname) });

// Plugins already registered by eslint-config-expo — must not be re-registered
const EXPO_PLUGINS = new Set(['import', 'expo', '@typescript-eslint', 'react', 'react-hooks']);

// Strip only the plugins that expo already owns; keep others (e.g. jsx-a11y)
const airbnbConfigs = compat.extends('airbnb', 'airbnb-typescript').map(config => {
  if (!config.plugins) return config;
  const filtered = Object.fromEntries(
    Object.entries(config.plugins).filter(([name]) => !EXPO_PLUGINS.has(name))
  );
  return { ...config, plugins: filtered };
});

// @typescript-eslint rules that were removed in v8 (all formatting — handled by Prettier)
const REMOVED_TS_RULES = {
  '@typescript-eslint/brace-style': 'off',
  '@typescript-eslint/comma-dangle': 'off',
  '@typescript-eslint/comma-spacing': 'off',
  '@typescript-eslint/func-call-spacing': 'off',
  '@typescript-eslint/indent': 'off',
  '@typescript-eslint/keyword-spacing': 'off',
  '@typescript-eslint/lines-between-class-members': 'off',
  '@typescript-eslint/no-extra-parens': 'off',
  '@typescript-eslint/no-extra-semi': 'off',
  '@typescript-eslint/space-before-blocks': 'off',
  '@typescript-eslint/no-throw-literal': 'off',
  '@typescript-eslint/quotes': 'off',
  '@typescript-eslint/semi': 'off',
  '@typescript-eslint/space-before-function-paren': 'off',
  '@typescript-eslint/space-infix-ops': 'off',
  '@typescript-eslint/object-curly-spacing': 'off',
};

module.exports = defineConfig([
  // Expo base config (React, React Native, TypeScript, import)
  expoConfig,

  // Airbnb + TypeScript rules (conflicting plugin registrations stripped)
  ...airbnbConfigs,

  // React Native specific rules
  {
    plugins: {
      'react-native': reactNativePlugin,
    },
    rules: {
      'react-native/no-inline-styles': 'warn',
      // no-unused-styles produces false positives for the createStyles(colors) factory
      // pattern used throughout this codebase — the plugin cannot statically resolve
      // styles returned from a memoized factory function.
      'react-native/no-unused-styles': 'off',
      'react-native/sort-styles': 'off',
      'react-native/no-color-literals': 'warn',
    },
  },

  // Prettier — disables conflicting formatting rules, then enforces Prettier
  prettierConfig,
  {
    plugins: {
      prettier: prettierPlugin,
    },
    rules: {
      'prettier/prettier': 'error',
    },
  },

  // Project-specific overrides
  {
    rules: {
      // Silence stale @typescript-eslint v7 rules removed in v8 (Prettier handles them)
      ...REMOVED_TS_RULES,

      // Disable typed-linting rules — these require parserOptions.project which
      // is not configured for React Native / Expo bundler setups
      '@typescript-eslint/dot-notation': 'off',
      '@typescript-eslint/no-implied-eval': 'off',
      '@typescript-eslint/return-await': 'off',
      '@typescript-eslint/naming-convention': 'off',

      // Console statements are common in RN development
      'no-console': 'warn',

      // React 17+ JSX transform: no need to import React in every file
      'react/react-in-jsx-scope': 'off',
      'react/jsx-uses-react': 'off',

      // Expo Router requires default exports for screen files
      'import/prefer-default-export': 'off',

      // Bundler handles resolution; don't require explicit extensions
      'import/extensions': [
        'error',
        'ignorePackages',
        { ts: 'never', tsx: 'never', js: 'never', jsx: 'never' },
      ],

      // Prop spreading is a common pattern in RN component wrappers
      'react/jsx-props-no-spreading': 'off',

      // React Native is not a web environment — skip web-only a11y rules
      'jsx-a11y/accessible-emoji': 'off',

      // TypeScript enforces prop types at compile time — runtime defaultProps are redundant
      'react/require-default-props': 'off',

      // The createStyles(colors) factory pattern defines styles after the component
      // that uses them — this is idiomatic React Native and perfectly readable
      '@typescript-eslint/no-use-before-define': ['error', { functions: false, classes: false }],

      // GraphQL __typename fields and Apollo private fields use underscore prefix
      'no-underscore-dangle': ['error', { allow: ['__typename'] }],

      // for...of loops are idiomatic JS/TS; the regenerator concern doesn't apply to RN (Hermes)
      'no-restricted-syntax': [
        'error',
        {
          selector: 'LabeledStatement',
          message: 'Labels are a form of GOTO; using them makes code hard to read and maintain.',
        },
        {
          selector: 'WithStatement',
          message: '`with` is disallowed in strict mode because it makes code unpredictable.',
        },
      ],

      // continue is used in imperative loops in lib code — warn but don't error
      'no-continue': 'warn',

      // Nested ternaries are common in RN conditional rendering — warn only
      'no-nested-ternary': 'warn',

      // Inline tabBarIcon functions in Expo Router/React Navigation are passed as
      // render props (not rendered directly), so allowAsProps is required here.
      'react/no-unstable-nested-components': ['error', { allowAsProps: true }],

      // parseInt is always called with a radix in this codebase via the fix pass
      radix: 'error',
    },
  },

  // Ignore generated and build output
  {
    ignores: ['dist/*', 'node_modules/**', 'src/__generated__/**', '__generated__/**'],
  },
]);

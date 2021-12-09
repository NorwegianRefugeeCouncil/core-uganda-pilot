module.exports = {
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 2020,
  },
  settings: {
    react: {
      version: 'detect',
    },
    'import/resolver': {
      typescript: {},
      node: {
        paths: ['src'],
        extensions: ['.js', '.jsx', '.ts', '.tsx'],
      },
    },
  },
  plugins: ['jest', 'prettier', 'react', 'react-hooks', '@typescript-eslint'],
  extends: [
    'plugin:import/errors',
    'plugin:import/warnings',
    'plugin:import/typescript',
    'plugin:@typescript-eslint/recommended',
    'prettier/@typescript-eslint',
    'airbnb',
    'plugin:prettier/recommended',
    'plugin:react/recommended',
    'prettier/react',
  ],
  globals: {
    Atomics: 'readonly',
    SharedArrayBuffer: 'readonly',
  },
  rules: {
    'comma-dangle': ['error', 'always-multiline'],
    'no-await-in-loop': [0],
    quotes: ['error', 'single', { avoidEscape: true }],
    'object-curly-spacing': ['error', 'always'],
    'import/no-extraneous-dependencies': 'off',
    'import/order': [
      'error',
      {
        groups: ['builtin', 'external', 'internal', 'parent', 'sibling', 'index'],
        'newlines-between': 'always',
      },
    ],
    'max-len': ['error', { code: 100, tabWidth: 2, ignorePattern: '^import\\s.+\\sfrom\\s.+;$' }],
    'no-plusplus': ['error', { allowForLoopAfterthoughts: true }],
    'no-unused-vars': ['error', { varsIgnorePattern: '_' }],
    'import/prefer-default-export': 0,
    'react/jsx-filename-extension': 0,
    'react/require-default-props': 0,
    'import/extensions': [
      'error',
      'ignorePackages',
      {
        js: 'never',
        mjs: 'never',
        jsx: 'never',
        ts: 'never',
        tsx: 'never',
      },
    ],
    'react/jsx-props-no-spreading': 'off',
    'react-hooks/rules-of-hooks': 'error', // Checks rules of Hooks
    'react/no-array-index-key': 'off',
    'react/sort-comp': [
      'error',
      {
        order: [
          '/props/',
          '/state/',
          'everything-else',
          'instance-variables',
          'static-methods',
          'lifecycle',
          'getters',
          'setters',
          'instance-methods',
          'render',
        ],
      },
    ],
    'react/state-in-constructor': 'off',
    'react/static-property-placement': 'off',
  },
};

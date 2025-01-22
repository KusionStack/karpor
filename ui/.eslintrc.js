module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:react/recommended',
    'plugin:prettier/recommended',
  ],
  overrides: [
    {
      files: ['.eslintrc.{js,cjs}'],
      parserOptions: {
        sourceType: 'script',
      },
    },
  ],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
  },
  plugins: ['@typescript-eslint', 'react', 'react-hooks'],
  settings: {
    react: {
      version: 'detect',
    },
  },
  globals: {},
  rules: {
    // js rules：http://eslint.cn/docs/rules/
    /** @js */
    quotes: 'off',
    semi: 'off',
    'no-undef': 'off',
    'no-var': 'off',
    'no-debugger': 'off',
    'no-console': 'off',
    'no-empty': 'off',
    'no-unsafe-optional-chaining': 'off',
    'no-unused-vars': 'off',

    // ts rules：https://typescript-eslint.io/rules/
    /** @typescript */
    '@typescript-eslint/no-unused-vars': 'off',
    '@typescript-eslint/no-explicit-any': 'off',

    /** @react */
    'react-refresh/only-export-components': 'off',
    'react/no-unescaped-entities': 'off',
    'react-hooks/exhaustive-deps': 'off',
    'react/display-name': 'off',
    'react/prop-types': 'off',
    'jsx-a11y/anchor-has-rel': 'off',
  },
}

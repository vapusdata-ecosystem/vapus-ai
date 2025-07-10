import { FlatCompat } from "@eslint/eslintrc";
import { dirname } from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({
  baseDirectory: __dirname,
});

const eslintConfig = [
    ...compat.config({
    extends: ['next'],
    rules: {
      // Ignore missing deps for hooks
    "react-hooks/exhaustive-deps": "off",
    // Ignore using <img>
    "@next/next/no-img-element": "off",
    // Ignore synchronous scripts
    "@next/next/no-sync-scripts": "off",
    // Ignore unescaped entities in JSX
    "react/no-unescaped-entities": "off",
    },
  }),
];

export default eslintConfig;

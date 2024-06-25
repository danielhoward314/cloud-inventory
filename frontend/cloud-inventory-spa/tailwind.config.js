import daisyui from "daisyui"
import typography from '@tailwindcss/typography';
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './public/**/*.html',
    './src/**/*.{js,jsx,ts,tsx,vue}',
  ],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: ["nord"],
  },
  plugins: [
    typography,
    daisyui,
  ],
}


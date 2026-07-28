import type { Config } from 'tailwindcss'
const config: Config = {
  darkMode: 'class',
  content: ['./app/**/*.{js,ts,jsx,tsx,mdx}', './components/**/*.{js,ts,jsx,tsx,mdx}'],
  theme: {
    extend: {
      colors: {
        primary: { 50: '#f0f9ff', 100: '#e0f2fe', 500: '#0ea5e9', 600: '#0284c7', 700: '#0369a1', 800: '#075985', 900: '#0c3d66' },
        success: { 50: '#f0fdf4', 600: '#16a34a', 700: '#15803d' },
        warning: { 50: '#fffbeb', 600: '#f59e0b', 700: '#d97706' },
        error: { 50: '#fef2f2', 600: '#dc2626', 700: '#b91c1c' },
        neutral: { 50: '#fafafa', 100: '#f5f5f5', 200: '#e5e5e5', 300: '#d4d4d4', 400: '#a3a3a3', 500: '#737373', 600: '#525252', 700: '#404040', 800: '#262626', 900: '#171717' },
      },
    },
  },
  plugins: [],
}
export default config

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{ts,tsx,js,jsx,mdx}'],
  theme: {
    extend: {
      colors: {
        bg: '#0b0b10',
        panel: '#13131c',
        muted: '#1c1c28',
        brand: '#ff2e63',
        ink: '#e8e8ea',
        sub: '#9aa0a6',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
};

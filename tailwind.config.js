/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.html",
    "./static/js/**/*.js",
  ],
  theme: {
    extend: {
      colors: {
        'sage': {
          50: '#F5F7F5',
          100: '#E8ECE9',
          200: '#C8D4CC',
          300: '#A7BBAF',
          400: '#7C9885',
          500: '#5A6B61',
          600: '#475650',
          700: '#384441',
          800: '#2B3432',
          900: '#1F2523',
        },
        'blush': {
          400: '#E88D9D',
          500: '#E17589',
        },
        'cream': '#FDFCF8',
      },
      fontFamily: {
        'sans': ['Inter', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'sans-serif'],
      },
      maxWidth: {
        'content': '650px',
      },
    },
  },
  plugins: [],
}
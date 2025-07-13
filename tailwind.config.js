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
          600: '#4A7C59',
          700: '#3E6B4A',
          800: '#2B3432',
          900: '#1F2523',
        },
        'blush': {
          400: '#E88D9D',
          500: '#E17589',
        },
        'terracotta': {
          100: '#FDF2F0',
          200: '#F9E1DD',
          300: '#E6A898',
          400: '#D4805F',
          500: '#B85C3E',
          600: '#9A4E34',
        },
        'slate-blue': {
          100: '#F1F5F9',
          200: '#E2E8F0',
          300: '#94A3B8',
          400: '#5B7FBD',
          500: '#4F68A6',
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
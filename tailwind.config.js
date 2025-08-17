/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.html",
    "./static/js/**/*.js",
  ],
  theme: {
    extend: {
      fontFamily: {
        'sans': ['Berkeley Mono', 'monospace'],
      },
      maxWidth: {
        'content': '650px',
      },
    },
  },
  plugins: [],
}
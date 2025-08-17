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
    },
  },
  plugins: [],
}
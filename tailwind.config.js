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
      colors: {
        jv: {
          DEFAULT: '#6b7c4c',
          light: '#7d9a78',
        },
      },
    },
  },
  plugins: [],
}
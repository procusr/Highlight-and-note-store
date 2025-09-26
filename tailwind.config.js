/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/**/*.{html,templ}",
    "./static/**/*.{html,js}",
    "./internal/**/*.go",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};

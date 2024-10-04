/** @type {import('tailwindcss').Config} */
module.exports = {
  content: {
    files: ["./views/**/*.go", "./components/**/*.go"],
  },
  safelist: [
    {
      pattern: /text-[2-9]xl/,
    },
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}

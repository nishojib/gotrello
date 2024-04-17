import plugin from 'tailwindcss/plugin';

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./**/*.html', './**/*.templ', './**/*.go'],
  theme: {
    extend: {},
  },
  plugins: [
    require('tailwind-scrollbar'),
    require('daisyui'),
    plugin(function ({ addVariant }) {
      addVariant('htmx-settling', ['&.htmx-settling', '.htmx-settling &']);
      addVariant('htmx-request', ['&.htmx-request', '.htmx-request &']);
      addVariant('htmx-swapping', ['&.htmx-swapping', '.htmx-swapping &']);
      addVariant('htmx-added', ['&.htmx-added', '.htmx-added &']);
    }),
  ],
  daisyui: {
    themes: ['night'],
  },
};

const colors = require('tailwindcss/colors')

module.exports = {
    purge: [
        "../*.go",
        "../*.html"
    ],
    mode: 'jit',
    darkMode: 'media', // or 'media' or 'class'
    theme: {
        extend: {
            colors: {
                gray: colors.trueGray
            }
        },
    },
    variants: {},
    plugins: [],
}
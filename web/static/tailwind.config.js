const colors = require('tailwindcss/colors')
const dc = require('tailwindcss/defaultConfig')

module.exports = {
    purge: [
        "../*.go",
        "../*.html"
    ],
    mode: 'jit',
    darkMode: 'media', // or 'media' or 'class'
    theme: {
        extend: {
            fontFamily: {
                serif: ['Times'].concat(dc.theme.fontFamily.serif)
            },
            colors: {
                gray: colors.trueGray
            }
        },
    },
    variants: {},
    plugins: [],
}
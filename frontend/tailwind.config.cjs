/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                'fox-orange': '#FF8C00',
                'fox-charcoal': '#141419', // Dark background
                'fox-dark': '#2C2C35', // Panel background
                'fox-light': '#F5F5F5', // Text
            },
            fontFamily: {
                sans: ['Inter', 'sans-serif'],
            },
        },
    },
    plugins: [],
}

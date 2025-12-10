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
                'fox-charcoal': '#141419',
                'fox-dark': '#2C2C35',
                'fox-light': '#F5F5F5',
                sudoku: {
                    board: '#FDF6E3',
                    'cell-bg': '#FDF6E3',
                    'cell-given': '#E5E5E5',
                    primary: '#FF9F43',
                    'primary-dark': '#D68D38',
                    'primary-darker': '#B57025',
                    'primary-light': '#FFD28F',
                    highlight: '#FFEAA7',
                    peer: '#FFE0B2',
                    text: '#2D3436',
                    'text-secondary': '#636e72',
                    panel: '#323846',
                    'panel-dark': '#1E222D',
                    teal: '#00b894',
                    'teal-light': '#00cec9',
                    grid: '#B2BEC3',
                }
            },
            fontFamily: {
                sans: ['Inter', 'sans-serif'],
            },
        },
    },
    plugins: [],
}

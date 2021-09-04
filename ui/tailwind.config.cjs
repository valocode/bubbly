
const colors = require('tailwindcss/colors');
const config = {
	mode: "jit",
	purge: [
		'tailwind-safelist.txt',
		"./src/**/*.{html,js,svelte,ts}",
	],
	theme: {
		extend: {
			colors: {
				cyan: colors.cyan,
				orange: colors.orange,
				amber: colors.amber,
				bubbly: {
					blue: '#0000b4',
					light: '#14f3ff',
					dark: '#292a5f',
					darker: '#000038'
				}
			}
		}
	},
	plugins: [require('@tailwindcss/forms')]
};

module.exports = config;

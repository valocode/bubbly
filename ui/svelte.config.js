import preprocess from 'svelte-preprocess';
import adapter from '@sveltejs/adapter-static';
import path from 'path';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://github.com/sveltejs/svelte-preprocess
	// for more information about preprocessors
	preprocess: [preprocess({
		"postcss": true
	})],

	kit: {
		// The UI is hosted at path /ui in the bubbly server, and this tells
		// svelte to use /ui as the prefix so that files can be found (e.g. js/css).
		// Due to some current issues in sveltekit, these lines should be commented
		// out during dev mode
		// paths: {
		// 	base: '/ui',
		// },
		adapter: adapter({
			// Run in "true" SPA mode
			fallback: 'index.html'
		}),
		// hydrate the <div id="svelte"> element in src/app.html
		target: '#svelte',
		vite: {
			resolve: {
				alias: {
					$stores: path.resolve('src/stores'),
					$schema: path.resolve('src/schema'),
					$types: path.resolve('src/types'),
					$utils: path.resolve('src/utils'),
				}
			},
			optimizeDeps: {
				exclude: ['@urql/svelte']
			}
		}
	},
};

export default config;

import preprocess from 'svelte-preprocess';
import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://github.com/sveltejs/svelte-preprocess
	// for more information about preprocessors
	preprocess: preprocess(),

	kit: {
		// The UI is hosted at path /ui in the bubbly server, and this tells
		// svelte to use /ui as the prefix so that files can be found (e.g. js/css)
		paths: {
			base: '/ui',
		},
		adapter: adapter({
			// Run in "true" SPA mode
			fallback: 'index.html'
		}),
		// hydrate the <div id="svelte"> element in src/app.html
		target: '#svelte'
	}
};

export default config;

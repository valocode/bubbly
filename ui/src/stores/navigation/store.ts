import { writable } from 'svelte/store';

function createNavigationStore() {
	const { subscribe, set, update } = writable<string>('Insight');

	return {
		subscribe,
		set,
		update
	};
}

export const navigationStore = createNavigationStore();

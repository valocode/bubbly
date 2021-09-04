import { writable } from 'svelte/store';

export interface notification {
	type?: string;
	name?: string;
	value?: boolean;
	title?: string;
	status?: boolean;
	show?: boolean;
}

function createNotificationStore() {
	const { subscribe, set, update } = writable<notification>({});

	return {
		subscribe,
		set,
		update
	};
}

export const notificationStore = createNotificationStore();

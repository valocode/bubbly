<script lang="ts">
	import '../app.postcss';
	import { createClient, setClient } from '@urql/svelte';
	import { bubblyAPI } from '$stores/http/env';
	import Notification from '$lib/Notification.svelte';
	import { notificationStore } from '$stores/notification/store';

	import Sidebar from '$lib/Sidebar.svelte';

	// Setup the urql graphql client
	const client = createClient({
		url: bubblyAPI + '/graphql'
		// fetchOptions: () => {
		// 	return {
		// 		headers: {
		// 			authorization: `Bearer ${token}`
		// 		}
		// 	};
		// }
	});
	setClient(client);
</script>

{#if $notificationStore}
	{#if $notificationStore.show}
		<Notification settings={$notificationStore} />
	{/if}
{/if}
<!-- <div class="h-screen flex overflow-hidden bg-gray-100"> -->
<Sidebar>
	<slot />
</Sidebar>
<!-- </div> -->

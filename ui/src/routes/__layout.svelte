<script lang="ts">
	import '../app.postcss';
	import { bubblyAPI } from '$stores/http/env';
	import Notification from '$lib/Notification.svelte';
	import { notificationStore } from '$stores/notification/store';

	import Sidebar from '$lib/Sidebar.svelte';
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

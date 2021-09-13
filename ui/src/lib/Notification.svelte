<script lang="ts">
	import { afterUpdate } from 'svelte';
	import { fly } from 'svelte/transition';
	import { notificationStore } from '$stores/notification/store';
	export let settings;

	const name = settings.name;
	const type = settings.type;
	const status = settings.status;
	let title = settings?.title;
	export let timeout = 4000;

	if (!title) {
		title = status ? `${type} '${name}' created` : `${type} '${name}' could not be created`;
	}

	afterUpdate(() => {
		// after some time, remove the notification from view and reset the
		// notification store
		setTimeout(function () {
			notificationStore.set({});
		}, timeout);
	});
</script>

<!-- Global notification live region, render this permanently at the end of the document -->
<div
	aria-live="assertive"
	class="z-50 fixed inset-0 flex items-end px-4 py-6 pointer-events-none sm:p-6 sm:items-start"
>
	<div
		in:fly={{ duration: 100 }}
		out:fly={{ duration: 300 }}
		class="w-full flex flex-col items-center space-y-4 sm:items-end"
	>
		<div
			class="max-w-sm w-full bg-white shadow-lg rounded-lg pointer-events-auto ring-1 ring-black ring-opacity-5 overflow-hidden"
		>
			<div class="p-4">
				<div class="flex items-start">
					<div class="flex-shrink-0">
						{#if status}
							<!-- Heroicon name: outline/check-circle -->
							<svg
								class="h-6 w-6 text-green-400"
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								aria-hidden="true"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
								/>
							</svg>
						{:else}
							<!-- Heroicon name: outline/x-circle -->
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-6 w-6 text-red-400"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
								aria-hidden="true"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
								/>
							</svg>
						{/if}
					</div>
					<div class="ml-3 w-0 flex-1 pt-0.5">
						<p class="text-sm font-medium text-gray-900">{title}</p>
					</div>
					<div class="ml-4 flex-shrink-0 flex">
						<button
							on:click={() => {
								notificationStore.update((value) => {
									value.show = false;
									return value;
								});
							}}
							class="bg-white rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
						>
							<span class="sr-only">Close</span>
							<!-- Heroicon name: solid/x -->
							<svg
								class="h-5 w-5"
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 20 20"
								fill="currentColor"
								aria-hidden="true"
							>
								<path
									fill-rule="evenodd"
									d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
									clip-rule="evenodd"
								/>
							</svg>
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

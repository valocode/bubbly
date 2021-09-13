<script lang="ts">
	import Icon from '$lib/Icon.svelte';
	import { onMount } from 'svelte';
	import { SortFields, sortReleases } from './stores';

	let show: boolean = false;
	let sortMenu = null;
	let sortMenuList = null;

	// TODO: we should streamline this and not duplicate it in every component
	// that needs it
	// source: https://codechips.me/tailwind-ui-react-vs-svelte/
	onMount(() => {
		const handleOutsideClick = (event) => {
			if (show && !sortMenu.contains(event.target) && !sortMenuList.contains(event.target)) {
				show = false;
			}
		};

		const handleEscape = (event) => {
			if (show && event.key === 'Escape') {
				show = false;
			}
		};

		// add events when element is added to the DOM
		document.addEventListener('click', handleOutsideClick, false);
		document.addEventListener('keyup', handleEscape, false);

		// remove events when element is removed from the DOM
		return () => {
			document.removeEventListener('click', handleOutsideClick, false);
			document.removeEventListener('keyup', handleEscape, false);
		};
	});
</script>

<button
	type="button"
	class="w-full bg-white border border-gray-300 rounded-md shadow-sm px-4 py-2 inline-flex justify-center text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
	id="sort-menu-button"
	aria-expanded="false"
	aria-haspopup="true"
	on:click={() => (show = !show)}
	bind:this={sortMenu}
>
	<Icon name="solid-sort-ascending" class="mr-3 h-5 w-5 text-gray-400" />
	Sort
	<Icon name="solid-chevron-down" class="ml-2.5 -mr-1.5 h-5 w-5 text-gray-400" />
</button>
<!-- Dropdown menu, show/hide based on menu state. -->
{#if show}
	<div
		class="origin-top-right z-10 absolute right-0 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none"
		role="menu"
		aria-orientation="vertical"
		aria-labelledby="sort-menu-button"
		tabindex="-1"
		bind:this={sortMenuList}
	>
		<div class="py-1" role="none">
			{#each SortFields as field}
				<button
					type="button"
					href=""
					class="{field === sortReleases.sortByField()
						? 'bg-gray-100 text-gray-900'
						: 'text-gray-700'} block px-4 py-2 text-sm w-full text-left"
					role="menuitem"
					on:click={() => {
						sortReleases.set(field);
						show = !show;
					}}
					tabindex="-1"
					id="sort-menu-item-0"
					>{field}
				</button>
			{/each}
		</div>
	</div>
{/if}

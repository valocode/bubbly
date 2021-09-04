<!--
  Custom select controls like this require a considerable amount of JS to implement from scratch. We're planning
  to build some low-level libraries to make this easier with popular frameworks like React, Vue, and even Alpine.js
  in the near future, but in the mean time we recommend these reference guides when building your implementation:

  https://www.w3.org/TR/wai-aria-practices/#Listbox
  https://www.w3.org/TR/wai-aria-practices/examples/listbox/listbox-collapsible.html
-->
<script lang="ts">
	import type { Project, Project_Json } from '$schema/schema_gen';

	import { onMount } from 'svelte';

	import { filterProjects } from './stores';
	import type { ProjectSelectMap } from './stores';
	import Icon from '$lib/Icon.svelte';
	import { operationStore, query } from '@urql/svelte';

	let show: boolean = false;
	let projectMenu = null;
	let projectMenuList = null;
	let projectSelect = [];
	let selectedProjects: string[] = [];

	filterProjects.subscribe((projects) => {
		projectSelect = Object.entries(projects);
		selectedProjects = filterProjects.selectedProjects();
	});

	// TODO: we should streamline this and not duplicate it in every component
	// that needs it
	// source: https://codechips.me/tailwind-ui-react-vs-svelte/
	onMount(() => {
		const handleOutsideClick = (event) => {
			if (show && !projectMenu.contains(event.target) && !projectMenuList.contains(event.target)) {
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

	const projectQuery = `
	{
		project {
			name
		}
	}`;
	const store = operationStore<Project_Json>(projectQuery);
	query(store);

	store.subscribe((value) => {
		if (!value.fetching && !value.error) {
			let projectList: ProjectSelectMap = {};
			value.data.project.forEach((p) => {
				projectList[p.name] = false;
			});
			filterProjects.set(projectList);
		}
	});
</script>

<span class="ml-3" id="release-filter-head-only">
	<span class="text-sm font-medium text-gray-900">Project: </span>
</span>
<button
	type="button"
	on:click={() => (show = !show)}
	bind:this={projectMenu}
	class="relative bg-white w-60 border border-gray-300 rounded-md shadow-sm pl-3 pr-10 py-2 text-left cursor-default focus:outline-none focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
	aria-haspopup="listbox"
	aria-expanded="true"
	aria-labelledby="listbox-label"
>
	<span class="block truncate">
		{selectedProjects.length === 0 ? 'All Projects' : selectedProjects.join(', ')}
	</span>
	<span class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
		<!-- Heroicon name: solid/selector -->
		<Icon name="solid-selector" class="h-5 w-5 text-gray-400" />
	</span>
</button>

<!--
        Select popover, show/hide based on select state.
  
        Entering: ""
          From: ""
          To: ""
        Leaving: "transition ease-in duration-100"
          From: "opacity-100"
          To: "opacity-0"
      -->
{#if show}
	<ul
		class="absolute z-10 mt-1 w-96 bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm"
		bind:this={projectMenuList}
		tabindex="-1"
		role="listbox"
		aria-labelledby="listbox-label"
		aria-activedescendant="listbox-option-3"
	>
		<!--
          Select option, manage highlight styles based on mouseenter/mouseleave and keyboard navigation.
  
          Highlighted: "text-white bg-indigo-600", Not Highlighted: "text-gray-900"
        -->
		{#each projectSelect as [name, selected]}
			<li
				class="text-gray-900 cursor-default select-none relative py-2 pl-8 pr-4"
				id="listbox-option-0"
				role="option"
				on:click={() => filterProjects.toggle(name)}
			>
				<!-- Selected: "font-semibold", Not Selected: "font-normal" -->
				<span class="{selected ? 'font-semibold' : 'font-normal'} block truncate">
					{name}
				</span>

				<!--
            Checkmark, only display for selected option.
  
            Highlighted: "text-white", Not Highlighted: "text-indigo-600"
          -->
				<span
					class="{selected
						? 'text-indigo-600'
						: 'text-white'} absolute inset-y-0 left-0 flex items-center pl-1.5"
				>
					<Icon name="solid-check" class="h-5 w-5" />
				</span>
			</li>
		{/each}
	</ul>
{/if}

<script lang="ts">
	import Icon from '$lib/Icon.svelte';

	import { httpStore } from '$stores/http/store';
	import HeadFilterToggle from './_headFilterToggle.svelte';
	import ReleaseSort from './_releaseSort.svelte';
	import Loading from '$lib/Loading.svelte';
	import ReleaseRow from './_releaseRow.svelte';
	import { filterHeadOnly, filterProjects, sortReleases } from './stores';
	import ProjectFilter from './projectFilter.svelte';

	export let selectedRelease;

	const queryParams = (): object => {
		let params = {
			head_only: filterHeadOnly.isEnabled(),
			projects: filterProjects.selectedProjects().join(',')
		};
		if (sortReleases.sortByField() !== null) {
			params['sort_by'] = sortReleases.sortByField().toLowerCase();
		}
		return params;
	};

	const store = httpStore(null);
	store.get('/releases', queryParams()).catch((err) => {});

	// Subscribe to changes from the different filter stores we have
	filterHeadOnly.subscribe((b) => {
		store.get('/releases', queryParams()).catch((err) => {});
	});

	filterProjects.subscribe((p) => {
		store.get('/releases', queryParams()).catch((err) => {});
	});

	sortReleases.subscribe((s) => {
		store.get('/releases', queryParams()).catch((err) => {});
	});
</script>

<div class="my-6 mx-8 flex items-center justify-between">
	<!-- Release Filters -->
	<h3 class="md:visible text-gray-800 text-2xl">Releases</h3>
	<div class="flex items-center justify-end space-x-8">
		<div>
			<HeadFilterToggle />
		</div>
		<div>
			<ProjectFilter />
		</div>
		<div>
			<ReleaseSort />
		</div>
	</div>
</div>
<div class="mt-8 max-w-6xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
	{#if $store.fetching}
		<Loading text="Loading Releases" height={10} width={10} />
	{:else if $store.error}
		<h1 class="mt-5 md:mt-20 text-xl font-bold text-center leading-tight text-red-700">
			Error fetching data: {$store.error.toString()}
		</h1>
	{:else if $store.data.releases.length === 0}
		<div class="mt-12 flex justify-center">
			<div
				class="relative block w-96 border-2 border-gray-300 border-dashed rounded-lg p-12 text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
			>
				<Icon name="search" class="mx-auto h-12 w-12 text-gray-400" />
				<span class="mt-2 block text-sm font-medium text-gray-900"> No releases found. </span>
			</div>
		</div>
	{:else}
		<div class="bg-white shadow overflow-hidden sm:rounded-md">
			<ul class="divide-y divide-gray-200">
				{#each $store.data.releases as release}
					<ReleaseRow {release} bind:selectedRelease />
				{/each}
			</ul>
		</div>
	{/if}
</div>

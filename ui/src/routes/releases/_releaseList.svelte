<script lang="ts">
	import Icon from '$lib/Icon.svelte';

	import { operationStore, query } from '@urql/svelte';
	import HeadFilterToggle from './_headFilterToggle.svelte';
	import ReleaseSort from './_releaseSort.svelte';
	import SelectWithSearch from './_selectWithSearch.svelte';
	import type {
		Project,
		Project_Json,
		Release,
		Release_Conn,
		Release_Relay
	} from '$schema/schema_gen';
	import Loading from '$lib/Loading.svelte';
	import ReleaseRow from './_releaseRow.svelte';
	import { writable } from 'svelte/store';
	import { filterHeadOnly, filterProjects, sortReleases } from './stores';
	import type { ProjectSelectMap } from './stores';
	import { base } from '$app/paths';

	export let selectedRelease: Release;

	const releaseQuery = (): string => {
		let releaseCursor = null;
		let queryPredicates: string[] = [];

		if (filterHeadOnly.isEnabled()) {
			queryPredicates.push('has_head_of: true');
		}
		const selProjects = filterProjects.selectedProjects();
		if (selProjects.length > 0) {
			queryPredicates.push(
				`has_commit_with: {has_repo_with: {has_project_with: {name_in: ["${selProjects.join(
					`","`
				)}"]}}}`
			);
		}

		let queryArgs: string[] = [
			'first: 20',
			`after: ${releaseCursor ? '"' + releaseCursor + '"' : 'null'}`
		];
		if (sortReleases.sortByField() !== null) {
			queryArgs.push(
				`order_by: {field: ${sortReleases.sortByField().toLowerCase()}, direction: DESC}`
			);
		}
		if (queryPredicates.length > 0) {
			queryArgs.push(`where: { ${queryPredicates.join(', ')} }`);
		}
		return `
    {
	
    release_connection(${queryArgs.join(', ')}) {
	    totalCount
	    pageInfo {
	      hasNextPage
	      startCursor
	      endCursor
	    }
	    edges {
	      cursor
	      node {
          id
          name
	        version
          violations {
            type
            severity
          }
	        commit {
	          hash
	          branch
	          tag
	          time
			  repo {
				  name
			  }
	        }
	      }
	    }
	  }
	}`;
	};

	const store = operationStore<Release_Relay>(releaseQuery());

	// Subscribe to changes from the different filter stores we have
	filterHeadOnly.subscribe((b) => {
		store.query = releaseQuery();
	});

	filterProjects.subscribe((p) => {
		store.query = releaseQuery();
	});

	sortReleases.subscribe((s) => {
		store.query = releaseQuery();
	});

	query(store);
</script>

<div class="my-6 mx-8 flex items-center justify-between">
	<!-- Release Filters -->
	<h3 class="text-gray-800 text-2xl">Releases</h3>
	<div class="flex items-center justify-end space-x-8">
		<div>
			<HeadFilterToggle />
		</div>
		<div>
			<SelectWithSearch />
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
	{:else if $store.data.release_connection.totalCount === 0}
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
				{#each $store.data.release_connection.edges as release}
					<ReleaseRow release={release.node} bind:selectedRelease />
				{/each}
			</ul>
		</div>
	{/if}
</div>

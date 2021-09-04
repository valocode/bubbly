<script lang="ts">
	import dayjs from 'dayjs';
	import { operationStore, query } from '@urql/svelte';
	import type { Release, Release_Json } from '$schema/schema_gen';
	import Loading from '$lib/Loading.svelte';
	import Icon from '$lib/Icon.svelte';
	import ReleaseLog from './_releaseLog.svelte';
	import ViolationList from './_violationList.svelte';
	import ReleaseStatus from './_releaseStatus.svelte';

	export let release: Release;
	let detailRelease: Release = null;
	const store = operationStore<Release_Json>(null);
	$: if (release !== null) {
		const releaseQuery = `
{
  release(where:{id: ${release.id}}) {
    name
    version
    commit {
      hash
      branch
      tag
      time
      repo {
        name
		project {
		  name
		}
      }
    }
    head_of {
      id
    }
    log(order_by: {field: time, direction: ASC}) {
      type
      time
      artifact {
        name
        type
      }
      code_scan {
        tool
      }
      test_run {
        tool
      }
    }
    violations {
      message
      type
      severity
    }
  }
}
		`;
		store.query = releaseQuery;
		query(store);
		store.subscribe((value) => {
			if (!value.fetching && !value.error) {
				if (value.data.release !== null) {
					detailRelease = value.data.release[0];
				}
			}
		});
	}
</script>

<div class="h-full bg-white p-8 overflow-y-auto">
	<h2 class="sr-only">Details</h2>
	{#if release === null}
		<div class="mt-12 text-center">
			<h3 class="mt-2 text-lg font-medium text-gray-500">No release selected</h3>
		</div>
	{:else if $store.fetching}
		<Loading text="Loading Releases" height={10} width={10} />
	{:else if $store.error}
		<h1 class="mt-5 md:mt-20 text-xl font-bold text-center leading-tight text-red-700">
			Error fetching data: {$store.error.toString()}
		</h1>
	{:else if detailRelease !== null}
		<div class="space-y-5">
			<div class="flex items-center space-x-2">
				<ReleaseStatus release={detailRelease} />
			</div>
			<div class="flex items-center space-x-2">
				<Icon name="hashtag" class="h-5 w-5 text-gray-400" />
				<span class="text-gray-900 text-sm font-medium truncate"
					>Version: {detailRelease.version}</span
				>
			</div>
			<div class="flex items-center space-x-2">
				<Icon name="folder" class="h-5 w-5 text-gray-400" />
				<span class="text-gray-900 text-sm font-medium"
					>Repository: {detailRelease.commit.repo.name}</span
				>
			</div>
			<div class="flex items-center space-x-2">
				<Icon name="collection" class="h-5 w-5 text-gray-400" />
				<span class="text-gray-900 text-sm font-medium"
					>Project: {detailRelease.commit.repo.project.name}</span
				>
			</div>
			<div class="flex items-center space-x-2">
				<!-- Heroicon name: solid/calendar -->
				<Icon name="calendar" class="h-5 w-5 text-gray-400" />
				<span class="text-gray-900 text-sm font-medium"
					><time datetime="2020-12-02"
						>{dayjs(detailRelease.commit.time).format('YYYY-MM-DD hh:mm:ss')}</time
					></span
				>
			</div>
		</div>
		<div class="mt-6 border-t border-gray-200 py-6 space-y-8">
			<div>
				<h2 class="text-sm font-medium text-gray-500">Violations</h2>
				<ViolationList violations={detailRelease.violations} />
			</div>
		</div>
		<div class="mt-6 border-t border-b border-gray-200 py-6 space-y-8">
			<div>
				<h2 class="text-sm mb-8 font-medium text-gray-500">Release Log</h2>
				<ReleaseLog release={detailRelease} />
			</div>
		</div>
	{/if}
</div>

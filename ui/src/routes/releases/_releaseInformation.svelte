<script lang="ts">
	import dayjs from 'dayjs';
	import { httpStore } from '$stores/http/store';
	import Loading from '$lib/Loading.svelte';
	import Icon from '$lib/Icon.svelte';
	import ReleaseLog from './_releaseLog.svelte';
	import ViolationList from './_violationList.svelte';
	import ReleaseStatus from './_releaseStatus.svelte';

	export let release;
	let detailRelease = null;
	const store = httpStore(null);
	$: if (release !== null) {
		store.get(`/releases/${release.release.id}`).catch((err) => {});
		store.subscribe((value) => {
			if (!value.fetching && !value.error) {
				if (value.data.release !== null) {
					detailRelease = value.data.release;
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
		{console.log(detailRelease)}
		<div class="space-y-5">
			<div class="flex items-center space-x-2">
				<ReleaseStatus release={detailRelease} />
			</div>
			<div class="flex items-center space-x-2">
				<Icon name="hashtag" class="h-5 w-5 text-gray-400" />
				<span class="text-gray-900 text-sm font-medium truncate"
					>Version: {detailRelease.release.version}</span
				>
			</div>
			<div class="flex items-center space-x-2">
				<Icon name="folder" class="h-5 w-5 text-gray-400" />
				<span class="text-gray-900 text-sm font-medium">Repository: {detailRelease.repo.name}</span>
			</div>
			<div class="flex items-center space-x-2">
				<Icon name="collection" class="h-5 w-5 text-gray-400" />
				<span class="text-gray-900 text-sm font-medium">Project: {detailRelease.project.name}</span>
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

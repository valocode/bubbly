<script lang="ts">
	import dayjs from 'dayjs';

	import type { GitCommit, Release, ReleasePolicyViolation } from '$schema/schema_gen';
	import Icon from '$lib/Icon.svelte';
	import ReleaseStatus from './_releaseStatus.svelte';

	export let release: Release;
	export let selectedRelease: Release;
	const commit: GitCommit = release.commit;
	const time: Date = commit.time;
</script>

<!-- This example requires Tailwind CSS v2.0+ -->
<li>
	<a
		href="javascript:void(null);"
		class="block hover:bg-gray-50"
		on:click={() => (selectedRelease = release)}
	>
		<div class="flex items-center px-4 py-4 sm:px-6">
			<div class="min-w-0 flex-1 flex items-center">
				<div class="min-w-0 px-4 grid grid-cols-2 gap-4 flex-1">
					<div>
						<p class="text-sm font-medium text-bubbly-blue truncate">
							{release.version}
						</p>
						<div class="flex items-center space-x-2 mt-2">
							<!-- <Icon name="folder" class="h-4 w-4 text-gray-500" /> -->
							<p class="flex items-center text-sm text-gray-500">
								<span class="truncate">Repo: {release.commit.repo.name}</span>
							</p>
						</div>
						{#if time && commit}
							<p class="mt-2 flex items-center text-sm text-gray-500">
								<span class="truncate"
									>Created: {dayjs(release.commit.time).format('YYYY-MM-DD hh:mm:ss')}</span
								>
							</p>
						{/if}
					</div>
					<div class="">
						<div>
							<div class="flex items-center space-x-3">
								<ReleaseStatus {release} />
							</div>
							<p class="mt-2 flex items-center text-sm text-gray-800">
								<span class="truncate">{release.violations.length} policy violations</span>
							</p>
						</div>
					</div>
				</div>
			</div>
			<div>
				<Icon name="chevron-right" class="h-5 w-5 text-gray-400" />
			</div>
		</div>
	</a>
</li>

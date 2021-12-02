<script lang="ts">
	import dayjs from 'dayjs';
	import { ReleaseEntryType } from '$schema/schema_gen';
	import Icon from '$lib/Icon.svelte';
	export let release;

	let previousEvent: Date = null;

	const eventTimeDiff = (event: Date): string => {
		if (previousEvent === null) {
			previousEvent = release.commit.time;
		}

		const diffMs = new Date(event).getTime() - new Date(previousEvent).getTime();
		previousEvent = event;
		return `${Math.floor(diffMs / 60000)}m`;
	};
</script>

<div class="flow-root">
	<ul class="-mb-8">
		<li>
			<div class="relative pb-10">
				<span class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" aria-hidden="true" />
				<div class="relative flex space-x-3">
					<div>
						<span
							class="h-8 w-8 rounded-full bg-green-500 flex items-center justify-center ring-8 ring-white"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-5 w-5 text-white"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14"
								/>
							</svg>
						</span>
					</div>
					<div class="min-w-0 flex-1 pt-1.5 flex justify-between space-x-4">
						<div>
							<p class="text-sm text-gray-500">Commit created</p>
						</div>
						<div class="text-right text-sm whitespace-nowrap text-gray-500">
							{dayjs(release.commit.time).format('YYYY-MM-DD hh:mm:ss')}
						</div>
					</div>
				</div>
			</div>
		</li>

		{#each release.entries as entry, i}
			<li>
				<div class="relative pb-10">
					<!-- If the last entry, don't create a line -->
					{#if i + 1 !== release.entries.length}
						<span
							class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"
							aria-hidden="true"
						/>
					{/if}
					<div class="relative flex space-x-3 items-center">
						<div>
							{#if entry.type == ReleaseEntryType.artifact}
								<span
									class="h-8 w-8 rounded-full bg-bubbly-dark flex items-center justify-center ring-8 ring-white"
								>
									<Icon name="cube" class="h-5 w-5 text-white" />
								</span>
							{:else}
								<span
									class="h-8 w-8 rounded-full bg-blue-500 flex items-center justify-center ring-8 ring-white"
								>
									<Icon name="cog" class="h-5 w-5 text-white" />
								</span>
							{/if}
						</div>
						<div class="min-w-0 flex-1 pt-1.5 flex justify-between items-center space-x-4">
							<div>
								<p class="text-sm text-gray-500">
									Event: <span class="font-medium text-gray-900">{entry.type}</span>
								</p>
								<p class="text-sm text-gray-500">
									After: {eventTimeDiff(entry.time)}
								</p>
							</div>
							<div class="text-right text-sm whitespace-nowrap text-gray-500">
								{dayjs(entry.time).format('YYYY-MM-DD hh:mm:ss')}
							</div>
						</div>
					</div>
				</div>
			</li>
		{/each}
	</ul>
</div>

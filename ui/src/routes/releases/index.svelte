<script lang="ts">
	import { writable } from 'svelte/store';

	import type { Release } from '$schema/schema_gen';
	import ReleaseInformation from './_releaseInformation.svelte';
	import ReleaseList from './_releaseList.svelte';

	const selectedRelease = writable<Release>(null);
</script>

<div class="h-screen flex overflow-hidden bg-white">
	<div class="flex flex-col min-w-0 flex-1 overflow-hidden">
		<div class="flex-1 relative z-0 flex overflow-hidden">
			<main class="flex-1 relative z-0 overflow-y-auto focus:outline-none">
				<!-- Start main area-->
				<div class="absolute inset-0 py-6 px-4 sm:px-6 lg:px-8">
					<div class="h-full rounded-lg">
						<ReleaseList bind:selectedRelease={$selectedRelease} />
					</div>
				</div>
				<!-- End main area -->
			</main>
			<aside
				class="hidden relative xl:flex xl:flex-col flex-shrink-0 w-4/12 border-r border-gray-200"
			>
				<!-- Start secondary column (hidden on smaller screens) -->
				<div class="absolute inset-0 py-6 px-4 sm:px-6 lg:px-8">
					<div class="h-full border-2 border-gray-200 border-dashed rounded-lg">
						<ReleaseInformation release={$selectedRelease} />
					</div>
				</div>
				<!-- End secondary column -->
			</aside>
		</div>
	</div>
</div>

<script lang="ts">
	import { filterHeadOnly } from './stores';

	let enabled: boolean = false;

	filterHeadOnly.subscribe((b) => {
		enabled = b;
	});

	$: colorClass = (): string => {
		if (enabled) {
			return 'bg-bubbly-blue';
		}
		return 'bg-gray-200';
	};
	$: translateClass = (): string => {
		if (enabled) {
			return 'translate-x-5';
		}
		return 'translate-x-0';
	};
</script>

<div class="flex items-center">
	<span class="mr-3" id="release-filter-head-only">
		<span class="text-sm font-medium text-gray-900">Head only</span>
	</span>
	<!-- Enabled: "bg-indigo-600", Not Enabled: "bg-gray-200" -->
	<button
		type="button"
		on:click={filterHeadOnly.toggled}
		class="{colorClass()} relative inline-flex flex-shrink-0 h-6 w-11 border-2 border-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
		role="switch"
		aria-checked="false"
		aria-labelledby="release-filter-head-only"
	>
		<!-- Enabled: "translate-x-5", Not Enabled: "translate-x-0" -->
		<span
			aria-hidden="true"
			class="{translateClass()} pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow transform ring-0 transition ease-in-out duration-200"
		/>
	</button>
</div>

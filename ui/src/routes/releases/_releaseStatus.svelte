<script lang="ts">
	import { ReleasePolicyViolationSeverity } from '$schema/schema_gen';
	import type { Release } from '$schema/schema_gen';
	import { worstSeverity } from './release';

	export let release: Release;

	const color = () => {
		switch (severity) {
			case null:
				return 'green';
			case ReleasePolicyViolationSeverity.blocking:
				return 'red';
			case ReleasePolicyViolationSeverity.warning:
				return 'orange';
			case ReleasePolicyViolationSeverity.suggestion:
				return 'yellow';
			default:
				return 'gray';
		}
	};

	const text = () => {
		switch (severity) {
			case null:
				return 'Release Ready';
			case ReleasePolicyViolationSeverity.blocking:
				return 'Release Blocked';
			case ReleasePolicyViolationSeverity.warning:
				return 'Ready with warnings';
			case ReleasePolicyViolationSeverity.suggestion:
				return 'Ready with suggestions';
			default:
				return 'Unknown';
		}
	};

	const severity = worstSeverity(release);
	const releaseColor = color();
</script>

<span
	class="h-4 w-4 bg-{releaseColor}-100 rounded-full flex items-center justify-center"
	aria-hidden="true"
>
	<span class="h-2 w-2 bg-{releaseColor}-400 rounded-full" />
</span>
<span class="block">
	<h2 class="text-sm font-medium">
		{text()}
	</h2>
</span>

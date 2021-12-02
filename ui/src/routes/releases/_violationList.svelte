<script lang="ts">
	import { ReleasePolicyViolationSeverity } from '$schema/schema_gen';

	export let violations = [];

	const violationColor = (severity: ReleasePolicyViolationSeverity): string => {
		switch (severity) {
			case ReleasePolicyViolationSeverity.blocking:
				return 'red';
			case ReleasePolicyViolationSeverity.warning:
				return 'orange';
			case ReleasePolicyViolationSeverity.suggestion:
				return 'yellow';
		}
	};
</script>

<div>
	<div class="flow-root mt-6">
		<ul role="list" class="-my-5 divide-y divide-gray-200">
			{#each violations as violation}
				<li class="py-5">
					<div class="relative focus-within:ring-2 focus-within:ring-indigo-500">
						<div class="flex items-center justify-between">
							<span
								class="inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium bg-{violationColor(
									violation.severity
								)}-100 text-{violationColor(violation.severity)}-800"
							>
								{violation.severity}
							</span>
							<p class="text-sm font-semibold text-gray-800">
								{violation.type}
							</p>
						</div>
						<p class="mt-1 text-sm text-gray-600 line-clamp-2">
							{violation.message}
						</p>
					</div>
				</li>
			{/each}
		</ul>
	</div>
</div>

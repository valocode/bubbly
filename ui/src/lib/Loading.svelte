<script>
	export let text;
	export let height = 5;
	export let width = 5;
</script>

<!-- TODO: use bg-bubbly-blue for fill -->
<div class="flex justify-center items-center py-10">
	<svg class="animate-spin h-{height} w-{width} mr-3 ..." viewBox="0 0 24 24">
		<circle class="opacity-25" cx="12" cy="12" r="10" fill="white" stroke-width="4" />
		<path
			class="opacity-75"
			fill="#5408e8"
			d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
		/>
	</svg>
	{text}
</div>

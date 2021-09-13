<script context="module">
	export function load({ error, status }) {
		return {
			props: {
				error: error,
				status: status
			}
		};
	}
</script>

<script lang="ts">
	import { base } from '$app/paths';

	export let error;
	export let status: number;
</script>

<div class="py-16 bg-gray-50 overflow-hidden lg:py-24">
	<div class="relative max-w-xl mx-auto px-4 sm:px-6 lg:px-8 lg:max-w-7xl">
		<div class="relative">
			<h1 class="my-8 text-center text-6xl font-semibold tracking-tight text-gray-900 sm:text-9xl">
				{status}
			</h1>
			<h2
				class="text-center text-3xl leading-8 font-semibold tracking-tight text-gray-900 sm:text-4xl"
			>
				{error.message}
			</h2>
			<a href="{base}/">
				<p
					class="mt-16 w-52 py-3 mx-auto border-dashed border-4 border-purple hover:border-darkblue text-center text-xl text-purple hover:text-darkblue"
				>
					Take me home
				</p>
			</a>
		</div>
	</div>
</div>

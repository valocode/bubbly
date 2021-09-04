<!-- 
	Icon allows you to use pre-styled icons present in $lib/icons.ts

	TODO: This abstraction was nice when solving the Navigation and Block components, 
	but elsewhere icon SVGs may require certain fill/stroke/other customizations
	that may overcomplicate this component.
-->
<script lang="ts">
	export let name: string;

	import { icons } from '$lib/icons';

	let displayIcon = icons.find((e) => e.name === name);
	let isFill = name.startsWith('solid-');
</script>

<svg
	class={$$props.class}
	xmlns="http://www.w3.org/2000/svg"
	fill={isFill ? 'currentColor' : 'none'}
	stroke={isFill ? 'none' : 'currentColor'}
	aria-hidden="true"
	viewBox="0 0 {displayIcon.box} {displayIcon.box}">{@html displayIcon.svg}</svg
>
<span class="sr-only">{name}</span>

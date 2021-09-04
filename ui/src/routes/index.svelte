<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';

	// For now just redirect to releases
	goto(`${base}/releases`);
</script>

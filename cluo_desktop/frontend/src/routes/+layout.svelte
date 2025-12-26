<script lang="ts">
	import { onMount } from 'svelte';
	import { theme } from '$lib/stores/theme';
	import '../app.css';

	// Initialize theme on mount
	onMount(() => {
		theme.set($theme);
	});
</script>

<slot />

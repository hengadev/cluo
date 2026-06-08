<script lang="ts">
	// /users has been merged into /people.
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";

	onMount(() => {
		goto("/people", { replaceState: true });
	});
</script>

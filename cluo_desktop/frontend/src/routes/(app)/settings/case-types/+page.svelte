<script lang="ts">
	// /settings/case-types has moved to the "Types d'affaire" tab on /cases.
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";

	onMount(() => {
		goto("/cases", { replaceState: true });
	});
</script>

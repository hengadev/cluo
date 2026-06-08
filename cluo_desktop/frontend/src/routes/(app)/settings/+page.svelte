<script lang="ts">
	// /settings has been removed. Its tabs migrated:
	// - "Types d'affaire" → /cases (tab on the Affaires page)
	// - "Utilisateurs" → /people (unified personnes list)
	// Redirect to /cases as a safe fallback.
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";

	onMount(() => {
		goto("/cases", { replaceState: true });
	});
</script>

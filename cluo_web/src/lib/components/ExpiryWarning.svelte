<script lang="ts">
	import { daysUntil } from '$lib/utils/days-until';

	interface Props {
		expiresAt: string;
	}

	let { expiresAt }: Props = $props();

	const remaining = $derived(daysUntil(expiresAt));
	const show = $derived(remaining >= 0 && remaining <= 7);
	const label = $derived(remaining === 1 ? 'jour' : 'jours');
</script>

{#if show}
	<div class="rounded-card border border-yellow-500/50 bg-yellow-50 text-yellow-800 px-4 py-3 text-sm">
		Votre lien expire dans {remaining} {label}
	</div>
{/if}

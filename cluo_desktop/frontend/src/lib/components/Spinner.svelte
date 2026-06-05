<script lang="ts">
	type Size = "sm" | "md" | "lg";

	interface Props {
		size?: Size;
	}

	let { size = "md" }: Props = $props();

	const sizes: Record<Size, string> = {
		sm: "h-4 w-4 border-2",
		md: "h-6 w-6 border-2",
		lg: "h-8 w-8 border-[3px]",
	};
</script>

<div
	class="inline-block animate-spin rounded-full border-muted border-t-foreground {sizes[size]}"
	role="status"
	aria-label="Chargement"
>
	<span class="sr-only">Chargement…</span>
</div>

<style>
	@media (prefers-reduced-motion: reduce) {
		div {
			opacity: 0.5;
			animation: pulse 1.5s ease-in-out infinite;
		}
	}

	@keyframes pulse {
		0%,
		100% {
			opacity: 0.4;
		}
		50% {
			opacity: 1;
		}
	}
</style>

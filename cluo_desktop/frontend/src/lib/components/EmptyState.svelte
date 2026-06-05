<script lang="ts">
	import type { Snippet, Component } from "svelte";

	type Props = {
		icon: Component;
		message: string;
		children?: Snippet;
	};

	let { icon: Icon, message, children }: Props = $props();
</script>

<div
	class="border border-dashed border-border rounded-lg p-6 bg-muted/20 flex flex-col items-center gap-3"
>
	<Icon class="w-8 h-8 text-muted-foreground" />
	<p class="text-sm text-muted-foreground text-center">{message}</p>
	{#if children}
		{@render children()}
	{/if}
</div>

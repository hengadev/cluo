<script lang="ts">
	import {
		MessageSquareMore,
		AudioLines,
	} from '@lucide/svelte';
	import AIChatPanel from './_aiChatPanel.svelte';
	import AiAudioPanel from './_aiAudioPanel.svelte';

	interface Props {
		width: number;
		minWidth: number;
		maxWidth: number;
	}

	let { width, minWidth, maxWidth }: Props = $props();
	let selected = $state(0);

	type AIButton = {
		icon: typeof import('@lucide/svelte').Icon;
		title: string;
	};

	const aiButtons: AIButton[] = [
		{ icon: MessageSquareMore, title: 'Chat' },
		{ icon: AudioLines, title: 'Audio' },
	];
</script>

<div
	class="grid grid-rows-[auto_1fr] border-l border-border-card overflow-x-hidden bg-background"
	style="width: {width}px; min-width: {minWidth}px; max-width: {maxWidth}px;"
>
	<div
		class="flex justify-center border-b border-border-card text-center"
		style="gap: {width <= minWidth ? '0.5rem' : '1.5rem'};"
	>
		{#each aiButtons as item, index}
			{@render button(index, item)}
		{/each}
	</div>

	<!-- Content Area -->
	<div class="overflow-hidden">
		{#if selected === 0}
			<AIChatPanel {width} />
		{:else if selected === 1}
			<AiAudioPanel {width} />
		{/if}
	</div>
</div>

{#snippet button(index: number, item: AIButton)}
	{@const Icon = item.icon}
	{@const isSelected = index === selected}
	{@const isCompact = width <= minWidth}
	<button
		class="flex items-center gap-2 transition-interactive relative"
		class:opacity-60={!isSelected}
		style="padding: {isCompact ? '0.75rem 0.5rem' : '1rem'};"
		onclick={() => (selected = index)}
	>
		{#if !isCompact}
			<Icon
				size={20}
				strokeWidth={1.5}
				class={isSelected ? 'text-primary' : 'text-muted-foreground'}
			/>
		{/if}
		<p class={isSelected ? 'text-foreground font-medium' : 'text-muted-foreground'}>
			{item.title}
		</p>
		{#if isSelected}
			<div class="absolute bottom-0 left-1/2 -translate-x-1/2 w-8 h-0.5 bg-primary rounded-full"></div>
		{/if}
	</button>
{/snippet}

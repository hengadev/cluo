<script lang="ts">
	import {
		MessageSquareMore,
		Lightbulb,
		PencilLine,
		AudioLines,
	} from '@lucide/svelte';
	import AIChatPanel from './_aiChatPanel.svelte';

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
		{ icon: Lightbulb, title: 'Idées' },
		{ icon: PencilLine, title: 'Relecture' },
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
			<div
				class="flex items-center justify-center h-full px-4 py-8 text-center"
				style="height: calc(100vh - 200px);"
			>
				<div>
					<Lightbulb class="w-12 h-12 mx-auto mb-4 text-muted-foreground/50" />
					<p class="text-muted-foreground text-sm">Idées - Coming Soon</p>
				</div>
			</div>
		{:else if selected === 2}
			<div
				class="flex items-center justify-center h-full px-4 py-8 text-center"
				style="height: calc(100vh - 200px);"
			>
				<div>
					<PencilLine class="w-12 h-12 mx-auto mb-4 text-muted-foreground/50" />
					<p class="text-muted-foreground text-sm">Relecture - Coming Soon</p>
				</div>
			</div>
		{:else if selected === 3}
			<div
				class="flex items-center justify-center h-full px-4 py-8 text-center"
				style="height: calc(100vh - 200px);"
			>
				<div>
					<AudioLines class="w-12 h-12 mx-auto mb-4 text-muted-foreground/50" />
					<p class="text-muted-foreground text-sm">Audio - Coming Soon</p>
				</div>
			</div>
		{/if}
	</div>
</div>

{#snippet button(index: number, item: AIButton)}
	{@const Icon = item.icon}
	{@const isSelected = index === selected}
	{@const isCompact = width <= minWidth}
	<button
		class="flex items-center gap-2 transition-all relative"
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

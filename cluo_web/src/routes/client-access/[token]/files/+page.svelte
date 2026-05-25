<script lang="ts">
	import type { PageData } from './$types';
	import { cn } from '$lib/utils/design-system';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	type TabId = 'documents' | 'rapport' | 'medias';

	let activeTab: TabId = $state('documents');

	const tabs: { id: TabId; label: string }[] = $derived(
		[
			{ id: 'documents' as const, label: 'Documents' },
			{ id: 'rapport' as const, label: 'Rapport' },
			...(data.hasMedia ? [{ id: 'medias' as const, label: 'Médias' }] : [])
		]
	);

</script>

<svelte:head>
	<title>{data.caseData.title} — Dossier client</title>
</svelte:head>

<div class="min-h-screen bg-background">
	<!-- Header -->
	<header class="border-b border-border-card px-6 py-4">
		<div class="max-w-4xl mx-auto flex items-center justify-between">
			<div class="flex items-center gap-4">
				<div class="w-8 h-8 border border-dark-300 flex items-center justify-center">
					<span class="font-mono text-xxs tracking-widest text-foreground">A·I·R</span>
				</div>
				<div>
					<h1 class="font-serif text-foreground text-lg leading-tight">{data.caseData.title}</h1>
					{#if data.caseData.externalReference}
						<p class="text-foreground-alt text-xs">{data.caseData.externalReference}</p>
					{/if}
				</div>
			</div>
			<a
				href="/client-access/{data.token}"
				class="text-foreground-alt text-sm hover:text-foreground transition-colors"
			>
				← Retour
			</a>
		</div>
	</header>

	<!-- Tab navigation -->
	<div class="border-b border-border-card">
		<nav class="max-w-4xl mx-auto px-6 flex gap-0">
			{#each tabs as tab (tab.id)}
				<button
					class={cn(
						'px-5 py-3 text-sm font-medium transition-colors border-b-2 -mb-px',
						activeTab === tab.id
							? 'border-foreground text-foreground'
							: 'border-transparent text-foreground-alt hover:text-foreground hover:border-dark-200'
					)}
					onclick={() => (activeTab = tab.id)}
				>
					{tab.label}
				</button>
			{/each}
		</nav>
	</div>

	<!-- Tab content -->
	<main class="max-w-4xl mx-auto px-6 py-10">
		{#if activeTab === 'documents'}
			<div class="flex flex-col items-center justify-center py-16 text-center">
				<p class="text-foreground-alt text-sm">Les documents du dossier apparaîtront ici.</p>
				<p class="text-foreground-alt/60 text-xs mt-2">Contenu à venir.</p>
			</div>
		{:else if activeTab === 'rapport'}
			{#if data.rapportHtml}
				<div class="rapport-content">
					{@html data.rapportHtml}
				</div>
			{:else if data.rapportError}
				<div class="flex flex-col items-center justify-center py-16 text-center">
					<p class="text-foreground-alt text-sm">Le rapport n'a pas pu être chargé. Veuillez réessayer.</p>
				</div>
			{:else}
				<div class="flex flex-col items-center justify-center py-16 text-center">
					<p class="text-foreground-alt text-sm">Aucun rapport n'est disponible pour ce dossier.</p>
				</div>
			{/if}
		{:else if activeTab === 'medias'}
			<div class="flex flex-col items-center justify-center py-16 text-center">
				<p class="text-foreground-alt text-sm">Les médias du dossier apparaîtront ici.</p>
				<p class="text-foreground-alt/60 text-xs mt-2">Contenu à venir.</p>
			</div>
		{/if}
	</main>
</div>

<style>
	.rapport-content {
		@apply text-foreground leading-relaxed;
	}

	.rapport-content :global(h1) {
		@apply font-serif text-2xl text-foreground mt-8 mb-4;
	}

	.rapport-content :global(h2) {
		@apply font-serif text-xl text-foreground mt-6 mb-3;
	}

	.rapport-content :global(h3) {
		@apply font-serif text-lg text-foreground mt-4 mb-2;
	}

	.rapport-content :global(p) {
		@apply text-foreground text-sm leading-relaxed mb-4;
	}

	.rapport-content :global(strong) {
		@apply font-semibold text-foreground;
	}

	.rapport-content :global(em) {
		@apply italic;
	}

	.rapport-content :global(u) {
		@apply underline;
	}

	.rapport-content :global(blockquote) {
		@apply border-l-2 border-dark-200 pl-4 my-4 italic text-foreground-alt;
	}

	.rapport-content :global(ul) {
		@apply list-disc list-outside ml-5 mb-4 space-y-1;
	}

	.rapport-content :global(ol) {
		@apply list-decimal list-outside ml-5 mb-4 space-y-1;
	}

	.rapport-content :global(li) {
		@apply text-sm text-foreground;
	}

	.rapport-content :global(li p) {
		@apply mb-1;
	}
</style>

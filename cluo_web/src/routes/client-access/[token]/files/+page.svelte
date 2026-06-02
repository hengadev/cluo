<script lang="ts">
	import type { PageData } from './$types';
	import { cn, buttonVariants } from '$lib/utils/design-system';
	import type { DocumentSummaryResponse, MediaResponse } from '$lib/server/client-access';
	import ExpiryWarning from '$lib/components/ExpiryWarning.svelte';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	type TabId = 'documents' | 'rapport' | 'medias';

	let activeTab: TabId = $state((data.activeTab as TabId) ?? 'documents');

	/** Track which document sections are expanded */
	let expandedDocs = $state<Set<string>>(new Set());

	/** Track which document previews are open */
	let previewOpen = $state<Set<string>>(new Set());

	function toggleDoc(id: string) {
		const next = new Set(expandedDocs);
		if (next.has(id)) {
			next.delete(id);
		} else {
			next.add(id);
		}
		expandedDocs = next;
	}

	function togglePreview(id: string) {
		const next = new Set(previewOpen);
		if (next.has(id)) {
			next.delete(id);
		} else {
			next.add(id);
		}
		previewOpen = next;
	}

	/** Canonical display order for document types */
	const TYPE_ORDER: Record<string, number> = {
		estimate: 0,
		mandate: 1,
		contract: 2,
		invoice: 3
	};

	const TYPE_LABELS: Record<string, string> = {
		estimate: 'Devis',
		mandate: 'Mandat',
		contract: 'Contrat',
		invoice: 'Facture'
	};

	const STATUS_LABELS: Record<string, string> = {
		sent: 'Envoyé',
		signed: 'Signé',
		active: 'Actif',
		archived: 'Archivé'
	};

	/** Sort documents in canonical order, filtering out types not in the known set */
	const sortedDocuments: DocumentSummaryResponse[] = $derived(
		[...data.documents]
			.filter((d) => d.type in TYPE_ORDER)
			.sort((a, b) => TYPE_ORDER[a.type] - TYPE_ORDER[b.type])
	);

	// ---- Media helpers ----

	const images: MediaResponse[] = $derived(data.media.filter((m) => m.type === 'image'));
	const videos: MediaResponse[] = $derived(data.media.filter((m) => m.type === 'video'));

	let lightboxOpen = $state(false);
	let lightboxIndex = $state(0);

	function openLightbox(index: number) {
		lightboxIndex = index;
		lightboxOpen = true;
	}

	function closeLightbox() {
		lightboxOpen = false;
	}

	function prevImage() {
		lightboxIndex = (lightboxIndex - 1 + images.length) % images.length;
	}

	function nextImage() {
		lightboxIndex = (lightboxIndex + 1) % images.length;
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return `${bytes} o`;
		if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} Ko`;
		return `${(bytes / (1024 * 1024)).toFixed(1)} Mo`;
	}

	const formatDate = (iso: string) =>
		new Intl.DateTimeFormat('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' }).format(new Date(iso));

	const tabs: { id: TabId; label: string }[] = $derived(
		[
			{ id: 'documents' as const, label: 'Documents' },
			{ id: 'rapport' as const, label: 'Rapport' },
			...(data.hasMedia ? [{ id: 'medias' as const, label: 'Médias' }] : [])
		]
	);

	function handleKeydown(e: KeyboardEvent) {
		if (!lightboxOpen) return;
		if (e.key === 'Escape') closeLightbox();
		if (e.key === 'ArrowLeft') prevImage();
		if (e.key === 'ArrowRight') nextImage();
	}
</script>

<svelte:window onkeydown={handleKeydown} />

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
					<ExpiryWarning expiresAt={data.caseData.tokenExpiresAt} />
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
			{#if data.documentsError}
				<div class="flex flex-col items-center justify-center py-16 text-center">
					<p class="text-foreground-alt text-sm">Les documents n'ont pas pu être chargés. Veuillez réessayer.</p>
				</div>
			{:else if sortedDocuments.length === 0}
				<div class="flex flex-col items-center justify-center py-16 text-center">
					<p class="text-foreground-alt text-sm">Aucun document n'est disponible pour ce dossier.</p>
				</div>
			{:else}
				<div class="flex flex-col gap-3">
					{#each sortedDocuments as doc (doc.id)}
						<div class="rounded-card border border-border-card bg-background shadow-card overflow-hidden">
							<!-- Header row -->
							<button
								class="w-full text-left flex items-center justify-between p-5"
								onclick={() => toggleDoc(doc.id)}
							>
								<div class="flex items-center gap-3">
									<h3 class="font-serif text-foreground text-lg">{TYPE_LABELS[doc.type]}</h3>
									<span class="inline-flex items-center rounded-9px px-2.5 py-0.5 text-xs font-medium bg-muted text-foreground-alt">
										{STATUS_LABELS[doc.status] ?? doc.status}
									</span>
								</div>
								<svg
									class={cn(
										'w-4 h-4 text-foreground-alt transition-transform duration-200',
										expandedDocs.has(doc.id) && 'rotate-180'
									)}
									fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"
								>
									<path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7" />
								</svg>
							</button>

							<!-- Expanded detail -->
							{#if expandedDocs.has(doc.id)}
								<div class="border-t border-border-card px-5 py-4">
									<div class="grid grid-cols-2 gap-4 mb-4">
										<div>
											<p class="text-xs tracking-widest uppercase text-foreground-alt mb-1">Référence</p>
											<p class="text-foreground text-sm">{doc.document_ref}</p>
										</div>
										<div>
											<p class="text-xs tracking-widest uppercase text-foreground-alt mb-1">Date de création</p>
											<p class="text-foreground text-sm">{formatDate(doc.created_at)}</p>
										</div>
										<div>
											<p class="text-xs tracking-widest uppercase text-foreground-alt mb-1">Dernière mise à jour</p>
											<p class="text-foreground text-sm">{formatDate(doc.updated_at)}</p>
										</div>
									</div>
								<div class="flex gap-2">
									<button
										class={buttonVariants({ variant: 'outline', size: 'sm' })}
										onclick={() => togglePreview(doc.id)}
									>
										Prévisualiser
									</button>
									<a
										href="/client-access/{data.token}/documents/{doc.type}/pdf"
										class={buttonVariants({ variant: 'outline', size: 'sm' })}
										download
									>
										<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
											<path stroke-linecap="round" stroke-linejoin="round" d="M12 10v6m0 0l-3-3m3 3l3-3M3 17v3a2 2 0 002 2h14a2 2 0 002-2v-3" />
										</svg>
										Télécharger le PDF
									</a>
								</div>

								{#if previewOpen.has(doc.id)}
									<div class="mt-4">
										<iframe
											src="/client-access/{data.token}/documents/{doc.type}/pdf"
											class="w-full border border-border-card rounded"
											style="height: 600px;"
											title="Prévisualisation PDF"
										></iframe>
									</div>
								{/if}
							</div>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		{:else if activeTab === 'rapport'}
			{#if data.rapportHtml}
				<div class="mb-6 flex flex-col gap-2">
					<a
						href="/client-access/{data.token}/rapport-pdf"
						class={cn(buttonVariants({ variant: 'outline', size: 'sm' }))}
					>
						<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M12 10v6m0 0l-3-3m3 3l3-3M3 17v3a2 2 0 002 2h14a2 2 0 002-2v-3" />
						</svg>
						Télécharger le rapport (PDF)
					</a>
					{#if data.pdfError}
						<p class="text-sm text-red-600">
							{data.pdfError === 'not_found'
								? 'Aucun rapport disponible pour ce dossier.'
								: 'Le téléchargement du PDF a échoué. Veuillez réessayer.'}
						</p>
					{/if}
				</div>
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
			{#if data.media.length === 0}
				<div class="flex flex-col items-center justify-center py-16 text-center">
					<p class="text-foreground-alt text-sm">Aucun média n'est disponible pour ce dossier.</p>
				</div>
			{:else}
				<!-- "Télécharger tous les médias" -->
				<div class="mb-8">
					<a
						href="/client-access/{data.token}/media-archive"
						download
						class={cn(
							buttonVariants({ variant: 'outline', size: 'sm' })
						)}
					>
						<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2M7 10l5 5 5-5M12 15V3" />
						</svg>
						Télécharger tous les médias
					</a>
				</div>

				<!-- Photos -->
				{#if images.length > 0}
					<section class="mb-10">
						<h2 class="font-serif text-foreground text-lg mb-4">Photos</h2>
						<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
							{#each images as img, i (img.id)}
								<div class="group relative rounded-card overflow-hidden border border-border-card bg-muted aspect-square">
									<button
										class="w-full h-full cursor-pointer"
										onclick={() => openLightbox(i)}
									>
										<img
											src={img.url}
											alt={img.caption || img.fileName}
											class="w-full h-full object-cover transition-transform duration-200 group-hover:scale-105"
											loading="lazy"
										/>
									</button>
									<!-- Overlay with download -->
									<div class="absolute bottom-0 inset-x-0 bg-gradient-to-t from-black/60 to-transparent p-2 opacity-0 group-hover:opacity-100 transition-opacity">
										<a
											href="/client-access/{data.token}/media/{img.id}/download"
											download
											class={cn(
												buttonVariants({ variant: 'ghost', size: 'sm' }),
												'text-white hover:bg-white/20 h-7 px-2 text-xs'
											)}
										>
											<svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
												<path stroke-linecap="round" stroke-linejoin="round" d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2M7 10l5 5 5-5M12 15V3" />
											</svg>
											{formatFileSize(img.fileSize)}
										</a>
									</div>
								</div>
							{/each}
						</div>
					</section>
				{/if}

				<!-- Videos -->
				{#if videos.length > 0}
					<section class="mb-10">
						<h2 class="font-serif text-foreground text-lg mb-4">Vidéos</h2>
						<div class="flex flex-col gap-4">
							{#each videos as vid (vid.id)}
								<div class="rounded-card border border-border-card bg-background overflow-hidden">
									<div class="aspect-video bg-black">
										<video
											controls
											preload="metadata"
											class="w-full h-full"
										>
											<source src={vid.url} type={vid.mimeType} />
											Votre navigateur ne supporte pas la lecture vidéo.
										</video>
									</div>
									<div class="px-4 py-3 flex items-center justify-between">
										<div class="min-w-0">
											{#if vid.caption}
												<p class="text-foreground text-sm truncate">{vid.caption}</p>
											{/if}
											<p class="text-foreground-alt text-xs">{vid.fileName} · {formatFileSize(vid.fileSize)}</p>
										</div>
										<a
											href="/client-access/{data.token}/media/{vid.id}/download"
											download
											class={cn(buttonVariants({ variant: 'outline', size: 'sm' }), 'shrink-0')}
										>
											<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
												<path stroke-linecap="round" stroke-linejoin="round" d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2M7 10l5 5 5-5M12 15V3" />
											</svg>
											Télécharger
										</a>
									</div>
								</div>
							{/each}
						</div>
					</section>
				{/if}


			{/if}
		{/if}
	</main>

	<!-- Lightbox overlay -->
	{#if lightboxOpen && images.length > 0}
		{@const current = images[lightboxIndex]}
		<div
			class="fixed inset-0 z-50 bg-black/90 flex items-center justify-center"
			role="dialog"
			aria-modal="true"
			aria-label="Visualiseur de photos"
			onclick={closeLightbox}
		>
			<!-- Close -->
			<button
				class="absolute top-4 right-4 text-white/80 hover:text-white transition-colors z-10"
				onclick={closeLightbox}
				aria-label="Fermer"
			>
				<svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
					<path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
				</svg>
			</button>

			<!-- Prev -->
			{#if images.length > 1}
				<button
					class="absolute left-4 top-1/2 -translate-y-1/2 text-white/80 hover:text-white transition-colors z-10"
					onclick={prevImage}
					aria-label="Précédent"
				>
					<svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
					</svg>
				</button>

				<!-- Next -->
				<button
					class="absolute right-4 top-1/2 -translate-y-1/2 text-white/80 hover:text-white transition-colors z-10"
					onclick={nextImage}
					aria-label="Suivant"
				>
					<svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
					</svg>
				</button>
			{/if}

			<!-- Image -->
			<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
			<div class="max-w-[90vw] max-h-[85vh] flex flex-col items-center" onclick={(e) => e.stopPropagation()}>
				<img
					src={current.url}
					alt={current.caption || current.fileName}
					class="max-w-full max-h-[75vh] object-contain select-none"
				/>
				<!-- Caption + download -->
				<div class="mt-3 flex items-center gap-4">
					{#if current.caption}
						<p class="text-white/80 text-sm">{current.caption}</p>
					{/if}
					<a
						href="/client-access/{data.token}/media/{current.id}/download"
						download
						class="text-white/80 hover:text-white text-sm underline-offset-2 hover:underline transition-colors"
					>
						Télécharger ({formatFileSize(current.fileSize)})
					</a>
				</div>
				<!-- Counter -->
				{#if images.length > 1}
					<p class="text-white/50 text-xs mt-2">{lightboxIndex + 1} / {images.length}</p>
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	@reference "../../../app.css";

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

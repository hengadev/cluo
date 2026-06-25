<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { currentCase } from '$lib/stores/case';
	import {
		fetchCaseMedia,
		getTranscriptionByMediaFile,
		getAnalysisByTranscriptionId,
	} from '$lib/services/api';
	import type { TranscriptAnalysis } from '$lib/types/entities';
	import {
		AudioLines,
		ChevronDown,
		ChevronRight,
		Loader2,
		FileText,
		Tag,
	} from '@lucide/svelte';

	interface Props {
		width: number;
	}
	let { width }: Props = $props();

	interface AnalysisItem {
		analysis: TranscriptAnalysis;
		fileName: string;
		expanded: boolean;
	}

	let caseId = $derived.by(() => $currentCase.id || $page.params.id);
	let items = $state<AnalysisItem[]>([]);
	let loading = $state(false);
	let loaded = $state(false);

	async function load() {
		if (!caseId) return;
		loading = true;
		try {
			const mediaResp = await fetchCaseMedia(caseId, 'audio');
			const results: AnalysisItem[] = [];
			for (const m of mediaResp.media) {
				try {
					const tResp = await getTranscriptionByMediaFile(m.id);
					const transcription = tResp.transcriptions[0];
					if (!transcription) continue;
					const analysis = await getAnalysisByTranscriptionId(transcription.id);
					if (analysis) {
						results.push({ analysis, fileName: m.fileName, expanded: false });
					}
				} catch {
					// Skip media without a usable transcription/analysis
				}
			}
			items = results;
		} finally {
			loading = false;
			loaded = true;
		}
	}

	onMount(load);

	function parseTopics(analysis: TranscriptAnalysis): string[] {
		try {
			const parsed = JSON.parse(analysis.topics);
			return Array.isArray(parsed) ? parsed.filter((t) => typeof t === 'string') : [];
		} catch {
			return [];
		}
	}

	function sentimentLabel(sentiment: string): string {
		switch (sentiment) {
			case 'positive':
				return 'Positif';
			case 'neutral':
				return 'Neutre';
			case 'negative':
				return 'Négatif';
			case 'mixed':
				return 'Mixte';
			default:
				return sentiment;
		}
	}

	function sentimentClass(sentiment: string): string {
		switch (sentiment) {
			case 'positive':
				return 'bg-success/15 text-success';
			case 'negative':
				return 'bg-destructive/15 text-destructive';
			case 'mixed':
				return 'bg-primary/15 text-primary';
			default:
				return 'bg-muted text-muted-foreground';
		}
	}
</script>

<div
	class="flex flex-col h-full"
	style="height: calc(100vh - 200px);"
>
	{#if loading}
		<div class="flex items-center justify-center h-full">
			<Loader2 size={20} class="animate-spin text-muted-foreground" />
		</div>
	{:else if items.length === 0}
		<div class="flex flex-col items-center justify-center h-full px-4 py-8 text-center gap-3">
			<AudioLines class="w-10 h-10 text-muted-foreground/50" />
			<p class="text-muted-foreground text-sm">
				Aucune analyse disponible. Lancez une transcription puis analysez-la depuis la page Enregistrements.
			</p>
		</div>
	{:else}
		<div class="flex-1 min-h-0 overflow-y-auto p-3 space-y-3">
			{#each items as item (item.analysis.id)}
				<div class="rounded-card border border-border-card bg-background p-3">
					<!-- Header: file name + sentiment -->
					<div class="flex items-start justify-between gap-2">
						<span class="inline-flex items-center gap-1.5 text-sm font-medium text-foreground truncate">
							<FileText size={14} class="flex-shrink-0 text-muted-foreground" />
							<span class="truncate">{item.fileName}</span>
						</span>
						<span
							class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium flex-shrink-0 {sentimentClass(item.analysis.sentiment)}"
						>
							{sentimentLabel(item.analysis.sentiment)}
						</span>
					</div>

					<!-- Summary (always visible) -->
					{#if item.analysis.summary}
						<p class="mt-2 text-sm text-foreground leading-relaxed">{item.analysis.summary}</p>
					{/if}

					<!-- Topics (always visible) -->
					{#if parseTopics(item.analysis).length > 0}
						<div class="mt-2 flex flex-wrap gap-1.5">
							{#each parseTopics(item.analysis) as topic}
								<span
									class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs bg-muted text-muted-foreground border border-border-card"
								>
									<Tag size={10} />
									{topic}
								</span>
							{/each}
						</div>
					{/if}

					<!-- Key findings (collapsed by default) -->
					{#if item.analysis.keyFindings || item.analysis.suggestedActions}
						<button
							type="button"
							onclick={() => (item.expanded = !item.expanded)}
							class="mt-2 inline-flex items-center gap-1 text-xs text-foreground hover:text-muted-foreground cursor-pointer"
						>
							{#if item.expanded}
								<ChevronDown size={12} />
							{:else}
								<ChevronRight size={12} />
							{/if}
							{item.expanded ? 'Masquer les détails' : 'Voir les détails'}
						</button>
						{#if item.expanded}
							<div class="mt-2 flex flex-col gap-3">
								{#if item.analysis.keyFindings}
									<div>
										<p class="text-xs font-medium text-muted-foreground mb-1">Points clés</p>
										<p class="text-sm text-foreground whitespace-pre-wrap leading-relaxed">
											{item.analysis.keyFindings}
										</p>
									</div>
								{/if}
								{#if item.analysis.suggestedActions}
									<div>
										<p class="text-xs font-medium text-muted-foreground mb-1">Actions suggérées</p>
										<p class="text-sm text-foreground whitespace-pre-wrap leading-relaxed">
											{item.analysis.suggestedActions}
										</p>
									</div>
								{/if}
							</div>
						{/if}
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>

<script lang="ts">
	import { ChevronLeft, Download } from "@lucide/svelte";
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { getAnalysis } from "$lib/api";
	import type { AnalysisResult } from "$lib/types/recording";

	let analysis = $state<AnalysisResult | null>(null);
	let isLoading = $state(true);
	let error = $state<string | null>(null);
	let recordingId = $state("");

	onMount(() => {
		const pathParts = window.location.pathname.split("/");
		recordingId = pathParts[pathParts.length - 2];
		loadAnalysis();
	});

	async function loadAnalysis() {
		try {
			isLoading = true;
			error = null;
			analysis = await getAnalysis(recordingId);
		} catch (err) {
			error = err instanceof Error ? err.message : "Failed to load analysis";
		} finally {
			isLoading = false;
		}
	}

	function parsedTopics(): string[] {
		if (!analysis?.topics) return [];
		try {
			return JSON.parse(analysis.topics) as string[];
		} catch {
			return analysis.topics ? [analysis.topics] : [];
		}
	}

	function sentimentColor(sentiment: string): string {
		switch (sentiment?.toLowerCase()) {
			case "positive": return "bg-green-100 text-green-800";
			case "negative": return "bg-red-100 text-red-800";
			default: return "bg-dark-100 text-dark-700";
		}
	}

	function handleExport() {
		if (!analysis) return;

		const topics = parsedTopics();
		let text = `AI Analysis\n${"=".repeat(40)}\n\n`;
		text += `Sentiment: ${analysis.sentiment}\n\n`;
		if (topics.length) text += `Topics: ${topics.join(", ")}\n\n`;
		text += `Summary\n${"-".repeat(20)}\n${analysis.summary}\n\n`;
		text += `Key Findings\n${"-".repeat(20)}\n${analysis.keyFindings}\n\n`;
		text += `Suggested Actions\n${"-".repeat(20)}\n${analysis.suggestedActions}\n`;

		const blob = new Blob([text], { type: "text/plain" });
		const url = URL.createObjectURL(blob);
		const a = document.createElement("a");
		a.href = url;
		a.download = `analysis-${recordingId}.txt`;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}

	function goBack() {
		if (history.length > 0) history.back();
		else goto(`/recording/${recordingId}`);
	}
</script>

<div class="min-h-screen flex flex-col gap-6 pb-24 mt-8 px-4">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<button onclick={goBack} class="text-dark-700 hover:text-dark-900">
			<ChevronLeft />
		</button>
		<p class="text-dark-900 font-bold text-lg">AI Analysis</p>
		<div class="w-6"></div>
	</div>

	{#if isLoading}
		<div class="flex items-center justify-center p-12">
			<p class="text-dark-600">Loading analysis...</p>
		</div>
	{:else if error}
		<div class="flex flex-col items-center gap-4 p-8 bg-red-50 rounded-2xl">
			<p class="text-red-700 font-semibold">Error</p>
			<p class="text-red-600 text-sm text-center">{error}</p>
			<button
				onclick={loadAnalysis}
				class="px-6 py-3 bg-red-600 hover:bg-red-500 text-white rounded-xl transition-colors font-medium"
			>
				Retry
			</button>
		</div>
	{:else if analysis}
		<!-- Sentiment + topics -->
		<div class="flex flex-wrap items-center gap-2">
			<span class="px-3 py-1 rounded-full text-sm font-medium {sentimentColor(analysis.sentiment)} capitalize">
				{analysis.sentiment || "neutral"}
			</span>
			{#each parsedTopics() as topic}
				<span class="px-3 py-1 rounded-full text-sm bg-dark-50 text-dark-700 border border-dark-100">
					{topic}
				</span>
			{/each}
		</div>

		<!-- Summary -->
		{#if analysis.summary}
			<div class="flex flex-col gap-2 p-4 bg-dark-50 rounded-2xl">
				<p class="text-dark-800 font-semibold text-sm uppercase tracking-wide">Summary</p>
				<p class="text-dark-700 text-sm leading-relaxed">{analysis.summary}</p>
			</div>
		{/if}

		<!-- Key findings -->
		{#if analysis.keyFindings}
			<div class="flex flex-col gap-2 p-4 border border-dark-100 rounded-2xl">
				<p class="text-dark-800 font-semibold text-sm uppercase tracking-wide">Key Findings</p>
				<p class="text-dark-700 text-sm leading-relaxed whitespace-pre-line">{analysis.keyFindings}</p>
			</div>
		{/if}

		<!-- Suggested actions -->
		{#if analysis.suggestedActions}
			<div class="flex flex-col gap-2 p-4 border border-dark-100 rounded-2xl">
				<p class="text-dark-800 font-semibold text-sm uppercase tracking-wide">Suggested Actions</p>
				<p class="text-dark-700 text-sm leading-relaxed whitespace-pre-line">{analysis.suggestedActions}</p>
			</div>
		{/if}

		<!-- Actions -->
		<div class="flex flex-col gap-3 mt-2">
			<button
				onclick={handleExport}
				class="flex items-center justify-center gap-2 px-6 py-4 bg-dark-700 hover:bg-dark-600 text-foreground rounded-xl transition-colors font-semibold"
			>
				<Download size={18} />
				<span>Export Analysis</span>
			</button>
			<button
				onclick={() => goto(`/recording/${recordingId}`)}
				class="flex items-center justify-center gap-2 px-6 py-4 bg-dark-100 hover:bg-dark-200 text-dark-700 rounded-xl transition-colors font-semibold"
			>
				Back to Recording
			</button>
		</div>

		<div class="flex items-center justify-center p-4 bg-dark-50 rounded-2xl mt-2">
			<p class="text-dark-600 text-sm text-center">
				Analysis is processed on private infrastructure
			</p>
		</div>
	{/if}
</div>

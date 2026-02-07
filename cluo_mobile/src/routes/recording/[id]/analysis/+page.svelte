<script lang="ts">
	import { ChevronLeft, Sparkles, Download, Share2 } from "@lucide/svelte";
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import AnalysisSection from "$lib/components/AnalysisSection.svelte";
	import { getAnalysis } from "$lib/api";
	import type { AnalysisResult, Suggestion, SuggestionCategory } from "$lib/types/recording";

	let analysis = $state<AnalysisResult | null>(null);
	let isLoading = $state(true);
	let error = $state<string | null>(null);

	// Get recording ID from URL
	let recordingId = $state("");

	// Group suggestions by category
	let groupedSuggestions = $derived(() => {
		if (!analysis) return null;
		const groups: Record<SuggestionCategory, Suggestion[]> = {
			observations: [],
			statements: [],
			actions: [],
			unclear: [],
		};
		analysis.suggestions.forEach((s) => {
			groups[s.category].push(s);
		});
		return groups;
	});

	// Count total selected suggestions
	let totalSelected = $derived(() => {
		if (!groupedSuggestions) return 0;
		return (
			groupedSuggestions.observations.filter((s) => s.selected).length +
			groupedSuggestions.statements.filter((s) => s.selected).length +
			groupedSuggestions.actions.filter((s) => s.selected).length +
			groupedSuggestions.unclear.filter((s) => s.selected).length
		);
	});

	onMount(() => {
		// Extract recording ID from URL path
		const pathParts = window.location.pathname.split("/");
		recordingId = pathParts[pathParts.length - 2]; // Get the ID before "analysis"

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

	function getSelectedSuggestions(): Suggestion[] {
		if (!analysis) return [];
		return analysis.suggestions.filter((s) => s.selected);
	}

	function handleExport() {
		const selected = getSelectedSuggestions();
		if (selected.length === 0) {
			alert("No suggestions selected");
			return;
		}

		// Group by category for export
		const grouped: Record<string, string[]> = {
			Observations: [],
			Statements: [],
			Actions: [],
			Unclear: [],
		};

		selected.forEach((s) => {
			const category =
				s.category.charAt(0).toUpperCase() +
				s.category.slice(1);
			if (category === "Unclear") {
				grouped.Unclear.push(s.text);
			} else if (grouped[category]) {
				grouped[category].push(s.text);
			}
		});

		// Create formatted text
		let exportText = `Analysis Results - Recording ${recordingId}\n\n`;
		Object.entries(grouped).forEach(([category, items]) => {
			if (items.length > 0) {
				exportText += `## ${category}\n\n`;
				items.forEach((item, i) => {
					exportText += `${i + 1}. ${item}\n`;
				});
				exportText += "\n";
			}
		});

		// Download as text file
		const blob = new Blob([exportText], { type: "text/plain" });
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
	{:else if analysis && groupedSuggestions}
		<!-- Summary -->
		<div class="flex items-center justify-between p-4 bg-dark-50 rounded-xl">
			<div class="flex items-center gap-3">
				<Sparkles class="text-accent" size={20} />
				<div>
					<p class="text-dark-900 font-semibold text-sm">
						{analysis.suggestions.length} suggestions
					</p>
					<p class="text-dark-600 text-xs">
						{totalSelected} selected
					</p>
				</div>
			</div>
		</div>

		<!-- Analysis Sections -->
		<div class="flex flex-col gap-4">
			<AnalysisSection
				category="observations"
				suggestions={groupedSuggestions.observations}
			/>
			<AnalysisSection
				category="statements"
				suggestions={groupedSuggestions.statements}
			/>
			<AnalysisSection
				category="actions"
				suggestions={groupedSuggestions.actions}
			/>
			<AnalysisSection
				category="unclear"
				suggestions={groupedSuggestions.unclear}
			/>
		</div>

		<!-- Action Buttons -->
		<div class="flex flex-col gap-3">
			<button
				onclick={handleExport}
				disabled={totalSelected === 0}
				class="flex items-center justify-center gap-2 px-6 py-4 bg-dark-700 hover:bg-dark-600 disabled:bg-dark-400 text-foreground rounded-xl transition-colors font-semibold"
			>
				<Download size={18} />
				<span>Export Selected ({totalSelected})</span>
			</button>

			<button
				onclick={() => goto(`/recording/${recordingId}`)}
				class="flex items-center justify-center gap-2 px-6 py-4 bg-dark-100 hover:bg-dark-200 text-dark-700 rounded-xl transition-colors font-semibold"
			>
				<span>Back to Recording</span>
			</button>
		</div>

		<!-- Privacy Notice -->
		<div class="flex items-center justify-center p-4 bg-dark-50 rounded-2xl mt-2">
			<p class="text-dark-600 text-sm text-center">
				Analysis is processed on private infrastructure
			</p>
		</div>
	{/if}
</div>

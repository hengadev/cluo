<script lang="ts">
	import { ChevronDown, ChevronUp, CheckSquare, Square } from "@lucide/svelte";
	import type { Suggestion, SuggestionCategory } from "$lib/types/recording";

	interface Props {
		category: SuggestionCategory;
		suggestions: Suggestion[];
		readonly?: boolean;
		class?: string;
	}

	let {
		category,
		suggestions,
		readonly = false,
		class: className = "",
	}: Props = $props();

	let isExpanded = $state(true);
	let selectedCount = $derived(
		suggestions.filter((s) => s.selected).length,
	);

	// Category display names
	const categoryNames: Record<SuggestionCategory, string> = {
		observations: "Observations",
		statements: "Déclarations",
		actions: "Actions",
		unclear: "Flou / Nécessite une confirmation",
	};

	// Category icons (simple emoji for now, could use Lucide icons)
	const categoryIcons: Record<SuggestionCategory, string> = {
		observations: "👁️",
		statements: "💬",
		actions: "✅",
		unclear: "❓",
	};

	function toggleSelectAll() {
		if (readonly) return;
		const allSelected = suggestions.every((s) => s.selected);
		suggestions.forEach((s) => {
			s.selected = !allSelected;
		});
	}

	function toggleSuggestion(index: number) {
		if (readonly) return;
		suggestions[index].selected = !suggestions[index].selected;
	}
</script>

<div class="analysis-section flex flex-col gap-3 {className}">
	<!-- Section Header -->
	<button
		onclick={() => (isExpanded = !isExpanded)}
		class="flex items-center justify-between p-4 bg-dark-50 hover:bg-dark-100 rounded-xl transition-colors"
	>
		<div class="flex items-center gap-3">
			<span class="text-2xl">{categoryIcons[category]}</span>
			<div class="flex flex-col items-start">
				<p class="text-dark-900 font-semibold">{categoryNames[category]}</p>
				<p class="text-dark-600 text-xs">{selectedCount} / {suggestions.length} sélectionné{suggestions.length > 1 ? 's' : ''}</p>
			</div>
		</div>
		{#if isExpanded}
			<ChevronUp class="text-dark-600" size={20} />
		{:else}
			<ChevronDown class="text-dark-600" size={20} />
		{/if}
	</button>

	{#if isExpanded}
		<!-- Suggestions List -->
		<div class="flex flex-col gap-2">
			{#if suggestions.length === 0}
				<p class="text-dark-500 text-sm italic p-4">Aucune suggestion dans cette catégorie.</p>
			{:else}
				<!-- Select All Button -->
				{#if !readonly}
					<button
						onclick={toggleSelectAll}
						class="flex items-center gap-2 p-2 text-dark-600 hover:text-dark-900 text-sm transition-colors"
					>
						{#if suggestions.every((s) => s.selected)}
							<CheckSquare size={16} />
							<span>Tout désélectionner</span>
						{:else}
							<Square size={16} />
							<span>Tout sélectionner</span>
						{/if}
					</button>
				{/if}

				<!-- Individual Suggestions -->
				{#each suggestions as suggestion, index}
					<div
						class="flex items-start gap-3 p-3 bg-background border-1 border-dark-100 rounded-lg hover:bg-dark-50 transition-colors {suggestion.selected
							? 'border-accent'
							: ''}"
					>
						<button
							onclick={() => toggleSuggestion(index)}
							class="mt-0.5 text-dark-600 hover:text-dark-900 transition-colors"
							aria-label="Toggle selection"
						>
							{#if suggestion.selected}
								<CheckSquare size={18} class="text-accent" />
							{:else}
								<Square size={18} />
							{/if}
						</button>
						<p class="flex-1 text-dark-800 text-sm leading-relaxed">
							{suggestion.text}
						</p>
					</div>
				{/each}
			{/if}
		</div>
	{/if}
</div>

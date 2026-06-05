<script lang="ts">
	import { currentCase } from "$lib/stores/case";
	import { onMount } from "svelte";
	import { fetchAllCases } from "$lib/services/api";
	import type { Case, CaseStatus } from "$lib/types/entities";

	let cases: Case[] = [];
	let loading = true;
	let error: string | null = null;

	const STATUS_LABELS: Record<CaseStatus, string> = {
		draft: "Brouillon",
		in_progress: "En cours",
		ready: "Prêt",
		released: "Clôturé"
	};

	const STATUS_BADGE_CLASSES: Record<CaseStatus, string> = {
		draft: "bg-gray-100 text-gray-800",
		in_progress: "bg-blue-100 text-blue-800",
		ready: "bg-green-100 text-green-800",
		released: "bg-purple-100 text-purple-800"
	};

	onMount(async () => {
		try {
			cases = await fetchAllCases();
		} catch (e) {
			error = e instanceof Error ? e.message : "Erreur lors du chargement des dossiers";
		} finally {
			loading = false;
		}
	});

	function selectCase(caseId: string) {
		currentCase.setCase(caseId);
	}
</script>

<div class="p-8 flex flex-col gap-6">
    <div class="animate-fade-in">
        <h1 class="text-3xl font-bold">Tableau de bord</h1>
        <h2 class="text-xl font-semibold mt-2" style="animation-delay: 100ms;">
            Dossiers récents
        </h2>
    </div>
		<section>
			{#if loading}
				<p class="text-muted-foreground">Chargement...</p>
			{:else if error}
				<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
					{error}
				</div>
			{:else if cases.length === 0}
				<p class="text-muted-foreground">Aucun dossier disponible</p>
			{:else}
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each cases as caseItem, index}
						<button
							class="border border-border-card rounded-card p-4 bg-background hover:border-border-input-hover hover:shadow-md hover:-translate-y-1 transition-all duration-300 text-left animate-fade-in cursor-pointer"
							style="animation-delay: {200 + index * 100}ms;"
							onclick={() => selectCase(caseItem.id)}
						>
							<h3 class="font-semibold text-foreground">{caseItem.title}</h3>
							<p class="text-sm text-muted-foreground">{caseItem.id}</p>
							{#if caseItem.status}
								<span
									class="inline-block mt-2 px-2 py-1 text-xs rounded-full {STATUS_BADGE_CLASSES[caseItem.status] || 'bg-gray-100 text-gray-800'}"
								>
									{STATUS_LABELS[caseItem.status] || caseItem.status}
								</span>
							{/if}
						</button>
					{/each}
				</div>
			{/if}
		</section>
</div>

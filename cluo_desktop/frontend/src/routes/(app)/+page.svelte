<script lang="ts">
	import { currentCase } from "$lib/stores/case";
	import { onMount } from "svelte";
	import { fetchAllCases } from "$lib/services/api";
	import type { Case, CaseStatus } from "$lib/types/entities";
	import { caseStatusBadge } from "$lib/utils/badgeVariants";
	import Spinner from "$lib/components/Spinner.svelte";
	import EmptyState from "$lib/components/EmptyState.svelte";
	import { FolderOpen } from "@lucide/svelte";

	let cases: Case[] = [];
	let loading = true;
	let error: string | null = null;

	const STATUS_LABELS: Record<CaseStatus, string> = {
		in_progress: "En cours",
		ready: "Prêt",
		released: "Clôturé"
	};

	onMount(async () => {
		try {
			cases = (await fetchAllCases()).cases;
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

<div class="p-8 flex flex-col gap-6 flex-1">
    <div class="animate-fade-in">
        <h1 class="text-3xl font-bold">Tableau de bord</h1>
        <h2 class="text-xl font-semibold mt-2 animate-fade-in" style="animation-delay: 100ms;">
            Dossiers récents
        </h2>
    </div>
		<section class="flex flex-col flex-1">
			{#if loading}
				<div class="flex flex-1 items-center justify-center">
					<Spinner size="lg" />
				</div>
			{:else if error}
				<div class="alert-error">
					{error}
				</div>
			{:else if cases.length === 0}
				<EmptyState icon={FolderOpen} message="Aucun dossier disponible" />
			{:else}
				<div class="flex gap-4 mb-2">
					<div class="bg-muted rounded-card-sm px-4 py-3 flex flex-col">
						<span class="text-xs text-muted-foreground">Total</span>
						<span class="text-lg font-semibold text-foreground">{cases.length}</span>
					</div>
					<div class="bg-muted rounded-card-sm px-4 py-3 flex flex-col">
						<span class="text-xs text-muted-foreground">En cours</span>
						<span class="text-lg font-semibold text-foreground">{cases.filter(c => c.status === 'in_progress').length}</span>
					</div>
					<div class="bg-muted rounded-card-sm px-4 py-3 flex flex-col">
						<span class="text-xs text-muted-foreground">Prêts</span>
						<span class="text-lg font-semibold text-foreground">{cases.filter(c => c.status === 'ready').length}</span>
					</div>
				</div>
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each cases as caseItem, index}
						<button
							class="border border-border-card rounded-card p-4 bg-background hover:border-border-input-hover hover:shadow-md hover:-translate-y-1 transition-all duration-300 text-left animate-fade-in cursor-pointer"
							style="animation-delay: {200 + index * 100}ms;"
							onclick={() => selectCase(caseItem.id)}
						>
							<h3 class="font-semibold text-foreground">{caseItem.title}</h3>
							{#if caseItem.externalReference}
								<p class="text-sm text-muted-foreground">{caseItem.externalReference}</p>
							{/if}
							{#if caseItem.status}
								<span
									class="inline-block mt-2 px-2 py-1 text-xs rounded-full {caseStatusBadge(caseItem.status)}"
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

<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import { Plus, Search, FolderOpen } from "@lucide/svelte";
	import { Tabs } from "bits-ui";
	import {
		fetchAllCases,
	} from "$lib/services/api";
	import { currentCase } from "$lib/stores/case";
	import { caseStatusBadge } from "$lib/utils/badgeVariants";
	import type { Case, CaseStatus } from "$lib/types/entities";
	import Spinner from "$lib/components/Spinner.svelte";
	import EmptyState from "$lib/components/EmptyState.svelte";
	import NewCase from "$lib/custom/header/NewCase.svelte";
	import CaseTypesTab from "../settings/CaseTypesTab.svelte";

	const STATUS_LABELS: Record<CaseStatus, string> = {
		in_progress: "En cours",
		ready: "Prêt",
		released: "Clôturé",
	};

	type StatusFilter = "all" | CaseStatus;

	const STATUS_FILTERS: { value: StatusFilter; label: string }[] = [
		{ value: "all", label: "Tous" },
		{ value: "in_progress", label: "En cours" },
		{ value: "ready", label: "Prêt" },
		{ value: "released", label: "Clôturé" },
	];

	let cases: Case[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let statusFilter: StatusFilter = $state("all");
	let searchQuery = $state("");
	let newCaseOpen = $state(false);

	onMount(async () => {
		try {
			const response = await fetchAllCases();
			cases = response.cases;
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Erreur lors du chargement des dossiers";
		} finally {
			loading = false;
		}
	});

	let filteredCases = $derived.by(() => {
		let result = cases;

		// Status filter
		if (statusFilter !== "all") {
			result = result.filter((c) => c.status === statusFilter);
		}

		// Search filter (client-side: title or externalReference)
		const q = searchQuery.trim().toLowerCase();
		if (q) {
			result = result.filter(
				(c) =>
					c.title.toLowerCase().includes(q) ||
					(c.externalReference ?? "").toLowerCase().includes(q),
			);
		}

		return result;
	});

	let statusCounts = $derived.by(() => {
		const counts: Record<string, number> = {
			all: cases.length,
			in_progress: 0,
			ready: 0,
			released: 0,
		};
		for (const c of cases) {
			counts[c.status] = (counts[c.status] ?? 0) + 1;
		}
		return counts;
	});

	function selectCase(caseId: string) {
		currentCase.setCase(caseId);
		goto(`/cases/${caseId}`);
	}
</script>

<div class="page-content">
	<!-- Header row -->
	<div
		class="flex items-center justify-between animate-fade-in"
	>
		<h1 class="text-3xl font-bold">Affaires</h1>
		<button
			type="button"
			onclick={() => (newCaseOpen = true)}
			class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/90 inline-flex items-center justify-center gap-2 px-6 text-sm font-semibold active:scale-[0.98] cursor-pointer transition-all"
		>
			<Plus size={18} strokeWidth={2} />
			Nouvelle affaire
		</button>
	</div>

	<Tabs.Root value="cases">
		<div class="overflow-x-auto w-full border-b border-border">
			<Tabs.List
				class="inline-flex items-center gap-0 text-sm font-medium bg-transparent p-0 h-auto rounded-none"
			>
				<Tabs.Trigger
					value="cases"
					class="px-4 py-2 rounded-none bg-transparent border-b-2 mb-[-1px] data-[state=active]:border-foreground data-[state=active]:text-foreground data-[state=active]:shadow-none data-[state=inactive]:border-transparent data-[state=inactive]:text-muted-foreground data-[state=inactive]:hover:bg-transparent data-[state=inactive]:hover:text-foreground transition-colors cursor-pointer whitespace-nowrap"
				>
					Affaires
				</Tabs.Trigger>
				<Tabs.Trigger
					value="case-types"
					class="px-4 py-2 rounded-none bg-transparent border-b-2 mb-[-1px] data-[state=active]:border-foreground data-[state=active]:text-foreground data-[state=active]:shadow-none data-[state=inactive]:border-transparent data-[state=inactive]:text-muted-foreground data-[state=inactive]:hover:bg-transparent data-[state=inactive]:hover:text-foreground transition-colors cursor-pointer whitespace-nowrap"
				>
					Types d'affaire
				</Tabs.Trigger>
			</Tabs.List>
		</div>

		<!-- Cases tab -->
		<Tabs.Content value="cases" class="pt-6">
			{#if loading}
				<div class="flex flex-1 items-center justify-center py-16">
					<Spinner size="lg" />
				</div>
			{:else if error}
				<div class="alert-error">
					{error}
				</div>
			{:else}
				<div class="flex flex-wrap gap-4 pb-6">
					<!-- Search input -->
					<div class="relative max-w-[480px]">
						<Search
							class="text-muted-foreground absolute start-3 top-1/2 size-4 -translate-y-1/2 pointer-events-none"
						/>
						<input
							type="text"
							placeholder="Rechercher par titre ou référence..."
							bind:value={searchQuery}
							class="h-9 rounded-input border border-border-input bg-background pl-9 pr-4 text-sm placeholder:text-foreground-alt/40 hover:border-border-input-hover focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors w-full"
						/>
					</div>

					<!-- Status filter chips -->
					<div class="flex flex-wrap items-center gap-2">
						{#each STATUS_FILTERS as filter}
							<button
								type="button"
								onclick={() => (statusFilter = filter.value)}
								class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm font-medium transition-all cursor-pointer {statusFilter === filter.value
									? 'bg-foreground text-background'
									: 'bg-muted text-muted-foreground hover:bg-foreground/10'}"
							>
								{filter.label}
								<span
									class="text-xs {statusFilter === filter.value
										? 'text-background/70'
										: 'text-muted-foreground/60'}"
								>
									{statusCounts[filter.value] ?? 0}
								</span>
							</button>
						{/each}
					</div>
				</div>

				<!-- Case list -->
				{#if filteredCases.length === 0}
					<EmptyState
						icon={FolderOpen}
						message={searchQuery || statusFilter !== "all"
							? "Aucun dossier ne correspond aux filtres"
							: "Aucun dossier disponible"}
					/>
				{:else}
					<div
						class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
					>
						{#each filteredCases as caseItem, index}
							<button
								class="border border-border-card rounded-card p-4 bg-background hover:border-border-input-hover hover:shadow-card hover:-translate-y-1 transition-all duration-300 text-left animate-fade-in cursor-pointer"
								style="animation-delay: {200 + index * 50}ms;"
								onclick={() => selectCase(caseItem.id)}
							>
								<h3
									class="font-semibold text-foreground"
								>
									{caseItem.title}
								</h3>
								{#if caseItem.externalReference}
									<p
										class="text-sm text-muted-foreground mt-0.5"
									>
										{caseItem.externalReference}
									</p>
								{/if}
								{#if caseItem.status}
									<span
										class="inline-block mt-2 px-2 py-1 text-xs rounded-full {caseStatusBadge(caseItem.status)}"
									>
										{STATUS_LABELS[caseItem.status] ||
											caseItem.status}
									</span>
								{/if}
							</button>
						{/each}
					</div>
				{/if}
			{/if}
		</Tabs.Content>

		<!-- Case Types tab -->
		<Tabs.Content value="case-types" class="pt-6">
			<CaseTypesTab />
		</Tabs.Content>
	</Tabs.Root>
</div>

<NewCase bind:open={newCaseOpen} />

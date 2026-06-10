<script lang="ts">
	import { currentCase, recentCases } from "$lib/stores/case";
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { Briefcase, ArrowRight } from "@lucide/svelte";
	import { fetchAllCases } from "$lib/services/api";
	import type { Case, CaseStatus } from "$lib/types/entities";
	import { caseStatusBadge } from "$lib/utils/badgeVariants";
	import Spinner from "$lib/components/Spinner.svelte";

	const STATUS_LABELS: Record<CaseStatus, string> = {
		in_progress: "En cours",
		ready: "Prêt",
		released: "Clôturé",
	};

	let allCases: Case[] = $state([]);
	let loading = $state(true);

	onMount(async () => {
		// If there's a current case and no recent cases loaded yet, redirect to cases page
		if (!$currentCase.id && $recentCases.length === 0) {
			goto("/cases");
		}

		try {
			const response = await fetchAllCases();
			allCases = response.cases;
		} catch {
			// Stats are non-critical; show dashboard anyway
		} finally {
			loading = false;
		}
	});

	let stats = $derived({
		total: allCases.length,
		inProgress: allCases.filter((c) => c.status === "in_progress").length,
		ready: allCases.filter((c) => c.status === "ready").length,
	});

	function resumeCase(entry: { id: string; title: string; status: CaseStatus }) {
		currentCase.setCase(entry.id);
		goto(`/cases/${entry.id}`);
	}

	function goToCases() {
		goto("/cases");
	}
</script>

<div class="p-8 flex flex-col gap-6 flex-1">
	<div class="animate-fade-in">
		<h1 class="text-3xl font-bold">Tableau de bord</h1>
		<p class="text-sm text-muted-foreground mt-1">Vue d'ensemble de l'activité</p>
	</div>

	<!-- Stat strip -->
	{#if !loading}
		<div class="grid grid-cols-3 gap-4 max-w-xl animate-fade-in" style="animation-delay: 100ms;">
			<div class="border border-border-card rounded-card p-4 bg-background hover:shadow-popover transition-shadow">
				<p class="text-xs text-muted-foreground font-medium uppercase tracking-wider">Total</p>
				<p class="text-2xl font-bold text-foreground mt-1">{stats.total}</p>
			</div>
			<div class="border border-border-card rounded-card p-4 bg-background hover:shadow-popover transition-shadow">
				<p class="text-xs text-muted-foreground font-medium uppercase tracking-wider">En cours</p>
				<p class="text-2xl font-bold text-foreground mt-1">{stats.inProgress}</p>
			</div>
			<div class="border border-border-card rounded-card p-4 bg-background hover:shadow-popover transition-shadow">
				<p class="text-xs text-muted-foreground font-medium uppercase tracking-wider">Prêt</p>
				<p class="text-2xl font-bold text-foreground mt-1">{stats.ready}</p>
			</div>
		</div>
	{:else}
		<div class="grid grid-cols-3 gap-4 max-w-xl">
			{#each { length: 3 } as _}
				<div class="border border-border-card rounded-card p-4 bg-background h-[72px] flex items-center justify-center">
					<Spinner size="sm" />
				</div>
			{/each}
		</div>
	{/if}

	<div class="flex flex-col gap-6 max-w-2xl">
		<!-- Quick resume section -->
		{#if $recentCases.length > 0}
			<section>
				<h2 class="text-xl font-semibold mb-4 animate-fade-in" style="animation-delay: 200ms;">
					Reprendre
				</h2>
				<div class="flex flex-col gap-3">
					{#each $recentCases.slice(0, 3) as entry, index}
						<button
							class="flex items-center gap-4 border border-border-card rounded-card p-4 bg-background hover:border-border-input-hover hover:shadow-popover transition-all duration-300 text-left animate-fade-in cursor-pointer group"
							style="animation-delay: {300 + index * 100}ms;"
							onclick={() => resumeCase(entry)}
						>
							<div class="p-2 rounded-input bg-foreground/10">
								<Briefcase size={18} class="text-foreground" />
							</div>
							<div class="flex-1 min-w-0">
								<h3 class="font-semibold text-foreground truncate">{entry.title}</h3>
								<span
									class="inline-block mt-1 px-2 py-0.5 text-xs rounded-full {caseStatusBadge(entry.status)}"
								>
									{STATUS_LABELS[entry.status] || entry.status}
								</span>
							</div>
							<ArrowRight
								size={18}
								class="text-muted-foreground group-hover:text-foreground group-hover:translate-x-1 transition-all"
							/>
						</button>
					{/each}
				</div>
			</section>
		{/if}

		<!-- Link to full cases page -->
		<section class="animate-fade-in" style="animation-delay: 400ms;">
			<button
				type="button"
				onclick={goToCases}
				class="inline-flex items-center gap-2 text-sm font-medium text-foreground-alt hover:text-foreground transition-colors cursor-pointer"
			>
				<Briefcase size={16} />
				Voir toutes les affaires
				<ArrowRight size={14} />
			</button>
		</section>
	</div>
</div>

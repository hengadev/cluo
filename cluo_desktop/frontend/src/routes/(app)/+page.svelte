<script lang="ts">
	import { currentCase, recentCases } from "$lib/stores/case";
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { Briefcase, ArrowRight } from "@lucide/svelte";
	import type { CaseStatus } from "$lib/types/entities";
	import { caseStatusBadge } from "$lib/utils/badgeVariants";

	const STATUS_LABELS: Record<CaseStatus, string> = {
		in_progress: "En cours",
		ready: "Prêt",
		released: "Clôturé",
	};

	onMount(() => {
		// If there's a current case and no recent cases loaded yet, redirect to cases page
		if (!$currentCase.id && $recentCases.length === 0) {
			goto("/cases");
		}
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
	</div>

	<div class="flex flex-col gap-6 max-w-2xl">
		<!-- Quick resume section -->
		{#if $recentCases.length > 0}
			<section>
				<h2 class="text-xl font-semibold mb-4 animate-fade-in" style="animation-delay: 100ms;">
					Reprendre
				</h2>
				<div class="flex flex-col gap-3">
					{#each $recentCases.slice(0, 3) as entry, index}
						<button
							class="flex items-center gap-4 border border-border-card rounded-card p-4 bg-background hover:border-border-input-hover hover:shadow-md transition-all duration-300 text-left animate-fade-in cursor-pointer group"
							style="animation-delay: {200 + index * 100}ms;"
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
		<section class="animate-fade-in" style="animation-delay: 300ms;">
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

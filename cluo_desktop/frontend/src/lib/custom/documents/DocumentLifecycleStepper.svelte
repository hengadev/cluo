<!-- Hallmark · component: document-lifecycle-stepper · genre: modern-minimal · theme: project-system
     states: upcoming · current · completed · terminal
     contrast: pass (foreground/muted-foreground/success on background) -->
<script lang="ts">
	import { Check } from "@lucide/svelte";
	import { documentStatusBadge } from "$lib/utils/badgeVariants";
	import type { DocumentStatus } from "$lib/types/entities";

	interface Step {
		key: string;
		label: string;
	}

	interface Props {
		steps: Step[];
		status: string;
		statusLabel?: string;
		note?: string;
	}

	let { steps, status, statusLabel, note }: Props = $props();

	const currentIndex = $derived(steps.findIndex((s) => s.key === status));
	const isTerminal = $derived(currentIndex === -1);
</script>

{#if isTerminal}
	<div class="flex items-center gap-3">
		<span class="px-3 py-1.5 inline-flex text-sm leading-5 font-semibold rounded-full {documentStatusBadge(status as DocumentStatus)}">
			{statusLabel ?? status}
		</span>
		{#if note}
			<span class="text-sm text-muted-foreground">{note}</span>
		{/if}
	</div>
{:else}
	<div class="flex items-center gap-4 flex-wrap">
		<ol class="flex items-center">
			{#each steps as step, i (step.key)}
				<li class="flex items-center {i < steps.length - 1 ? 'flex-1' : ''}">
					<div class="flex items-center gap-2 shrink-0">
						<span
							class="flex items-center justify-center w-6 h-6 rounded-full text-xs font-semibold shrink-0 {i < currentIndex
								? 'bg-success text-success-foreground'
								: i === currentIndex
									? 'bg-foreground text-background'
									: 'bg-muted text-muted-foreground'}"
						>
							{#if i < currentIndex}
								<Check size={13} />
							{:else}
								{i + 1}
							{/if}
						</span>
						<span
							class="text-sm font-medium whitespace-nowrap {i === currentIndex
								? 'text-foreground'
								: i < currentIndex
									? 'text-foreground-alt'
									: 'text-muted-foreground'}"
						>
							{step.label}
						</span>
					</div>
					{#if i < steps.length - 1}
						<div class="flex-1 h-px mx-3 {i < currentIndex ? 'bg-success' : 'bg-border'}"></div>
					{/if}
				</li>
			{/each}
		</ol>
		{#if note}
			<span class="text-sm text-muted-foreground whitespace-nowrap">{note}</span>
		{/if}
	</div>
{/if}

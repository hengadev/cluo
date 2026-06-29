<script lang="ts">
	import { Dialog } from "bits-ui";
	import { X, Loader2, Check, X as XIcon, RotateCcw } from "@lucide/svelte";
	import type { AITextOperation } from "$lib/services/api";
	import { AI_OPERATION_LABELS } from "$lib/services/api";

	interface Props {
		open: boolean;
		operation: AITextOperation | null;
		originalText: string;
		suggestedText: string | null;
		isLoading?: boolean;
		error: string | null;
		onOpenChange: (open: boolean) => void;
		onAccept: () => void;
		onRetry: () => void;
	}

	let {
		open,
		operation,
		originalText,
		suggestedText,
		isLoading = false,
		error,
		onOpenChange,
		onAccept,
		onRetry
	}: Props = $props();

	// Get operation label
	function getOperationLabel() {
		if (!operation) return "";
		return AI_OPERATION_LABELS[operation].label;
	}
</script>

<Dialog.Root open={open} onOpenChange={onOpenChange}>
	<Dialog.Portal>
		<Dialog.Overlay
			class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
		/>
		<Dialog.Content
			class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-6 sm:max-w-[600px] md:w-full max-h-[90vh] overflow-y-auto"
		>
			<!-- Header -->
			<div class="flex items-start justify-between mb-4">
				<div>
					<Dialog.Title class="text-lg font-semibold tracking-tight flex items-center gap-2">
						{#if isLoading}
							<Loader2 class="w-5 h-5 animate-spin" />
						{:else if error}
							<XIcon class="w-5 h-5 text-destructive" />
						{:else}
							<Check class="w-5 h-5 text-primary" />
						{/if}
						{getOperationLabel()}
					</Dialog.Title>
					<Dialog.Description class="text-foreground-alt !mt-1 text-sm">
						{#if isLoading}
							Traitement en cours...
						{:else if error}
							Une erreur s'est produite
						{:else}
							Examinez la suggestion ci-dessous
						{/if}
					</Dialog.Description>
				</div>
				<Dialog.Close
					class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-4 top-4 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer hover:bg-muted"
				>
					<div>
						<X class="text-foreground size-5" />
						<span class="sr-only">Fermer</span>
					</div>
				</Dialog.Close>
			</div>

			<!-- Safety Message -->
			{#if isLoading}
				<div class="mb-4 p-3 bg-primary/10 border border-primary/20 rounded-lg flex items-start gap-2 text-sm">
					<Loader2 class="w-4 h-4 text-primary mt-0.5 flex-shrink-0 animate-spin" />
					<p class="text-primary">
						L'assistance IA est traitée sur une infrastructure privée.
					</p>
				</div>
			{/if}

			<!-- Loading State -->
			{#if isLoading}
				<div class="py-12 text-center">
					<Loader2 class="w-12 h-12 animate-spin mx-auto mb-4 text-muted-foreground" />
					<p class="text-muted-foreground">Génération de la suggestion...</p>
				</div>
			{/if}

			<!-- Error State -->
			{#if error}
				<div class="mb-4 p-4 bg-destructive/10 border border-destructive/20 rounded-lg">
					<p class="text-destructive text-sm">{error}</p>
				</div>
			{/if}

			<!-- Comparison View -->
			{#if suggestedText && !isLoading && !error}
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
					<!-- Original Text -->
					<div class="border border-border-card rounded-lg p-3">
						<h3 class="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-2">
							Original
						</h3>
						<p class="text-sm whitespace-pre-wrap">{originalText}</p>
					</div>

					<!-- Suggested Text -->
					<div class="border border-primary/30 rounded-lg p-3 bg-primary/5">
						<h3 class="text-xs font-semibold uppercase tracking-wide text-primary mb-2">
							Suggestion
						</h3>
						<p class="text-sm whitespace-pre-wrap">{suggestedText}</p>
					</div>
				</div>
			{/if}

			<!-- Actions -->
			<div class="flex justify-end gap-2 !mt-6">
				{#if isLoading}
					<Dialog.Close
						class="h-input rounded-input bg-transparent text-foreground hover:bg-muted focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] border cursor-pointer"
					>
						Annuler
					</Dialog.Close>
				{:else if error}
					<Dialog.Close
						class="h-input rounded-input bg-transparent text-foreground hover:bg-muted focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] border cursor-pointer"
					>
						Fermer
					</Dialog.Close>
					<button
						onclick={onRetry}
						class="h-input rounded-input bg-primary text-primary-foreground shadow-mini hover:bg-primary/90 focus-visible:ring-primary focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
						type="button"
					>
						<RotateCcw class="w-4 h-4" />
						Réessayer
					</button>
				{:else}
					<Dialog.Close
						class="h-input rounded-input bg-transparent text-foreground hover:bg-muted focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] border cursor-pointer"
					>
						<XIcon class="w-4 h-4" />
						Rejeter
					</Dialog.Close>
					<button
						onclick={onAccept}
						class="h-input rounded-input bg-primary text-primary-foreground shadow-mini hover:bg-primary/90 focus-visible:ring-primary focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center gap-2 px-4 py-2 text-sm font-medium focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
						type="button"
					>
						<Check class="w-4 h-4" />
						Appliquer
					</button>
				{/if}
			</div>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>

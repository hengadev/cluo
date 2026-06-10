<script lang="ts">
	import { currentCase } from "$lib/stores/case";
	import { page } from "$app/stores";
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { fetchCaseInvoices, fetchCaseMandates, fetchCaseEstimates, fetchCaseContracts } from "$lib/services/api";
	import { ReceiptEuro, Handshake, FileText, ShieldCheck } from "@lucide/svelte";

	const caseId = $derived($page.params.id);

	// Update the current case store when navigating to a case's documents
	$effect(() => {
		if (caseId && caseId !== $currentCase.id) {
			currentCase.setCase(caseId);
		}
	});

	let invoices: any[] = [];
	let mandates: any[] = [];
	let estimates: any[] = [];
	let contracts: any[] = [];
	let loading = true;
	let error: string | null = null;

	onMount(async () => {
		if (!caseId) {
			loading = false;
			return;
		}

		try {
			const [invoicesData, mandatesData, estimatesData, contractsData] = await Promise.all([
				fetchCaseInvoices(caseId),
				fetchCaseMandates(caseId),
				fetchCaseEstimates(caseId),
				fetchCaseContracts(caseId)
			]);

			invoices = invoicesData;
			mandates = mandatesData;
			estimates = estimatesData;
			contracts = contractsData;
		} catch (e) {
			error = e instanceof Error ? e.message : "Erreur lors du chargement des documents";
		} finally {
			loading = false;
		}
	});

	const documentTypes = $derived([
		{
			type: "facture",
			title: "Factures",
			icon: ReceiptEuro,
			count: invoices.length,
			description: "Factures du dossier"
		},
		{
			type: "mandate",
			title: "Mandats",
			icon: Handshake,
			count: mandates.length,
			description: "Mandats du dossier"
		},
		{
			type: "estimate",
			title: "Devis",
			icon: FileText,
			count: estimates.length,
			description: "Devis du dossier"
		},
		{
			type: "contract",
			title: "Contrats",
			icon: ShieldCheck,
			count: contracts.length,
			description: "Contrats du dossier"
		}
	]);
</script>

<div class="p-8 flex flex-col gap-6">
	<div>
		<h1 class="text-3xl font-bold">Documents</h1>
		<p class="text-muted-foreground mt-1">Gestion des documents du dossier</p>
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<p class="text-muted-foreground">Chargement...</p>
		</div>
	{:else if error}
		<div class="alert-error">
			{error}
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each documentTypes as doc}
				<button
					class="border border-border-card rounded-card p-6 bg-background hover:border-border-input-hover hover:shadow-card hover:-translate-y-1 transition-all duration-300 text-left"
					onclick={() => goto(`/cases/${caseId}/documents/${doc.type}`)}
				>
					<div class="flex items-start justify-between mb-4">
						<svelte:component this={doc.icon} class="w-8 h-8 text-foreground" strokeWidth={1.5} />
						{#if doc.count > 0}
							<span class="bg-muted-foreground text-background px-2 py-1 rounded-full text-xs font-medium">
								{doc.count}
							</span>
						{/if}
					</div>
					<h3 class="font-semibold text-foreground text-lg">{doc.title}</h3>
					<p class="text-sm text-muted-foreground mt-2">{doc.description}</p>
				</button>
			{/each}
		</div>

		{#if invoices.length === 0 && mandates.length === 0 && estimates.length === 0 && contracts.length === 0}
			<div class="mt-12 text-center">
				<p class="text-muted-foreground">Aucun document disponible pour ce dossier</p>
			</div>
		{/if}
	{/if}
</div>

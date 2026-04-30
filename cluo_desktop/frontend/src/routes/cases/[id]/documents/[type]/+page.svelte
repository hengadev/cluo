<script lang="ts">
	import { currentCase } from "$lib/stores/case";
	import { page } from "$app/stores";
	import { onMount } from "svelte";
	import {
		fetchCaseInvoices,
		fetchCaseMandates,
		fetchCaseEstimates,
		fetchCaseContracts
	} from "$lib/services/api";
	import type { Invoice, Mandate, Estimate, Contract } from "$lib/types/entities";

	const caseId = $derived($page.params.id);
	const docType = $derived($page.params.type);

	// Update the current case store when navigating to a case's document type
	$effect(() => {
		if (caseId && caseId !== $currentCase.id) {
			currentCase.setCase(caseId);
		}
	});

	// Map document types to their display names and API functions
	const docTypeConfig: Record<string, { title: string; fetchFn: (id: string) => Promise<any[]> }> = {
		facture: { title: "Factures", fetchFn: fetchCaseInvoices },
		mandat: { title: "Mandats", fetchFn: fetchCaseMandates },
		devis: { title: "Devis", fetchFn: fetchCaseEstimates },
		contrat: { title: "Contrats", fetchFn: fetchCaseContracts }
	};

	let documents: any[] = [];
	let loading = true;
	let error: string | null = null;

	const displayName = $derived(docTypeConfig[docType]?.title || docType);

	onMount(async () => {
		if (!caseId || !docType || !docTypeConfig[docType]) {
			error = "Type de document non reconnu";
			loading = false;
			return;
		}

		try {
			documents = await docTypeConfig[docType].fetchFn(caseId);
		} catch (e) {
			error = e instanceof Error ? e.message : "Erreur lors du chargement des documents";
		} finally {
			loading = false;
		}
	});

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString("fr-FR", {
			day: "2-digit",
			month: "2-digit",
			year: "numeric"
		});
	}

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat("fr-FR", {
			style: "currency",
			currency: "EUR"
		}).format(amount);
	}

	function getDocumentId(doc: any): string {
		return doc.invoiceNumber || doc.mandateNumber || doc.estimateNumber || doc.contractNumber || doc.id;
	}

	function getDocumentStatus(doc: any): string | null {
		return doc.paymentStatus || doc.status || null;
	}

	function getDocumentAmount(doc: any): number | null {
		return doc.totalAmount || doc.estimatedTotal || doc.contractValue || null;
	}
</script>

<div class="p-8">
	<h1 class="text-2xl font-bold mb-2">{displayName}</h1>
	<p class="text-muted-foreground mb-6">Liste des {displayName.toLowerCase()} du dossier</p>

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<p class="text-muted-foreground">Chargement...</p>
		</div>
	{:else if error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
			{error}
		</div>
	{:else if documents.length === 0}
		<div class="text-center py-12">
			<p class="text-muted-foreground">Aucun {displayName.toLowerCase().slice(0, -1)} disponible pour ce dossier</p>
		</div>
	{:else}
		<div class="border border-border-card rounded-lg overflow-hidden">
			<table class="w-full">
				<thead class="bg-muted">
					<tr>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
							Référence
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
							Date
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
							Description
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
							Montant
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
							Statut
						</th>
					</tr>
				</thead>
				<tbody class="bg-background divide-y divide-border">
					{#each documents as doc}
						<tr class="hover:bg-muted/50 transition-colors">
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="text-sm font-medium text-foreground">{getDocumentId(doc)}</div>
								<div class="text-sm text-muted-foreground">{doc.id}</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="text-sm text-foreground">{formatDate(doc.issueDate || doc.startDate)}</div>
							</td>
							<td class="px-6 py-4">
								<div class="text-sm text-foreground max-w-xs truncate">
									{doc.scopeOfWork || doc.scopeOfServices || doc.notes || doc.termsConditions || '-'}
								</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								{#if getDocumentAmount(doc) !== null}
									<div class="text-sm font-medium text-foreground">
										{formatCurrency(getDocumentAmount(doc))}
									</div>
								{:else}
									<div class="text-sm text-muted-foreground">-</div>
								{/if}
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								{#if getDocumentStatus(doc)}
									<span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-blue-100 text-blue-800">
										{getDocumentStatus(doc)}
									</span>
								{:else}
									<span class="text-sm text-muted-foreground">-</span>
								{/if}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

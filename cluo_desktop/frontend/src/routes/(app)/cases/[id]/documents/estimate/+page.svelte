<script lang="ts">
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { goto } from "$app/navigation";
	import {
		Pencil,
		Send,
		CheckCircle,
		Plus,
		Trash2,
		X,
		Save,
		FileText,
		ChevronLeft,
		AlertTriangle,
		Printer,
	} from "@lucide/svelte";
	import { Dialog } from "bits-ui";
	import {
		fetchCase,
		fetchClient,
		fetchCaseEstimates,
		createEstimate,
		updateEstimate,
		deleteDocument,
		sendDocument,
		acceptEstimate,
		openDocumentPDF,
	} from "$lib/services/api";
	import { currentCase } from "$lib/stores/case";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import DocumentLifecycleStepper from "$lib/custom/documents/DocumentLifecycleStepper.svelte";
	import { documentStatusBadge } from "$lib/utils/badgeVariants";
	import type { Case, Client, Estimate, EstimateItem, Mandate, DocumentStatus } from "$lib/types/entities";

	const toastState = getToastContext();

	// Route params
	const caseId = $derived($page.params.id);

	// Update the current case store
	$effect(() => {
		if (caseId && caseId !== $currentCase.id) {
			currentCase.setCase(caseId);
		}
	});

	// Core state
	let caseData: Case | null = $state(null);
	let client: Client | null = $state(null);
	let estimates: Estimate[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	// Currently selected estimate for viewing/editing
	let selectedEstimate: Estimate | null = $state(null);
	let viewMode: "list" | "detail" | "edit" = $state("list");
	let showCreateModal = $state(false);

	// Create/Edit form state
	let formIssueDate = $state(todayISO());
	let formValidUntil = $state("");
	let formNotes = $state("");
	let formLineItems = $state<FormItem[]>([{ description: "", quantity: 1, unit_price: 0 }]);
	let formSaving = $state(false);

	// Lifecycle action state
	let sendingEstimate = $state(false);
	let acceptingEstimate = $state(false);
	let previewingEstimateId: string | null = $state(null);
	let deletingEstimateId: string | null = $state(null);

	// Status labels and colors
	const STATUS_LABELS: Record<string, string> = {
		draft: "Brouillon",
		sent: "Envoyé",
		signed: "Signé",
		active: "Actif",
		archived: "Archivé",
		cancelled: "Annulé",
		rejected: "Rejeté",
		expired: "Expiré",
	};

	const ESTIMATE_STEPS = [
		{ key: "draft", label: "Brouillon" },
		{ key: "sent", label: "Envoyé" },
		{ key: "signed", label: "Accepté" },
	];

	function statusNote(est: Estimate): string {
		if (est.status !== "draft") return "Ce devis ne peut plus être modifié.";
		return "";
	}

	interface FormItem {
		description: string;
		quantity: number;
		unit_price: number;
	}

	function todayISO(): string {
		return new Date().toISOString().split("T")[0];
	}

	function formatDate(dateStr: string): string {
		if (!dateStr) return "—";
		return new Date(dateStr).toLocaleDateString("fr-FR", {
			day: "2-digit",
			month: "short",
			year: "numeric",
		});
	}

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat("fr-FR", {
			style: "currency",
			currency: "EUR",
		}).format(amount);
	}

	function formTotal(): number {
		return formLineItems.reduce(
			(sum, item) => sum + item.quantity * item.unit_price,
			0,
		);
	}

	function isDraft(est: Estimate): boolean {
		return est.status === "draft";
	}

	function canSend(est: Estimate): boolean {
		return est.status === "draft";
	}

	function canAccept(est: Estimate): boolean {
		return est.status === "sent";
	}

	function canDelete(est: Estimate): boolean {
		return est.status === "draft";
	}

	// =========================================================================
	// Data loading
	// =========================================================================

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		if (!caseId) {
			loading = false;
			return;
		}
		loading = true;
		error = null;

		try {
			const [c, ests] = await Promise.all([
				fetchCase(caseId),
				fetchCaseEstimates(caseId),
			]);
			caseData = c;
			estimates = ests;

			if (c?.clientId) {
				client = await fetchClient(c.clientId);
			}
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Erreur lors du chargement des devis";
		} finally {
			loading = false;
		}
	}

	// =========================================================================
	// Navigation helpers
	// =========================================================================

	function showList() {
		selectedEstimate = null;
		viewMode = "list";
		formLineItems = [{ description: "", quantity: 1, unit_price: 0 }];
		formNotes = "";
		formIssueDate = todayISO();
		formValidUntil = "";
	}

	function showCreate() {
		formIssueDate = todayISO();
		formValidUntil = "";
		formNotes = "";
		formLineItems = [{ description: "", quantity: 1, unit_price: 0 }];
		showCreateModal = true;
	}

	function showDetail(est: Estimate) {
		selectedEstimate = est;
		viewMode = "detail";
	}

	async function handlePreview(est: Estimate) {
		previewingEstimateId = est.id;
		try {
			await openDocumentPDF(est.id, "estimate");
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible d'afficher l'aperçu du devis.",
			);
		} finally {
			previewingEstimateId = null;
		}
	}

	async function handleDelete(est: Estimate) {
		if (!canDelete(est)) return;
		deletingEstimateId = est.id;
		try {
			await deleteDocument(est.id, "estimate");
			estimates = estimates.filter((x) => x.id !== est.id);
			if (selectedEstimate?.id === est.id) showList();
			toastState.add(TOAST_LEVELS.Info, "Devis supprimé", "Le devis a été supprimé.");
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de supprimer le devis.",
			);
		} finally {
			deletingEstimateId = null;
		}
	}

	function showEdit(est: Estimate) {
		selectedEstimate = est;
		viewMode = "edit";
		formIssueDate = est.issue_date ? est.issue_date.split("T")[0] : todayISO();
		formValidUntil = est.valid_until ? est.valid_until.split("T")[0] : "";
		formNotes = est.notes || "";
		formLineItems = est.line_items.map((li) => ({
			description: li.description,
			quantity: li.quantity,
			unit_price: li.unit_price,
		}));
		if (formLineItems.length === 0) {
			formLineItems = [{ description: "", quantity: 1, unit_price: 0 }];
		}
	}

	// =========================================================================
	// Line item helpers
	// =========================================================================

	function addLineItem() {
		formLineItems = [...formLineItems, { description: "", quantity: 1, unit_price: 0 }];
	}

	function removeLineItem(index: number) {
		if (formLineItems.length <= 1) return;
		formLineItems = formLineItems.filter((_, i) => i !== index);
	}

	// =========================================================================
	// Create Estimate
	// =========================================================================

	async function handleCreate() {
		if (!caseData) return;
		if (formLineItems.some((li) => !li.description.trim())) {
			toastState.add(TOAST_LEVELS.Error, "Erreur", "Toutes les lignes doivent avoir une description.");
			return;
		}

		formSaving = true;
		try {
			const lineItems = formLineItems.map((li) => ({
				description: li.description.trim(),
				quantity: li.quantity,
				unit_price: li.unit_price,
				subtotal: li.quantity * li.unit_price,
			}));

			const estimateNumber = `DEV-${new Date().getFullYear()}-${String(estimates.length + 1).padStart(3, "0")}`;

			const payload: Partial<Estimate> = {
				case_id: caseData.id,
				client_id: caseData.clientId,
				estimate_number: estimateNumber,
				issue_date: new Date(formIssueDate).toISOString(),
				valid_until: formValidUntil ? new Date(formValidUntil).toISOString() : undefined,
				line_items: lineItems as EstimateItem[],
				estimated_total: formTotal(),
				notes: formNotes.trim() || undefined,
				status: "draft" as DocumentStatus,
			} as Partial<Estimate>;

			const result = await createEstimate(payload as Estimate);
			if (result.data) {
				estimates = [...estimates, result.data];
				selectedEstimate = result.data;
				showCreateModal = false;
				viewMode = "detail";
				toastState.add(TOAST_LEVELS.Info, "Devis créé", "Le devis a été créé en brouillon.");
			}
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de créer le devis.",
			);
		} finally {
			formSaving = false;
		}
	}

	// =========================================================================
	// Edit Estimate (only draft)
	// =========================================================================

	async function handleEdit() {
		if (!selectedEstimate || !isDraft(selectedEstimate)) return;
		if (formLineItems.some((li) => !li.description.trim())) {
			toastState.add(TOAST_LEVELS.Error, "Erreur", "Toutes les lignes doivent avoir une description.");
			return;
		}

		formSaving = true;
		try {
			const lineItems = formLineItems.map((li) => ({
				description: li.description.trim(),
				quantity: li.quantity,
				unit_price: li.unit_price,
				subtotal: li.quantity * li.unit_price,
			}));

			const result = await updateEstimate(selectedEstimate.id, lineItems);
			if (result.success) {
				const freshEstimates = await fetchCaseEstimates(caseId);
				estimates = freshEstimates;
				const updated = freshEstimates.find((e) => e.id === selectedEstimate!.id);
				if (updated) selectedEstimate = updated;
				viewMode = "detail";
				toastState.add(TOAST_LEVELS.Info, "Devis mis à jour", "Les modifications ont été enregistrées.");
			}
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de mettre à jour le devis.",
			);
		} finally {
			formSaving = false;
		}
	}

	// =========================================================================
	// Send Estimate (draft → sent)
	// =========================================================================

	async function handleSend() {
		if (!selectedEstimate || !canSend(selectedEstimate)) return;
		sendingEstimate = true;
		try {
			const result = await sendDocument(selectedEstimate.id, "estimate", {
				recipients: [],
				subject: `Devis ${selectedEstimate.estimate_number}`,
				message: "",
				send_email: true,
				send_sms: false,
			});

			if (result.success) {
				const freshEstimates = await fetchCaseEstimates(caseId);
				estimates = freshEstimates;
				const updated = freshEstimates.find((e) => e.id === selectedEstimate!.id);
				if (updated) selectedEstimate = updated;
				toastState.add(
					TOAST_LEVELS.Info,
					"Devis envoyé",
					"Le devis a été marqué comme envoyé.",
				);
			}
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible d'envoyer le devis.",
			);
		} finally {
			sendingEstimate = false;
		}
	}

	// =========================================================================
	// Accept Estimate (sent → signed; derives a Mandate)
	// =========================================================================

	async function handleAccept() {
		if (!selectedEstimate || !canAccept(selectedEstimate)) return;
		acceptingEstimate = true;
		try {
			const result = await acceptEstimate(selectedEstimate.id);
			if (result.success) {
				const freshEstimates = await fetchCaseEstimates(caseId);
				estimates = freshEstimates;

				const mandate = result.data as Partial<Mandate> | null;
				const mandateRef = mandate?.mandate_number;
				toastState.add(
					TOAST_LEVELS.Info,
					"Devis accepté",
					mandateRef
						? `Un mandat (${mandateRef}) a été créé.`
						: "Le devis a été accepté et un mandat a été créé.",
				);
				await goto(`/cases/${caseId}/documents/mandate`);
			}
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible d'accepter le devis.",
			);
		} finally {
			acceptingEstimate = false;
		}
	}
</script>

<div class="p-8 flex flex-col flex-1 min-h-0 gap-6">
	{#if loading}
		<div class="flex items-center justify-center py-12">
			<p class="text-muted-foreground">Chargement...</p>
		</div>
	{:else if error}
		<div class="alert-error">
			{error}
		</div>
	{:else}
		<!-- Header -->
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-3">
				{#if viewMode !== "list"}
					<button
						onclick={showList}
						class="p-2 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
						title="Retour à la liste"
					>
						<ChevronLeft size={18} />
					</button>
				{/if}
				<div>
					<h1 class="text-2xl font-bold text-foreground">
						{#if viewMode === "edit" && selectedEstimate}
							Modifier le devis {selectedEstimate.estimate_number}
						{:else if viewMode === "detail" && selectedEstimate}
							Devis {selectedEstimate.estimate_number}
						{:else}
							Devis
						{/if}
					</h1>
					<p class="text-sm text-muted-foreground">
						{#if viewMode === "list"}
							Liste des devis du dossier
						{:else if client}
							Client : {client.name}
						{/if}
					</p>
				</div>
			</div>

			{#if viewMode === "list" && estimates.length > 0}
				<button
					type="button"
					onclick={showCreate}
					class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
				>
					Nouveau devis
				</button>
			{/if}

			{#if viewMode === "detail" && selectedEstimate}
				<div class="flex items-center gap-2">
					<div class="flex items-center gap-1">
						<button
							type="button"
							onclick={() => selectedEstimate && handlePreview(selectedEstimate)}
							disabled={previewingEstimateId === selectedEstimate.id}
							class="h-input rounded-input bg-transparent text-foreground hover:bg-muted inline-flex items-center justify-center gap-1.5 px-3 text-sm font-medium active:scale-[0.98] border border-border-input cursor-pointer disabled:opacity-50 transition-interactive duration-150"
						>
							<Printer size={14} />
							Aperçu
						</button>
						{#if isDraft(selectedEstimate)}
							<button
								type="button"
								onclick={() => { if (selectedEstimate) showEdit(selectedEstimate); }}
								class="h-input rounded-input bg-transparent text-foreground hover:bg-muted inline-flex items-center justify-center gap-1.5 px-3 text-sm font-medium active:scale-[0.98] border border-border-input cursor-pointer transition-interactive duration-150"
							>
								<Pencil size={14} />
								Modifier
							</button>
						{/if}
					</div>

					{#if canDelete(selectedEstimate) || canSend(selectedEstimate) || canAccept(selectedEstimate)}
						<div class="w-px h-5 bg-border"></div>
					{/if}

					{#if canDelete(selectedEstimate)}
						<ConfirmDialog
							title="Supprimer le devis"
							description="Le devis sera définitivement supprimé. Cette action est irréversible."
							confirmLabel="Supprimer"
							onConfirm={() => { if (selectedEstimate) return handleDelete(selectedEstimate); }}
						>
							<button
								type="button"
								disabled={deletingEstimateId === selectedEstimate.id}
								class="h-input rounded-input inline-flex items-center justify-center gap-1.5 px-3 text-sm font-medium btn-ghost-destructive cursor-pointer disabled:opacity-50 transition-interactive duration-150"
							>
								<Trash2 size={14} />
								Supprimer
							</button>
						</ConfirmDialog>
					{/if}
					{#if canSend(selectedEstimate)}
						<ConfirmDialog
							title="Envoyer le devis"
							description="Le devis sera marqué comme envoyé et un email sera envoyé au client. Cette action est irréversible."
							confirmLabel="Envoyer"
							onConfirm={handleSend}
						>
							<button
								type="button"
								disabled={sendingEstimate}
								class="h-input rounded-input bg-accent text-accent-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-1.5 px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<Send size={14} />
								{sendingEstimate ? "Envoi..." : "Envoyer"}
							</button>
						</ConfirmDialog>
					{/if}
					{#if canAccept(selectedEstimate)}
						<ConfirmDialog
							title="Accepter le devis"
							description="Le devis sera marqué comme accepté et un mandat sera automatiquement créé. Cette action est irréversible."
							confirmLabel="Accepter"
							onConfirm={handleAccept}
						>
							<button
								type="button"
								disabled={acceptingEstimate}
								class="h-input rounded-input bg-success text-success-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-1.5 px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<CheckCircle size={14} />
								{acceptingEstimate ? "Acceptation..." : "Accepter"}
							</button>
						</ConfirmDialog>
					{/if}
				</div>
			{/if}
		</div>

		<!-- ================================================================ -->
		<!-- LIST VIEW -->
		<!-- ================================================================ -->
		{#if viewMode === "list"}
			{#if estimates.length === 0}
				<div class="border border-dashed border-border rounded-lg p-12 bg-muted/20 flex flex-col items-center justify-center gap-4 flex-1 min-h-[50vh]">
					<FileText class="w-12 h-12 text-muted-foreground" />
					<p class="text-muted-foreground text-center">Aucun devis pour ce dossier.</p>
					<button
						type="button"
						onclick={showCreate}
						class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
					>
						Créer un devis
					</button>
				</div>
			{:else}
				<div class="border border-border-card rounded-lg overflow-hidden">
					<table class="w-full">
						<thead class="bg-muted">
							<tr>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Référence</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Date d'émission</th>
								<th class="px-6 py-3 text-right text-xs font-medium text-muted-foreground uppercase tracking-wider">Montant</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Statut</th>
								<th class="px-6 py-3 w-12"></th>
							</tr>
						</thead>
						<tbody class="bg-background divide-y divide-border">
							{#each estimates as est (est.id)}
								<tr
									class="hover:shadow-mini hover:relative transition-interactive duration-150 cursor-pointer"
									onclick={() => showDetail(est)}
								>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm font-medium text-foreground">{est.estimate_number}</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm text-foreground">{formatDate(est.issue_date)}</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-right">
										<div class="text-sm font-medium text-foreground">{formatCurrency(est.estimated_total)}</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {documentStatusBadge(est.status as DocumentStatus)}">
											{STATUS_LABELS[est.status] || est.status}
										</span>
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-right">
										<div class="flex items-center justify-end gap-1">
											<button
												type="button"
												onclick={(e) => { e.stopPropagation(); handlePreview(est); }}
												disabled={previewingEstimateId === est.id}
												class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer disabled:opacity-50"
												title="Aperçu / Imprimer"
											>
												<Printer size={16} />
											</button>
											{#if isDraft(est)}
												<button
													type="button"
													onclick={(e) => { e.stopPropagation(); showEdit(est); }}
													class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
													title="Modifier"
												>
													<Pencil size={16} />
												</button>
											{/if}
											{#if canDelete(est)}
												<ConfirmDialog
													title="Supprimer le devis"
													description="Le devis sera définitivement supprimé. Cette action est irréversible."
													confirmLabel="Supprimer"
													onConfirm={() => handleDelete(est)}
												>
													<button
														type="button"
														onclick={(e) => e.stopPropagation()}
														disabled={deletingEstimateId === est.id}
														class="p-1.5 rounded btn-ghost-destructive cursor-pointer disabled:opacity-50"
														title="Supprimer"
													>
														<Trash2 size={16} />
													</button>
												</ConfirmDialog>
											{/if}
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}

		<!-- ================================================================ -->
		<!-- EDIT FORM (inline) -->
		<!-- ================================================================ -->
		{:else if viewMode === "edit"}
			<div class="border border-border-card rounded-lg p-6 max-w-3xl animate-fade-in">
				<!-- Dates row -->
				<div class="grid grid-cols-2 gap-4 mb-6">
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Date d'émission *</label>
						<input
							type="date"
							bind:value={formIssueDate}
							disabled
							class="h-input rounded-input border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2 disabled:opacity-50"
						/>
					</div>
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Valide jusqu'au</label>
						<input
							type="date"
							bind:value={formValidUntil}
							disabled
							class="h-input rounded-input border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2 disabled:opacity-50"
						/>
					</div>
				</div>

				<!-- Line items -->
				<div class="mb-6">
					<div class="flex items-center justify-between mb-3">
						<label class="text-sm font-medium text-foreground">Lignes de devis</label>
						<button
							type="button"
							onclick={addLineItem}
							class="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
						>
							<Plus size={14} />
							Ajouter une ligne
						</button>
					</div>

					<div class="space-y-3">
						{#each formLineItems as item, i (i)}
							<div class="flex items-start gap-3 border border-border rounded-lg p-3 bg-muted/30">
								<div class="flex-1">
									<input
										type="text"
										bind:value={item.description}
										placeholder="Description de la prestation"
										class="h-input rounded-input border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
									/>
								</div>
								<div class="w-24">
									<label class="text-xs text-muted-foreground mb-1 block">Quantité</label>
									<input
										type="number"
										bind:value={item.quantity}
										min="1"
										step="1"
										class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
									/>
								</div>
								<div class="w-32">
									<label class="text-xs text-muted-foreground mb-1 block">Prix unitaire (€)</label>
									<input
										type="number"
										bind:value={item.unit_price}
										min="0"
										step="0.01"
										class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
									/>
								</div>
								<div class="w-28 pt-5 text-right">
									<span class="text-sm font-medium text-foreground">
										{formatCurrency(item.quantity * item.unit_price)}
									</span>
								</div>
								{#if formLineItems.length > 1}
									<button
										type="button"
										onclick={() => removeLineItem(i)}
										class="p-2 mt-3 rounded btn-ghost-destructive cursor-pointer"
										title="Supprimer la ligne"
									>
										<Trash2 size={14} />
									</button>
								{/if}
							</div>
						{/each}
					</div>

					<!-- Total -->
					<div class="flex justify-end mt-4 pt-4 border-t border-border">
						<div class="text-right">
							<p class="text-sm text-muted-foreground">Total estimé</p>
							<p class="text-2xl font-bold text-foreground">{formatCurrency(formTotal())}</p>
						</div>
					</div>
				</div>

				<!-- Notes -->
				<div class="mb-6">
					<label class="text-xs text-muted-foreground mb-1 block">Notes</label>
					<textarea
						bind:value={formNotes}
						placeholder="Notes internes ou conditions..."
						rows="3"
						disabled
						class="rounded-input border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 py-2 text-sm focus:ring-2 focus:ring-offset-2 resize-none disabled:opacity-50"
					></textarea>
				</div>

				<!-- Actions -->
				<div class="flex justify-end gap-2">
					<button
						type="button"
						onclick={showList}
						class="h-input rounded-input bg-transparent text-dark hover:bg-muted inline-flex items-center justify-center gap-1.5 px-4 text-sm font-semibold active:scale-[0.98] border-2 border-border-input cursor-pointer"
					>
						<X size={14} />
						Annuler
					</button>
					<button
						type="button"
						onclick={handleEdit}
						disabled={formSaving}
						class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-1.5 px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
					>
						<Save size={14} />
						{formSaving ? "Enregistrement..." : "Enregistrer"}
					</button>
				</div>
			</div>

		<!-- ================================================================ -->
		<!-- DETAIL VIEW -->
		<!-- ================================================================ -->
		{:else if viewMode === "detail" && selectedEstimate}
			<div class="max-w-5xl animate-fade-in flex flex-col gap-6">
				<DocumentLifecycleStepper
					steps={ESTIMATE_STEPS}
					status={selectedEstimate.status}
					statusLabel={STATUS_LABELS[selectedEstimate.status] || selectedEstimate.status}
					note={statusNote(selectedEstimate)}
				/>

				<div class="grid grid-cols-1 lg:grid-cols-[1fr_320px] gap-6">
					<!-- Primary content -->
					<div class="flex flex-col gap-6 min-w-0">
						<div class="border border-border-card rounded-lg p-6">
							{#if selectedEstimate.notes}
								<div class="mb-4">
									<p class="text-xs text-muted-foreground mb-1">Notes</p>
									<p class="text-sm text-muted-foreground italic">{selectedEstimate.notes}</p>
								</div>
							{/if}

							<!-- Line items table -->
							<div class="border border-border rounded-lg overflow-hidden">
								<table class="w-full">
									<thead class="bg-muted">
										<tr>
											<th class="px-4 py-2 text-left text-xs font-medium text-muted-foreground uppercase">Description</th>
											<th class="px-4 py-2 text-right text-xs font-medium text-muted-foreground uppercase w-20">Qté</th>
											<th class="px-4 py-2 text-right text-xs font-medium text-muted-foreground uppercase w-28">Prix unit.</th>
											<th class="px-4 py-2 text-right text-xs font-medium text-muted-foreground uppercase w-28">Sous-total</th>
										</tr>
									</thead>
									<tbody class="divide-y divide-border">
										{#each selectedEstimate.line_items as item, i (i)}
											<tr>
												<td class="px-4 py-3 text-sm text-foreground">{item.description}</td>
												<td class="px-4 py-3 text-sm text-foreground text-right">{item.quantity}</td>
												<td class="px-4 py-3 text-sm text-foreground text-right">{formatCurrency(item.unit_price)}</td>
												<td class="px-4 py-3 text-sm font-medium text-foreground text-right">{formatCurrency(item.subtotal)}</td>
											</tr>
										{/each}
									</tbody>
									<tfoot class="bg-muted/50 border-t border-border">
										<tr>
											<td colspan="3" class="px-4 py-3 text-sm font-semibold text-foreground text-right">Total</td>
											<td class="px-4 py-3 text-lg font-bold text-foreground text-right">
												{formatCurrency(selectedEstimate.estimated_total)}
											</td>
										</tr>
									</tfoot>
								</table>
							</div>
						</div>

						<!-- Lifecycle warning for non-actionable states -->
						{#if !canSend(selectedEstimate) && !canAccept(selectedEstimate) && !isDraft(selectedEstimate)}
							<div class="flex items-start gap-2 text-sm text-muted-foreground bg-muted/30 border border-border rounded-lg p-4">
								<AlertTriangle size={16} class="flex-shrink-0 mt-0.5" />
								<p>
									Ce devis est dans l'état <strong>{STATUS_LABELS[selectedEstimate.status]}</strong>.
									Aucune action n'est disponible pour le moment.
								</p>
							</div>
						{/if}
					</div>

					<!-- Metadata rail -->
					<div class="flex flex-col gap-6">
						<div class="border border-border-card rounded-lg p-6">
							<p class="text-xs text-muted-foreground mb-3 font-medium uppercase tracking-wider">Détails</p>
							<div class="flex flex-col gap-4">
								<div>
									<p class="text-xs text-muted-foreground mb-1">Référence</p>
									<p class="text-sm font-semibold text-foreground">{selectedEstimate.estimate_number}</p>
								</div>
								<div>
									<p class="text-xs text-muted-foreground mb-1">Devise</p>
									<p class="text-sm text-foreground">EUR (€)</p>
								</div>
								<div>
									<p class="text-xs text-muted-foreground mb-1">Date d'émission</p>
									<p class="text-sm text-foreground">{formatDate(selectedEstimate.issue_date)}</p>
								</div>
								<div>
									<p class="text-xs text-muted-foreground mb-1">Valide jusqu'au</p>
									<p class="text-sm text-foreground">
										{selectedEstimate.valid_until ? formatDate(selectedEstimate.valid_until) : "—"}
									</p>
								</div>
							</div>
						</div>

						{#if selectedEstimate.accepted && selectedEstimate.accepted_at}
							<div class="bg-success/10 border border-success/30 rounded-lg p-3 flex items-center gap-2">
								<CheckCircle size={16} class="text-success flex-shrink-0" />
								<p class="text-sm text-success">
									Accepté le {formatDate(selectedEstimate.accepted_at)}
								</p>
							</div>
						{/if}
					</div>
				</div>
			</div>
		{/if}
	{/if}
</div>

<!-- ================================================================ -->
<!-- CREATE MODAL -->
<!-- Hallmark · component: modal · genre: modern-minimal · theme: system
 * states: default · hover · focus · active · disabled · saving
 * contrast: pass
 -->
<!-- ================================================================ -->
<Dialog.Root bind:open={showCreateModal}>
	<Dialog.Portal>
		<Dialog.Overlay
			class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/60 backdrop-blur-[2px]"
		/>
		<Dialog.Content
			class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-2xl translate-x-[-50%] translate-y-[-50%] border flex flex-col max-h-[90vh]"
		>
			<!-- Modal header -->
			<div class="flex-shrink-0 px-8 pt-7 pb-5 border-b border-border-card">
				<Dialog.Title class="text-base font-semibold tracking-tight text-foreground">
					Nouveau devis
				</Dialog.Title>
				{#if client}
					<p class="text-sm text-muted-foreground mt-0.5">{client.name}</p>
				{/if}
				<Dialog.Close
					class="absolute right-5 top-6 rounded-md text-muted-foreground hover:text-foreground transition-interactive duration-150 cursor-pointer p-0.5"
				>
					<X class="size-4" />
					<span class="sr-only">Fermer</span>
				</Dialog.Close>
			</div>

			<!-- Modal body (scrollable) -->
			<div class="flex-1 min-h-0 overflow-y-auto px-8 py-6 flex flex-col gap-5">

				<!-- Dates -->
				<div class="grid grid-cols-2 gap-4">
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Date d'émission</label>
						<input
							type="date"
							bind:value={formIssueDate}
							class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">
							Valide jusqu'au
							<span class="font-normal text-muted-foreground">(optionnel)</span>
						</label>
						<input
							type="date"
							bind:value={formValidUntil}
							class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
				</div>

				<!-- Prestations block -->
				<div class="flex flex-col gap-3">
					<!-- Inline section separator -->
					<div class="flex items-center gap-3">
						<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider whitespace-nowrap">Prestations</span>
						<div class="h-px flex-1 bg-border-input"></div>
						<button
							type="button"
							onclick={addLineItem}
							class="inline-flex items-center gap-1 text-xs font-medium text-accent hover:text-accent/70 transition-interactive duration-150 cursor-pointer shrink-0"
						>
							<Plus size={12} />
							Ajouter une ligne
						</button>
					</div>

					<!-- Column headers -->
					<div class="grid items-center gap-2 px-0.5" style="grid-template-columns: 1fr 4.5rem 7.5rem 6.5rem 1.75rem">
						<span class="text-xs text-muted-foreground">Description</span>
						<span class="text-xs text-muted-foreground text-right">Qté</span>
						<span class="text-xs text-muted-foreground text-right">Prix unit. (€)</span>
						<span class="text-xs text-muted-foreground text-right">Sous-total</span>
						<span></span>
					</div>

					<!-- Line item rows -->
					<div class="flex flex-col gap-2">
						{#each formLineItems as item, i (i)}
							<div class="grid items-center gap-2" style="grid-template-columns: 1fr 4.5rem 7.5rem 6.5rem 1.75rem">
								<input
									type="text"
									bind:value={item.description}
									placeholder="Description de la prestation"
									class="h-input rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2"
								/>
								<input
									type="number"
									bind:value={item.quantity}
									min="1"
									step="1"
									class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-2 text-sm text-right tabular-nums focus:ring-2 focus:ring-offset-2"
								/>
								<input
									type="number"
									bind:value={item.unit_price}
									min="0"
									step="0.01"
									class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-2 text-sm text-right tabular-nums focus:ring-2 focus:ring-offset-2"
								/>
								<div class="h-input flex items-center justify-end pr-1">
									<span class="text-sm font-medium text-foreground tabular-nums">
										{formatCurrency(item.quantity * item.unit_price)}
									</span>
								</div>
								{#if formLineItems.length > 1}
									<button
										type="button"
										onclick={() => removeLineItem(i)}
										class="size-7 rounded flex items-center justify-center btn-ghost-destructive cursor-pointer"
										title="Supprimer la ligne"
									>
										<Trash2 size={13} />
									</button>
								{:else}
									<div class="size-7"></div>
								{/if}
							</div>
						{/each}
					</div>

					<!-- Total -->
					<div class="flex items-center justify-between pt-3 border-t border-border-input">
						<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider">Total estimé</span>
						<span class="text-xl font-semibold text-foreground tabular-nums" style="font-family: var(--font-display)">
							{formatCurrency(formTotal())}
						</span>
					</div>
				</div>

				<!-- Notes block -->
				<div class="flex flex-col gap-3">
					<!-- Inline section separator -->
					<div class="flex items-center gap-3">
						<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider whitespace-nowrap">Notes</span>
						<div class="h-px flex-1 bg-border-input"></div>
						<span class="text-xs text-muted-foreground">(optionnel)</span>
					</div>

					<textarea
						bind:value={formNotes}
						placeholder="Conditions particulières, délais, modalités de paiement..."
						rows="3"
						class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"
					></textarea>
				</div>

			</div>

			<!-- Modal footer -->
			<div class="flex items-center justify-between px-8 py-4 border-t border-border-card shrink-0">
				<p class="text-xs text-muted-foreground">Le devis sera créé en brouillon.</p>
				<div class="flex items-center gap-2">
					<Dialog.Close
						class="h-input rounded-input bg-transparent text-foreground hover:bg-muted inline-flex items-center justify-center px-4 text-sm font-medium active:scale-[0.98] border border-border-input cursor-pointer transition-interactive duration-150 focus:outline-none"
					>
						Annuler
					</Dialog.Close>
					<button
						type="button"
						onclick={handleCreate}
						disabled={formSaving}
						class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-1.5 px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50 transition-interactive duration-150"
					>
						<Save size={14} />
						{formSaving ? "Enregistrement..." : "Créer le devis"}
					</button>
				</div>
			</div>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>

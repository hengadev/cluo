<script lang="ts">
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import {
		Send,
		Banknote,
		Ban,
		ChevronLeft,
		AlertTriangle,
		FileText,
		Plus,
		X,
		Save,
		Trash2,
	} from "@lucide/svelte";
	import { Dialog } from "bits-ui";
	import {
		fetchCase,
		fetchClient,
		fetchCaseInvoices,
		createInvoice,
		sendDocument,
		processPayment,
		voidInvoice,
		ConflictError,
	} from "$lib/services/api";
	import { currentCase } from "$lib/stores/case";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import { documentStatusBadge, paymentStatusBadge } from "$lib/utils/badgeVariants";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import type { Case, Client, Invoice, InvoiceItem, PaymentRequest, DocumentStatus, PaymentStatus } from "$lib/types/entities";

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
	let invoices: Invoice[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	// Currently selected invoice for viewing
	let selectedInvoice: Invoice | null = $state(null);
	let viewMode: "list" | "detail" = $state("list");
	let showCreateModal = $state(false);

	// Create form state
	interface FormItem {
		description: string;
		quantity: number;
		unit_price: number;
	}
	let formIssueDate = $state(todayISO());
	let formDueDate = $state("");
	let formTaxRate = $state(0);
	let formNotes = $state("");
	let formPaymentTerms = $state("");
	let formLineItems = $state<FormItem[]>([{ description: "", quantity: 1, unit_price: 0 }]);
	let formSaving = $state(false);

	// Lifecycle action state
	let sendingInvoice = $state(false);
	let voidingInvoice = $state(false);

	// Payment form state
	let showPaymentForm = $state(false);
	let paymentAmount = $state("");
	let paymentMethod = $state("bank_transfer");
	let paymentSubmitting = $state(false);

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
		void: "Annulé",
	};

	const PAYMENT_STATUS_LABELS: Record<string, string> = {
		unpaid: "Non payée",
		paid: "Payée",
		partially_paid: "Partiellement payée",
		overdue: "En retard",
		refunded: "Remboursée",
		void: "Annulée",
	};

	const PAYMENT_METHODS: Record<string, string> = {
		bank_transfer: "Virement bancaire",
		cheque: "Chèque",
		cash: "Espèces",
		card: "Carte bancaire",
		online: "Paiement en ligne",
		other: "Autre",
	};

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

	function formatCurrency(amount: number, currency: string = "EUR"): string {
		return new Intl.NumberFormat("fr-FR", {
			style: "currency",
			currency: currency || "EUR",
		}).format(amount);
	}

	// =========================================================================
	// Lifecycle guards
	// =========================================================================

	function canSend(inv: Invoice): boolean {
		return inv.status === "draft";
	}

	function canPay(inv: Invoice): boolean {
		return (
			(inv.status === "sent" || inv.status === "active") &&
			inv.payment_status !== "paid" &&
			inv.payment_status !== "void"
		);
	}

	function canVoid(inv: Invoice): boolean {
		return (
			inv.payment_status !== "paid" &&
			inv.payment_status !== "void" &&
			inv.status !== "cancelled" &&
			inv.status !== "archived"
		);
	}

	function hasNoActions(inv: Invoice): boolean {
		return !canSend(inv) && !canPay(inv) && !canVoid(inv);
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
			const [c, invs] = await Promise.all([
				fetchCase(caseId),
				fetchCaseInvoices(caseId),
			]);
			caseData = c;
			invoices = invs;

			if (c?.clientId) {
				client = await fetchClient(c.clientId);
			}
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Erreur lors du chargement des factures";
		} finally {
			loading = false;
		}
	}

	async function refreshInvoices() {
		const id = selectedInvoice?.id;
		if (!id) return;
		const freshInvoices = await fetchCaseInvoices(caseId);
		invoices = freshInvoices;
		const updated = freshInvoices.find((inv) => inv.id === id);
		if (updated) selectedInvoice = updated;
	}

	// =========================================================================
	// Navigation / form helpers
	// =========================================================================

	function formTotal(): number {
		const subtotal = formLineItems.reduce((sum, item) => sum + item.quantity * item.unit_price, 0);
		return subtotal + subtotal * (formTaxRate / 100);
	}

	function addLineItem() {
		formLineItems = [...formLineItems, { description: "", quantity: 1, unit_price: 0 }];
	}

	function removeLineItem(index: number) {
		if (formLineItems.length <= 1) return;
		formLineItems = formLineItems.filter((_, i) => i !== index);
	}

	function showList() {
		selectedInvoice = null;
		viewMode = "list";
		showPaymentForm = false;
	}

	function showCreate() {
		formLineItems = [{ description: "", quantity: 1, unit_price: 0 }];
		formIssueDate = todayISO();
		formDueDate = "";
		formTaxRate = 0;
		formNotes = "";
		formPaymentTerms = "";
		showCreateModal = true;
	}

	async function handleCreate() {
		if (!caseData) return;
		if (formLineItems.some((li) => !li.description.trim())) {
			toastState.add(TOAST_LEVELS.Error, "Erreur", "Toutes les lignes doivent avoir une description.");
			return;
		}
		if (!formDueDate) {
			toastState.add(TOAST_LEVELS.Error, "Erreur", "La date d'échéance est obligatoire.");
			return;
		}

		formSaving = true;
		try {
			const subtotal = formLineItems.reduce((sum, item) => sum + item.quantity * item.unit_price, 0);
			const taxAmount = subtotal * (formTaxRate / 100);
			const totalAmount = subtotal + taxAmount;
			const invoiceNumber = `FAC-${new Date().getFullYear()}-${String(invoices.length + 1).padStart(3, "0")}`;

			const payload = {
				case_id: caseData.id,
				client_id: caseData.clientId,
				invoice_number: invoiceNumber,
				issue_date: new Date(formIssueDate).toISOString(),
				due_date: new Date(formDueDate).toISOString(),
				line_items: formLineItems.map((li) => ({
					description: li.description.trim(),
					quantity: li.quantity,
					unit_price: li.unit_price,
					subtotal: li.quantity * li.unit_price,
				})),
				total_amount: totalAmount,
				tax_rate: formTaxRate,
				tax_amount: taxAmount,
				payment_status: "unpaid" as const,
				currency: "EUR",
				notes: formNotes.trim() || undefined,
				payment_terms: formPaymentTerms.trim() || undefined,
				status: "draft" as const,
			} as Invoice;

			const result = await createInvoice(payload);
			if (result.data) {
				invoices = [...invoices, result.data];
				selectedInvoice = result.data;
				showCreateModal = false;
				viewMode = "detail";
				toastState.add(TOAST_LEVELS.Info, "Facture créée", "La facture a été créée en brouillon.");
			}
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de créer la facture.",
			);
		} finally {
			formSaving = false;
		}
	}

	function showDetail(inv: Invoice) {
		selectedInvoice = inv;
		viewMode = "detail";
		showPaymentForm = false;
		paymentAmount = "";
		paymentMethod = "bank_transfer";
	}

	// =========================================================================
	// Send Invoice (draft → sent)
	// =========================================================================

	async function handleSend() {
		if (!selectedInvoice || !canSend(selectedInvoice)) return;
		sendingInvoice = true;
		try {
			const result = await sendDocument(selectedInvoice.id, "invoice", {
				recipients: [],
				subject: `Facture ${selectedInvoice.invoice_number}`,
				message: "",
				send_email: true,
				send_sms: false,
			});

			if (result.success) {
				await refreshInvoices();
				toastState.add(
					TOAST_LEVELS.Info,
					"Facture envoyée",
					"La facture a été marquée comme envoyée et un email a été envoyé au client.",
				);
			}
		} catch (e) {
			const msg =
				e instanceof ConflictError
					? e.message
					: e instanceof Error
						? e.message
						: "Impossible d'envoyer la facture.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			sendingInvoice = false;
		}
	}

	// =========================================================================
	// Record payment
	// =========================================================================

	function openPaymentForm() {
		if (!selectedInvoice) return;
		const remaining =
			selectedInvoice.total_amount - (selectedInvoice.paid_amount || 0);
		paymentAmount = remaining.toFixed(2);
		showPaymentForm = true;
	}

	function cancelPayment() {
		showPaymentForm = false;
		paymentAmount = "";
		paymentMethod = "bank_transfer";
	}

	async function handlePayment() {
		if (!selectedInvoice) return;
		const amount = parseFloat(paymentAmount);
		if (!amount || amount <= 0) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				"Le montant du paiement doit être supérieur à 0.",
			);
			return;
		}

		paymentSubmitting = true;
		try {
			const request: PaymentRequest = {
				amount,
				payment_method: paymentMethod,
			};

			const result = await processPayment(selectedInvoice.id, request);
			if (result.success) {
				await refreshInvoices();
				showPaymentForm = false;
				paymentAmount = "";
				paymentMethod = "bank_transfer";
				toastState.add(
					TOAST_LEVELS.Info,
					"Paiement enregistré",
					`Paiement de ${formatCurrency(amount, selectedInvoice.currency)} enregistré.`,
				);
			}
		} catch (e) {
			const msg =
				e instanceof ConflictError
					? e.message
					: e instanceof Error
						? e.message
						: "Impossible d'enregistrer le paiement.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			paymentSubmitting = false;
		}
	}

	// =========================================================================
	// Void Invoice
	// =========================================================================

	async function handleVoid() {
		if (!selectedInvoice || !canVoid(selectedInvoice)) return;
		voidingInvoice = true;
		try {
			const result = await voidInvoice(selectedInvoice.id);
			if (result.success) {
				await refreshInvoices();
				toastState.add(
					TOAST_LEVELS.Info,
					"Facture annulée",
					"La facture a été annulée.",
				);
			}
		} catch (e) {
			const msg =
				e instanceof ConflictError
					? e.message
					: e instanceof Error
						? e.message
						: "Impossible d'annuler la facture.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			voidingInvoice = false;
		}
	}

	// =========================================================================
	// Computed helpers
	// =========================================================================

	function remainingAmount(inv: Invoice): number {
		return inv.total_amount - (inv.paid_amount || 0);
	}

	function paymentProgress(inv: Invoice): number {
		if (!inv.total_amount) return 0;
		return Math.min(100, ((inv.paid_amount || 0) / inv.total_amount) * 100);
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
						{#if viewMode === "detail" && selectedInvoice}
							Facture {selectedInvoice.invoice_number}
						{:else}
							Factures
						{/if}
					</h1>
					<p class="text-sm text-muted-foreground">
						{#if viewMode === "list"}
							Liste des factures du dossier
						{:else if client}
							Client : {client.name}
						{/if}
					</p>
				</div>
			</div>

			{#if viewMode === "list" && invoices.length > 0}
				<button
					type="button"
					onclick={showCreate}
					class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
				>
					Nouvelle facture
				</button>
			{/if}

			{#if viewMode === "detail" && selectedInvoice}
				<div class="flex items-center gap-2">
					{#if canSend(selectedInvoice)}
						<ConfirmDialog
							title="Envoyer la facture"
							description="La facture sera marquée comme envoyée et un email sera envoyé au client. Cette action est irréversible."
							confirmLabel="Envoyer"
							onConfirm={handleSend}
						>
							<button
								type="button"
								disabled={sendingInvoice}
								class="h-input rounded-input bg-accent text-accent-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<Send size={14} class="mr-1" />
								{sendingInvoice ? "Envoi..." : "Envoyer"}
							</button>
						</ConfirmDialog>
					{/if}
					{#if canPay(selectedInvoice) && !showPaymentForm}
						<button
							type="button"
							onclick={openPaymentForm}
							class="h-input rounded-input bg-success text-success-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer"
						>
							<Banknote size={14} class="mr-1" />
							Enregistrer un paiement
						</button>
					{/if}
					{#if canVoid(selectedInvoice)}
						<ConfirmDialog
							title="Annuler la facture"
							description="Attention : cette facture sera annulée définitivement. Cette action est irréversible."
							confirmLabel="Annuler la facture"
							onConfirm={handleVoid}
						>
							<button
								type="button"
								disabled={voidingInvoice}
								class="h-input rounded-input bg-destructive text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<Ban size={14} class="mr-1" />
								{voidingInvoice ? "Annulation..." : "Annuler la facture"}
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
			{#if invoices.length === 0}
				<div class="border border-dashed border-border rounded-lg p-12 bg-muted/20 flex flex-col items-center justify-center gap-4 flex-1 min-h-[50vh]">
					<FileText class="w-12 h-12 text-muted-foreground" />
					<p class="text-muted-foreground text-center">Aucune facture pour ce dossier.</p>
					<p class="text-sm text-muted-foreground text-center">
						Une facture peut être créée à partir d'un contrat actif, ou vous pouvez en créer une manuellement.
					</p>
					<button
						type="button"
						onclick={showCreate}
						class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
					>
						Créer une facture
					</button>
				</div>
			{:else}
				<div class="border border-border-card rounded-lg overflow-hidden">
					<table class="w-full">
						<thead class="bg-muted">
							<tr>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Référence</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Date d'émission</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Échéance</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Montant</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Paiement</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Statut</th>
							</tr>
						</thead>
						<tbody class="bg-background divide-y divide-border">
							{#each invoices as inv (inv.id)}
								<tr
									class="hover:bg-muted/50 transition-colors cursor-pointer"
									onclick={() => showDetail(inv)}
								>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm font-medium text-foreground">{inv.invoice_number}</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm text-foreground">{formatDate(inv.issue_date)}</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm text-foreground">
											{inv.due_date ? formatDate(inv.due_date) : "—"}
										</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm font-medium text-foreground">
											{formatCurrency(inv.total_amount, inv.currency)}
										</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {paymentStatusBadge(inv.payment_status as PaymentStatus)}">
											{PAYMENT_STATUS_LABELS[inv.payment_status] || inv.payment_status}
										</span>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {documentStatusBadge(inv.status as DocumentStatus)}">
											{STATUS_LABELS[inv.status] || inv.status}
										</span>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}

		<!-- ================================================================ -->
		<!-- DETAIL VIEW -->
		<!-- ================================================================ -->
		{:else if viewMode === "detail" && selectedInvoice}
			<div class="max-w-3xl animate-fade-in">
				<!-- Status banner -->
				<div class="flex items-center gap-3 mb-6">
					<span class="px-3 py-1.5 inline-flex text-sm leading-5 font-semibold rounded-full {documentStatusBadge(selectedInvoice.status as DocumentStatus)}">
						{STATUS_LABELS[selectedInvoice.status] || selectedInvoice.status}
					</span>
					<span class="px-3 py-1.5 inline-flex text-sm leading-5 font-semibold rounded-full {paymentStatusBadge(selectedInvoice.payment_status as PaymentStatus)}">
						{PAYMENT_STATUS_LABELS[selectedInvoice.payment_status] || selectedInvoice.payment_status}
					</span>
					{#if selectedInvoice.payment_status === "paid"}
						<span class="text-sm text-success">
							Facture réglée.
						</span>
					{:else if selectedInvoice.payment_status === "overdue"}
						<span class="text-sm text-destructive">
							Paiement en retard — l'échéance est dépassée.
						</span>
					{:else if selectedInvoice.payment_status === "partially_paid"}
						<span class="text-sm text-tertiary">
							Paiement partiel — reste {formatCurrency(remainingAmount(selectedInvoice), selectedInvoice.currency)}.
						</span>
					{:else if selectedInvoice.status === "sent"}
						<span class="text-sm text-accent-subtle-foreground">
							Envoyée au client — en attente de paiement.
						</span>
					{:else if selectedInvoice.status === "draft"}
						<span class="text-sm text-muted-foreground">
							Brouillon — à envoyer au client.
						</span>
					{:else if selectedInvoice.payment_status === "void"}
						<span class="text-sm text-destructive">
							Facture annulée.
						</span>
					{/if}
				</div>

				<!-- Invoice info card -->
				<div class="border border-border-card rounded-lg p-6 mb-6">
					<!-- Reference, dates, currency -->
					<div class="grid grid-cols-2 gap-4 mb-6">
						<div>
							<p class="text-xs text-muted-foreground mb-1">Référence</p>
							<p class="text-sm font-semibold text-foreground">{selectedInvoice.invoice_number}</p>
						</div>
						<div>
							<p class="text-xs text-muted-foreground mb-1">Devise</p>
							<p class="text-sm text-foreground">{selectedInvoice.currency || "EUR"}</p>
						</div>
						<div>
							<p class="text-xs text-muted-foreground mb-1">Date d'émission</p>
							<p class="text-sm text-foreground">{formatDate(selectedInvoice.issue_date)}</p>
						</div>
						<div>
							<p class="text-xs text-muted-foreground mb-1">Date d'échéance</p>
							<p class="text-sm text-foreground">
								{selectedInvoice.due_date ? formatDate(selectedInvoice.due_date) : "—"}
							</p>
						</div>
						{#if selectedInvoice.linked_contract_id}
							<div>
								<p class="text-xs text-muted-foreground mb-1">Contrat lié</p>
								<p class="text-sm text-foreground">{selectedInvoice.linked_contract_id}</p>
							</div>
						{/if}
						{#if selectedInvoice.payment_terms}
							<div>
								<p class="text-xs text-muted-foreground mb-1">Conditions de paiement</p>
								<p class="text-sm text-foreground">{selectedInvoice.payment_terms}</p>
							</div>
						{/if}
					</div>

					{#if selectedInvoice.notes}
						<div class="mb-4">
							<p class="text-xs text-muted-foreground mb-1">Notes</p>
							<p class="text-sm text-muted-foreground italic">{selectedInvoice.notes}</p>
						</div>
					{/if}

					<!-- Line items table -->
					{#if selectedInvoice.line_items && selectedInvoice.line_items.length > 0}
						<div class="border border-border rounded-lg overflow-hidden mb-6">
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
									{#each selectedInvoice.line_items as item, i (i)}
										<tr>
											<td class="px-4 py-3 text-sm text-foreground">{item.description}</td>
											<td class="px-4 py-3 text-sm text-foreground text-right">{item.quantity}</td>
											<td class="px-4 py-3 text-sm text-foreground text-right">{formatCurrency(item.unit_price, selectedInvoice.currency)}</td>
											<td class="px-4 py-3 text-sm font-medium text-foreground text-right">{formatCurrency(item.subtotal, selectedInvoice.currency)}</td>
										</tr>
									{/each}
								</tbody>
								<tfoot class="bg-muted/50 border-t border-border">
									{#if selectedInvoice.tax_rate > 0}
										<tr>
											<td colspan="3" class="px-4 py-2 text-sm text-muted-foreground text-right">Sous-total</td>
											<td class="px-4 py-2 text-sm font-medium text-foreground text-right">
												{formatCurrency(selectedInvoice.total_amount - selectedInvoice.tax_amount, selectedInvoice.currency)}
											</td>
										</tr>
										<tr>
											<td colspan="3" class="px-4 py-2 text-sm text-muted-foreground text-right">Taxe ({selectedInvoice.tax_rate}%)</td>
											<td class="px-4 py-2 text-sm font-medium text-foreground text-right">
												{formatCurrency(selectedInvoice.tax_amount, selectedInvoice.currency)}
											</td>
										</tr>
									{/if}
									<tr>
										<td colspan="3" class="px-4 py-3 text-sm font-semibold text-foreground text-right">Total</td>
										<td class="px-4 py-3 text-lg font-bold text-foreground text-right">
											{formatCurrency(selectedInvoice.total_amount, selectedInvoice.currency)}
										</td>
									</tr>
								</tfoot>
							</table>
						</div>
					{/if}

					<!-- Payment summary -->
					<div class="border border-border rounded-lg p-4">
						<p class="text-xs text-muted-foreground mb-3 font-medium uppercase tracking-wider">Paiement</p>

						<!-- Payment progress bar -->
						<div class="mb-3">
							<div class="flex justify-between text-sm mb-1">
								<span class="text-muted-foreground">Payé</span>
								<span class="font-medium text-foreground">
									{formatCurrency(selectedInvoice.paid_amount || 0, selectedInvoice.currency)} / {formatCurrency(selectedInvoice.total_amount, selectedInvoice.currency)}
								</span>
							</div>
							<div class="w-full bg-muted rounded-full h-2.5">
								<div
									class="h-2.5 rounded-full transition-interactive duration-300 {selectedInvoice.payment_status === 'paid' ? 'bg-success' : selectedInvoice.payment_status === 'overdue' ? 'bg-destructive' : 'bg-accent'}"
									style="width: {paymentProgress(selectedInvoice)}%"
								></div>
							</div>
						</div>

						<div class="grid grid-cols-2 gap-4">
							{#if selectedInvoice.paid_at}
								<div>
									<p class="text-xs text-muted-foreground mb-1">Date de paiement</p>
									<p class="text-sm text-foreground">{formatDate(selectedInvoice.paid_at)}</p>
								</div>
							{/if}
							{#if selectedInvoice.payment_method}
								<div>
									<p class="text-xs text-muted-foreground mb-1">Mode de paiement</p>
									<p class="text-sm text-foreground">{PAYMENT_METHODS[selectedInvoice.payment_method] || selectedInvoice.payment_method}</p>
								</div>
							{/if}
							{#if selectedInvoice.payment_status !== "paid" && selectedInvoice.payment_status !== "void"}
								<div>
									<p class="text-xs text-muted-foreground mb-1">Reste à payer</p>
									<p class="text-sm font-semibold text-foreground">
										{formatCurrency(remainingAmount(selectedInvoice), selectedInvoice.currency)}
									</p>
								</div>
							{/if}
							{#if selectedInvoice.late_fee}
								<div>
									<p class="text-xs text-muted-foreground mb-1">Pénalité de retard</p>
									<p class="text-sm text-destructive">{formatCurrency(selectedInvoice.late_fee, selectedInvoice.currency)}</p>
								</div>
							{/if}
						</div>
					</div>
				</div>

				<!-- PAYMENT FORM -->
				{#if showPaymentForm}
					<div class="border border-success/30 bg-success/10 rounded-lg p-6 mb-6 animate-fade-in">
						<h3 class="text-sm font-semibold text-foreground mb-4">Enregistrer un paiement</h3>

						<div class="grid grid-cols-2 gap-4 mb-4">
							<div>
								<label class="text-xs text-muted-foreground mb-1 block">Montant ({selectedInvoice.currency || "EUR"}) *</label>
								<input
									type="number"
									bind:value={paymentAmount}
									min="0.01"
									step="0.01"
									placeholder="0.00"
									class="h-input rounded-input border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
								/>
								<p class="text-xs text-muted-foreground mt-1">
									Reste : {formatCurrency(remainingAmount(selectedInvoice), selectedInvoice.currency)}
								</p>
							</div>
							<div>
								<label class="text-xs text-muted-foreground mb-1 block">Mode de paiement *</label>
								<select
									bind:value={paymentMethod}
									class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
								>
									{#each Object.entries(PAYMENT_METHODS) as [value, label]}
										<option value={value}>{label}</option>
									{/each}
								</select>
							</div>
						</div>

						<div class="flex justify-end gap-2">
							<button
								type="button"
								onclick={cancelPayment}
								class="h-input rounded-input bg-transparent text-dark hover:bg-muted inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] border-2 border-border-input cursor-pointer"
							>
								<X size={14} class="mr-1" />
								Annuler
							</button>
							<button
								type="button"
								onclick={handlePayment}
								disabled={paymentSubmitting}
								class="h-input rounded-input bg-success text-success-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<Save size={14} class="mr-1" />
								{paymentSubmitting ? "Enregistrement..." : "Enregistrer le paiement"}
							</button>
						</div>
					</div>
				{/if}

				<!-- Lifecycle info for non-actionable states -->
				{#if hasNoActions(selectedInvoice) && !showPaymentForm}
					<div class="flex items-start gap-2 text-sm text-muted-foreground bg-muted/30 border border-border rounded-lg p-4">
						<AlertTriangle size={16} class="flex-shrink-0 mt-0.5" />
						<p>
							Cette facture est dans l'état <strong>{STATUS_LABELS[selectedInvoice.status]}</strong>
							({PAYMENT_STATUS_LABELS[selectedInvoice.payment_status] || selectedInvoice.payment_status}).
							{#if selectedInvoice.payment_status === "paid"}
								La facture a été réglée intégralement.
							{:else if selectedInvoice.status === "archived"}
								Cette facture a été archivée et n'est plus modifiable.
							{:else if selectedInvoice.payment_status === "void"}
								Cette facture a été annulée.
							{:else}
								Aucune action n'est disponible pour le moment.
							{/if}
						</p>
					</div>
				{/if}
			</div>
		{/if}
	{/if}
</div>

<!-- ================================================================ -->
<!-- CREATE MODAL -->
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
			<div class="flex items-center justify-between px-6 py-4 border-b border-border shrink-0">
				<Dialog.Title class="text-lg font-semibold text-foreground">
					Nouvelle facture
				</Dialog.Title>
				<Dialog.Close
					class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
				>
					<X class="text-foreground size-5" />
					<span class="sr-only">Fermer</span>
				</Dialog.Close>
			</div>

			<!-- Modal body (scrollable) -->
			<div class="px-6 py-5 overflow-y-auto flex-1 space-y-6">
				<!-- Dates -->
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Date d'émission *</label>
						<input
							type="date"
							bind:value={formIssueDate}
							class="h-input rounded-input border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Date d'échéance *</label>
						<input
							type="date"
							bind:value={formDueDate}
							class="h-input rounded-input border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
				</div>

				<!-- Line items -->
				<div>
					<div class="flex items-center justify-between mb-3">
						<label class="text-sm font-medium text-foreground">Lignes de facture</label>
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

					<div class="flex items-end justify-end gap-6 mt-4 pt-4 border-t border-border">
						<div class="flex items-center gap-3">
							<label class="text-xs text-muted-foreground">TVA (%)</label>
							<input
								type="number"
								bind:value={formTaxRate}
								min="0"
								max="100"
								step="0.1"
								class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-20 px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
						<div class="text-right">
							<p class="text-sm text-muted-foreground">Total TTC</p>
							<p class="text-2xl font-bold text-foreground">{formatCurrency(formTotal())}</p>
						</div>
					</div>
				</div>

				<!-- Extra fields -->
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Conditions de paiement</label>
						<input
							type="text"
							bind:value={formPaymentTerms}
							placeholder="Ex : Paiement à 30 jours"
							class="h-input rounded-input border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Notes</label>
						<input
							type="text"
							bind:value={formNotes}
							placeholder="Notes internes..."
							class="h-input rounded-input border border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
				</div>
			</div>

			<!-- Modal footer -->
			<div class="flex justify-end gap-2 px-6 py-4 border-t border-border shrink-0">
				<Dialog.Close
					class="h-input rounded-input bg-transparent text-dark hover:bg-muted inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] border-2 border-border-input cursor-pointer"
				>
					Annuler
				</Dialog.Close>
				<button
					type="button"
					onclick={handleCreate}
					disabled={formSaving}
					class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
				>
					<Save size={14} class="mr-1" />
					{formSaving ? "Enregistrement..." : "Créer la facture"}
				</button>
			</div>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>

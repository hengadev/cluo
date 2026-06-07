<script lang="ts">
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { goto } from "$app/navigation";
	import {
		Send,
		CheckCircle,
		ShieldCheck,
		Receipt,
		ChevronLeft,
		AlertTriangle,
		FileText,
		X,
		Save,
	} from "@lucide/svelte";
	import { Dialog } from "bits-ui";
	import {
		fetchCase,
		fetchClient,
		fetchCaseContracts,
		createContract,
		sendDocument,
		signContract,
		activateContract,
		createInvoiceFromContract,
		ConflictError,
	} from "$lib/services/api";
	import { currentCase } from "$lib/stores/case";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import type { Case, Client, Contract } from "$lib/types/entities";

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
	let contracts: Contract[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	// Currently selected contract for viewing
	let selectedContract: Contract | null = $state(null);
	let viewMode: "list" | "detail" = $state("list");
	let showCreateModal = $state(false);

	// Create form state
	let formStartDate = $state(todayISO());
	let formEndDate = $state("");
	let formScopeOfServices = $state("");
	let formPaymentTerms = $state("");
	let formConfidentiality = $state("");
	let formTerminationClause = $state("");
	let formContractValue = $state("");
	let formCurrency = $state("EUR");
	let formSaving = $state(false);

	// Lifecycle action state
	let sendingContract = $state(false);
	let signingContract = $state(false);
	let activatingContract = $state(false);
	let creatingInvoice = $state(false);

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

	const STATUS_COLORS: Record<string, string> = {
		draft: "bg-muted text-muted-foreground",
		sent: "bg-accent text-accent-foreground",
		signed: "bg-success/15 text-success",
		active: "bg-success text-success-foreground",
		archived: "bg-muted text-muted-foreground",
		cancelled: "bg-destructive/15 text-destructive",
		rejected: "bg-destructive/15 text-destructive",
		expired: "bg-tertiary/15 text-tertiary",
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

	function canSend(c: Contract): boolean {
		return c.status === "draft";
	}

	function canSign(c: Contract): boolean {
		return c.status === "sent";
	}

	function canActivate(c: Contract): boolean {
		return c.status === "signed";
	}

	function canCreateInvoice(c: Contract): boolean {
		return c.status === "active";
	}

	function hasNoActions(c: Contract): boolean {
		return !canSend(c) && !canSign(c) && !canActivate(c) && !canCreateInvoice(c);
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
			const [c, conts] = await Promise.all([
				fetchCase(caseId),
				fetchCaseContracts(caseId),
			]);
			caseData = c;
			contracts = conts;

			if (c?.clientId) {
				client = await fetchClient(c.clientId);
			}
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Erreur lors du chargement des contrats";
		} finally {
			loading = false;
		}
	}

	async function refreshContracts() {
		const id = selectedContract?.id;
		if (!id) return;
		const freshContracts = await fetchCaseContracts(caseId);
		contracts = freshContracts;
		const updated = freshContracts.find((c) => c.id === id);
		if (updated) selectedContract = updated;
	}

	// =========================================================================
	// Navigation helpers
	// =========================================================================

	function showList() {
		selectedContract = null;
		viewMode = "list";
	}

	function showCreate() {
		formScopeOfServices = "";
		formPaymentTerms = "";
		formConfidentiality = "";
		formTerminationClause = "";
		formStartDate = todayISO();
		formEndDate = "";
		formContractValue = "";
		formCurrency = "EUR";
		showCreateModal = true;
	}

	function showDetail(c: Contract) {
		selectedContract = c;
		viewMode = "detail";
	}

	async function handleCreate() {
		if (!caseData) return;
		if (!formScopeOfServices.trim() || !formPaymentTerms.trim() || !formConfidentiality.trim() || !formTerminationClause.trim() || !formStartDate) {
			toastState.add(TOAST_LEVELS.Error, "Erreur", "Veuillez remplir tous les champs obligatoires.");
			return;
		}

		formSaving = true;
		try {
			const contractNumber = `CTR-${new Date().getFullYear()}-${String(contracts.length + 1).padStart(3, "0")}`;
			const payload = {
				case_id: caseData.id,
				client_id: caseData.clientId,
				contract_number: contractNumber,
				start_date: new Date(formStartDate).toISOString(),
				end_date: formEndDate ? new Date(formEndDate).toISOString() : undefined,
				scope_of_services: formScopeOfServices.trim(),
				payment_terms: formPaymentTerms.trim(),
				confidentiality: formConfidentiality.trim(),
				termination_clause: formTerminationClause.trim(),
				contract_value: formContractValue ? parseFloat(formContractValue) : undefined,
				currency: formCurrency || "EUR",
				signatures: [],
				status: "draft" as const,
			} as Contract;

			const result = await createContract(payload);
			if (result.data) {
				contracts = [...contracts, result.data];
				selectedContract = result.data;
				showCreateModal = false;
				viewMode = "detail";
				toastState.add(TOAST_LEVELS.Info, "Contrat créé", "Le contrat a été créé en brouillon.");
			}
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de créer le contrat.",
			);
		} finally {
			formSaving = false;
		}
	}

	// =========================================================================
	// Send Contract (draft → sent)
	// =========================================================================

	async function handleSend() {
		if (!selectedContract || !canSend(selectedContract)) return;
		sendingContract = true;
		try {
			const result = await sendDocument(selectedContract.id, "contract", {
				recipients: [],
				subject: `Contrat ${selectedContract.contract_number}`,
				message: "",
				send_email: true,
				send_sms: false,
			});

			if (result.success) {
				await refreshContracts();
				toastState.add(
					TOAST_LEVELS.Info,
					"Contrat envoyé",
					"Le contrat a été marqué comme envoyé.",
				);
			}
		} catch (e) {
			const msg = e instanceof ConflictError
				? e.message
				: e instanceof Error
					? e.message
					: "Impossible d'envoyer le contrat.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			sendingContract = false;
		}
	}

	// =========================================================================
	// Sign Contract (sent → signed)
	// =========================================================================

	async function handleSign() {
		if (!selectedContract || !canSign(selectedContract)) return;
		signingContract = true;
		try {
			const result = await signContract(selectedContract.id, {
				signer_name: "Enquêteur",
				signer_role: "investigator",
				method: "e-sign",
			});

			if (result.success) {
				await refreshContracts();
				toastState.add(
					TOAST_LEVELS.Info,
					"Contrat signé",
					"La signature a été enregistrée.",
				);
			}
		} catch (e) {
			const msg = e instanceof ConflictError
				? e.message
				: e instanceof Error
					? e.message
					: "Impossible de signer le contrat.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			signingContract = false;
		}
	}

	// =========================================================================
	// Activate Contract (signed → active)
	// =========================================================================

	async function handleActivate() {
		if (!selectedContract || !canActivate(selectedContract)) return;
		activatingContract = true;
		try {
			const result = await activateContract(selectedContract.id);
			if (result.success) {
				await refreshContracts();
				toastState.add(
					TOAST_LEVELS.Info,
					"Contrat activé",
					"Le contrat est maintenant actif. L'accord commercial est en vigueur.",
				);
			}
		} catch (e) {
			const msg = e instanceof ConflictError
				? e.message
				: e instanceof Error
					? e.message
					: "Impossible d'activer le contrat.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			activatingContract = false;
		}
	}

	// =========================================================================
	// Create Invoice (from active contract)
	// =========================================================================

	async function handleCreateInvoice() {
		if (!selectedContract || !canCreateInvoice(selectedContract)) return;
		creatingInvoice = true;
		try {
			const result = await createInvoiceFromContract(selectedContract.id);
			if (result.success) {
				toastState.add(
					TOAST_LEVELS.Info,
					"Facture créée",
					"Une facture a été générée à partir du contrat.",
				);
				await goto(`/cases/${caseId}/documents/facture`);
			}
		} catch (e) {
			const msg = e instanceof ConflictError
				? e.message
				: e instanceof Error
					? e.message
					: "Impossible de créer la facture.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			creatingInvoice = false;
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
						{#if viewMode === "detail" && selectedContract}
							Contrat {selectedContract.contract_number}
						{:else}
							Contrats
						{/if}
					</h1>
					<p class="text-sm text-muted-foreground">
						{#if viewMode === "list"}
							Liste des contrats du dossier
						{:else if client}
							Client : {client.name}
						{/if}
					</p>
				</div>
			</div>

			{#if viewMode === "list" && contracts.length > 0}
				<button
					type="button"
					onclick={showCreate}
					class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
				>
					Nouveau contrat
				</button>
			{/if}

			{#if viewMode === "detail" && selectedContract}
				<div class="flex items-center gap-2">
					{#if canSend(selectedContract)}
						<ConfirmDialog
							title="Envoyer le contrat"
							description="Le contrat sera marqué comme envoyé au client. Cette action est irréversible."
							confirmLabel="Envoyer"
							onConfirm={handleSend}
						>
							<button
								type="button"
								disabled={sendingContract}
								class="h-input rounded-input bg-accent text-accent-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<Send size={14} class="mr-1" />
								{sendingContract ? "Envoi..." : "Envoyer"}
							</button>
						</ConfirmDialog>
					{/if}
					{#if canSign(selectedContract)}
						<ConfirmDialog
							title="Signer le contrat"
							description="Enregistrez la signature du contrat. Le contrat passera en état « Signé »."
							confirmLabel="Signer"
							onConfirm={handleSign}
						>
							<button
								type="button"
								disabled={signingContract}
								class="h-input rounded-input bg-success text-success-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<CheckCircle size={14} class="mr-1" />
								{signingContract ? "Signature..." : "Signer"}
							</button>
						</ConfirmDialog>
					{/if}
					{#if canActivate(selectedContract)}
						<ConfirmDialog
							title="Activer le contrat"
							description="Le contrat sera activé, mettant en vigueur l'accord commercial. Cette action est irréversible."
							confirmLabel="Activer"
							onConfirm={handleActivate}
						>
							<button
								type="button"
								disabled={activatingContract}
								class="h-input rounded-input bg-success text-success-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<ShieldCheck size={14} class="mr-1" />
								{activatingContract ? "Activation..." : "Activer"}
							</button>
						</ConfirmDialog>
					{/if}
					{#if canCreateInvoice(selectedContract)}
						<ConfirmDialog
							title="Créer une facture"
							description="Une facture sera générée à partir de ce contrat actif et vous serez redirigé vers la liste des factures."
							confirmLabel="Créer la facture"
							onConfirm={handleCreateInvoice}
						>
							<button
								type="button"
								disabled={creatingInvoice}
								class="h-input rounded-input bg-tertiary text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<Receipt size={14} class="mr-1" />
								{creatingInvoice ? "Création..." : "Créer une facture"}
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
			{#if contracts.length === 0}
				<div class="border border-dashed border-border rounded-lg p-12 bg-muted/20 flex flex-col items-center justify-center gap-4 flex-1 min-h-[50vh]">
					<FileText class="w-12 h-12 text-muted-foreground" />
					<p class="text-muted-foreground text-center">Aucun contrat pour ce dossier.</p>
					<p class="text-sm text-muted-foreground text-center">
						Un contrat peut être créé à partir d'un mandat activé, ou vous pouvez en créer un manuellement.
					</p>
					<button
						type="button"
						onclick={showCreate}
						class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
					>
						Créer un contrat
					</button>
				</div>
			{:else}
				<div class="border border-border-card rounded-lg overflow-hidden">
					<table class="w-full">
						<thead class="bg-muted">
							<tr>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Référence</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Début</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Fin</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Montant</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Statut</th>
							</tr>
						</thead>
						<tbody class="bg-background divide-y divide-border">
							{#each contracts as c (c.id)}
								<tr
									class="hover:bg-muted/50 transition-colors cursor-pointer"
									onclick={() => showDetail(c)}
								>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm font-medium text-foreground">{c.contract_number}</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm text-foreground">{formatDate(c.start_date)}</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm text-foreground">
											{c.end_date ? formatDate(c.end_date) : "—"}
										</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										{#if c.contract_value}
											<div class="text-sm font-medium text-foreground">
												{formatCurrency(c.contract_value, c.currency)}
											</div>
										{:else}
											<div class="text-sm text-muted-foreground">—</div>
										{/if}
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {STATUS_COLORS[c.status] || 'bg-muted text-muted-foreground'}">
											{STATUS_LABELS[c.status] || c.status}
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
		{:else if viewMode === "detail" && selectedContract}
			<div class="max-w-3xl animate-fade-in">
				<!-- Status banner -->
				<div class="flex items-center gap-3 mb-6">
					<span class="px-3 py-1.5 inline-flex text-sm leading-5 font-semibold rounded-full {STATUS_COLORS[selectedContract.status] || 'bg-muted text-muted-foreground'}">
						{STATUS_LABELS[selectedContract.status] || selectedContract.status}
					</span>
					{#if selectedContract.status === "active"}
						<span class="text-sm text-success">
							Accord commercial en vigueur.
						</span>
					{:else if selectedContract.status === "signed"}
						<span class="text-sm text-success">
							Signé — en attente d'activation.
						</span>
					{:else if selectedContract.status === "sent"}
						<span class="text-sm text-accent-foreground">
							Envoyé au client — en attente de signature.
						</span>
					{:else if selectedContract.status === "draft"}
						<span class="text-sm text-muted-foreground">
							Brouillon — à envoyer pour signature.
						</span>
					{:else if selectedContract.status === "archived"}
						<span class="text-sm text-muted-foreground">
							Archivé — ce contrat n'est plus actif.
						</span>
					{/if}
				</div>

				<!-- Contract info card -->
				<div class="border border-border-card rounded-lg p-6 mb-6">
					<!-- Reference and dates -->
					<div class="grid grid-cols-2 gap-4 mb-6">
						<div>
							<p class="text-xs text-muted-foreground mb-1">Référence</p>
							<p class="text-sm font-semibold text-foreground">{selectedContract.contract_number}</p>
						</div>
						{#if selectedContract.currency}
							<div>
								<p class="text-xs text-muted-foreground mb-1">Devise</p>
								<p class="text-sm text-foreground">{selectedContract.currency}</p>
							</div>
						{/if}
						<div>
							<p class="text-xs text-muted-foreground mb-1">Date de début</p>
							<p class="text-sm text-foreground">{formatDate(selectedContract.start_date)}</p>
						</div>
						<div>
							<p class="text-xs text-muted-foreground mb-1">Date de fin</p>
							<p class="text-sm text-foreground">
								{selectedContract.end_date ? formatDate(selectedContract.end_date) : "—"}
							</p>
						</div>
						{#if selectedContract.contract_value}
							<div>
								<p class="text-xs text-muted-foreground mb-1">Montant du contrat</p>
								<p class="text-lg font-bold text-foreground">
									{formatCurrency(selectedContract.contract_value, selectedContract.currency)}
								</p>
							</div>
						{/if}
						{#if selectedContract.linked_mandate_id}
							<div>
								<p class="text-xs text-muted-foreground mb-1">Mandat lié</p>
								<p class="text-sm text-foreground">{selectedContract.linked_mandate_id}</p>
							</div>
						{/if}
					</div>

					<!-- Scope of services -->
					<div class="mb-4">
						<p class="text-xs text-muted-foreground mb-1">Objet des prestations</p>
						<p class="text-sm text-foreground">{selectedContract.scope_of_services}</p>
					</div>

					<!-- Payment terms -->
					<div class="mb-4">
						<p class="text-xs text-muted-foreground mb-1">Conditions de paiement</p>
						<p class="text-sm text-foreground">{selectedContract.payment_terms}</p>
					</div>

					<!-- Confidentiality -->
					<div class="mb-4">
						<p class="text-xs text-muted-foreground mb-1">Clause de confidentialité</p>
						<p class="text-sm text-muted-foreground">{selectedContract.confidentiality}</p>
					</div>

					<!-- Termination clause -->
					<div class="mb-4">
						<p class="text-xs text-muted-foreground mb-1">Clause de résiliation</p>
						<p class="text-sm text-muted-foreground">{selectedContract.termination_clause}</p>
					</div>

					{#if selectedContract.renewal_terms}
						<div class="mb-4">
							<p class="text-xs text-muted-foreground mb-1">Conditions de renouvellement</p>
							<p class="text-sm text-muted-foreground">{selectedContract.renewal_terms}</p>
						</div>
					{/if}

					{#if selectedContract.governing_law}
						<div class="mb-4">
							<p class="text-xs text-muted-foreground mb-1">Droit applicable</p>
							<p class="text-sm text-muted-foreground">{selectedContract.governing_law}</p>
						</div>
					{/if}

					<!-- Signatures section -->
					<div class="mt-6 pt-4 border-t border-border">
						<p class="text-xs text-muted-foreground mb-3 font-medium uppercase tracking-wider">Signatures</p>
						{#if selectedContract.signatures && selectedContract.signatures.length > 0}
							<div class="grid grid-cols-2 gap-4">
								{#each selectedContract.signatures as sig, i (i)}
									<div class="border border-border rounded-lg p-4">
										<p class="text-xs text-muted-foreground mb-2">
											{sig.role === 'investigator' ? 'Enquêteur' : sig.role === 'client' ? 'Client' : `Signataire ${i + 1}`}
										</p>
										<div class="flex items-center gap-2">
											<CheckCircle size={16} class="text-success flex-shrink-0" />
											<div>
												<p class="text-sm font-medium text-foreground">{sig.name}</p>
												<p class="text-xs text-muted-foreground">
													Signé le {formatDate(sig.signed_at)}
												</p>
											</div>
										</div>
									</div>
								{/each}
							</div>
						{:else}
							<p class="text-sm text-muted-foreground">Aucune signature enregistrée.</p>
						{/if}
					</div>
				</div>

				<!-- Lifecycle info for non-actionable states -->
				{#if hasNoActions(selectedContract)}
					<div class="flex items-start gap-2 text-sm text-muted-foreground bg-muted/30 border border-border rounded-lg p-4">
						<AlertTriangle size={16} class="flex-shrink-0 mt-0.5" />
						<p>
							Ce contrat est dans l'état <strong>{STATUS_LABELS[selectedContract.status]}</strong>.
							{#if selectedContract.status === "archived"}
								Ce contrat a été archivé et n'est plus modifiable.
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
					Nouveau contrat
				</Dialog.Title>
				<Dialog.Close
					class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
				>
					<X class="text-foreground size-5" />
					<span class="sr-only">Fermer</span>
				</Dialog.Close>
			</div>

			<!-- Modal body (scrollable) -->
			<div class="px-6 py-5 overflow-y-auto flex-1 space-y-4">
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Date de début *</label>
						<input
							type="date"
							bind:value={formStartDate}
							class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Date de fin</label>
						<input
							type="date"
							bind:value={formEndDate}
							class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Montant du contrat (€)</label>
						<input
							type="number"
							bind:value={formContractValue}
							min="0"
							step="0.01"
							placeholder="0.00"
							class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Devise</label>
						<input
							type="text"
							bind:value={formCurrency}
							placeholder="EUR"
							class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
				</div>

				<div>
					<label class="text-xs text-muted-foreground mb-1 block">Objet des prestations *</label>
					<textarea
						bind:value={formScopeOfServices}
						placeholder="Décrivez l'objet des prestations..."
						rows="3"
						class="rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 py-2 text-sm focus:ring-2 focus:ring-offset-2 resize-none"
					></textarea>
				</div>

				<div>
					<label class="text-xs text-muted-foreground mb-1 block">Conditions de paiement *</label>
					<textarea
						bind:value={formPaymentTerms}
						placeholder="Ex : Paiement à 30 jours..."
						rows="2"
						class="rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 py-2 text-sm focus:ring-2 focus:ring-offset-2 resize-none"
					></textarea>
				</div>

				<div>
					<label class="text-xs text-muted-foreground mb-1 block">Clause de confidentialité *</label>
					<textarea
						bind:value={formConfidentiality}
						placeholder="Clause de confidentialité..."
						rows="2"
						class="rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 py-2 text-sm focus:ring-2 focus:ring-offset-2 resize-none"
					></textarea>
				</div>

				<div>
					<label class="text-xs text-muted-foreground mb-1 block">Clause de résiliation *</label>
					<textarea
						bind:value={formTerminationClause}
						placeholder="Conditions de résiliation..."
						rows="2"
						class="rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 py-2 text-sm focus:ring-2 focus:ring-offset-2 resize-none"
					></textarea>
				</div>
			</div>

			<!-- Modal footer -->
			<div class="flex justify-end gap-2 px-6 py-4 border-t border-border shrink-0">
				<Dialog.Close
					class="h-input rounded-input bg-transparent text-dark hover:bg-[#fafafa] inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] border-2 border-[#dedede] cursor-pointer"
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
					{formSaving ? "Enregistrement..." : "Créer le contrat"}
				</button>
			</div>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>

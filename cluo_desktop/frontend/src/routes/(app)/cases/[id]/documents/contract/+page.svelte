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
		Printer,
		Pencil,
		Trash2,
	} from "@lucide/svelte";
	import { Dialog } from "bits-ui";
	import {
		fetchCase,
		fetchClient,
		fetchCaseContracts,
		createContract,
		updateDocument,
		deleteDocument,
		sendDocument,
		signContract,
		activateContract,
		createInvoiceFromContract,
		openDocumentPDF,
		ConflictError,
	} from "$lib/services/api";
	import { currentCase } from "$lib/stores/case";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import DocumentLifecycleStepper from "$lib/custom/documents/DocumentLifecycleStepper.svelte";
	import { documentStatusBadge } from "$lib/utils/badgeVariants";
	import type { Case, Client, Contract, DocumentStatus } from "$lib/types/entities";

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
	let showEditModal = $state(false);

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
	let previewingContractId: string | null = $state(null);
	let deletingContractId: string | null = $state(null);

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

	const CONTRACT_STEPS = [
		{ key: "draft", label: "Brouillon" },
		{ key: "sent", label: "Envoyé" },
		{ key: "signed", label: "Signé" },
		{ key: "active", label: "Actif" },
	];

	function statusNote(c: Contract): string {
		switch (c.status) {
			case "active":
				return "Accord commercial en vigueur.";
			case "signed":
				return "Signé — en attente d'activation.";
			case "sent":
				return "Envoyé au client — en attente de signature.";
			case "draft":
				return "Brouillon — à envoyer pour signature.";
			case "archived":
				return "Archivé — ce contrat n'est plus actif.";
			default:
				return "";
		}
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

	function canEdit(c: Contract): boolean {
		return c.status === "draft";
	}

	function canDelete(c: Contract): boolean {
		return c.status === "draft";
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

	async function handlePreview(c: Contract) {
		previewingContractId = c.id;
		try {
			await openDocumentPDF(c.id, "contract");
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible d'afficher l'aperçu du contrat.",
			);
		} finally {
			previewingContractId = null;
		}
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
	// Edit Contract (draft only)
	// =========================================================================

	function showEdit(c: Contract) {
		if (!canEdit(c)) return;
		selectedContract = c;
		formStartDate = c.start_date ? c.start_date.split("T")[0] : todayISO();
		formEndDate = c.end_date ? c.end_date.split("T")[0] : "";
		formScopeOfServices = c.scope_of_services || "";
		formPaymentTerms = c.payment_terms || "";
		formConfidentiality = c.confidentiality || "";
		formTerminationClause = c.termination_clause || "";
		formContractValue = c.contract_value != null ? String(c.contract_value) : "";
		formCurrency = c.currency || "EUR";
		showEditModal = true;
	}

	async function handleEditSave() {
		if (!selectedContract || !canEdit(selectedContract)) return;
		if (!formScopeOfServices.trim() || !formPaymentTerms.trim() || !formConfidentiality.trim() || !formTerminationClause.trim() || !formStartDate) {
			toastState.add(TOAST_LEVELS.Error, "Erreur", "Veuillez remplir tous les champs obligatoires.");
			return;
		}

		formSaving = true;
		try {
			const data = {
				start_date: new Date(formStartDate).toISOString(),
				end_date: formEndDate ? new Date(formEndDate).toISOString() : undefined,
				scope_of_services: formScopeOfServices.trim(),
				payment_terms: formPaymentTerms.trim(),
				confidentiality: formConfidentiality.trim(),
				termination_clause: formTerminationClause.trim(),
				contract_value: formContractValue ? parseFloat(formContractValue) : undefined,
				currency: formCurrency || "EUR",
			};

			const result = await updateDocument(selectedContract.id, "contract", { type: "contract", data });
			if (result.data) {
				const updated = { ...selectedContract, ...data } as Contract;
				contracts = contracts.map((c) => (c.id === updated.id ? updated : c));
				selectedContract = updated;
				showEditModal = false;
				toastState.add(TOAST_LEVELS.Info, "Contrat modifié", "Les modifications ont été enregistrées.");
			}
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de modifier le contrat.",
			);
		} finally {
			formSaving = false;
		}
	}

	// =========================================================================
	// Delete Contract (draft only)
	// =========================================================================

	async function handleDelete(c: Contract) {
		if (!canDelete(c)) return;
		deletingContractId = c.id;
		try {
			await deleteDocument(c.id, "contract");
			contracts = contracts.filter((x) => x.id !== c.id);
			if (selectedContract?.id === c.id) showList();
			toastState.add(TOAST_LEVELS.Info, "Contrat supprimé", "Le contrat a été supprimé.");
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de supprimer le contrat.",
			);
		} finally {
			deletingContractId = null;
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
					<div class="flex items-center gap-1">
						<button
							type="button"
							onclick={() => selectedContract && handlePreview(selectedContract)}
							disabled={previewingContractId === selectedContract.id}
							class="h-input rounded-input bg-transparent text-foreground hover:bg-muted inline-flex items-center justify-center gap-1.5 px-3 text-sm font-medium active:scale-[0.98] border border-border-input cursor-pointer disabled:opacity-50 transition-interactive duration-150"
						>
							<Printer size={14} />
							Aperçu
						</button>
						{#if canEdit(selectedContract)}
							<button
								type="button"
								onclick={() => selectedContract && showEdit(selectedContract)}
								class="h-input rounded-input bg-transparent text-foreground hover:bg-muted inline-flex items-center justify-center gap-1.5 px-3 text-sm font-medium active:scale-[0.98] border border-border-input cursor-pointer transition-interactive duration-150"
							>
								<Pencil size={14} />
								Modifier
							</button>
						{/if}
					</div>

					{#if canDelete(selectedContract) || canSend(selectedContract) || canSign(selectedContract) || canActivate(selectedContract) || canCreateInvoice(selectedContract)}
						<div class="w-px h-5 bg-border"></div>
					{/if}

					{#if canDelete(selectedContract)}
						<ConfirmDialog
							title="Supprimer le contrat"
							description="Le contrat sera définitivement supprimé. Cette action est irréversible."
							confirmLabel="Supprimer"
							onConfirm={() => { if (selectedContract) return handleDelete(selectedContract); }}
						>
							<button
								type="button"
								disabled={deletingContractId === selectedContract.id}
								class="h-input rounded-input inline-flex items-center justify-center gap-1.5 px-3 text-sm font-medium btn-ghost-destructive cursor-pointer disabled:opacity-50 transition-interactive duration-150"
							>
								<Trash2 size={14} />
								Supprimer
							</button>
						</ConfirmDialog>
					{/if}
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
								class="h-input rounded-input bg-accent text-accent-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-1.5 px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<Send size={14} />
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
								class="h-input rounded-input bg-success text-success-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-1.5 px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<CheckCircle size={14} />
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
								class="h-input rounded-input bg-success text-success-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-1.5 px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<ShieldCheck size={14} />
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
								class="h-input rounded-input bg-tertiary text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-1.5 px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<Receipt size={14} />
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
								<th class="px-6 py-3 text-right text-xs font-medium text-muted-foreground uppercase tracking-wider">Montant</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Statut</th>
								<th class="px-6 py-3 w-12"></th>
							</tr>
						</thead>
						<tbody class="bg-background divide-y divide-border">
							{#each contracts as c (c.id)}
								<tr
									class="hover:shadow-mini hover:relative transition-interactive duration-150 cursor-pointer"
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
									<td class="px-6 py-4 whitespace-nowrap text-right">
										{#if c.contract_value}
											<div class="text-sm font-medium text-foreground">
												{formatCurrency(c.contract_value, c.currency)}
											</div>
										{:else}
											<div class="text-sm text-muted-foreground">—</div>
										{/if}
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {documentStatusBadge(c.status as DocumentStatus)}">
											{STATUS_LABELS[c.status] || c.status}
										</span>
									</td>
									<td class="px-6 py-4 whitespace-nowrap text-right">
										<div class="flex items-center justify-end gap-1">
											<button
												type="button"
												onclick={(e) => { e.stopPropagation(); handlePreview(c); }}
												disabled={previewingContractId === c.id}
												class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer disabled:opacity-50"
												title="Aperçu / Imprimer"
											>
												<Printer size={16} />
											</button>
											{#if canEdit(c)}
												<button
													type="button"
													onclick={(e) => { e.stopPropagation(); showEdit(c); }}
													class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
													title="Modifier"
												>
													<Pencil size={16} />
												</button>
											{/if}
											{#if canDelete(c)}
												<ConfirmDialog
													title="Supprimer le contrat"
													description="Le contrat sera définitivement supprimé. Cette action est irréversible."
													confirmLabel="Supprimer"
													onConfirm={() => handleDelete(c)}
												>
													<button
														type="button"
														onclick={(e) => e.stopPropagation()}
														disabled={deletingContractId === c.id}
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
		<!-- DETAIL VIEW -->
		<!-- ================================================================ -->
		{:else if viewMode === "detail" && selectedContract}
			<div class="max-w-5xl animate-fade-in flex flex-col gap-6">
				<DocumentLifecycleStepper
					steps={CONTRACT_STEPS}
					status={selectedContract.status}
					statusLabel={STATUS_LABELS[selectedContract.status] || selectedContract.status}
					note={statusNote(selectedContract)}
				/>

				<div class="grid grid-cols-1 lg:grid-cols-[1fr_320px] gap-6">
				<!-- Primary content -->
				<div class="flex flex-col gap-6 min-w-0">
				<div class="border border-border-card rounded-lg p-6">
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

				<!-- Metadata + signatures rail -->
				<div class="flex flex-col gap-6">
					<div class="border border-border-card rounded-lg p-6">
						<p class="text-xs text-muted-foreground mb-3 font-medium uppercase tracking-wider">Détails</p>
						<div class="flex flex-col gap-4">
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
					</div>

					<div class="border border-border-card rounded-lg p-6">
						<p class="text-xs text-muted-foreground mb-3 font-medium uppercase tracking-wider">Signatures</p>
						{#if selectedContract.signatures && selectedContract.signatures.length > 0}
							<div class="flex flex-col gap-3">
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
				</div>
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
			<div class="flex-shrink-0 px-8 pt-7 pb-5 border-b border-border-card relative">
				<Dialog.Title class="text-base font-semibold tracking-tight text-foreground">
					Nouveau contrat
				</Dialog.Title>
				{#if client}
					<p class="text-sm text-muted-foreground mt-0.5">{client.name}</p>
				{/if}
				<Dialog.Close class="absolute right-5 top-6 rounded-md text-muted-foreground hover:text-foreground transition-interactive duration-150 cursor-pointer p-0.5">
					<X class="size-4" />
					<span class="sr-only">Fermer</span>
				</Dialog.Close>
			</div>

			<!-- Modal body -->
			<div class="flex-1 min-h-0 overflow-y-auto px-8 py-6 flex flex-col gap-5">

				<!-- Dates + financials -->
				<div class="grid grid-cols-2 gap-4">
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Date de début</label>
						<input type="date" bind:value={formStartDate} class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2" />
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Date de fin <span class="font-normal text-muted-foreground">(optionnel)</span></label>
						<input type="date" bind:value={formEndDate} class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2" />
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Montant <span class="font-normal text-muted-foreground">(€, optionnel)</span></label>
						<input type="number" bind:value={formContractValue} min="0" step="0.01" placeholder="0.00" class="h-input rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2" />
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Devise</label>
						<input type="text" bind:value={formCurrency} placeholder="EUR" class="h-input rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2" />
					</div>
				</div>

				<!-- Objet section -->
				<div class="flex flex-col gap-3">
					<div class="flex items-center gap-3">
						<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider whitespace-nowrap">Objet</span>
						<div class="h-px flex-1 bg-border-input"></div>
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Objet des prestations</label>
						<textarea bind:value={formScopeOfServices} placeholder="Décrivez l'objet des prestations..." rows="3" class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
					</div>
				</div>

				<!-- Clauses section -->
				<div class="flex flex-col gap-3">
					<div class="flex items-center gap-3">
						<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider whitespace-nowrap">Clauses</span>
						<div class="h-px flex-1 bg-border-input"></div>
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Conditions de paiement</label>
						<textarea bind:value={formPaymentTerms} placeholder="Ex : Paiement à 30 jours..." rows="2" class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Clause de confidentialité</label>
						<textarea bind:value={formConfidentiality} placeholder="Clause de confidentialité..." rows="2" class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Clause de résiliation</label>
						<textarea bind:value={formTerminationClause} placeholder="Conditions de résiliation..." rows="2" class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
					</div>
				</div>

			</div>

			<!-- Modal footer -->
			<div class="flex items-center justify-between px-8 py-4 border-t border-border-card shrink-0">
				<p class="text-xs text-muted-foreground">Le contrat sera créé en brouillon.</p>
				<div class="flex items-center gap-2">
					<Dialog.Close class="h-input rounded-input bg-transparent text-foreground hover:bg-muted inline-flex items-center justify-center px-4 text-sm font-medium active:scale-[0.98] border border-border-input cursor-pointer transition-interactive duration-150 focus:outline-none">Annuler</Dialog.Close>
					<button type="button" onclick={handleCreate} disabled={formSaving} class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-1.5 px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50 transition-interactive duration-150">
						<Save size={14} />
						{formSaving ? "Enregistrement..." : "Créer le contrat"}
					</button>
				</div>
			</div>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>

<!-- ================================================================ -->
<!-- EDIT MODAL -->
<!-- ================================================================ -->
<Dialog.Root bind:open={showEditModal}>
	<Dialog.Portal>
		<Dialog.Overlay
			class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/60 backdrop-blur-[2px]"
		/>
		<Dialog.Content
			class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-2xl translate-x-[-50%] translate-y-[-50%] border flex flex-col max-h-[90vh]"
		>
			<!-- Modal header -->
			<div class="flex-shrink-0 px-8 pt-7 pb-5 border-b border-border-card relative">
				<Dialog.Title class="text-base font-semibold tracking-tight text-foreground">
					Modifier le contrat
				</Dialog.Title>
				{#if selectedContract}
					<p class="text-sm text-muted-foreground mt-0.5">{selectedContract.contract_number}</p>
				{/if}
				<Dialog.Close class="absolute right-5 top-6 rounded-md text-muted-foreground hover:text-foreground transition-interactive duration-150 cursor-pointer p-0.5">
					<X class="size-4" />
					<span class="sr-only">Fermer</span>
				</Dialog.Close>
			</div>

			<!-- Modal body -->
			<div class="flex-1 min-h-0 overflow-y-auto px-8 py-6 flex flex-col gap-5">

				<div class="grid grid-cols-2 gap-4">
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Date de début</label>
						<input type="date" bind:value={formStartDate} class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2" />
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Date de fin <span class="font-normal text-muted-foreground">(optionnel)</span></label>
						<input type="date" bind:value={formEndDate} class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2" />
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Montant <span class="font-normal text-muted-foreground">(€, optionnel)</span></label>
						<input type="number" bind:value={formContractValue} min="0" step="0.01" placeholder="0.00" class="h-input rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2" />
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Devise</label>
						<input type="text" bind:value={formCurrency} placeholder="EUR" class="h-input rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2" />
					</div>
				</div>

				<div class="flex flex-col gap-1.5">
					<label class="text-sm font-medium text-foreground">Objet des prestations</label>
					<textarea bind:value={formScopeOfServices} placeholder="Décrivez l'objet des prestations..." rows="3" class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
				</div>

				<div class="flex flex-col gap-1.5">
					<label class="text-sm font-medium text-foreground">Conditions de paiement</label>
					<textarea bind:value={formPaymentTerms} placeholder="Ex : Paiement à 30 jours..." rows="2" class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
				</div>

				<div class="flex flex-col gap-1.5">
					<label class="text-sm font-medium text-foreground">Clause de confidentialité</label>
					<textarea bind:value={formConfidentiality} placeholder="Clause de confidentialité..." rows="2" class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
				</div>

				<div class="flex flex-col gap-1.5">
					<label class="text-sm font-medium text-foreground">Clause de résiliation</label>
					<textarea bind:value={formTerminationClause} placeholder="Conditions de résiliation..." rows="2" class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
				</div>

			</div>

			<!-- Modal footer -->
			<div class="flex items-center justify-between px-8 py-4 border-t border-border-card shrink-0">
				<p class="text-xs text-muted-foreground">Seuls les contrats en brouillon peuvent être modifiés.</p>
				<div class="flex items-center gap-2">
					<Dialog.Close class="h-input rounded-input bg-transparent text-foreground hover:bg-muted inline-flex items-center justify-center px-4 text-sm font-medium active:scale-[0.98] border border-border-input cursor-pointer transition-interactive duration-150 focus:outline-none">Annuler</Dialog.Close>
					<button type="button" onclick={handleEditSave} disabled={formSaving} class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-1.5 px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50 transition-interactive duration-150">
						<Save size={14} />
						{formSaving ? "Enregistrement..." : "Enregistrer"}
					</button>
				</div>
			</div>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>

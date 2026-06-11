<script lang="ts">
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import {
		Send,
		CheckCircle,
		ShieldCheck,
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
		fetchCaseMandates,
		createMandate,
		sendDocument,
		signMandate,
		activateMandate,
		ConflictError,
	} from "$lib/services/api";
	import { currentCase } from "$lib/stores/case";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import { documentStatusBadge } from "$lib/utils/badgeVariants";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import type { Case, Client, DocumentStatus, Mandate } from "$lib/types/entities";

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
	let mandates: Mandate[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	// Currently selected mandate for viewing
	let selectedMandate: Mandate | null = $state(null);
	let viewMode: "list" | "detail" = $state("list");
	let showCreateModal = $state(false);

	// Create form state
	let formScopeOfWork = $state("");
	let formTermsConditions = $state("");
	let formValidFrom = $state(todayISO());
	let formValidUntil = $state("");
	let formJurisdiction = $state("");
	let formSpecialInstructions = $state("");
	let formSaving = $state(false);

	// Lifecycle action state
	let sendingMandate = $state(false);
	let signingMandate = $state(false);
	let activatingMandate = $state(false);

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

	// =========================================================================
	// Lifecycle guards
	// =========================================================================

	function canSend(m: Mandate): boolean {
		return m.status === "draft";
	}

	function canSign(m: Mandate): boolean {
		return m.status === "sent";
	}

	function canActivate(m: Mandate): boolean {
		return m.status === "signed";
	}

	function hasNoActions(m: Mandate): boolean {
		return m.status !== "active" && !canSend(m) && !canSign(m) && !canActivate(m);
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
			const [c, mands] = await Promise.all([
				fetchCase(caseId),
				fetchCaseMandates(caseId),
			]);
			caseData = c;
			mandates = mands;

			if (c?.clientId) {
				client = await fetchClient(c.clientId);
			}
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Erreur lors du chargement des mandats";
		} finally {
			loading = false;
		}
	}

	async function refreshMandates() {
		const id = selectedMandate?.id;
		if (!id) return;
		const freshMandates = await fetchCaseMandates(caseId);
		mandates = freshMandates;
		const updated = freshMandates.find((m) => m.id === id);
		if (updated) selectedMandate = updated;
	}

	// =========================================================================
	// Navigation helpers
	// =========================================================================

	function showList() {
		selectedMandate = null;
		viewMode = "list";
	}

	function showCreate() {
		formScopeOfWork = "";
		formTermsConditions = "";
		formValidFrom = todayISO();
		formValidUntil = "";
		formJurisdiction = "";
		formSpecialInstructions = "";
		showCreateModal = true;
	}

	function showDetail(m: Mandate) {
		selectedMandate = m;
		viewMode = "detail";
	}

	async function handleCreate() {
		if (!caseData) return;
		if (!formScopeOfWork.trim() || !formTermsConditions.trim() || !formValidFrom) {
			toastState.add(TOAST_LEVELS.Error, "Erreur", "Veuillez remplir tous les champs obligatoires.");
			return;
		}

		formSaving = true;
		try {
			const mandateNumber = `MAN-${new Date().getFullYear()}-${String(mandates.length + 1).padStart(3, "0")}`;
			const payload = {
				case_id: caseData.id,
				client_id: caseData.clientId,
				mandate_number: mandateNumber,
				issue_date: new Date().toISOString(),
				scope_of_work: formScopeOfWork.trim(),
				terms_conditions: formTermsConditions.trim(),
				valid_from: new Date(formValidFrom).toISOString(),
				valid_until: formValidUntil ? new Date(formValidUntil).toISOString() : undefined,
				jurisdiction: formJurisdiction.trim() || undefined,
				special_instructions: formSpecialInstructions.trim() || undefined,
				status: "draft" as const,
			} as Mandate;

			const result = await createMandate(payload);
			if (result.data) {
				mandates = [...mandates, result.data];
				selectedMandate = result.data;
				showCreateModal = false;
				viewMode = "detail";
				toastState.add(TOAST_LEVELS.Info, "Mandat créé", "Le mandat a été créé en brouillon.");
			}
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de créer le mandat.",
			);
		} finally {
			formSaving = false;
		}
	}

	// =========================================================================
	// Send Mandate (draft → sent)
	// =========================================================================

	async function handleSend() {
		if (!selectedMandate || !canSend(selectedMandate)) return;
		sendingMandate = true;
		try {
			const result = await sendDocument(selectedMandate.id, "mandate", {
				recipients: [],
				subject: `Mandat ${selectedMandate.mandate_number}`,
				message: "",
				send_email: true,
				send_sms: false,
			});

			if (result.success) {
				await refreshMandates();
				toastState.add(
					TOAST_LEVELS.Info,
					"Mandat envoyé",
					"Le mandat a été marqué comme envoyé.",
				);
			}
		} catch (e) {
			const msg = e instanceof ConflictError
				? e.message
				: e instanceof Error
					? e.message
					: "Impossible d'envoyer le mandat.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			sendingMandate = false;
		}
	}

	// =========================================================================
	// Sign Mandate (sent/draft → signed)
	// =========================================================================

	async function handleSign() {
		if (!selectedMandate || !canSign(selectedMandate)) return;
		signingMandate = true;
		try {
			const result = await signMandate(selectedMandate.id, {
				signer_name: "Enquêteur",
				signer_role: "investigator",
				method: "e-sign",
			});

			if (result.success) {
				await refreshMandates();
				toastState.add(
					TOAST_LEVELS.Info,
					"Mandat signé",
					"La signature a été enregistrée.",
				);
			}
		} catch (e) {
			const msg = e instanceof ConflictError
				? e.message
				: e instanceof Error
					? e.message
					: "Impossible de signer le mandat.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			signingMandate = false;
		}
	}

	// =========================================================================
	// Activate Mandate (signed → active)
	// =========================================================================

	async function handleActivate() {
		if (!selectedMandate || !canActivate(selectedMandate)) return;
		activatingMandate = true;
		try {
			const result = await activateMandate(selectedMandate.id);
			if (result.success) {
				await refreshMandates();
				toastState.add(
					TOAST_LEVELS.Info,
					"Mandat activé",
					"Le mandat est maintenant actif. L'enquête est autorisée.",
				);
			}
		} catch (e) {
			const msg = e instanceof ConflictError
				? e.message
				: e instanceof Error
					? e.message
					: "Impossible d'activer le mandat.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			activatingMandate = false;
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
						{#if viewMode === "detail" && selectedMandate}
							Mandat {selectedMandate.mandate_number}
						{:else}
							Mandats
						{/if}
					</h1>
					<p class="text-sm text-muted-foreground">
						{#if viewMode === "list"}
							Liste des mandats du dossier
						{:else if client}
							Client : {client.name}
						{/if}
					</p>
				</div>
			</div>

			{#if viewMode === "list" && mandates.length > 0}
				<button
					type="button"
					onclick={showCreate}
					class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
				>
					Nouveau mandat
				</button>
			{/if}

			{#if viewMode === "detail" && selectedMandate}
				<div class="flex items-center gap-2">
					{#if canSend(selectedMandate)}
						<ConfirmDialog
							title="Envoyer le mandat"
							description="Le mandat sera marqué comme envoyé au client. Cette action est irréversible."
							confirmLabel="Envoyer"
							onConfirm={handleSend}
						>
							<button
								type="button"
								disabled={sendingMandate}
								class="h-input rounded-input bg-accent text-accent-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<Send size={14} class="mr-1" />
								{sendingMandate ? "Envoi..." : "Envoyer"}
							</button>
						</ConfirmDialog>
					{/if}
					{#if canSign(selectedMandate)}
						<ConfirmDialog
							title="Signer le mandat"
							description="Enregistrez la signature du mandat. Le mandat passera en état « Signé » si toutes les signatures sont collectées."
							confirmLabel="Signer"
							onConfirm={handleSign}
						>
							<button
								type="button"
								disabled={signingMandate}
								class="h-input rounded-input bg-success text-success-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<CheckCircle size={14} class="mr-1" />
								{signingMandate ? "Signature..." : "Signer"}
							</button>
						</ConfirmDialog>
					{/if}
					{#if canActivate(selectedMandate)}
						<ConfirmDialog
							title="Activer le mandat"
							description="Le mandat sera activé, autorisant formellement le début de l'enquête. Cette action est irréversible."
							confirmLabel="Activer"
							onConfirm={handleActivate}
						>
							<button
								type="button"
								disabled={activatingMandate}
								class="h-input rounded-input bg-success text-success-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<ShieldCheck size={14} class="mr-1" />
								{activatingMandate ? "Activation..." : "Activer"}
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
			{#if mandates.length === 0}
				<div class="border border-dashed border-border rounded-lg p-12 bg-muted/20 flex flex-col items-center justify-center gap-4 flex-1 min-h-[50vh]">
					<FileText class="w-12 h-12 text-muted-foreground" />
					<p class="text-muted-foreground text-center">Aucun mandat pour ce dossier.</p>
					<p class="text-sm text-muted-foreground text-center">
						Un mandat est automatiquement créé lorsqu'un devis est accepté, ou vous pouvez en créer un manuellement.
					</p>
					<button
						type="button"
						onclick={showCreate}
						class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
					>
						Créer un mandat
					</button>
				</div>
			{:else}
				<div class="border border-border-card rounded-lg overflow-hidden">
					<table class="w-full">
						<thead class="bg-muted">
							<tr>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Référence</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Date d'émission</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Validité</th>
								<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">Statut</th>
							</tr>
						</thead>
						<tbody class="bg-background divide-y divide-border">
							{#each mandates as m (m.id)}
								<tr
									class="hover:bg-muted/50 transition-colors cursor-pointer"
									onclick={() => showDetail(m)}
								>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm font-medium text-foreground">{m.mandate_number}</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm text-foreground">{formatDate(m.issue_date)}</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<div class="text-sm text-foreground">
											{formatDate(m.valid_from)} — {m.valid_until ? formatDate(m.valid_until) : "—"}
										</div>
									</td>
									<td class="px-6 py-4 whitespace-nowrap">
										<span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {documentStatusBadge(m.status as DocumentStatus)}">
											{STATUS_LABELS[m.status] || m.status}
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
		{:else if viewMode === "detail" && selectedMandate}
			<div class="max-w-3xl animate-fade-in">
				<!-- Status banner -->
				<div class="flex items-center gap-3 mb-6">
					<span class="px-3 py-1.5 inline-flex text-sm leading-5 font-semibold rounded-full {documentStatusBadge(selectedMandate.status as DocumentStatus)}">
						{STATUS_LABELS[selectedMandate.status] || selectedMandate.status}
					</span>
					{#if selectedMandate.status === "active"}
						<span class="text-sm text-success">
							Enquête autorisée — mandat en vigueur.
						</span>
					{:else if selectedMandate.status === "signed"}
						<span class="text-sm text-success">
							Signé — en attente d'activation.
						</span>
					{:else if selectedMandate.status === "sent"}
						<span class="text-sm text-accent-subtle-foreground">
							Envoyé au client — en attente de signature.
						</span>
					{:else if selectedMandate.status === "draft"}
						<span class="text-sm text-muted-foreground">
							Brouillon — à envoyer pour signature.
						</span>
					{:else if selectedMandate.status === "archived"}
						<span class="text-sm text-muted-foreground">
							Archivé — ce mandat n'est plus actif.
						</span>
					{/if}
				</div>

				<!-- Mandate info card -->
				<div class="border border-border-card rounded-lg p-6 mb-6">
					<div class="grid grid-cols-2 gap-4 mb-6">
						<div>
							<p class="text-xs text-muted-foreground mb-1">Référence</p>
							<p class="text-sm font-semibold text-foreground">{selectedMandate.mandate_number}</p>
						</div>
						<div>
							<p class="text-xs text-muted-foreground mb-1">Date d'émission</p>
							<p class="text-sm text-foreground">{formatDate(selectedMandate.issue_date)}</p>
						</div>
						<div>
							<p class="text-xs text-muted-foreground mb-1">Valide du</p>
							<p class="text-sm text-foreground">{formatDate(selectedMandate.valid_from)}</p>
						</div>
						<div>
							<p class="text-xs text-muted-foreground mb-1">Valide jusqu'au</p>
							<p class="text-sm text-foreground">
								{selectedMandate.valid_until ? formatDate(selectedMandate.valid_until) : "—"}
							</p>
						</div>
						{#if selectedMandate.jurisdiction}
							<div>
								<p class="text-xs text-muted-foreground mb-1">Juridiction</p>
								<p class="text-sm text-foreground">{selectedMandate.jurisdiction}</p>
							</div>
						{/if}
						{#if selectedMandate.linked_estimate_id}
							<div>
								<p class="text-xs text-muted-foreground mb-1">Devis lié</p>
								<p class="text-sm text-foreground">{selectedMandate.linked_estimate_id}</p>
							</div>
						{/if}
					</div>

					<!-- Scope of work -->
					<div class="mb-4">
						<p class="text-xs text-muted-foreground mb-1">Objet de la mission</p>
						<p class="text-sm text-foreground">{selectedMandate.scope_of_work}</p>
					</div>

					<!-- Terms and conditions -->
					<div class="mb-4">
						<p class="text-xs text-muted-foreground mb-1">Conditions</p>
						<p class="text-sm text-muted-foreground">{selectedMandate.terms_conditions}</p>
					</div>

					{#if selectedMandate.special_instructions}
						<div class="mb-4">
							<p class="text-xs text-muted-foreground mb-1">Instructions spéciales</p>
							<p class="text-sm text-muted-foreground italic">{selectedMandate.special_instructions}</p>
						</div>
					{/if}

					<!-- Signatures section -->
					<div class="mt-6 pt-4 border-t border-border">
						<p class="text-xs text-muted-foreground mb-3 font-medium uppercase tracking-wider">Signatures</p>
						<div class="grid grid-cols-2 gap-4">
							<!-- Client signature -->
							<div class="border border-border rounded-lg p-4">
								<p class="text-xs text-muted-foreground mb-2">Client</p>
								{#if selectedMandate.client_signature}
									<div class="flex items-center gap-2">
										<CheckCircle size={16} class="text-success flex-shrink-0" />
										<div>
											<p class="text-sm font-medium text-foreground">{selectedMandate.client_signature.name}</p>
											<p class="text-xs text-muted-foreground">
												Signé le {formatDate(selectedMandate.client_signature.signed_at)}
											</p>
										</div>
									</div>
								{:else}
									<p class="text-sm text-muted-foreground">En attente de signature</p>
								{/if}
							</div>
							<!-- Investigator signature -->
							<div class="border border-border rounded-lg p-4">
								<p class="text-xs text-muted-foreground mb-2">Enquêteur</p>
								{#if selectedMandate.investigator_signature}
									<div class="flex items-center gap-2">
										<CheckCircle size={16} class="text-success flex-shrink-0" />
										<div>
											<p class="text-sm font-medium text-foreground">{selectedMandate.investigator_signature.name}</p>
											<p class="text-xs text-muted-foreground">
												Signé le {formatDate(selectedMandate.investigator_signature.signed_at)}
											</p>
										</div>
									</div>
								{:else}
									<p class="text-sm text-muted-foreground">En attente de signature</p>
								{/if}
							</div>
						</div>
					</div>
				</div>

				<!-- Lifecycle info for non-actionable states -->
				{#if hasNoActions(selectedMandate)}
					<div class="flex items-start gap-2 text-sm text-muted-foreground bg-muted/30 border border-border rounded-lg p-4">
						<AlertTriangle size={16} class="flex-shrink-0 mt-0.5" />
						<p>
							Ce mandat est dans l'état <strong>{STATUS_LABELS[selectedMandate.status]}</strong>.
							{#if selectedMandate.status === "active"}
								L'enquête est autorisée et en cours.
							{:else if selectedMandate.status === "archived"}
								Ce mandat a été archivé et n'est plus modifiable.
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
			<div class="flex-shrink-0 px-8 pt-7 pb-5 border-b border-border-card relative">
				<Dialog.Title class="text-base font-semibold tracking-tight text-foreground">
					Nouveau mandat
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

				<div class="grid grid-cols-2 gap-4">
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Valide du</label>
						<input type="date" bind:value={formValidFrom} class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2" />
					</div>
					<div class="flex flex-col gap-1.5">
						<label class="text-sm font-medium text-foreground">Valide jusqu'au <span class="font-normal text-muted-foreground">(optionnel)</span></label>
						<input type="date" bind:value={formValidUntil} class="h-input rounded-input border border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2" />
					</div>
				</div>

				<div class="flex flex-col gap-1.5">
					<label class="text-sm font-medium text-foreground">Objet de la mission</label>
					<textarea bind:value={formScopeOfWork} placeholder="Décrivez l'objet de la mission..." rows="3" class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
				</div>

				<div class="flex flex-col gap-1.5">
					<label class="text-sm font-medium text-foreground">Conditions</label>
					<textarea bind:value={formTermsConditions} placeholder="Conditions générales du mandat..." rows="3" class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
				</div>

				<div class="flex flex-col gap-1.5">
					<label class="text-sm font-medium text-foreground">Juridiction <span class="font-normal text-muted-foreground">(optionnel)</span></label>
					<input type="text" bind:value={formJurisdiction} placeholder="Ex : France, Paris..." class="h-input rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2" />
				</div>

				<div class="flex flex-col gap-1.5">
					<label class="text-sm font-medium text-foreground">Instructions spéciales <span class="font-normal text-muted-foreground">(optionnel)</span></label>
					<textarea bind:value={formSpecialInstructions} placeholder="Instructions particulières..." rows="2" class="rounded-input border border-border-input bg-background placeholder:text-muted-foreground/40 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
				</div>

			</div>

			<!-- Modal footer -->
			<div class="flex items-center justify-between px-8 py-4 border-t border-border-card shrink-0">
				<p class="text-xs text-muted-foreground">Le mandat sera créé en brouillon.</p>
				<div class="flex items-center gap-2">
					<Dialog.Close class="h-input rounded-input bg-transparent text-foreground hover:bg-muted inline-flex items-center justify-center px-4 text-sm font-medium active:scale-[0.98] border border-border-input cursor-pointer transition-interactive duration-150 focus:outline-none">Annuler</Dialog.Close>
					<button type="button" onclick={handleCreate} disabled={formSaving} class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50 transition-interactive duration-150">
						<Save size={14} class="mr-1.5" />
						{formSaving ? "Enregistrement..." : "Créer le mandat"}
					</button>
				</div>
			</div>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>

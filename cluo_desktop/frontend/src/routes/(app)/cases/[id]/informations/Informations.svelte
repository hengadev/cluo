<script lang="ts">
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { Briefcase, MapPin, Pencil, Check, X, Tag, User, UserPlus, Trash2 } from "@lucide/svelte";
	import {
		fetchCase,
		fetchClient,
		fetchContact,
		fetchCaseSubject,
		fetchAllClients,
		fetchAllCaseTypes,
		fetchClientContacts,
		updateCase,
		createCaseSubject,
		updateCaseSubject,
		deleteCaseSubject,
	} from "$lib/services/api";
	import { recentCases } from "$lib/stores/case";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import DocumentWorkflowSummary from "./DocumentWorkflowSummary.svelte";
	import type {
		Case,
		CaseType,
		Client,
		Contact,
		CaseSubject,
		CaseStatus,
	} from "$lib/types/entities";

	const toastState = getToastContext();

	const STATUS_LABELS: Record<CaseStatus, string> = {
		in_progress: "En cours",
		ready: "Prêt",
		released: "Clôturé",
	};

	let caseData: Case | null = $state(null);
	let client: Client | null = $state(null);
	let contact: Contact | null = $state(null);
	let subject: CaseSubject | null = $state(null);
	let caseTypeName: string | null = $state(null);
	let allCaseTypes: CaseType[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	// CaseSubject management
	let showSubjectForm = $state(false);
	let editingSubject = $state(false);
	let savingSubject = $state(false);
	let detachingSubject = $state(false);
	let subjectForm = $state({
		firstname: "",
		lastname: "",
		email: "",
		phone: "",
		address1: "",
		address2: "",
		city: "",
		postalCode: "",
		occupation: "",
		notes: "",
	});

	// CaseType edit mode
	let editingCaseType = $state(false);
	let selectedCaseTypeId = $state("");
	let savingCaseType = $state(false);

	// Client/Contact edit mode
	let editingClient = $state(false);
	let allClients: Client[] = $state([]);
	let selectedClientId = $state("");
	let contactsForClient: Contact[] = $state([]);
	let selectedContactId = $state("");
	let clientSearchQuery = $state("");
	let savingClient = $state(false);
	let loadingContacts = $state(false);

	const caseId = $derived($page.params.id);

	let filteredClients = $derived(
		clientSearchQuery.trim() === ""
			? allClients
			: allClients.filter((c) =>
					c.name.toLowerCase().includes(clientSearchQuery.toLowerCase()),
				),
	);

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
			caseData = await fetchCase(caseId);

			if (!caseData) {
				error = "Dossier introuvable";
				loading = false;
				return;
			}

			recentCases.push({ id: caseData.id, title: caseData.title, status: caseData.status });

			const [clientData, contactData, subjectData, typesData] = await Promise.all([
				fetchClient(caseData.clientId),
				caseData.assignedContactID
					? fetchContact(caseData.assignedContactID)
					: Promise.resolve(null),
				caseData.caseSubjectId
					? fetchCaseSubject(caseData.caseSubjectId)
					: Promise.resolve(null),
				fetchAllCaseTypes(),
			]);

			client = clientData;
			contact = contactData;
			subject = subjectData;
			allCaseTypes = typesData;
			if (caseData?.caseTypeId) {
				const ct = typesData.find((t: CaseType) => t.id === caseData!.caseTypeId);
				caseTypeName = ct ? ct.name : null;
			}
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Erreur lors du chargement des données";
		} finally {
			loading = false;
		}
	}

	async function startClientEdit() {
		editingClient = true;
		selectedClientId = caseData?.clientId || "";
		selectedContactId = caseData?.assignedContactID || "";
		clientSearchQuery = "";

		if (allClients.length === 0) {
			try {
				allClients = await fetchAllClients();
			} catch (e) {
				toastState.add(
					TOAST_LEVELS.Error,
					"Erreur",
					"Impossible de charger la liste des clients.",
				);
			}
		}

		if (selectedClientId) {
			await loadContactsForClient(selectedClientId);
		}
	}

	async function loadContactsForClient(clientId: string) {
		if (!clientId) {
			contactsForClient = [];
			return;
		}
		contactsForClient = [];
		loadingContacts = true;
		try {
			contactsForClient = await fetchClientContacts(clientId);
		} catch {
			contactsForClient = [];
		} finally {
			loadingContacts = false;
		}
	}

	function cancelClientEdit() {
		editingClient = false;
		clientSearchQuery = "";
	}

	async function saveClientEdit() {
		if (!caseData || !selectedClientId) return;
		savingClient = true;
		try {
			caseData = await updateCase(caseData.id, {
				clientId: selectedClientId,
				assignedContactID: selectedContactId || undefined,
			});
			// Refresh display data
			client = await fetchClient(selectedClientId);
			contact = selectedContactId
				? await fetchContact(selectedContactId)
				: null;
			editingClient = false;
			toastState.add(
				TOAST_LEVELS.Info,
				"Dossier mis à jour",
				"Le client et le contact ont été mis à jour.",
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de mettre à jour",
			);
		} finally {
			savingClient = false;
		}
	}

	// CaseType edit
	function startCaseTypeEdit() {
		editingCaseType = true;
		selectedCaseTypeId = caseData?.caseTypeId || "";
	}

	function cancelCaseTypeEdit() {
		editingCaseType = false;
	}

	async function saveCaseTypeEdit() {
		if (!caseData) return;
		savingCaseType = true;
		try {
			caseData = await updateCase(caseData.id, {
				caseTypeId: selectedCaseTypeId || null,
			});
			if (selectedCaseTypeId) {
				const ct = allCaseTypes.find((t) => t.id === selectedCaseTypeId);
				caseTypeName = ct ? ct.name : null;
			} else {
				caseTypeName = null;
			}
			editingCaseType = false;
			toastState.add(
				TOAST_LEVELS.Info,
				"Dossier mis à jour",
				"Le type d'affaire a été mis à jour.",
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de mettre à jour",
			);
		} finally {
			savingCaseType = false;
		}
	}

	// =========================================================================
	// CaseSubject management
	// =========================================================================

	function startCreateSubject() {
		showSubjectForm = true;
		editingSubject = false;
		subjectForm = {
			firstname: "",
			lastname: "",
			email: "",
			phone: "",
			address1: "",
			address2: "",
			city: "",
			postalCode: "",
			occupation: "",
			notes: "",
		};
	}

	function cancelSubjectForm() {
		showSubjectForm = false;
		editingSubject = false;
	}

	function startEditSubject() {
		if (!subject) return;
		editingSubject = true;
		showSubjectForm = true;
		subjectForm = {
			firstname: subject.firstname,
			lastname: subject.lastname,
			email: subject.email || "",
			phone: subject.phone || "",
			address1: subject.address1 || "",
			address2: subject.address2 || "",
			city: subject.city || "",
			postalCode: subject.postalCode || "",
			occupation: subject.occupation || "",
			notes: subject.notes || "",
		};
	}

	async function saveSubjectCreate() {
		if (!caseData) return;
		if (!subjectForm.lastname.trim() || !subjectForm.firstname.trim()) {
			toastState.add(TOAST_LEVELS.Error, "Erreur", "Le nom et le prénom sont requis.");
			return;
		}
		savingSubject = true;
		try {
			// 1. Create the CaseSubject
			const newSubject = await createCaseSubject({
				firstname: subjectForm.firstname.trim(),
				lastname: subjectForm.lastname.trim(),
				email: subjectForm.email.trim() || undefined,
				phone: subjectForm.phone.trim() || undefined,
				address1: subjectForm.address1.trim() || undefined,
				address2: subjectForm.address2.trim() || undefined,
				city: subjectForm.city.trim() || undefined,
				postalCode: subjectForm.postalCode.trim() || undefined,
				occupation: subjectForm.occupation.trim() || undefined,
				notes: subjectForm.notes.trim() || undefined,
			});
			// 2. Attach to the Case
			caseData = await updateCase(caseData.id, {
				caseSubjectId: newSubject.id,
			});
			subject = newSubject;
			showSubjectForm = false;
			toastState.add(
				TOAST_LEVELS.Info,
				"Personne ajoutée",
				"La personne a été créée et associée au dossier.",
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de créer la personne.",
			);
		} finally {
			savingSubject = false;
		}
	}

	async function saveSubjectEdit() {
		if (!subject) return;
		savingSubject = true;
		try {
			const updated = await updateCaseSubject(subject.id, {
				firstname: subjectForm.firstname.trim(),
				lastname: subjectForm.lastname.trim(),
				email: subjectForm.email.trim() || undefined,
				phone: subjectForm.phone.trim() || undefined,
				address1: subjectForm.address1.trim() || undefined,
				address2: subjectForm.address2.trim() || undefined,
				city: subjectForm.city.trim() || undefined,
				postalCode: subjectForm.postalCode.trim() || undefined,
				occupation: subjectForm.occupation.trim() || undefined,
				notes: subjectForm.notes.trim() || undefined,
			});
			subject = updated;
			editingSubject = false;
			showSubjectForm = false;
			toastState.add(
				TOAST_LEVELS.Info,
				"Personne mise à jour",
				"Les informations ont été enregistrées.",
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de mettre à jour.",
			);
		} finally {
			savingSubject = false;
		}
	}

	async function detachSubject() {
		if (!caseData || !subject) return;
		detachingSubject = true;
		try {
			// 1. Detach from the Case
			caseData = await updateCase(caseData.id, {
				caseSubjectId: null,
			});
			subject = null;
			toastState.add(
				TOAST_LEVELS.Info,
				"Personne détachée",
				"La personne a été détachée du dossier.",
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de détacher la personne.",
			);
		} finally {
			detachingSubject = false;
		}
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString("fr-FR", {
			day: "2-digit",
			month: "short",
			year: "numeric",
		});
	}
</script>

<div class="p-8">
	{#if loading}
		<div class="flex items-center justify-center py-12">
			<p class="text-muted-foreground">Chargement...</p>
		</div>
	{:else if error}
		<div
			class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg"
		>
			{error}
		</div>
	{:else if caseData}
		<div class="flex gap-8">
			<!-- Case Header -->
			<div
				class="grid gap-5 p-6 border border-border-card rounded-card flex-1 animate-fade-in hover:shadow-md transition-shadow duration-300"
				style="animation-delay: 100ms;"
			>
				<div class="flex gap-4 items-center">
					<span
						class="bg-blue-100 text-blue-800 px-2 py-1 rounded-card text-sm font-medium"
					>
						STATUT: {STATUS_LABELS[caseData.status] || caseData.status}
					</span>
					<p class="text-muted-foreground text-sm">
						Créé le {formatDate(caseData.createdAt)}
					</p>
				</div>
				<h2 class="text-3xl font-bold text-foreground">
					{caseData.title}
				</h2>
				<div class="flex gap-4 text-lg items-center">
					<span class="text-muted-foreground">ID de dossier:</span>
					<span class="font-mono text-foreground">#{caseData.id}</span>
				</div>
				{#if caseData.externalReference}
					<div class="flex gap-4 text-lg items-center">
						<span class="text-muted-foreground"
							>Référence externe:</span
						>
						<span class="text-foreground"
							>{caseData.externalReference}</span
						>
					</div>
				{/if}
				{#if caseData.description}
					<div class="mt-4">
						<p class="text-sm text-muted-foreground mb-1">Description</p>
						<p class="text-foreground">{caseData.description}</p>
					</div>
				{/if}
			</div>

			<div class="flex flex-col gap-6">
				<!-- Client Details -->
				<div
					class="border border-border-card rounded-card p-6 grid gap-4 animate-fade-in hover:shadow-md transition-shadow duration-300 w-80"
					style="animation-delay: 200ms;"
				>
					<div class="flex justify-between items-center">
						<p class="text-muted-foreground text-sm font-medium">CLIENT</p>
						<div class="flex items-center gap-2">
							{#if !editingClient}
								<button
									onclick={startClientEdit}
									class="p-1 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
									title="Modifier le client"
								>
									<Pencil size={14} />
								</button>
							{/if}
							<Briefcase class="w-5 h-5 text-muted-foreground" />
						</div>
					</div>

					{#if editingClient}
						<!-- Edit mode: Client selector -->
						<div class="flex flex-col gap-3">
							<div class="relative">
								<input
									type="text"
									placeholder="Rechercher un client..."
									bind:value={clientSearchQuery}
									class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
								/>
							</div>
							<select
								bind:value={selectedClientId}
								onchange={() => {
									selectedContactId = "";
									loadContactsForClient(selectedClientId);
								}}
								class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2 cursor-pointer"
							>
								<option value="">-- Sélectionner un client --</option>
								{#each filteredClients as c}
									<option value={c.id}>{c.name}</option>
								{/each}
							</select>

							<!-- Contact selector (shown when client is selected) -->
							{#if selectedClientId}
								{#if loadingContacts}
									<p class="text-xs text-muted-foreground">Chargement des contacts...</p>
								{:else if contactsForClient.length > 0}
									<select
										bind:value={selectedContactId}
										class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2 cursor-pointer"
									>
										<option value="">-- Aucun contact --</option>
										{#each contactsForClient as c}
											<option value={c.id}
												>{c.firstname} {c.lastname}</option
											>
										{/each}
									</select>
								{/if}
							{/if}

							<div class="flex justify-end gap-2">
								<button
									type="button"
									onclick={cancelClientEdit}
									class="h-input rounded-input bg-transparent text-dark hover:bg-[#fafafa] inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] border-2 border-[#dedede] cursor-pointer"
								>
									<X size={14} />
								</button>
								<button
									type="button"
									onclick={saveClientEdit}
									disabled={savingClient || !selectedClientId}
									class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
								>
									<Check size={14} />
								</button>
							</div>
						</div>
					{:else}
						<!-- Display mode -->
						{#if client}
							<div>
								<p class="font-semibold text-foreground">{client.name}</p>
							</div>
							{#if contact}
								<div class="border-t border-border pt-4 mt-2">
									<p class="text-sm text-muted-foreground mb-2"
										>Contact principal</p
									>
									<p class="font-medium text-foreground">
										{contact.firstname} {contact.lastname}
									</p>
									{#if contact.position}
										<p class="text-sm text-muted-foreground">
											{contact.position}
										</p>
									{/if}
									{#if contact.email}
										<p class="text-sm text-muted-foreground">
											{contact.email}
										</p>
									{/if}
									{#if contact.phone}
										<p class="text-sm text-muted-foreground">
											{contact.phone}
										</p>
									{/if}
								</div>
							{/if}
						{:else}
							<p class="text-sm text-muted-foreground"
								>Aucun client associé</p
							>
						{/if}
					{/if}
				</div>

				<!-- Location -->
				<div
					class="border border-border-card rounded-card p-6 grid gap-4 animate-fade-in hover:shadow-md transition-shadow duration-300 w-80"
					style="animation-delay: 300ms;"
				>
					<div class="flex justify-between items-center">
						<p class="text-muted-foreground text-sm font-medium">LIEU</p>
						<MapPin class="w-5 h-5 text-muted-foreground" />
					</div>
					{#if caseData.placename}
						<p class="font-semibold text-foreground">{caseData.placename}</p>
					{/if}
					{#if caseData.address1}
						<p class="text-sm text-foreground">
							{caseData.address1}
							{#if caseData.address2}<br />{caseData.address2}{/if}
						</p>
						<p class="text-sm text-foreground">
							{caseData.postalCode} {caseData.city}
						</p>
						{#if caseData.country}
							<p class="text-sm text-muted-foreground">{caseData.country}</p>
						{/if}
					{/if}
					{#if caseData.locationNotes}
						<p class="text-sm text-muted-foreground mt-2 italic">
							{caseData.locationNotes}
						</p>
					{/if}
				</div>

				<!-- Case Type -->
				<div
					class="border border-border-card rounded-card p-6 grid gap-4 animate-fade-in hover:shadow-md transition-shadow duration-300 w-80"
					style="animation-delay: 350ms;"
				>
					<div class="flex justify-between items-center">
						<p class="text-muted-foreground text-sm font-medium">TYPE D'AFFAIRE</p>
						<div class="flex items-center gap-2">
							{#if !editingCaseType}
								<button
									onclick={startCaseTypeEdit}
									class="p-1 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
									title="Modifier le type"
								>
									<Pencil size={14} />
								</button>
							{/if}
							<Tag class="w-5 h-5 text-muted-foreground" />
						</div>
					</div>

					{#if editingCaseType}
						<div class="flex flex-col gap-3">
							<select
								bind:value={selectedCaseTypeId}
								class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2 cursor-pointer"
							>
								<option value="">-- Aucun type --</option>
								{#each allCaseTypes as ct}
									<option value={ct.id}>{ct.name}</option>
								{/each}
							</select>
							<div class="flex justify-end gap-2">
								<button
									type="button"
									onclick={cancelCaseTypeEdit}
									class="h-input rounded-input bg-transparent text-dark hover:bg-[#fafafa] inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] border-2 border-[#dedede] cursor-pointer"
								>
									<X size={14} />
								</button>
								<button
									type="button"
									onclick={saveCaseTypeEdit}
									disabled={savingCaseType}
									class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
								>
									<Check size={14} />
								</button>
							</div>
						</div>
					{:else}
						{#if caseTypeName}
							<span class="bg-emerald-100 text-emerald-800 px-2 py-1 rounded-card text-sm font-medium w-fit">{caseTypeName}</span>
						{:else}
							<p class="text-sm text-muted-foreground">Non défini</p>
						{/if}
					{/if}
				</div>
			</div>
		</div>

		<!-- Case Subject -->
		<div
			class="mt-8 border border-border-card rounded-card p-6 animate-fade-in"
			style="animation-delay: 400ms;"
		>
			<div class="flex justify-between items-center mb-4">
				<h3 class="text-lg font-semibold text-foreground">
					Personne impliquée
				</h3>
				{#if subject && !showSubjectForm}
					<div class="flex items-center gap-2">
						<button
							onclick={startEditSubject}
							class="p-1 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
							title="Modifier la personne"
						>
							<Pencil size={14} />
						</button>
						<ConfirmDialog
							title="Détacher la personne"
							description="Voulez-vous détacher cette personne du dossier ? Les données de la personne seront conservées."
							onConfirm={detachSubject}
						>
							<button
								class="p-1 rounded hover:bg-red-50 text-muted-foreground hover:text-red-600 transition-colors cursor-pointer"
								title="Détacher la personne"
							>
								<Trash2 size={14} />
							</button>
						</ConfirmDialog>
					</div>
				{/if}
			</div>

			{#if showSubjectForm}
				<!-- Creation / Edit form -->
				<div class="border border-border rounded-lg p-5 bg-muted/30 max-w-2xl grid gap-4">
					<p class="text-sm font-medium text-muted-foreground col-span-2">
						{editingSubject ? "Modifier la personne" : "Ajouter une personne"}
					</p>
					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="text-xs text-muted-foreground mb-1 block">Prénom *</label>
							<input
								type="text"
								bind:value={subjectForm.firstname}
								placeholder="Prénom"
								class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
						<div>
							<label class="text-xs text-muted-foreground mb-1 block">Nom *</label>
							<input
								type="text"
								bind:value={subjectForm.lastname}
								placeholder="Nom"
								class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
					</div>
					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="text-xs text-muted-foreground mb-1 block">Email</label>
							<input
								type="email"
								bind:value={subjectForm.email}
								placeholder="email@example.com"
								class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
						<div>
							<label class="text-xs text-muted-foreground mb-1 block">Téléphone</label>
							<input
								type="tel"
								bind:value={subjectForm.phone}
								placeholder="+33 6 00 00 00 00"
								class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
					</div>
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Adresse</label>
						<input
							type="text"
								bind:value={subjectForm.address1}
								placeholder="Adresse ligne 1"
								class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
					<div>
						<input
							type="text"
							bind:value={subjectForm.address2}
							placeholder="Adresse ligne 2"
							class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
						/>
					</div>
					<div class="grid grid-cols-3 gap-3">
						<div>
							<label class="text-xs text-muted-foreground mb-1 block">Code postal</label>
							<input
								type="text"
								bind:value={subjectForm.postalCode}
								placeholder="75001"
								class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
						<div>
							<label class="text-xs text-muted-foreground mb-1 block">Ville</label>
							<input
								type="text"
								bind:value={subjectForm.city}
								placeholder="Paris"
								class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
						<div>
							<label class="text-xs text-muted-foreground mb-1 block">Profession</label>
							<input
								type="text"
								bind:value={subjectForm.occupation}
								placeholder="Profession"
								class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
					</div>
					<div>
						<label class="text-xs text-muted-foreground mb-1 block">Notes</label>
						<textarea
							bind:value={subjectForm.notes}
							placeholder="Notes..."
							rows="2"
							class="rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 py-2 text-sm focus:ring-2 focus:ring-offset-2 resize-none"
						></textarea>
					</div>
					<div class="flex justify-end gap-2 mt-2">
						<button
							type="button"
							onclick={cancelSubjectForm}
							class="h-input rounded-input bg-transparent text-dark hover:bg-[#fafafa] inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] border-2 border-[#dedede] cursor-pointer"
						>
							<X size={14} class="mr-1" />
							Annuler
						</button>
						<button
							type="button"
							onclick={editingSubject ? saveSubjectEdit : saveSubjectCreate}
							disabled={savingSubject}
							class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
						>
							<Check size={14} class="mr-1" />
							{savingSubject ? "Enregistrement..." : (editingSubject ? "Enregistrer" : "Ajouter")}
						</button>
					</div>
				</div>
			{:else if subject}
				<!-- Display mode: existing subject -->
				<div class="border border-border rounded-lg p-4 bg-muted/30 max-w-2xl">
					<div class="flex items-start gap-4">
						<div class="w-10 h-10 rounded-full bg-blue-100 text-blue-800 flex items-center justify-center flex-shrink-0">
							<User size={20} />
						</div>
						<div class="flex-1">
							<p class="text-sm font-semibold text-foreground">
								{subject.firstname} {subject.lastname}
							</p>
							{#if subject.occupation}
								<p class="text-sm text-muted-foreground">{subject.occupation}</p>
							{/if}
						</div>
					</div>
					<div class="grid grid-cols-2 gap-x-8 gap-y-2 mt-4">
						{#if subject.email}
							<div>
								<p class="text-xs text-muted-foreground">Email</p>
								<p class="text-sm text-foreground">{subject.email}</p>
							</div>
						{/if}
						{#if subject.phone}
							<div>
								<p class="text-xs text-muted-foreground">Téléphone</p>
								<p class="text-sm text-foreground">{subject.phone}</p>
							</div>
						{/if}
					</div>
					{#if subject.address1 || subject.city}
						<div class="mt-3">
							<p class="text-xs text-muted-foreground">Adresse</p>
							<p class="text-sm text-foreground">
								{subject.address1}{#if subject.address2}<br />{subject.address2}{/if}
							</p>
							{#if subject.postalCode || subject.city}
								<p class="text-sm text-foreground">
									{subject.postalCode} {subject.city}
								</p>
							{/if}
						</div>
					{/if}
					{#if subject.notes}
						<div class="mt-3">
							<p class="text-xs text-muted-foreground">Notes</p>
							<p class="text-sm text-muted-foreground italic">{subject.notes}</p>
						</div>
					{/if}
				</div>
			{:else}
				<!-- No subject: show Add button -->
				<div class="border border-dashed border-border rounded-lg p-6 bg-muted/20 max-w-md flex flex-col items-center gap-3">
					<UserPlus class="w-8 h-8 text-muted-foreground" />
					<p class="text-sm text-muted-foreground text-center">
						Aucune personne associée à ce dossier.
					</p>
					<button
						type="button"
						onclick={startCreateSubject}
						class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
					>
						<UserPlus size={16} class="mr-2" />
						Ajouter une personne
					</button>
				</div>
			{/if}
		</div>

		<!-- Document Workflow Summary -->
		{#if caseData}
			<DocumentWorkflowSummary caseId={caseId} />
		{/if}
	{/if}
</div>

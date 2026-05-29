<script lang="ts">
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { Briefcase, MapPin, Pencil, Check, X, Tag } from "@lucide/svelte";
	import {
		fetchCase,
		fetchClient,
		fetchContact,
		fetchCaseSubject,
		fetchAllClients,
		fetchAllCaseTypes,
		fetchClientContacts,
		updateCase,
	} from "$lib/services/api";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
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
		{#if subject}
			<div
				class="mt-8 border border-border-card rounded-card p-6 animate-fade-in"
				style="animation-delay: 400ms;"
			>
				<h3 class="text-lg font-semibold text-foreground mb-4">
					Personne impliquée
				</h3>
				<div class="border border-border rounded-lg p-4 bg-muted/30 max-w-md">
					<div class="mb-2">
						<span class="text-sm font-medium text-foreground">
							{subject.firstname} {subject.lastname}
						</span>
					</div>
					{#if subject.occupation}
						<p class="text-sm text-muted-foreground">{subject.occupation}</p>
					{/if}
					{#if subject.email}
						<p class="text-sm text-muted-foreground">{subject.email}</p>
					{/if}
					{#if subject.phone}
						<p class="text-sm text-muted-foreground">{subject.phone}</p>
					{/if}
					{#if subject.address1}
						<p class="text-sm text-muted-foreground">
							{subject.address1}{subject.address2
								? " " + subject.address2
								: ""}
						</p>
					{/if}
				</div>
			</div>
		{/if}
	{/if}
</div>

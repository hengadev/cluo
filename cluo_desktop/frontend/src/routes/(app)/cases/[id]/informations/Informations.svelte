<script lang="ts">
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { MapPin, Pencil, Check, X, User, UserPlus, Trash2, Send, Building2, Layers } from "@lucide/svelte";
	import {
		fetchCase,
		fetchClient,
		fetchContact,
		fetchCaseSubject,
		fetchAllCaseTypes,
		createCaseSubject,
		updateCase,
		updateCaseSubject,
	} from "$lib/services/api";
	import { releaseCaseAndNotify } from "$lib/services/caseActions";
	import { recentCases } from "$lib/stores/case";
	import { caseStatusBadge } from "$lib/utils/badgeVariants";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import CaseClientDialog from "$lib/custom/global/CaseClientDialog.svelte";
	import CaseTypeDialog from "$lib/custom/global/CaseTypeDialog.svelte";
	import CaseLocationDialog from "$lib/custom/global/CaseLocationDialog.svelte";
	import Spinner from "$lib/components/Spinner.svelte";
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
	let releasing = $state(false);

	// Modal open states
	let clientDialogOpen = $state(false);
	let caseTypeDialogOpen = $state(false);
	let locationDialogOpen = $state(false);

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

	const caseId = $derived($page.params.id);

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

	// =========================================================================
	// Case release (status: ready -> released)
	// =========================================================================

	async function releaseDossier() {
		if (!caseData || caseData.status !== 'ready') return;
		releasing = true;
		try {
			await releaseCaseAndNotify(caseData.id, caseData.title);
			caseData.status = 'released';
			recentCases.push({
				id: caseData.id,
				title: caseData.title,
				status: 'released',
			});
			toastState.add(
				TOAST_LEVELS.Info,
				"Dossier publié",
				"Le lien du portail client a été généré.",
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de publier le dossier.",
			);
		} finally {
			releasing = false;
		}
	}

	// =========================================================================
	// Modal save callbacks
	// =========================================================================

	function onClientSaved(savedClient: Client, savedContact: Contact | null) {
		client = savedClient;
		contact = savedContact;
		if (caseData) {
			caseData.clientId = savedClient.id;
			caseData.assignedContactID = savedContact?.id ?? null;
		}
	}

	function onCaseTypeSaved(updated: Case) {
		caseData = updated;
		if (updated.caseTypeId) {
			const ct = allCaseTypes.find((t) => t.id === updated.caseTypeId);
			caseTypeName = ct ? ct.name : null;
		} else {
			caseTypeName = null;
		}
	}

	function onLocationSaved(updated: Case) {
		caseData = updated;
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

	let hasLocation = $derived(
		!!caseData && !!(caseData.placename || caseData.address1 || caseData.city),
	);
</script>

{#if caseData}
	<CaseClientDialog
		bind:open={clientDialogOpen}
		caseId={caseData.id}
		currentClientId={caseData.clientId}
		currentContactId={caseData.assignedContactID}
		onSaved={onClientSaved}
	/>
	<CaseTypeDialog
		bind:open={caseTypeDialogOpen}
		caseId={caseData.id}
		currentCaseTypeId={caseData.caseTypeId}
		{allCaseTypes}
		onSaved={onCaseTypeSaved}
	/>
	<CaseLocationDialog
		bind:open={locationDialogOpen}
		caseId={caseData.id}
		{caseData}
		onSaved={onLocationSaved}
	/>
{/if}

<div class="page-content">
	{#if loading}
		<div class="flex items-center justify-center py-12">
			<Spinner size="lg" />
		</div>
	{:else if error}
		<div class="alert-error">
			{error}
		</div>
	{:else if caseData}
		<div class="flex gap-8">
			<!-- Case Header -->
			<div
				class="flex flex-col gap-3 p-6 border border-border-card rounded-card flex-1 hover:shadow-card transition-shadow duration-300"
			>
				<div class="flex items-center gap-3">
					<span
						class="{caseStatusBadge(caseData.status)} px-2 py-0.5 rounded-card text-xs font-medium"
					>
						{STATUS_LABELS[caseData.status] || caseData.status}
					</span>
					<p class="text-muted-foreground text-xs">
						Créé le {formatDate(caseData.createdAt)}
					</p>
					{#if caseData.status === 'ready'}
						<ConfirmDialog
							title="Publier le dossier"
							description="Le dossier sera publié sur le portail client. Un lien d'accès sécurisé sera généré et envoyé au client. Cette action est irréversible."
							confirmLabel="Publier"
							onConfirm={releaseDossier}
						>
							<button
								type="button"
								disabled={releasing}
								class="ml-auto h-input rounded-input bg-accent text-accent-foreground shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
							>
								<Send size={14} />
								{releasing ? "Publication..." : "Publier le dossier"}
							</button>
						</ConfirmDialog>
					{/if}
				</div>

				<h2 class="text-2xl font-bold text-foreground leading-tight">
					{caseData.title}
				</h2>

				<div class="flex flex-wrap gap-2">
					<span class="inline-flex items-center gap-1 bg-muted px-2 py-0.5 rounded text-xs text-muted-foreground font-mono">
						#{caseData.id}
					</span>
					{#if caseData.externalReference}
						<span class="inline-flex items-center gap-1 bg-muted px-2 py-0.5 rounded text-xs text-muted-foreground">
							Réf. {caseData.externalReference}
						</span>
					{/if}
				</div>

				{#if caseData.description}
					<p class="text-sm text-foreground">{caseData.description}</p>
				{/if}
			</div>

			<div class="flex flex-col gap-4">
				<!-- Client Details -->
				<div
					class="border border-border-card rounded-card p-6 grid gap-4 hover:shadow-card transition-shadow duration-300 w-80"
				>
					<div class="flex justify-between items-center">
						<p class="text-muted-foreground text-sm font-medium">CLIENT</p>
						<div class="flex items-center gap-2">
							<button
								onclick={() => (clientDialogOpen = true)}
								class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
								title="Modifier le client"
							>
								<Pencil size={14} />
							</button>
							<Building2 class="w-5 h-5 text-muted-foreground" />
						</div>
					</div>

					{#if client}
						<div>
							<p class="font-semibold text-foreground">{client.name}</p>
						</div>
						{#if contact}
							<div class="border-t border-border pt-4">
								<p class="text-sm text-muted-foreground mb-2">Interlocuteur principal</p>
								<p class="font-medium text-foreground">
									{contact.firstname} {contact.lastname}
								</p>
								{#if contact.position}
									<p class="text-sm text-muted-foreground">{contact.position}</p>
								{/if}
								{#if contact.email}
									<p class="text-sm text-muted-foreground">{contact.email}</p>
								{/if}
								{#if contact.phone}
									<p class="text-sm text-muted-foreground">{contact.phone}</p>
								{/if}
							</div>
						{:else}
							<div class="border-t border-border pt-4">
								<p class="text-sm text-muted-foreground">Aucun interlocuteur</p>
								<button
									type="button"
									onclick={() => (clientDialogOpen = true)}
									class="mt-2 text-xs text-foreground underline hover:no-underline cursor-pointer"
								>
									Ajouter un interlocuteur
								</button>
							</div>
						{/if}
					{:else}
						<p class="text-sm text-muted-foreground">Aucun client associé</p>
					{/if}
				</div>

				<!-- Location -->
				<div
					class="border border-border-card rounded-card p-6 grid gap-4 hover:shadow-card transition-shadow duration-300 w-80"
				>
					<div class="flex justify-between items-center">
						<p class="text-muted-foreground text-sm font-medium">LIEU</p>
						<div class="flex items-center gap-2">
							<button
								onclick={() => (locationDialogOpen = true)}
								class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
								title="Modifier le lieu"
							>
								<Pencil size={14} />
							</button>
							<MapPin class="w-5 h-5 text-muted-foreground" />
						</div>
					</div>
					{#if hasLocation}
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
							<p class="text-sm text-muted-foreground italic">{caseData.locationNotes}</p>
						{/if}
					{:else}
						<p class="text-sm text-muted-foreground">Aucun lieu renseigné</p>
					{/if}
				</div>

				<!-- Case Type -->
				<div
					class="border border-border-card rounded-card p-6 grid gap-4 hover:shadow-card transition-shadow duration-300 w-80"
				>
					<div class="flex justify-between items-center">
						<p class="text-muted-foreground text-sm font-medium">TYPE D'AFFAIRE</p>
						<div class="flex items-center gap-2">
							<button
								onclick={() => (caseTypeDialogOpen = true)}
								class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
								title="Modifier le type"
							>
								<Pencil size={14} />
							</button>
							<Layers class="w-5 h-5 text-muted-foreground" />
						</div>
					</div>

					{#if caseTypeName}
						<span class="bg-success/15 text-success px-2 py-1 rounded-card text-sm font-medium w-fit">
							{caseTypeName}
						</span>
					{:else}
						<p class="text-sm text-muted-foreground">Non défini</p>
					{/if}
				</div>
			</div>
		</div>

		<!-- Case Subject -->
		<div class="border border-border-card rounded-card p-6 grid gap-4">
			<div class="flex justify-between items-center">
				<h3 class="text-lg font-semibold text-foreground">Personne impliquée</h3>
				{#if subject && !showSubjectForm}
					<div class="flex items-center gap-2">
						<button
							onclick={startEditSubject}
							class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
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
								class="p-1.5 rounded btn-ghost-destructive cursor-pointer"
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
				<div class="border border-border rounded-lg p-6 bg-muted/30 max-w-2xl grid gap-4">
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
					<div class="flex justify-end gap-2">
						<button
							type="button"
							onclick={cancelSubjectForm}
							class="h-input rounded-input bg-transparent text-dark hover:bg-muted inline-flex items-center justify-center gap-1 px-4 text-sm font-semibold active:scale-[0.98] border border-border-input cursor-pointer"
						>
							<X size={14} />
							Annuler
						</button>
						<button
							type="button"
							onclick={editingSubject ? saveSubjectEdit : saveSubjectCreate}
							disabled={savingSubject}
							class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-1 px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
						>
							<Check size={14} />
							{savingSubject ? "Enregistrement..." : editingSubject ? "Enregistrer" : "Ajouter"}
						</button>
					</div>
				</div>
			{:else if subject}
				<!-- Display mode: existing subject -->
				<div class="border border-border rounded-lg p-5 bg-muted/30 max-w-2xl flex flex-col gap-5">
					<div class="flex items-start gap-4">
						<div class="w-10 h-10 rounded-full bg-accent text-accent-foreground flex items-center justify-center flex-shrink-0">
							<User size={20} />
						</div>
						<div class="flex-1 flex flex-col gap-1">
							<p class="text-sm font-semibold text-foreground">
								{subject.firstname} {subject.lastname}
							</p>
							{#if subject.occupation}
								<p class="text-sm text-muted-foreground">{subject.occupation}</p>
							{/if}
						</div>
					</div>
					{#if subject.email || subject.phone}
						<div class="grid grid-cols-2 gap-x-8 gap-y-4">
							{#if subject.email}
								<div class="flex flex-col gap-1">
									<p class="text-xs text-muted-foreground">Email</p>
									<p class="text-sm text-foreground">{subject.email}</p>
								</div>
							{/if}
							{#if subject.phone}
								<div class="flex flex-col gap-1">
									<p class="text-xs text-muted-foreground">Téléphone</p>
									<p class="text-sm text-foreground">{subject.phone}</p>
								</div>
							{/if}
						</div>
					{/if}
					{#if subject.address1 || subject.city}
						<div class="flex flex-col gap-1">
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
						<div class="flex flex-col gap-1">
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
						class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
					>
						<UserPlus size={16} />
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

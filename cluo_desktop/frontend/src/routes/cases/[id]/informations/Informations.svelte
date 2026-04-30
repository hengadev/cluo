<script lang="ts">
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { Briefcase, MapPin } from "@lucide/svelte";
	import {
		fetchCase,
		fetchClient,
		fetchContact,
		fetchCaseSubjects
	} from "$lib/services/api";
	import type { Case, Client, Contact, CaseSubject, CaseStatus, SubjectRole } from "$lib/types/entities";

	const STATUS_LABELS: Record<CaseStatus, string> = {
		draft: "Brouillon",
		in_progress: "En cours",
		ready: "Prêt",
		released: "Clôturé"
	};

	const SUBJECT_ROLE_LABELS: Record<SubjectRole, string> = {
		victim: "Victime",
		suspect: "Suspect",
		witness: "Témoin",
		claimant: "Demandeur",
		representative: "Représentant"
	};

	let caseData: Case | null = null;
	let client: Client | null = null;
	let contact: Contact | null = null;
	let subjects: Array<{ subject: CaseSubject; role: string }> = [];
	let loading = true;
	let error: string | null = null;

	const caseId = $derived($page.params.id);

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		if (!caseId) return;

		loading = true;
		error = null;

		try {
			// Fetch case data
			caseData = await fetchCase(caseId);

			if (!caseData) {
				error = "Dossier introuvable";
				loading = false;
				return;
			}

			// Fetch related data in parallel
			const [clientData, contactData, subjectsData] = await Promise.all([
				fetchClient(caseData.clientId),
				caseData.assignedContactId ? fetchContact(caseData.assignedContactId) : Promise.resolve(null),
				fetchCaseSubjects(caseId)
			]);

			client = clientData;
			contact = contactData;
			subjects = subjectsData;
		} catch (e) {
			error = e instanceof Error ? e.message : "Erreur lors du chargement des données";
		} finally {
			loading = false;
		}
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString("fr-FR", {
			day: "2-digit",
			month: "short",
			year: "numeric"
		});
	}

	function getLocationIcon() {
		return MapPin;
	}
</script>

<div class="p-8">
	{#if loading}
		<div class="flex items-center justify-center py-12">
			<p class="text-muted-foreground">Chargement...</p>
		</div>
	{:else if error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
			{error}
		</div>
	{:else if caseData}
		<div class="flex gap-8">
			<!-- Case Header -->
			<div class="grid gap-5 p-6 border border-border-card rounded-card flex-1 animate-fade-in hover:shadow-md transition-shadow duration-300" style="animation-delay: 100ms;">
				<div class="flex gap-4 items-center">
					<span class="bg-blue-100 text-blue-800 px-2 py-1 rounded-card text-sm font-medium">
						STATUT: {STATUS_LABELS[caseData.status] || caseData.status}
					</span>
					<p class="text-muted-foreground text-sm">
						Créé le {formatDate(caseData.createdAt)}
					</p>
				</div>
				<h2 class="text-3xl font-bold text-foreground">{caseData.title}</h2>
				<div class="flex gap-4 text-lg items-center">
					<span class="text-muted-foreground">ID de dossier:</span>
					<span class="font-mono text-foreground">#{caseData.id}</span>
				</div>
				{#if caseData.externalReference}
					<div class="flex gap-4 text-lg items-center">
						<span class="text-muted-foreground">Référence externe:</span>
						<span class="text-foreground">{caseData.externalReference}</span>
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
				<div class="border border-border-card rounded-card p-6 grid gap-4 animate-fade-in hover:shadow-md transition-shadow duration-300 w-80" style="animation-delay: 200ms;">
					<div class="flex justify-between items-center">
						<p class="text-muted-foreground text-sm font-medium">CLIENT</p>
						<Briefcase class="w-5 h-5 text-muted-foreground" />
					</div>
					{#if client}
						<div>
							<p class="font-semibold text-foreground">{client.name}</p>
						</div>
						{#if contact}
							<div class="border-t border-border pt-4 mt-2">
								<p class="text-sm text-muted-foreground mb-2">Contact principal</p>
								<p class="font-medium text-foreground">
									{contact.firstname} {contact.lastname}
								</p>
								<p class="text-sm text-muted-foreground">{contact.position}</p>
								<p class="text-sm text-muted-foreground">{contact.email}</p>
								<p class="text-sm text-muted-foreground">{contact.phone}</p>
							</div>
						{/if}
					{:else}
						<p class="text-sm text-muted-foreground">Aucun client associé</p>
					{/if}
				</div>

				<!-- Location -->
				<div class="border border-border-card rounded-card p-6 grid gap-4 animate-fade-in hover:shadow-md transition-shadow duration-300 w-80" style="animation-delay: 300ms;">
					<div class="flex justify-between items-center">
						<p class="text-muted-foreground text-sm font-medium">LIEU</p>
						<MapPin class="w-5 h-5 text-muted-foreground" />
					</div>
					{#if caseData.placename}
						<p class="font-semibold text-foreground">{caseData.placename}</p>
					{/if}
					<p class="text-sm text-foreground">
						{caseData.address1}
						{#if caseData.address2}<br />{caseData.address2}{/if}
					</p>
					<p class="text-sm text-foreground">
						{caseData.postalCode} {caseData.city}
					</p>
					<p class="text-sm text-muted-foreground">{caseData.country}</p>
					{#if caseData.locationNotes}
						<p class="text-sm text-muted-foreground mt-2 italic">{caseData.locationNotes}</p>
					{/if}
				</div>
			</div>
		</div>

		<!-- Case Subjects -->
		{#if subjects.length > 0}
			<div class="mt-8 border border-border-card rounded-card p-6 animate-fade-in" style="animation-delay: 400ms;">
				<h3 class="text-lg font-semibold text-foreground mb-4">Personnes impliquées</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each subjects as { subject, role }}
						<div class="border border-border rounded-lg p-4 bg-muted/30">
							<div class="flex items-center justify-between mb-2">
								<span class="text-sm font-medium text-foreground">
									{subject.firstName} {subject.lastName}
								</span>
								<span class="px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800">
									{SUBJECT_ROLE_LABELS[role as SubjectRole] || role}
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
							{#if subject.address}
								<p class="text-sm text-muted-foreground">{subject.address}</p>
							{/if}
						</div>
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>

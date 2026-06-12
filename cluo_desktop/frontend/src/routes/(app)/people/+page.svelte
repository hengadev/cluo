<script lang="ts">
	import { onMount } from "svelte";
	import { page } from "$app/stores";
	import { goto } from "$app/navigation";
	import {
		Plus,
		Search,
		Building2,
		User,
		Scale,
		Shield,
		Landmark,
		Users,
		Mail,
		Phone,
		Trash2,
		X,
		Check,
		Pencil,
		UserPlus,
		ExternalLink,
		MapPin,
		RefreshCw,
	} from "@lucide/svelte";
	import Spinner from "$lib/components/Spinner.svelte";
	import EmptyState from "$lib/components/EmptyState.svelte";
	import {
		fetchAllClients,
		deleteClient,
		fetchAllUsers,
		fetchAllCaseSubjects,
		updateCaseSubject,
		deleteCaseSubject,
		fetchAllCases,
		fetchClientContacts,
	} from "$lib/services/api";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import { clientTypeBadge, userRoleBadge } from "$lib/utils/badgeVariants";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import NewClientDialog from "$lib/custom/global/NewClientDialog.svelte";
	import NewSubjectDialog from "$lib/custom/global/NewSubjectDialog.svelte";
	import type { Client, ClientType, AuthUser, CaseSubject, Case } from "$lib/types/entities";

	const toastState = getToastContext();

	// =========================================================================
	// TAB STATE (URL-driven)
	// =========================================================================

	type Tab = "clients" | "subjects" | "users";

	let activeTab = $derived.by<Tab>(() => {
		const t = $page.url.searchParams.get("tab");
		return t === "subjects" ? "subjects" : t === "users" ? "users" : "clients";
	});

	function switchTab(tab: Tab) {
		goto(`/people?tab=${tab}`, { replaceState: true });
	}

	// =========================================================================
	// CLIENTS TAB
	// =========================================================================

	let clients: Client[] = $state([]);
	let contactCounts: Record<string, number> = $state({});
	let newClientOpen = $state(false);
	let clientSearchQuery = $state("");
	let activeClientFilter = $state<ClientType | "all">("all");

	let clientBeingEdited: Client | null = $state(null);
	let editClientOpen = $state(false);

	function startEditClient(c: Client) {
		clientBeingEdited = c;
		editClientOpen = true;
	}

	const CLIENT_TYPE_LABELS: Record<ClientType, string> = {
		person: "Particulier",
		insurance: "Assurance",
		lawyer: "Cabinet juridique",
		company: "Entreprise",
		government: "Administration",
	};

	const CLIENT_TYPE_ICONS: Record<ClientType, typeof User> = {
		person: User,
		insurance: Shield,
		lawyer: Scale,
		company: Building2,
		government: Landmark,
	};

	const CLIENT_FILTERS: { value: ClientType | "all"; label: string }[] = [
		{ value: "all", label: "Tous" },
		{ value: "person", label: "Particulier" },
		{ value: "insurance", label: "Assurance" },
		{ value: "lawyer", label: "Cabinet juridique" },
		{ value: "company", label: "Entreprise" },
		{ value: "government", label: "Administration" },
	];

	let filteredClients = $derived.by(() => {
		let result = clients;
		if (activeClientFilter !== "all") {
			result = result.filter((c) => c.type === activeClientFilter);
		}
		if (clientSearchQuery.trim()) {
			const q = clientSearchQuery.toLowerCase();
			result = result.filter((c) => c.name.toLowerCase().includes(q));
		}
		return result;
	});

	async function handleDeleteClient(client: Client) {
		try {
			await deleteClient(client.id);
			clients = clients.filter((c) => c.id !== client.id);
			toastState.add(TOAST_LEVELS.Info, "Client supprimé", `« ${client.name} » a été supprimé.`);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de supprimer le client",
			);
		}
	}

	// =========================================================================
	// SUBJECTS TAB
	// =========================================================================

	let subjects: CaseSubject[] = $state([]);
	let casesBySubjectId: Record<string, Case> = $state({});
	let subjectsLoading = $state(false);
	let subjectsError: string | null = $state(null);
	let subjectsLoaded = $state(false);
	let subjectSearchQuery = $state("");

	let newSubjectOpen = $state(false);
	let editingSubjectId: string | null = $state(null);
	let savingSubject = $state(false);

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

	let filteredSubjects = $derived.by(() => {
		if (!subjectSearchQuery.trim()) return subjects;
		const q = subjectSearchQuery.toLowerCase();
		return subjects.filter(
			(s) =>
				`${s.firstname} ${s.lastname}`.toLowerCase().includes(q) ||
				(s.occupation?.toLowerCase().includes(q) ?? false),
		);
	});

	$effect(() => {
		if (activeTab === "subjects" && !subjectsLoaded) {
			loadSubjectsData();
		}
	});

	async function loadSubjectsData() {
		subjectsLoading = true;
		subjectsError = null;
		subjectsLoaded = true;
		try {
			const [subjectsData, casesData] = await Promise.all([
				fetchAllCaseSubjects(),
				fetchAllCases({ pageSize: 500 }),
			]);
			subjects = subjectsData;
			const map: Record<string, Case> = {};
			for (const c of casesData.cases) {
				if (c.caseSubjectId) map[c.caseSubjectId] = c;
			}
			casesBySubjectId = map;
		} catch (e) {
			subjectsError = e instanceof Error ? e.message : "Erreur lors du chargement";
			subjectsLoaded = false;
		} finally {
			subjectsLoading = false;
		}
	}

	function startCreateSubject() {
		newSubjectOpen = true;
	}

	function handleSubjectCreated(created: CaseSubject) {
		subjects = [...subjects, created];
	}

	function startEditSubject(s: CaseSubject) {
		editingSubjectId = s.id;
		subjectForm = {
			firstname: s.firstname,
			lastname: s.lastname,
			email: s.email || "",
			phone: s.phone || "",
			address1: s.address1 || "",
			address2: s.address2 || "",
			city: s.city || "",
			postalCode: s.postalCode || "",
			occupation: s.occupation || "",
			notes: s.notes || "",
		};
	}

	function cancelSubjectForm() {
		editingSubjectId = null;
	}



	async function saveSubjectEdit() {
		if (!editingSubjectId) return;
		if (!subjectForm.lastname.trim() || !subjectForm.firstname.trim()) {
			toastState.add(TOAST_LEVELS.Error, "Erreur", "Le nom et le prénom sont requis.");
			return;
		}
		savingSubject = true;
		try {
			const updated = await updateCaseSubject(editingSubjectId, {
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
			subjects = subjects.map((s) => (s.id === editingSubjectId ? updated : s));
			editingSubjectId = null;
			toastState.add(TOAST_LEVELS.Info, "Personne mise à jour", "Les informations ont été enregistrées.");
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

	async function handleDeleteSubject(subject: CaseSubject) {
		try {
			await deleteCaseSubject(subject.id);
			subjects = subjects.filter((s) => s.id !== subject.id);
			const newMap = { ...casesBySubjectId };
			delete newMap[subject.id];
			casesBySubjectId = newMap;
			toastState.add(
				TOAST_LEVELS.Info,
				"Personne supprimée",
				`« ${subject.firstname} ${subject.lastname} » a été supprimée.`,
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de supprimer.",
			);
		}
	}

	// =========================================================================
	// USERS TAB
	// =========================================================================

	let users: AuthUser[] = $state([]);

	const ROLE_LABELS: Record<string, string> = {
		admin: "Administrateur",
		investigator: "Enquêteur",
		viewer: "Lecteur",
	};

	// =========================================================================
	// INITIAL LOAD
	// =========================================================================

	let loading = $state(true);

	onMount(async () => {
		try {
			const [clientsData, usersData] = await Promise.all([
				fetchAllClients(),
				fetchAllUsers(),
			]);
			clients = clientsData;
			users = usersData;
		} catch (_) {
			// individual tab error states handle display
		} finally {
			loading = false;
		}

		// Fire contact count fetches after the client list is already rendered
		Promise.all(
			clients.map(async (c) => {
				try {
					const contacts = await fetchClientContacts(c.id);
					return [c.id, contacts.length] as const;
				} catch {
					return [c.id, 0] as const;
				}
			}),
		).then((entries) => {
			contactCounts = Object.fromEntries(entries);
		});
	});
</script>

<div class="page-content">
	<!-- Page header -->
	<div class="flex items-center justify-between">
		<h1 class="text-3xl font-bold">Personnes</h1>
		{#if activeTab === "clients"}
			<NewClientDialog bind:open={newClientOpen} onSaved={(c) => { clients = [...clients, c]; }}>
				<button
					class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-5 text-[15px] font-semibold active:scale-[0.98] cursor-pointer transition-interactive duration-200"
				>
					<Plus size={18} />
					Nouveau client
				</button>
			</NewClientDialog>
		{:else if activeTab === "subjects"}
			<button
				type="button"
				onclick={startCreateSubject}
				class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-5 text-[15px] font-semibold active:scale-[0.98] cursor-pointer transition-interactive duration-200"
			>
				<Plus size={18} />
				Nouvelle personne
			</button>
		{/if}
	</div>

	<!-- Tab bar -->
	<div class="flex border-b border-border">
		{#each [
			{ key: "clients" as Tab, label: "Clients", count: loading ? null : clients.length },
			{ key: "subjects" as Tab, label: "Personnes impliquées", count: subjectsLoaded ? subjects.length : null },
			{ key: "users" as Tab, label: "Utilisateurs", count: loading ? null : users.length },
		] as tab}
			<button
				type="button"
				onclick={() => switchTab(tab.key)}
				class="px-5 py-3 text-sm font-medium border-b-2 -mb-px transition-colors duration-150 cursor-pointer {activeTab === tab.key
					? 'border-foreground text-foreground'
					: 'border-transparent text-muted-foreground hover:text-foreground hover:border-border-input'}"
			>
				{tab.label}{#if tab.count !== null && tab.count > 0}<span
						class="ml-3 inline-flex items-center justify-center rounded-full px-2 py-0.5 text-xs font-medium tabular-nums {activeTab === tab.key
							? 'bg-foreground/10 text-foreground'
							: 'bg-muted text-muted-foreground'}">{tab.count}</span
					>{/if}
			</button>
		{/each}
	</div>

	<!-- ======================================================================= -->
	<!-- CLIENTS TAB                                                              -->
	<!-- ======================================================================= -->
	{#if activeTab === "clients"}
		{#if loading}
			<div class="flex items-center justify-center py-12">
				<Spinner />
			</div>
		{:else}
			<div class="relative">
				<Search
					size={18}
					class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
				/>
				<input
					type="text"
					placeholder="Rechercher un client..."
					bind:value={clientSearchQuery}
					class="h-input rounded-input border border-border-card bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full pl-10 pr-4 text-sm focus:ring-2 focus:ring-offset-2"
				/>
			</div>

			<div class="flex flex-wrap gap-2">
				{#each CLIENT_FILTERS as chip}
					<button
						type="button"
						onclick={() => (activeClientFilter = chip.value)}
						class="px-3 py-1.5 text-xs font-medium rounded-full border transition-interactive duration-200 cursor-pointer {activeClientFilter === chip.value
							? 'bg-foreground text-background border-foreground'
							: 'bg-background text-muted-foreground border-border-card hover:border-border-input-hover hover:text-foreground'}"
					>
						{chip.label}
					</button>
				{/each}
			</div>

			{#if filteredClients.length === 0}
				{#if clientSearchQuery || activeClientFilter !== "all"}
					<EmptyState icon={Search} message="Aucun client trouvé pour cette recherche." />
				{:else}
					<EmptyState icon={Users} message="Aucun client enregistré.">
						<button
							onclick={() => (newClientOpen = true)}
							class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-5 text-[15px] font-semibold active:scale-[0.98] cursor-pointer transition-interactive duration-200"
						>
							<Plus size={18} />
							Nouveau client
						</button>
					</EmptyState>
				{/if}
			{:else}
				<div class="grid grid-cols-2 lg:grid-cols-3 gap-4">
					{#each filteredClients as client}
						{@const Icon = CLIENT_TYPE_ICONS[client.type] || User}
						{@const badge = clientTypeBadge(client.type)}
						{@const label = CLIENT_TYPE_LABELS[client.type] || client.type}
						<div class="border border-border-card rounded-card p-4 bg-background hover:shadow-card transition-interactive duration-300 group cursor-pointer"
							onclick={() => goto(`/clients/${client.id}`)}
						>
							<div class="flex items-start justify-between">
								<div class="flex items-center gap-3 min-w-0">
									<div class="flex-shrink-0 w-10 h-10 rounded-full bg-muted flex items-center justify-center">
										<Icon size={20} class="text-muted-foreground" />
									</div>
									<div class="min-w-0">
										<h3 class="font-semibold text-foreground truncate">{client.name}</h3>
										<div class="flex items-center gap-2 mt-1">
											<span class="inline-block px-2 py-0.5 text-xs rounded-full {badge}">
												{label}
											</span>
											{#if contactCounts[client.id]}
												<span class="inline-block px-2 py-0.5 text-xs rounded-full bg-muted text-muted-foreground">
													{contactCounts[client.id]} {contactCounts[client.id] === 1 ? 'interlocuteur' : 'interlocuteurs'}
												</span>
											{/if}
										</div>
									</div>
								</div>
								<div class="flex gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity duration-150 flex-shrink-0">
									<button
										onclick={(e) => { e.stopPropagation(); startEditClient(client); }}
										class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
										title="Modifier"
									>
										<Pencil size={13} />
									</button>
									<ConfirmDialog
										title="Supprimer le client"
										description="Voulez-vous vraiment supprimer {client.name} ? Cette action est irréversible."
										onConfirm={() => handleDeleteClient(client)}
									>
										<button onclick={(e) => e.stopPropagation()} class="p-1.5 rounded btn-ghost-destructive cursor-pointer" title="Supprimer">
											<Trash2 size={13} />
										</button>
									</ConfirmDialog>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		{/if}

	<!-- ======================================================================= -->
	<!-- SUBJECTS TAB                                                             -->
	<!-- ======================================================================= -->
	{:else if activeTab === "subjects"}
		{#if subjectsLoading}
			<div class="flex items-center justify-center py-12">
				<Spinner />
			</div>
		{:else if subjectsError}
			<div class="flex flex-col items-center gap-4 py-12">
				<p class="text-sm text-destructive">{subjectsError}</p>
				<button
					type="button"
					onclick={loadSubjectsData}
					class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-input border border-border-input bg-background hover:bg-muted cursor-pointer transition-interactive duration-200"
				>
					<RefreshCw size={14} />
					Réessayer
				</button>
			</div>
		{:else}
			<div class="relative">
				<Search
					size={18}
					class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
				/>
				<input
					type="text"
					placeholder="Rechercher par nom ou profession..."
					bind:value={subjectSearchQuery}
					class="h-input rounded-input border border-border-card bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full pl-10 pr-4 text-sm focus:ring-2 focus:ring-offset-2"
				/>
			</div>

			{#if filteredSubjects.length === 0}
				{#if subjectSearchQuery}
					<EmptyState icon={Search} message="Aucune personne trouvée pour cette recherche." />
				{:else}
					<EmptyState icon={UserPlus} message="Aucune personne impliquée enregistrée.">
						<button
							type="button"
							onclick={startCreateSubject}
							class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-5 text-[15px] font-semibold active:scale-[0.98] cursor-pointer transition-interactive duration-200"
						>
							<Plus size={18} />
							Nouvelle personne
						</button>
					</EmptyState>
				{/if}
			{:else}
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each filteredSubjects as subject}
						{#if editingSubjectId === subject.id}
							<div class="border border-border-card rounded-card p-5 bg-background">
								<p class="text-sm font-medium text-muted-foreground mb-4">Modifier la personne</p>
								<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
								<div>
									<label class="text-xs text-muted-foreground mb-1 block">Prénom *</label>
									<input type="text" bind:value={subjectForm.firstname} placeholder="Prénom" class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2" />
								</div>
								<div>
									<label class="text-xs text-muted-foreground mb-1 block">Nom *</label>
									<input type="text" bind:value={subjectForm.lastname} placeholder="Nom" class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2" />
								</div>
								<div>
									<label class="text-xs text-muted-foreground mb-1 block">Email</label>
									<input type="email" bind:value={subjectForm.email} placeholder="email@example.com" class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2" />
								</div>
								<div>
									<label class="text-xs text-muted-foreground mb-1 block">Téléphone</label>
									<input type="tel" bind:value={subjectForm.phone} placeholder="+33 6 00 00 00 00" class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2" />
								</div>
								<div>
									<label class="text-xs text-muted-foreground mb-1 block">Adresse</label>
									<input type="text" bind:value={subjectForm.address1} placeholder="Adresse ligne 1" class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2" />
								</div>
								<div>
									<label class="text-xs text-muted-foreground mb-1 block">Adresse (suite)</label>
									<input type="text" bind:value={subjectForm.address2} placeholder="Adresse ligne 2" class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2" />
								</div>
								<div>
									<label class="text-xs text-muted-foreground mb-1 block">Code postal</label>
									<input type="text" bind:value={subjectForm.postalCode} placeholder="75001" class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2" />
								</div>
								<div>
									<label class="text-xs text-muted-foreground mb-1 block">Ville</label>
									<input type="text" bind:value={subjectForm.city} placeholder="Paris" class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2" />
								</div>
								<div class="md:col-span-2">
									<label class="text-xs text-muted-foreground mb-1 block">Profession</label>
									<input type="text" bind:value={subjectForm.occupation} placeholder="Profession" class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2" />
								</div>
								<div class="md:col-span-2">
									<label class="text-xs text-muted-foreground mb-1 block">Notes</label>
									<textarea bind:value={subjectForm.notes} placeholder="Notes..." rows="2" class="rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 py-2 text-sm focus:ring-2 focus:ring-offset-2 resize-none"></textarea>
								</div>
							</div>
								<div class="flex justify-end gap-2 mt-4">
									<button
										type="button"
										onclick={cancelSubjectForm}
										class="h-input rounded-input bg-transparent hover:bg-muted inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] border border-border-input cursor-pointer"
									>
										<X size={14} class="mr-1" /> Annuler
									</button>
									<button
										type="button"
										onclick={saveSubjectEdit}
										disabled={savingSubject}
										class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
									>
										<Check size={14} class="mr-1" />
										{savingSubject ? "Enregistrement..." : "Enregistrer"}
									</button>
								</div>
							</div>
						{:else}
							<div class="border border-border-card rounded-card p-4 bg-background hover:shadow-card transition-interactive duration-200 group flex flex-col">
								<!-- Zone 1 · Identity -->
								<div class="flex items-start gap-3.5">
									<div class="w-10 h-10 rounded-full bg-accent-subtle flex items-center justify-center flex-shrink-0 text-xs font-semibold text-accent-subtle-foreground select-none">
										{(subject.firstname?.[0] ?? '').toUpperCase()}{(subject.lastname?.[0] ?? '').toUpperCase()}
									</div>
									<div class="min-w-0 flex-1">
										<p class="font-semibold text-foreground text-sm leading-snug">
											{subject.firstname} {subject.lastname}
										</p>
										{#if subject.occupation}
											<p class="text-xs text-foreground-alt mt-1.5 leading-snug">{subject.occupation}</p>
										{/if}
									</div>
									<div class="flex gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity duration-150 flex-shrink-0 -mt-0.5">
										<button
											onclick={() => startEditSubject(subject)}
											class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
											title="Modifier"
										>
											<Pencil size={13} />
										</button>
										<ConfirmDialog
											title="Supprimer la personne"
											description="Voulez-vous vraiment supprimer {subject.firstname} {subject.lastname} ? Cette action est irréversible."
											onConfirm={() => handleDeleteSubject(subject)}
										>
											<button
												class="p-1.5 rounded btn-ghost-destructive cursor-pointer"
												title="Supprimer"
											>
												<Trash2 size={13} />
											</button>
										</ConfirmDialog>
									</div>
								</div>

								<!-- Zone 2 · Contact details -->
								{#if subject.email || subject.phone || subject.city}
									<div class="mt-6 pt-4 border-t border-border-card flex flex-col gap-3">
										{#if subject.email}
											<span class="text-xs text-foreground-alt inline-flex items-center gap-2 min-w-0">
												<Mail size={12} class="flex-shrink-0 opacity-50" />
												<span class="truncate">{subject.email}</span>
											</span>
										{/if}
										{#if subject.phone}
											<span class="text-xs text-foreground-alt inline-flex items-center gap-2">
												<Phone size={12} class="flex-shrink-0 opacity-50" />{subject.phone}
											</span>
										{/if}
										{#if subject.city}
											<span class="text-xs text-foreground-alt inline-flex items-center gap-2">
												<MapPin size={12} class="flex-shrink-0 opacity-50" />{subject.city}
											</span>
										{/if}
									</div>
								{/if}

								<!-- Zone 3 · Case link -->
								{#if casesBySubjectId[subject.id]}
									{@const linkedCase = casesBySubjectId[subject.id]}
									<button
										type="button"
										onclick={() => goto(`/cases/${linkedCase.id}`)}
										class="w-full flex items-center gap-2 px-3 pt-4 pb-2.5 rounded-input bg-muted hover:bg-surface text-xs transition-interactive duration-150 cursor-pointer text-left mt-6 border-t border-border-card"
									>
										<span class="text-foreground-alt font-medium flex-shrink-0">Affaire</span>
										<span class="truncate text-foreground font-medium flex-1">{linkedCase.title}</span>
										<ExternalLink size={11} class="flex-shrink-0 opacity-50" />
									</button>
								{/if}
							</div>
						{/if}
					{/each}
				</div>
			{/if}
		{/if}

	<!-- ======================================================================= -->
	<!-- USERS TAB                                                                -->
	<!-- ======================================================================= -->
	{:else if activeTab === "users"}
		{#if loading}
			<div class="flex items-center justify-center py-12">
				<Spinner />
			</div>
		{:else if users.length === 0}
			<EmptyState icon={Users} message="Aucun utilisateur trouvé." />
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				{#each users as user}
					{@const badge = userRoleBadge(user.role)}
					{@const roleLabel = ROLE_LABELS[user.role] || user.role}
					<div class="border border-border-card rounded-card p-4 bg-background">
						<div class="flex items-center gap-3">
							<div
								class="flex-shrink-0 w-10 h-10 rounded-full bg-muted flex items-center justify-center"
							>
								<User size={20} class="text-muted-foreground" />
							</div>
							<div class="min-w-0">
								<h3 class="font-semibold text-foreground truncate">{user.email}</h3>
								<span class="inline-block mt-1 px-2 py-0.5 text-xs rounded-full {badge}">
									{roleLabel}
								</span>
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	{/if}
</div>

<NewSubjectDialog bind:open={newSubjectOpen} onCreated={handleSubjectCreated} />
<NewClientDialog
	bind:open={editClientOpen}
	client={clientBeingEdited ?? undefined}
	onSaved={(updated) => { clients = clients.map((c) => (c.id === updated.id ? updated : c)); }}
/>

<script lang="ts">
	import { onMount } from "svelte";
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
	} from "@lucide/svelte";
	import Spinner from "$lib/components/Spinner.svelte";
	import EmptyState from "$lib/components/EmptyState.svelte";
	import { fetchAllClients, deleteClient, fetchAllUsers } from "$lib/services/api";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import { clientTypeBadge, userRoleBadge } from "$lib/utils/badgeVariants";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import NewClientDialog from "$lib/custom/global/NewClientDialog.svelte";
	import type { Client, ClientType, AuthUser, UserRole } from "$lib/types/entities";

	const toastState = getToastContext();

	let newClientOpen = $state(false);

	interface PersonEntry {
		id: string;
		name: string;
		kind: "client" | "user";
		type: ClientType | "user";
		role?: UserRole;
	}

	let clients: Client[] = $state([]);
	let users: AuthUser[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let searchQuery = $state("");
	let activeFilter = $state<ClientType | "user" | "all">("all");

	const TYPE_LABELS: Record<ClientType | "user", string> = {
		person: "Particulier",
		insurance: "Assurance",
		lawyer: "Cabinet juridique",
		company: "Entreprise",
		government: "Administration",
		user: "Utilisateur",
	};

	const TYPE_ICONS: Record<ClientType | "user", typeof User> = {
		person: User,
		insurance: Shield,
		lawyer: Scale,
		company: Building2,
		government: Landmark,
		user: Users,
	};

	const FILTER_CHIPS: { value: ClientType | "user" | "all"; label: string }[] = [
		{ value: "all", label: "Tous" },
		{ value: "person", label: "Particulier" },
		{ value: "insurance", label: "Assurance" },
		{ value: "lawyer", label: "Cabinet juridique" },
		{ value: "company", label: "Entreprise" },
		{ value: "government", label: "Administration" },
		{ value: "user", label: "Utilisateur" },
	];

	let persons: PersonEntry[] = $derived([
		...clients.map((c) => ({
			id: c.id,
			name: c.name,
			kind: "client" as const,
			type: c.type,
		})),
		...users.map((u) => ({
			id: u.id,
			name: u.email,
			kind: "user" as const,
			type: "user" as const,
			role: u.role,
		})),
	]);

	let filteredPersons = $derived.by(() => {
		let result = persons;

		if (activeFilter !== "all") {
			result = result.filter((p) => p.type === activeFilter);
		}

		if (searchQuery.trim() !== "") {
			const q = searchQuery.toLowerCase();
			result = result.filter((p) => p.name.toLowerCase().includes(q));
		}

		return result;
	});

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		loading = true;
		error = null;
		try {
			const [clientsData, usersData] = await Promise.all([
				fetchAllClients(),
				fetchAllUsers(),
			]);
			clients = clientsData;
			users = usersData;
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Erreur lors du chargement des personnes";
		} finally {
			loading = false;
		}
	}

	async function handleDelete(person: PersonEntry) {
		if (person.kind !== "client") return;
		try {
			await deleteClient(person.id);
			clients = clients.filter((c) => c.id !== person.id);
			toastState.add(
				TOAST_LEVELS.Info,
				"Client supprimé",
				`${person.name} a été supprimé.`,
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de supprimer le client",
			);
		}
	}

	function handleClick(person: PersonEntry) {
		if (person.kind === "client") {
			goto(`/clients/${person.id}`);
		}
	}

	function getBadge(person: PersonEntry): string {
		if (person.kind === "user" && person.role) {
			return userRoleBadge(person.role);
		}
		return clientTypeBadge(person.type as ClientType);
	}

	function getBadgeLabel(person: PersonEntry): string {
		if (person.kind === "user" && person.role) {
			const ROLE_LABELS: Record<string, string> = {
				admin: "Administrateur",
				investigator: "Enquêteur",
				viewer: "Lecteur",
			};
			return ROLE_LABELS[person.role] || person.role;
		}
		return TYPE_LABELS[person.type] || person.type;
	}
</script>

<div class="page-content">
	<div class="flex items-center justify-between">
		<h1 class="text-3xl font-bold">Personnes</h1>
		<NewClientDialog bind:open={newClientOpen}>
			<button
				class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-5 text-[15px] font-semibold active:scale-[0.98] cursor-pointer transition-interactive duration-200"
			>
				<Plus size={18} />
				Nouvelle personne
			</button>
		</NewClientDialog>
	</div>

	<!-- Search bar -->
	<div class="relative">
		<Search
			size={18}
			class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
		/>
		<input
			type="text"
			placeholder="Rechercher une personne par nom..."
			bind:value={searchQuery}
			class="h-input rounded-input border border-border-card bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full pl-10 pr-4 text-sm focus:ring-2 focus:ring-offset-2"
		/>
	</div>

	<!-- Filter chips -->
	<div
		class="flex flex-wrap gap-2"
	>
		{#each FILTER_CHIPS as chip}
			<button
				class="px-3 py-1.5 text-xs font-medium rounded-full border transition-interactive duration-200 cursor-pointer {activeFilter === chip.value
					? 'bg-foreground text-background border-foreground'
					: 'bg-background text-muted-foreground border-border-card hover:border-border-input-hover hover:text-foreground'}"
				onclick={() => (activeFilter = chip.value)}
				type="button"
			>
				{chip.label}
			</button>
		{/each}
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<Spinner />
		</div>
	{:else if error}
		<div class="alert-error">
			{error}
		</div>
	{:else if filteredPersons.length === 0}
		{#if searchQuery || activeFilter !== "all"}
			<EmptyState icon={Search} message="Aucune personne trouvée pour cette recherche." />
		{:else}
			<EmptyState icon={Users} message="Aucune personne enregistrée.">
				<button
					class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-5 text-[15px] font-semibold active:scale-[0.98] cursor-pointer transition-interactive duration-200"
					onclick={() => (newClientOpen = true)}
				>
					<Plus size={18} />
					Nouvelle personne
				</button>
			</EmptyState>
		{/if}
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each filteredPersons as person, index}
				{@const Icon = TYPE_ICONS[person.type] || User}
				{@const badge = getBadge(person)}
				{@const badgeLabel = getBadgeLabel(person)}
				{#if person.kind === "client"}
					<!-- Client card — clickable, navigates to detail -->
					<button
						class="border border-border-card rounded-card p-4 bg-background hover:shadow-card transition-interactive duration-300 cursor-pointer group text-left"
						onclick={() => handleClick(person)}
					>
						{@render cardContent(person, Icon, badge, badgeLabel, true)}
					</button>
				{:else}
					<!-- User card — non-interactive display -->
					<div
						class="border border-border-card rounded-card p-4 bg-background transition-interactive duration-300"
					>
						{@render cardContent(person, Icon, badge, badgeLabel, false)}
					</div>
				{/if}
			{/each}
		</div>
	{/if}
</div>

{#snippet cardContent(person: PersonEntry, Icon: typeof User, badge: string, badgeLabel: string, interactive: boolean)}
	<div class="flex items-start justify-between">
		<div class="flex items-center gap-3 min-w-0">
			<div
				class="flex-shrink-0 w-10 h-10 rounded-full bg-muted flex items-center justify-center"
			>
				<Icon size={20} class="text-muted-foreground" />
			</div>
			<div class="min-w-0">
				<h3 class="font-semibold text-foreground truncate">
					{person.name}
				</h3>
				<span
					class="inline-block mt-1 px-2 py-0.5 text-xs rounded-full {badge}"
				>
					{badgeLabel}
				</span>
			</div>
		</div>
		{#if interactive}
			<div
				class="opacity-0 group-hover:opacity-100 transition-opacity"
			>
				<ConfirmDialog
					title="Supprimer le client"
					description="Voulez-vous vraiment supprimer {person.name} ? Cette action est irréversible."
					onConfirm={() => handleDelete(person)}
				>
					<button
						class="p-1.5 rounded btn-ghost-destructive"
						onclick={(e: MouseEvent) => e.stopPropagation()}
						type="button"
					>
						✕
					</button>
				</ConfirmDialog>
			</div>
		{/if}
	</div>
{/snippet}

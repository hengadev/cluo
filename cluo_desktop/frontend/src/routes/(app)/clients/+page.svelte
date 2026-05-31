<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import { Plus, Search, Building2, User, Scale, Shield, Landmark } from "@lucide/svelte";
	import { fetchAllClients, deleteClient } from "$lib/services/api";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import type { Client, ClientType } from "$lib/types/entities";

	const toastState = getToastContext();

	let clients: Client[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let searchQuery = $state("");

	const TYPE_LABELS: Record<ClientType, string> = {
		person: "Particulier",
		insurance: "Assurance",
		lawyer: "Cabinet juridique",
		company: "Entreprise",
		government: "Administration",
	};

	const TYPE_ICONS: Record<ClientType, typeof User> = {
		person: User,
		insurance: Shield,
		lawyer: Scale,
		company: Building2,
		government: Landmark,
	};

	const TYPE_BADGE: Record<ClientType, string> = {
		person: "bg-violet-100 text-violet-800",
		insurance: "bg-blue-100 text-blue-800",
		lawyer: "bg-amber-100 text-amber-800",
		company: "bg-emerald-100 text-emerald-800",
		government: "bg-red-100 text-red-800",
	};

	let filteredClients = $derived(
		searchQuery.trim() === ""
			? clients
			: clients.filter((c) =>
					c.name.toLowerCase().includes(searchQuery.toLowerCase()),
				),
	);

	onMount(async () => {
		await loadClients();
	});

	async function loadClients() {
		loading = true;
		error = null;
		try {
			clients = await fetchAllClients();
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Erreur lors du chargement des clients";
		} finally {
			loading = false;
		}
	}

	async function handleDelete(client: Client) {
		try {
			await deleteClient(client.id);
			clients = clients.filter((c) => c.id !== client.id);
			toastState.add(
				TOAST_LEVELS.Info,
				"Client supprimé",
				`${client.name} a été supprimé.`,
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de supprimer le client",
			);
		}
	}
</script>

<div class="p-8 flex flex-col gap-6">
	<div class="flex items-center justify-between animate-fade-in">
		<h1 class="text-3xl font-bold">Clients</h1>
		<button
			class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-5 text-[15px] font-semibold active:scale-[0.98] cursor-pointer transition-all duration-200"
			onclick={() => goto("/clients/new")}
		>
			<Plus size={18} />
			Nouveau client
		</button>
	</div>

	<!-- Search bar -->
	<div class="relative animate-fade-in" style="animation-delay: 100ms;">
		<Search
			size={18}
			class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground"
		/>
		<input
			type="text"
			placeholder="Rechercher un client par nom..."
			bind:value={searchQuery}
			class="h-input rounded-input border border-border-card bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full pl-10 pr-4 text-sm focus:ring-2 focus:ring-offset-2"
		/>
	</div>

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
	{:else if filteredClients.length === 0}
		<div class="text-center py-12">
			<p class="text-muted-foreground">
				{searchQuery
					? "Aucun client trouvé pour cette recherche."
					: "Aucun client enregistré. Créez votre premier client !"}
			</p>
		</div>
	{:else}
		<div
			class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
		>
			{#each filteredClients as client, index}
				{@const Icon = TYPE_ICONS[client.type] || Building2}
				{@const badge = TYPE_BADGE[client.type] || "bg-gray-100 text-gray-800"}
				<div
					class="border border-border-card rounded-card p-5 bg-background hover:border-border-input-hover hover:shadow-md hover:-translate-y-0.5 transition-all duration-300 animate-fade-in cursor-pointer group"
					style="animation-delay: {200 + index * 50}ms;"
					role="button"
					tabindex="0"
					onclick={() => goto(`/clients/${client.id}`)}
					onkeydown={(e) => e.key === "Enter" && goto(`/clients/${client.id}`)}
				>
					<div class="flex items-start justify-between">
						<div class="flex items-center gap-3 min-w-0">
							<div
								class="flex-shrink-0 w-10 h-10 rounded-full bg-muted flex items-center justify-center"
							>
								<Icon size={20} class="text-muted-foreground" />
							</div>
							<div class="min-w-0">
								<h3
									class="font-semibold text-foreground truncate"
								>
									{client.name}
								</h3>
								<span
									class="inline-block mt-1 px-2 py-0.5 text-xs rounded-full {badge}"
								>
									{TYPE_LABELS[client.type] || client.type}
								</span>
							</div>
						</div>
						<div
							class="opacity-0 group-hover:opacity-100 transition-opacity"
						>
							<ConfirmDialog
								title="Supprimer le client"
								description="Voulez-vous vraiment supprimer {client.name} ? Cette action est irréversible."
								onConfirm={() => handleDelete(client)}
							>
								<button
									class="p-1.5 rounded hover:bg-red-50 text-muted-foreground hover:text-red-600 transition-colors"
									onclick={(e: MouseEvent) => e.stopPropagation()}
									type="button"
								>
									✕
								</button>
							</ConfirmDialog>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

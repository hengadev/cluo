<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import {
		ArrowLeft,
		Plus,
		Pencil,
		Trash2,
		Mail,
		Phone,
		User,
	} from "@lucide/svelte";
	import {
		fetchClient,
		deleteClient,
		fetchClientContacts,
		deleteContact,
	} from "$lib/services/api";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import NewClientDialog from "$lib/custom/global/NewClientDialog.svelte";
	import ContactDialog from "$lib/custom/global/ContactDialog.svelte";
	import type { Client, Contact } from "$lib/types/entities";

	const toastState = getToastContext();

	let client: Client | null = $state(null);
	let contacts: Contact[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	// Client edit dialog
	let editClientOpen = $state(false);

	// Contact create/edit dialog
	let contactDialogOpen = $state(false);
	let contactBeingEdited: Contact | null = $state(null);

	const clientId = $derived($page.params.id);

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		if (!clientId) {
			loading = false;
			return;
		}
		loading = true;
		error = null;
		try {
			const [clientData, contactsData] = await Promise.all([
				fetchClient(clientId),
				fetchClientContacts(clientId),
			]);
			client = clientData;
			contacts = contactsData;
		} catch (e) {
			error = e instanceof Error ? e.message : "Erreur lors du chargement";
		} finally {
			loading = false;
		}
	}

	async function handleDeleteClient() {
		if (!client) return;
		try {
			await deleteClient(client.id);
			toastState.add(
				TOAST_LEVELS.Info,
				"Client supprimé",
				`« ${client.name} » a été supprimé.`,
			);
			goto("/people");
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de supprimer",
			);
		}
	}

	function openCreateContact() {
		contactBeingEdited = null;
		contactDialogOpen = true;
	}

	function openEditContact(c: Contact) {
		contactBeingEdited = c;
		contactDialogOpen = true;
	}

	async function handleDeleteContact(contactId: string) {
		try {
			await deleteContact(contactId);
			contacts = contacts.filter((c) => c.id !== contactId);
			toastState.add(TOAST_LEVELS.Info, "Interlocuteur supprimé", "L'interlocuteur a été supprimé.");
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de supprimer",
			);
		}
	}

	const TYPE_LABELS: Record<string, string> = {
		person: "Particulier",
		insurance: "Assurance",
		lawyer: "Cabinet juridique",
		company: "Entreprise",
		government: "Administration",
	};
</script>

<div class="page-content">
	<button
		class="inline-flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
		onclick={() => goto("/people")}
	>
		<ArrowLeft size={18} />
		<span class="text-sm font-medium">Retour aux personnes</span>
	</button>

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<p class="text-muted-foreground">Chargement...</p>
		</div>
	{:else if error}
		<div class="alert-error">{error}</div>
	{:else if client}
		<!-- Client header -->
		<div class="border border-border-card rounded-card p-6 animate-fade-in hover:shadow-card transition-shadow duration-300 self-center">
			<div class="flex items-start gap-6">
				<div>
					<div class="flex items-center gap-4">
						<h1 class="text-3xl font-bold">{client.name}</h1>
						<span class="px-2 py-0.5 text-xs rounded-full bg-accent-subtle text-accent-subtle-foreground">
							{TYPE_LABELS[client.type] || client.type}
						</span>
					</div>
					<p class="text-sm text-muted-foreground mt-2">ID : {client.id}</p>
				</div>
				<div class="flex gap-2">
					<button
						onclick={() => (editClientOpen = true)}
						class="p-2 rounded-lg hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
						title="Modifier"
					>
						<Pencil size={18} />
					</button>
					<ConfirmDialog
						title="Supprimer le client"
						description="Voulez-vous vraiment supprimer {client.name} ? Cette action est irréversible."
						onConfirm={handleDeleteClient}
					>
						<button class="p-2 rounded-lg btn-ghost-destructive cursor-pointer" title="Supprimer">
							<Trash2 size={18} />
						</button>
					</ConfirmDialog>
				</div>
			</div>
		</div>

		<!-- Interlocuteurs section -->
		<div class="animate-fade-in flex flex-col gap-4" style="animation-delay: 150ms;">
			<div class="flex items-center justify-between">
				<h2 class="text-xl font-semibold">
					Interlocuteurs
					<span class="text-muted-foreground text-base font-normal ml-2">({contacts.length})</span>
				</h2>
				<button
					onclick={openCreateContact}
					class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-5 text-[15px] font-semibold active:scale-[0.98] cursor-pointer transition-interactive duration-200"
				>
					<Plus size={16} />
					Ajouter un interlocuteur
				</button>
			</div>

			{#if contacts.length === 0}
				<div class="text-center py-8 text-muted-foreground">
					<User size={32} class="mx-auto mb-2 opacity-50" />
					<p>Aucun interlocuteur enregistré pour ce client.</p>
				</div>
			{:else}
				<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each contacts as contact}
						<div class="border border-border-card rounded-card p-5 bg-background hover:shadow-card transition-shadow duration-200 group flex flex-col gap-3">
							<!-- Identity row -->
							<div class="flex items-start justify-between gap-2">
								<div class="flex items-center gap-3 min-w-0">
									<div class="w-9 h-9 rounded-full bg-accent-subtle flex items-center justify-center flex-shrink-0 text-xs font-semibold text-accent-subtle-foreground select-none">
										{(contact.firstname[0] ?? '').toUpperCase()}{(contact.lastname[0] ?? '').toUpperCase()}
									</div>
									<div class="min-w-0">
										<p class="font-semibold text-foreground text-sm leading-snug truncate">
											{contact.firstname} {contact.lastname}
										</p>
										{#if contact.position}
											<p class="text-xs text-muted-foreground mt-0.5 truncate">{contact.position}</p>
										{/if}
									</div>
								</div>
								<div class="flex gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity duration-150 flex-shrink-0">
									<button
										onclick={() => openEditContact(contact)}
										class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
										title="Modifier"
									>
										<Pencil size={13} />
									</button>
									<ConfirmDialog
										title="Supprimer l'interlocuteur"
										description="Voulez-vous vraiment supprimer {contact.firstname} {contact.lastname} ?"
										onConfirm={() => handleDeleteContact(contact.id)}
									>
										<button class="p-1.5 rounded btn-ghost-destructive cursor-pointer" title="Supprimer">
											<Trash2 size={13} />
										</button>
									</ConfirmDialog>
								</div>
							</div>
							<!-- Contact details -->
							{#if contact.email || contact.phone}
								<div class="pt-3 border-t border-border-card flex flex-col gap-1.5">
									{#if contact.email}
										<span class="text-xs text-foreground-alt inline-flex items-center gap-2 min-w-0">
											<Mail size={12} class="flex-shrink-0 opacity-60" />
											<span class="truncate">{contact.email}</span>
										</span>
									{/if}
									{#if contact.phone}
										<span class="text-xs text-foreground-alt inline-flex items-center gap-2">
											<Phone size={12} class="flex-shrink-0 opacity-60" />{contact.phone}
										</span>
									{/if}
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>

<NewClientDialog
	bind:open={editClientOpen}
	client={client ?? undefined}
	onSaved={(updated) => { client = updated; }}
/>

<ContactDialog
	bind:open={contactDialogOpen}
	contact={contactBeingEdited ?? undefined}
	clientId={clientId}
	onSaved={(saved) => {
		if (contactBeingEdited) {
			contacts = contacts.map((c) => (c.id === saved.id ? saved : c));
		} else {
			contacts = [...contacts, saved];
		}
		contactBeingEdited = null;
	}}
/>

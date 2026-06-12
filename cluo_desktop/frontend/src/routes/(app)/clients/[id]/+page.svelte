<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import {
		ArrowLeft,
		Plus,
		Pencil,
		Check,
		X,
		Trash2,
		Mail,
		Phone,
		Briefcase,
		User,
	} from "@lucide/svelte";
	import {
		fetchClient,
		updateClient,
		deleteClient,
		fetchClientContacts,
		createContact,
		updateContact,
		deleteContact,
		ConflictError,
	} from "$lib/services/api";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import type { Client, ClientType, Contact } from "$lib/types/entities";

	const toastState = getToastContext();

	let client: Client | null = $state(null);
	let contacts: Contact[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	// Edit client state
	let editing = $state(false);
	let editName = $state("");
	let editType = $state<ClientType>("company");
	let saving = $state(false);

	// New contact state
	let showNewContact = $state(false);
	let newFirstname = $state("");
	let newLastname = $state("");
	let newEmail = $state("");
	let newPhone = $state("");
	let newPosition = $state("");

	// Edit contact state
	let editingContactId = $state<string | null>(null);
	let editContactFirstname = $state("");
	let editContactLastname = $state("");
	let editContactEmail = $state("");
	let editContactPhone = $state("");
	let editContactPosition = $state("");

	const clientId = $derived($page.params.id);

	const TYPE_OPTIONS: { value: ClientType; label: string }[] = [
		{ value: "person", label: "Particulier" },
		{ value: "company", label: "Entreprise" },
		{ value: "lawyer", label: "Cabinet juridique" },
		{ value: "insurance", label: "Assurance" },
		{ value: "government", label: "Administration" },
	];

	const TYPE_LABELS: Record<ClientType, string> = {
		person: "Particulier",
		insurance: "Assurance",
		lawyer: "Cabinet juridique",
		company: "Entreprise",
		government: "Administration",
	};

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
			if (client) {
				editName = client.name;
				editType = client.type;
			}
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Erreur lors du chargement";
		} finally {
			loading = false;
		}
	}

	// --- Client edit ---

	function startEditing() {
		if (!client) return;
		editName = client.name;
		editType = client.type;
		editing = true;
	}

	function cancelEditing() {
		if (!client) return;
		editName = client.name;
		editType = client.type;
		editing = false;
	}

	async function saveClient() {
		if (!client) return;
		const trimmed = editName.trim();
		if (!trimmed) {
			toastState.add(
				TOAST_LEVELS.Warning,
				"Champ requis",
				"Le nom du client est obligatoire.",
			);
			return;
		}
		saving = true;
		try {
			client = await updateClient(client.id, {
				name: trimmed,
				type: editType,
			});
			editing = false;
			toastState.add(
				TOAST_LEVELS.Info,
				"Client mis à jour",
				"Les modifications ont été enregistrées.",
			);
		} catch (e) {
			if (e instanceof ConflictError) {
				toastState.add(
					TOAST_LEVELS.Error,
					"Conflit",
					"Un client avec ce nom existe déjà.",
				);
			} else {
				toastState.add(
					TOAST_LEVELS.Error,
					"Erreur",
					e instanceof Error ? e.message : "Erreur lors de la mise à jour",
				);
			}
		} finally {
			saving = false;
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

	// --- Contacts ---

	function resetNewContact() {
		newFirstname = "";
		newLastname = "";
		newEmail = "";
		newPhone = "";
		newPosition = "";
	}

	async function handleCreateContact(e: SubmitEvent) {
		e.preventDefault();
		if (!client) return;
		const trimmedLast = newLastname.trim();
		const trimmedFirst = newFirstname.trim();
		if (!trimmedLast || !trimmedFirst) {
			toastState.add(
				TOAST_LEVELS.Warning,
				"Champs requis",
				"Le nom et le prénom du contact sont obligatoires.",
			);
			return;
		}
		saving = true;
		try {
			const contact = await createContact({
				clientID: client.id,
				lastname: trimmedLast,
				firstname: trimmedFirst,
				email: newEmail.trim(),
				phone: newPhone.trim(),
				position: newPosition.trim(),
			});
			contacts = [...contacts, contact];
			showNewContact = false;
			resetNewContact();
			toastState.add(
				TOAST_LEVELS.Info,
				"Contact ajouté",
				`« ${trimmedFirst} ${trimmedLast} » a été ajouté.`,
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible d'ajouter le contact",
			);
		} finally {
			saving = false;
		}
	}

	function startEditingContact(c: Contact) {
		editingContactId = c.id;
		editContactFirstname = c.firstname;
		editContactLastname = c.lastname;
		editContactEmail = c.email;
		editContactPhone = c.phone;
		editContactPosition = c.position;
	}

	function cancelEditingContact() {
		editingContactId = null;
	}

	async function saveContact(contactId: string) {
		const trimmedLast = editContactLastname.trim();
		const trimmedFirst = editContactFirstname.trim();
		if (!trimmedLast || !trimmedFirst) {
			toastState.add(
				TOAST_LEVELS.Warning,
				"Champs requis",
				"Le nom et le prénom sont obligatoires.",
			);
			return;
		}
		saving = true;
		try {
			const updated = await updateContact(contactId, {
				lastname: trimmedLast,
				firstname: trimmedFirst,
				email: editContactEmail.trim(),
				phone: editContactPhone.trim(),
				position: editContactPosition.trim(),
			});
			contacts = contacts.map((c) => (c.id === contactId ? updated : c));
			editingContactId = null;
			toastState.add(
				TOAST_LEVELS.Info,
				"Contact mis à jour",
				"Les modifications ont été enregistrées.",
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de mettre à jour",
			);
		} finally {
			saving = false;
		}
	}

	async function handleDeleteContact(contactId: string) {
		try {
			await deleteContact(contactId);
			contacts = contacts.filter((c) => c.id !== contactId);
			toastState.add(
				TOAST_LEVELS.Info,
				"Contact supprimé",
				"Le contact a été supprimé.",
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de supprimer",
			);
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
		<div
			class="alert-error"
		>
			{error}
		</div>
	{:else if client}
		<!-- Client header -->
		<div
			class="border border-border-card rounded-card p-6 animate-fade-in hover:shadow-card transition-shadow duration-300"
		>
			<div class="flex items-start justify-between">
				<div class="flex-1">
					{#if editing}
						<form
							onsubmit={(e) => { e.preventDefault(); saveClient(); }}
							class="flex flex-wrap items-end gap-4"
						>
							<div class="flex flex-col gap-2 flex-1 min-w-[200px]">
								<label for="edit-name" class="text-sm font-medium"
									>Nom</label
								>
								<input
									id="edit-name"
									type="text"
									bind:value={editName}
									class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2"
								/>
							</div>
							<div class="flex flex-col gap-2">
								<label for="edit-type" class="text-sm font-medium"
									>Type</label
								>
								<select
									id="edit-type"
									bind:value={editType}
									class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden px-4 text-sm focus:ring-2 focus:ring-offset-2 cursor-pointer"
								>
									{#each TYPE_OPTIONS as opt}
										<option value={opt.value}>{opt.label}</option>
									{/each}
								</select>
							</div>
							<div class="flex gap-2">
								<button
									type="submit"
									disabled={saving}
									class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
								>
									<Check size={16} />
								</button>
								<button
									type="button"
									onclick={cancelEditing}
									class="h-input rounded-input bg-muted text-foreground inline-flex items-center justify-center px-4 text-sm active:scale-[0.98] cursor-pointer"
								>
									<X size={16} />
								</button>
							</div>
						</form>
					{:else}
						<div class="flex items-center gap-4">
							<h1 class="text-3xl font-bold">{client.name}</h1>
							<span
								class="px-2 py-0.5 text-xs rounded-full bg-accent-subtle text-accent-subtle-foreground"
							>
								{TYPE_LABELS[client.type] || client.type}
							</span>
						</div>
						<p class="text-sm text-muted-foreground mt-2">
							ID : {client.id}
						</p>
					{/if}
				</div>
				{#if !editing}
					<div class="flex gap-2">
						<button
							onclick={startEditing}
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
							<button
								class="p-2 rounded-lg btn-ghost-destructive cursor-pointer"
								title="Supprimer"
							>
								<Trash2 size={18} />
							</button>
						</ConfirmDialog>
					</div>
				{/if}
			</div>
		</div>

		<!-- Contacts section -->
		<div class="animate-fade-in" style="animation-delay: 150ms;">
			<div class="flex items-center justify-between mb-4">
				<h2 class="text-xl font-semibold">
					Contacts
					<span class="text-muted-foreground text-base font-normal ml-2"
						>({contacts.length})</span
					>
				</h2>
				{#if !showNewContact}
					<button
						onclick={() => (showNewContact = true)}
						class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-input border border-border-input bg-background hover:bg-muted active:scale-[0.98] cursor-pointer transition-interactive duration-200"
					>
						<Plus size={16} />
						Ajouter un contact
					</button>
				{/if}
			</div>

			<!-- New contact form -->
			{#if showNewContact}
				<form
					onsubmit={handleCreateContact}
					class="border border-accent/50 rounded-card p-5 mb-4 bg-accent/30"
				>
					<h3 class="text-sm font-semibold mb-4">Nouveau contact</h3>
					<div
						class="grid grid-cols-1 md:grid-cols-2 gap-4"
					>
						<div class="flex flex-col gap-1">
							<label for="new-firstname" class="text-xs text-muted-foreground">Prénom *</label>
							<input
								id="new-firstname"
								type="text"
								bind:value={newFirstname}
								required
								class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
						<div class="flex flex-col gap-1">
							<label for="new-lastname" class="text-xs text-muted-foreground">Nom *</label>
							<input
								id="new-lastname"
								type="text"
								bind:value={newLastname}
								required
								class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
						<div class="flex flex-col gap-1">
							<label for="new-email" class="text-xs text-muted-foreground flex items-center gap-1"><Mail size={12} /> Email</label>
							<input
								id="new-email"
								type="email"
								bind:value={newEmail}
								class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
						<div class="flex flex-col gap-1">
							<label for="new-phone" class="text-xs text-muted-foreground flex items-center gap-1"><Phone size={12} /> Téléphone</label>
							<input
								id="new-phone"
								type="tel"
								bind:value={newPhone}
								class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
						<div class="flex flex-col gap-1 md:col-span-2">
							<label for="new-position" class="text-xs text-muted-foreground flex items-center gap-1"><Briefcase size={12} /> Poste</label>
							<input
								id="new-position"
								type="text"
								bind:value={newPosition}
								class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
							/>
						</div>
					</div>
					<div class="flex justify-end gap-2 mt-4">
						<button
							type="button"
							onclick={() => { showNewContact = false; resetNewContact(); }}
							class="h-input rounded-input bg-transparent text-dark hover:bg-muted inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] border-2 border-border-input cursor-pointer"
						>
							Annuler
						</button>
						<button
							type="submit"
							disabled={saving}
							class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
						>
							{saving ? "Enregistrement..." : "Ajouter"}
						</button>
					</div>
				</form>
			{/if}

			<!-- Contact list -->
			{#if contacts.length === 0 && !showNewContact}
				<div class="text-center py-8 text-muted-foreground">
					<User size={32} class="mx-auto mb-2 opacity-50" />
					<p>Aucun contact enregistré pour ce client.</p>
				</div>
			{:else}
				<div class="grid gap-3">
					{#each contacts as contact}
						<div
							class="border border-border-card rounded-card p-4 bg-background hover:shadow-mini transition-shadow duration-200"
						>
							{#if editingContactId === contact.id}
								<!-- Inline edit form -->
								<form
									onsubmit={(e) => { e.preventDefault(); saveContact(contact.id); }}
									class="grid grid-cols-1 md:grid-cols-2 gap-3"
								>
									<div class="flex flex-col gap-1">
										<label class="text-xs text-muted-foreground">Prénom</label>
										<input
											type="text"
											bind:value={editContactFirstname}
											class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
										/>
									</div>
									<div class="flex flex-col gap-1">
										<label class="text-xs text-muted-foreground">Nom</label>
										<input
											type="text"
											bind:value={editContactLastname}
											class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
										/>
									</div>
									<div class="flex flex-col gap-1">
										<label class="text-xs text-muted-foreground">Email</label>
										<input
											type="email"
											bind:value={editContactEmail}
											class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
										/>
									</div>
									<div class="flex flex-col gap-1">
										<label class="text-xs text-muted-foreground">Téléphone</label>
										<input
											type="tel"
											bind:value={editContactPhone}
											class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
										/>
									</div>
									<div class="flex flex-col gap-1 md:col-span-2">
										<label class="text-xs text-muted-foreground">Poste</label>
										<input
											type="text"
											bind:value={editContactPosition}
											class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 text-sm focus:ring-2 focus:ring-offset-2"
										/>
									</div>
									<div class="md:col-span-2 flex justify-end gap-2">
										<button
											type="button"
											onclick={cancelEditingContact}
											class="h-input rounded-input bg-transparent text-dark hover:bg-muted inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] border-2 border-border-input cursor-pointer"
										>
											Annuler
										</button>
										<button
											type="submit"
											disabled={saving}
											class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
										>
											<Check size={16} class="mr-1" />
											Enregistrer
										</button>
									</div>
								</form>
							{:else}
								<!-- Display mode -->
								<div class="flex items-start justify-between">
									<div class="flex items-start gap-4">
										<div
											class="w-10 h-10 rounded-full bg-muted flex items-center justify-center flex-shrink-0"
										>
											<User size={18} class="text-muted-foreground" />
										</div>
										<div>
											<p class="font-semibold text-foreground">
												{contact.firstname} {contact.lastname}
											</p>
											{#if contact.position}
												<p class="text-sm text-muted-foreground">
													{contact.position}
												</p>
											{/if}
											<div
												class="flex flex-wrap gap-x-4 gap-y-1 mt-1"
											>
												{#if contact.email}
													<span
														class="text-sm text-muted-foreground inline-flex items-center gap-1"
													>
														<Mail size={12} />
														{contact.email}
													</span>
												{/if}
												{#if contact.phone}
													<span
														class="text-sm text-muted-foreground inline-flex items-center gap-1"
													>
														<Phone size={12} />
														{contact.phone}
													</span>
												{/if}
											</div>
										</div>
									</div>
									<div class="flex gap-1">
										<button
											onclick={() => startEditingContact(contact)}
											class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
											title="Modifier"
										>
											<Pencil size={14} />
										</button>
										<ConfirmDialog
											title="Supprimer le contact"
											description="Voulez-vous vraiment supprimer {contact.firstname} {contact.lastname} ?"
											onConfirm={() => handleDeleteContact(contact.id)}
										>
											<button
												class="p-1.5 rounded btn-ghost-destructive cursor-pointer"
												title="Supprimer"
											>
												<Trash2 size={14} />
											</button>
										</ConfirmDialog>
									</div>
								</div>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>

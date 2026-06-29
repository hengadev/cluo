<script lang="ts">
    import { Dialog, Label, Separator } from "bits-ui";
    import { X, Loader2, UserPlus, Pencil } from "@lucide/svelte";
    import { fetchAllClients, fetchClientContacts, updateCase, fetchClient, fetchContact } from "$lib/services/api";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import ContactDialog from "$lib/custom/global/ContactDialog.svelte";
    import type { Client, Contact } from "$lib/types/entities";

    const toastState = getToastContext();

    type Props = {
        open?: boolean;
        caseId: string;
        currentClientId: string;
        currentContactId: string | null;
        onSaved?: (client: Client, contact: Contact | null) => void;
    };
    let { open = $bindable(false), caseId, currentClientId, currentContactId, onSaved }: Props = $props();

    let allClients: Client[] = $state([]);
    let contactsForClient: Contact[] = $state([]);
    let selectedClientId = $state("");
    let selectedContactId = $state("");
    let clientSearchQuery = $state("");
    let saving = $state(false);
    let loadingClients = $state(false);
    let loadingContacts = $state(false);
    let contactDialogOpen = $state(false);
    let contactToEdit = $state<Contact | undefined>(undefined);

    let filteredClients = $derived(
        clientSearchQuery.trim() === ""
            ? allClients
            : allClients.filter((c) =>
                  c.name.toLowerCase().includes(clientSearchQuery.toLowerCase()),
              ),
    );

    let selectedContact = $derived(contactsForClient.find((c) => c.id === selectedContactId));

    $effect(() => {
        if (open) {
            selectedClientId = currentClientId;
            selectedContactId = currentContactId ?? "";
            clientSearchQuery = "";
            init();
        }
    });

    async function init() {
        if (allClients.length === 0) {
            loadingClients = true;
            try {
                allClients = await fetchAllClients();
            } catch {
                toastState.add(TOAST_LEVELS.Error, "Erreur", "Impossible de charger la liste des clients.");
            } finally {
                loadingClients = false;
            }
        }
        if (currentClientId) {
            await loadContactsForClient(currentClientId);
        }
    }

    async function loadContactsForClient(clientId: string) {
        if (!clientId) {
            contactsForClient = [];
            return;
        }
        loadingContacts = true;
        contactsForClient = [];
        try {
            contactsForClient = await fetchClientContacts(clientId);
        } catch {
            contactsForClient = [];
        } finally {
            loadingContacts = false;
        }
    }

    async function handleSubmit() {
        if (!selectedClientId) return;
        saving = true;
        try {
            await updateCase(caseId, {
                clientId: selectedClientId,
                assignedContactID: selectedContactId || undefined,
            });
            const [client, contact] = await Promise.all([
                fetchClient(selectedClientId),
                selectedContactId ? fetchContact(selectedContactId) : Promise.resolve(null),
            ]);
            toastState.add(TOAST_LEVELS.Info, "Dossier mis à jour", "Le client et le contact ont été mis à jour.");
            open = false;
            onSaved?.(client, contact);
        } catch (e) {
            toastState.add(TOAST_LEVELS.Error, "Erreur", e instanceof Error ? e.message : "Impossible de mettre à jour.");
        } finally {
            saving = false;
        }
    }

    function openCreateContact() {
        contactToEdit = undefined;
        contactDialogOpen = true;
    }

    function openEditContact() {
        contactToEdit = selectedContact;
        contactDialogOpen = true;
    }

    async function onContactSaved(contact: Contact) {
        await loadContactsForClient(selectedClientId);
        selectedContactId = contact.id;
    }

    const inputClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors";
    const selectClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors cursor-pointer";
</script>

<ContactDialog
    bind:open={contactDialogOpen}
    contact={contactToEdit}
    clientId={selectedClientId}
    onSaved={onContactSaved}
/>

<Dialog.Root bind:open>
    <Dialog.Portal>
        <Dialog.Overlay
            class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
        />
        <Dialog.Content
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border flex flex-col sm:max-w-[480px] md:w-full"
        >
            <div class="flex-shrink-0 px-8 pt-8 pb-6">
                <Dialog.Title class="text-base font-semibold tracking-tight">
                    Client et interlocuteur
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt text-sm mt-1">
                    Modifiez le client associé et son interlocuteur principal.
                </Dialog.Description>
            </div>

            <Separator.Root class="bg-border-input mx-0 !m-0 block h-px flex-shrink-0" />

            <div class="px-8 py-6">
                <form
                    class="flex flex-col gap-4"
                    onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}
                >
                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="client-search" class="text-sm font-medium">
                            Client
                        </Label.Root>
                        <input
                            id="client-search"
                            type="text"
                            placeholder="Rechercher un client..."
                            bind:value={clientSearchQuery}
                            class={inputClass}
                        />
                        {#if loadingClients}
                            <p class="text-xs text-muted-foreground">Chargement…</p>
                        {:else}
                            <select
                                bind:value={selectedClientId}
                                onchange={() => {
                                    selectedContactId = "";
                                    loadContactsForClient(selectedClientId);
                                }}
                                class={selectClass}
                            >
                                <option value="">— Sélectionner un client —</option>
                                {#each filteredClients as c}
                                    <option value={c.id}>{c.name}</option>
                                {/each}
                            </select>
                        {/if}
                    </div>

                    {#if selectedClientId}
                        <div class="flex flex-col gap-1.5">
                            <div class="flex items-center justify-between">
                                <Label.Root for="contact-select" class="text-sm font-medium">
                                    Interlocuteur
                                </Label.Root>
                                <div class="flex items-center gap-1.5">
                                    {#if selectedContactId && selectedContact}
                                        <button
                                            type="button"
                                            onclick={openEditContact}
                                            class="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                                        >
                                            <Pencil size={11} />
                                            Modifier
                                        </button>
                                    {/if}
                                    <button
                                        type="button"
                                        onclick={openCreateContact}
                                        class="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                                    >
                                        <UserPlus size={11} />
                                        Nouveau
                                    </button>
                                </div>
                            </div>

                            {#if loadingContacts}
                                <p class="text-xs text-muted-foreground">Chargement des interlocuteurs…</p>
                            {:else if contactsForClient.length > 0}
                                <select
                                    id="contact-select"
                                    bind:value={selectedContactId}
                                    class={selectClass}
                                >
                                    <option value="">— Aucun interlocuteur —</option>
                                    {#each contactsForClient as c}
                                        <option value={c.id}>{c.firstname} {c.lastname}</option>
                                    {/each}
                                </select>
                            {:else}
                                <p class="text-sm text-muted-foreground">
                                    Aucun interlocuteur pour ce client.
                                    <button
                                        type="button"
                                        onclick={openCreateContact}
                                        class="text-foreground underline cursor-pointer hover:no-underline"
                                    >
                                        En créer un
                                    </button>
                                </p>
                            {/if}
                        </div>
                    {/if}

                    <div class="flex justify-end gap-2 pt-2">
                        <Dialog.Close
                            class="h-input rounded-input bg-transparent text-foreground border border-border-input hover:bg-muted inline-flex items-center justify-center px-5 text-sm font-semibold active:scale-[0.98] cursor-pointer transition-colors"
                        >
                            Annuler
                        </Dialog.Close>
                        <button
                            type="submit"
                            disabled={saving || !selectedClientId}
                            class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/90 inline-flex items-center justify-center gap-2 px-6 text-sm font-semibold active:scale-[0.98] disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer transition-interactive"
                        >
                            {#if saving}
                                <Loader2 size={14} class="animate-spin" />
                            {/if}
                            Enregistrer
                        </button>
                    </div>
                </form>
            </div>

            <Dialog.Close
                class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-6 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            >
                <div>
                    <X class="text-foreground-alt size-4" />
                    <span class="sr-only">Close</span>
                </div>
            </Dialog.Close>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

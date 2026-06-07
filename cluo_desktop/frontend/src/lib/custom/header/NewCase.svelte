<script lang="ts">
    import { Dialog, Label, Separator, Collapsible } from "bits-ui";
    import { X, Loader2, ChevronRight } from "@lucide/svelte";
    import { createCase, fetchAllClients, fetchAllCaseTypes, fetchClientContacts, fetchAllCaseSubjects } from "$lib/services/api";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import { goto } from "$app/navigation";
    import type { CaseStatus, CaseType, CaseSubject, Client, Contact, LocationType } from "$lib/types/entities";

    const toastState = getToastContext();

    type Props = { children?: import("svelte").Snippet; open?: boolean };
    let { children, open = $bindable(false) }: Props = $props();
    let loading = $state(false);

    // Form state
    let title = $state("");
    let description = $state("");
    let clientId = $state("");
    let caseTypeId = $state("");
    let status: CaseStatus = $state("in_progress");
    let isHistorical = $state(false);

    // Additional details state
    let additionalExpanded = $state(false);
    let assignedContactId = $state("");
    let caseSubjectId = $state("");
    let externalReference = $state("");

    // Location state
    let locationExpanded = $state(false);
    let placename = $state("");
    let address1 = $state("");
    let address2 = $state("");
    let city = $state("");
    let postalCode = $state("");
    let country = $state("");
    let locationType: LocationType | "" = $state("");
    let latitude = $state("");
    let longitude = $state("");
    let locationNotes = $state("");

    // Reference data
    let clients: Client[] = $state([]);
    let caseTypes: CaseType[] = $state([]);
    let caseSubjects: CaseSubject[] = $state([]);
    let clientContacts: Contact[] = $state([]);
    let loadingRefData = $state(false);
    let refDataLoaded = $state(false);
    let loadingContacts = $state(false);

    // Computed: count of filled additional fields
    let additionalFilledCount = $derived(
        (assignedContactId ? 1 : 0) +
        (caseSubjectId ? 1 : 0) +
        (externalReference.trim() ? 1 : 0)
    );

    // Computed: count of filled location fields
    let locationFilledCount = $derived(
        (placename.trim() ? 1 : 0) +
        (address1.trim() ? 1 : 0) +
        (address2.trim() ? 1 : 0) +
        (city.trim() ? 1 : 0) +
        (postalCode.trim() ? 1 : 0) +
        (country.trim() ? 1 : 0) +
        (locationType ? 1 : 0) +
        (latitude.trim() ? 1 : 0) +
        (longitude.trim() ? 1 : 0) +
        (locationNotes.trim() ? 1 : 0)
    );

    $effect(() => {
        if (open && !refDataLoaded && !loadingRefData) {
            loadReferenceData();
        }
    });

    // Reload contacts when client changes
    $effect(() => {
        const currentClientId = clientId;
        assignedContactId = "";
        clientContacts = [];
        if (currentClientId) {
            loadClientContacts(currentClientId);
        }
    });

    async function loadReferenceData() {
        loadingRefData = true;
        try {
            const [clientsData, caseTypesData, caseSubjectsData] = await Promise.all([
                fetchAllClients(),
                fetchAllCaseTypes(),
                fetchAllCaseSubjects(),
            ]);
            clients = clientsData;
            caseTypes = caseTypesData;
            caseSubjects = caseSubjectsData;
        } catch (e) {
            toastState.add(
                TOAST_LEVELS.Error,
                "Erreur",
                "Impossible de charger les données de référence.",
            );
        } finally {
            loadingRefData = false;
            refDataLoaded = true;
        }
    }

    async function loadClientContacts(cId: string) {
        loadingContacts = true;
        try {
            clientContacts = await fetchClientContacts(cId);
        } catch (e) {
            clientContacts = [];
        } finally {
            loadingContacts = false;
        }
    }

    async function handleSubmit() {
        if (!title.trim() || !clientId) return;
        loading = true;
        try {
            const newCase = await createCase({
                title: title.trim(),
                description: description.trim(),
                clientId,
                status,
                caseTypeId: caseTypeId || undefined,
                assignedContactID: assignedContactId || undefined,
                caseSubjectId: caseSubjectId || undefined,
                externalReference: externalReference.trim() || undefined,
                placename: placename.trim() || undefined,
                address1: address1.trim() || undefined,
                address2: address2.trim() || undefined,
                city: city.trim() || undefined,
                postalCode: postalCode.trim() || undefined,
                country: country.trim() || undefined,
                locationType: locationType || undefined,
                latitude: latitude.trim() || undefined,
                longitude: longitude.trim() || undefined,
                locationNotes: locationNotes.trim() || undefined,
            });
            open = false;
            resetForm();
            toastState.add(
                TOAST_LEVELS.Info,
                "Dossier créé",
                `"${newCase.title}" a été créé avec succès.`,
            );
            goto(`/cases/${newCase.id}`);
        } catch (e) {
            toastState.add(
                TOAST_LEVELS.Error,
                "Erreur",
                e instanceof Error ? e.message : "Impossible de créer le dossier",
            );
        } finally {
            loading = false;
        }
    }

    function resetForm() {
        title = "";
        description = "";
        clientId = "";
        caseTypeId = "";
        status = "in_progress";
        isHistorical = false;
        additionalExpanded = false;
        assignedContactId = "";
        caseSubjectId = "";
        externalReference = "";
        locationExpanded = false;
        placename = "";
        address1 = "";
        address2 = "";
        city = "";
        postalCode = "";
        country = "";
        locationType = "";
        latitude = "";
        longitude = "";
        locationNotes = "";
    }
</script>

<Dialog.Root bind:open>
    {#if children}
        <Dialog.Trigger>
            {@render children()}
        </Dialog.Trigger>
    {/if}
    <Dialog.Portal>
        <Dialog.Overlay
            class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
        />
        <Dialog.Content
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border flex flex-col max-h-[90vh] sm:max-w-[720px] md:w-full"
        >
            <div class="flex-shrink-0 px-8 pt-8">
                <Dialog.Title
                    class="flex w-full items-center text-lg font-semibold tracking-tight"
                >
                    Créer un nouveau dossier
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt text-sm"
                    >Remplissez les champs ci-dessous pour créer le dossier. Les informations pourront être modifiées ultérieurement.</Dialog.Description
                >
            </div>
            <Separator.Root class="bg-muted mx-5 !mb-0 !mt-5 block h-px flex-shrink-0" />

            <div class="flex-1 min-h-0 overflow-y-auto px-8 pb-8 pt-6">
            {#if loadingRefData || !refDataLoaded}
                <div class="flex items-center justify-center py-8">
                    <Loader2 class="size-5 animate-spin text-muted-foreground" />
                </div>
            {:else}
                <form class="flex flex-col items-start gap-4" onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
                    <div class="flex flex-col gap-2 w-full">
                        <Label.Root for="title" class="text-sm font-medium">Titre du dossier</Label.Root>
                        <input
                            id="title"
                            name="title"
                            placeholder="Entrez le titre du dossier"
                            required
                            bind:value={title}
                            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
                        />
                    </div>
                    <div class="flex flex-col gap-2 w-full">
                        <Label.Root for="description" class="text-sm font-medium">Description</Label.Root>
                        <textarea
                            id="description"
                            name="description"
                            placeholder="Description du dossier (optionnel)"
                            rows={3}
                            bind:value={description}
                            class="rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm resize-none"
                        ></textarea>
                    </div>

                    <div class="flex justify-between gap-4 w-full">
                        <!-- Client selector -->
                        <div class="flex flex-col gap-2 w-full">
                            <Label.Root for="clientId" class="text-sm font-medium">Client</Label.Root>
                            <select
                                id="clientId"
                                bind:value={clientId}
                                required
                                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm cursor-pointer"
                            >
                                <option value="">-- Sélectionner un client --</option>
                                {#each clients as c}
                                    <option value={c.id}>{c.name}</option>
                                {/each}
                            </select>
                        </div>

                        <!-- CaseType selector -->
                        <div class="flex flex-col gap-2 w-full">
                            <Label.Root for="caseTypeId" class="text-sm font-medium">Type d'affaire</Label.Root>
                            <select
                                id="caseTypeId"
                                bind:value={caseTypeId}
                                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm cursor-pointer"
                            >
                                <option value="">-- Aucun type --</option>
                                {#each caseTypes as ct}
                                    <option value={ct.id}>{ct.name}</option>
                                {/each}
                            </select>
                        </div>
                    </div>

                    <!-- Historical case toggle -->
                    <div class="flex items-center gap-2 w-full mt-1">
                        <input
                            id="isHistorical"
                            type="checkbox"
                            bind:checked={isHistorical}
                            onchange={() => {
                                if (isHistorical) {
                                    status = "released";
                                } else {
                                    status = "in_progress";
                                }
                            }}
                            class="size-4 rounded border-border-input bg-background accent-dark cursor-pointer"
                        />
                        <Label.Root for="isHistorical" class="text-sm font-medium cursor-pointer">C'est un dossier historique</Label.Root>
                    </div>

                    {#if isHistorical}
                        <div class="flex flex-col gap-2 w-full">
                            <Label.Root for="status" class="text-sm font-medium">Statut</Label.Root>
                            <select
                                id="status"
                                bind:value={status}
                                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm cursor-pointer"
                            >
                                <option value="in_progress">En cours</option>
                                <option value="ready">Prêt</option>
                                <option value="released">Publié</option>
                            </select>
                        </div>
                    {/if}

                    <!-- Additional details (collapsible) -->
                    <Collapsible.Root bind:open={additionalExpanded} class="w-full">
                        <Collapsible.Trigger
                            class="flex w-full items-center gap-2 rounded-card-sm px-3 py-2 text-sm font-medium text-foreground-alt hover:bg-muted/50 transition-colors cursor-pointer"
                        >
                            <ChevronRight
                                size={16}
                                class="transition-transform duration-200 {additionalExpanded ? 'rotate-90' : ''}"
                            />
                            <span>Détails supplémentaires</span>
                            {#if !additionalExpanded && additionalFilledCount > 0}
                                <span class="text-xs text-foreground-alt/70">· {additionalFilledCount} rempli{additionalFilledCount > 1 ? 's' : ''}</span>
                            {/if}
                        </Collapsible.Trigger>
                        <Collapsible.Content>
                            <div class="flex flex-col gap-4 pt-4 pb-1">
                                <!-- Assigned contact -->
                                <div class="flex flex-col gap-2 w-full">
                                    <Label.Root for="assignedContactId" class="text-sm font-medium">Contact assigné</Label.Root>
                                    {#if !clientId}
                                        <select
                                            id="assignedContactId"
                                            disabled
                                            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 w-full px-4 text-base sm:text-sm cursor-not-allowed opacity-50"
                                        >
                                            <option value="">Sélectionnez d'abord un client</option>
                                        </select>
                                    {:else if loadingContacts}
                                        <div class="h-input flex items-center px-4 text-sm text-foreground-alt">
                                            <Loader2 size={14} class="animate-spin mr-2" /> Chargement…
                                        </div>
                                    {:else}
                                        <select
                                            id="assignedContactId"
                                            bind:value={assignedContactId}
                                            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm cursor-pointer"
                                        >
                                            <option value="">-- Aucun contact --</option>
                                            {#each clientContacts as contact}
                                                <option value={contact.id}>{contact.firstname} {contact.lastname}</option>
                                            {/each}
                                        </select>
                                    {/if}
                                </div>

                                <!-- Case subject -->
                                <div class="flex flex-col gap-2 w-full">
                                    <Label.Root for="caseSubjectId" class="text-sm font-medium">Sujet du dossier</Label.Root>
                                    <select
                                        id="caseSubjectId"
                                        bind:value={caseSubjectId}
                                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm cursor-pointer"
                                    >
                                        <option value="">-- Aucun sujet --</option>
                                        {#each caseSubjects as cs}
                                            <option value={cs.id}>{cs.firstname} {cs.lastname}</option>
                                        {/each}
                                    </select>
                                </div>

                                <!-- External reference -->
                                <div class="flex flex-col gap-2 w-full">
                                    <Label.Root for="externalReference" class="text-sm font-medium">Référence externe</Label.Root>
                                    <input
                                        id="externalReference"
                                        name="externalReference"
                                        placeholder="Référence externe (optionnel)"
                                        bind:value={externalReference}
                                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
                                    />
                                </div>
                            </div>
                        </Collapsible.Content>
                    </Collapsible.Root>

                    <!-- Location (collapsible) -->
                    <Collapsible.Root bind:open={locationExpanded} class="w-full">
                        <Collapsible.Trigger
                            class="flex w-full items-center gap-2 rounded-card-sm px-3 py-2 text-sm font-medium text-foreground-alt hover:bg-muted/50 transition-colors cursor-pointer"
                        >
                            <ChevronRight
                                size={16}
                                class="transition-transform duration-200 {locationExpanded ? 'rotate-90' : ''}"
                            />
                            <span>Localisation</span>
                            {#if !locationExpanded && locationFilledCount > 0}
                                <span class="text-xs text-foreground-alt/70">· {locationFilledCount} rempli{locationFilledCount > 1 ? 's' : ''}</span>
                            {/if}
                        </Collapsible.Trigger>
                        <Collapsible.Content>
                            <div class="flex flex-col gap-4 pt-4 pb-1">
                                <!-- Placename -->
                                <div class="flex flex-col gap-2 w-full">
                                    <Label.Root for="placename" class="text-sm font-medium">Nom du lieu</Label.Root>
                                    <input
                                        id="placename"
                                        name="placename"
                                        placeholder="Nom du lieu (optionnel)"
                                        bind:value={placename}
                                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
                                    />
                                </div>

                                <!-- Address lines -->
                                <div class="flex justify-between gap-4 w-full">
                                    <div class="flex flex-col gap-2 w-full">
                                        <Label.Root for="address1" class="text-sm font-medium">Adresse 1</Label.Root>
                                        <input
                                            id="address1"
                                            name="address1"
                                            placeholder="Adresse ligne 1"
                                            bind:value={address1}
                                            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
                                        />
                                    </div>
                                    <div class="flex flex-col gap-2 w-full">
                                        <Label.Root for="address2" class="text-sm font-medium">Adresse 2</Label.Root>
                                        <input
                                            id="address2"
                                            name="address2"
                                            placeholder="Adresse ligne 2"
                                            bind:value={address2}
                                            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
                                        />
                                    </div>
                                </div>

                                <!-- City, Postal code, Country -->
                                <div class="flex justify-between gap-4 w-full">
                                    <div class="flex flex-col gap-2 w-full">
                                        <Label.Root for="city" class="text-sm font-medium">Ville</Label.Root>
                                        <input
                                            id="city"
                                            name="city"
                                            placeholder="Ville"
                                            bind:value={city}
                                            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
                                        />
                                    </div>
                                    <div class="flex flex-col gap-2 w-full">
                                        <Label.Root for="postalCode" class="text-sm font-medium">Code postal</Label.Root>
                                        <input
                                            id="postalCode"
                                            name="postalCode"
                                            placeholder="Code postal"
                                            bind:value={postalCode}
                                            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
                                        />
                                    </div>
                                    <div class="flex flex-col gap-2 w-full">
                                        <Label.Root for="country" class="text-sm font-medium">Pays</Label.Root>
                                        <input
                                            id="country"
                                            name="country"
                                            placeholder="Pays"
                                            bind:value={country}
                                            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
                                        />
                                    </div>
                                </div>

                                <!-- Location type -->
                                <div class="flex flex-col gap-2 w-full">
                                    <Label.Root for="locationType" class="text-sm font-medium">Type de lieu</Label.Root>
                                    <select
                                        id="locationType"
                                        bind:value={locationType}
                                        class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm cursor-pointer"
                                    >
                                        <option value="">-- Aucun type --</option>
                                        <option value="home">Domicile</option>
                                        <option value="business">Entreprise</option>
                                        <option value="public">Public</option>
                                        <option value="vehicle">Véhicule</option>
                                        <option value="other">Autre</option>
                                    </select>
                                </div>

                                <!-- Latitude, Longitude -->
                                <div class="flex justify-between gap-4 w-full">
                                    <div class="flex flex-col gap-2 w-full">
                                        <Label.Root for="latitude" class="text-sm font-medium">Latitude</Label.Root>
                                        <input
                                            id="latitude"
                                            name="latitude"
                                            placeholder="ex: 48.8566"
                                            bind:value={latitude}
                                            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
                                        />
                                    </div>
                                    <div class="flex flex-col gap-2 w-full">
                                        <Label.Root for="longitude" class="text-sm font-medium">Longitude</Label.Root>
                                        <input
                                            id="longitude"
                                            name="longitude"
                                            placeholder="ex: 2.3522"
                                            bind:value={longitude}
                                            class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm"
                                        />
                                    </div>
                                </div>

                                <!-- Location notes -->
                                <div class="flex flex-col gap-2 w-full">
                                    <Label.Root for="locationNotes" class="text-sm font-medium">Notes de localisation</Label.Root>
                                    <textarea
                                        id="locationNotes"
                                        name="locationNotes"
                                        placeholder="Notes sur la localisation (optionnel)"
                                        rows={2}
                                        bind:value={locationNotes}
                                        class="rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 py-3 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm resize-none"
                                    ></textarea>
                                </div>
                            </div>
                        </Collapsible.Content>
                    </Collapsible.Root>

                    <div class="flex w-full justify-end mt-4">
                        <button
                            type="submit"
                            disabled={loading || !title.trim() || !clientId}
                            class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center gap-2 px-[50px] text-[15px] font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] disabled:opacity-50 cursor-pointer"
                        >
                            {#if loading}
                                <Loader2 size={16} class="animate-spin" />
                            {/if}
                            Sauvegarder
                        </button>
                    </div>
                </form>
            {/if}
            </div>

            <Dialog.Close
                class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-8 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            >
                <div>
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </div>
            </Dialog.Close>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>


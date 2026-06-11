<script lang="ts">
    import { Dialog, Label, Separator, Collapsible } from "bits-ui";
    import { X, Loader2, ChevronDown, MapPin, SlidersHorizontal } from "@lucide/svelte";
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

    let additionalFilledCount = $derived(
        (assignedContactId ? 1 : 0) +
        (caseSubjectId ? 1 : 0) +
        (externalReference.trim() ? 1 : 0)
    );

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

    const inputClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors";
    const selectClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors cursor-pointer";
    const textareaClass = "w-full rounded-card-sm border border-border-input bg-background px-4 py-3 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors resize-none";
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
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border flex flex-col max-h-[90vh] sm:max-w-[680px] md:w-full"
        >
            <!-- Header -->
            <div class="flex-shrink-0 px-8 pt-8 pb-6">
                <Dialog.Title class="text-base font-semibold tracking-tight">
                    Créer un nouveau dossier
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt text-sm mt-1">
                    Les champs marqués d'un <span class="text-foreground font-medium">*</span> sont obligatoires.
                </Dialog.Description>
            </div>

            <Separator.Root class="bg-border-input mx-0 !m-0 block h-px flex-shrink-0" />

            <div class="flex-1 min-h-0 overflow-y-auto px-8 py-6">
            {#if loadingRefData || !refDataLoaded}
                <div class="flex items-center justify-center py-8">
                    <Loader2 class="size-5 animate-spin text-muted-foreground" />
                </div>
            {:else}
                <form class="flex flex-col gap-4" onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>

                    <!-- Title -->
                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="title" class="text-sm font-medium">
                            Titre du dossier <span class="text-foreground-alt font-normal">*</span>
                        </Label.Root>
                        <input
                            id="title"
                            name="title"
                            placeholder="Entrez le titre du dossier"
                            required
                            bind:value={title}
                            class={inputClass}
                        />
                    </div>

                    <!-- Description -->
                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="description" class="text-sm font-medium text-foreground-alt">Description</Label.Root>
                        <textarea
                            id="description"
                            name="description"
                            placeholder="Description du dossier (optionnel)"
                            rows={3}
                            bind:value={description}
                            class={textareaClass}
                        ></textarea>
                    </div>

                    <!-- Client + Case type -->
                    <div class="grid grid-cols-2 gap-4">
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="clientId" class="text-sm font-medium">
                                Client <span class="text-foreground-alt font-normal">*</span>
                            </Label.Root>
                            <select id="clientId" bind:value={clientId} required class={selectClass}>
                                <option value="">Sélectionner un client</option>
                                {#each clients as c}
                                    <option value={c.id}>{c.name}</option>
                                {/each}
                            </select>
                        </div>

                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="caseTypeId" class="text-sm font-medium text-foreground-alt">Type d'affaire</Label.Root>
                            <select id="caseTypeId" bind:value={caseTypeId} class={selectClass}>
                                <option value="">Aucun type</option>
                                {#each caseTypes as ct}
                                    <option value={ct.id}>{ct.name}</option>
                                {/each}
                            </select>
                        </div>
                    </div>

                    <!-- Historical toggle -->
                    <label class="flex items-center gap-3 cursor-pointer w-fit group">
                        <div class="relative flex items-center">
                            <input
                                id="isHistorical"
                                type="checkbox"
                                bind:checked={isHistorical}
                                onchange={() => { status = isHistorical ? "released" : "in_progress"; }}
                                class="peer size-4 cursor-pointer appearance-none rounded border border-border-input bg-background checked:bg-dark checked:border-dark transition-colors"
                            />
                            <svg
                                class="pointer-events-none absolute inset-0 m-auto size-2.5 text-background opacity-0 peer-checked:opacity-100 transition-opacity"
                                viewBox="0 0 10 8" fill="none"
                            >
                                <path d="M1 4L3.5 6.5L9 1" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                        </div>
                        <span class="text-sm font-medium select-none">C'est un dossier historique</span>
                    </label>

                    {#if isHistorical}
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="status" class="text-sm font-medium">Statut</Label.Root>
                            <select id="status" bind:value={status} class={selectClass}>
                                <option value="in_progress">En cours</option>
                                <option value="ready">Prêt</option>
                                <option value="released">Publié</option>
                            </select>
                        </div>
                    {/if}

                    <!-- Optional sections separator -->
                    <div class="flex items-center gap-3 pt-1">
                        <span class="text-xs font-medium text-foreground-alt/60 uppercase tracking-wider whitespace-nowrap">Informations optionnelles</span>
                        <div class="h-px flex-1 bg-border-input"></div>
                    </div>

                    <!-- Additional details (collapsible) -->
                    <Collapsible.Root bind:open={additionalExpanded} class="w-full rounded-card-sm border border-border-input overflow-hidden">
                        <Collapsible.Trigger
                            class="flex w-full items-center gap-3 px-4 py-3 bg-muted/60 hover:bg-muted transition-colors cursor-pointer text-left"
                        >
                            <SlidersHorizontal size={14} class="text-foreground-alt shrink-0" />
                            <span class="text-sm font-medium flex-1">Détails supplémentaires</span>
                            {#if !additionalExpanded && additionalFilledCount > 0}
                                <span class="text-xs font-medium bg-dark/10 text-foreground-alt px-2 py-0.5 rounded-full">
                                    {additionalFilledCount} rempli{additionalFilledCount > 1 ? 's' : ''}
                                </span>
                            {/if}
                            <ChevronDown
                                size={14}
                                class="text-foreground-alt transition-transform duration-200 shrink-0 {additionalExpanded ? 'rotate-180' : ''}"
                            />
                        </Collapsible.Trigger>
                        <Collapsible.Content>
                            <div class="flex flex-col gap-4 p-4 border-t border-border-input">
                                <!-- Assigned contact -->
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="assignedContactId" class="text-sm font-medium text-foreground-alt">Contact assigné</Label.Root>
                                    {#if !clientId}
                                        <select disabled class="{selectClass} opacity-40 cursor-not-allowed">
                                            <option value="">Sélectionnez d'abord un client</option>
                                        </select>
                                    {:else if loadingContacts}
                                        <div class="h-input flex items-center px-4 text-sm text-foreground-alt border border-border-input rounded-card-sm bg-background">
                                            <Loader2 size={13} class="animate-spin mr-2 text-foreground-alt" /> Chargement…
                                        </div>
                                    {:else}
                                        <select id="assignedContactId" bind:value={assignedContactId} class={selectClass}>
                                            <option value="">Aucun contact</option>
                                            {#each clientContacts as contact}
                                                <option value={contact.id}>{contact.firstname} {contact.lastname}</option>
                                            {/each}
                                        </select>
                                    {/if}
                                </div>

                                <!-- Case subject -->
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="caseSubjectId" class="text-sm font-medium text-foreground-alt">Sujet du dossier</Label.Root>
                                    <select id="caseSubjectId" bind:value={caseSubjectId} class={selectClass}>
                                        <option value="">Aucun sujet</option>
                                        {#each caseSubjects as cs}
                                            <option value={cs.id}>{cs.firstname} {cs.lastname}</option>
                                        {/each}
                                    </select>
                                </div>

                                <!-- External reference -->
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="externalReference" class="text-sm font-medium text-foreground-alt">Référence externe</Label.Root>
                                    <input
                                        id="externalReference"
                                        name="externalReference"
                                        placeholder="Référence externe (optionnel)"
                                        bind:value={externalReference}
                                        class={inputClass}
                                    />
                                </div>
                            </div>
                        </Collapsible.Content>
                    </Collapsible.Root>

                    <!-- Location (collapsible) -->
                    <Collapsible.Root bind:open={locationExpanded} class="w-full rounded-card-sm border border-border-input overflow-hidden">
                        <Collapsible.Trigger
                            class="flex w-full items-center gap-3 px-4 py-3 bg-muted/60 hover:bg-muted transition-colors cursor-pointer text-left"
                        >
                            <MapPin size={14} class="text-foreground-alt shrink-0" />
                            <span class="text-sm font-medium flex-1">Localisation</span>
                            {#if !locationExpanded && locationFilledCount > 0}
                                <span class="text-xs font-medium bg-dark/10 text-foreground-alt px-2 py-0.5 rounded-full">
                                    {locationFilledCount} rempli{locationFilledCount > 1 ? 's' : ''}
                                </span>
                            {/if}
                            <ChevronDown
                                size={14}
                                class="text-foreground-alt transition-transform duration-200 shrink-0 {locationExpanded ? 'rotate-180' : ''}"
                            />
                        </Collapsible.Trigger>
                        <Collapsible.Content>
                            <div class="flex flex-col gap-4 p-4 border-t border-border-input">
                                <!-- Placename -->
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="placename" class="text-sm font-medium text-foreground-alt">Nom du lieu</Label.Root>
                                    <input
                                        id="placename"
                                        name="placename"
                                        placeholder="Nom du lieu (optionnel)"
                                        bind:value={placename}
                                        class={inputClass}
                                    />
                                </div>

                                <!-- Address lines -->
                                <div class="grid grid-cols-2 gap-4">
                                    <div class="flex flex-col gap-1.5">
                                        <Label.Root for="address1" class="text-sm font-medium text-foreground-alt">Adresse 1</Label.Root>
                                        <input id="address1" name="address1" placeholder="Adresse ligne 1" bind:value={address1} class={inputClass} />
                                    </div>
                                    <div class="flex flex-col gap-1.5">
                                        <Label.Root for="address2" class="text-sm font-medium text-foreground-alt">Adresse 2</Label.Root>
                                        <input id="address2" name="address2" placeholder="Adresse ligne 2" bind:value={address2} class={inputClass} />
                                    </div>
                                </div>

                                <!-- City, Postal code, Country -->
                                <div class="grid grid-cols-3 gap-4">
                                    <div class="flex flex-col gap-1.5">
                                        <Label.Root for="city" class="text-sm font-medium text-foreground-alt">Ville</Label.Root>
                                        <input id="city" name="city" placeholder="Ville" bind:value={city} class={inputClass} />
                                    </div>
                                    <div class="flex flex-col gap-1.5">
                                        <Label.Root for="postalCode" class="text-sm font-medium text-foreground-alt">Code postal</Label.Root>
                                        <input id="postalCode" name="postalCode" placeholder="Code postal" bind:value={postalCode} class={inputClass} />
                                    </div>
                                    <div class="flex flex-col gap-1.5">
                                        <Label.Root for="country" class="text-sm font-medium text-foreground-alt">Pays</Label.Root>
                                        <input id="country" name="country" placeholder="Pays" bind:value={country} class={inputClass} />
                                    </div>
                                </div>

                                <!-- Location type -->
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="locationType" class="text-sm font-medium text-foreground-alt">Type de lieu</Label.Root>
                                    <select id="locationType" bind:value={locationType} class={selectClass}>
                                        <option value="">Aucun type</option>
                                        <option value="home">Domicile</option>
                                        <option value="business">Entreprise</option>
                                        <option value="public">Public</option>
                                        <option value="vehicle">Véhicule</option>
                                        <option value="other">Autre</option>
                                    </select>
                                </div>

                                <!-- Latitude, Longitude -->
                                <div class="grid grid-cols-2 gap-4">
                                    <div class="flex flex-col gap-1.5">
                                        <Label.Root for="latitude" class="text-sm font-medium text-foreground-alt">Latitude</Label.Root>
                                        <input id="latitude" name="latitude" placeholder="ex: 48.8566" bind:value={latitude} class={inputClass} />
                                    </div>
                                    <div class="flex flex-col gap-1.5">
                                        <Label.Root for="longitude" class="text-sm font-medium text-foreground-alt">Longitude</Label.Root>
                                        <input id="longitude" name="longitude" placeholder="ex: 2.3522" bind:value={longitude} class={inputClass} />
                                    </div>
                                </div>

                                <!-- Location notes -->
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="locationNotes" class="text-sm font-medium text-foreground-alt">Notes de localisation</Label.Root>
                                    <textarea
                                        id="locationNotes"
                                        name="locationNotes"
                                        placeholder="Notes sur la localisation (optionnel)"
                                        rows={2}
                                        bind:value={locationNotes}
                                        class={textareaClass}
                                    ></textarea>
                                </div>
                            </div>
                        </Collapsible.Content>
                    </Collapsible.Root>

                    <!-- Submit -->
                    <div class="flex w-full justify-end pt-2">
                        <button
                            type="submit"
                            disabled={loading || !title.trim() || !clientId}
                            class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/90 inline-flex items-center justify-center gap-2 px-10 text-sm font-semibold active:scale-[0.98] disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer transition-interactive"
                        >
                            {#if loading}
                                <Loader2 size={14} class="animate-spin" />
                            {/if}
                            Créer le dossier
                        </button>
                    </div>
                </form>
            {/if}
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

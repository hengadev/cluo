<script lang="ts">
    import { Dialog, Label, Separator } from "bits-ui";
    import { X, Loader2 } from "@lucide/svelte";
    import { createCase, fetchAllClients, fetchAllCaseTypes } from "$lib/services/api";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import { goto } from "$app/navigation";
    import type { CaseStatus, CaseType, Client } from "$lib/types/entities";

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

    // Reference data
    let clients: Client[] = $state([]);
    let caseTypes: CaseType[] = $state([]);
    let loadingRefData = $state(false);
    let refDataLoaded = $state(false);

    $effect(() => {
        if (open && !refDataLoaded && !loadingRefData) {
            loadReferenceData();
        }
    });

    async function loadReferenceData() {
        loadingRefData = true;
        try {
            const [clientsData, caseTypesData] = await Promise.all([
                fetchAllClients(),
                fetchAllCaseTypes(),
            ]);
            clients = clientsData;
            caseTypes = caseTypesData;
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
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[720px] md:w-full"
        >
            <Dialog.Title
                class="flex w-full items-center text-lg font-semibold tracking-tight"
            >
                Créer un nouveau dossier
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt text-sm"
                >Remplissez les champs ci-dessous pour créer le dossier. Les informations pourront être modifiées ultérieurement.</Dialog.Description
            >
            <Separator.Root class="bg-muted mx-5 !mb-6 !mt-5 block h-px" />

            {#if loadingRefData || !refDataLoaded}
                <div class="flex items-center justify-center py-8">
                    <Loader2 class="size-5 animate-spin text-muted-foreground" />
                </div>
            {:else}
                <form class="flex flex-col items-start gap-4 pb-4" onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
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

                    <div class="flex justify-between gap-4 w-full">
                        <!-- Status selector -->
                        <div class="flex flex-col gap-2 w-full">
                            <Label.Root for="status" class="text-sm font-medium">Statut</Label.Root>
                            <select
                                id="status"
                                bind:value={status}
                                class="h-input rounded-card-sm border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-dark-40 focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-base focus:ring-2 focus:ring-offset-2 sm:text-sm cursor-pointer"
                            >
                                <option value="in_progress">En cours</option>
                                <option value="ready">Prêt</option>
                            </select>
                        </div>
                        <div class="w-full"></div>
                    </div>

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


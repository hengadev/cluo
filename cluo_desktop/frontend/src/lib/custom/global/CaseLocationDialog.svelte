<script lang="ts">
    import { Dialog, Label, Separator } from "bits-ui";
    import { X, Loader2 } from "@lucide/svelte";
    import { updateCase } from "$lib/services/api";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import type { Case, LocationType } from "$lib/types/entities";

    const toastState = getToastContext();

    type Props = {
        open?: boolean;
        caseId: string;
        caseData: Case;
        onSaved?: (updated: Case) => void;
    };
    let { open = $bindable(false), caseId, caseData, onSaved }: Props = $props();

    let placename = $state("");
    let address1 = $state("");
    let address2 = $state("");
    let city = $state("");
    let postalCode = $state("");
    let country = $state("");
    let locationType = $state<LocationType | "">("");
    let locationNotes = $state("");
    let saving = $state(false);

    $effect(() => {
        if (open) {
            placename = caseData.placename ?? "";
            address1 = caseData.address1 ?? "";
            address2 = caseData.address2 ?? "";
            city = caseData.city ?? "";
            postalCode = caseData.postalCode ?? "";
            country = caseData.country ?? "";
            locationType = caseData.locationType ?? "";
            locationNotes = caseData.locationNotes ?? "";
        }
    });

    const LOCATION_TYPE_OPTIONS: { value: LocationType; label: string }[] = [
        { value: "home", label: "Domicile" },
        { value: "business", label: "Entreprise" },
        { value: "public", label: "Lieu public" },
        { value: "vehicle", label: "Véhicule" },
        { value: "other", label: "Autre" },
    ];

    const inputClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors";
    const selectClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors cursor-pointer";

    async function handleSubmit() {
        saving = true;
        try {
            const updated = await updateCase(caseId, {
                placename: placename.trim() || undefined,
                address1: address1.trim() || undefined,
                address2: address2.trim() || undefined,
                city: city.trim() || undefined,
                postalCode: postalCode.trim() || undefined,
                country: country.trim() || undefined,
                locationType: (locationType as LocationType) || undefined,
                locationNotes: locationNotes.trim() || undefined,
            });
            toastState.add(TOAST_LEVELS.Info, "Lieu mis à jour", "Les informations de lieu ont été enregistrées.");
            open = false;
            onSaved?.(updated);
        } catch (e) {
            toastState.add(TOAST_LEVELS.Error, "Erreur", e instanceof Error ? e.message : "Impossible de mettre à jour le lieu.");
        } finally {
            saving = false;
        }
    }
</script>

<Dialog.Root bind:open>
    <Dialog.Portal>
        <Dialog.Overlay
            class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
        />
        <Dialog.Content
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border flex flex-col sm:max-w-[520px] md:w-full"
        >
            <div class="flex-shrink-0 px-8 pt-8 pb-6">
                <Dialog.Title class="text-base font-semibold tracking-tight">
                    Lieu de l'affaire
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt text-sm mt-1">
                    Adresse et informations de localisation du dossier.
                </Dialog.Description>
            </div>

            <Separator.Root class="bg-border-input mx-0 !m-0 block h-px flex-shrink-0" />

            <div class="px-8 py-6 overflow-y-auto">
                <form
                    class="flex flex-col gap-4"
                    onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}
                >
                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="loc-placename" class="text-sm font-medium text-foreground-alt">
                            Nom du lieu
                        </Label.Root>
                        <input
                            id="loc-placename"
                            type="text"
                            bind:value={placename}
                            placeholder="Ex : Résidence Les Pins, Centre commercial…"
                            class={inputClass}
                        />
                    </div>

                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="loc-type" class="text-sm font-medium text-foreground-alt">
                            Type de lieu
                        </Label.Root>
                        <select id="loc-type" bind:value={locationType} class={selectClass}>
                            <option value="">— Non défini —</option>
                            {#each LOCATION_TYPE_OPTIONS as opt}
                                <option value={opt.value}>{opt.label}</option>
                            {/each}
                        </select>
                    </div>

                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="loc-address1" class="text-sm font-medium text-foreground-alt">
                            Adresse
                        </Label.Root>
                        <input
                            id="loc-address1"
                            type="text"
                            bind:value={address1}
                            placeholder="Adresse ligne 1"
                            class={inputClass}
                        />
                    </div>

                    <input
                        type="text"
                        bind:value={address2}
                        placeholder="Adresse ligne 2 (optionnel)"
                        class={inputClass}
                    />

                    <div class="grid grid-cols-2 gap-3">
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="loc-postal" class="text-sm font-medium text-foreground-alt">
                                Code postal
                            </Label.Root>
                            <input
                                id="loc-postal"
                                type="text"
                                bind:value={postalCode}
                                placeholder="75001"
                                class={inputClass}
                            />
                        </div>
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="loc-city" class="text-sm font-medium text-foreground-alt">
                                Ville
                            </Label.Root>
                            <input
                                id="loc-city"
                                type="text"
                                bind:value={city}
                                placeholder="Paris"
                                class={inputClass}
                            />
                        </div>
                    </div>

                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="loc-country" class="text-sm font-medium text-foreground-alt">
                            Pays
                        </Label.Root>
                        <input
                            id="loc-country"
                            type="text"
                            bind:value={country}
                            placeholder="France"
                            class={inputClass}
                        />
                    </div>

                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="loc-notes" class="text-sm font-medium text-foreground-alt">
                            Notes sur le lieu
                        </Label.Root>
                        <textarea
                            id="loc-notes"
                            bind:value={locationNotes}
                            placeholder="Digicode, étage, instructions d'accès…"
                            rows="2"
                            class="w-full rounded-card-sm border border-border-input bg-background px-4 py-2 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors resize-none"
                        ></textarea>
                    </div>

                    <div class="flex justify-end gap-2 pt-2">
                        <Dialog.Close
                            class="h-input rounded-input bg-transparent text-foreground border border-border-input hover:bg-muted inline-flex items-center justify-center px-5 text-sm font-semibold active:scale-[0.98] cursor-pointer transition-colors"
                        >
                            Annuler
                        </Dialog.Close>
                        <button
                            type="submit"
                            disabled={saving}
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

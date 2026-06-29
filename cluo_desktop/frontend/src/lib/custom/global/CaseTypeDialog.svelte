<script lang="ts">
    import { Dialog, Label, Separator } from "bits-ui";
    import { X, Loader2 } from "@lucide/svelte";
    import { updateCase } from "$lib/services/api";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import type { Case, CaseType } from "$lib/types/entities";

    const toastState = getToastContext();

    type Props = {
        open?: boolean;
        caseId: string;
        currentCaseTypeId: string | null;
        allCaseTypes: CaseType[];
        onSaved?: (updated: Case) => void;
    };
    let { open = $bindable(false), caseId, currentCaseTypeId, allCaseTypes, onSaved }: Props = $props();

    let selectedId = $state("");
    let saving = $state(false);

    $effect(() => {
        if (open) {
            selectedId = currentCaseTypeId ?? "";
        }
    });

    const selectClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors cursor-pointer";

    async function handleSubmit() {
        saving = true;
        try {
            const updated = await updateCase(caseId, {
                caseTypeId: selectedId || null,
            });
            toastState.add(TOAST_LEVELS.Info, "Type mis à jour", "Le type d'affaire a été enregistré.");
            open = false;
            onSaved?.(updated);
        } catch (e) {
            toastState.add(TOAST_LEVELS.Error, "Erreur", e instanceof Error ? e.message : "Impossible de mettre à jour.");
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
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border flex flex-col sm:max-w-[400px] md:w-full"
        >
            <div class="flex-shrink-0 px-8 pt-8 pb-6">
                <Dialog.Title class="text-base font-semibold tracking-tight">
                    Type d'affaire
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt text-sm mt-1">
                    Sélectionnez le type qui correspond le mieux à ce dossier.
                </Dialog.Description>
            </div>

            <Separator.Root class="bg-border-input mx-0 !m-0 block h-px flex-shrink-0" />

            <div class="px-8 py-6">
                <form
                    class="flex flex-col gap-4"
                    onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}
                >
                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="case-type-select" class="text-sm font-medium">
                            Type d'affaire
                        </Label.Root>
                        <select id="case-type-select" bind:value={selectedId} class={selectClass}>
                            <option value="">— Aucun type —</option>
                            {#each allCaseTypes as ct}
                                <option value={ct.id}>{ct.name}</option>
                            {/each}
                        </select>
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

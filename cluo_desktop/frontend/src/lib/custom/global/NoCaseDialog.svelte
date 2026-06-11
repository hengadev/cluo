<script lang="ts">
    import { Dialog, Separator } from "bits-ui";
    import { X, FolderPlus, FolderSearch, ArrowRight } from "@lucide/svelte";
    import { goto } from "$app/navigation";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import { fetchAllCases } from "$lib/services/api";

    type Props = {
        open?: boolean;
        onCreateCase: () => void;
    };
    let { open = $bindable(false), onCreateCase }: Props = $props();

    const toastState = getToastContext();

    let hasCases: boolean | null = $state(null);

    $effect(() => {
        if (open) {
            hasCases = null;
            fetchAllCases({ pageSize: 1 })
                .then(res => { hasCases = res.cases.length > 0; })
                .catch(() => { hasCases = false; });
        }
    });

    function handleCreateCase() {
        open = false;
        onCreateCase();
    }

    function handleExistingCases() {
        open = false;
        toastState.add(TOAST_LEVELS.Info, "Dossiers", "Accédez à vos dossiers depuis la page dédiée.");
        goto("/cases");
    }
</script>

<Dialog.Root bind:open>
    <Dialog.Portal>
        <Dialog.Overlay
            class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
        />
        <Dialog.Content
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border sm:max-w-[380px] md:w-full"
        >
            <div class="px-8 pt-8 pb-6">
                <Dialog.Title class="text-base font-semibold tracking-tight">
                    Choisissez un dossier
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt text-sm mt-1">
                    Créez un nouveau dossier ou reprenez un dossier existant.
                </Dialog.Description>
            </div>

            <Separator.Root class="bg-border-input mx-0 !m-0 block h-px flex-shrink-0" />

            <div class="flex flex-col gap-2.5 p-8 pt-6">
                <button
                    onclick={handleCreateCase}
                    class="group flex items-center gap-3.5 rounded-input bg-dark text-background p-4 text-left transition-interactive duration-200 hover:bg-dark/90 active:scale-[0.98] cursor-pointer shadow-mini"
                >
                    <div class="flex items-center justify-center size-8 rounded-[6px] bg-background/[0.08] shrink-0">
                        <FolderPlus size={16} strokeWidth={1.75} />
                    </div>
                    <div class="flex-1 min-w-0">
                        <p class="text-sm font-semibold">Nouveau dossier</p>
                        <p class="text-xs text-background/60 mt-0.5">Créer un dossier depuis zéro</p>
                    </div>
                    <ArrowRight size={14} class="shrink-0 opacity-0 group-hover:opacity-100 group-hover:translate-x-0.5 transition-interactive duration-200" />
                </button>

                {#if hasCases}
                <button
                    onclick={handleExistingCases}
                    class="group flex items-center gap-3.5 rounded-input border border-border-input bg-transparent text-foreground p-4 text-left transition-interactive duration-200 hover:bg-foreground/[0.04] active:scale-[0.98] cursor-pointer"
                >
                    <div class="flex items-center justify-center size-8 rounded-[6px] bg-foreground/[0.06] shrink-0">
                        <FolderSearch size={16} strokeWidth={1.75} />
                    </div>
                    <div class="min-w-0 flex-1">
                        <p class="text-sm font-semibold">Dossiers existants</p>
                        <p class="text-xs text-foreground-alt mt-0.5">Ouvrir un dossier déjà créé</p>
                    </div>
                    <ArrowRight size={14} class="shrink-0 opacity-0 group-hover:opacity-100 group-hover:translate-x-0.5 transition-interactive duration-200" />
                </button>
                {/if}
            </div>

            <Dialog.Close
                class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-6 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            >
                <div>
                    <X class="text-foreground-alt size-4" />
                    <span class="sr-only">Fermer</span>
                </div>
            </Dialog.Close>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

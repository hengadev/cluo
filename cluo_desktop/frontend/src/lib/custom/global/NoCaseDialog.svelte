<script lang="ts">
    import { Dialog } from "bits-ui";
    import { X, FolderPlus, FolderSearch, ArrowRight } from "@lucide/svelte";
    import { goto } from "$app/navigation";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";

    type Props = {
        open?: boolean;
        onCreateCase: () => void;
    };
    let { open = $bindable(false), onCreateCase }: Props = $props();

    const toastState = getToastContext();

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
            class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/60 backdrop-blur-[2px]"
        />
        <Dialog.Content
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[380px] md:w-full"
        >
            <Dialog.Title
                class="text-lg font-semibold tracking-tight text-center"
            >
                Pour continuer, choisissez un dossier
            </Dialog.Title>
            <Dialog.Description class="sr-only">
                Créer un nouveau dossier ou accéder à vos dossiers existants.
            </Dialog.Description>

            <div class="flex flex-col gap-3 mt-10">
                <!-- Nouveau dossier (primary card) -->
                <button
                    onclick={handleCreateCase}
                    class="group flex items-center gap-3 rounded-input bg-dark text-background p-4 text-left transition-all duration-200 hover:bg-dark/90 active:scale-[0.98] cursor-pointer animate-in fade-in slide-in-from-bottom-2 duration-200"
                >
                    <FolderPlus size={20} strokeWidth={1.75} class="shrink-0" />
                    <div class="flex-1 min-w-0">
                        <p class="text-sm font-semibold">Nouveau dossier</p>
                        <p class="text-xs text-background/70 mt-0.5">Créer un nouveau dossier depuis zéro</p>
                    </div>
                    <ArrowRight size={16} class="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity duration-150" />
                </button>

                <!-- Dossiers existants (secondary card) -->
                <button
                    onclick={handleExistingCases}
                    class="group flex items-center gap-3 rounded-input border-2 border-border-input bg-transparent text-foreground p-4 text-left transition-all duration-200 hover:bg-foreground/5 active:scale-[0.98] cursor-pointer animate-in fade-in slide-in-from-bottom-2 duration-200"
                    style="animation-delay: 80ms;"
                >
                    <FolderSearch size={20} strokeWidth={1.75} class="shrink-0" />
                    <div class="flex-1 min-w-0">
                        <p class="text-sm font-semibold">Dossiers existants</p>
                        <p class="text-xs text-foreground-alt mt-0.5">Ouvrir un dossier déjà créé</p>
                    </div>
                    <ArrowRight size={16} class="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity duration-150" />
                </button>
            </div>

            <Dialog.Close
                class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            >
                <div>
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </div>
            </Dialog.Close>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

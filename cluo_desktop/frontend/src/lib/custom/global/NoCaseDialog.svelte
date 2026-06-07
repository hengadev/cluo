<script lang="ts">
    import { Dialog } from "bits-ui";
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
            class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/65 backdrop-blur-[2px]"
        />
        <Dialog.Content
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] sm:max-w-[400px] overflow-hidden border"
        >
            <!-- Dot-grid texture -->
            <div class="dot-grid pointer-events-none absolute inset-0" aria-hidden="true" />

            <div class="relative p-8 pt-7">
                <!-- Header -->
                <header class="mb-8">
                    <p class="mb-2.5 font-mono text-[10px] uppercase tracking-[0.18em] text-foreground-alt">
                        Espace de travail
                    </p>
                    <div class="editorial-rule mb-5" />
                    <Dialog.Title class="font-serif text-[2rem] leading-[1.15] tracking-tight text-foreground">
                        Choisissez<br />un dossier
                    </Dialog.Title>
                    <Dialog.Description class="sr-only">
                        Créer un nouveau dossier ou accéder à vos dossiers existants.
                    </Dialog.Description>
                </header>

                <!-- Actions -->
                <div class="flex flex-col gap-2.5">
                    <button
                        onclick={handleCreateCase}
                        class="primary-card group relative flex cursor-pointer items-center gap-3.5 rounded-input bg-dark p-4 text-left text-background transition-all duration-200 hover:bg-dark/90 active:scale-[0.99]"
                    >
                        <span class="absolute right-3.5 top-3 font-mono text-[10px] tabular-nums text-background/30">01</span>
                        <div class="flex size-8 shrink-0 items-center justify-center rounded-[6px] bg-background/[0.08]">
                            <FolderPlus size={16} strokeWidth={1.75} />
                        </div>
                        <div class="min-w-0 flex-1">
                            <p class="mb-1 text-[13px] font-semibold leading-none tracking-wide">Nouveau dossier</p>
                            <p class="text-[11px] leading-snug text-background/55">Créer un dossier depuis zéro</p>
                        </div>
                        <ArrowRight size={14} class="shrink-0 translate-x-0 opacity-0 transition-all duration-200 group-hover:translate-x-0.5 group-hover:opacity-100" />
                    </button>

                    {#if hasCases}
                    <button
                        onclick={handleExistingCases}
                        class="card-secondary group relative flex cursor-pointer items-center gap-3.5 rounded-input border border-border-input bg-transparent p-4 text-left text-foreground transition-all duration-200 hover:bg-foreground/[0.04] active:scale-[0.99]"
                    >
                        <span class="absolute right-3.5 top-3 font-mono text-[10px] tabular-nums text-foreground/25">02</span>
                        <div class="flex size-8 shrink-0 items-center justify-center rounded-[6px] bg-foreground/[0.06]">
                            <FolderSearch size={16} strokeWidth={1.75} />
                        </div>
                        <div class="min-w-0 flex-1">
                            <p class="mb-1 text-[13px] font-semibold leading-none tracking-wide">Dossiers existants</p>
                            <p class="text-[11px] leading-snug text-foreground-alt">Ouvrir un dossier déjà créé</p>
                        </div>
                        <ArrowRight size={14} class="shrink-0 translate-x-0 opacity-0 transition-all duration-200 group-hover:translate-x-0.5 group-hover:opacity-100" />
                    </button>
                    {/if}
                </div>

                <!-- Footer -->
                <div class="mt-7 border-t border-border-card pt-4">
                    <p class="font-mono text-[9px] uppercase tracking-[0.15em] text-foreground/20">
                        Cluo — Gestion de dossiers
                    </p>
                </div>
            </div>

            <Dialog.Close
                class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-4 top-4 z-10 flex size-7 cursor-pointer items-center justify-center rounded-[5px] text-foreground/40 transition-all duration-150 hover:bg-foreground/[0.07] hover:text-foreground active:scale-95 focus-visible:ring-2 focus-visible:ring-offset-2"
            >
                <X size={14} />
                <span class="sr-only">Fermer</span>
            </Dialog.Close>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<style>
    .dot-grid {
        background-image: radial-gradient(circle, var(--foreground) 0.5px, transparent 0.5px);
        background-size: 20px 20px;
        opacity: 0.025;
    }

    .editorial-rule {
        position: relative;
        height: 6px;
    }
    .editorial-rule::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 3px;
        background: var(--tertiary);
    }
    .editorial-rule::after {
        content: "";
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        height: 1px;
        background: var(--foreground);
        opacity: 0.2;
    }

    .primary-card {
        border-left: 3px solid var(--tertiary);
        animation: cardEnter 0.45s cubic-bezier(0.16, 1, 0.3, 1) both;
    }

    .card-secondary {
        animation: cardEnter 0.45s cubic-bezier(0.16, 1, 0.3, 1) 75ms both;
    }

    @keyframes cardEnter {
        from {
            opacity: 0;
            transform: translateY(10px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
</style>

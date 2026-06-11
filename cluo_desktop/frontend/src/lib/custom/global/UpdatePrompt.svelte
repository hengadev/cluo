<script lang="ts">
    import { X, Download, Loader2, AlertCircle, Check } from "@lucide/svelte";
    import { onMount, onDestroy } from "svelte";
    import { EventsOn, EventsOff } from "$lib/wailsjs/runtime/runtime";
    type UpdateInfo = {
        available: boolean;
        version: string;
        release_notes: string;
        download_url: string;
    };

    type ProgressEvent = {
        downloaded: number;
        total: number;
        percent: number;
    };

    type PromptPhase = "hidden" | "available" | "downloading" | "installing" | "ready" | "error";

    // State
    let phase: PromptPhase = $state("hidden");
    let newVersion = $state<string>("");
    let progress = $state<ProgressEvent>({ downloaded: 0, total: 0, percent: 0 });
    let errorMessage = $state<string>("");

    // Wails updater bindings (lazy-loaded)
    let Updater: {
        CheckForUpdate: () => Promise<UpdateInfo>;
        DownloadAndInstall: () => Promise<void>;
        RestartApp: () => Promise<void>;
    } | null = $state(null);

    let visible = $derived(phase !== "hidden");

    function dismiss() {
        phase = "hidden";
    }

    onMount(async () => {
        // Import the Updater bindings dynamically
        try {
            const module = await import("$lib/wailsjs/go/updater/Updater");
            Updater = module;
        } catch (e) {
            console.warn("Updater bindings not available:", e);
            return;
        }

        // Listen for updater events
        EventsOn("updater:progress", (data: ProgressEvent) => {
            progress = data;
        });

        EventsOn("updater:status", (status: string) => {
            if (status === "downloading" && phase !== "hidden") {
                phase = "downloading";
            } else if (status === "installing" && phase !== "hidden") {
                phase = "installing";
            } else if (status === "ready" && phase !== "hidden") {
                phase = "ready";
                // Auto-restart once the update is ready
                autoRestart();
            }
        });

        EventsOn("updater:error", (message: string) => {
            errorMessage = message;
            if (phase !== "hidden") {
                phase = "error";
            }
        });

        // Check for update on launch
        try {
            const info = await Updater.CheckForUpdate();
            if (info.available) {
                newVersion = info.version;
                phase = "available";
            }
        } catch {
            // ManifestURL not configured (dev build) or network error — silently skip
        }
    });

    onDestroy(() => {
        EventsOff("updater:progress", "updater:status", "updater:error");
    });

    async function installAndRestart() {
        if (!Updater) return;

        try {
            await Updater.DownloadAndInstall();
            // The "ready" status event triggers autoRestart()
        } catch (e) {
            errorMessage = e instanceof Error ? e.message : String(e);
            phase = "error";
        }
    }

    async function autoRestart() {
        if (!Updater) return;
        try {
            await Updater.RestartApp();
        } catch (e) {
            errorMessage = e instanceof Error ? e.message : String(e);
            phase = "error";
        }
    }

    function formatBytes(bytes: number): string {
        if (bytes === 0) return "0 B";
        const k = 1024;
        const sizes = ["B", "KB", "MB", "GB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
    }
</script>

{#if visible}
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        class="fixed bottom-5 right-5 z-50 w-[380px] animate-in slide-in-from-bottom-4 fade-in duration-300"
        onkeydown={(e) => { if (e.key === "Escape") dismiss(); }}
    >
        <div class="rounded-card-lg bg-background border border-border-card shadow-popover">
            <!-- Header -->
            <div class="flex items-start justify-between px-5 pt-5 pb-3">
                <div class="flex items-center gap-2 text-sm font-semibold tracking-tight">
                    {#if phase === "available"}
                        <Download class="size-4" />
                        Mise à jour disponible
                    {:else if phase === "downloading"}
                        <Loader2 class="size-4 animate-spin" />
                        Téléchargement…
                    {:else if phase === "installing"}
                        <Loader2 class="size-4 animate-spin" />
                        Installation…
                    {:else if phase === "ready"}
                        <Check class="size-4 text-success" />
                        Redémarrage en cours…
                    {:else if phase === "error"}
                        <AlertCircle class="size-4 text-destructive" />
                        Erreur
                    {/if}
                </div>
                {#if phase === "available" || phase === "error"}
                    <button
                        onclick={dismiss}
                        class="text-foreground-alt hover:text-foreground -mt-0.5 -mr-1 rounded-md p-1 transition-colors cursor-pointer"
                    >
                        <X class="size-3.5" />
                    </button>
                {/if}
            </div>

            <!-- Body -->
            <div class="px-5 pb-4 space-y-3">
                {#if phase === "available"}
                    <p class="text-foreground-alt text-sm">
                        La version <strong class="text-foreground">v{newVersion}</strong> est disponible.
                    </p>
                {:else if phase === "downloading"}
                    <div class="space-y-1.5">
                        <div class="flex justify-between text-xs">
                            <span class="text-foreground-alt">
                                {formatBytes(progress.downloaded)}{#if progress.total > 0} / {formatBytes(progress.total)}{/if}
                            </span>
                            <span class="font-medium">{progress.percent.toFixed(1)}%</span>
                        </div>
                        <div class="h-1 w-full overflow-hidden rounded-full bg-muted">
                            <div
                                class="h-full bg-dark transition-interactive duration-200 rounded-full"
                                style="width: {progress.percent}%"
                            ></div>
                        </div>
                    </div>
                {:else if phase === "installing"}
                    <p class="text-foreground-alt text-sm">
                        Installation de la mise à jour…
                    </p>
                {:else if phase === "ready"}
                    <p class="text-foreground-alt text-sm">
                        L'application va redémarrer pour appliquer la mise à jour.
                    </p>
                {:else if phase === "error"}
                    <p class="text-sm text-destructive">
                        {errorMessage || "Une erreur inconnue est survenue."}
                    </p>
                {/if}
            </div>

            <!-- Actions -->
            {#if phase === "available"}
                <div class="flex justify-end gap-2 border-t border-border-card px-5 py-4">
                    <button
                        onclick={dismiss}
                        class="h-input rounded-input bg-transparent text-foreground border border-border-input hover:bg-muted inline-flex items-center justify-center px-4 text-sm font-medium cursor-pointer transition-colors"
                    >
                        Plus tard
                    </button>
                    <button
                        onclick={installAndRestart}
                        class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/90 inline-flex items-center justify-center gap-2 px-4 text-sm font-semibold cursor-pointer transition-interactive"
                    >
                        <Download size={14} />
                        Installer et redémarrer
                    </button>
                </div>
            {/if}
        </div>
    </div>
{/if}

<script lang="ts">
    import { currentCase } from "$lib/stores/case";
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import {
        fetchCasePieces,
        uploadPiece,
        deletePiece,
        getPieceDownloadUrl,
    } from "$lib/services/api";
    import type { Piece } from "$lib/types/entities";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import Spinner from "$lib/components/Spinner.svelte";
    import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
    import { FileText, Upload, Trash2, Download, File } from "@lucide/svelte";

    const toastState = getToastContext();

    let pieces = $state<Piece[]>([]);
    let loading = $state(true);
    let uploading = $state(false);
    let fileInput = $state<HTMLInputElement>();
    let notesInput = $state("");
    let showUploadForm = $state(false);

    $effect(() => {
        const caseId = $page.params.id;
        if (caseId && caseId !== $currentCase.id) {
            currentCase.setCase(caseId);
        }
    });

    onMount(loadPieces);

    async function loadPieces() {
        const caseId = $page.params.id;
        if (!caseId) return;
        loading = true;
        try {
            const response = await fetchCasePieces(caseId);
            pieces = response.pieces;
        } catch (err) {
            toastState.add(TOAST_LEVELS.Error, "Erreur", "Impossible de charger les pièces.");
        } finally {
            loading = false;
        }
    }

    async function handleUpload() {
        const target = fileInput;
        const files = target?.files;
        if (!files || files.length === 0) return;

        const caseId = $page.params.id;
        uploading = true;

        for (const file of Array.from(files)) {
            try {
                await uploadPiece(caseId, file, notesInput || undefined);
            } catch (err) {
                toastState.add(
                    TOAST_LEVELS.Error,
                    "Erreur",
                    `Impossible d'ajouter ${file.name}.`,
                );
            }
        }

        notesInput = "";
        showUploadForm = false;
        if (target) target.value = "";
        uploading = false;
        await loadPieces();
    }

    async function handleDelete(piece: Piece) {
        const caseId = $page.params.id;
        try {
            await deletePiece(caseId, piece.id);
            pieces = pieces.filter((p) => p.id !== piece.id);
        } catch (err) {
            toastState.add(TOAST_LEVELS.Error, "Erreur", "Impossible de supprimer la pièce.");
        }
    }

    function handleDownload(piece: Piece) {
        const caseId = $page.params.id;
        const url = getPieceDownloadUrl(caseId, piece.id);
        // Open in a new tab — the browser will use the authenticated session
        // cookie to download the file via the API.
        window.open(url, "_blank");
    }

    function formatFileSize(bytes: number): string {
        if (bytes < 1024) return `${bytes} o`;
        if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} Ko`;
        return `${(bytes / (1024 * 1024)).toFixed(1)} Mo`;
    }

    function formatDate(dateStr: string): string {
        return new Date(dateStr).toLocaleDateString("fr-FR", {
            day: "2-digit",
            month: "short",
            year: "numeric",
        });
    }

    function mimeTypeIcon(_mime: string): typeof FileText {
        return FileText;
    }
</script>

<div class="page-content flex-1 min-h-0">
    <!-- Hidden File Input -->
    <input
        type="file"
        bind:this={fileInput}
        multiple
        onchange={handleUpload}
        class="hidden"
    />

    <!-- Header -->
    <div class="flex items-center justify-between">
        <div>
            <h1 class="text-2xl font-bold text-foreground">Pièces</h1>
            <p class="text-sm text-muted-foreground">
                {pieces.length} pièce{pieces.length !== 1 ? "s" : ""} jointe{pieces.length !== 1 ? "s" : ""}
            </p>
        </div>
        <button
            type="button"
            onclick={() => {
                showUploadForm = !showUploadForm;
            }}
            class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
        >
            <Upload size={16} class="mr-2" />
            Ajouter des pièces
        </button>
    </div>

    <!-- Upload form -->
    {#if showUploadForm}
        <div class="border border-border-card rounded-card p-4 bg-background-alt animate-fade-in">
            <p class="text-sm font-medium text-foreground mb-3">Ajouter des fichiers</p>
            <div class="flex flex-col gap-3">
                <textarea
                    bind:value={notesInput}
                    placeholder="Notes (optionnel)…"
                    rows="2"
                    class="rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-3 py-2 text-sm focus:ring-2 focus:ring-offset-2 resize-none"
                ></textarea>
                <div class="flex justify-end gap-2">
                    <button
                        type="button"
                        onclick={() => { showUploadForm = false; notesInput = ""; }}
                        class="h-input rounded-input bg-transparent text-foreground hover:bg-muted inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] border border-border-input cursor-pointer"
                    >
                        Annuler
                    </button>
                    <button
                        type="button"
                        onclick={() => fileInput?.click()}
                        disabled={uploading}
                        class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50"
                    >
                        {uploading ? "Envoi en cours…" : "Choisir des fichiers"}
                    </button>
                </div>
            </div>
        </div>
    {/if}

    <!-- Content -->
    {#if loading}
        <div class="flex items-center justify-center py-12">
            <Spinner size="lg" />
        </div>
    {:else if pieces.length === 0}
        <div class="border border-dashed border-border rounded-lg bg-muted/20 flex flex-col items-center justify-center flex-1 gap-4 min-h-[40vh]">
            <FileText class="w-12 h-12 text-muted-foreground" />
            <p class="text-muted-foreground">Aucune pièce jointe pour ce dossier.</p>
            <button
                type="button"
                onclick={() => { showUploadForm = true; }}
                class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
            >
                <Upload size={16} class="mr-2" />
                Ajouter des pièces
            </button>
        </div>
    {:else}
        <div class="flex-1 min-h-0 overflow-y-auto">
            <div class="flex flex-col divide-y divide-border-input">
                {#each pieces as piece, index (piece.id)}
                    <div
                        class="flex items-center gap-4 py-3 px-2 hover:bg-surface/50 transition-colors animate-fade-in group"
                        style="animation-delay: {index * 30}ms;"
                    >
                        <!-- File icon -->
                        <div class="p-2 rounded-input bg-muted">
                            <svelte:component this={mimeTypeIcon(piece.mimeType)} size={20} class="text-muted-foreground" />
                        </div>

                        <!-- File info -->
                        <div class="flex-1 min-w-0">
                            <p class="text-sm font-medium text-foreground truncate">{piece.filename}</p>
                            <div class="flex items-center gap-3 mt-0.5">
                                <span class="text-xs text-muted-foreground">{formatFileSize(piece.sizeBytes)}</span>
                                <span class="text-xs text-muted-foreground">{piece.mimeType}</span>
                                <span class="text-xs text-muted-foreground">{formatDate(piece.createdAt)}</span>
                            </div>
                            {#if piece.notes}
                                <p class="text-xs text-muted-foreground mt-1 truncate">{piece.notes}</p>
                            {/if}
                        </div>

                        <!-- Actions -->
                        <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                            <button
                                type="button"
                                onclick={() => handleDownload(piece)}
                                class="p-2 rounded-input hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                                title="Télécharger"
                            >
                                <Download size={16} />
                            </button>
                            <ConfirmDialog
                                title="Supprimer la pièce"
                                description="Voulez-vous vraiment supprimer « {piece.filename} » ? Cette action est irréversible."
                                onConfirm={() => handleDelete(piece)}
                            >
                                <button
                                    type="button"
                                    class="p-2 rounded-input btn-ghost-destructive cursor-pointer"
                                    title="Supprimer"
                                >
                                    <Trash2 size={16} />
                                </button>
                            </ConfirmDialog>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    {/if}
</div>

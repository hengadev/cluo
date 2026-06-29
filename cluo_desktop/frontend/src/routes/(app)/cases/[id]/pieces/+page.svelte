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
    import { Dialog, Separator } from "bits-ui";
    import { FileText, Upload, Trash2, Download, X, Loader2 } from "@lucide/svelte";

    const toastState = getToastContext();

    let pieces = $state<Piece[]>([]);
    let loading = $state(true);
    let uploading = $state(false);
    let fileInput = $state<HTMLInputElement>();
    let notesInput = $state("");
    let showUploadModal = $state(false);
    let isDragging = $state(false);

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

    async function uploadFiles(files: FileList) {
        const caseId = $page.params.id;
        uploading = true;
        let hasError = false;
        for (const file of Array.from(files)) {
            try {
                await uploadPiece(caseId, file, notesInput || undefined);
            } catch (err) {
                hasError = true;
                toastState.add(TOAST_LEVELS.Error, "Erreur", `Impossible d'ajouter « ${file.name} ».`);
            }
        }
        uploading = false;
        if (!hasError) {
            notesInput = "";
            showUploadModal = false;
        }
        await loadPieces();
    }

    async function handleUpload() {
        const target = fileInput;
        const files = target?.files;
        if (!files || files.length === 0) return;
        await uploadFiles(files);
        if (target) target.value = "";
    }

    function handleDragOver(e: DragEvent) {
        e.preventDefault();
        isDragging = true;
    }

    function handleDragLeave(e: DragEvent) {
        if (!(e.currentTarget as HTMLElement).contains(e.relatedTarget as Node)) {
            isDragging = false;
        }
    }

    async function handleDrop(e: DragEvent) {
        e.preventDefault();
        isDragging = false;
        const files = e.dataTransfer?.files;
        if (files && files.length > 0) await uploadFiles(files);
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

    <!-- Upload Modal -->
    <Dialog.Root bind:open={showUploadModal}>
        <Dialog.Portal>
            <Dialog.Overlay
                class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
            />
            <Dialog.Content
                class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border flex flex-col sm:max-w-[480px] md:w-full"
            >
                <div class="flex-shrink-0 px-8 pt-8 pb-6">
                    <Dialog.Title class="text-base font-semibold tracking-tight">
                        Ajouter des pièces
                    </Dialog.Title>
                    <Dialog.Description class="text-foreground-alt text-sm mt-1">
                        Ajoutez une note optionnelle, puis sélectionnez les fichiers à joindre.
                    </Dialog.Description>
                </div>

                <Separator.Root class="bg-border-input mx-0 !m-0 block h-px flex-shrink-0" />

                <div class="px-8 py-6">
                    <div class="flex flex-col gap-4">
                        <!-- Drop zone -->
                        <div
                            role="button"
                            tabindex="0"
                            ondragover={handleDragOver}
                            ondragleave={handleDragLeave}
                            ondrop={handleDrop}
                            onclick={() => !uploading && fileInput?.click()}
                            onkeydown={(e) => e.key === "Enter" && !uploading && fileInput?.click()}
                            class="rounded-card border-2 border-dashed flex flex-col items-center justify-center gap-2 py-8 transition-colors {isDragging ? 'border-dark bg-dark/5 cursor-copy' : 'border-border-input hover:border-dark-40 hover:bg-muted/30 cursor-pointer'} {uploading ? 'pointer-events-none' : ''}"
                        >
                            {#if uploading}
                                <Loader2 size={24} class="animate-spin text-muted-foreground" />
                                <p class="text-sm text-muted-foreground">Envoi en cours…</p>
                            {:else if isDragging}
                                <Upload size={24} class="text-foreground" />
                                <p class="text-sm font-medium text-foreground">Déposez les fichiers ici</p>
                            {:else}
                                <Upload size={24} class="text-muted-foreground" />
                                <p class="text-sm text-muted-foreground">
                                    Glissez-déposez ou <span class="text-foreground font-medium underline underline-offset-2">parcourez</span>
                                </p>
                            {/if}
                        </div>

                        <!-- Notes -->
                        <textarea
                            bind:value={notesInput}
                            placeholder="Notes (optionnel)…"
                            rows="2"
                            disabled={uploading}
                            class="rounded-card-sm border border-border-input bg-background placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 w-full px-4 py-3 text-sm transition-colors resize-none disabled:opacity-50"
                        ></textarea>

                        <div class="flex justify-end">
                            <Dialog.Close
                                onclick={() => { notesInput = ""; }}
                                disabled={uploading}
                                class="h-input rounded-input bg-transparent text-foreground border border-border-input hover:bg-muted inline-flex items-center justify-center px-5 text-sm font-semibold active:scale-[0.98] cursor-pointer transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
                            >
                                Annuler
                            </Dialog.Close>
                        </div>
                    </div>
                </div>

                <Dialog.Close
                    onclick={() => { notesInput = ""; }}
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
            onclick={() => { showUploadModal = true; }}
            class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
        >
            <Upload size={16} />
            Ajouter des pièces
        </button>
    </div>

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
                onclick={() => { showUploadModal = true; }}
                class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-4 text-sm font-semibold active:scale-[0.98] cursor-pointer"
            >
                <Upload size={16} />
                Ajouter des pièces
            </button>
        </div>
    {:else}
        <div class="flex-1 min-h-0 overflow-y-auto">
            <div class="flex flex-col divide-y divide-border-input">
                {#each pieces as piece, index (piece.id)}
                    <div
                        class="flex items-center gap-4 py-3 px-2 hover:bg-surface/50 transition-colors group"
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

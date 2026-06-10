<script lang="ts">
    import AISidebar from "./AISidebar.svelte";
    import TextEditor from "./TextEditor.svelte";
    import { requestAITextOperation, type AITextOperation, AI_CONFIG } from "$lib/services/api";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte.js";
    import AITextPreviewModal from "./_aiTextPreviewModal.svelte";

    // Toast context for notifications
    const toast = getToastContext();

    // Panel resizing state
    const MIN_AI_PANEL_WIDTH = 400;
    const MAX_AI_PANEL_WIDTH = 800;
    const DEFAULT_AI_PANEL_WIDTH = 400;

    let aiPanelWidth = $state(DEFAULT_AI_PANEL_WIDTH);
    let isDragging = $state(false);
    let containerRef: HTMLDivElement;

    // AI preview modal state
    let aiModalOpen = $state(false);
    let aiModalOperation = $state<AITextOperation | null>(null);
    let aiModalOriginalText = $state("");
    let aiModalSuggestedText = $state<string | null>(null);
    let aiModalIsLoading = $state(false);
    let aiModalError = $state<string | null>(null);
    let aiSelectionRange = $state<{ from: number; to: number } | null>(null);

    // Text editor reference for calling replaceSelectedText
    let textEditorRef: {
        replaceSelectedText: (text: string, range: { from: number; to: number }) => void;
    } | undefined;

    function handleMouseDown(e: MouseEvent) {
        isDragging = true;
        e.preventDefault();
    }

    function handleMouseMove(e: MouseEvent) {
        if (!isDragging || !containerRef) return;

        const containerRect = containerRef.getBoundingClientRect();
        const newWidth = containerRect.right - e.clientX;

        aiPanelWidth = Math.max(
            MIN_AI_PANEL_WIDTH,
            Math.min(MAX_AI_PANEL_WIDTH, newWidth),
        );
    }

    function handleMouseUp() {
        isDragging = false;
    }

    // AI operation handlers
    async function handleAIOperation(
        operation: AITextOperation,
        selectedText: string,
        selectionRange: { from: number; to: number }
    ) {
        const length = selectedText.trim().length;
        if (length < AI_CONFIG.MIN_SELECTION_LENGTH || length > AI_CONFIG.MAX_SELECTION_LENGTH) {
            toast.add(
                "error",
                "Sélection invalide",
                "Veuillez sélectionner entre 3 et 5000 caractères."
            );
            return;
        }

        // Store selection range for later replacement
        aiSelectionRange = selectionRange;

        // Open modal with loading state
        aiModalOpen = true;
        aiModalOperation = operation;
        aiModalOriginalText = selectedText;
        aiModalSuggestedText = null;
        aiModalIsLoading = true;
        aiModalError = null;

        try {
            const response = await requestAITextOperation({
                text: aiModalOriginalText,
                operation,
                language: "fr", // TODO: get from user preferences
            });

            aiModalSuggestedText = response.result;
        } catch (error) {
            const message = error instanceof Error ? error.message : "Une erreur inconnue s'est produite";
            aiModalError = message;
            toast.add("error", "Opération IA échouée", message);
        } finally {
            aiModalIsLoading = false;
        }
    }

    function handleAIAccept() {
        if (aiModalSuggestedText && textEditorRef && aiSelectionRange) {
            textEditorRef.replaceSelectedText(aiModalSuggestedText, aiSelectionRange);
            toast.add("success", "Texte mis à jour", "La suggestion a été appliquée.");
        }
        aiModalOpen = false;
    }

    function handleAIRetry() {
        if (aiModalOperation && aiModalOriginalText && aiSelectionRange) {
            handleAIOperation(aiModalOperation, aiModalOriginalText, aiSelectionRange);
        }
    }
</script>

<svelte:window onmousemove={handleMouseMove} onmouseup={handleMouseUp} />

<div class="flex h-full" bind:this={containerRef}>
    <!-- Left panel: Text editor (flexible width) -->
    <div class="flex-1 min-w-0">
        <TextEditor
            bind:this={textEditorRef}
            onAIOperation={handleAIOperation}
        />
    </div>

    <!-- Draggable divider -->
    <div
        class="w-1 cursor-col-resize hover:bg-surface-active transition-colors relative group"
        onmousedown={handleMouseDown}
        role="separator"
        aria-orientation="vertical"
        tabindex="0"
        aria-label="Resize panels"
    >
        <div
            class="absolute inset-y-0 -left-1 -right-1 group-hover:bg-dark-10"
        ></div>
    </div>

    <!-- Right panel: AI Sidebar -->
    <AISidebar
        width={aiPanelWidth}
        minWidth={MIN_AI_PANEL_WIDTH}
        maxWidth={MAX_AI_PANEL_WIDTH}
    />
</div>

<!-- AI Text Preview Modal -->
<AITextPreviewModal
    open={aiModalOpen}
    operation={aiModalOperation}
    originalText={aiModalOriginalText}
    suggestedText={aiModalSuggestedText}
    isLoading={aiModalIsLoading}
    error={aiModalError}
    onOpenChange={(open) => {
        if (!open) aiModalOpen = false;
    }}
    onAccept={handleAIAccept}
    onRetry={handleAIRetry}
/>

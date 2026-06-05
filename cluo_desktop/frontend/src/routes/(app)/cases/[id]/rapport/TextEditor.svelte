<script lang="ts">
    import { Editor } from "@tiptap/core";
    import StarterKit from "@tiptap/starter-kit";
    import Placeholder from "@tiptap/extension-placeholder";
    import Underline from "@tiptap/extension-underline";
    import { PaginationPlus } from "tiptap-pagination-plus";
    import { onMount, onDestroy } from "svelte";
    import {
        Bold,
        Italic,
        Underline as UnderlineIcon,
        Strikethrough,
        List,
        ListOrdered,
        Quote,
        Undo,
        Redo,
        Heading1,
        Heading2,
        Heading3,
    } from "@lucide/svelte";
    import type { AITextOperation } from "$lib/services/api";
    import AIFloatingMenu from "./_aiFloatingMenu.svelte";

    interface Props {
        onAIOperation?: (
            operation: AITextOperation,
            selectedText: string,
            selectionRange: { from: number; to: number }
        ) => void;
    }

    let { onAIOperation }: Props = $props();

    let editorElement: HTMLElement;
    let editor: Editor;
    let editorState = $state({
        isBold: false,
        isItalic: false,
        isUnderline: false,
        isStrike: false,
        isBulletList: false,
        isOrderedList: false,
        isBlockquote: false,
        isH1: false,
        isH2: false,
        isH3: false,
    });

    // AI operation state
    let selectedText = $state("");
    let selectionRange = $state<{ from: number; to: number } | null>(null);

    onMount(() => {
        editor = new Editor({
            element: editorElement,
            extensions: [
                StarterKit.configure({
                    heading: {
                        levels: [1, 2, 3],
                    },
                }),
                Placeholder.configure({
                    placeholder: "Start writing your rapport...",
                }),
                Underline,
                PaginationPlus.configure({
                    pageWidth: 794,
                    pageHeight: 1123,
                    pageGap: 48,
                    marginTop: 76,
                    marginBottom: 76,
                    marginLeft: 95,
                    marginRight: 95,
                    // NOTE: The pageBreakBackground property should be in sync with the div that contains the text editor (it has an overflow-y property)
                    pageBreakBackground: "var(--color-muted)",
                }),
            ],
            content: "",
            editorProps: {
                attributes: {
                    class: "prose prose-sm max-w-none focus:outline-none editor-content",
                },
            },
            onUpdate: ({ editor }) => {
                updateEditorState();
            },
            onSelectionUpdate: ({ editor }) => {
                updateEditorState();
            },
        });
    });

    function updateEditorState() {
        if (!editor) return;
        editorState.isBold = editor.isActive("bold");
        editorState.isItalic = editor.isActive("italic");
        editorState.isUnderline = editor.isActive("underline");
        editorState.isStrike = editor.isActive("strike");
        editorState.isBulletList = editor.isActive("bulletList");
        editorState.isOrderedList = editor.isActive("orderedList");
        editorState.isBlockquote = editor.isActive("blockquote");
        editorState.isH1 = editor.isActive("heading", { level: 1 });
        editorState.isH2 = editor.isActive("heading", { level: 2 });
        editorState.isH3 = editor.isActive("heading", { level: 3 });

        // Track selection for AI operations
        const { from, to, empty } = editor.state.selection;
        if (!empty && from !== to) {
            selectionRange = { from, to };
            selectedText = editor.state.doc.textBetween(from, to, " ");
        } else {
            selectionRange = null;
            selectedText = "";
        }
    }

    onDestroy(() => {
        if (editor) {
            editor.destroy();
        }
    });

    function toggleBold() {
        editor.chain().focus().toggleBold().run();
    }

    function toggleItalic() {
        editor.chain().focus().toggleItalic().run();
    }

    function toggleUnderline() {
        editor.chain().focus().toggleUnderline().run();
    }

    function toggleStrike() {
        editor.chain().focus().toggleStrike().run();
    }

    function toggleBulletList() {
        editor.chain().focus().toggleBulletList().run();
    }

    function toggleOrderedList() {
        editor.chain().focus().toggleOrderedList().run();
    }

    function toggleBlockquote() {
        editor.chain().focus().toggleBlockquote().run();
    }

    function setHeading(level: 1 | 2 | 3) {
        editor.chain().focus().toggleHeading({ level }).run();
    }

    function undo() {
        editor.chain().focus().undo().run();
    }

    function redo() {
        editor.chain().focus().redo().run();
    }

    // Replace selected text with AI suggestion (callable from parent)
    export function replaceSelectedText(newText: string, range: { from: number; to: number }) {
        editor
            .chain()
            .focus()
            .insertContentAt({ from: range.from, to: range.to }, newText)
            .run();
    }
</script>

<div class="flex flex-col h-full w-full">
    <!-- Toolbar -->
    <div
        class="flex justify-center items-center gap-1 p-2 border-b border-border-card bg-background sticky top-0 z-10"
    >
        <!-- Text formatting -->
        <div class="flex items-center gap-0.5">
            <button
                onclick={toggleBold}
                class="p-1.5 rounded hover:bg-muted transition-colors {editorState.isBold
                    ? 'bg-muted'
                    : ''}"
                title="Bold"
                type="button"
            >
                <Bold class="w-5 h-5" />
            </button>
            <button
                onclick={toggleItalic}
                class="p-1.5 rounded hover:bg-muted transition-colors {editorState.isItalic
                    ? 'bg-muted'
                    : ''}"
                title="Italic"
                type="button"
            >
                <Italic class="w-5 h-5" />
            </button>
            <button
                onclick={toggleUnderline}
                class="p-1.5 rounded hover:bg-muted transition-colors {editorState.isUnderline
                    ? 'bg-muted'
                    : ''}"
                title="Underline"
                type="button"
            >
                <UnderlineIcon class="w-5 h-5" />
            </button>
            <button
                onclick={toggleStrike}
                class="p-1.5 rounded hover:bg-muted transition-colors {editorState.isStrike
                    ? 'bg-muted'
                    : ''}"
                title="Strikethrough"
                type="button"
            >
                <Strikethrough class="w-5 h-5" />
            </button>
        </div>

        <div class="w-px h-6 bg-border-card mx-1"></div>

        <!-- Headings -->
        <div class="flex items-center gap-0.5">
            <button
                onclick={() => setHeading(1)}
                class="p-1.5 rounded hover:bg-muted transition-colors {editorState.isH1
                    ? 'bg-muted'
                    : ''}"
                title="Heading 1"
                type="button"
            >
                <Heading1 class="w-5 h-5" />
            </button>
            <button
                onclick={() => setHeading(2)}
                class="p-1.5 rounded hover:bg-muted transition-colors {editorState.isH2
                    ? 'bg-muted'
                    : ''}"
                title="Heading 2"
                type="button"
            >
                <Heading2 class="w-5 h-5" />
            </button>
            <button
                onclick={() => setHeading(3)}
                class="p-1.5 rounded hover:bg-muted transition-colors {editorState.isH3
                    ? 'bg-muted'
                    : ''}"
                title="Heading 3"
                type="button"
            >
                <Heading3 class="w-5 h-5" />
            </button>
        </div>

        <div class="w-px h-6 bg-border-card mx-1"></div>

        <!-- Lists and quotes -->
        <div class="flex items-center gap-0.5">
            <button
                onclick={toggleBulletList}
                class="p-1.5 rounded hover:bg-muted transition-colors {editorState.isBulletList
                    ? 'bg-muted'
                    : ''}"
                title="Bullet List"
                type="button"
            >
                <List class="w-5 h-5" />
            </button>
            <button
                onclick={toggleOrderedList}
                class="p-1.5 rounded hover:bg-muted transition-colors {editorState.isOrderedList
                    ? 'bg-muted'
                    : ''}"
                title="Ordered List"
                type="button"
            >
                <ListOrdered class="w-5 h-5" />
            </button>
            <button
                onclick={toggleBlockquote}
                class="p-1.5 rounded hover:bg-muted transition-colors {editorState.isBlockquote
                    ? 'bg-muted'
                    : ''}"
                title="Quote"
                type="button"
            >
                <Quote class="w-5 h-5" />
            </button>
        </div>

        <div class="w-px h-6 bg-border-card mx-1"></div>

        <!-- Undo/Redo -->
        <div class="flex items-center gap-0.5">
            <button
                onclick={undo}
                class="p-1.5 rounded hover:bg-muted transition-colors"
                title="Undo"
                type="button"
            >
                <Undo class="w-5 h-5" />
            </button>
            <button
                onclick={redo}
                class="p-1.5 rounded hover:bg-muted transition-colors"
                title="Redo"
                type="button"
            >
                <Redo class="w-5 h-5" />
            </button>
        </div>
    </div>

    <!-- AI Floating Menu -->
    {#if selectionRange && selectedText.trim().length >= 3}
        <AIFloatingMenu
            {editor}
            onOperationClick={(op) => onAIOperation?.(op, selectedText, selectionRange)}
        />
    {/if}

    <!-- Editor -->
    <!-- <div class="flex-1 overflow-y-auto bg-muted"> -->
    <div class="flex-1 overflow-y-auto bg-muted">
        <div class="flex justify-center py-8">
            <div bind:this={editorElement} class="tiptap-editor"></div>
        </div>
    </div>
</div>

<style>
    .tiptap-editor {
        background: var(--background);
        box-shadow: 0 0 10px var(--color-dark-10);
    }

    :global(.tiptap-editor .ProseMirror) {
        outline: none;
        color: var(--color-foreground);
        font-family: var(--font-sans);
        line-height: 1.6;
    }

    /* Style pages created by pagination extension */
    :global(.tiptap-editor .page) {
        background: var(--background) !important;
        box-shadow: 0 0 10px var(--color-dark-10);
        color: var(--color-foreground);
    }

    /* Style page gaps/breaks to match container background */
    :global(.tiptap-editor .page-break) {
        background: var(--color-muted) !important;
    }

    :global(.tiptap-editor .ProseMirror p.is-editor-empty:first-child::before) {
        content: attr(data-placeholder);
        color: var(--color-muted-foreground);
        float: left;
        height: 0;
        pointer-events: none;
    }

    :global(.tiptap-editor .ProseMirror h1) {
        font-size: 2em;
        font-weight: 700;
        margin-top: 1em;
        margin-bottom: 0.5em;
        color: var(--color-foreground);
        line-height: 1.2;
    }

    :global(.tiptap-editor .ProseMirror h2) {
        font-size: 1.5em;
        font-weight: 600;
        margin-top: 0.8em;
        margin-bottom: 0.4em;
        color: var(--color-foreground);
        line-height: 1.3;
    }

    :global(.tiptap-editor .ProseMirror h3) {
        font-size: 1.25em;
        font-weight: 600;
        margin-top: 0.6em;
        margin-bottom: 0.3em;
        color: var(--color-foreground);
        line-height: 1.4;
    }

    :global(.tiptap-editor .ProseMirror p) {
        margin-top: 0.75em;
        margin-bottom: 0.75em;
    }

    :global(.tiptap-editor .ProseMirror ul),
    :global(.tiptap-editor .ProseMirror ol) {
        padding-left: 1.5em;
        margin-top: 0.75em;
        margin-bottom: 0.75em;
    }

    :global(.tiptap-editor .ProseMirror ul) {
        list-style-type: disc;
    }

    :global(.tiptap-editor .ProseMirror ol) {
        list-style-type: decimal;
    }

    :global(.tiptap-editor .ProseMirror li) {
        margin-top: 0.25em;
        margin-bottom: 0.25em;
    }

    :global(.tiptap-editor .ProseMirror blockquote) {
        border-left: 3px solid var(--color-border-input);
        padding-left: 1em;
        margin-left: 0;
        margin-top: 1em;
        margin-bottom: 1em;
        color: var(--color-foreground-alt);
        font-style: italic;
    }

    :global(.tiptap-editor .ProseMirror code) {
        background-color: var(--color-muted);
        color: var(--color-foreground);
        padding: 0.2em 0.4em;
        border-radius: 3px;
        font-family: var(--font-mono);
        font-size: 0.9em;
    }

    :global(.tiptap-editor .ProseMirror pre) {
        background-color: var(--color-muted);
        color: var(--color-foreground);
        padding: 1em;
        border-radius: 6px;
        overflow-x: auto;
        margin-top: 1em;
        margin-bottom: 1em;
    }

    :global(.tiptap-editor .ProseMirror pre code) {
        background-color: transparent;
        padding: 0;
        font-size: 0.9em;
    }

    :global(.tiptap-editor .ProseMirror strong) {
        font-weight: 700;
        color: var(--color-foreground);
    }

    :global(.tiptap-editor .ProseMirror em) {
        font-style: italic;
    }

    :global(.tiptap-editor .ProseMirror u) {
        text-decoration: underline;
    }

    :global(.tiptap-editor .ProseMirror s) {
        text-decoration: line-through;
    }
</style>

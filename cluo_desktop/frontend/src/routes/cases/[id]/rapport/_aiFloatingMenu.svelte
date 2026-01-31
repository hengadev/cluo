<script lang="ts">
	import { Editor } from "@tiptap/core";
	import { onMount } from "svelte";
	import type { AITextOperation } from "$lib/services/api";
	import { AI_OPERATION_LABELS } from "$lib/services/api";
	import { RefreshCw, FileText, Briefcase, Sparkles } from "@lucide/svelte";
	import AIOperationButton from "./_aiOperationButton.svelte";

	interface Props {
		editor: Editor;
		onOperationClick: (operation: AITextOperation) => void;
	}

	let { editor, onOperationClick }: Props = $props();

	let menuPosition = $state({ top: 0, left: 0 });
	let isVisible = $state(false);
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;

	// Get icon for operation
	function getIcon(operation: AITextOperation) {
		switch (operation) {
			case "reword":
				return RefreshCw;
			case "summarize":
				return FileText;
			case "formalize":
				return Briefcase;
			case "clarify":
				return Sparkles;
		}
	}

	// Calculate position from Tiptap selection coordinates
	function updatePosition() {
		const { from, to } = editor.state.selection;
		if (from === to) {
			isVisible = false;
			return;
		}

		// Get bounding rect of selection
		const start = editor.view.coordsAtPos(from);
		const end = editor.view.coordsAtPos(to);

		// Position menu above selection, centered
		const editorRect = editor.view.dom.getBoundingClientRect();
		const menuWidth = 200; // approximate width
		const menuHeight = 50; // approximate height

		menuPosition = {
			top: start.top - editorRect.top - menuHeight - 8, // 8px gap
			left:
				Math.min(start.left, end.left) - editorRect.left + Math.abs(end.left - start.left) / 2 - menuWidth / 2
		};

		// Clamp to viewport bounds
		menuPosition.left = Math.max(8, Math.min(menuPosition.left, editorRect.width - menuWidth - 8));
		menuPosition.top = Math.max(8, menuPosition.top);
	}

	// Debounced visibility
	$effect(() => {
		const { empty } = editor.state.selection;
		if (debounceTimer) clearTimeout(debounceTimer);

		if (!empty) {
			debounceTimer = setTimeout(() => {
				updatePosition();
				isVisible = true;
			}, 200);
		} else {
			isVisible = false;
		}
	});

	onMount(() => {
		// Cleanup timer on unmount
		return () => {
			if (debounceTimer) clearTimeout(debounceTimer);
		};
	});

	function handleOperationClick(operation: AITextOperation) {
		onOperationClick(operation);
		isVisible = false;
	}
</script>

{#if isVisible}
	<div
		class="ai-floating-menu"
		style="top: {menuPosition.top}px; left: {menuPosition.left}px;"
	>
		<div class="ai-floating-menu__inner">
			{#each Object.entries(AI_OPERATION_LABELS) as [operation, { label }]}
				<AIOperationButton
					operation={operation as AITextOperation}
					{label}
					icon={getIcon(operation as AITextOperation)}
					onclick={() => handleOperationClick(operation as AITextOperation)}
				/>
			{/each}
		</div>
	</div>
{/if}

<style>
	.ai-floating-menu {
		position: absolute;
		z-index: 100;
		pointer-events: none;
	}

	.ai-floating-menu__inner {
		display: flex;
		gap: 0.25rem;
		padding: 0.375rem;
		background: var(--color-background);
		border: 1px solid var(--color-border-card);
		border-radius: 0.75rem;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
		pointer-events: auto;
		animation: fadeIn 0.15s ease-out;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(-4px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	:global(.dark) .ai-floating-menu__inner {
		background: var(--color-background-alt);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
	}
</style>

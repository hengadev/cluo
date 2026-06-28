<script lang="ts">
	import { Check, X } from "@lucide/svelte";

	interface Props {
		text: string;
		readonly?: boolean;
		class?: string;
		onSave?: (text: string) => void;
		onCancel?: () => void;
	}

	let {
		text,
		readonly = false,
		class: className = "",
		onSave,
		onCancel,
	}: Props = $props();

	let editedText = $state(text);
	let wordCount = $derived(editedText.trim().split(/\s+/).filter(Boolean).length);
	let charCount = $derived(editedText.length);
	const maxChars = 100000;

	// Update edited text when prop changes
	$effect(() => {
		editedText = text;
	});

	function handleSave() {
		if (onSave) {
			onSave(editedText);
		}
	}

	function handleCancel() {
		editedText = text;
		if (onCancel) {
			onCancel();
		}
	}

	// Method to get current text (can be called by parent via ref)
	function getText() {
		return editedText;
	}

	// Export methods for parent component access
	export {
		getText,
	};
</script>

<div class="transcript-editor flex flex-col gap-3 {className}">
	<div class="flex flex-col gap-2">
		<div class="flex justify-between items-center">
			<label for="transcript-textarea" class="text-dark-700 font-semibold text-sm">
				Transcription
			</label>
			<div class="flex gap-3 text-dark-500 text-xs">
				<span>{wordCount} mots</span>
				<span>{charCount} / {maxChars} caractères</span>
			</div>
		</div>
		<textarea
			bind:value={editedText}
			id="transcript-textarea"
			readonly={readonly}
			maxlength={maxChars}
			class="w-full min-h-48 p-4 border border-dark-200 rounded-xl bg-background text-dark-800 text-base leading-relaxed resize-y focus:outline-none focus:ring-2 focus:ring-foreground placeholder:text-dark-400"
			placeholder="La transcription apparaîtra ici..."
			class:readonly={readonly}
		></textarea>
	</div>

	{#if !readonly}
		<div class="flex gap-3">
			<button
				onclick={handleSave}
				class="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-accent hover:bg-accent/90 text-accent-foreground rounded-xl transition-colors font-medium"
			>
				<Check size={18} />
				<span>Enregistrer</span>
			</button>
			<button
				onclick={handleCancel}
				class="flex-1 flex items-center justify-center gap-2 px-4 py-3 bg-dark-100 hover:bg-dark-200 text-dark-700 rounded-xl transition-colors font-medium"
			>
				<X size={18} />
				<span>Annuler</span>
			</button>
		</div>
	{/if}
</div>

<style>
	.transcript-editor textarea.readonly {
		cursor: default;
		background-color: var(--color-dark-50);
	}
</style>

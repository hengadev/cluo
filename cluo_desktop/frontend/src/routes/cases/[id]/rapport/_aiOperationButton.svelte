<script lang="ts">
	import type { Component } from "svelte";
	import type { AITextOperation } from "$lib/services/api";
	import { Loader2 } from "@lucide/svelte";

	interface Props {
		operation: AITextOperation;
		label: string;
		icon: Component;
		isLoading?: boolean;
		disabled?: boolean;
		onclick: () => void;
	}

	let {
		operation,
		label,
		icon: Icon,
		isLoading = false,
		disabled = false,
		onclick
	}: Props = $props();
</script>

<button
	class="ai-operation-button"
	class:loading={isLoading}
	class:disabled={disabled || isLoading}
	{disabled}
	onclick={onclick}
	type="button"
	title="{label}"
>
	{#if isLoading}
		<Loader2 class="spinner" size={16} />
	{:else}
		<Icon size={16} />
	{/if}
	<span>{label}</span>
</button>

<style>
	.ai-operation-button {
		display: flex;
		align-items: center;
		gap: 0.375rem;
		padding: 0.5rem 0.75rem;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		font-weight: 500;
		transition: all 0.2s ease;
		cursor: pointer;
		border: none;
		background: transparent;
		color: var(--color-foreground);
	}

	.ai-operation-button:hover:not(.disabled):not(.loading) {
		background: var(--color-muted);
		transform: scale(1.02);
	}

	.ai-operation-button:active:not(.disabled):not(.loading) {
		transform: scale(0.98);
	}

	.ai-operation-button.disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.ai-operation-button.loading {
		opacity: 0.8;
		cursor: wait;
	}

	.spinner {
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}
</style>

<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { MessageSquare, Send, Trash2, Plus, Loader2 } from '@lucide/svelte';
	import { chatStore, currentMessages, currentConversation } from '$lib/stores/chat';
	import { currentCase } from '$lib/stores/case';
	import { getToastContext } from '$lib/custom/global/toast/state.svelte.js';
	import type { ChatMessage } from '$lib/types/chat';

	interface Props {
		width: number;
	}

	let { width }: Props = $props();
	const toast = getToastContext();

	// Local state
	let inputMessage = $state('');
	let textareaElement: HTMLTextAreaElement;
	let messagesEnd: HTMLElement;
	let isDeleting = $state<string | null>(null);

	// Subscribe to stores
	let conversations = $derived.by(() => $chatStore.conversations);
	let currentConversationId = $derived.by(() => $chatStore.currentConversationId);
	let messages = $derived($currentMessages);
	let isLoading = $derived.by(() => $chatStore.isLoading);
	let isSending = $derived.by(() => $chatStore.isSending);
	let error = $derived.by(() => $chatStore.error);

	// Get current case ID
	let caseId = $derived.by(() => $currentCase.id || $page.params.id);

	// Load conversations on mount
	onMount(() => {
		if (caseId) {
			chatStore.loadConversations(caseId);
		}
	});

	// Watch for case changes
	$effect(() => {
		if (caseId) {
			chatStore.loadConversations(caseId);
		}
	});

	// Scroll to bottom when messages change
	$effect(() => {
		if (messagesEnd && messages.length > 0) {
			messagesEnd.scrollIntoView({ behavior: 'smooth' });
		}
	});

	function handleSend() {
		if (!inputMessage.trim() || isSending || !caseId) return;

		const message = inputMessage;
		inputMessage = '';

		// Reset textarea height
		if (textareaElement) {
			textareaElement.style.height = 'auto';
		}

		chatStore.sendMessage(caseId, message).catch((err) => {
			toast.add('error', "Échec de l'envoi", err.message);
		});
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			handleSend();
		}
	}

	function autoResizeTextarea(element: HTMLTextAreaElement) {
		element.style.height = 'auto';
		element.style.height = element.scrollHeight + 'px';
	}

	function selectConversation(conversationId: string) {
		chatStore.setCurrentConversation(conversationId);
		chatStore.loadConversation(conversationId);
	}

	function startNewChat() {
		chatStore.setCurrentConversation(null);
	}

	function deleteConversation(conversationId: string) {
		isDeleting = conversationId;
		chatStore.deleteConversation(conversationId).finally(() => {
			isDeleting = null;
		});
	}

	function getRoleClass(role: ChatMessage['role']) {
		switch (role) {
			case 'user':
				return 'bg-primary text-primary-foreground ml-auto';
			case 'assistant':
				return 'bg-muted text-foreground';
			default:
				return 'bg-muted-foreground/50 text-foreground';
		}
	}

	function formatTime(isoString: string): string {
		const date = new Date(isoString);
		return date.toLocaleTimeString('fr-FR', {
			hour: '2-digit',
			minute: '2-digit',
		});
	}
</script>

<div class="flex flex-col h-full bg-background">
	<!-- Header -->
	<div class="flex items-center justify-between px-4 py-3 border-b border-border-card bg-muted/30">
		<h2 class="font-semibold text-base flex items-center gap-2">
			<MessageSquare class="w-4 h-4 text-primary" />
			Assistant IA
		</h2>
		{#if currentConversationId}
			<button
				onclick={startNewChat}
				class="p-1.5 rounded-lg hover:bg-muted transition-colors"
				title="Nouvelle conversation"
			>
				<Plus class="w-4 h-4" />
			</button>
		{/if}
	</div>

	<!-- Conversations List (when no conversation selected) -->
	{#if !currentConversationId && conversations.length > 0}
		<div class="flex-1 overflow-y-auto p-3">
			<h3 class="text-xs font-medium text-muted-foreground px-2 mb-2 uppercase tracking-wide">
				Conversations
			</h3>
			{#each conversations as conv (conv.id)}
				<div
					class="group flex items-center justify-between p-3 rounded-lg hover:bg-muted cursor-pointer mb-1 transition-colors"
					class:bg-muted={currentConversationId === conv.id}
					onclick={() => selectConversation(conv.id)}
				>
					<div class="flex-1 min-w-0">
						<p class="font-medium text-sm truncate">{conv.title}</p>
						<p class="text-xs text-muted-foreground">
							{conv.messageCount} message{conv.messageCount > 1 ? 's' : ''}
						</p>
					</div>
					<button
						onclick={(e) => {
							e.stopPropagation();
							deleteConversation(conv.id);
						}}
						disabled={isDeleting === conv.id}
						class="p-1 rounded-lg hover:bg-destructive/10 text-destructive opacity-0 group-hover:opacity-100 transition-opacity"
					>
						{#if isDeleting === conv.id}
							<Loader2 class="w-3.5 h-3.5 animate-spin" />
						{:else}
							<Trash2 class="w-3.5 h-3.5" />
						{/if}
					</button>
				</div>
			{/each}
		</div>
	{:else}
		<!-- Messages Area -->
		<div class="flex-1 overflow-y-auto p-4 space-y-4">
			{#if isLoading}
				<div class="flex items-center justify-center h-full">
					<Loader2 class="w-8 h-8 animate-spin text-muted-foreground" />
				</div>
			{:else if messages.length === 0}
				<div class="flex items-center justify-center h-full text-center px-4">
					<div>
						<MessageSquare class="w-12 h-12 mx-auto mb-4 text-muted-foreground/50" />
						<p class="text-muted-foreground text-sm">
							{currentConversationId
								? 'Commencez la conversation...'
								: 'Entrez votre message ci-dessous pour commencer.'}
						</p>
					</div>
				</div>
			{:else}
				{#each messages as msg (msg.id)}
					<div class="flex flex-col gap-1.5">
						<div class="flex items-center gap-2 text-xs text-muted-foreground px-1">
							{msg.role === 'user' ? 'Vous' : 'Assistant'}
						</div>
						<div
							class="max-w-[85%] p-3 rounded-2xl {getRoleClass(msg.role)}"
						>
							<p class="text-sm whitespace-pre-wrap break-words leading-relaxed">
								{msg.content}
							</p>
						</div>
					</div>
				{/each}

				<!-- Loading indicator for assistant response -->
				{#if isSending}
					<div class="flex flex-col gap-1.5">
						<div class="flex items-center gap-2 text-xs text-muted-foreground px-1">
							Assistant
						</div>
						<div class="bg-muted text-foreground max-w-[85%] p-3 rounded-2xl">
							<div class="flex items-center gap-2">
								<Loader2 class="w-4 h-4 animate-spin" />
								<p class="text-sm">En train d'écrire...</p>
							</div>
						</div>
					</div>
				{/if}
			{/if}

			{#if error}
				<div class="p-3 bg-destructive/10 text-destructive rounded-lg text-sm">
					{error}
				</div>
			{/if}

			<div bind:this={messagesEnd}></div>
		</div>

		<!-- Input Area -->
		<div class="p-3 border-t border-border-card bg-muted/20">
			<div class="flex gap-2 items-end">
				<textarea
					bind:this={textareaElement}
					bind:value={inputMessage}
					onkeydown={handleKeydown}
					oninput={(e) => autoResizeTextarea(e.currentTarget)}
					placeholder="Écrivez votre message... (Entrée pour envoyer)"
					disabled={isSending}
					rows="1"
					class="flex-1 resize-none p-3 rounded-xl border border-border-input bg-background focus:outline-none focus:ring-2 focus:ring-primary/50 min-h-[48px] max-h-[200px] text-sm"
				></textarea>
				<button
					onclick={handleSend}
					disabled={!inputMessage.trim() || isSending || !caseId}
					class="p-3 rounded-xl bg-primary text-primary-foreground hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-sm"
					type="button"
				>
					{#if isSending}
						<Loader2 class="w-5 h-5 animate-spin" />
					{:else}
						<Send class="w-5 h-5" />
					{/if}
				</button>
			</div>
		</div>
	{/if}
</div>

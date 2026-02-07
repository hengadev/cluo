import { writable, derived, get } from 'svelte/store';
import type { ChatConversation, ChatMessage } from '$lib/types/chat';
import * as api from '$lib/services/api';

interface ChatState {
	conversations: ChatConversation[];
	currentConversationId: string | null;
	messages: Map<string, ChatMessage[]>; // conversationId -> messages
	isLoading: boolean;
	isSending: boolean;
	isStreaming: boolean;
	pendingMessage: string; // For streaming responses
	error: string | null;
}

const initialState: ChatState = {
	conversations: [],
	currentConversationId: null,
	messages: new Map(),
	isLoading: false,
	isSending: false,
	isStreaming: false,
	pendingMessage: '',
	error: null,
};

function createChatStore() {
	const { subscribe, update, set } = writable<ChatState>(initialState);

	return {
		subscribe,

		// Actions
		async loadConversations(caseId: string) {
			update((s) => ({ ...s, isLoading: true, error: null }));
			try {
				const { conversations } = await api.listChatConversations(caseId);
				update((s) => ({
					...s,
					conversations,
					isLoading: false,
				}));
			} catch (error) {
				update((s) => ({
					...s,
					error: error instanceof Error ? error.message : 'Failed to load conversations',
					isLoading: false,
				}));
			}
		},

		async loadConversation(conversationId: string) {
			update((s) => ({ ...s, isLoading: true, error: null }));
			try {
				const { conversation, messages } = await api.getChatConversation(conversationId);

				update((s) => {
					const newMessages = new Map(s.messages);
					newMessages.set(conversationId, messages);

					return {
						...s,
						currentConversationId: conversationId,
						messages: newMessages,
						isLoading: false,
					};
				});
			} catch (error) {
				update((s) => ({
					...s,
					error: error instanceof Error ? error.message : 'Failed to load conversation',
					isLoading: false,
				}));
			}
		},

		async sendMessage(caseId: string, message: string) {
			const currentState = get({ subscribe });
			const currentConvId = currentState.currentConversationId;

			update((s) => ({ ...s, isSending: true, error: null }));

			// Optimistically add user message
			const tempUserMessage: ChatMessage = {
				id: `temp-${Date.now()}`,
				conversationId: currentConvId || '',
				role: 'user',
				content: message,
				createdAt: new Date().toISOString(),
			};

			// Add temporary user message to store
			update((s) => {
				const newMessages = new Map(s.messages);
				const convId = s.currentConversationId;
				if (convId) {
					const existing = newMessages.get(convId) || [];
					newMessages.set(convId, [...existing, tempUserMessage]);
				}
				return { ...s, messages: newMessages };
			});

			try {
				const response = await api.sendChatMessage(caseId, {
					conversationId: currentConvId || undefined,
					message,
				});

				update((s) => {
					const newMessages = new Map(s.messages);

					// Replace temp user message with real one and add assistant message
					const convId = response.conversationId;
					const existing = newMessages.get(convId) || [];

					// Filter out temp message and add real messages
					const filtered = existing.filter((m) => m.id !== tempUserMessage.id);
					newMessages.set(convId, [...filtered, response.assistantMessage]);

					// Update or add conversation
					let conversations = s.conversations;
					if (response.conversation) {
						const index = conversations.findIndex((c) => c.id === response.conversationId);
						if (index >= 0) {
							conversations = conversations.map((c, i) =>
								i === index ? response.conversation! : c,
							);
						} else {
							conversations = [response.conversation, ...conversations];
						}
					}

					return {
						...s,
						conversations,
						currentConversationId: response.conversationId,
						messages: newMessages,
						isSending: false,
					};
				});
			} catch (error) {
				update((s) => ({
					...s,
					error: error instanceof Error ? error.message : 'Failed to send message',
					isSending: false,
				}));
				throw error;
			}
		},

		async deleteConversation(conversationId: string) {
			try {
				await api.deleteChatConversation(conversationId);

				update((s) => {
					const newMessages = new Map(s.messages);
					newMessages.delete(conversationId);

					return {
						...s,
						conversations: s.conversations.filter((c) => c.id !== conversationId),
						messages: newMessages,
						currentConversationId: s.currentConversationId === conversationId ? null : s.currentConversationId,
					};
				});
			} catch (error) {
				update((s) => ({
					...s,
					error: error instanceof Error ? error.message : 'Failed to delete conversation',
				}));
			}
		},

		setCurrentConversation(conversationId: string | null) {
			update((s) => ({ ...s, currentConversationId: conversationId }));
		},

		clearError() {
			update((s) => ({ ...s, error: null }));
		},

		// For streaming support (future enhancement)
		setStreaming(isStreaming: boolean, pendingMessage: string = '') {
			update((s) => ({ ...s, isStreaming, pendingMessage }));
		},
	};
}

export const chatStore = createChatStore();

// Derived store for current messages
export const currentMessages = derived(chatStore, ($chat) => {
	if (!$chat.currentConversationId) return [];
	return $chat.messages.get($chat.currentConversationId) || [];
});

// Derived store for current conversation
export const currentConversation = derived(chatStore, ($chat) => {
	if (!$chat.currentConversationId) return null;
	return $chat.conversations.find((c) => c.id === $chat.currentConversationId) || null;
});

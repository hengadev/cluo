/**
 * Chat types for AI chat feature
 */

export type ChatRole = 'user' | 'assistant' | 'system';

export interface ChatMessage {
	id: string;
	conversationId: string;
	role: ChatRole;
	content: string;
	createdAt: string;
	tokenCount?: number;
}

export interface ChatConversation {
	id: string;
	caseId: string;
	title: string;
	createdAt: string;
	updatedAt: string;
	createdBy: string;
	messageCount: number;
	totalTokens: number;
}

export interface ChatContext {
	caseId: string;
	caseTitle: string;
	caseType: string;
	status: string;
	clientName: string;
	location: string;
	subjects: Array<{ role: string; name: string }>;
	recentNotes: string[];
}

// API Request/Response types
export interface SendMessageRequest {
	conversationId?: string; // Optional for first message
	message: string;
}

export interface SendMessageResponse {
	conversationId: string;
	userMessageId: string;
	assistantMessage: ChatMessage;
	conversation?: ChatConversation; // Included if new conversation created
}

export interface GetConversationRequest {
	conversationId: string;
}

export interface GetConversationResponse {
	conversation: ChatConversation;
	messages: ChatMessage[];
}

export interface ListConversationsResponse {
	conversations: ChatConversation[];
}

// Stream types
export type StreamChatChunk =
	| { type: 'token'; token: string }
	| { type: 'done'; messageId?: string }
	| { type: 'error'; error: string };

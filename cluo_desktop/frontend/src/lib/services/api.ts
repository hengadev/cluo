/**
 * API Service
 * 
 * When VITE_USE_MOCK_DATA is true, returns mock data for development.
 * When false, makes actual API calls to the backend.
 */

import { isMockEnabled } from '../config';
import type {
	User,
	Client,
	Contact,
	Case,
	CaseSubject,
	CaseSubjectAssignment,
	Estimate,
	Mandate,
	Contract,
	Invoice,
} from '../types/entities';

// Import mock data
import * as mockData from '../mockData';

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

/**
 * Simulates API delay for realistic mock behavior
 */
async function mockDelay(ms: number = 100): Promise<void> {
	await new Promise(resolve => setTimeout(resolve, ms));
}

// =============================================================================
// USERS
// =============================================================================

/**
 * Fetch all users
 */
export async function fetchAllUsers(): Promise<User[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getAllUsers();
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch a user by ID
 */
export async function fetchUser(id: string): Promise<User | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getUserById(id) || null;
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

// =============================================================================
// CLIENTS
// =============================================================================

/**
 * Fetch all clients (global list for sidebar)
 */
export async function fetchAllClients(): Promise<Client[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getAllClients();
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch a client by ID
 */
export async function fetchClient(id: string): Promise<Client | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getClientById(id) || null;
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch contacts for a specific client
 */
export async function fetchClientContacts(clientId: string): Promise<Contact[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getContactsByClientId(clientId);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

// =============================================================================
// CONTACTS
// =============================================================================

/**
 * Fetch a contact by ID
 */
export async function fetchContact(id: string): Promise<Contact | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getContactById(id) || null;
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

// =============================================================================
// CASES
// =============================================================================

/**
 * Fetch all cases
 */
export async function fetchAllCases(): Promise<Case[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getAllCases();
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch all cases (alias for backward compatibility)
 * @deprecated Use fetchAllCases instead
 */
export const fetchCases = fetchAllCases;

/**
 * Fetch a case by ID with full details
 */
export async function fetchCase(id: string): Promise<Case | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getCaseById(id) || null;
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch cases by status
 */
export async function fetchCasesByStatus(status: string): Promise<Case[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getCasesByStatus(status as any);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch cases by client ID
 */
export async function fetchCasesByClient(clientId: string): Promise<Case[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getCasesByClientId(clientId);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

// =============================================================================
// CASE SUBJECTS
// =============================================================================

/**
 * Fetch all case subjects (global list)
 */
export async function fetchAllCaseSubjects(): Promise<CaseSubject[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getAllCaseSubjects();
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch a case subject by ID
 */
export async function fetchCaseSubject(id: string): Promise<CaseSubject | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getCaseSubjectById(id) || null;
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch subjects for a specific case with their roles
 */
export async function fetchCaseSubjects(caseId: string): Promise<Array<{ subject: CaseSubject; role: string }>> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getCaseSubjectsWithRoles(caseId);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch subject assignments (junction table)
 */
export async function fetchCaseSubjectAssignments(): Promise<CaseSubjectAssignment[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.caseSubjectAssignments;
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

// =============================================================================
// ESTIMATES
// =============================================================================

/**
 * Fetch all estimates (global list for sidebar)
 */
export async function fetchAllEstimates(): Promise<Estimate[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getAllEstimates();
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch estimates for a specific case
 */
export async function fetchCaseEstimates(caseId: string): Promise<Estimate[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getEstimatesByCaseId(caseId);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch an estimate by ID
 */
export async function fetchEstimate(id: string): Promise<Estimate | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getEstimateById(id) || null;
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch estimates by client ID
 */
export async function fetchClientEstimates(clientId: string): Promise<Estimate[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getEstimatesByClientId(clientId);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

// =============================================================================
// MANDATES
// =============================================================================

/**
 * Fetch all mandates (global list for sidebar)
 */
export async function fetchAllMandates(): Promise<Mandate[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getAllMandates();
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch mandates for a specific case
 */
export async function fetchCaseMandates(caseId: string): Promise<Mandate[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getMandatesByCaseId(caseId);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch a mandate by ID
 */
export async function fetchMandate(id: string): Promise<Mandate | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getMandateById(id) || null;
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch mandates by client ID
 */
export async function fetchClientMandates(clientId: string): Promise<Mandate[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getMandatesByClientId(clientId);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

// =============================================================================
// CONTRACTS
// =============================================================================

/**
 * Fetch all contracts (global list for sidebar)
 */
export async function fetchAllContracts(): Promise<Contract[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getAllContracts();
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch contracts for a specific case
 */
export async function fetchCaseContracts(caseId: string): Promise<Contract[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getContractsByCaseId(caseId);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch a contract by ID
 */
export async function fetchContract(id: string): Promise<Contract | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getContractById(id) || null;
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch contracts by client ID
 */
export async function fetchClientContracts(clientId: string): Promise<Contract[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getContractsByClientId(clientId);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

// =============================================================================
// INVOICES
// =============================================================================

/**
 * Fetch all invoices (global list for sidebar)
 */
export async function fetchAllInvoices(): Promise<Invoice[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getAllInvoices();
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch invoices for a specific case
 */
export async function fetchCaseInvoices(caseId: string): Promise<Invoice[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getInvoicesByCaseId(caseId);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch an invoice by ID
 */
export async function fetchInvoice(id: string): Promise<Invoice | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getInvoiceById(id) || null;
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch invoices by client ID
 */
export async function fetchClientInvoices(clientId: string): Promise<Invoice[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getInvoicesByClientId(clientId);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

/**
 * Fetch invoices by payment status
 */
export async function fetchInvoicesByPaymentStatus(paymentStatus: string): Promise<Invoice[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getInvoicesByPaymentStatus(paymentStatus as any);
	}
	// TODO: Implement actual API call
	throw new Error('API not implemented');
}

// =============================================================================
// IMAGES (Legacy - from existing implementation)
// =============================================================================

export interface ApiImage {
	id: string;
	url: string;
	// Add other image fields as needed
}

/**
 * Fetch images for a specific case
 * TODO: Implement actual API call
 */
export async function fetchCaseImages(caseId: string): Promise<ApiImage[]> {
	if (isMockEnabled()) {
		await mockDelay();
		// Use the existing mock data for images
		// This is handled by the photos/mockData.ts file in routes
		return [];
	}
	// TODO: Implement actual API call
	return [];
}

// =============================================================================
// AI TEXT OPERATIONS (Existing implementation)
// =============================================================================

/**
 * AI text operation types
 */
export type AITextOperation = "reword" | "summarize" | "formalize" | "clarify";
export type AILanguage = "en" | "fr";

/**
 * Request payload for AI text operation
 */
export interface AITextOperationRequest {
	text: string;
	operation: AITextOperation;
	language: AILanguage;
}

/**
 * Response from AI text operation
 */
export interface AITextOperationResponse {
	result: string;
}

/**
 * Configuration constants for AI operations
 */
export const AI_CONFIG = {
	MAX_SELECTION_LENGTH: 5000,
	MIN_SELECTION_LENGTH: 3,
	DEFAULT_TIMEOUT: 30000,
	DEFAULT_LANGUAGE: "fr" as AILanguage
} as const;

/**
 * Operation labels for UI (French)
 */
export const AI_OPERATION_LABELS: Record<AITextOperation, { label: string; description: string }> = {
	reword: { label: "Reformuler", description: "Réécrire avec d'autres mots" },
	summarize: { label: "Résumer", description: "Condenser le texte" },
	formalize: { label: "Formaliser", description: "Rendre plus professionnel" },
	clarify: { label: "Clarifier", description: "Simplifier et clarifier" }
};

/**
 * Send AI text operation request to backend
 * Returns plain-text suggestion
 *
 * @param request - The operation request with text, operation type, and language
 * @returns Promise with the AI-generated result
 * @throws Error if request fails or times out
 */
export async function requestAITextOperation(
	request: AITextOperationRequest
): Promise<AITextOperationResponse> {
	const baseURL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
	const timeout = AI_CONFIG.DEFAULT_TIMEOUT;

	const controller = new AbortController();
	const timeoutId = setTimeout(() => controller.abort(), timeout);

	try {
		const response = await fetch(`${baseURL}/api/ai/text`, {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify(request),
			signal: controller.signal,
		});

		clearTimeout(timeoutId);

		if (!response.ok) {
			throw new Error(`API error: ${response.status} ${response.statusText}`);
		}

		return await response.json();
	} catch (error) {
		clearTimeout(timeoutId);
		if (error instanceof DOMException && error.name === "AbortError") {
			throw new Error("Request timed out. Please try again.");
		}
		throw error;
	}
}

// =============================================================================
// AI CHAT OPERATIONS
// =============================================================================

import type {
	SendMessageRequest,
	SendMessageResponse,
	GetConversationResponse,
	ListConversationsResponse,
	ChatConversation,
	ChatMessage,
} from '../types/chat';

/**
 * Send a chat message
 */
export async function sendChatMessage(
	caseId: string,
	request: SendMessageRequest,
): Promise<SendMessageResponse> {
	const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
	const url = new URL(`${baseURL}/api/ai/chat/message`);
	url.searchParams.set('case_id', caseId);

	const response = await fetch(url.toString(), {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		credentials: 'include',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		const error = await response.json().catch(() => ({ message: response.statusText }));
		throw new Error(error.message || `API error: ${response.status}`);
	}

	return response.json();
}

/**
 * Get a conversation with messages
 */
export async function getChatConversation(
	conversationId: string,
): Promise<GetConversationResponse> {
	const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
	const response = await fetch(`${baseURL}/api/ai/chat/conversations/${conversationId}`, {
		credentials: 'include',
	});

	if (!response.ok) {
		throw new Error(`Failed to get conversation: ${response.status}`);
	}

	return response.json();
}

/**
 * List conversations for a case
 */
export async function listChatConversations(
	caseId: string,
): Promise<ListConversationsResponse> {
	const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
	const url = new URL(`${baseURL}/api/ai/chat/conversations`);
	url.searchParams.set('case_id', caseId);

	const response = await fetch(url.toString(), {
		credentials: 'include',
	});

	if (!response.ok) {
		throw new Error(`Failed to list conversations: ${response.status}`);
	}

	return response.json();
}

/**
 * Delete a conversation
 */
export async function deleteChatConversation(conversationId: string): Promise<void> {
	const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';
	const response = await fetch(`${baseURL}/api/ai/chat/conversations/${conversationId}`, {
		method: 'DELETE',
		credentials: 'include',
	});

	if (!response.ok) {
		throw new Error(`Failed to delete conversation: ${response.status}`);
	}
}

// =============================================================================
// LEGACY TYPES (for backward compatibility)
// =============================================================================

/**
 * @deprecated Use Case interface instead
 */
export interface ApiCase {
	id: string;
	title: string;
	description?: string;
	// Add other case fields as needed
}

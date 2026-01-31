/**
 * API Service
 * Placeholder for future API integration
 * When mock flag is disabled, these functions will fetch real data
 */

// TODO: Implement actual API calls using Wails bindings or fetch

export interface ApiCase {
	id: string;
	title: string;
	description?: string;
	// Add other case fields as needed
}

export interface ApiImage {
	id: string;
	url: string;
	// Add other image fields as needed
}

/**
 * Fetch all cases from the API
 * TODO: Implement actual API call
 */
export async function fetchCases(): Promise<ApiCase[]> {
	// Placeholder - implement when ready
	return [];
}

/**
 * Fetch a specific case by ID
 * TODO: Implement actual API call
 */
export async function fetchCase(id: string): Promise<ApiCase | null> {
	// Placeholder - implement when ready
	return null;
}

/**
 * Fetch images for a specific case
 * TODO: Implement actual API call
 */
export async function fetchCaseImages(caseId: string): Promise<ApiImage[]> {
	// Placeholder - implement when ready
	return [];
}

// =============================================================================
// AI Text Operations
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

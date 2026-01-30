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

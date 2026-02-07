/**
 * Application configuration
 * Central place for environment-based feature flags and settings
 */

/**
 * Base URL for API requests
 */
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * Check if mock data is enabled via environment variable
 * @returns true if VITE_USE_MOCK_DATA is set to 'true', false otherwise
 */
export function isMockEnabled(): boolean {
	return import.meta.env.VITE_USE_MOCK_DATA === 'true';
}

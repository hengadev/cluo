/**
 * Application configuration
 * Central place for environment-based feature flags and settings
 */

/**
 * Check if mock data is enabled via environment variable
 * @returns true if VITE_USE_MOCK_DATA is set to 'true', false otherwise
 */
export function isMockEnabled(): boolean {
	return import.meta.env.VITE_USE_MOCK_DATA === 'true';
}

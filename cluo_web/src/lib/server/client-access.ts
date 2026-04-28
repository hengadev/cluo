export type TokenValidation =
	| { valid: true; caseId: string }
	| { valid: false; reason: 'expired' | 'invalid' | 'unavailable' };

export type CaseSummary = {
	title: string;
	message?: string;
	filesSummary: string[];
	expiresAt: Date;
};

import { env } from '$env/dynamic/private';

const USE_MOCK_DATA = env.USE_MOCK_DATA === 'true';

const MOCK_CASE_ID = 'mock-case-001';

const MOCK_SUMMARY: CaseSummary = {
	title: 'Succession — Famille Martin',
	message: 'Tous les documents ont été vérifiés et compilés pour votre dossier.',
	filesSummary: ['Rapport final', 'Photographies', 'Documents complémentaires'],
	expiresAt: new Date('2026-05-12')
};

export async function validateClientToken(token: string): Promise<TokenValidation> {
	if (USE_MOCK_DATA) {
		if (token === 'expired') return { valid: false, reason: 'expired' };
		if (token === 'unavailable') return { valid: false, reason: 'unavailable' };
		if (token === 'invalid') return { valid: false, reason: 'invalid' };
		return { valid: true, caseId: MOCK_CASE_ID };
	}
	throw new Error('Not implemented');
}

export async function createClientSession(caseId: string): Promise<void> {
	if (USE_MOCK_DATA) return;
	throw new Error('Not implemented');
}

export async function getClientCaseSummary(caseId: string): Promise<CaseSummary> {
	if (USE_MOCK_DATA) return MOCK_SUMMARY;
	throw new Error('Not implemented');
}

export async function streamCaseArchive(caseId: string): Promise<ReadableStream> {
	if (USE_MOCK_DATA) {
		const content = `Mock archive — case: ${caseId}\nDevelopment mode only.`;
		const bytes = new TextEncoder().encode(content);
		return new ReadableStream({
			start(controller) {
				controller.enqueue(bytes);
				controller.close();
			}
		});
	}
	throw new Error('Not implemented');
}

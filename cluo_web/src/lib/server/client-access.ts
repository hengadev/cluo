export type TokenValidation =
	| { valid: true; caseId: string }
	| { valid: false; reason: 'expired' | 'invalid' | 'unavailable' };

export type CaseSummary = {
	title: string;
	message?: string;
	filesSummary: string[];
	expiresAt: Date;
};

export async function validateClientToken(token: string): Promise<TokenValidation> {
	throw new Error('Not implemented');
}

export async function createClientSession(caseId: string): Promise<void> {
	throw new Error('Not implemented');
}

export async function getClientCaseSummary(caseId: string): Promise<CaseSummary> {
	throw new Error('Not implemented');
}

export async function streamCaseArchive(caseId: string): Promise<ReadableStream> {
	throw new Error('Not implemented');
}

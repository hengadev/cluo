/**
 * Unit tests for api.ts — real-API paths only.
 *
 * Every test forces isMockEnabled() to return false so we exercise
 * the code that talks to the live backend.  apiFetch is mocked so
 * no real HTTP requests are made.
 */
import { describe, it, expect, vi, beforeEach } from 'vitest';

// ---------------------------------------------------------------------------
// Mocks — must be set up before importing the module under test
// ---------------------------------------------------------------------------

// Stub the env variable so isMockEnabled() returns false
vi.stubEnv('VITE_USE_MOCK_DATA', 'false');
vi.stubEnv('VITE_API_BASE_URL', 'http://localhost:8080');

// Mock apiFetch — every test can override globalThis.fetch
const mockApiFetch = vi.fn();
vi.mock('../apiFetch', () => ({
	apiFetch: (...args: any[]) => mockApiFetch(...args),
}));

// Mock mockData — not used when mock is disabled but imported by the module
vi.mock('../../mockData', () => ({
	getAllUsers: vi.fn(() => []),
	getUserById: vi.fn(() => undefined),
	getAllClients: vi.fn(() => []),
	getClientById: vi.fn(() => undefined),
	getContactsByClientId: vi.fn(() => []),
	getContactById: vi.fn(() => undefined),
	getAllCases: vi.fn(() => []),
	getCaseById: vi.fn(() => undefined),
	getCasesByClientId: vi.fn(() => []),
	getAllCaseSubjects: vi.fn(() => []),
	getCaseSubjectById: vi.fn(() => undefined),
	getAllCaseTypes: vi.fn(() => []),
	getCaseTypeById: vi.fn(() => undefined),
	getAllEstimates: vi.fn(() => []),
	getEstimatesByCaseId: vi.fn(() => []),
	getEstimatesByClientId: vi.fn(() => []),
	getEstimateById: vi.fn(() => undefined),
	getAllMandates: vi.fn(() => []),
	getMandatesByCaseId: vi.fn(() => []),
	getMandatesByClientId: vi.fn(() => []),
	getMandateById: vi.fn(() => undefined),
	getAllContracts: vi.fn(() => []),
	getContractsByCaseId: vi.fn(() => []),
	getContractsByClientId: vi.fn(() => []),
	getContractById: vi.fn(() => undefined),
	getAllInvoices: vi.fn(() => []),
	getInvoicesByCaseId: vi.fn(() => []),
	getInvoicesByClientId: vi.fn(() => []),
	getInvoiceById: vi.fn(() => undefined),
	getInvoicesByPaymentStatus: vi.fn(() => []),
}));

// Import after mocks
const api = await import('../api');

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const BASE = 'http://localhost:8080';

function jsonResponse(body: any, status = 200) {
	return {
		ok: status >= 200 && status < 300,
		status,
		json: async () => body,
	};
}

// ---------------------------------------------------------------------------
// fetchAllUsers
// ---------------------------------------------------------------------------

describe('fetchAllUsers', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('calls GET /auth/me and returns a one-element array', async () => {
		const user = { id: 'u-1', email: 'test@example.com', role: 'admin' };
		mockApiFetch.mockResolvedValueOnce(jsonResponse(user));

		const result = await api.fetchAllUsers();

		expect(mockApiFetch).toHaveBeenCalledWith(`${BASE}/auth/me`);
		expect(result).toEqual([user]);
		expect(result).toHaveLength(1);
	});

	it('throws on non-OK response', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({}, 401));

		await expect(api.fetchAllUsers()).rejects.toThrow('Failed to fetch current user');
	});
});

// ---------------------------------------------------------------------------
// fetchUser
// ---------------------------------------------------------------------------

describe('fetchUser', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('calls GET /auth/me and returns the user when ID matches', async () => {
		const user = { id: 'u-1', email: 'test@example.com', role: 'admin' };
		mockApiFetch.mockResolvedValueOnce(jsonResponse(user));

		const result = await api.fetchUser('u-1');

		expect(mockApiFetch).toHaveBeenCalledWith(`${BASE}/auth/me`);
		expect(result).toEqual(user);
	});

	it('returns null when the ID does not match the current user', async () => {
		const user = { id: 'u-1', email: 'test@example.com', role: 'admin' };
		mockApiFetch.mockResolvedValueOnce(jsonResponse(user));

		const result = await api.fetchUser('different-id');

		expect(result).toBeNull();
	});

	it('throws on non-OK response', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({}, 500));

		await expect(api.fetchUser('u-1')).rejects.toThrow('Failed to fetch current user');
	});
});

// ---------------------------------------------------------------------------
// fetchCaseImages
// ---------------------------------------------------------------------------

describe('fetchCaseImages', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('calls GET /case/{caseId}/media?type=image and maps to ApiImage[]', async () => {
		const mediaResponse = {
			media: [
				{
					id: 'm-1',
					caseId: 'c-1',
					url: 'https://cdn.example.com/photo1.jpg',
					type: 'image',
					mimeType: 'image/jpeg',
					fileName: 'photo1.jpg',
					fileSize: 1024,
					caption: 'Evidence photo',
					isPublished: true,
					createdAt: '2024-01-01T00:00:00Z',
				},
				{
					id: 'm-2',
					caseId: 'c-1',
					url: 'https://cdn.example.com/photo2.jpg',
					type: 'image',
					mimeType: 'image/jpeg',
					fileName: 'photo2.jpg',
					fileSize: 2048,
					caption: 'Another photo',
					isPublished: false,
					createdAt: '2024-01-02T00:00:00Z',
				},
			],
			pagination: { page: 1, pageSize: 20, totalItems: 2, totalPages: 1 },
		};
		mockApiFetch.mockResolvedValueOnce(jsonResponse(mediaResponse));

		const result = await api.fetchCaseImages('c-1');

		expect(mockApiFetch).toHaveBeenCalledWith(
			'http://localhost:8080/case/c-1/media?type=image',
		);
		expect(result).toEqual([
			{ id: 'm-1', url: 'https://cdn.example.com/photo1.jpg' },
			{ id: 'm-2', url: 'https://cdn.example.com/photo2.jpg' },
		]);
	});

	it('returns empty array when no media items exist', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({
				media: [],
				pagination: { page: 1, pageSize: 20, totalItems: 0, totalPages: 1 },
			}),
		);

		const result = await api.fetchCaseImages('c-1');

		expect(result).toEqual([]);
	});

	it('throws on non-OK response', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({}, 500));

		await expect(api.fetchCaseImages('c-1')).rejects.toThrow('Failed to fetch case images');
	});
});

// ---------------------------------------------------------------------------
// fetchDocumentWorkflow — URL audit
// ---------------------------------------------------------------------------

describe('fetchDocumentWorkflow', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('calls GET /cases/{caseId}/document-workflow', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: [] }),
		);

		await api.fetchDocumentWorkflow('case-123');

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/cases/case-123/document-workflow`,
		);
	});

	it('returns the parsed response', async () => {
		const expected = { success: true, data: [{ id: 'd-1', type: 'estimate', status: 'draft' }] };
		mockApiFetch.mockResolvedValueOnce(jsonResponse(expected));

		const result = await api.fetchDocumentWorkflow('case-123');

		expect(result).toEqual(expected);
	});

	it('throws on non-OK response', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({}, 500));

		await expect(api.fetchDocumentWorkflow('case-123')).rejects.toThrow(
			'Failed to fetch document workflow',
		);
	});
});

// ---------------------------------------------------------------------------
// fetchOverdueInvoices — URL audit
// ---------------------------------------------------------------------------

describe('fetchOverdueInvoices', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('calls GET /invoices/overdue with pagination params', async () => {
		const expected = {
			success: true,
			data: [],
			total: 0,
			page: 1,
			per_page: 20,
		};
		mockApiFetch.mockResolvedValueOnce(jsonResponse(expected));

		const result = await api.fetchOverdueInvoices(1, 20);

		expect(mockApiFetch).toHaveBeenCalledWith(
			'http://localhost:8080/invoices/overdue?page=1&per_page=20',
		);
		expect(result).toEqual(expected);
	});

	it('throws on non-OK response', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({}, 500));

		await expect(api.fetchOverdueInvoices()).rejects.toThrow(
			'Failed to fetch overdue invoices',
		);
	});
});

// ---------------------------------------------------------------------------
// Document endpoints — URL shape audit
// ---------------------------------------------------------------------------

describe('Document endpoint URLs', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	// fetchDocuments
	it('fetchDocuments calls GET /documents with query params', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: [], total: 0, page: 1, per_page: 20 }),
		);

		await api.fetchDocuments({ type: 'estimate', case_id: 'c-1', status: 'draft', page: 2, per_page: 10 });

		const calledUrl = mockApiFetch.mock.calls[0][0] as string;
		expect(calledUrl).toContain('/documents?');
		expect(calledUrl).toContain('type=estimate');
		expect(calledUrl).toContain('case_id=c-1');
		expect(calledUrl).toContain('status=draft');
		expect(calledUrl).toContain('page=2');
		expect(calledUrl).toContain('per_page=10');
	});

	// fetchDocument
	it('fetchDocument calls GET /documents/{id}/{type}', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: { id: 'd-1', type: 'estimate' } }),
		);

		await api.fetchDocument('d-1', 'estimate');

		expect(mockApiFetch).toHaveBeenCalledWith(`${BASE}/documents/d-1/estimate`);
	});

	// updateDocument
	it('updateDocument calls PATCH /documents/{id}/{type}', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: { id: 'd-1' } }),
		);

		await api.updateDocument('d-1', 'mandate', { data: { notes: 'updated' } });

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/documents/d-1/mandate`,
			expect.objectContaining({ method: 'PATCH' }),
		);
	});

	// deleteDocument
	it('deleteDocument calls DELETE /documents/{id}/{type}', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: {} }),
		);

		await api.deleteDocument('d-1', 'contract');

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/documents/d-1/contract`,
			expect.objectContaining({ method: 'DELETE' }),
		);
	});

	// sendDocument
	it('sendDocument calls POST /documents/{id}/{type}/send', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: {} }),
		);

		await api.sendDocument('d-1', 'invoice', {
			recipients: ['test@example.com'],
			send_email: true,
			send_sms: false,
		});

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/documents/d-1/invoice/send`,
			expect.objectContaining({ method: 'POST' }),
		);
	});

	// signDocument
	it('signDocument calls POST /documents/{id}/{type}/sign', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: {} }),
		);

		await api.signDocument('d-1', 'mandate', {
			signer_name: 'John',
			signer_role: 'client',
			method: 'e-sign',
		});

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/documents/d-1/mandate/sign`,
			expect.objectContaining({ method: 'POST' }),
		);
	});

	// archiveDocument
	it('archiveDocument calls POST /documents/{id}/{type}/archive', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: {} }),
		);

		await api.archiveDocument('d-1', 'estimate');

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/documents/d-1/estimate/archive`,
			expect.objectContaining({ method: 'POST' }),
		);
	});

	// fetchDocumentHistory
	it('fetchDocumentHistory calls GET /documents/{id}/{type}/history', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: [], total: 0, page: 1, per_page: 20 }),
		);

		await api.fetchDocumentHistory('d-1', 'contract');

		expect(mockApiFetch).toHaveBeenCalledWith(
			expect.stringContaining('/documents/d-1/contract/history'),
		);
	});
});

// ---------------------------------------------------------------------------
// Client endpoint URLs — URL audit (singular /client)
// ---------------------------------------------------------------------------

describe('Client endpoint URLs', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('fetchAllClients calls GET /client', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse([]));
		await api.fetchAllClients();
		expect(mockApiFetch).toHaveBeenCalledWith(`${BASE}/client`);
	});

	it('fetchClient calls GET /client/{id}', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({ id: 'c-1' }));
		await api.fetchClient('c-1');
		expect(mockApiFetch).toHaveBeenCalledWith(`${BASE}/client/c-1`);
	});

	it('createClient calls POST /client', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({ id: 'c-1' }));
		await api.createClient({ name: 'Test' });
		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/client`,
			expect.objectContaining({ method: 'POST' }),
		);
	});

	it('updateClient calls PATCH /client/{id}', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({ id: 'c-1' }));
		await api.updateClient('c-1', { name: 'Updated' });
		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/client/c-1`,
			expect.objectContaining({ method: 'PATCH' }),
		);
	});

	it('deleteClient calls DELETE /client/{id}', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse(null, 204));
		await api.deleteClient('c-1');
		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/client/c-1`,
			expect.objectContaining({ method: 'DELETE' }),
		);
	});

	it('fetchClientContacts calls GET /client/{id}/contact', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse([]));
		await api.fetchClientContacts('c-1');
		expect(mockApiFetch).toHaveBeenCalledWith(`${BASE}/client/c-1/contact`);
	});
});

// ---------------------------------------------------------------------------
// Contact endpoint URLs — URL audit (singular /contact, create nested)
// ---------------------------------------------------------------------------

describe('Contact endpoint URLs', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('fetchContact calls GET /contact/{id}', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({ id: 'ct-1' }));
		await api.fetchContact('ct-1');
		expect(mockApiFetch).toHaveBeenCalledWith(`${BASE}/contact/ct-1`);
	});

	it('createContact calls POST /client/{clientId}/contact', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({ id: 'ct-1' }));
		await api.createContact({ clientID: 'c-1', lastname: 'Doe', firstname: 'John' });
		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/client/c-1/contact`,
			expect.objectContaining({ method: 'POST' }),
		);
	});

	it('updateContact calls PATCH /contact/{id}', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({ id: 'ct-1' }));
		await api.updateContact('ct-1', { firstname: 'Jane' });
		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/contact/ct-1`,
			expect.objectContaining({ method: 'PATCH' }),
		);
	});

	it('deleteContact calls DELETE /contact/{id}', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse(null, 204));
		await api.deleteContact('ct-1');
		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/contact/ct-1`,
			expect.objectContaining({ method: 'DELETE' }),
		);
	});
});

// ---------------------------------------------------------------------------
// Case endpoint URLs — URL audit
// ---------------------------------------------------------------------------

describe('Case endpoint URLs', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('fetchAllCases calls GET /cases', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ cases: [], pagination: { page: 1, pageSize: 50, totalItems: 0, totalPages: 1 } }),
		);
		await api.fetchAllCases();
		const calledUrl = mockApiFetch.mock.calls[0][0] as string;
		expect(calledUrl).toContain('/cases');
		expect(calledUrl).not.toContain('/api/');
	});

	it('fetchCase calls GET /cases/{id}', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({ id: 'cs-1' }));
		await api.fetchCase('cs-1');
		expect(mockApiFetch).toHaveBeenCalledWith(`${BASE}/cases/cs-1`);
	});

	it('fetchCasesByClient calls GET /clients/{clientId}/cases', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ cases: [], pagination: { page: 1, pageSize: 50, totalItems: 0, totalPages: 1 } }),
		);
		await api.fetchCasesByClient('c-1');
		const calledUrl = mockApiFetch.mock.calls[0][0] as string;
		expect(calledUrl).toContain('/clients/c-1/cases');
		expect(calledUrl).not.toContain('/api/');
	});

	it('createCase calls POST /cases', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({ id: 'cs-1' }));
		await api.createCase({ title: 'Test', description: '', clientId: 'c-1', status: 'in_progress' });
		const [url, opts] = mockApiFetch.mock.calls[0];
		expect(url).toContain('/cases');
		expect(opts.method).toBe('POST');
	});

	it('markCaseReady calls POST /cases/{id}/mark-ready', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse(null, 204));
		await api.markCaseReady('cs-1');
		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/cases/cs-1/mark-ready`,
			expect.objectContaining({ method: 'POST' }),
		);
	});

	it('releaseCase calls POST /cases/{id}/release', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse({ caseId: 'cs-1', tokenId: 't-1', rawToken: 'tok', portalUrl: 'https://p.com', expiresAt: '2025-01-01' }));
		await api.releaseCase('cs-1');
		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/cases/cs-1/release`,
			expect.objectContaining({ method: 'POST' }),
		);
	});
});

// ---------------------------------------------------------------------------
// Mandate-specific routes — URL audit
// ---------------------------------------------------------------------------

describe('Mandate endpoint URLs', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('createContractFromMandate calls POST /mandates/{id}/create-contract', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: { id: 'c-1' } }),
		);

		await api.createContractFromMandate('m-1', { id: 'c-1' } as any);

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/mandates/m-1/create-contract`,
			expect.objectContaining({ method: 'POST' }),
		);
	});
});

// ---------------------------------------------------------------------------
// Contract-specific routes — URL audit
// ---------------------------------------------------------------------------

describe('Contract endpoint URLs', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('createInvoiceFromContract calls POST /contracts/{id}/create-invoice', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: { id: 'inv-1' } }),
		);

		await api.createInvoiceFromContract('c-1', { id: 'inv-1' } as any);

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/contracts/c-1/create-invoice`,
			expect.objectContaining({ method: 'POST' }),
		);
	});
});

// ---------------------------------------------------------------------------
// Invoice-specific routes — URL audit
// ---------------------------------------------------------------------------

describe('Invoice endpoint URLs', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('processPayment calls POST /invoices/{id}/pay', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: {} }),
		);

		await api.processPayment('inv-1', { amount: 100, payment_method: 'card' });

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/invoices/inv-1/pay`,
			expect.objectContaining({ method: 'POST' }),
		);
	});

	it('voidInvoice calls POST /invoices/{id}/void', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ success: true, data: {} }),
		);

		await api.voidInvoice('inv-1');

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/invoices/inv-1/void`,
			expect.objectContaining({ method: 'POST' }),
		);
	});
});

// ---------------------------------------------------------------------------
// AI text endpoint — URL audit
// ---------------------------------------------------------------------------

describe('AI text endpoint URL', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('requestAITextOperation calls POST /ai/text/transform', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ result: 'transformed text' }),
		);

		await api.requestAITextOperation({
			text: 'hello',
			operation: 'reword',
			language: 'fr',
		});

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/ai/text/transform`,
			expect.objectContaining({ method: 'POST' }),
		);
	});
});

// ---------------------------------------------------------------------------
// AI chat endpoints — URL audit (these keep /api/ prefix)
// ---------------------------------------------------------------------------

describe('AI chat endpoint URLs', () => {
	beforeEach(() => {
		mockApiFetch.mockReset();
	});

	it('sendChatMessage calls POST /api/ai/chat/message', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ conversationId: 'conv-1', userMessageId: 'msg-1', assistantMessage: {} }),
		);

		await api.sendChatMessage('cs-1', { message: 'hello' });

		const calledUrl = mockApiFetch.mock.calls[0][0] as string;
		expect(calledUrl).toContain('/api/ai/chat/message');
		expect(calledUrl).toContain('case_id=cs-1');
	});

	it('getChatConversation calls GET /api/ai/chat/conversations/{id}', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ conversation: {}, messages: [] }),
		);

		await api.getChatConversation('conv-1');

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/api/ai/chat/conversations/conv-1`,
		);
	});

	it('listChatConversations calls GET /api/ai/chat/conversations', async () => {
		mockApiFetch.mockResolvedValueOnce(
			jsonResponse({ conversations: [] }),
		);

		await api.listChatConversations('cs-1');

		const calledUrl = mockApiFetch.mock.calls[0][0] as string;
		expect(calledUrl).toContain('/api/ai/chat/conversations');
		expect(calledUrl).toContain('case_id=cs-1');
	});

	it('deleteChatConversation calls DELETE /api/ai/chat/conversations/{id}', async () => {
		mockApiFetch.mockResolvedValueOnce(jsonResponse(null, 204));

		await api.deleteChatConversation('conv-1');

		expect(mockApiFetch).toHaveBeenCalledWith(
			`${BASE}/api/ai/chat/conversations/conv-1`,
			expect.objectContaining({ method: 'DELETE' }),
		);
	});
});

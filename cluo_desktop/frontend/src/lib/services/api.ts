/**
 * API Service
 * 
 * When VITE_USE_MOCK_DATA is true, returns mock data for development.
 * When false, makes actual API calls to the backend.
 */

import { isMockEnabled, API_BASE_URL } from '../config';
import { apiFetch } from './apiFetch';
import type {
	AuthUser,
	Case,
	CaseSubject,
	CaseType,
	Client,
	Contact,
	Contract,
	CreateCaseRequest,
	CreateDocumentRequest,
	DocumentAPIResponse,
	DocumentHistoryResponse,
	DocumentListResponse,
	DocumentSummary,
	DocumentWorkflowResponse,
	Estimate,
	Invoice,
	ListCasesResponse,
	Mandate,
	OverdueInvoicesResponse,
	PaymentRequest,
	ReleaseResponse,
	SendDocumentRequest,
	SignDocumentRequest,
	UpdateDocumentRequest,
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
export async function fetchAllUsers(): Promise<AuthUser[]> {
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
export async function fetchUser(id: string): Promise<AuthUser | null> {
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
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/clients`);
	if (!response.ok) {
		throw new Error(`Failed to fetch clients: ${response.status}`);
	}
	return response.json();
}

/**
 * Fetch a client by ID
 */
export async function fetchClient(id: string): Promise<Client | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getClientById(id) || null;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/clients/${id}`);
	if (!response.ok) {
		if (response.status === 404) return null;
		throw new Error(`Failed to fetch client: ${response.status}`);
	}
	return response.json();
}

/**
 * Create a new client
 */
export async function createClient(request: {
	name: string;
	type?: string;
}): Promise<Client> {
	if (isMockEnabled()) {
		await mockDelay();
		const newClient: Client = {
			id: `mock-${Date.now()}`,
			name: request.name,
			type: (request.type || 'company') as any,
		};
		return newClient;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/clients`, {
		method: 'POST',
		body: JSON.stringify(request),
	});
	if (!response.ok) {
		throw new Error(`Failed to create client: ${response.status}`);
	}
	return response.json();
}

/**
 * Update an existing client
 */
export async function updateClient(id: string, request: {
	name?: string;
	type?: string;
}): Promise<Client> {
	if (isMockEnabled()) {
		await mockDelay();
		const existing = mockData.getClientById(id);
		if (!existing) throw new Error(`Client ${id} not found`);
		const updated: Client = {
			id: existing.id,
			name: request.name ?? existing.name,
			type: (request.type ?? existing.type) as any,
		};
		return updated;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/clients/${id}`, {
		method: 'PATCH',
		body: JSON.stringify(request),
	});
	if (!response.ok) {
		throw new Error(`Failed to update client: ${response.status}`);
	}
	return response.json();
}

/**
 * Delete a client
 */
export async function deleteClient(id: string): Promise<void> {
	if (isMockEnabled()) {
		await mockDelay();
		return;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/clients/${id}`, {
		method: 'DELETE',
	});
	if (!response.ok) {
		throw new Error(`Failed to delete client: ${response.status}`);
	}
}

/**
 * Fetch contacts for a specific client
 */
export async function fetchClientContacts(clientId: string): Promise<Contact[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getContactsByClientId(clientId);
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/clients/${clientId}/contacts`);
	if (!response.ok) {
		throw new Error(`Failed to fetch client contacts: ${response.status}`);
	}
	return response.json();
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
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/contacts/${id}`);
	if (!response.ok) {
		if (response.status === 404) return null;
		throw new Error(`Failed to fetch contact: ${response.status}`);
	}
	return response.json();
}

/**
 * Create a new contact
 */
export async function createContact(request: {
	clientID: string;
	lastname: string;
	firstname: string;
	email?: string;
	phone?: string;
	position?: string;
}): Promise<Contact> {
	if (isMockEnabled()) {
		await mockDelay();
		const newContact: Contact = {
			id: `mock-${Date.now()}`,
			clientID: request.clientID,
			lastname: request.lastname,
			firstname: request.firstname,
			email: request.email ?? '',
			phone: request.phone ?? '',
			position: request.position ?? '',
			createdAt: new Date().toISOString(),
		};
		return newContact;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/contacts`, {
		method: 'POST',
		body: JSON.stringify(request),
	});
	if (!response.ok) {
		throw new Error(`Failed to create contact: ${response.status}`);
	}
	return response.json();
}

/**
 * Update an existing contact
 */
export async function updateContact(id: string, request: {
	lastname?: string;
	firstname?: string;
	email?: string;
	phone?: string;
	position?: string;
}): Promise<Contact> {
	if (isMockEnabled()) {
		await mockDelay();
		const existing = mockData.getContactById(id);
		if (!existing) throw new Error(`Contact ${id} not found`);
		const updated = { ...existing, ...request };
		return updated;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/contacts/${id}`, {
		method: 'PATCH',
		body: JSON.stringify(request),
	});
	if (!response.ok) {
		throw new Error(`Failed to update contact: ${response.status}`);
	}
	return response.json();
}

/**
 * Delete a contact
 */
export async function deleteContact(id: string): Promise<void> {
	if (isMockEnabled()) {
		await mockDelay();
		return;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/contacts/${id}`, {
		method: 'DELETE',
	});
	if (!response.ok) {
		throw new Error(`Failed to delete contact: ${response.status}`);
	}
}

// =============================================================================
// CASES
// =============================================================================

interface FetchCasesParams {
	page?: number;
	pageSize?: number;
	status?: string;
}

/**
 * Fetch all cases with optional pagination and filters
 */
export async function fetchAllCases(params?: FetchCasesParams): Promise<ListCasesResponse> {
		if (isMockEnabled()) {
			await mockDelay();
			return {
				cases: mockData.getAllCases(),
				pagination: { page: 1, pageSize: 50, totalItems: mockData.getAllCases().length, totalPages: 1 }
			};
		}

		const baseURL = API_BASE_URL;
		const url = new URL(`${baseURL}/api/cases`);

		if (params?.page) url.searchParams.set('page', params.page.toString());
		if (params?.pageSize) url.searchParams.set('page_size', params.pageSize.toString());
		if (params?.status) url.searchParams.set('status', params.status);

		const response = await apiFetch(url.toString());
		if (!response.ok) {
			throw new Error(`Failed to fetch cases: ${response.status}`);
		}

		return response.json();
	}

	/**
	 * Fetch a case by ID with full details
	 */
	export async function fetchCase(id: string): Promise<Case> {
		if (isMockEnabled()) {
			await mockDelay();
			const caseData = mockData.getCaseById(id);
			if (!caseData) throw new Error(`Case ${id} not found`);
			return caseData;
		}

		const baseURL = API_BASE_URL;
		const response = await apiFetch(`${baseURL}/api/cases/${id}`);
		if (!response.ok) {
			throw new Error(`Failed to fetch case: ${response.status}`);
		}

		return response.json();
	}

	/**
	 * Fetch cases by status
	 */
	export async function fetchCasesByStatus(status: string): Promise<ListCasesResponse> {
		return fetchAllCases({ status });
	}

	/**
	 * Fetch cases by client ID
	 */
	export async function fetchCasesByClient(clientId: string, params?: Omit<FetchCasesParams, 'status'>): Promise<ListCasesResponse> {
		if (isMockEnabled()) {
			await mockDelay();
			return {
				cases: mockData.getCasesByClientId(clientId),
				pagination: { page: 1, pageSize: 50, totalItems: mockData.getCasesByClientId(clientId).length, totalPages: 1 }
			};
		}

		const baseURL = API_BASE_URL;
		const url = new URL(`${baseURL}/api/clients/${clientId}/cases`);

		if (params?.page) url.searchParams.set('page', params.page.toString());
		if (params?.pageSize) url.searchParams.set('page_size', params.pageSize.toString());

		const response = await apiFetch(url.toString());
		if (!response.ok) {
			throw new Error(`Failed to fetch cases for client: ${response.status}`);
		}

		return response.json();
	}

	/**
	 * Create a new case
	 */
	export async function createCase(request: CreateCaseRequest): Promise<Case> {
		if (isMockEnabled()) {
			await mockDelay();
			const newCase: Case = {
				id: `mock-${Date.now()}`,
				title: request.title,
				description: request.description,
				clientId: request.clientId,
				assignedContactID: request.assignedContactID ?? null,
				caseSubjectId: request.caseSubjectId ?? null,
				externalReference: request.externalReference ?? null,
				caseTypeId: request.caseTypeId ?? null,
				status: request.status,
				placename: request.placename ?? null,
				address1: request.address1 ?? null,
				address2: request.address2 ?? null,
				city: request.city ?? null,
				postalCode: request.postalCode ?? null,
				country: request.country ?? null,
				latitude: request.latitude ?? null,
				longitude: request.longitude ?? null,
				locationType: request.locationType ?? null,
				locationNotes: request.locationNotes ?? null,
				createdAt: new Date().toISOString(),
				updatedAt: new Date().toISOString(),
			};
			return newCase;
		}

		const baseURL = API_BASE_URL;
		const response = await apiFetch(`${baseURL}/api/cases`, {
			method: 'POST',
			body: JSON.stringify(request),
		});

		if (!response.ok) {
			throw new Error(`Failed to create case: ${response.status}`);
		}

		return response.json();
	}

	/**
	 * Update an existing case
	 */
	export async function updateCase(id: string, request: Partial<CreateCaseRequest>): Promise<Case> {
		if (isMockEnabled()) {
			await mockDelay();
			const existing = mockData.getCaseById(id);
			if (!existing) throw new Error(`Case ${id} not found`);
			const updated = { ...existing, ...request, updatedAt: new Date().toISOString() };
			return updated;
		}

		const baseURL = API_BASE_URL;
		const response = await apiFetch(`${baseURL}/api/cases/${id}`, {
			method: 'PATCH',
			body: JSON.stringify(request),
		});

		if (!response.ok) {
			throw new Error(`Failed to update case: ${response.status}`);
		}

		return response.json();
	}

	/**
	 * Delete a case
	 */
	export async function deleteCase(id: string): Promise<void> {
		if (isMockEnabled()) {
			await mockDelay();
			return;
		}

		const baseURL = API_BASE_URL;
		const response = await apiFetch(`${baseURL}/api/cases/${id}`, {
			method: 'DELETE',
		});

		if (!response.ok) {
			throw new Error(`Failed to delete case: ${response.status}`);
		}
	}

	/**
	 * Mark a case as ready
	 */
	export async function markCaseReady(id: string): Promise<void> {
		if (isMockEnabled()) {
			await mockDelay();
			return;
		}

		const baseURL = API_BASE_URL;
		const response = await apiFetch(`${baseURL}/api/cases/${id}/mark-ready`, {
			method: 'POST',
		});

		if (!response.ok) {
			throw new Error(`Failed to mark case as ready: ${response.status}`);
		}
	}

	/**
	 * Release a case and generate portal access token
	 */
	export async function releaseCase(id: string): Promise<ReleaseResponse> {
		if (isMockEnabled()) {
			await mockDelay();
			return {
				caseId: id,
				tokenId: `mock-token-${Date.now()}`,
				rawToken: 'mock-raw-token',
				portalUrl: 'https://portal.example.com',
				expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
			};
		}

		const baseURL = API_BASE_URL;
		const response = await apiFetch(`${baseURL}/api/cases/${id}/release`, {
			method: 'POST',
		});

		if (!response.ok) {
			throw new Error(`Failed to release case: ${response.status}`);
		}

		return response.json();
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
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/subjects`);
	if (!response.ok) {
		throw new Error(`Failed to fetch case subjects: ${response.status}`);
	}
	return response.json();
}

/**
 * Fetch a case subject by ID
 */
export async function fetchCaseSubject(id: string): Promise<CaseSubject | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getCaseSubjectById(id) || null;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/subjects/${id}`);
	if (!response.ok) {
		if (response.status === 404) return null;
		throw new Error(`Failed to fetch case subject: ${response.status}`);
	}
	return response.json();
}

/**
 * Create a new case subject
 */
export async function createCaseSubject(request: {
	lastname: string;
	firstname: string;
	email?: string;
	phone?: string;
	city?: string;
	postalCode?: string;
	address1?: string;
	address2?: string;
	occupation?: string;
	notes?: string;
}): Promise<CaseSubject> {
	if (isMockEnabled()) {
		await mockDelay();
		const newSubject: CaseSubject = {
			id: `mock-${Date.now()}`,
			firstname: request.firstname,
			lastname: request.lastname,
			email: request.email ?? '',
			phone: request.phone ?? '',
			address1: request.address1 ?? '',
			address2: request.address2 ?? '',
			city: request.city ?? '',
			postalCode: request.postalCode ?? '',
			occupation: request.occupation ?? '',
			notes: request.notes ?? '',
			createdAt: new Date().toISOString(),
		};
		return newSubject;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/subjects`, {
		method: 'POST',
		body: JSON.stringify(request),
	});
	if (!response.ok) {
		throw new Error(`Failed to create case subject: ${response.status}`);
	}
	return response.json();
}

/**
 * Update an existing case subject
 */
export async function updateCaseSubject(id: string, request: {
	lastname?: string;
	firstname?: string;
	email?: string;
	phone?: string;
	city?: string;
	postalCode?: string;
	address1?: string;
	address2?: string;
	occupation?: string;
	notes?: string;
}): Promise<CaseSubject> {
	if (isMockEnabled()) {
		await mockDelay();
		const existing = mockData.getCaseSubjectById(id);
		if (!existing) throw new Error(`Case subject ${id} not found`);
		const updated = { ...existing, ...request };
		return updated;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/subjects/${id}`, {
		method: 'PATCH',
		body: JSON.stringify(request),
	});
	if (!response.ok) {
		throw new Error(`Failed to update case subject: ${response.status}`);
	}
	return response.json();
}

/**
 * Delete a case subject
 */
export async function deleteCaseSubject(id: string): Promise<void> {
	if (isMockEnabled()) {
		await mockDelay();
		return;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/subjects/${id}`, {
		method: 'DELETE',
	});
	if (!response.ok) {
		throw new Error(`Failed to delete case subject: ${response.status}`);
	}
}

// =============================================================================
// CASE TYPES
// =============================================================================

/**
 * Fetch all case types (global list)
 */
export async function fetchAllCaseTypes(): Promise<CaseType[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getAllCaseTypes() as unknown as CaseType[];
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/case-types`);
	if (!response.ok) {
		throw new Error(`Failed to fetch case types: ${response.status}`);
	}
	return response.json();
}

/**
 * Fetch a case type by ID
 */
export async function fetchCaseType(id: string): Promise<CaseType | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getCaseTypeById(id) as unknown as CaseType || null;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/case-types/${id}`);
	if (!response.ok) {
		if (response.status === 404) return null;
		throw new Error(`Failed to fetch case type: ${response.status}`);
	}
	return response.json();
}

/**
 * Create a new case type
 */
export async function createCaseType(request: { name: string }): Promise<CaseType> {
	if (isMockEnabled()) {
		await mockDelay();
		const newType: CaseType = {
			id: `mock-${Date.now()}`,
			name: request.name,
			createdAt: new Date().toISOString(),
			updatedAt: new Date().toISOString(),
		};
		return newType;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/case-types`, {
		method: 'POST',
		body: JSON.stringify(request),
	});
	if (!response.ok) {
		throw new Error(`Failed to create case type: ${response.status}`);
	}
	return response.json();
}

/**
 * Update an existing case type
 */
export async function updateCaseType(id: string, request: { name: string }): Promise<CaseType> {
	if (isMockEnabled()) {
		await mockDelay();
		const existing = mockData.getCaseTypeById(id) as unknown as CaseType | undefined;
		if (!existing) throw new Error(`Case type ${id} not found`);
		const updated = { ...existing, name: request.name, updatedAt: new Date().toISOString() };
		return updated;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/case-types/${id}`, {
		method: 'PATCH',
		body: JSON.stringify(request),
	});
	if (!response.ok) {
		throw new Error(`Failed to update case type: ${response.status}`);
	}
	return response.json();
}

/**
 * Delete a case type
 */
export async function deleteCaseType(id: string): Promise<void> {
	if (isMockEnabled()) {
		await mockDelay();
		return;
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/case-types/${id}`, {
		method: 'DELETE',
	});
	if (!response.ok) {
		throw new Error(`Failed to delete case type: ${response.status}`);
	}
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
		return mockData.getAllEstimates() as unknown as Estimate[];
	}
	const result = await fetchDocuments({ type: 'estimate' });
	return result.data as unknown as Estimate[];
}

/**
 * Fetch estimates for a specific case
 */
export async function fetchCaseEstimates(caseId: string): Promise<Estimate[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getEstimatesByCaseId(caseId) as unknown as Estimate[];
	}
	const result = await fetchDocuments({ type: 'estimate', case_id: caseId });
	return result.data as unknown as Estimate[];
}

/**
 * Fetch an estimate by ID
 */
export async function fetchEstimate(id: string): Promise<Estimate | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getEstimateById(id) as unknown as Estimate || null;
	}
	const result = await fetchDocument(id, 'estimate');
	return result.data as Estimate || null;
}

/**
 * Fetch estimates by client ID
 */
export async function fetchClientEstimates(clientId: string): Promise<Estimate[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getEstimatesByClientId(clientId) as unknown as Estimate[];
	}
	// Use the generic fetchDocuments with client filter
	const result = await fetchDocuments({ type: 'estimate' });
	// Filter by client_id in the response
	return (result.data as unknown as Estimate[]).filter(est => (est as Estimate).client_id === clientId);
}

/**
 * Create a new estimate
 */
export async function createEstimate(estimate: Estimate): Promise<DocumentAPIResponse<Estimate>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: estimate };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/estimates`, {
		method: 'POST',
		body: JSON.stringify(estimate),
	});

	if (!response.ok) {
		throw new Error(`Failed to create estimate: ${response.status}`);
	}

	return response.json();
}

/**
 * Update an estimate (line items)
 */
export async function updateEstimate(id: string, lineItems: any[]): Promise<DocumentAPIResponse<Estimate>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { id, line_items: lineItems } as any };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/estimates/${id}`, {
		method: 'PATCH',
		body: JSON.stringify({ line_items: lineItems }),
	});

	if (!response.ok) {
		throw new Error(`Failed to update estimate: ${response.status}`);
	}

	return response.json();
}

/**
 * Accept an estimate
 */
export async function acceptEstimate(id: string, acceptedBy: string): Promise<DocumentAPIResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { message: 'Estimate accepted successfully' } };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/estimates/${id}/accept`, {
		method: 'POST',
		body: JSON.stringify({ accepted_by: acceptedBy }),
	});

	if (!response.ok) {
		throw new Error(`Failed to accept estimate: ${response.status}`);
	}

	return response.json();
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
		return mockData.getAllMandates() as unknown as Mandate[];
	}
	const result = await fetchDocuments({ type: 'mandate' });
	return result.data as unknown as Mandate[];
}

/**
 * Fetch mandates for a specific case
 */
export async function fetchCaseMandates(caseId: string): Promise<Mandate[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getMandatesByCaseId(caseId) as unknown as Mandate[];
	}
	const result = await fetchDocuments({ type: 'mandate', case_id: caseId });
	return result.data as unknown as Mandate[];
}

/**
 * Fetch a mandate by ID
 */
export async function fetchMandate(id: string): Promise<Mandate | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getMandateById(id) as unknown as Mandate || null;
	}
	const result = await fetchDocument(id, 'mandate');
	return result.data as Mandate || null;
}

/**
 * Fetch mandates by client ID
 */
export async function fetchClientMandates(clientId: string): Promise<Mandate[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getMandatesByClientId(clientId) as unknown as Mandate[];
	}
	const result = await fetchDocuments({ type: 'mandate' });
	return (result.data as unknown as Mandate[]).filter(mand => (mand as Mandate).client_id === clientId);
}

/**
 * Create a new mandate
 */
export async function createMandate(mandate: Mandate): Promise<DocumentAPIResponse<Mandate>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: mandate };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/mandates`, {
		method: 'POST',
		body: JSON.stringify(mandate),
	});

	if (!response.ok) {
		throw new Error(`Failed to create mandate: ${response.status}`);
	}

	return response.json();
}

/**
 * Sign a mandate
 */
export async function signMandate(id: string, request: SignDocumentRequest): Promise<DocumentAPIResponse<Mandate>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { id } as any };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/mandates/${id}/sign`, {
		method: 'POST',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		throw new Error(`Failed to sign mandate: ${response.status}`);
	}

	return response.json();
}

/**
 * Activate a mandate
 */
export async function activateMandate(id: string): Promise<DocumentAPIResponse<Mandate>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { id } as any };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/mandates/${id}/activate`, {
		method: 'POST',
	});

	if (!response.ok) {
		throw new Error(`Failed to activate mandate: ${response.status}`);
	}

	return response.json();
}

/**
 * Create a contract from a mandate
 */
export async function createContractFromMandate(mandateId: string, contract: Contract): Promise<DocumentAPIResponse<Contract>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: contract };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/mandates/${mandateId}/create-contract`, {
		method: 'POST',
		body: JSON.stringify(contract),
	});

	if (!response.ok) {
		throw new Error(`Failed to create contract from mandate: ${response.status}`);
	}

	return response.json();
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
		return mockData.getAllContracts() as unknown as Contract[];
	}
	const result = await fetchDocuments({ type: 'contract' });
	return result.data as unknown as Contract[];
}

/**
 * Fetch contracts for a specific case
 */
export async function fetchCaseContracts(caseId: string): Promise<Contract[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getContractsByCaseId(caseId) as unknown as Contract[];
	}
	const result = await fetchDocuments({ type: 'contract', case_id: caseId });
	return result.data as unknown as Contract[];
}

/**
 * Fetch a contract by ID
 */
export async function fetchContract(id: string): Promise<Contract | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getContractById(id) as unknown as Contract || null;
	}
	const result = await fetchDocument(id, 'contract');
	return result.data as Contract || null;
}

/**
 * Fetch contracts by client ID
 */
export async function fetchClientContracts(clientId: string): Promise<Contract[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getContractsByClientId(clientId) as unknown as Contract[];
	}
	const result = await fetchDocuments({ type: 'contract' });
	return (result.data as unknown as Contract[]).filter(cont => (cont as Contract).client_id === clientId);
}

/**
 * Create a new contract
 */
export async function createContract(contract: Contract): Promise<DocumentAPIResponse<Contract>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: contract };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/contracts`, {
		method: 'POST',
		body: JSON.stringify(contract),
	});

	if (!response.ok) {
		throw new Error(`Failed to create contract: ${response.status}`);
	}

	return response.json();
}

/**
 * Sign a contract
 */
export async function signContract(id: string, request: SignDocumentRequest): Promise<DocumentAPIResponse<Contract>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { id } as any };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/contracts/${id}/sign`, {
		method: 'POST',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		throw new Error(`Failed to sign contract: ${response.status}`);
	}

	return response.json();
}

/**
 * Activate a contract
 */
export async function activateContract(id: string): Promise<DocumentAPIResponse<Contract>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { id } as any };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/contracts/${id}/activate`, {
		method: 'POST',
	});

	if (!response.ok) {
		throw new Error(`Failed to activate contract: ${response.status}`);
	}

	return response.json();
}

/**
 * Create an invoice from a contract
 */
export async function createInvoiceFromContract(contractId: string, invoice: Invoice): Promise<DocumentAPIResponse<Invoice>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: invoice };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/contracts/${contractId}/create-invoice`, {
		method: 'POST',
		body: JSON.stringify(invoice),
	});

	if (!response.ok) {
		throw new Error(`Failed to create invoice from contract: ${response.status}`);
	}

	return response.json();
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
		return mockData.getAllInvoices() as unknown as Invoice[];
	}
	const result = await fetchDocuments({ type: 'invoice' });
	return result.data as unknown as Invoice[];
}

/**
 * Fetch invoices for a specific case
 */
export async function fetchCaseInvoices(caseId: string): Promise<Invoice[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getInvoicesByCaseId(caseId) as unknown as Invoice[];
	}
	const result = await fetchDocuments({ type: 'invoice', case_id: caseId });
	return result.data as unknown as Invoice[];
}

/**
 * Fetch an invoice by ID
 */
export async function fetchInvoice(id: string): Promise<Invoice | null> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getInvoiceById(id) as unknown as Invoice || null;
	}
	const result = await fetchDocument(id, 'invoice');
	return result.data as Invoice || null;
}

/**
 * Fetch invoices by client ID
 */
export async function fetchClientInvoices(clientId: string): Promise<Invoice[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getInvoicesByClientId(clientId) as unknown as Invoice[];
	}
	const result = await fetchDocuments({ type: 'invoice' });
	return (result.data as unknown as Invoice[]).filter(inv => (inv as Invoice).client_id === clientId);
}

/**
 * Fetch invoices by payment status
 */
export async function fetchInvoicesByPaymentStatus(paymentStatus: string): Promise<Invoice[]> {
	if (isMockEnabled()) {
		await mockDelay();
		return mockData.getInvoicesByPaymentStatus(paymentStatus as any) as unknown as Invoice[];
	}
	const result = await fetchDocuments({ type: 'invoice' });
	return (result.data as unknown as Invoice[]).filter(inv => (inv as Invoice).payment_status === paymentStatus);
}

/**
 * Create a new invoice
 */
export async function createInvoice(invoice: Invoice): Promise<DocumentAPIResponse<Invoice>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: invoice };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/invoices`, {
		method: 'POST',
		body: JSON.stringify(invoice),
	});

	if (!response.ok) {
		throw new Error(`Failed to create invoice: ${response.status}`);
	}

	return response.json();
}

/**
 * Fetch overdue invoices
 */
export async function fetchOverdueInvoices(page: number = 1, perPage: number = 20): Promise<OverdueInvoicesResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		return {
			success: true,
			data: [],
			total: 0,
			page,
			per_page: perPage
		};
	}
	const baseURL = API_BASE_URL;
	const url = new URL(`${baseURL}/api/invoices/overdue`);
	url.searchParams.set('page', page.toString());
	url.searchParams.set('per_page', perPage.toString());

	const response = await apiFetch(url.toString());
	if (!response.ok) {
		throw new Error(`Failed to fetch overdue invoices: ${response.status}`);
	}

	return response.json();
}

/**
 * Process a payment on an invoice
 */
export async function processPayment(id: string, request: PaymentRequest): Promise<DocumentAPIResponse<Invoice>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { id } as any };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/invoices/${id}/pay`, {
		method: 'POST',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		throw new Error(`Failed to process payment: ${response.status}`);
	}

	return response.json();
}

/**
 * Void an invoice
 */
export async function voidInvoice(id: string): Promise<DocumentAPIResponse<Invoice>> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { id } as any };
	}
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/invoices/${id}/void`, {
		method: 'POST',
	});

	if (!response.ok) {
		throw new Error(`Failed to void invoice: ${response.status}`);
	}

	return response.json();
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
	const baseURL = API_BASE_URL;
	const timeout = AI_CONFIG.DEFAULT_TIMEOUT;

	const controller = new AbortController();
	const timeoutId = setTimeout(() => controller.abort(), timeout);

	try {
		const response = await apiFetch(`${baseURL}/api/ai/text`, {
			method: "POST",
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
	const baseURL = API_BASE_URL;
	const url = new URL(`${baseURL}/api/ai/chat/message`);
	url.searchParams.set('case_id', caseId);

	const response = await apiFetch(url.toString(), {
		method: 'POST',
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
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/ai/chat/conversations/${conversationId}`);

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
	const baseURL = API_BASE_URL;
	const url = new URL(`${baseURL}/api/ai/chat/conversations`);
	url.searchParams.set('case_id', caseId);

	const response = await apiFetch(url.toString());

	if (!response.ok) {
		throw new Error(`Failed to list conversations: ${response.status}`);
	}

	return response.json();
}

/**
 * Delete a conversation
 */
export async function deleteChatConversation(conversationId: string): Promise<void> {
	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/ai/chat/conversations/${conversationId}`, {
		method: 'DELETE',
	});

	if (!response.ok) {
		throw new Error(`Failed to delete conversation: ${response.status}`);
	}
}

// =============================================================================
// DOCUMENTS (Generic)
// =============================================================================

interface FetchDocumentsParams {
	case_id?: string;
	type?: string;
	status?: string;
	page?: number;
	per_page?: number;
}

/**
 * Fetch documents with optional filters and pagination
 */
export async function fetchDocuments(params?: FetchDocumentsParams): Promise<DocumentListResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		// Return mock data for now
		return {
			success: true,
			data: [],
			total: 0,
			page: 1,
			per_page: 20
		};
	}

	const baseURL = API_BASE_URL;
	const url = new URL(`${baseURL}/api/documents`);

	if (params?.case_id) url.searchParams.set('case_id', params.case_id);
	if (params?.type) url.searchParams.set('type', params.type);
	if (params?.status) url.searchParams.set('status', params.status);
	if (params?.page) url.searchParams.set('page', params.page.toString());
	if (params?.per_page) url.searchParams.set('per_page', params.per_page.toString());

	const response = await apiFetch(url.toString());
	if (!response.ok) {
		throw new Error(`Failed to fetch documents: ${response.status}`);
	}

	return response.json();
}

/**
 * Create a new document
 */
export async function createDocument(request: CreateDocumentRequest): Promise<DocumentAPIResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { id: `mock-${Date.now()}` } };
	}

	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/documents`, {
		method: 'POST',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		throw new Error(`Failed to create document: ${response.status}`);
	}

	return response.json();
}

/**
 * Fetch a single document by ID and type
 */
export async function fetchDocument(id: string, type: string): Promise<DocumentAPIResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { id, type } };
	}

	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/documents/${id}/${type}`);
	if (!response.ok) {
		throw new Error(`Failed to fetch document: ${response.status}`);
	}

	return response.json();
}

/**
 * Update an existing document
 */
export async function updateDocument(id: string, type: string, request: UpdateDocumentRequest): Promise<DocumentAPIResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { id, type } };
	}

	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/documents/${id}/${type}`, {
		method: 'PATCH',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		throw new Error(`Failed to update document: ${response.status}`);
	}

	return response.json();
}

/**
 * Delete a document
 */
export async function deleteDocument(id: string, type: string): Promise<DocumentAPIResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { message: 'Document deleted successfully' } };
	}

	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/documents/${id}/${type}`, {
		method: 'DELETE',
	});

	if (!response.ok) {
		throw new Error(`Failed to delete document: ${response.status}`);
	}

	return response.json();
}

/**
 * Send a document to recipients
 */
export async function sendDocument(id: string, type: string, request: SendDocumentRequest): Promise<DocumentAPIResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { message: 'Document sent successfully' } };
	}

	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/documents/${id}/${type}/send`, {
		method: 'POST',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		throw new Error(`Failed to send document: ${response.status}`);
	}

	return response.json();
}

/**
 * Sign a document
 */
export async function signDocument(id: string, type: string, request: SignDocumentRequest): Promise<DocumentAPIResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { id, type } };
	}

	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/documents/${id}/${type}/sign`, {
		method: 'POST',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		throw new Error(`Failed to sign document: ${response.status}`);
	}

	return response.json();
}

/**
 * Archive a document
 */
export async function archiveDocument(id: string, type: string): Promise<DocumentAPIResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: { message: 'Document archived successfully' } };
	}

	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/documents/${id}/${type}/archive`, {
		method: 'POST',
	});

	if (!response.ok) {
		throw new Error(`Failed to archive document: ${response.status}`);
	}

	return response.json();
}

/**
 * Fetch document history (versions)
 */
export async function fetchDocumentHistory(id: string, type: string, page: number = 1, perPage: number = 20): Promise<DocumentHistoryResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		return {
			success: true,
			data: [],
			total: 0,
			page,
			per_page: perPage
		};
	}

	const baseURL = API_BASE_URL;
	const url = new URL(`${baseURL}/api/documents/${id}/${type}/history`);
	url.searchParams.set('page', page.toString());
	url.searchParams.set('per_page', perPage.toString());

	const response = await apiFetch(url.toString());
	if (!response.ok) {
		throw new Error(`Failed to fetch document history: ${response.status}`);
	}

	return response.json();
}

/**
 * Fetch document workflow for a case (full document chain)
 */
export async function fetchDocumentWorkflow(caseId: string): Promise<DocumentWorkflowResponse> {
	if (isMockEnabled()) {
		await mockDelay();
		return { success: true, data: [] };
	}

	const baseURL = API_BASE_URL;
	const response = await apiFetch(`${baseURL}/api/documents/workflow/${caseId}`);
	if (!response.ok) {
		throw new Error(`Failed to fetch document workflow: ${response.status}`);
	}

	return response.json();
}


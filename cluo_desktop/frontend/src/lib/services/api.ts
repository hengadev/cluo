/**
 * API Service
 *
 * All calls go to the real backend API, except when VITE_MOCK_USER_ROLE is set
 * — in that case every exported function is handled by mockApi.ts instead.
 */

import { API_BASE_URL } from '../config';
import { apiFetch } from './apiFetch';
import * as mock from './mockApi';

const MOCK = import.meta.env.VITE_MOCK_USER_ROLE as string | undefined;
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
	EstimateItem,
	Invoice,
	ListCasesResponse,
	Mandate,
	OverdueInvoicesResponse,
	PaymentRequest,
	ReleaseResponse,
	SearchResult,
	SendDocumentRequest,
	SignDocumentRequest,
	UpdateDocumentRequest,
} from '../types/entities';

const BASE_URL = API_BASE_URL;

/**
 * Error thrown when the API returns 409 Conflict.
 * Used for lifecycle violations (e.g. activating an unsigned mandate).
 */
export class ConflictError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'ConflictError';
	}
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// =============================================================================
// USERS
// =============================================================================

/**
 * Fetch all users.
 * The API is single-PI so there is no "list users" endpoint.
 * We call GET /auth/me and return a one-element array.
 */
export async function fetchAllUsers(): Promise<AuthUser[]> {
	if (MOCK) return mock.fetchAllUsers();
	const response = await apiFetch(`${BASE_URL}/auth/me`);
	if (!response.ok) {
		throw new Error(`Failed to fetch current user: ${response.status}`);
	}
	const user: AuthUser = await response.json();
	return [user];
}

/**
 * Fetch a user by ID.
 * The API is single-PI so we call GET /auth/me and return null if the
 * requested ID does not match the authenticated user.
 */
export async function fetchUser(id: string): Promise<AuthUser | null> {
	if (MOCK) return mock.fetchUser(id);
	const response = await apiFetch(`${BASE_URL}/auth/me`);
	if (!response.ok) {
		throw new Error(`Failed to fetch current user: ${response.status}`);
	}
	const user: AuthUser = await response.json();
	return user.id === id ? user : null;
}

// =============================================================================
// CLIENTS
// =============================================================================

/**
 * Fetch all clients (global list for sidebar)
 */
export async function fetchAllClients(): Promise<Client[]> {
	if (MOCK) return mock.fetchAllClients();
	const response = await apiFetch(`${BASE_URL}/client`);
	if (!response.ok) {
		throw new Error(`Failed to fetch clients: ${response.status}`);
	}
	return response.json();
}

/**
 * Fetch a client by ID
 */
export async function fetchClient(id: string): Promise<Client | null> {
	if (MOCK) return mock.fetchClient(id);
	const response = await apiFetch(`${BASE_URL}/client/${id}`);
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
	if (MOCK) return mock.createClient(request);
	const response = await apiFetch(`${BASE_URL}/client`, {
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
	if (MOCK) return mock.updateClient(id, request);
	const response = await apiFetch(`${BASE_URL}/client/${id}`, {
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
	if (MOCK) return mock.deleteClient(id);
	const response = await apiFetch(`${BASE_URL}/client/${id}`, {
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
	if (MOCK) return mock.fetchClientContacts(clientId);
	const response = await apiFetch(`${BASE_URL}/client/${clientId}/contact`);
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
	if (MOCK) return mock.fetchContact(id);
	const response = await apiFetch(`${BASE_URL}/contact/${id}`);
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
	if (MOCK) return mock.createContact(request);
	const response = await apiFetch(`${BASE_URL}/client/${request.clientID}/contact`, {
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
	if (MOCK) return mock.updateContact(id, request);
	const response = await apiFetch(`${BASE_URL}/contact/${id}`, {
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
	if (MOCK) return mock.deleteContact(id);
	const response = await apiFetch(`${BASE_URL}/contact/${id}`, {
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
	if (MOCK) return mock.fetchAllCases(params);
	const url = new URL(`${BASE_URL}/cases`);

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
	if (MOCK) return mock.fetchCase(id);
	const response = await apiFetch(`${BASE_URL}/cases/${id}`);
	if (!response.ok) {
		throw new Error(`Failed to fetch case: ${response.status}`);
	}

	return response.json();
}

/**
 * Fetch cases by status
 */
export async function fetchCasesByStatus(status: string): Promise<ListCasesResponse> {
	if (MOCK) return mock.fetchCasesByStatus(status);
	return fetchAllCases({ status });
}

/**
 * Fetch cases by client ID
 */
export async function fetchCasesByClient(clientId: string, params?: Omit<FetchCasesParams, 'status'>): Promise<ListCasesResponse> {
	if (MOCK) return mock.fetchCasesByClient(clientId, params);
	const url = new URL(`${BASE_URL}/clients/${clientId}/cases`);

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
	if (MOCK) return mock.createCase(request);
	const response = await apiFetch(`${BASE_URL}/cases`, {
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
	if (MOCK) return mock.updateCase(id, request);
	const response = await apiFetch(`${BASE_URL}/cases/${id}`, {
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
	if (MOCK) return mock.deleteCase(id);
	const response = await apiFetch(`${BASE_URL}/cases/${id}`, {
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
	if (MOCK) return mock.markCaseReady(id);
	const response = await apiFetch(`${BASE_URL}/cases/${id}/mark-ready`, {
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
	if (MOCK) return mock.releaseCase(id);
	const response = await apiFetch(`${BASE_URL}/cases/${id}/release`, {
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
	if (MOCK) return mock.fetchAllCaseSubjects();
	const response = await apiFetch(`${BASE_URL}/subjects`);
	if (!response.ok) {
		throw new Error(`Failed to fetch case subjects: ${response.status}`);
	}
	return response.json();
}

/**
 * Fetch a case subject by ID
 */
export async function fetchCaseSubject(id: string): Promise<CaseSubject | null> {
	if (MOCK) return mock.fetchCaseSubject(id);
	const response = await apiFetch(`${BASE_URL}/subjects/${id}`);
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
	if (MOCK) return mock.createCaseSubject(request);
	const response = await apiFetch(`${BASE_URL}/subjects`, {
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
	if (MOCK) return mock.updateCaseSubject(id, request);
	const response = await apiFetch(`${BASE_URL}/subjects/${id}`, {
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
	if (MOCK) return mock.deleteCaseSubject(id);
	const response = await apiFetch(`${BASE_URL}/subjects/${id}`, {
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
	if (MOCK) return mock.fetchAllCaseTypes();
	const response = await apiFetch(`${BASE_URL}/case-types`);
	if (!response.ok) {
		throw new Error(`Failed to fetch case types: ${response.status}`);
	}
	return response.json();
}

/**
 * Fetch a case type by ID
 */
export async function fetchCaseType(id: string): Promise<CaseType | null> {
	if (MOCK) return mock.fetchCaseType(id);
	const response = await apiFetch(`${BASE_URL}/case-types/${id}`);
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
	if (MOCK) return mock.createCaseType(request);
	const response = await apiFetch(`${BASE_URL}/case-types`, {
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
	if (MOCK) return mock.updateCaseType(id, request);
	const response = await apiFetch(`${BASE_URL}/case-types/${id}`, {
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
	if (MOCK) return mock.deleteCaseType(id);
	const response = await apiFetch(`${BASE_URL}/case-types/${id}`, {
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
	if (MOCK) return mock.fetchAllEstimates();
	const result = await fetchDocuments({ type: 'estimate' });
	return result.data as unknown as Estimate[];
}

/**
 * Fetch estimates for a specific case
 */
export async function fetchCaseEstimates(caseId: string): Promise<Estimate[]> {
	if (MOCK) return mock.fetchCaseEstimates(caseId);
	const result = await fetchDocuments({ type: 'estimate', case_id: caseId });
	return result.data as unknown as Estimate[];
}

/**
 * Fetch an estimate by ID
 */
export async function fetchEstimate(id: string): Promise<Estimate | null> {
	if (MOCK) return mock.fetchEstimate(id);
	const result = await fetchDocument(id, 'estimate');
	return result.data as Estimate || null;
}

/**
 * Fetch estimates by client ID
 */
export async function fetchClientEstimates(clientId: string): Promise<Estimate[]> {
	if (MOCK) return mock.fetchClientEstimates(clientId);
	const result = await fetchDocuments({ type: 'estimate' });
	return (result.data as unknown as Estimate[]).filter(est => (est as Estimate).client_id === clientId);
}

/**
 * Create a new estimate
 */
export async function createEstimate(estimate: Estimate): Promise<DocumentAPIResponse<Estimate>> {
	if (MOCK) return mock.createEstimate(estimate);
	const response = await apiFetch(`${BASE_URL}/estimates`, {
		method: 'POST',
		body: JSON.stringify(estimate),
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		throw new Error(errorBody?.error || `Failed to create estimate: ${response.status}`);
	}

	return response.json();
}

/**
 * Update an estimate (line items only — the backend does not support updating dates or notes)
 */
export async function updateEstimate(id: string, lineItems: EstimateItem[]): Promise<DocumentAPIResponse<Estimate>> {
	if (MOCK) return mock.updateEstimate(id, lineItems);
	const response = await apiFetch(`${BASE_URL}/estimates/${id}`, {
		method: 'PATCH',
		body: JSON.stringify({ line_items: lineItems }),
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		throw new Error(errorBody?.error || `Failed to update estimate: ${response.status}`);
	}

	return response.json();
}

/**
 * Accept an estimate — the backend derives the accepting user from the session
 */
export async function acceptEstimate(id: string): Promise<DocumentAPIResponse> {
	if (MOCK) return mock.acceptEstimate(id);
	const response = await apiFetch(`${BASE_URL}/estimates/${id}/accept`, {
		method: 'POST',
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		throw new Error(errorBody?.error || `Failed to accept estimate: ${response.status}`);
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
	if (MOCK) return mock.fetchAllMandates();
	const result = await fetchDocuments({ type: 'mandate' });
	return result.data as unknown as Mandate[];
}

/**
 * Fetch mandates for a specific case
 */
export async function fetchCaseMandates(caseId: string): Promise<Mandate[]> {
	if (MOCK) return mock.fetchCaseMandates(caseId);
	const result = await fetchDocuments({ type: 'mandate', case_id: caseId });
	return result.data as unknown as Mandate[];
}

/**
 * Fetch a mandate by ID
 */
export async function fetchMandate(id: string): Promise<Mandate | null> {
	if (MOCK) return mock.fetchMandate(id);
	const result = await fetchDocument(id, 'mandate');
	return result.data as Mandate || null;
}

/**
 * Fetch mandates by client ID
 */
export async function fetchClientMandates(clientId: string): Promise<Mandate[]> {
	if (MOCK) return mock.fetchClientMandates(clientId);
	const result = await fetchDocuments({ type: 'mandate' });
	return (result.data as unknown as Mandate[]).filter(mand => (mand as Mandate).client_id === clientId);
}

/**
 * Create a new mandate
 */
export async function createMandate(mandate: Mandate): Promise<DocumentAPIResponse<Mandate>> {
	if (MOCK) return mock.createMandate(mandate);
	const response = await apiFetch(`${BASE_URL}/mandates`, {
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
	if (MOCK) return mock.signMandate(id, request);
	const response = await apiFetch(`${BASE_URL}/mandates/${id}/sign`, {
		method: 'POST',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		if (response.status === 409) {
			throw new ConflictError(errorBody?.error || 'Ce mandat ne peut pas être signé dans son état actuel.');
		}
		throw new Error(errorBody?.error || `Failed to sign mandate: ${response.status}`);
	}

	return response.json();
}

/**
 * Activate a mandate
 */
export async function activateMandate(id: string): Promise<DocumentAPIResponse<Mandate>> {
	if (MOCK) return mock.activateMandate(id);
	const response = await apiFetch(`${BASE_URL}/mandates/${id}/activate`, {
		method: 'POST',
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		if (response.status === 409) {
			throw new ConflictError(errorBody?.error || 'Ce mandat ne peut pas être activé dans son état actuel.');
		}
		throw new Error(errorBody?.error || `Failed to activate mandate: ${response.status}`);
	}

	return response.json();
}

/**
 * Create a contract from a mandate
 */
export async function createContractFromMandate(mandateId: string, contract: Contract): Promise<DocumentAPIResponse<Contract>> {
	if (MOCK) return mock.createContractFromMandate(mandateId, contract);
	const response = await apiFetch(`${BASE_URL}/mandates/${mandateId}/create-contract`, {
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
	if (MOCK) return mock.fetchAllContracts();
	const result = await fetchDocuments({ type: 'contract' });
	return result.data as unknown as Contract[];
}

/**
 * Fetch contracts for a specific case
 */
export async function fetchCaseContracts(caseId: string): Promise<Contract[]> {
	if (MOCK) return mock.fetchCaseContracts(caseId);
	const result = await fetchDocuments({ type: 'contract', case_id: caseId });
	return result.data as unknown as Contract[];
}

/**
 * Fetch a contract by ID
 */
export async function fetchContract(id: string): Promise<Contract | null> {
	if (MOCK) return mock.fetchContract(id);
	const result = await fetchDocument(id, 'contract');
	return result.data as Contract || null;
}

/**
 * Fetch contracts by client ID
 */
export async function fetchClientContracts(clientId: string): Promise<Contract[]> {
	if (MOCK) return mock.fetchClientContracts(clientId);
	const result = await fetchDocuments({ type: 'contract' });
	return (result.data as unknown as Contract[]).filter(cont => (cont as Contract).client_id === clientId);
}

/**
 * Create a new contract
 */
export async function createContract(contract: Contract): Promise<DocumentAPIResponse<Contract>> {
	if (MOCK) return mock.createContract(contract);
	const response = await apiFetch(`${BASE_URL}/contracts`, {
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
	if (MOCK) return mock.signContract(id, request);
	const response = await apiFetch(`${BASE_URL}/contracts/${id}/sign`, {
		method: 'POST',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		if (response.status === 409) {
			throw new ConflictError(errorBody?.error || 'Ce contrat ne peut pas être signé dans son état actuel.');
		}
		throw new Error(errorBody?.error || `Failed to sign contract: ${response.status}`);
	}

	return response.json();
}

/**
 * Activate a contract
 */
export async function activateContract(id: string): Promise<DocumentAPIResponse<Contract>> {
	if (MOCK) return mock.activateContract(id);
	const response = await apiFetch(`${BASE_URL}/contracts/${id}/activate`, {
		method: 'POST',
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		if (response.status === 409) {
			throw new ConflictError(errorBody?.error || 'Ce contrat ne peut pas être activé dans son état actuel.');
		}
		throw new Error(errorBody?.error || `Failed to activate contract: ${response.status}`);
	}

	return response.json();
}

/**
 * Create an invoice from a contract
 */
export async function createInvoiceFromContract(contractId: string): Promise<DocumentAPIResponse<Invoice>> {
	if (MOCK) return mock.createInvoiceFromContract(contractId);
	const response = await apiFetch(`${BASE_URL}/contracts/${contractId}/create-invoice`, {
		method: 'POST',
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		if (response.status === 409) {
			throw new ConflictError(errorBody?.error || 'Impossible de créer une facture depuis ce contrat dans son état actuel.');
		}
		throw new Error(errorBody?.error || `Failed to create invoice from contract: ${response.status}`);
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
	if (MOCK) return mock.fetchAllInvoices();
	const result = await fetchDocuments({ type: 'invoice' });
	return result.data as unknown as Invoice[];
}

/**
 * Fetch invoices for a specific case
 */
export async function fetchCaseInvoices(caseId: string): Promise<Invoice[]> {
	if (MOCK) return mock.fetchCaseInvoices(caseId);
	const result = await fetchDocuments({ type: 'invoice', case_id: caseId });
	return result.data as unknown as Invoice[];
}

/**
 * Fetch an invoice by ID
 */
export async function fetchInvoice(id: string): Promise<Invoice | null> {
	if (MOCK) return mock.fetchInvoice(id);
	const result = await fetchDocument(id, 'invoice');
	return result.data as Invoice || null;
}

/**
 * Fetch invoices by client ID
 */
export async function fetchClientInvoices(clientId: string): Promise<Invoice[]> {
	if (MOCK) return mock.fetchClientInvoices(clientId);
	const result = await fetchDocuments({ type: 'invoice' });
	return (result.data as unknown as Invoice[]).filter(inv => (inv as Invoice).client_id === clientId);
}

/**
 * Fetch invoices by payment status
 */
export async function fetchInvoicesByPaymentStatus(paymentStatus: string): Promise<Invoice[]> {
	if (MOCK) return mock.fetchInvoicesByPaymentStatus(paymentStatus);
	const result = await fetchDocuments({ type: 'invoice' });
	return (result.data as unknown as Invoice[]).filter(inv => (inv as Invoice).payment_status === paymentStatus);
}

/**
 * Create a new invoice
 */
export async function createInvoice(invoice: Invoice): Promise<DocumentAPIResponse<Invoice>> {
	if (MOCK) return mock.createInvoice(invoice);
	const response = await apiFetch(`${BASE_URL}/invoices`, {
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
	if (MOCK) return mock.fetchOverdueInvoices(page, perPage);
	const url = new URL(`${BASE_URL}/invoices/overdue`);
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
	if (MOCK) return mock.processPayment(id, request);
	const response = await apiFetch(`${BASE_URL}/documents/${id}/pay`, {
		method: 'POST',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		if (response.status === 409) {
			throw new ConflictError(errorBody?.error || 'Impossible d\'enregistrer un paiement sur cette facture dans son état actuel.');
		}
		throw new Error(errorBody?.error || `Failed to process payment: ${response.status}`);
	}

	return response.json();
}

/**
 * Void an invoice
 */
export async function voidInvoice(id: string): Promise<DocumentAPIResponse<Invoice>> {
	if (MOCK) return mock.voidInvoice(id);
	const response = await apiFetch(`${BASE_URL}/documents/${id}/void`, {
		method: 'POST',
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		if (response.status === 409) {
			throw new ConflictError(errorBody?.error || 'Impossible d\'annuler cette facture dans son état actuel.');
		}
		throw new Error(errorBody?.error || `Failed to void invoice: ${response.status}`);
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
 * Fetch images for a specific case.
 * Calls GET /case/{caseId}/media?type=image and maps the response to ApiImage.
 */
export async function fetchCaseImages(caseId: string): Promise<ApiImage[]> {
	if (MOCK) return mock.fetchCaseImages(caseId);
	const url = new URL(`${BASE_URL}/case/${caseId}/media`);
	url.searchParams.set('type', 'image');

	const response = await apiFetch(url.toString());
	if (!response.ok) {
		throw new Error(`Failed to fetch case images: ${response.status}`);
	}

	const data = await response.json();
	const mediaItems = data.media ?? data;

	return mediaItems.map((item: any) => ({
		id: item.id,
		url: item.url,
	}));
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
	if (MOCK) return mock.requestAITextOperation(request);
	const timeout = AI_CONFIG.DEFAULT_TIMEOUT;

	const controller = new AbortController();
	const timeoutId = setTimeout(() => controller.abort(), timeout);

	try {
		const response = await apiFetch(`${BASE_URL}/ai/text/transform`, {
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
	if (MOCK) return mock.sendChatMessage(caseId, request);
	const url = new URL(`${BASE_URL}/api/ai/chat/message`);
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
	if (MOCK) return mock.getChatConversation(conversationId);
	const response = await apiFetch(`${BASE_URL}/api/ai/chat/conversations/${conversationId}`);

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
	if (MOCK) return mock.listChatConversations(caseId);
	const url = new URL(`${BASE_URL}/api/ai/chat/conversations`);
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
	if (MOCK) return mock.deleteChatConversation(conversationId);
	const response = await apiFetch(`${BASE_URL}/api/ai/chat/conversations/${conversationId}`, {
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
	if (MOCK) return mock.fetchDocuments(params);
	const url = new URL(`${BASE_URL}/documents`);

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
	if (MOCK) return mock.createDocument(request);
	const response = await apiFetch(`${BASE_URL}/documents`, {
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
	if (MOCK) return mock.fetchDocument(id, type);
	const response = await apiFetch(`${BASE_URL}/documents/${id}/${type}`);
	if (!response.ok) {
		throw new Error(`Failed to fetch document: ${response.status}`);
	}

	return response.json();
}

/**
 * Update an existing document
 */
export async function updateDocument(id: string, type: string, request: UpdateDocumentRequest): Promise<DocumentAPIResponse> {
	if (MOCK) return mock.updateDocument(id, type, request);
	const response = await apiFetch(`${BASE_URL}/documents/${id}/${type}`, {
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
	if (MOCK) return mock.deleteDocument(id, type);
	const response = await apiFetch(`${BASE_URL}/documents/${id}/${type}`, {
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
	if (MOCK) return mock.sendDocument(id, type, request);
	const response = await apiFetch(`${BASE_URL}/documents/${id}/${type}/send`, {
		method: 'POST',
		body: JSON.stringify(request),
	});

	if (!response.ok) {
		const errorBody = await response.json().catch(() => null);
		if (response.status === 409) {
			throw new ConflictError(errorBody?.error || 'Ce document ne peut pas être envoyé dans son état actuel.');
		}
		throw new Error(errorBody?.error || `Failed to send document: ${response.status}`);
	}

	return response.json();
}

/**
 * Sign a document
 */
export async function signDocument(id: string, type: string, request: SignDocumentRequest): Promise<DocumentAPIResponse> {
	if (MOCK) return mock.signDocument(id, type, request);
	const response = await apiFetch(`${BASE_URL}/documents/${id}/${type}/sign`, {
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
	if (MOCK) return mock.archiveDocument(id, type);
	const response = await apiFetch(`${BASE_URL}/documents/${id}/${type}/archive`, {
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
	if (MOCK) return mock.fetchDocumentHistory(id, type, page, perPage);
	const url = new URL(`${BASE_URL}/documents/${id}/${type}/history`);
	url.searchParams.set('page', page.toString());
	url.searchParams.set('per_page', perPage.toString());

	const response = await apiFetch(url.toString());
	if (!response.ok) {
		throw new Error(`Failed to fetch document history: ${response.status}`);
	}

	return response.json();
}

// =============================================================================
// SEARCH
// =============================================================================

/**
 * Search across cases, clients, and contacts using fuzzy matching.
 */
export async function searchAll(query: string): Promise<SearchResult[]> {
	if (MOCK) return mock.searchAll(query);
	// TODO: implement backend endpoint GET /search?q=<query>
	const response = await apiFetch(`${BASE_URL}/search?q=${encodeURIComponent(query)}`);
	if (!response.ok) throw new Error(`Search failed: ${response.status}`);
	const data: { results: SearchResult[] } = await response.json();
	return data.results;
}

/**
 * Fetch document workflow for a case (full document chain)
 */
export async function fetchDocumentWorkflow(caseId: string): Promise<DocumentWorkflowResponse> {
	if (MOCK) return mock.fetchDocumentWorkflow(caseId);
	const response = await apiFetch(`${BASE_URL}/cases/${caseId}/document-workflow`);
	if (!response.ok) {
		throw new Error(`Failed to fetch document workflow: ${response.status}`);
	}

	return response.json();
}

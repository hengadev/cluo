import type {
	AuthUser,
	Case,
	CaseSubject,
	CaseType,
	Client,
	Contact,
	Contract,
	CreateCaseRequest,
	DocumentAPIResponse,
	DocumentHistoryResponse,
	DocumentListResponse,
	DocumentWorkflowResponse,
	Estimate,
	EstimateItem,
	Invoice,
	ListCasesResponse,
	Mandate,
	OverdueInvoicesResponse,
	PaymentRequest,
	ReleaseResponse,
	SendDocumentRequest,
	SignDocumentRequest,
	UpdateDocumentRequest,
	CreateDocumentRequest,
} from '../types/entities';
import type {
	SendMessageRequest,
	SendMessageResponse,
	GetConversationResponse,
	ListConversationsResponse,
} from '../types/chat';

import { users } from '../mockData/users';
import { getAllClients, getClientById } from '../mockData/clients';
import { getAllContacts, getContactById, getContactsByClientId } from '../mockData/contacts';
import { caseSubjects } from '../mockData/caseSubjects';
import { caseTypes, getCaseTypeById } from '../mockData/caseTypes';
import { getAllCases, getCaseById, getCasesByStatus, getCasesByClientId } from '../mockData/cases';

// Inline types to avoid circular import with api.ts
interface AITextOperationRequest { text: string; operation: string; language: string; }
interface AITextOperationResponse { result: string; }
interface ApiImage { id: string; url: string; }
interface FetchCasesParams { page?: number; pageSize?: number; status?: string; }

const MOCK_ROLE = import.meta.env.VITE_MOCK_USER_ROLE as AuthUser['role'] | undefined;

function mockUser(): AuthUser {
	return { id: 'mock-user', email: 'dev@cluo.local', role: MOCK_ROLE ?? 'admin' };
}

function paginate<T>(items: T[], page = 1, pageSize = 20): { items: T[]; total: number } {
	const start = (page - 1) * pageSize;
	return { items: items.slice(start, start + pageSize), total: items.length };
}

function casesToResponse(items: Case[], page = 1, pageSize = 20): ListCasesResponse {
	const { items: sliced, total } = paginate(items, page, pageSize);
	return {
		cases: sliced as unknown as import('../types/entities').Case[],
		pagination: { page, pageSize, totalItems: total, totalPages: Math.ceil(total / pageSize) },
	};
}

function ok<T>(data?: T): DocumentAPIResponse<T> {
	return { success: true, data };
}

// =============================================================================
// USERS
// =============================================================================

export async function fetchAllUsers(): Promise<AuthUser[]> {
	return [mockUser(), ...users.filter(u => u.id !== 'mock-user')] as AuthUser[];
}

export async function fetchUser(id: string): Promise<AuthUser | null> {
	if (id === 'mock-user') return mockUser();
	return (users.find(u => u.id === id) as AuthUser | undefined) ?? null;
}

// =============================================================================
// CLIENTS
// =============================================================================

export async function fetchAllClients(): Promise<Client[]> {
	return getAllClients() as unknown as Client[];
}

export async function fetchClient(id: string): Promise<Client | null> {
	return (getClientById(id) as unknown as Client) ?? null;
}

export async function createClient(request: { name: string; type?: string }): Promise<Client> {
	return { id: crypto.randomUUID(), name: request.name, type: (request.type ?? 'company') as Client['type'] };
}

export async function updateClient(id: string, request: { name?: string; type?: string }): Promise<Client> {
	const existing = getClientById(id);
	return {
		id,
		name: request.name ?? existing?.name ?? '',
		type: (request.type ?? existing?.type ?? 'company') as Client['type'],
	};
}

export async function deleteClient(_id: string): Promise<void> {}

export async function fetchClientContacts(clientId: string): Promise<Contact[]> {
	return getContactsByClientId(clientId) as unknown as Contact[];
}

// =============================================================================
// CONTACTS
// =============================================================================

export async function fetchContact(id: string): Promise<Contact | null> {
	return (getContactById(id) as unknown as Contact) ?? null;
}

export async function createContact(request: {
	clientID: string; lastname: string; firstname: string;
	email?: string; phone?: string; position?: string;
}): Promise<Contact> {
	return {
		id: crypto.randomUUID(),
		clientID: request.clientID,
		lastname: request.lastname,
		firstname: request.firstname,
		email: request.email ?? '',
		phone: request.phone ?? '',
		position: request.position ?? '',
		createdAt: new Date().toISOString(),
	};
}

export async function updateContact(id: string, request: {
	lastname?: string; firstname?: string; email?: string; phone?: string; position?: string;
}): Promise<Contact> {
	const existing = getContactById(id);
	return {
		id,
		clientID: existing?.clientID ?? '',
		lastname: request.lastname ?? existing?.lastname ?? '',
		firstname: request.firstname ?? existing?.firstname ?? '',
		email: request.email ?? existing?.email ?? '',
		phone: request.phone ?? existing?.phone ?? '',
		position: request.position ?? existing?.position ?? '',
		createdAt: existing?.createdAt ?? new Date().toISOString(),
	};
}

export async function deleteContact(_id: string): Promise<void> {}

// =============================================================================
// CASES
// =============================================================================

export async function fetchAllCases(params?: FetchCasesParams): Promise<ListCasesResponse> {
	const all = params?.status ? getCasesByStatus(params.status as any) : getAllCases();
	return casesToResponse(all as any, params?.page, params?.pageSize);
}

export async function fetchCase(id: string): Promise<Case> {
	const found = getCaseById(id);
	if (!found) throw new Error(`Case not found: ${id}`);
	return found as unknown as Case;
}

export async function fetchCasesByStatus(status: string): Promise<ListCasesResponse> {
	return casesToResponse(getCasesByStatus(status as any) as any);
}

export async function fetchCasesByClient(
	clientId: string,
	params?: { page?: number; pageSize?: number },
): Promise<ListCasesResponse> {
	return casesToResponse(getCasesByClientId(clientId) as any, params?.page, params?.pageSize);
}

export async function createCase(request: CreateCaseRequest): Promise<Case> {
	const now = new Date().toISOString();
	return {
		id: crypto.randomUUID(),
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
		createdAt: now,
		updatedAt: now,
	};
}

export async function updateCase(id: string, request: Partial<CreateCaseRequest>): Promise<Case> {
	const existing = getCaseById(id) as unknown as Case;
	if (!existing) throw new Error(`Case not found: ${id}`);
	return { ...existing, ...request, id, updatedAt: new Date().toISOString() } as Case;
}

export async function deleteCase(_id: string): Promise<void> {}

export async function markCaseReady(_id: string): Promise<void> {}

export async function releaseCase(id: string): Promise<ReleaseResponse> {
	return {
		caseId: id,
		tokenId: crypto.randomUUID(),
		rawToken: 'mock-token-' + id,
		portalUrl: `https://portal.cluo.local/cases/${id}`,
		expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(),
	};
}

// =============================================================================
// CASE SUBJECTS
// =============================================================================

export async function fetchAllCaseSubjects(): Promise<CaseSubject[]> {
	return caseSubjects;
}

export async function fetchCaseSubject(id: string): Promise<CaseSubject | null> {
	return caseSubjects.find(s => s.id === id) ?? null;
}

export async function createCaseSubject(request: {
	lastname: string; firstname: string; email?: string; phone?: string;
	city?: string; postalCode?: string; address1?: string; address2?: string;
	occupation?: string; notes?: string;
}): Promise<CaseSubject> {
	return {
		id: crypto.randomUUID(),
		lastname: request.lastname,
		firstname: request.firstname,
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
}

export async function updateCaseSubject(id: string, request: {
	lastname?: string; firstname?: string; email?: string; phone?: string;
	city?: string; postalCode?: string; address1?: string; address2?: string;
	occupation?: string; notes?: string;
}): Promise<CaseSubject> {
	const existing = caseSubjects.find(s => s.id === id);
	return {
		id,
		lastname: request.lastname ?? existing?.lastname ?? '',
		firstname: request.firstname ?? existing?.firstname ?? '',
		email: request.email ?? existing?.email ?? '',
		phone: request.phone ?? existing?.phone ?? '',
		address1: request.address1 ?? existing?.address1 ?? '',
		address2: request.address2 ?? existing?.address2 ?? '',
		city: request.city ?? existing?.city ?? '',
		postalCode: request.postalCode ?? existing?.postalCode ?? '',
		occupation: request.occupation ?? existing?.occupation ?? '',
		notes: request.notes ?? existing?.notes ?? '',
		createdAt: existing?.createdAt ?? new Date().toISOString(),
	};
}

export async function deleteCaseSubject(_id: string): Promise<void> {}

// =============================================================================
// CASE TYPES
// =============================================================================

export async function fetchAllCaseTypes(): Promise<CaseType[]> {
	return caseTypes;
}

export async function fetchCaseType(id: string): Promise<CaseType | null> {
	return getCaseTypeById(id) ?? null;
}

export async function createCaseType(request: { name: string }): Promise<CaseType> {
	const now = new Date().toISOString();
	return { id: crypto.randomUUID(), name: request.name, createdAt: now, updatedAt: now };
}

export async function updateCaseType(id: string, request: { name: string }): Promise<CaseType> {
	const existing = getCaseTypeById(id);
	return {
		id,
		name: request.name,
		createdAt: existing?.createdAt ?? new Date().toISOString(),
		updatedAt: new Date().toISOString(),
	};
}

export async function deleteCaseType(_id: string): Promise<void> {}

// =============================================================================
// ESTIMATES — stub (empty lists, no-op writes)
// =============================================================================

export async function fetchAllEstimates(): Promise<Estimate[]> { return []; }
export async function fetchCaseEstimates(_caseId: string): Promise<Estimate[]> { return []; }
export async function fetchEstimate(_id: string): Promise<Estimate | null> { return null; }
export async function fetchClientEstimates(_clientId: string): Promise<Estimate[]> { return []; }

export async function createEstimate(estimate: Estimate): Promise<DocumentAPIResponse<Estimate>> {
	return ok({ ...estimate, id: crypto.randomUUID() });
}

export async function updateEstimate(id: string, lineItems: EstimateItem[]): Promise<DocumentAPIResponse<Estimate>> {
	return ok({ id, line_items: lineItems } as unknown as Estimate);
}

export async function acceptEstimate(_id: string): Promise<DocumentAPIResponse> {
	return ok();
}

// =============================================================================
// MANDATES — stub
// =============================================================================

export async function fetchAllMandates(): Promise<Mandate[]> { return []; }
export async function fetchCaseMandates(_caseId: string): Promise<Mandate[]> { return []; }
export async function fetchMandate(_id: string): Promise<Mandate | null> { return null; }
export async function fetchClientMandates(_clientId: string): Promise<Mandate[]> { return []; }

export async function createMandate(mandate: Mandate): Promise<DocumentAPIResponse<Mandate>> {
	return ok({ ...mandate, id: crypto.randomUUID() });
}

export async function signMandate(_id: string, _request: SignDocumentRequest): Promise<DocumentAPIResponse<Mandate>> {
	return ok();
}

export async function activateMandate(_id: string): Promise<DocumentAPIResponse<Mandate>> {
	return ok();
}

export async function createContractFromMandate(_mandateId: string, contract: Contract): Promise<DocumentAPIResponse<Contract>> {
	return ok({ ...contract, id: crypto.randomUUID() });
}

// =============================================================================
// CONTRACTS — stub
// =============================================================================

export async function fetchAllContracts(): Promise<Contract[]> { return []; }
export async function fetchCaseContracts(_caseId: string): Promise<Contract[]> { return []; }
export async function fetchContract(_id: string): Promise<Contract | null> { return null; }
export async function fetchClientContracts(_clientId: string): Promise<Contract[]> { return []; }

export async function createContract(contract: Contract): Promise<DocumentAPIResponse<Contract>> {
	return ok({ ...contract, id: crypto.randomUUID() });
}

export async function signContract(_id: string, _request: SignDocumentRequest): Promise<DocumentAPIResponse<Contract>> {
	return ok();
}

export async function activateContract(_id: string): Promise<DocumentAPIResponse<Contract>> {
	return ok();
}

export async function createInvoiceFromContract(_contractId: string): Promise<DocumentAPIResponse<Invoice>> {
	return ok();
}

// =============================================================================
// INVOICES — stub
// =============================================================================

export async function fetchAllInvoices(): Promise<Invoice[]> { return []; }
export async function fetchCaseInvoices(_caseId: string): Promise<Invoice[]> { return []; }
export async function fetchInvoice(_id: string): Promise<Invoice | null> { return null; }
export async function fetchClientInvoices(_clientId: string): Promise<Invoice[]> { return []; }
export async function fetchInvoicesByPaymentStatus(_status: string): Promise<Invoice[]> { return []; }

export async function createInvoice(invoice: Invoice): Promise<DocumentAPIResponse<Invoice>> {
	return ok({ ...invoice, id: crypto.randomUUID() });
}

export async function fetchOverdueInvoices(_page = 1, _perPage = 20): Promise<OverdueInvoicesResponse> {
	return { success: true, data: [], total: 0, page: _page, per_page: _perPage };
}

export async function processPayment(_id: string, _request: PaymentRequest): Promise<DocumentAPIResponse<Invoice>> {
	return ok();
}

export async function voidInvoice(_id: string): Promise<DocumentAPIResponse<Invoice>> {
	return ok();
}

// =============================================================================
// IMAGES — stub
// =============================================================================

export async function fetchCaseImages(_caseId: string): Promise<ApiImage[]> {
	return [];
}

// =============================================================================
// AI TEXT OPERATIONS — stub
// =============================================================================

export async function requestAITextOperation(
	request: AITextOperationRequest,
): Promise<AITextOperationResponse> {
	return { result: `[Mode démo] ${request.text}` };
}

// =============================================================================
// AI CHAT OPERATIONS — stub
// =============================================================================

export async function sendChatMessage(
	_caseId: string,
	_request: SendMessageRequest,
): Promise<SendMessageResponse> {
	const conversationId = crypto.randomUUID();
	return {
		conversationId,
		userMessageId: crypto.randomUUID(),
		assistantMessage: {
			id: crypto.randomUUID(),
			conversationId,
			role: 'assistant',
			content: 'Fonctionnalité non disponible en mode démo.',
			createdAt: new Date().toISOString(),
		},
	};
}

export async function getChatConversation(_conversationId: string): Promise<GetConversationResponse> {
	return { conversation: null as any, messages: [] };
}

export async function listChatConversations(_caseId: string): Promise<ListConversationsResponse> {
	return { conversations: [] };
}

export async function deleteChatConversation(_conversationId: string): Promise<void> {}

// =============================================================================
// DOCUMENTS (generic) — stub
// =============================================================================

export async function fetchDocuments(_params?: unknown): Promise<DocumentListResponse> {
	return { success: true, data: [], total: 0, page: 1, per_page: 20 };
}

export async function createDocument(_request: CreateDocumentRequest): Promise<DocumentAPIResponse> {
	return ok();
}

export async function fetchDocument(_id: string, _type: string): Promise<DocumentAPIResponse> {
	return ok();
}

export async function updateDocument(
	_id: string,
	_type: string,
	_request: UpdateDocumentRequest,
): Promise<DocumentAPIResponse> {
	return ok();
}

export async function deleteDocument(_id: string, _type: string): Promise<DocumentAPIResponse> {
	return ok();
}

export async function sendDocument(
	_id: string,
	_type: string,
	_request: SendDocumentRequest,
): Promise<DocumentAPIResponse> {
	return ok();
}

export async function signDocument(
	_id: string,
	_type: string,
	_request: SignDocumentRequest,
): Promise<DocumentAPIResponse> {
	return ok();
}

export async function archiveDocument(_id: string, _type: string): Promise<DocumentAPIResponse> {
	return ok();
}

export async function fetchDocumentHistory(
	_id: string,
	_type: string,
	page = 1,
	perPage = 20,
): Promise<DocumentHistoryResponse> {
	return { success: true, data: [], total: 0, page, per_page: perPage };
}

export async function fetchDocumentWorkflow(_caseId: string): Promise<DocumentWorkflowResponse> {
	return { success: true, data: [] };
}

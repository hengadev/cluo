/**
 * Type definitions for all entities in the investigation system.
 * These types match the backend data models.
 */

// =============================================================================
// USER TYPES
// =============================================================================

export type UserRole = 'admin' | 'investigator' | 'viewer';

export interface AuthUser {
	id: string;
	email: string;
	role: UserRole;
}

// =============================================================================
// CLIENT TYPES
// =============================================================================

export type ClientType = 'person' | 'insurance' | 'lawyer' | 'company' | 'government';

export interface Client {
	id: string;
	name: string;
	type: ClientType;
}

// =============================================================================
// CONTACT TYPES
// =============================================================================

export interface Contact {
	id: string;
	clientID: string;
	lastname: string;
	firstname: string;
	email: string;
	phone: string;
	position: string;
	createdAt: string;
}

// =============================================================================
// CASE SUBJECT TYPES
// =============================================================================

export type SubjectRole = 'victim' | 'suspect' | 'witness' | 'claimant' | 'representative';

export interface CaseSubject {
	id: string;
	firstname: string;
	lastname: string;
	email: string;
	phone: string;
	address1: string;
	address2: string;
	city: string;
	postalCode: string;
	occupation: string;
	notes: string;
	createdAt: string;
}

// =============================================================================
// CASE TYPES
// =============================================================================

export type CaseStatus = 'in_progress' | 'ready' | 'released';

export interface CaseType {
	id: string;
	name: string;
	createdAt: string;
	updatedAt: string;
}
export type LocationType = 'home' | 'business' | 'public' | 'vehicle' | 'other';

export interface Case {
	id: string;
	title: string;
	description: string;
	clientId: string;
	assignedContactID: string | null;
	caseSubjectId: string | null;
	externalReference: string | null;
	caseTypeId: string | null;
	status: CaseStatus;
	// Location fields
	placename: string | null;
	address1: string | null;
	address2: string | null;
	city: string | null;
	postalCode: string | null;
	country: string | null;
	latitude: string | null;
	longitude: string | null;
	locationType: LocationType | null;
	locationNotes: string | null;
	createdAt: string;
	updatedAt: string;
}

export interface PaginationInfo {
	page: number;
	pageSize: number;
	totalItems: number;
	totalPages: number;
}

export interface ListCasesResponse {
	cases: Case[];
	pagination: PaginationInfo;
}

export interface CreateCaseRequest {
	title: string;
	description: string;
	clientId: string;
	status: CaseStatus;
	assignedContactID?: string;
	caseSubjectId?: string | null;
	caseTypeId?: string | null;
	externalReference?: string;
	placename?: string;
	address1?: string;
	address2?: string;
	city?: string;
	postalCode?: string;
	country?: string;
	latitude?: string;
	longitude?: string;
	locationType?: LocationType;
	locationNotes?: string;
}

export interface ReleaseResponse {
	caseId: string;
	tokenId: string;
	rawToken: string;
	portalUrl: string;
	expiresAt: string;
}

// =============================================================================
// DOCUMENT TYPES (Common)
// =============================================================================

export type DocumentStatus = 'draft' | 'sent' | 'signed' | 'active' | 'archived' | 'cancelled' | 'rejected' | 'expired';

export interface Signature {
	id: string;
	name: string;
	role: string;
	signature_file_url?: string;
	method?: string;
	signed_at: string;
}

// =============================================================================
// DOCUMENT LIST RESPONSE
// =============================================================================

export interface DocumentSummary {
	id: string;
	case_id: string;
	client_id: string;
	type: string;
	status: DocumentStatus;
	document_ref: string;
	created_at: string;
	updated_at: string;
}

export interface DocumentListResponse {
	success: boolean;
	data: DocumentSummary[];
	total: number;
	page: number;
	per_page: number;
}

// =============================================================================
// ESTIMATE TYPES
// =============================================================================

export interface EstimateItem {
	description: string;
	quantity: number;
	unit_price: number;
	subtotal: number;
}

export interface Estimate {
	id: string;
	case_id: string;
	client_id: string;
	estimate_number: string;
	issue_date: string;
	valid_until?: string;
	line_items: EstimateItem[];
	estimated_total: number;
	notes?: string;
	accepted: boolean;
	accepted_at?: string;
	accepted_by?: string;
	status: DocumentStatus;
	created_at: string;
	updated_at: string;
}

// =============================================================================
// MANDATE TYPES
// =============================================================================

export interface Mandate {
	id: string;
	case_id: string;
	client_id: string;
	mandate_number: string;
	issue_date: string;
	scope_of_work: string;
	valid_from: string;
	valid_until?: string;
	terms_conditions: string;
	client_signature?: Signature;
	investigator_signature?: Signature;
	linked_estimate_id?: string;
	special_instructions?: string;
	jurisdiction?: string;
	status: DocumentStatus;
	created_at: string;
	updated_at: string;
}

// =============================================================================
// CONTRACT TYPES
// =============================================================================

export interface Contract {
	id: string;
	case_id: string;
	client_id: string;
	contract_number: string;
	start_date: string;
	end_date?: string;
	scope_of_services: string;
	payment_terms: string;
	confidentiality: string;
	termination_clause: string;
	signatures: Signature[];
	linked_mandate_id?: string;
	contract_value?: number;
	currency?: string;
	renewal_terms?: string;
	governing_law?: string;
	status: DocumentStatus;
	created_at: string;
	updated_at: string;
}

// =============================================================================
// INVOICE TYPES
// =============================================================================

export type PaymentStatus = 'unpaid' | 'paid' | 'partially_paid' | 'overdue' | 'refunded' | 'void';

export interface InvoiceItem {
	description: string;
	quantity: number;
	unit_price: number;
	subtotal: number;
}

export interface Invoice {
	id: string;
	case_id: string;
	client_id: string;
	invoice_number: string;
	issue_date: string;
	due_date: string;
	line_items: InvoiceItem[];
	total_amount: number;
	tax_rate: number;
	tax_amount: number;
	notes?: string;
	payment_status: PaymentStatus;
	paid_at?: string;
	paid_amount?: number;
	payment_method?: string;
	linked_contract_id?: string;
	currency?: string;
	payment_terms?: string;
	late_fee?: number;
	late_fee_rate?: number;
	status: DocumentStatus;
	created_at: string;
	updated_at: string;
}

// =============================================================================
// LEGACY/COMPATIBILITY TYPES
// =============================================================================

/**
 * @deprecated Use Case interface instead
 */
export interface ApiCase {
	id: string;
	title: string;
	description?: string;
}

/**
 * @deprecated Use specific image types instead
 */
export interface ApiImage {
	id: string;
	url: string;
}

// =============================================================================
// DOCUMENT API REQUEST/RESPONSE TYPES
// =============================================================================

export interface CreateDocumentRequest {
	type: string;
	case_id: string;
	client_id: string;
	data: any;
}

export interface UpdateDocumentRequest {
	data: any;
	reason?: string;
}

export interface SendDocumentRequest {
	recipients: string[];
	subject?: string;
	message?: string;
	send_email: boolean;
	send_sms: boolean;
}

export interface SignDocumentRequest {
	signer_name: string;
	signer_role: string;
	signature_file_url?: string;
	method: 'e-sign' | 'wet' | 'pdf-stamp' | 'third-party';
	ip_address?: string;
	user_agent?: string;
}

export interface PaymentRequest {
	amount: number;
	payment_method: string;
	transaction_id?: string;
	notes?: string;
}

export interface DocumentVersion {
	id: number;
	document_id: string;
	doc_type: string;
	version: number;
	author_id?: string;
	data: any;
	created_at: string;
	reason?: string;
}

export interface DocumentHistoryResponse {
	success: boolean;
	data: DocumentVersion[];
	total: number;
	page: number;
	per_page: number;
}

export interface DocumentWorkflowResponse {
	success: boolean;
	data: DocumentSummary[];
}

export interface OverdueInvoicesResponse {
	success: boolean;
	data: Invoice[];
	total: number;
	page: number;
	per_page: number;
}

export interface DocumentAPIResponse<T = any> {
	success: boolean;
	data?: T;
	error?: string;
}

// =============================================================================
// SEARCH TYPES
// =============================================================================

export interface SearchResultMatch {
	key: string;
	indices: readonly [number, number][];
	value?: string;
}

export interface SearchResult {
	type: 'case' | 'client' | 'contact';
	score: number;
	item: Case | Client | Contact;
	clientName?: string;
	matches?: SearchResultMatch[];
}

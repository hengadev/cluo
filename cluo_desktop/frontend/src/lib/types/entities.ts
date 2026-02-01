/**
 * Type definitions for all entities in the investigation system.
 * These types match the backend data models.
 */

// =============================================================================
// USER TYPES
// =============================================================================

export type UserRole = 'admin' | 'investigator' | 'viewer';

export interface User {
	id: string;
	email: string;
	firstName: string;
	lastName: string;
	role: UserRole;
	createdAt: string;
}

// =============================================================================
// CLIENT TYPES
// =============================================================================

export type ClientType = 'person' | 'insurance' | 'lawyer' | 'company' | 'government';

export interface Client {
	id: string;
	name: string;
	type: ClientType;
	createdAt: string;
}

// =============================================================================
// CONTACT TYPES
// =============================================================================

export interface Contact {
	id: string;
	clientId: string;
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
	firstName: string;
	lastName: string;
	email: string | null;
	phone: string | null;
	address: string | null;
	city: string | null;
	postalCode: string | null;
	country: string;
	occupation: string | null;
	notes: string | null;
	createdAt: string;
	updatedAt: string;
}

export interface CaseSubjectAssignment {
	caseId: string;
	subjectId: string;
	role: SubjectRole;
}

// =============================================================================
// CASE TYPES
// =============================================================================

export type CaseStatus = 'draft' | 'in_progress' | 'ready' | 'released';
export type LocationType = 'home' | 'business' | 'public' | 'vehicle' | 'other';

export interface Case {
	id: string;
	title: string;
	description: string;
	clientId: string;
	assignedContactId: string | null;
	caseSubjectIds: string[];
	externalReference: string | null;
	caseType: string;
	status: CaseStatus;
	// Location fields
	placename: string;
	address1: string;
	address2: string | null;
	city: string;
	postalCode: string;
	country: string;
	latitude: number;
	longitude: number;
	locationType: LocationType;
	locationNotes: string | null;
	createdAt: string;
	updatedAt: string;
}

// =============================================================================
// DOCUMENT TYPES (Common)
// =============================================================================

export type DocumentStatus = 'draft' | 'sent' | 'signed' | 'active' | 'archived' | 'cancelled' | 'rejected' | 'expired';

export interface DocumentLineItem {
	description: string;
	quantity: number;
	unitPrice: number;
	total: number;
}

// =============================================================================
// ESTIMATE TYPES
// =============================================================================

export interface Estimate {
	id: string;
	caseId: string;
	clientId: string;
	estimateNumber: string;
	issueDate: string;
	validUntil: string;
	lineItems: DocumentLineItem[];
	estimatedTotal: number;
	notes: string;
	accepted: boolean;
	acceptedAt: string | null;
	acceptedBy: string | null;
	status: DocumentStatus;
	createdAt: string;
	updatedAt: string;
}

// =============================================================================
// MANDATE TYPES
// =============================================================================

export interface MandateSignature {
	name: string;
	date: string;
}

export interface Mandate {
	id: string;
	caseId: string;
	clientId: string;
	mandateNumber: string;
	issueDate: string;
	scopeOfWork: string;
	validFrom: string;
	validUntil: string;
	termsConditions: string;
	clientSignature: MandateSignature | null;
	investigatorSignature: MandateSignature | null;
	linkedEstimateId: string | null;
	specialInstructions: string | null;
	jurisdiction: string;
	status: DocumentStatus;
	createdAt: string;
	updatedAt: string;
}

// =============================================================================
// CONTRACT TYPES
// =============================================================================

export interface ContractSignature {
	name: string;
	date: string;
	role: string;
}

export interface Contract {
	id: string;
	caseId: string;
	clientId: string;
	contractNumber: string;
	startDate: string;
	endDate: string;
	scopeOfServices: string;
	paymentTerms: string;
	confidentiality: string;
	terminationClause: string;
	signatures: ContractSignature[];
	linkedMandateId: string | null;
	contractValue: number;
	currency: string;
	renewalTerms: string;
	governingLaw: string;
	status: DocumentStatus;
	createdAt: string;
	updatedAt: string;
}

// =============================================================================
// INVOICE TYPES
// =============================================================================

export type PaymentStatus = 'unpaid' | 'paid' | 'partially_paid' | 'overdue' | 'refunded' | 'void';

export interface Invoice {
	id: string;
	caseId: string;
	clientId: string;
	invoiceNumber: string;
	issueDate: string;
	dueDate: string;
	lineItems: DocumentLineItem[];
	totalAmount: number;
	taxRate: number;
	taxAmount: number;
	paymentStatus: PaymentStatus;
	paidAt: string | null;
	paidAmount: number | null;
	paymentMethod: string | null;
	linkedContractId: string | null;
	currency: string;
	paymentTerms: string;
	lateFee: number | null;
	lateFeeRate: number | null;
	status: DocumentStatus;
	createdAt: string;
	updatedAt: string;
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

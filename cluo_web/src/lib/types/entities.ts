/**
 * Core entity type unions used across the application.
 * These types mirror the backend data models.
 */

// =============================================================================
// STATUS TYPES
// =============================================================================

export type CaseStatus = 'in_progress' | 'ready' | 'released';

export type DocumentStatus =
	| 'draft'
	| 'sent'
	| 'signed'
	| 'active'
	| 'archived'
	| 'cancelled'
	| 'rejected'
	| 'expired';

export type PaymentStatus = 'unpaid' | 'paid' | 'partially_paid' | 'overdue' | 'refunded' | 'void';

// =============================================================================
// ROLE TYPES
// =============================================================================

export type UserRole = 'admin' | 'investigator' | 'viewer';

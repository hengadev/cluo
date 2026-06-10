import type { CaseStatus, ClientType, DocumentStatus, PaymentStatus, UserRole } from '$lib/types/entities';

// ---------------------------------------------------------------------------
// Case status badges
// ---------------------------------------------------------------------------

const caseStatusClasses: Record<CaseStatus, string> = {
	in_progress: 'bg-accent text-accent-foreground',
	ready: 'bg-success text-success-foreground',
	released: 'bg-surface text-foreground-alt'
};

export function caseStatusBadge(status: CaseStatus): string {
	return caseStatusClasses[status];
}

// ---------------------------------------------------------------------------
// Document status badges
// ---------------------------------------------------------------------------

const documentStatusClasses: Record<DocumentStatus, string> = {
	draft: 'bg-surface text-foreground-alt',
	sent: 'bg-accent text-accent-foreground',
	signed: 'bg-success text-success-foreground',
	active: 'bg-success text-success-foreground',
	archived: 'bg-surface text-foreground-alt',
	cancelled: 'bg-destructive/10 text-destructive',
	rejected: 'bg-destructive/10 text-destructive',
	expired: 'bg-tertiary/20 text-foreground'
};

export function documentStatusBadge(status: DocumentStatus): string {
	return documentStatusClasses[status];
}

// ---------------------------------------------------------------------------
// Payment status badges
// ---------------------------------------------------------------------------

const paymentStatusClasses: Record<PaymentStatus, string> = {
	unpaid: 'bg-surface text-foreground-alt',
	paid: 'bg-success text-success-foreground',
	partially_paid: 'bg-tertiary/20 text-foreground',
	overdue: 'bg-destructive/10 text-destructive',
	refunded: 'bg-surface text-foreground-alt',
	void: 'bg-destructive/10 text-destructive'
};

export function paymentStatusBadge(status: PaymentStatus): string {
	return paymentStatusClasses[status];
}

// ---------------------------------------------------------------------------
// User role badges
// ---------------------------------------------------------------------------

const userRoleClasses: Record<UserRole, string> = {
	admin: 'bg-destructive/10 text-destructive',
	investigator: 'bg-accent text-accent-foreground',
	viewer: 'bg-surface text-foreground-alt'
};

export function userRoleBadge(role: UserRole): string {
	return userRoleClasses[role];
}

// ---------------------------------------------------------------------------
// Client type badges
// ---------------------------------------------------------------------------

const clientTypeClasses: Record<ClientType, string> = {
	person: 'bg-tertiary/15 text-tertiary',
	insurance: 'bg-accent text-accent-foreground',
	lawyer: 'bg-tertiary/15 text-tertiary',
	company: 'bg-success/15 text-success',
	government: 'bg-muted text-foreground',
};

export function clientTypeBadge(type: ClientType): string {
	return clientTypeClasses[type];
}

// ---------------------------------------------------------------------------
// Document workflow status dot colours
// ---------------------------------------------------------------------------

const documentStatusDotClasses: Record<DocumentStatus, string> = {
	draft: 'bg-muted-foreground',
	sent: 'bg-accent-foreground',
	signed: 'bg-success',
	active: 'bg-success',
	archived: 'bg-muted-foreground',
	cancelled: 'bg-destructive',
	rejected: 'bg-destructive',
	expired: 'bg-tertiary',
};

export function documentStatusDot(status: DocumentStatus): string {
	return documentStatusDotClasses[status];
}

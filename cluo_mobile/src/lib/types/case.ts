export interface Case {
	id: string;
	title: string;
	status: string; // "in_progress" | "ready" | "released" from backend
	externalReference?: string;
	clientId?: string;
	clientName?: string;
	clientNumber?: string;
}

export type CaseStatus = "in_progress" | "ready" | "released";

export interface Case {
	id: string;
	title: string;
	status: CaseStatus;
	externalReference?: string;
	clientId?: string;
	clientName?: string;
	clientNumber?: string;
}

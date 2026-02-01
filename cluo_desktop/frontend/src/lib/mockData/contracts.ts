/**
 * Mock contract data - Static, consistent across reloads.
 * Formal contracts for larger engagements. Not all cases require a contract.
 */

import { CONTRACT_SCOPES, PAYMENT_TERMS, CONFIDENTIALITY_CLAUSES, TERMINATION_CLAUSES, GOVERNING_LAWS, CURRENCY, DATES, type DocumentStatus } from './helpers';
import { CASE_IDS } from './cases';
import { CLIENT_IDS } from './clients';
import { MANDATE_IDS } from './mandates';

// UUIDs for contracts
export const CONTRACT_IDS = {
  // In progress cases (only one has a contract)
  inProgress2_contract1: 'c50e8400-e29b-41d4-a716-446655440102',

  // Ready cases (one has a contract)
  ready2_contract1: 'c50e8400-e29b-41d4-a716-446655440202',

  // Released cases (two have contracts)
  released4_contract1: 'c50e8400-e29b-41d4-a716-446655440304',
  released5_contract1: 'c50e8400-e29b-41d4-a716-446655440305',
} as const;

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

export const contracts: Contract[] = [
  // === IN PROGRESS CASES ===
  {
    id: CONTRACT_IDS.inProgress2_contract1,
    caseId: CASE_IDS.inProgress2,
    clientId: CLIENT_IDS.insurance1,
    contractNumber: 'CTR-2024-019',
    startDate: DATES.mar2024_25,
    endDate: DATES.mar2025_01,
    scopeOfServices: CONTRACT_SCOPES[1],
    paymentTerms: PAYMENT_TERMS[1],
    confidentiality: CONFIDENTIALITY_CLAUSES[1],
    terminationClause: TERMINATION_CLAUSES[0],
    signatures: [
      { name: 'Nathalie Fontaine', date: DATES.mar2024_25, role: 'Directrice Régionale AXA France' },
      { name: 'Sophie Martin', date: DATES.mar2024_25, role: 'Investigatrice Chef' }
    ],
    linkedMandateId: MANDATE_IDS.inProgress2_mand1,
    contractValue: 2345,
    currency: CURRENCY,
    renewalTerms: 'Renouvellement tacite pour une durée équivalente sauf préavis de 30 jours.',
    governingLaw: GOVERNING_LAWS[1],
    status: 'active',
    createdAt: DATES.mar2024_25,
    updatedAt: DATES.oct2024_25,
  },

  // === READY CASES ===
  {
    id: CONTRACT_IDS.ready2_contract1,
    caseId: CASE_IDS.ready2,
    clientId: CLIENT_IDS.company2,
    contractNumber: 'CTR-2024-050',
    startDate: DATES.apr2024_15,
    endDate: DATES.apr2025_01,
    scopeOfServices: CONTRACT_SCOPES[3],
    paymentTerms: PAYMENT_TERMS[2],
    confidentiality: CONFIDENTIALITY_CLAUSES[2],
    terminationClause: TERMINATION_CLAUSES[1],
    signatures: [
      { name: 'Olivier Michel', date: DATES.apr2024_15, role: 'Directeur des Opérations' },
      { name: 'Jean Pierre', date: DATES.apr2024_15, role: 'Investigateur Senior' }
    ],
    linkedMandateId: MANDATE_IDS.ready2_mand1,
    contractValue: 1030,
    currency: CURRENCY,
    renewalTerms: 'Non renouvelable - mission ponctuelle.',
    governingLaw: GOVERNING_LAWS[0],
    status: 'active',
    createdAt: DATES.apr2024_15,
    updatedAt: DATES.nov2024_22,
  },

  // === RELEASED CASES ===
  {
    id: CONTRACT_IDS.released4_contract1,
    caseId: CASE_IDS.released4,
    clientId: CLIENT_IDS.company1,
    contractNumber: 'CTR-2023-208',
    startDate: DATES.apr2024_25,
    endDate: DATES.nov2024_01,
    scopeOfServices: CONTRACT_SCOPES[1],
    paymentTerms: PAYMENT_TERMS[0],
    confidentiality: CONFIDENTIALITY_CLAUSES[0],
    terminationClause: TERMINATION_CLAUSES[2],
    signatures: [
      { name: 'Nicolas Richard', date: DATES.apr2024_25, role: 'Directeur Général' },
      { name: 'Sophie Martin', date: DATES.apr2024_25, role: 'Investigatrice Responsable' }
    ],
    linkedMandateId: MANDATE_IDS.released4_mand1,
    contractValue: 1760,
    currency: CURRENCY,
    renewalTerms: 'Non renouvelable - enquête terminée.',
    governingLaw: GOVERNING_LAWS[1],
    status: 'archived',
    createdAt: DATES.apr2024_25,
    updatedAt: DATES.nov2024_05,
  },
  {
    id: CONTRACT_IDS.released5_contract1,
    caseId: CASE_IDS.released5,
    clientId: CLIENT_IDS.lawyer1,
    contractNumber: 'CTR-2023-252',
    startDate: DATES.may2024_10,
    endDate: DATES.dec2024_01,
    scopeOfServices: CONTRACT_SCOPES[0],
    paymentTerms: PAYMENT_TERMS[1],
    confidentiality: CONFIDENTIALITY_CLAUSES[1],
    terminationClause: TERMINATION_CLAUSES[0],
    signatures: [
      { name: 'Jean-Pierre Aubry', date: DATES.may2024_10, role: 'Associé Fondateur' },
      { name: 'Jean Pierre', date: DATES.may2024_10, role: 'Investigateur Principal' }
    ],
    linkedMandateId: MANDATE_IDS.released5_mand1,
    contractValue: 1030,
    currency: CURRENCY,
    renewalTerms: 'Renouvellement possible selon avancement du dossier.',
    governingLaw: GOVERNING_LAWS[1],
    status: 'archived',
    createdAt: DATES.may2024_10,
    updatedAt: DATES.dec2024_10,
  },
];

// Helper functions
export function getAllContracts(): Contract[] {
  return contracts;
}

export function getContractById(id: string): Contract | undefined {
  return contracts.find(c => c.id === id);
}

export function getContractsByCaseId(caseId: string): Contract[] {
  return contracts.filter(c => c.caseId === caseId);
}

export function getContractsByClientId(clientId: string): Contract[] {
  return contracts.filter(c => c.clientId === clientId);
}

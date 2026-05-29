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
  // In progress cases
  inProgress1_contract1: 'c50e8400-e29b-41d4-a716-446655440101',
  inProgress2_contract1: 'c50e8400-e29b-41d4-a716-446655440102',
  inProgress3_contract1: 'c50e8400-e29b-41d4-a716-446655440103',
  inProgress4_contract1: 'c50e8400-e29b-41d4-a716-446655440104',
  inProgress5_contract1: 'c50e8400-e29b-41d4-a716-446655440105',

  // Ready cases
  ready1_contract1: 'c50e8400-e29b-41d4-a716-446655440201',
  ready2_contract1: 'c50e8400-e29b-41d4-a716-446655440202',
  ready3_contract1: 'c50e8400-e29b-41d4-a716-446655440203',

  // Released cases
  released1_contract1: 'c50e8400-e29b-41d4-a716-446655440301',
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
    id: CONTRACT_IDS.inProgress1_contract1,
    caseId: CASE_IDS.inProgress1,
    clientId: CLIENT_IDS.person1,
    contractNumber: 'CTR-2024-009',
    startDate: DATES.jan2024_30,
    endDate: DATES.jul2025_01,
    scopeOfServices: CONTRACT_SCOPES[0],
    paymentTerms: PAYMENT_TERMS[0],
    confidentiality: CONFIDENTIALITY_CLAUSES[0],
    terminationClause: TERMINATION_CLAUSES[0],
    signatures: [
      { name: 'Jean Dupont', date: DATES.jan2024_30, role: 'Client' },
      { name: 'Sophie Martin', date: DATES.jan2024_30, role: 'Investigatrice Chef' }
    ],
    linkedMandateId: MANDATE_IDS.inProgress1_mand1,
    contractValue: 1850,
    currency: CURRENCY,
    renewalTerms: 'Renouvellement tacite pour une durée équivalente sauf préavis de 30 jours.',
    governingLaw: GOVERNING_LAWS[0],
    status: 'active',
    createdAt: DATES.jan2024_30,
    updatedAt: DATES.oct2024_25,
  },
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
  {
    id: CONTRACT_IDS.inProgress3_contract1,
    caseId: CASE_IDS.inProgress3,
    clientId: CLIENT_IDS.lawyer1,
    contractNumber: 'CTR-2024-033',
    startDate: DATES.may2024_20,
    endDate: DATES.nov2025_01,
    scopeOfServices: CONTRACT_SCOPES[2],
    paymentTerms: PAYMENT_TERMS[1],
    confidentiality: CONFIDENTIALITY_CLAUSES[1],
    terminationClause: TERMINATION_CLAUSES[1],
    signatures: [
      { name: 'Jean-Pierre Aubry', date: DATES.may2024_20, role: 'Associé Fondateur' },
      { name: 'Sophie Martin', date: DATES.may2024_20, role: 'Investigatrice Responsable' }
    ],
    linkedMandateId: MANDATE_IDS.inProgress3_mand1,
    contractValue: 2100,
    currency: CURRENCY,
    renewalTerms: 'Renouvellement possible selon avancement de la recherche.',
    governingLaw: GOVERNING_LAWS[2],
    status: 'active',
    createdAt: DATES.may2024_20,
    updatedAt: DATES.oct2024_25,
  },
  {
    id: CONTRACT_IDS.inProgress4_contract1,
    caseId: CASE_IDS.inProgress4,
    clientId: CLIENT_IDS.insurance2,
    contractNumber: 'CTR-2024-046',
    startDate: DATES.jul2024_15,
    endDate: DATES.jan2025_15,
    scopeOfServices: CONTRACT_SCOPES[1],
    paymentTerms: PAYMENT_TERMS[0],
    confidentiality: CONFIDENTIALITY_CLAUSES[2],
    terminationClause: TERMINATION_CLAUSES[2],
    signatures: [
      { name: 'Christophe Bertrand', date: DATES.jul2024_15, role: 'Directeur Sinistres' },
      { name: 'Jean Pierre', date: DATES.jul2024_15, role: 'Investigateur Senior' }
    ],
    linkedMandateId: MANDATE_IDS.inProgress4_mand1,
    contractValue: 1680,
    currency: CURRENCY,
    renewalTerms: 'Non renouvelable - mission ponctuelle.',
    governingLaw: GOVERNING_LAWS[1],
    status: 'active',
    createdAt: DATES.jul2024_15,
    updatedAt: DATES.nov2024_22,
  },

  {
    id: CONTRACT_IDS.inProgress5_contract1,
    caseId: CASE_IDS.inProgress5,
    clientId: CLIENT_IDS.company1,
    contractNumber: 'CTR-2025-001',
    startDate: DATES.jan2024_30,
    endDate: '',
    scopeOfServices: CONTRACT_SCOPES[0],
    paymentTerms: PAYMENT_TERMS[0],
    confidentiality: CONFIDENTIALITY_CLAUSES[0],
    terminationClause: TERMINATION_CLAUSES[0],
    signatures: [],
    linkedMandateId: MANDATE_IDS.inProgress5_mand1,
    contractValue: 1500,
    currency: CURRENCY,
    renewalTerms: '',
    governingLaw: GOVERNING_LAWS[0],
    status: 'draft',
    createdAt: DATES.jan2024_30,
    updatedAt: DATES.jan2024_30,
  },

  // === READY CASES ===
  {
    id: CONTRACT_IDS.ready1_contract1,
    caseId: CASE_IDS.ready1,
    clientId: CLIENT_IDS.person2,
    contractNumber: 'CTR-2024-026',
    startDate: DATES.feb2024_20,
    endDate: DATES.aug2025_01,
    scopeOfServices: CONTRACT_SCOPES[4],
    paymentTerms: PAYMENT_TERMS[0],
    confidentiality: CONFIDENTIALITY_CLAUSES[0],
    terminationClause: TERMINATION_CLAUSES[0],
    signatures: [
      { name: 'Marie Leroy', date: DATES.feb2024_20, role: 'Cliente' },
      { name: 'Sophie Martin', date: DATES.feb2024_20, role: 'Investigatrice Chef' }
    ],
    linkedMandateId: MANDATE_IDS.ready1_mand1,
    contractValue: 1260,
    currency: CURRENCY,
    renewalTerms: 'Renouvellement tacite pour une durée équivalente.',
    governingLaw: GOVERNING_LAWS[0],
    status: 'active',
    createdAt: DATES.feb2024_20,
    updatedAt: DATES.nov2024_05,
  },
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
  {
    id: CONTRACT_IDS.ready3_contract1,
    caseId: CASE_IDS.ready3,
    clientId: CLIENT_IDS.lawyer2,
    contractNumber: 'CTR-2024-069',
    startDate: DATES.jun2024_25,
    endDate: DATES.dec2025_01,
    scopeOfServices: CONTRACT_SCOPES[3],
    paymentTerms: PAYMENT_TERMS[1],
    confidentiality: CONFIDENTIALITY_CLAUSES[1],
    terminationClause: TERMINATION_CLAUSES[1],
    signatures: [
      { name: 'Isabelle Fournier', date: DATES.jun2024_25, role: 'Avocate' },
      { name: 'Sophie Martin', date: DATES.jun2024_25, role: 'Investigatrice Responsable' }
    ],
    linkedMandateId: MANDATE_IDS.ready3_mand1,
    contractValue: 830,
    currency: CURRENCY,
    renewalTerms: 'Non renouvelable - enquête ponctuelle.',
    governingLaw: GOVERNING_LAWS[0],
    status: 'active',
    createdAt: DATES.jun2024_25,
    updatedAt: DATES.dec2024_10,
  },

  // === RELEASED CASES ===
  {
    id: CONTRACT_IDS.released1_contract1,
    caseId: CASE_IDS.released1,
    clientId: CLIENT_IDS.government1,
    contractNumber: 'CTR-2023-096',
    startDate: DATES.jan2024_30,
    endDate: DATES.jul2024_01,
    scopeOfServices: CONTRACT_SCOPES[3],
    paymentTerms: PAYMENT_TERMS[0],
    confidentiality: CONFIDENTIALITY_CLAUSES[2],
    terminationClause: TERMINATION_CLAUSES[0],
    signatures: [
      { name: 'Alain Blanc', date: DATES.jan2024_30, role: 'Directeur RH' },
      { name: 'Jean Pierre', date: DATES.jan2024_30, role: 'Investigateur Principal' }
    ],
    linkedMandateId: MANDATE_IDS.released1_mand1,
    contractValue: 1300,
    currency: CURRENCY,
    renewalTerms: 'Non renouvelable - enquête terminée.',
    governingLaw: GOVERNING_LAWS[1],
    status: 'archived',
    createdAt: DATES.jan2024_30,
    updatedAt: DATES.aug2024_18,
  },
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

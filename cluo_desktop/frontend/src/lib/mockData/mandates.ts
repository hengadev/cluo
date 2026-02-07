/**
 * Mock mandate data - Static, consistent across reloads.
 * Each case typically has one mandate (authorization to investigate).
 * Some cases may be in early stages without a mandate yet.
 */

import { MANDATE_SCOPES, GOVERNING_LAWS, JURISDICTIONS, SPECIAL_INSTRUCTIONS, DATES, type DocumentStatus } from './helpers';
import { CASE_IDS } from './cases';
import { CLIENT_IDS } from './clients';
import { ESTIMATE_IDS } from './estimates';

// UUIDs for mandates
export const MANDATE_IDS = {
  // In progress cases (most have mandates, one pending)
  inProgress1_mand1: 'b50e8400-e29b-41d4-a716-446655440101',
  inProgress2_mand1: 'b50e8400-e29b-41d4-a716-446655440102',
  inProgress3_mand1: 'b50e8400-e29b-41d4-a716-446655440103',
  inProgress4_mand1: 'b50e8400-e29b-41d4-a716-446655440104',
  // inProgress5 - no mandate yet (estimate only)

  // Ready cases (all have mandates)
  ready1_mand1: 'b50e8400-e29b-41d4-a716-446655440201',
  ready2_mand1: 'b50e8400-e29b-41d4-a716-446655440202',
  ready3_mand1: 'b50e8400-e29b-41d4-a716-446655440203',

  // Released cases (all have mandates)
  released1_mand1: 'b50e8400-e29b-41d4-a716-446655440301',
  released2_mand1: 'b50e8400-e29b-41d4-a716-446655440302',
  released3_mand1: 'b50e8400-e29b-41d4-a716-446655440303',
  released4_mand1: 'b50e8400-e29b-41d4-a716-446655440304',
  released5_mand1: 'b50e8400-e29b-41d4-a716-446655440305',
} as const;

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

export const mandates: Mandate[] = [
  // === IN PROGRESS CASES ===
  {
    id: MANDATE_IDS.inProgress1_mand1,
    caseId: CASE_IDS.inProgress1,
    clientId: CLIENT_IDS.person1,
    mandateNumber: 'MAN-2024-008',
    issueDate: DATES.jan2024_25,
    scopeOfWork: MANDATE_SCOPES[0],
    validFrom: DATES.jan2024_25,
    validUntil: DATES.jul2025_01,
    termsConditions: 'Prestation soumise aux conditions générales de la profession. Honoraires selon devis DEV-2024-008.',
    clientSignature: { name: 'Jean Dupont', date: DATES.jan2024_25 },
    investigatorSignature: { name: 'Sophie Martin', date: DATES.jan2024_25 },
    linkedEstimateId: ESTIMATE_IDS.inProgress1_est1,
    specialInstructions: null,
    jurisdiction: JURISDICTIONS[0],
    status: 'active',
    createdAt: DATES.jan2024_25,
    updatedAt: DATES.oct2024_25,
  },
  {
    id: MANDATE_IDS.inProgress2_mand1,
    caseId: CASE_IDS.inProgress2,
    clientId: CLIENT_IDS.insurance1,
    mandateNumber: 'MAN-2024-019',
    issueDate: DATES.mar2024_25,
    scopeOfWork: MANDATE_SCOPES[1],
    validFrom: DATES.mar2024_25,
    validUntil: DATES.sep2025_01,
    termsConditions: 'Mission de surveillance dans le cadre du contrat d\'assurance. Rapports hebdomadaires requis.',
    clientSignature: { name: 'Nathalie Fontaine', date: DATES.mar2024_25 },
    investigatorSignature: { name: 'Jean Pierre', date: DATES.mar2024_25 },
    linkedEstimateId: ESTIMATE_IDS.inProgress2_est1,
    specialInstructions: SPECIAL_INSTRUCTIONS[1],
    jurisdiction: JURISDICTIONS[1],
    status: 'active',
    createdAt: DATES.mar2024_25,
    updatedAt: DATES.oct2024_25,
  },
  {
    id: MANDATE_IDS.inProgress3_mand1,
    caseId: CASE_IDS.inProgress3,
    clientId: CLIENT_IDS.lawyer1,
    mandateNumber: 'MAN-2024-032',
    issueDate: DATES.may2024_20,
    scopeOfWork: MANDATE_SCOPES[3],
    validFrom: DATES.may2024_20,
    validUntil: DATES.nov2025_01,
    termsConditions: 'Recherche de personnes avec investigations terrain et rapports périodiques.',
    clientSignature: { name: 'Jean-Pierre Aubry', date: DATES.may2024_20 },
    investigatorSignature: { name: 'Sophie Martin', date: DATES.may2024_20 },
    linkedEstimateId: ESTIMATE_IDS.inProgress3_est1,
    specialInstructions: null,
    jurisdiction: JURISDICTIONS[2],
    status: 'active',
    createdAt: DATES.may2024_20,
    updatedAt: DATES.oct2024_25,
  },
  {
    id: MANDATE_IDS.inProgress4_mand1,
    caseId: CASE_IDS.inProgress4,
    clientId: CLIENT_IDS.insurance2,
    mandateNumber: 'MAN-2024-045',
    issueDate: DATES.jul2024_15,
    scopeOfWork: MANDATE_SCOPES[4],
    validFrom: DATES.jul2024_15,
    validUntil: DATES.jan2025_01,
    termsConditions: 'Enquête fraude présumée. Rapport détaillé avec conclusions attendu.',
    clientSignature: { name: 'Christophe Bertrand', date: DATES.jul2024_15 },
    investigatorSignature: { name: 'Jean Pierre', date: DATES.jul2024_15 },
    linkedEstimateId: ESTIMATE_IDS.inProgress4_est1,
    specialInstructions: SPECIAL_INSTRUCTIONS[0],
    jurisdiction: JURISDICTIONS[0],
    status: 'active',
    createdAt: DATES.jul2024_15,
    updatedAt: DATES.nov2024_22,
  },

  // === READY CASES ===
  {
    id: MANDATE_IDS.ready1_mand1,
    caseId: CASE_IDS.ready1,
    clientId: CLIENT_IDS.person2,
    mandateNumber: 'MAN-2024-025',
    issueDate: DATES.feb2024_20,
    scopeOfWork: MANDATE_SCOPES[5],
    validFrom: DATES.feb2024_20,
    validUntil: DATES.aug2025_01,
    termsConditions: 'Enquête conjointe avec médiation familiale possible.',
    clientSignature: { name: 'Marie Leroy', date: DATES.feb2024_20 },
    investigatorSignature: { name: 'Sophie Martin', date: DATES.feb2024_20 },
    linkedEstimateId: ESTIMATE_IDS.ready1_est1,
    specialInstructions: SPECIAL_INSTRUCTIONS[5],
    jurisdiction: JURISDICTIONS[0],
    status: 'active',
    createdAt: DATES.feb2024_20,
    updatedAt: DATES.nov2024_05,
  },
  {
    id: MANDATE_IDS.ready2_mand1,
    caseId: CASE_IDS.ready2,
    clientId: CLIENT_IDS.company2,
    mandateNumber: 'MAN-2024-050',
    issueDate: DATES.apr2024_15,
    scopeOfWork: MANDATE_SCOPES[3],
    validFrom: DATES.apr2024_15,
    validUntil: DATES.oct2025_01,
    termsConditions: 'Enquête interne avec confidentialité renforcée. Rapport final synthétique.',
    clientSignature: { name: 'Olivier Michel', date: DATES.apr2024_15 },
    investigatorSignature: { name: 'Jean Pierre', date: DATES.apr2024_15 },
    linkedEstimateId: ESTIMATE_IDS.ready2_est1,
    specialInstructions: SPECIAL_INSTRUCTIONS[3],
    jurisdiction: JURISDICTIONS[1],
    status: 'active',
    createdAt: DATES.apr2024_15,
    updatedAt: DATES.nov2024_22,
  },
  {
    id: MANDATE_IDS.ready3_mand1,
    caseId: CASE_IDS.ready3,
    clientId: CLIENT_IDS.lawyer2,
    mandateNumber: 'MAN-2024-068',
    issueDate: DATES.jun2024_25,
    scopeOfWork: MANDATE_SCOPES[6],
    validFrom: DATES.jun2024_25,
    validUntil: DATES.dec2025_01,
    termsConditions: 'Investigation vol propriété intellectuelle avec preuves recevables en justice.',
    clientSignature: { name: 'Isabelle Fournier', date: DATES.jun2024_25 },
    investigatorSignature: { name: 'Sophie Martin', date: DATES.jun2024_25 },
    linkedEstimateId: ESTIMATE_IDS.ready3_est1,
    specialInstructions: SPECIAL_INSTRUCTIONS[2],
    jurisdiction: JURISDICTIONS[0],
    status: 'active',
    createdAt: DATES.jun2024_25,
    updatedAt: DATES.dec2024_10,
  },

  // === RELEASED CASES ===
  {
    id: MANDATE_IDS.released1_mand1,
    caseId: CASE_IDS.released1,
    clientId: CLIENT_IDS.government1,
    mandateNumber: 'MAN-2023-095',
    issueDate: DATES.jan2024_30,
    scopeOfWork: MANDATE_SCOPES[3],
    validFrom: DATES.jan2024_30,
    validUntil: DATES.jul2024_01,
    termsConditions: 'Enquête harcèlement au travail selon protocole administratif.',
    clientSignature: { name: 'Alain Blanc', date: DATES.jan2024_30 },
    investigatorSignature: { name: 'Jean Pierre', date: DATES.jan2024_30 },
    linkedEstimateId: ESTIMATE_IDS.released1_est1,
    specialInstructions: SPECIAL_INSTRUCTIONS[3],
    jurisdiction: JURISDICTIONS[0],
    status: 'archived',
    createdAt: DATES.jan2024_30,
    updatedAt: DATES.aug2024_18,
  },
  {
    id: MANDATE_IDS.released2_mand1,
    caseId: CASE_IDS.released2,
    clientId: CLIENT_IDS.person3,
    mandateNumber: 'MAN-2023-118',
    issueDate: DATES.mar2024_05,
    scopeOfWork: MANDATE_SCOPES[1],
    validFrom: DATES.mar2024_05,
    validUntil: DATES.sep2024_01,
    termsConditions: 'Surveillance et documentation preuves accident de la route.',
    clientSignature: { name: 'Pierre Moreau', date: DATES.mar2024_05 },
    investigatorSignature: { name: 'Sophie Martin', date: DATES.mar2024_05 },
    linkedEstimateId: ESTIMATE_IDS.released2_est1,
    specialInstructions: null,
    jurisdiction: JURISDICTIONS[0],
    status: 'archived',
    createdAt: DATES.mar2024_05,
    updatedAt: DATES.sep2024_01,
  },
  {
    id: MANDATE_IDS.released3_mand1,
    caseId: CASE_IDS.released3,
    clientId: CLIENT_IDS.insurance1,
    mandateNumber: 'MAN-2023-162',
    issueDate: DATES.mar2024_10,
    scopeOfWork: MANDATE_SCOPES[0],
    validFrom: DATES.mar2024_10,
    validUntil: DATES.sep2024_01,
    termsConditions: 'Enquête filiation avec vérification liens familiaux.',
    clientSignature: { name: 'Vincent Roux', date: DATES.mar2024_10 },
    investigatorSignature: { name: 'Jean Pierre', date: DATES.mar2024_10 },
    linkedEstimateId: ESTIMATE_IDS.released3_est1,
    specialInstructions: null,
    jurisdiction: JURISDICTIONS[0],
    status: 'archived',
    createdAt: DATES.mar2024_10,
    updatedAt: DATES.oct2024_08,
  },
  {
    id: MANDATE_IDS.released4_mand1,
    caseId: CASE_IDS.released4,
    clientId: CLIENT_IDS.company1,
    mandateNumber: 'MAN-2023-208',
    issueDate: DATES.apr2024_25,
    scopeOfWork: MANDATE_SCOPES[1],
    validFrom: DATES.apr2024_25,
    validUntil: DATES.oct2024_01,
    termsConditions: 'Surveillance employee avec rapports périodiques et documentation photo.',
    clientSignature: { name: 'Nicolas Richard', date: DATES.apr2024_25 },
    investigatorSignature: { name: 'Sophie Martin', date: DATES.apr2024_25 },
    linkedEstimateId: ESTIMATE_IDS.released4_est1,
    specialInstructions: SPECIAL_INSTRUCTIONS[1],
    jurisdiction: JURISDICTIONS[1],
    status: 'archived',
    createdAt: DATES.apr2024_25,
    updatedAt: DATES.nov2024_05,
  },
  {
    id: MANDATE_IDS.released5_mand1,
    caseId: CASE_IDS.released5,
    clientId: CLIENT_IDS.lawyer1,
    mandateNumber: 'MAN-2023-252',
    issueDate: DATES.may2024_10,
    scopeOfWork: MANDATE_SCOPES[0],
    validFrom: DATES.may2024_10,
    validUntil: DATES.nov2024_01,
    termsConditions: 'Enquête successorale complet avec constitution dossier juridique.',
    clientSignature: { name: 'Jean-Pierre Aubry', date: DATES.may2024_10 },
    investigatorSignature: { name: 'Jean Pierre', date: DATES.may2024_10 },
    linkedEstimateId: ESTIMATE_IDS.released5_est1,
    specialInstructions: SPECIAL_INSTRUCTIONS[5],
    jurisdiction: JURISDICTIONS[0],
    status: 'archived',
    createdAt: DATES.may2024_10,
    updatedAt: DATES.dec2024_10,
  },
];

// Helper functions
export function getAllMandates(): Mandate[] {
  return mandates;
}

export function getMandateById(id: string): Mandate | undefined {
  return mandates.find(m => m.id === id);
}

export function getMandatesByCaseId(caseId: string): Mandate[] {
  return mandates.filter(m => m.caseId === caseId);
}

export function getMandatesByClientId(clientId: string): Mandate[] {
  return mandates.filter(m => m.clientId === clientId);
}

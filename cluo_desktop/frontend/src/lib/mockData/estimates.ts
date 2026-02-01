/**
 * Mock estimate data - Static, consistent across reloads.
 * Most cases start with an estimate. Some may have revised estimates.
 */

import { SERVICE_LINE_ITEMS, PAYMENT_TERMS, CURRENCY, DATES, type DocumentStatus } from './helpers';
import { CASE_IDS } from './cases';
import { CLIENT_IDS } from './clients';

// Re-export DocumentStatus for convenience
export type { DocumentStatus };

// UUIDs for estimates
export const ESTIMATE_IDS = {
  // For draft case (estimate created but not sent yet)
  draft1_est1: 'a50e8400-e29b-41d4-a716-446655440001',

  // For in progress cases
  inProgress1_est1: 'a50e8400-e29b-41d4-a716-446655440101',
  inProgress2_est1: 'a50e8400-e29b-41d4-a716-446655440102',
  inProgress3_est1: 'a50e8400-e29b-41d4-a716-446655440103',
  inProgress4_est1: 'a50e8400-e29b-41d4-a716-446655440104',
  inProgress4_est2: 'a50e8400-e29b-41d4-a716-446655440105', // Revised estimate
  inProgress5_est1: 'a50e8400-e29b-41d4-a716-446655440106',

  // For ready cases
  ready1_est1: 'a50e8400-e29b-41d4-a716-446655440201',
  ready2_est1: 'a50e8400-e29b-41d4-a716-446655440202',
  ready3_est1: 'a50e8400-e29b-41d4-a716-446655440203',

  // For released cases (all accepted, most have invoices)
  released1_est1: 'a50e8400-e29b-41d4-a716-446655440301',
  released2_est1: 'a50e8400-e29b-41d4-a716-446655440302',
  released2_est2: 'a50e8400-e29b-41d4-a716-446655440303', // Revised
  released3_est1: 'a50e8400-e29b-41d4-a716-446655440304',
  released4_est1: 'a50e8400-e29b-41d4-a716-446655440305',
  released5_est1: 'a50e8400-e29b-41d4-a716-446655440306',
} as const;

export interface EstimateLineItem {
  description: string;
  quantity: number;
  unitPrice: number;
  total: number;
}

export interface Estimate {
  id: string;
  caseId: string;
  clientId: string;
  estimateNumber: string;
  issueDate: string;
  validUntil: string;
  lineItems: EstimateLineItem[];
  estimatedTotal: number;
  notes: string;
  accepted: boolean;
  acceptedAt: string | null;
  acceptedBy: string | null;
  status: DocumentStatus;
  createdAt: string;
  updatedAt: string;
}

export const estimates: Estimate[] = [
  // === DRAFT CASE ===
  {
    id: ESTIMATE_IDS.draft1_est1,
    caseId: CASE_IDS.draft1,
    clientId: CLIENT_IDS.person1,
    estimateNumber: 'DEV-2024-089',
    issueDate: DATES.oct2024_25,
    validUntil: DATES.nov2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[0].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[0].unitPrice, total: 450 },
      { description: SERVICE_LINE_ITEMS[1].description, quantity: 2, unitPrice: SERVICE_LINE_ITEMS[1].unitPrice, total: 760 },
    ],
    estimatedTotal: 1210,
    notes: 'Devis préliminaire, sujet à modification selon l\'avancée de l\'enquête.',
    accepted: false,
    acceptedAt: null,
    acceptedBy: null,
    status: 'draft',
    createdAt: DATES.oct2024_25,
    updatedAt: DATES.oct2024_25,
  },

  // === IN PROGRESS CASES ===
  {
    id: ESTIMATE_IDS.inProgress1_est1,
    caseId: CASE_IDS.inProgress1,
    clientId: CLIENT_IDS.person1,
    estimateNumber: 'DEV-2024-008',
    issueDate: DATES.jan2024_10,
    validUntil: DATES.feb2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[8].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[8].unitPrice, total: 320 },
      { description: SERVICE_LINE_ITEMS[5].description, quantity: 2, unitPrice: SERVICE_LINE_ITEMS[5].unitPrice, total: 360 },
      { description: SERVICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[4].unitPrice, total: 220 },
    ],
    estimatedTotal: 900,
    notes: 'Enquête de filiation complète avec recherches administratives.',
    accepted: true,
    acceptedAt: DATES.jan2024_25,
    acceptedBy: 'Jean Dupont',
    status: 'signed',
    createdAt: DATES.jan2024_10,
    updatedAt: DATES.jan2024_25,
  },
  {
    id: ESTIMATE_IDS.inProgress2_est1,
    caseId: CASE_IDS.inProgress2,
    clientId: CLIENT_IDS.insurance1,
    estimateNumber: 'DEV-2024-019',
    issueDate: DATES.mar2024_22,
    validUntil: DATES.apr2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[1].description, quantity: 5, unitPrice: SERVICE_LINE_ITEMS[1].unitPrice, total: 1900 },
      { description: SERVICE_LINE_ITEMS[6].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[6].unitPrice, total: 120 },
      { description: SERVICE_LINE_ITEMS[7].description, quantity: 5, unitPrice: SERVICE_LINE_ITEMS[7].unitPrice, total: 325 },
    ],
    estimatedTotal: 2345,
    notes: 'Surveillance prolongée avec documentation photographique complète.',
    accepted: true,
    acceptedAt: DATES.mar2024_25,
    acceptedBy: 'Nathalie Fontaine',
    status: 'active',
    createdAt: DATES.mar2024_22,
    updatedAt: DATES.mar2024_25,
  },
  {
    id: ESTIMATE_IDS.inProgress3_est1,
    caseId: CASE_IDS.inProgress3,
    clientId: CLIENT_IDS.lawyer1,
    estimateNumber: 'DEV-2024-032',
    issueDate: DATES.may2024_15,
    validUntil: DATES.jun2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[2].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[2].unitPrice, total: 280 },
      { description: SERVICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[4].unitPrice, total: 220 },
      { description: SERVICE_LINE_ITEMS[5].description, quantity: 3, unitPrice: SERVICE_LINE_ITEMS[5].unitPrice, total: 540 },
    ],
    estimatedTotal: 1040,
    notes: 'Recheche de personne avec enquête terrain et investigations multiples.',
    accepted: true,
    acceptedAt: DATES.may2024_20,
    acceptedBy: 'Jean-Pierre Aubry',
    status: 'active',
    createdAt: DATES.may2024_15,
    updatedAt: DATES.may2024_20,
  },
  {
    id: ESTIMATE_IDS.inProgress4_est1,
    caseId: CASE_IDS.inProgress4,
    clientId: CLIENT_IDS.insurance2,
    estimateNumber: 'DEV-2024-045',
    issueDate: DATES.jul2024_10,
    validUntil: DATES.jul2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[0].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[0].unitPrice, total: 450 },
      { description: SERVICE_LINE_ITEMS[5].description, quantity: 2, unitPrice: SERVICE_LINE_ITEMS[5].unitPrice, total: 360 },
    ],
    estimatedTotal: 810,
    notes: 'Enquête fraude présumée - première phase.',
    accepted: true,
    acceptedAt: DATES.jul2024_15,
    acceptedBy: 'Christophe Bertrand',
    status: 'active',
    createdAt: DATES.jul2024_10,
    updatedAt: DATES.jul2024_15,
  },
  {
    id: ESTIMATE_IDS.inProgress4_est2,
    caseId: CASE_IDS.inProgress4,
    clientId: CLIENT_IDS.insurance2,
    estimateNumber: 'DEV-2024-045-BIS',
    issueDate: DATES.aug2024_05,
    validUntil: DATES.sep2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[10].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[10].unitPrice, total: 195 },
      { description: SERVICE_LINE_ITEMS[11].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[11].unitPrice, total: 250 },
      { description: SERVICE_LINE_ITEMS[7].description, quantity: 2, unitPrice: SERVICE_LINE_ITEMS[7].unitPrice, total: 130 },
    ],
    estimatedTotal: 575,
    notes: 'Devis complémentaire - investigations techniques supplémentaires.',
    accepted: true,
    acceptedAt: DATES.aug2024_10,
    acceptedBy: 'Christophe Bertrand',
    status: 'active',
    createdAt: DATES.aug2024_05,
    updatedAt: DATES.aug2024_10,
  },
  {
    id: ESTIMATE_IDS.inProgress5_est1,
    caseId: CASE_IDS.inProgress5,
    clientId: CLIENT_IDS.company1,
    estimateNumber: 'DEV-2024-059',
    issueDate: DATES.sep2024_12,
    validUntil: DATES.oct2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[2].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[2].unitPrice, total: 280 },
      { description: SERVICE_LINE_ITEMS[5].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[5].unitPrice, total: 180 },
      { description: SERVICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[4].unitPrice, total: 220 },
    ],
    estimatedTotal: 680,
    notes: 'Vérification approfondie des antécédents professionnels et diplômes.',
    accepted: true,
    acceptedAt: DATES.sep2024_18,
    acceptedBy: 'Nicolas Richard',
    status: 'active',
    createdAt: DATES.sep2024_12,
    updatedAt: DATES.sep2024_18,
  },

  // === READY CASES ===
  {
    id: ESTIMATE_IDS.ready1_est1,
    caseId: CASE_IDS.ready1,
    clientId: CLIENT_IDS.person2,
    estimateNumber: 'DEV-2024-025',
    issueDate: DATES.feb2024_14,
    validUntil: DATES.mar2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[1].description, quantity: 3, unitPrice: SERVICE_LINE_ITEMS[1].unitPrice, total: 1140 },
      { description: SERVICE_LINE_ITEMS[6].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[6].unitPrice, total: 120 },
    ],
    estimatedTotal: 1260,
    notes: 'Surveillance conjointe avec preuves photographiques.',
    accepted: true,
    acceptedAt: DATES.feb2024_20,
    acceptedBy: 'Marie Leroy',
    status: 'active',
    createdAt: DATES.feb2024_14,
    updatedAt: DATES.feb2024_20,
  },
  {
    id: ESTIMATE_IDS.ready2_est1,
    caseId: CASE_IDS.ready2,
    clientId: CLIENT_IDS.company2,
    estimateNumber: 'DEV-2024-050',
    issueDate: DATES.apr2024_08,
    validUntil: DATES.may2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[0].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[0].unitPrice, total: 450 },
      { description: SERVICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[4].unitPrice, total: 220 },
      { description: SERVICE_LINE_ITEMS[5].description, quantity: 2, unitPrice: SERVICE_LINE_ITEMS[5].unitPrice, total: 360 },
    ],
    estimatedTotal: 1030,
    notes: 'Enquête interne harcèlement avec interviews témoins.',
    accepted: true,
    acceptedAt: DATES.apr2024_15,
    acceptedBy: 'Olivier Michel',
    status: 'active',
    createdAt: DATES.apr2024_08,
    updatedAt: DATES.apr2024_15,
  },
  {
    id: ESTIMATE_IDS.ready3_est1,
    caseId: CASE_IDS.ready3,
    clientId: CLIENT_IDS.lawyer2,
    estimateNumber: 'DEV-2024-068',
    issueDate: DATES.jun2024_20,
    validUntil: DATES.jul2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[5].description, quantity: 2, unitPrice: SERVICE_LINE_ITEMS[5].unitPrice, total: 360 },
      { description: SERVICE_LINE_ITEMS[11].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[11].unitPrice, total: 250 },
      { description: SERVICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[4].unitPrice, total: 220 },
    ],
    estimatedTotal: 830,
    notes: 'Investigation vol propriété intellectuelle avec analyse numérique.',
    accepted: true,
    acceptedAt: DATES.jun2024_25,
    acceptedBy: 'Isabelle Fournier',
    status: 'active',
    createdAt: DATES.jun2024_20,
    updatedAt: DATES.jun2024_25,
  },

  // === RELEASED CASES ===
  {
    id: ESTIMATE_IDS.released1_est1,
    caseId: CASE_IDS.released1,
    clientId: CLIENT_IDS.government1,
    estimateNumber: 'DEV-2023-095',
    issueDate: DATES.jan2024_25,
    validUntil: DATES.feb2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[0].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[0].unitPrice, total: 450 },
      { description: SERVICE_LINE_ITEMS[3].description, quantity: 3, unitPrice: SERVICE_LINE_ITEMS[3].unitPrice, total: 450 },
      { description: SERVICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[4].unitPrice, total: 220 },
    ],
    estimatedTotal: 1120,
    notes: 'Enquête harcèlement au travail avec témoignages collègues.',
    accepted: true,
    acceptedAt: DATES.jan2024_30,
    acceptedBy: 'Alain Blanc',
    status: 'archived',
    createdAt: DATES.jan2024_25,
    updatedAt: DATES.aug2024_18,
  },
  {
    id: ESTIMATE_IDS.released2_est1,
    caseId: CASE_IDS.released2,
    clientId: CLIENT_IDS.person3,
    estimateNumber: 'DEV-2023-118',
    issueDate: DATES.feb2024_28,
    validUntil: DATES.mar2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[1].description, quantity: 2, unitPrice: SERVICE_LINE_ITEMS[1].unitPrice, total: 760 },
      { description: SERVICE_LINE_ITEMS[6].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[6].unitPrice, total: 120 },
      { description: SERVICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[4].unitPrice, total: 220 },
    ],
    estimatedTotal: 1100,
    notes: 'Surveillance et collecte de preuves accident.',
    accepted: true,
    acceptedAt: DATES.mar2024_05,
    acceptedBy: 'Pierre Moreau',
    status: 'archived',
    createdAt: DATES.feb2024_28,
    updatedAt: DATES.sep2024_01,
  },
  {
    id: ESTIMATE_IDS.released2_est2,
    caseId: CASE_IDS.released2,
    clientId: CLIENT_IDS.person3,
    estimateNumber: 'DEV-2023-118-BIS',
    issueDate: DATES.apr2024_18,
    validUntil: DATES.may2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[7].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[7].unitPrice, total: 120 },
    ],
    estimatedTotal: 120,
    notes: 'Devis complémentaire - photographs additionnelles.',
    accepted: true,
    acceptedAt: DATES.apr2024_22,
    acceptedBy: 'Pierre Moreau',
    status: 'archived',
    createdAt: DATES.apr2024_18,
    updatedAt: DATES.sep2024_01,
  },
  {
    id: ESTIMATE_IDS.released3_est1,
    caseId: CASE_IDS.released3,
    clientId: CLIENT_IDS.insurance1,
    estimateNumber: 'DEV-2023-162',
    issueDate: DATES.mar2024_05,
    validUntil: DATES.apr2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[8].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[8].unitPrice, total: 320 },
      { description: SERVICE_LINE_ITEMS[5].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[5].unitPrice, total: 180 },
      { description: SERVICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[4].unitPrice, total: 220 },
    ],
    estimatedTotal: 720,
    notes: 'Enquête filiation famille Aubert.',
    accepted: true,
    acceptedAt: DATES.mar2024_10,
    acceptedBy: 'Vincent Roux',
    status: 'archived',
    createdAt: DATES.mar2024_05,
    updatedAt: DATES.oct2024_08,
  },
  {
    id: ESTIMATE_IDS.released4_est1,
    caseId: CASE_IDS.released4,
    clientId: CLIENT_IDS.company1,
    estimateNumber: 'DEV-2023-208',
    issueDate: DATES.apr2024_18,
    validUntil: DATES.may2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[1].description, quantity: 4, unitPrice: SERVICE_LINE_ITEMS[1].unitPrice, total: 1520 },
      { description: SERVICE_LINE_ITEMS[6].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[6].unitPrice, total: 120 },
      { description: SERVICE_LINE_ITEMS[7].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[7].unitPrice, total: 120 },
    ],
    estimatedTotal: 1760,
    notes: 'Surveillance employee suspicion fraude.',
    accepted: true,
    acceptedAt: DATES.apr2024_25,
    acceptedBy: 'Nicolas Richard',
    status: 'archived',
    createdAt: DATES.apr2024_18,
    updatedAt: DATES.nov2024_05,
  },
  {
    id: ESTIMATE_IDS.released5_est1,
    caseId: CASE_IDS.released5,
    clientId: CLIENT_IDS.lawyer1,
    estimateNumber: 'DEV-2023-252',
    issueDate: DATES.may2024_02,
    validUntil: DATES.jun2025_01,
    lineItems: [
      { description: SERVICE_LINE_ITEMS[0].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[0].unitPrice, total: 450 },
      { description: SERVICE_LINE_ITEMS[5].description, quantity: 2, unitPrice: SERVICE_LINE_ITEMS[5].unitPrice, total: 360 },
      { description: SERVICE_LINE_ITEMS[4].description, quantity: 1, unitPrice: SERVICE_LINE_ITEMS[4].unitPrice, total: 220 },
    ],
    estimatedTotal: 1030,
    notes: 'Enquête successionale famille Dupont-Mercier.',
    accepted: true,
    acceptedAt: DATES.may2024_10,
    acceptedBy: 'Jean-Pierre Aubry',
    status: 'archived',
    createdAt: DATES.may2024_02,
    updatedAt: DATES.dec2024_10,
  },
];

// Helper functions
export function getAllEstimates(): Estimate[] {
  return estimates;
}

export function getEstimateById(id: string): Estimate | undefined {
  return estimates.find(e => e.id === id);
}

export function getEstimatesByCaseId(caseId: string): Estimate[] {
  return estimates.filter(e => e.caseId === caseId);
}

export function getEstimatesByClientId(clientId: string): Estimate[] {
  return estimates.filter(e => e.clientId === clientId);
}

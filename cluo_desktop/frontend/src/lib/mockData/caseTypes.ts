/**
 * Mock case type data - Static, consistent across reloads.
 * Categories/classifications for cases.
 */

import { DATES } from './helpers';
import type { CaseType } from '../types/entities';

// UUIDs for case types
export const CASE_TYPE_IDS = {
  civil: '960e8400-e29b-41d4-a716-446655440001',
  criminal: '960e8400-e29b-41d4-a716-446655440002',
  corporate: '960e8400-e29b-41d4-a716-446655440003',
  family: '960e8400-e29b-41d4-a716-446655440004',
  insurance: '960e8400-e29b-41d4-a716-446655440005',
  background: '960e8400-e29b-41d4-a716-446655440006',
  fraud: '960e8400-e29b-41d4-a716-446655440007',
  missing: '960e8400-e29b-41d4-a716-446655440008',
  surveillance: '960e8400-e29b-41d4-a716-446655440009',
  digital: '960e8400-e29b-41d4-a716-446655440010',
} as const;

export const caseTypes: CaseType[] = [
  {
    id: CASE_TYPE_IDS.civil,
    name: 'Affaire civile',
    createdAt: DATES.jan2023_15,
    updatedAt: DATES.jan2023_15,
  },
  {
    id: CASE_TYPE_IDS.criminal,
    name: 'Affaire pénale',
    createdAt: DATES.jan2023_15,
    updatedAt: DATES.jan2023_15,
  },
  {
    id: CASE_TYPE_IDS.corporate,
    name: 'Enquête corporate',
    createdAt: DATES.jan2023_15,
    updatedAt: DATES.jan2023_15,
  },
  {
    id: CASE_TYPE_IDS.family,
    name: 'Affaire familiale',
    createdAt: DATES.jan2023_15,
    updatedAt: DATES.jan2023_15,
  },
  {
    id: CASE_TYPE_IDS.insurance,
    name: 'Assurance',
    createdAt: DATES.jan2023_15,
    updatedAt: DATES.jan2023_15,
  },
  {
    id: CASE_TYPE_IDS.background,
    name: 'Vérification de background',
    createdAt: DATES.jan2023_15,
    updatedAt: DATES.jan2023_15,
  },
  {
    id: CASE_TYPE_IDS.fraud,
    name: 'Fraude',
    createdAt: DATES.jan2023_15,
    updatedAt: DATES.jan2023_15,
  },
  {
    id: CASE_TYPE_IDS.missing,
    name: 'Personne disparue',
    createdAt: DATES.jan2023_15,
    updatedAt: DATES.jan2023_15,
  },
  {
    id: CASE_TYPE_IDS.surveillance,
    name: 'Filature',
    createdAt: DATES.jan2023_15,
    updatedAt: DATES.jan2023_15,
  },
  {
    id: CASE_TYPE_IDS.digital,
    name: 'Enquête numérique',
    createdAt: DATES.jan2023_15,
    updatedAt: DATES.jan2023_15,
  },
];

// Helper functions
export function getAllCaseTypes(): CaseType[] {
  return caseTypes;
}

export function getCaseTypeById(id: string): CaseType | undefined {
  return caseTypes.find(t => t.id === id);
}

export function getCaseTypesByIds(ids: string[]): CaseType[] {
  return caseTypes.filter(t => ids.includes(t.id));
}

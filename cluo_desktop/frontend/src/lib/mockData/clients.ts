/**
 * Mock client data - Static, consistent across reloads.
 * Diverse client types: persons, insurance companies, lawyers, companies, government.
 */

import { FRENCH_FIRST_NAMES, FRENCH_LAST_NAMES, INSURANCE_COMPANIES, GOVERNMENT_ENTITIES, COMPANY_NAMES, FRENCH_COMPANY_SUFFIXES, DATES } from './helpers';

// UUIDs for clients
export const CLIENT_IDS = {
  // Individual persons
  person1: '650e8400-e29b-41d4-a716-446655440101',
  person2: '650e8400-e29b-41d4-a716-446655440102',
  person3: '650e8400-e29b-41d4-a716-446655440103',

  // Insurance companies
  insurance1: '650e8400-e29b-41d4-a716-446655440201',
  insurance2: '650e8400-e29b-41d4-a716-446655440202',

  // Lawyers
  lawyer1: '650e8400-e29b-41d4-a716-446655440301',
  lawyer2: '650e8400-e29b-41d4-a716-446655440302',

  // Companies
  company1: '650e8400-e29b-41d4-a716-446655440401',
  company2: '650e8400-e29b-41d4-a716-446655440402',

  // Government entities
  government1: '650e8400-e29b-41d4-a716-446655440501',
} as const;

export type ClientType = 'person' | 'insurance' | 'lawyer' | 'company' | 'government';

export interface Client {
  id: string;
  name: string;
  type: ClientType;
  createdAt: string;
}

export const clients: Client[] = [
  // === INDIVIDUAL PERSONS ===
  {
    id: CLIENT_IDS.person1,
    name: 'Jean Dupont',
    type: 'person',
    createdAt: DATES.jan2023_15,
  },
  {
    id: CLIENT_IDS.person2,
    name: 'Marie Leroy',
    type: 'person',
    createdAt: DATES.feb2023_28,
  },
  {
    id: CLIENT_IDS.person3,
    name: 'Pierre Moreau',
    type: 'person',
    createdAt: DATES.mar2023_10,
  },

  // === INSURANCE COMPANIES ===
  {
    id: CLIENT_IDS.insurance1,
    name: 'AXA France',
    type: 'insurance',
    createdAt: DATES.apr2023_22,
  },
  {
    id: CLIENT_IDS.insurance2,
    name: 'Groupama',
    type: 'insurance',
    createdAt: DATES.may2023_05,
  },

  // === LAWYERS ===
  {
    id: CLIENT_IDS.lawyer1,
    name: 'Cabinet Aubry & Associés',
    type: 'lawyer',
    createdAt: DATES.jun2023_18,
  },
  {
    id: CLIENT_IDS.lawyer2,
    name: 'Maître Isabelle Fournier',
    type: 'lawyer',
    createdAt: DATES.jul2023_30,
  },

  // === COMPANIES ===
  {
    id: CLIENT_IDS.company1,
    name: 'Tech Solutions SAS',
    type: 'company',
    createdAt: DATES.aug2023_12,
  },
  {
    id: CLIENT_IDS.company2,
    name: 'Logistics Express SA',
    type: 'company',
    createdAt: DATES.sep2023_25,
  },

  // === GOVERNMENT ENTITIES ===
  {
    id: CLIENT_IDS.government1,
    name: 'Mairie de Paris',
    type: 'government',
    createdAt: DATES.oct2023_08,
  },
];

// Helper functions
export function getAllClients(): Client[] {
  return clients;
}

export function getClientById(id: string): Client | undefined {
  return clients.find(c => c.id === id);
}

export function getClientsByType(type: ClientType): Client[] {
  return clients.filter(c => c.type === type);
}

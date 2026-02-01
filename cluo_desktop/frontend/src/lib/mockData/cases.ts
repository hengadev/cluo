/**
 * Mock case data - Static, consistent across reloads.
 * Diverse investigative cases with different statuses, types, and locations.
 */

import { CASE_TYPES, CITIES, STREET_NAMES, COUNTRY, CASE_DESCRIPTIONS, LOCATION_NOTES, DATES, type CaseStatus, type LocationType, type SubjectRole } from './helpers';
import { CLIENT_IDS } from './clients';
import { CONTACT_IDS } from './contacts';
import { SUBJECT_IDS, caseSubjects, type CaseSubject } from './caseSubjects';

// Re-export types from helpers for convenience
export type { CaseStatus, LocationType, SubjectRole };

// UUIDs for cases
export const CASE_IDS = {
  // Draft cases (early stage, no documents)
  draft1: '950e8400-e29b-41d4-a716-446655440001',

  // In progress cases
  inProgress1: '950e8400-e29b-41d4-a716-446655440101',
  inProgress2: '950e8400-e29b-41d4-a716-446655440102',
  inProgress3: '950e8400-e29b-41d4-a716-446655440103',
  inProgress4: '950e8400-e29b-41d4-a716-446655440104',
  inProgress5: '950e8400-e29b-41d4-a716-446655440105',

  // Ready cases
  ready1: '950e8400-e29b-41d4-a716-446655440201',
  ready2: '950e8400-e29b-41d4-a716-446655440202',
  ready3: '950e8400-e29b-41d4-a716-446655440203',

  // Released cases
  released1: '950e8400-e29b-41d4-a716-446655440301',
  released2: '950e8400-e29b-41d4-a716-446655440302',
  released3: '950e8400-e29b-41d4-a716-446655440303',
  released4: '950e8400-e29b-41d4-a716-446655440304',
  released5: '950e8400-e29b-41d4-a716-446655440305',
} as const;

export interface CaseSubjectAssignment {
  caseId: string;
  subjectId: string;
  role: SubjectRole;
}

export interface Case {
  id: string;
  title: string;
  description: string;
  clientId: string;
  assignedContactId: string | null;
  caseSubjectIds: string[];
  externalReference: string | null;
  caseType: string;
  status: CaseStatus;
  // Location fields
  placename: string;
  address1: string;
  address2: string | null;
  city: string;
  postalCode: string;
  country: string;
  latitude: number;
  longitude: number;
  locationType: LocationType;
  locationNotes: string | null;
  createdAt: string;
  updatedAt: string;
}

// Junction table for case-subject relationships with roles
export const caseSubjectAssignments: CaseSubjectAssignment[] = [
  { caseId: CASE_IDS.inProgress1, subjectId: SUBJECT_IDS.subject1, role: 'suspect' },
  { caseId: CASE_IDS.inProgress1, subjectId: SUBJECT_IDS.subject2, role: 'witness' },

  { caseId: CASE_IDS.inProgress2, subjectId: SUBJECT_IDS.subject3, role: 'suspect' },
  { caseId: CASE_IDS.inProgress2, subjectId: SUBJECT_IDS.subject4, role: 'victim' },

  { caseId: CASE_IDS.inProgress3, subjectId: SUBJECT_IDS.subject5, role: 'claimant' },
  { caseId: CASE_IDS.inProgress3, subjectId: SUBJECT_IDS.subject6, role: 'suspect' },

  { caseId: CASE_IDS.inProgress4, subjectId: SUBJECT_IDS.subject7, role: 'representative' },
  { caseId: CASE_IDS.inProgress4, subjectId: SUBJECT_IDS.subject8, role: 'witness' },

  { caseId: CASE_IDS.inProgress5, subjectId: SUBJECT_IDS.subject9, role: 'victim' },

  { caseId: CASE_IDS.ready1, subjectId: SUBJECT_IDS.subject10, role: 'suspect' },
  { caseId: CASE_IDS.ready1, subjectId: SUBJECT_IDS.subject11, role: 'witness' },
  { caseId: CASE_IDS.ready1, subjectId: SUBJECT_IDS.subject12, role: 'victim' },

  { caseId: CASE_IDS.ready2, subjectId: SUBJECT_IDS.subject13, role: 'claimant' },

  { caseId: CASE_IDS.ready3, subjectId: SUBJECT_IDS.subject14, role: 'suspect' },
  { caseId: CASE_IDS.ready3, subjectId: SUBJECT_IDS.subject15, role: 'witness' },

  { caseId: CASE_IDS.released1, subjectId: SUBJECT_IDS.subject16, role: 'victim' },

  { caseId: CASE_IDS.released2, subjectId: SUBJECT_IDS.subject17, role: 'suspect' },
  { caseId: CASE_IDS.released2, subjectId: SUBJECT_IDS.subject18, role: 'witness' },

  { caseId: CASE_IDS.released3, subjectId: SUBJECT_IDS.subject19, role: 'claimant' },

  { caseId: CASE_IDS.released4, subjectId: SUBJECT_IDS.subject20, role: 'suspect' },

  { caseId: CASE_IDS.released5, subjectId: SUBJECT_IDS.subject1, role: 'victim' },
  { caseId: CASE_IDS.released5, subjectId: SUBJECT_IDS.subject3, role: 'suspect' },
];

export const cases: Case[] = [
  // === DRAFT CASES ===
  {
    id: CASE_IDS.draft1,
    title: 'Affaire Mousquet',
    description: 'Nouvelle demande - informations à compléter',
    clientId: CLIENT_IDS.person1,
    assignedContactId: null,
    caseSubjectIds: [],
    externalReference: null,
    caseType: 'Enquête préliminaire',
    status: 'draft',
    placename: 'Résidence Dupont',
    address1: '15 Rue de la Paix',
    address2: 'Appt 4B',
    city: 'Paris',
    postalCode: '75002',
    country: COUNTRY,
    latitude: 48.8698,
    longitude: 2.3432,
    locationType: 'home',
    locationNotes: 'Zone résidentielle calme, facile d\'accès',
    createdAt: DATES.oct2024_25,
    updatedAt: DATES.oct2024_25,
  },

  // === IN PROGRESS CASES ===
  {
    id: CASE_IDS.inProgress1,
    title: 'Affaire Dupont',
    description: CASE_DESCRIPTIONS[0],
    clientId: CLIENT_IDS.person1,
    assignedContactId: CONTACT_IDS.person1_contact1,
    caseSubjectIds: [SUBJECT_IDS.subject1, SUBJECT_IDS.subject2],
    externalReference: 'REF-2024-001',
    caseType: CASE_TYPES[0],
    status: 'in_progress',
    placename: 'Domicile François Mercier',
    address1: '42 Rue de la République',
    address2: null,
    city: 'Lyon',
    postalCode: '69001',
    country: COUNTRY,
    latitude: 45.7640,
    longitude: 4.8357,
    locationType: 'home',
    locationNotes: LOCATION_NOTES[0],
    createdAt: DATES.jan2024_10,
    updatedAt: DATES.oct2024_25,
  },
  {
    id: CASE_IDS.inProgress2,
    title: 'Surveillance Conflent-Sainte-Honorine',
    description: CASE_DESCRIPTIONS[1],
    clientId: CLIENT_IDS.insurance1,
    assignedContactId: CONTACT_IDS.insurance1_contact2,
    caseSubjectIds: [SUBJECT_IDS.subject3, SUBJECT_IDS.subject4],
    externalReference: 'REF-2024-012',
    caseType: CASE_TYPES[1],
    status: 'in_progress',
    placename: 'Bureau Guillaume',
    address1: '8 Boulevard Haussmann',
    address2: '3ème étage',
    city: 'Paris',
    postalCode: '75009',
    country: COUNTRY,
    latitude: 48.8756,
    longitude: 2.3274,
    locationType: 'business',
    locationNotes: LOCATION_NOTES[3],
    createdAt: DATES.mar2024_22,
    updatedAt: DATES.oct2024_25,
  },
  {
    id: CASE_IDS.inProgress3,
    title: 'Recherche Personne Disparue - Lyon',
    description: CASE_DESCRIPTIONS[2],
    clientId: CLIENT_IDS.lawyer1,
    assignedContactId: CONTACT_IDS.lawyer1_contact1,
    caseSubjectIds: [SUBJECT_IDS.subject5, SUBJECT_IDS.subject6],
    externalReference: 'REF-2024-025',
    caseType: CASE_TYPES[2],
    status: 'in_progress',
    placename: 'Dernière localisation connue',
    address1: '56 Avenue de France',
    address2: null,
    city: 'Marseille',
    postalCode: '13002',
    country: COUNTRY,
    latitude: 43.2965,
    longitude: 5.3698,
    locationType: 'public',
    locationNotes: 'Zone à haute fréquentation',
    createdAt: DATES.may2024_15,
    updatedAt: DATES.oct2024_25,
  },
  {
    id: CASE_IDS.inProgress4,
    title: 'Fraude Assurance Auto - Dossier Martin',
    description: CASE_DESCRIPTIONS[3],
    clientId: CLIENT_IDS.insurance2,
    assignedContactId: CONTACT_IDS.insurance2_contact1,
    caseSubjectIds: [SUBJECT_IDS.subject7, SUBJECT_IDS.subject8],
    externalReference: 'REF-2024-038',
    caseType: CASE_TYPES[3],
    status: 'in_progress',
    placename: 'Parking Résidence Gaillard',
    address1: '78 Boulevard de Sébastopol',
    address2: 'Sous-sol',
    city: 'Paris',
    postalCode: '75004',
    country: COUNTRY,
    latitude: 48.8566,
    longitude: 2.3522,
    locationType: 'vehicle',
    locationNotes: 'Accès limité - badge nécessaire',
    createdAt: DATES.jul2024_10,
    updatedAt: DATES.nov2024_22,
  },
  {
    id: CASE_IDS.inProgress5,
    title: 'Vérification CV - Candidature Directeur',
    description: CASE_DESCRIPTIONS[4],
    clientId: CLIENT_IDS.company1,
    assignedContactId: CONTACT_IDS.company1_contact1,
    caseSubjectIds: [SUBJECT_IDS.subject9],
    externalReference: 'REF-2024-052',
    caseType: CASE_TYPES[4],
    status: 'in_progress',
    placename: 'Siège Social Tech Solutions',
    address1: '19 Avenue de la Grande-Armée',
    address2: null,
    city: 'Paris',
    postalCode: '75017',
    country: COUNTRY,
    latitude: 48.8758,
    longitude: 2.2984,
    locationType: 'business',
    locationNotes: 'Zone professionnelle, sécurité renforcée',
    createdAt: DATES.sep2024_12,
    updatedAt: DATES.nov2024_22,
  },

  // === READY CASES ===
  {
    id: CASE_IDS.ready1,
    title: 'Conflit Conjugal - Affaire Roussel',
    description: CASE_DESCRIPTIONS[5],
    clientId: CLIENT_IDS.person2,
    assignedContactId: CONTACT_IDS.person2_contact1,
    caseSubjectIds: [SUBJECT_IDS.subject10, SUBJECT_IDS.subject11, SUBJECT_IDS.subject12],
    externalReference: 'REF-2024-018',
    caseType: CASE_TYPES[5],
    status: 'ready',
    placename: 'Domicile conjugal',
    address1: '12 Rue de la Paix',
    address2: null,
    city: 'Paris',
    postalCode: '75002',
    country: COUNTRY,
    latitude: 48.8698,
    longitude: 2.3432,
    locationType: 'home',
    locationNotes: 'Discrétion absolue requise',
    createdAt: DATES.feb2024_14,
    updatedAt: DATES.nov2024_05,
  },
  {
    id: CASE_IDS.ready2,
    title: 'Accusation Abusive - Harcèlement',
    description: CASE_DESCRIPTIONS[6],
    clientId: CLIENT_IDS.company2,
    assignedContactId: CONTACT_IDS.company2_contact1,
    caseSubjectIds: [SUBJECT_IDS.subject13],
    externalReference: 'REF-2024-043',
    caseType: CASE_TYPES[6],
    status: 'ready',
    placename: 'Bureaux Logistics Express',
    address1: '23 Boulevard de Magenta',
    address2: 'Bâtiment C',
    city: 'Paris',
    postalCode: '75010',
    country: COUNTRY,
    latitude: 48.8765,
    longitude: 2.3541,
    locationType: 'business',
    locationNotes: 'Environnement professionnel sensible',
    createdAt: DATES.apr2024_08,
    updatedAt: DATES.nov2024_22,
  },
  {
    id: CASE_IDS.ready3,
    title: 'Vol Propriété Intellectuelle',
    description: CASE_DESCRIPTIONS[7],
    clientId: CLIENT_IDS.lawyer2,
    assignedContactId: CONTACT_IDS.lawyer2_contact1,
    caseSubjectIds: [SUBJECT_IDS.subject14, SUBJECT_IDS.subject15],
    externalReference: 'REF-2024-061',
    caseType: CASE_TYPES[7],
    status: 'ready',
    placename: 'Anciens bureaux',
    address1: '16 Rue de la Tour',
    address2: null,
    city: 'Nice',
    postalCode: '06000',
    country: COUNTRY,
    latitude: 43.7102,
    longitude: 7.2620,
    locationType: 'business',
    locationNotes: 'Zone industrielle - sécurité accrue',
    createdAt: DATES.jun2024_20,
    updatedAt: DATES.dec2024_10,
  },

  // === RELEASED CASES ===
  {
    id: CASE_IDS.released1,
    title: 'Harcèlement Travail - Dossier RH',
    description: CASE_DESCRIPTIONS[8],
    clientId: CLIENT_IDS.government1,
    assignedContactId: CONTACT_IDS.government1_contact1,
    caseSubjectIds: [SUBJECT_IDS.subject16],
    externalReference: 'REF-2023-089',
    caseType: CASE_TYPES[8],
    status: 'released',
    placename: 'Mairie - Bureau DGS',
    address1: '5 Rue Lobau',
    address2: '3ème étage',
    city: 'Paris',
    postalCode: '75004',
    country: COUNTRY,
    latitude: 48.8564,
    longitude: 2.3570,
    locationType: 'business',
    locationNotes: 'Administration publique - protocole strict',
    createdAt: DATES.jan2024_25,
    updatedAt: DATES.aug2024_18,
  },
  {
    id: CASE_IDS.released2,
    title: 'Infraction Code de la Route',
    description: CASE_DESCRIPTIONS[9],
    clientId: CLIENT_IDS.person3,
    assignedContactId: CONTACT_IDS.person3_contact1,
    caseSubjectIds: [SUBJECT_IDS.subject17, SUBJECT_IDS.subject18],
    externalReference: 'REF-2023-112',
    caseType: CASE_TYPES[9],
    status: 'released',
    placename: 'Lieu de l\'accident',
    address1: '72 Boulevard de Charonne',
    address2: 'Intersection',
    city: 'Paris',
    postalCode: '75011',
    country: COUNTRY,
    latitude: 48.8546,
    longitude: 2.3987,
    locationType: 'public',
    locationNotes: 'Lieu public - circulation dense',
    createdAt: DATES.feb2024_28,
    updatedAt: DATES.sep2024_01,
  },
  {
    id: CASE_IDS.released3,
    title: 'Enquête Filiation - Famille Aubert',
    description: CASE_DESCRIPTIONS[0],
    clientId: CLIENT_IDS.insurance1,
    assignedContactId: CONTACT_IDS.insurance1_contact3,
    caseSubjectIds: [SUBJECT_IDS.subject19],
    externalReference: 'REF-2023-156',
    caseType: CASE_TYPES[0],
    status: 'released',
    placename: 'Domicile Aubert',
    address1: '24 Avenue de la République',
    address2: null,
    city: 'Lyon',
    postalCode: '69006',
    country: COUNTRY,
    latitude: 45.7625,
    longitude: 4.8462,
    locationType: 'home',
    locationNotes: 'Quartier résidentiel familial',
    createdAt: DATES.mar2024_05,
    updatedAt: DATES.oct2024_08,
  },
  {
    id: CASE_IDS.released4,
    title: 'Surveillance Employee - Suspicion Fraude',
    description: CASE_DESCRIPTIONS[1],
    clientId: CLIENT_IDS.company1,
    assignedContactId: CONTACT_IDS.company1_contact2,
    caseSubjectIds: [SUBJECT_IDS.subject20],
    externalReference: 'REF-2023-201',
    caseType: CASE_TYPES[1],
    status: 'released',
    placename: 'Tech Solutions - Siège',
    address1: '19 Avenue de la Grande-Armée',
    address2: 'Zone Industrielle',
    city: 'Paris',
    postalCode: '75017',
    country: COUNTRY,
    latitude: 48.8758,
    longitude: 2.2984,
    locationType: 'business',
    locationNotes: 'Zone professionnelle - surveillance discrète',
    createdAt: DATES.apr2024_18,
    updatedAt: DATES.nov2024_05,
  },
  {
    id: CASE_IDS.released5,
    title: 'Conflit Succession - Famille Dupont-Mercier',
    description: 'Litige successoral complexe avec plusieurs parties prenantes. Enquête sur les testaments et les relations familiales.',
    clientId: CLIENT_IDS.lawyer1,
    assignedContactId: CONTACT_IDS.lawyer1_contact2,
    caseSubjectIds: [SUBJECT_IDS.subject1, SUBJECT_IDS.subject3],
    externalReference: 'REF-2023-245',
    caseType: 'Enquête successorale',
    status: 'released',
    placename: 'Étude Notariale',
    address1: '15 Rue de la Paix',
    address2: null,
    city: 'Paris',
    postalCode: '75002',
    country: COUNTRY,
    latitude: 48.8698,
    longitude: 2.3432,
    locationType: 'business',
    locationNotes: 'Étude notariale - confidentialité requise',
    createdAt: DATES.may2024_02,
    updatedAt: DATES.dec2024_10,
  },
];

// Helper functions
export function getAllCases(): Case[] {
  return cases;
}

export function getCaseById(id: string): Case | undefined {
  return cases.find(c => c.id === id);
}

export function getCasesByStatus(status: CaseStatus): Case[] {
  return cases.filter(c => c.status === status);
}

export function getCasesByClientId(clientId: string): Case[] {
  return cases.filter(c => c.clientId === clientId);
}

export function getSubjectIdsForCase(caseId: string): string[] {
  return caseSubjectAssignments
    .filter(a => a.caseId === caseId)
    .map(a => a.subjectId);
}

export function getSubjectRole(caseId: string, subjectId: string): SubjectRole | undefined {
  const assignment = caseSubjectAssignments.find(
    a => a.caseId === caseId && a.subjectId === subjectId
  );
  return assignment?.role;
}

export function getCaseSubjectsWithRoles(caseId: string): Array<{ subject: CaseSubject, role: SubjectRole }> {
  const assignments = caseSubjectAssignments.filter(a => a.caseId === caseId);
  return assignments.map(a => ({
    subject: caseSubjects.find((s: CaseSubject) => s.id === a.subjectId)!,
    role: a.role
  }));
}

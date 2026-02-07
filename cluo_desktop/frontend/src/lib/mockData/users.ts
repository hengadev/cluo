/**
 * Mock user data - Static, consistent across reloads.
 * Represents the users of the investigation system.
 */

import { FRENCH_FIRST_NAMES, FRENCH_LAST_NAMES, DATES, type UserRole } from './helpers';

// UUIDs for users
export const USER_IDS = {
  admin: '550e8400-e29b-41d4-a716-446655440001',
  investigator1: '550e8400-e29b-41d4-a716-446655440002',
  investigator2: '550e8400-e29b-41d4-a716-446655440003',
  viewer1: '550e8400-e29b-41d4-a716-446655440004',
  viewer2: '550e8400-e29b-41d4-a716-446655440005',
} as const;

export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: UserRole;
  createdAt: string;
}

export const users: User[] = [
  {
    id: USER_IDS.admin,
    email: 'martin.dubois@enquete.fr',
    firstName: 'Martin',
    lastName: 'Dubois',
    role: 'admin',
    createdAt: DATES.jan2023_15,
  },
  {
    id: USER_IDS.investigator1,
    email: 'sophie.martin@enquete.fr',
    firstName: 'Sophie',
    lastName: 'Martin',
    role: 'investigator',
    createdAt: DATES.feb2023_28,
  },
  {
    id: USER_IDS.investigator2,
    email: 'jean.pierre@enquete.fr',
    firstName: 'Jean',
    lastName: 'Pierre',
    role: 'investigator',
    createdAt: DATES.mar2023_10,
  },
  {
    id: USER_IDS.viewer1,
    email: 'marie.bernard@enquete.fr',
    firstName: 'Marie',
    lastName: 'Bernard',
    role: 'viewer',
    createdAt: DATES.apr2023_22,
  },
  {
    id: USER_IDS.viewer2,
    email: 'philippe.thomas@enquete.fr',
    firstName: 'Philippe',
    lastName: 'Thomas',
    role: 'viewer',
    createdAt: DATES.may2023_05,
  },
];

// Helper functions
export function getAllUsers(): User[] {
  return users;
}

export function getUserById(id: string): User | undefined {
  return users.find(u => u.id === id);
}

export function getUsersByRole(role: UserRole): User[] {
  return users.filter(u => u.role === role);
}

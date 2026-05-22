/**
 * Mock user data - Static, consistent across reloads.
 * Represents the users of the investigation system.
 */

import { FRENCH_FIRST_NAMES, FRENCH_LAST_NAMES, DATES, type UserRole } from './helpers';
import type { AuthUser } from '../types/entities';

// UUIDs for users
export const USER_IDS = {
  admin: '550e8400-e29b-41d4-a716-446655440001',
  investigator1: '550e8400-e29b-41d4-a716-446655440002',
  investigator2: '550e8400-e29b-41d4-a716-446655440003',
  viewer1: '550e8400-e29b-41d4-a716-446655440004',
  viewer2: '550e8400-e29b-41d4-a716-446655440005',
} as const;

export const users: AuthUser[] = [
  {
    id: USER_IDS.admin,
    email: 'martin.dubois@enquete.fr',
    role: 'admin',
  },
  {
    id: USER_IDS.investigator1,
    email: 'sophie.martin@enquete.fr',
    role: 'investigator',
  },
  {
    id: USER_IDS.investigator2,
    email: 'jean.pierre@enquete.fr',
    role: 'investigator',
  },
  {
    id: USER_IDS.viewer1,
    email: 'marie.bernard@enquete.fr',
    role: 'viewer',
  },
  {
    id: USER_IDS.viewer2,
    email: 'philippe.thomas@enquete.fr',
    role: 'viewer',
  },
];

// Helper functions
export function getAllUsers(): AuthUser[] {
  return users;
}

export function getUserById(id: string): AuthUser | undefined {
  return users.find(u => u.id === id);
}

export function getUsersByRole(role: UserRole): AuthUser[] {
  return users.filter(u => u.role === role);
}

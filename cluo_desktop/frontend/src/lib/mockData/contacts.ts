/**
 * Mock contact data - Static, consistent across reloads.
 * Each client has 2-4 contacts with position, email, and phone.
 */

import { FRENCH_FIRST_NAMES, FRENCH_LAST_NAMES, JOB_TITLES, CITIES, STREET_NAMES, COUNTRY, DATES } from './helpers';
import { CLIENT_IDS } from './clients';

// UUIDs for contacts
export const CONTACT_IDS = {
  // Client: Jean Dupont (person1) - 2 contacts
  person1_contact1: '750e8400-e29b-41d4-a716-446655440101',
  person1_contact2: '750e8400-e29b-41d4-a716-446655440102',

  // Client: Marie Leroy (person2) - 3 contacts
  person2_contact1: '750e8400-e29b-41d4-a716-446655440103',
  person2_contact2: '750e8400-e29b-41d4-a716-446655440104',
  person2_contact3: '750e8400-e29b-41d4-a716-446655440105',

  // Client: Pierre Moreau (person3) - 2 contacts
  person3_contact1: '750e8400-e29b-41d4-a716-446655440106',
  person3_contact2: '750e8400-e29b-41d4-a716-446655440107',

  // Client: AXA France (insurance1) - 4 contacts
  insurance1_contact1: '750e8400-e29b-41d4-a716-446655440201',
  insurance1_contact2: '750e8400-e29b-41d4-a716-446655440202',
  insurance1_contact3: '750e8400-e29b-41d4-a716-446655440203',
  insurance1_contact4: '750e8400-e29b-41d4-a716-446655440204',

  // Client: Groupama (insurance2) - 3 contacts
  insurance2_contact1: '750e8400-e29b-41d4-a716-446655440205',
  insurance2_contact2: '750e8400-e29b-41d4-a716-446655440206',
  insurance2_contact3: '750e8400-e29b-41d4-a716-446655440207',

  // Client: Cabinet Aubry & Associés (lawyer1) - 2 contacts
  lawyer1_contact1: '750e8400-e29b-41d4-a716-446655440301',
  lawyer1_contact2: '750e8400-e29b-41d4-a716-446655440302',

  // Client: Maître Isabelle Fournier (lawyer2) - 3 contacts
  lawyer2_contact1: '750e8400-e29b-41d4-a716-446655440303',
  lawyer2_contact2: '750e8400-e29b-41d4-a716-446655440304',
  lawyer2_contact3: '750e8400-e29b-41d4-a716-446655440305',

  // Client: Tech Solutions SAS (company1) - 3 contacts
  company1_contact1: '750e8400-e29b-41d4-a716-446655440401',
  company1_contact2: '750e8400-e29b-41d4-a716-446655440402',
  company1_contact3: '750e8400-e29b-41d4-a716-446655440403',

  // Client: Logistics Express SA (company2) - 2 contacts
  company2_contact1: '750e8400-e29b-41d4-a716-446655440404',
  company2_contact2: '750e8400-e29b-41d4-a716-446655440405',

  // Client: Mairie de Paris (government1) - 4 contacts
  government1_contact1: '750e8400-e29b-41d4-a716-446655440501',
  government1_contact2: '750e8400-e29b-41d4-a716-446655440502',
  government1_contact3: '750e8400-e29b-41d4-a716-446655440503',
  government1_contact4: '750e8400-e29b-41d4-a716-446655440504',
} as const;

export interface Contact {
  id: string;
  clientID: string;
  lastname: string;
  firstname: string;
  email: string;
  phone: string;
  position: string;
  createdAt: string;
}

export const contacts: Contact[] = [
  // === Jean Dupont (person) - 2 contacts ===
  {
    id: CONTACT_IDS.person1_contact1,
    clientID: CLIENT_IDS.person1,
    lastname: 'Dupont',
    firstname: 'Jean',
    email: 'jean.dupont@email.fr',
    phone: '06 12 34 56 78',
    position: 'Particulier',
    createdAt: DATES.jan2023_15,
  },
  {
    id: CONTACT_IDS.person1_contact2,
    clientID: CLIENT_IDS.person1,
    lastname: 'Dupont',
    firstname: 'Françoise',
    email: 'f.dupont@email.fr',
    phone: '06 98 76 54 32',
    position: 'Conjoint',
    createdAt: DATES.jan2023_15,
  },

  // === Marie Leroy (person) - 3 contacts ===
  {
    id: CONTACT_IDS.person2_contact1,
    clientID: CLIENT_IDS.person2,
    lastname: 'Leroy',
    firstname: 'Marie',
    email: 'marie.leroy@orange.fr',
    phone: '06 23 45 67 89',
    position: 'Particulier',
    createdAt: DATES.feb2023_28,
  },
  {
    id: CONTACT_IDS.person2_contact2,
    clientID: CLIENT_IDS.person2,
    lastname: 'Leroy',
    firstname: 'Philippe',
    email: 'p.leroy@societe.fr',
    phone: '06 45 67 89 01',
    position: 'Frère',
    createdAt: DATES.feb2023_28,
  },
  {
    id: CONTACT_IDS.person2_contact3,
    clientID: CLIENT_IDS.person2,
    lastname: 'Martinez',
    firstname: 'Carlos',
    email: 'carlos.m@avocat.fr',
    phone: '01 43 52 64 75',
    position: 'Conseiller Juridique',
    createdAt: DATES.mar2023_10,
  },

  // === Pierre Moreau (person) - 2 contacts ===
  {
    id: CONTACT_IDS.person3_contact1,
    clientID: CLIENT_IDS.person3,
    lastname: 'Moreau',
    firstname: 'Pierre',
    email: 'pierre.moreau@gmail.com',
    phone: '06 87 65 43 21',
    position: 'Particulier',
    createdAt: DATES.mar2023_10,
  },
  {
    id: CONTACT_IDS.person3_contact2,
    clientID: CLIENT_IDS.person3,
    lastname: 'Moreau',
    firstname: 'Claire',
    email: 'claire.moreau@gmail.com',
    phone: '06 11 22 33 44',
    position: 'Fille',
    createdAt: DATES.mar2023_10,
  },

  // === AXA France (insurance) - 4 contacts ===
  {
    id: CONTACT_IDS.insurance1_contact1,
    clientID: CLIENT_IDS.insurance1,
    lastname: 'Fontaine',
    firstname: 'Nathalie',
    email: 'nathalie.fontaine@axa.fr',
    phone: '01 53 42 67 89',
    position: 'Directrice Régionale',
    createdAt: DATES.apr2023_22,
  },
  {
    id: CONTACT_IDS.insurance1_contact2,
    clientID: CLIENT_IDS.insurance1,
    lastname: 'Roux',
    firstname: 'Vincent',
    email: 'vincent.roux@axa.fr',
    phone: '01 53 42 67 90',
    position: 'Gestionnaire de Sinistres',
    createdAt: DATES.apr2023_22,
  },
  {
    id: CONTACT_IDS.insurance1_contact3,
    clientID: CLIENT_IDS.insurance1,
    lastname: 'Lefebvre',
    firstname: 'Sophie',
    email: 'sophie.lefebvre@axa.fr',
    phone: '01 53 42 67 91',
    position: 'Chargée de Clientèle',
    createdAt: DATES.may2023_05,
  },
  {
    id: CONTACT_IDS.insurance1_contact4,
    clientID: CLIENT_IDS.insurance1,
    lastname: 'Garcia',
    firstname: 'David',
    email: 'david.garcia@axa.fr',
    phone: '01 53 42 67 92',
    position: 'Responsable Juridique',
    createdAt: DATES.jun2023_18,
  },

  // === Groupama (insurance) - 3 contacts ===
  {
    id: CONTACT_IDS.insurance2_contact1,
    clientID: CLIENT_IDS.insurance2,
    lastname: 'Bertrand',
    firstname: 'Christophe',
    email: 'c.bertrand@groupama.fr',
    phone: '01 34 56 78 90',
    position: 'Directeur des Sinistres',
    createdAt: DATES.may2023_05,
  },
  {
    id: CONTACT_IDS.insurance2_contact2,
    clientID: CLIENT_IDS.insurance2,
    lastname: 'Petit',
    firstname: 'Valérie',
    email: 'valerie.petit@groupama.fr',
    phone: '01 34 56 78 91',
    position: 'Gestionnaire Senior',
    createdAt: DATES.may2023_05,
  },
  {
    id: CONTACT_IDS.insurance2_contact3,
    clientID: CLIENT_IDS.insurance2,
    lastname: 'Morel',
    firstname: 'Stéphane',
    email: 'stephane.morel@groupama.fr',
    phone: '01 34 56 78 92',
    position: 'Attaché Commercial',
    createdAt: DATES.jun2023_18,
  },

  // === Cabinet Aubry & Associés (lawyer) - 2 contacts ===
  {
    id: CONTACT_IDS.lawyer1_contact1,
    clientID: CLIENT_IDS.lawyer1,
    lastname: 'Aubry',
    firstname: 'Jean-Pierre',
    email: 'jp.aubry@aubry-associes.fr',
    phone: '01 42 86 53 74',
    position: 'Fondateur & Associé',
    createdAt: DATES.jun2023_18,
  },
  {
    id: CONTACT_IDS.lawyer1_contact2,
    clientID: CLIENT_IDS.lawyer1,
    lastname: 'Dubois',
    firstname: 'Céline',
    email: 'celine.dubois@aubry-associes.fr',
    phone: '01 42 86 53 75',
    position: 'Associée',
    createdAt: DATES.jul2023_30,
  },

  // === Maître Isabelle Fournier (lawyer) - 3 contacts ===
  {
    id: CONTACT_IDS.lawyer2_contact1,
    clientID: CLIENT_IDS.lawyer2,
    lastname: 'Fournier',
    firstname: 'Isabelle',
    email: 'i.fournier@avocat.fr',
    phone: '01 44 73 82 91',
    position: 'Fondatrice',
    createdAt: DATES.jul2023_30,
  },
  {
    id: CONTACT_IDS.lawyer2_contact2,
    clientID: CLIENT_IDS.lawyer2,
    lastname: 'Marty',
    firstname: 'Patrick',
    email: 'patrick.marty@fournier.fr',
    phone: '01 44 73 82 92',
    position: 'Collaborateur',
    createdAt: DATES.aug2023_12,
  },
  {
    id: CONTACT_IDS.lawyer2_contact3,
    clientID: CLIENT_IDS.lawyer2,
    lastname: 'Lambert',
    firstname: 'Delphine',
    email: 'delphine.lambert@fournier.fr',
    phone: '01 44 73 82 93',
    position: 'Assistante Juridique',
    createdAt: DATES.aug2023_12,
  },

  // === Tech Solutions SAS (company) - 3 contacts ===
  {
    id: CONTACT_IDS.company1_contact1,
    clientID: CLIENT_IDS.company1,
    lastname: 'Richard',
    firstname: 'Nicolas',
    email: 'n.richard@techsolutions.fr',
    phone: '01 76 54 32 10',
    position: 'Directeur Général',
    createdAt: DATES.aug2023_12,
  },
  {
    id: CONTACT_IDS.company1_contact2,
    clientID: CLIENT_IDS.company1,
    lastname: 'André',
    firstname: 'Sandrine',
    email: 's.andre@techsolutions.fr',
    phone: '01 76 54 32 11',
    position: 'Responsable RH',
    createdAt: DATES.aug2023_12,
  },
  {
    id: CONTACT_IDS.company1_contact3,
    clientID: CLIENT_IDS.company1,
    lastname: 'Simon',
    firstname: 'Laurent',
    email: 'l.simon@techsolutions.fr',
    phone: '01 76 54 32 12',
    position: 'Responsable Sécurité',
    createdAt: DATES.sep2023_25,
  },

  // === Logistics Express SA (company) - 2 contacts ===
  {
    id: CONTACT_IDS.company2_contact1,
    clientID: CLIENT_IDS.company2,
    lastname: 'Michel',
    firstname: 'Olivier',
    email: 'o.michel@logistics-express.com',
    phone: '01 39 87 65 43',
    position: 'Directeur des Opérations',
    createdAt: DATES.sep2023_25,
  },
  {
    id: CONTACT_IDS.company2_contact2,
    clientID: CLIENT_IDS.company2,
    lastname: 'Girard',
    firstname: 'Charlotte',
    email: 'charlotte.girard@logistics-express.com',
    phone: '01 39 87 65 44',
    position: 'Assistante Direction',
    createdAt: DATES.oct2023_08,
  },

  // === Mairie de Paris (government) - 4 contacts ===
  {
    id: CONTACT_IDS.government1_contact1,
    clientID: CLIENT_IDS.government1,
    lastname: 'Blanc',
    firstname: 'Alain',
    email: 'alain.blanc@paris.fr',
    phone: '01 42 76 80 00',
    position: 'Chef de Service',
    createdAt: DATES.oct2023_08,
  },
  {
    id: CONTACT_IDS.government1_contact2,
    clientID: CLIENT_IDS.government1,
    lastname: 'Rousseau',
    firstname: 'Martine',
    email: 'martine.rousseau@paris.fr',
    phone: '01 42 76 80 01',
    position: 'Responsable Juridique',
    createdAt: DATES.oct2023_08,
  },
  {
    id: CONTACT_IDS.government1_contact3,
    clientID: CLIENT_IDS.government1,
    lastname: 'Vincent',
    firstname: 'Catherine',
    email: 'catherine.vincent@paris.fr',
    phone: '01 42 76 80 02',
    position: 'Déléguée Régionale',
    createdAt: DATES.nov2023_20,
  },
  {
    id: CONTACT_IDS.government1_contact4,
    clientID: CLIENT_IDS.government1,
    lastname: 'Fournier',
    firstname: 'Thierry',
    email: 'thierry.fournier@paris.fr',
    phone: '01 42 76 80 03',
    position: 'Coordinateur',
    createdAt: DATES.nov2023_20,
  },
];

// Helper functions
export function getAllContacts(): Contact[] {
  return contacts;
}

export function getContactById(id: string): Contact | undefined {
  return contacts.find(c => c.id === id);
}

export function getContactsByClientId(clientID: string): Contact[] {
  return contacts.filter(c => c.clientID === clientID);
}

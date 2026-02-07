/**
 * Helper utilities and static data pools for mock data generation.
 * All data is fixed and static - no random generation.
 */

// ============================================================================
// FRENCH DATA POOLS - Static arrays for realistic French names and addresses
// ============================================================================

export const FRENCH_FIRST_NAMES = [
  'Jean', 'Marie', 'Pierre', 'Sophie', 'Michel', 'Isabelle', 'Philippe', 'Nathalie',
  'Alain', 'Catherine', 'Nicolas', 'Sandrine', 'Christophe', 'Valérie', 'David',
  'Laura', 'Stéphane', 'Céline', 'Patrick', 'Martine', 'Julien', 'Aurélie',
  'Olivier', 'Charlotte', 'Ludovic', 'Camille', 'Thierry', 'Delphine', 'Sébastien'
];

export const FRENCH_LAST_NAMES = [
  'Martin', 'Bernard', 'Dubois', 'Thomas', 'Robert', 'Richard', 'Petit', 'Durand',
  'Leroy', 'Moreau', 'Simon', 'Laurent', 'Lefebvre', 'Michel', 'Garcia', 'David',
  'Bertrand', 'Roux', 'Vincent', 'Fournier', 'Morel', 'André', 'Girard', 'Bonnet',
  'Dupont', 'Lambert', 'Fontaine', 'Rousseau', 'Meyer', 'Blanc', 'Marty'
];

export const FRENCH_COMPANY_SUFFIXES = [
  'SARL', 'SA', 'SAS', 'EURL', 'GIE'
];

export const STREET_NAMES = [
  'Rue de la République', 'Avenue des Champs-Élysées', 'Boulevard Haussmann',
  'Rue de Rivoli', 'Avenue Montaigne', 'Rue du Faubourg Saint-Honoré',
  'Boulevard Saint-Germain', 'Rue de la Liberté', 'Avenue de France',
  'Rue de la Paix', 'Boulevard de Sébastopol', 'Rue de la Bourse',
  'Avenue des Ternes', 'Rue de Lévis', 'Boulevard de Magenta',
  'Rue de la Roquette', 'Avenue de la Grande-Armée', 'Rue de la Pompe',
  'Boulevard de la Villette', 'Rue de la Convention', 'Avenue de Breteuil',
  'Rue de la Tombe-Issoire', 'Boulevard du Montparnasse', 'Rue de la Tour',
  'Avenue de Suffren', 'Rue de la Université', 'Boulevard de Charonne',
  'Rue de la Chine', 'Avenue de la République', 'Rue de la Forge'
];

export const CITIES = [
  { name: 'Paris', postalCode: '75001', region: 'Île-de-France' },
  { name: 'Paris', postalCode: '75008', region: 'Île-de-France' },
  { name: 'Paris', postalCode: '75015', region: 'Île-de-France' },
  { name: 'Lyon', postalCode: '69001', region: 'Auvergne-Rhône-Alpes' },
  { name: 'Lyon', postalCode: '69006', region: 'Auvergne-Rhône-Alpes' },
  { name: 'Marseille', postalCode: '13001', region: 'Provence-Alpes-Côte d\'Azur' },
  { name: 'Marseille', postalCode: '13007', region: 'Provence-Alpes-Côte d\'Azur' },
  { name: 'Toulouse', postalCode: '31000', region: 'Occitanie' },
  { name: 'Nice', postalCode: '06000', region: 'Provence-Alpes-Côte d\'Azur' },
  { name: 'Nantes', postalCode: '44000', region: 'Pays de la Loire' },
  { name: 'Strasbourg', postalCode: '67000', region: 'Grand Est' },
  { name: 'Montpellier', postalCode: '34000', region: 'Occitanie' },
  { name: 'Bordeaux', postalCode: '33000', region: 'Nouvelle-Aquitaine' },
  { name: 'Lille', postalCode: '59000', region: 'Hauts-de-France' },
  { name: 'Rennes', postalCode: '35000', region: 'Bretagne' },
  { name: 'Reims', postalCode: '51100', region: 'Grand Est' },
  { name: 'Dijon', postalCode: '21000', region: 'Bourgogne-Franche-Comté' },
  { name: 'Metz', postalCode: '57000', region: 'Grand Est' },
  { name: 'Nancy', postalCode: '54000', region: 'Grand Est' },
  { name: 'Aix-en-Provence', postalCode: '13100', region: 'Provence-Alpes-Côte d\'Azur' }
];

export const COUNTRY = 'France';

export const JOB_TITLES = [
  'Directeur Général', 'Responsable Juridique', 'Chef de Service',
  'Gestionnaire de Sinistres', 'Conseiller Clientèle', 'Assistante Direction',
  'Chargé de Clientèle', 'Responsable Administration', 'Délégué Régional',
  'Attaché Commercial', 'Coordinateur', 'Responsable des Assurances'
];

export const OCCUPATIONS = [
  'Comptable', 'Médecin', 'Avocat', 'Architecte', 'Ingénieur',
  'Enseignant', 'Cadre Supérieur', 'Commerçant', 'Artisan', 'Retraité',
  'Sans Emploi', 'Chef d\'Entreprise', 'Consultant', 'Agent Immobilier',
  'Profession Libérale', 'Fonctionnaire', 'Assistante', 'Technicien'
];

export const INSURANCE_COMPANIES = [
  'AXA France', 'Allianz France', 'Generali', 'CNP Assurances',
  'Groupama', 'MACIF', 'MAAF', 'MATMUT', 'Credit Agricole Assurances'
];

export const GOVERNMENT_ENTITIES = [
  'Mairie de Paris', 'Conseil Départemental', 'Préfecture de Police',
  'Direction Départementale des Territoires', 'Caisse Primaire d\'Assurance Maladie'
];

// ============================================================================
// COMMON CONSTANTS
// ============================================================================

export const CURRENCY = 'EUR';

export const PAYMENT_TERMS = [
  'Paiement à 30 jours', 'Paiement à 45 jours', 'Paiement à 60 jours',
  'Paiement comptable', 'Paiement en 3 fois', 'Paiement en ligne'
];

export const CASE_TYPES = [
  'Enquête de filiation', 'Surveillance', 'Recherche de personnes',
  'Fraude assurance', 'Vérification de CV', 'Conflit conjugal',
  'Accusation abusive', 'Vol de propriété intellectuelle',
  'Harcèlement au travail', 'Infraction au code de la route'
];

export const CASE_STATUS_VALUES = ['draft', 'in_progress', 'ready', 'released'] as const;
export type CaseStatus = typeof CASE_STATUS_VALUES[number];

export const DOCUMENT_STATUSES = ['draft', 'sent', 'signed', 'active', 'archived', 'cancelled', 'rejected', 'expired'] as const;
export type DocumentStatus = typeof DOCUMENT_STATUSES[number];

export const PAYMENT_STATUSES = ['unpaid', 'paid', 'partially_paid', 'overdue', 'refunded', 'void'] as const;
export type PaymentStatus = typeof PAYMENT_STATUSES[number];

export const SUBJECT_ROLES = ['victim', 'suspect', 'witness', 'claimant', 'representative'] as const;
export type SubjectRole = typeof SUBJECT_ROLES[number];

export const LOCATION_TYPES = ['home', 'business', 'public', 'vehicle', 'other'] as const;
export type LocationType = typeof LOCATION_TYPES[number];

export const USER_ROLES = ['admin', 'investigator', 'viewer'] as const;
export type UserRole = typeof USER_ROLES[number];

// ============================================================================
// LEGAL/FORM TEXT SNIPPETS
// ============================================================================

export const MANDATE_SCOPES = [
  'Enquête de filiation avec vérification des liens familiaux et constitution de dossier complet',
  'Surveillance discrète et collecte de preuves photographiques',
  'Investigation approfondie sur les allégations de fraude avec interviews de témoins',
  'Recherche de personnes disparues ou injoignables avec enquête terrain',
  'Vérification des antécédents et constitution de dossier de preuves',
  'Enquête conjointe pour conflit familial avec médiation',
  'Investigation pour vol de propriété intellectuelle avec collecte de preuves'
];

export const CONTRACT_SCOPES = [
  'Prestation complète d\'investigation privée incluant surveillance, interviews et rapport détaillé',
  'Services d\'enquête pour fraude présumée avec vérification documentaire et enquête terrain',
  'Missions de recherche et localisation avec rapport de findings et recommandations',
  'Investigation pour litige civil avec collecte de preuves et témoignages',
  'Services de surveillance continue avec rapports périodiques et documentation photo/vidéo'
];

export const CONFIDENTIALITY_CLAUSES = [
  'Le Prestataire s\'engage à maintenir la confidentialité absolue de toutes les informations recueillies dans le cadre de cette prestation.',
  'Toutes les informations échangées sont couvertes par le secret professionnel et ne peuvent être divulguées sans accord préalable écrit.',
  'Les données collectées sont protégées conformément au RGPD et ne seront utilisées que dans le cadre strict de la mission.'
];

export const TERMINATION_CLAUSES = [
  'Chaque partie peut résilier ce contrat avec un préavis de 30 jours par lettre recommandée avec accusé de réception.',
  'En cas de résiliation anticipée, le Client s\'engage à payer les prestations effectuées au prorata des services rendus.',
  'Le Prestataire se réserve le droit de suspendre ou résilier le contrat en cas de non-paiement.'
];

export const GOVERNING_LAWS = ['Droit français', 'Droit français - Tribunal de Commerce de Paris', 'Droit français - Tribunal Judiciaire'];

export const JURISDICTIONS = [
  'Tribunal Judiciaire de Paris',
  'Tribunal de Commerce de Paris',
  'Tribunal Judiciaire de Lyon',
  'Tribunal de Commerce de Lyon',
  'Tribunal Judiciaire de Marseille'
];

// ============================================================================
// LINE ITEM DESCRIPTIONS
// ============================================================================

export const SERVICE_LINE_ITEMS = [
  { description: 'Enquête préliminaire et constitution de dossier', unitPrice: 450 },
  { description: 'Journée de surveillance terrain', unitPrice: 380 },
  { description: 'Recherche d\'informations et vérification background', unitPrice: 280 },
  { description: 'Interview de témoins', unitPrice: 150 },
  { description: 'Rédaction de rapport d\'enquête', unitPrice: 220 },
  { description: 'Analyse documentaire et expertise', unitPrice: 180 },
  { description: 'Photographie et documentation preuves', unitPrice: 120 },
  { description: 'Déplacements et frais kilométriques', unitPrice: 65 },
  { description: 'Recherche filiation et vérification actes', unitPrice: 320 },
  { description: 'Surveillance vidéo supplémentaire', unitPrice: 280 },
  { description: 'Consultation expert en graphologie', unitPrice: 195 },
  { description: 'Analyse numérique et récupération données', unitPrice: 250 }
];

export const INVOICE_LINE_ITEMS = [
  ...SERVICE_LINE_ITEMS,
  { description: 'Honoraires d\'investigation (tarif journalier)', unitPrice: 450 },
  { description: 'Prime de déplacement', unitPrice: 85 },
  { description: 'Frais de dossier et administratifs', unitPrice: 75 },
  { description: 'Urgence - intervention prioritaire', unitPrice: 120 }
];

// ============================================================================
// DATE HELPERS - Create ISO date strings for specific static dates
// ============================================================================

export function isoDate(year: number, month: number, day: number): string {
  const m = String(month).padStart(2, '0');
  const d = String(day).padStart(2, '0');
  return `${year}-${m}-${d}T00:00:00.000Z`;
}

// Common static dates for consistent mock data
export const DATES = {
  // 2023 dates
  jan2023_15: isoDate(2023, 1, 15),
  feb2023_28: isoDate(2023, 2, 28),
  mar2023_10: isoDate(2023, 3, 10),
  apr2023_22: isoDate(2023, 4, 22),
  may2023_05: isoDate(2023, 5, 5),
  jun2023_18: isoDate(2023, 6, 18),
  jul2023_30: isoDate(2023, 7, 30),
  aug2023_12: isoDate(2023, 8, 12),
  sep2023_25: isoDate(2023, 9, 25),
  oct2023_08: isoDate(2023, 10, 8),
  nov2023_20: isoDate(2023, 11, 20),
  dec2023_05: isoDate(2023, 12, 5),
  dec2023_18: isoDate(2023, 12, 18),

  // 2024 dates
  jan2024_10: isoDate(2024, 1, 10),
  jan2024_25: isoDate(2024, 1, 25),
  jan2024_30: isoDate(2024, 1, 30),
  feb2024_14: isoDate(2024, 2, 14),
  feb2024_20: isoDate(2024, 2, 20),
  feb2024_28: isoDate(2024, 2, 28),
  mar2024_05: isoDate(2024, 3, 5),
  mar2024_10: isoDate(2024, 3, 10),
  mar2024_22: isoDate(2024, 3, 22),
  mar2024_25: isoDate(2024, 3, 25),
  apr2024_08: isoDate(2024, 4, 8),
  apr2024_15: isoDate(2024, 4, 15),
  apr2024_18: isoDate(2024, 4, 18),
  apr2024_22: isoDate(2024, 4, 22),
  apr2024_25: isoDate(2024, 4, 25),
  may2024_02: isoDate(2024, 5, 2),
  may2024_10: isoDate(2024, 5, 10),
  may2024_15: isoDate(2024, 5, 15),
  may2024_20: isoDate(2024, 5, 20),
  jun2024_01: isoDate(2024, 6, 1),
  jun2024_05: isoDate(2024, 6, 5),
  jun2024_10: isoDate(2024, 6, 10),
  jun2024_15: isoDate(2024, 6, 15),
  jun2024_20: isoDate(2024, 6, 20),
  jun2024_25: isoDate(2024, 6, 25),
  jul2024_01: isoDate(2024, 7, 1),
  jul2024_10: isoDate(2024, 7, 10),
  jul2024_15: isoDate(2024, 7, 15),
  jul2024_20: isoDate(2024, 7, 20),
  jul2024_25: isoDate(2024, 7, 25),
  aug2024_05: isoDate(2024, 8, 5),
  aug2024_10: isoDate(2024, 8, 10),
  aug2024_15: isoDate(2024, 8, 15),
  aug2024_18: isoDate(2024, 8, 18),
  aug2024_20: isoDate(2024, 8, 20),
  aug2024_25: isoDate(2024, 8, 25),
  sep2024_01: isoDate(2024, 9, 1),
  sep2024_12: isoDate(2024, 9, 12),
  sep2024_18: isoDate(2024, 9, 18),
  oct2024_01: isoDate(2024, 10, 1),
  oct2024_05: isoDate(2024, 10, 5),
  oct2024_08: isoDate(2024, 10, 8),
  oct2024_20: isoDate(2024, 10, 20),
  oct2024_25: isoDate(2024, 10, 25),
  nov2024_01: isoDate(2024, 11, 1),
  nov2024_05: isoDate(2024, 11, 5),
  nov2024_22: isoDate(2024, 11, 22),
  dec2024_01: isoDate(2024, 12, 1),
  dec2024_10: isoDate(2024, 12, 10),
  dec2024_20: isoDate(2024, 12, 20),

  // Due dates and validity periods (future dates)
  jan2025_01: isoDate(2025, 1, 1),
  jan2025_15: isoDate(2025, 1, 15),
  feb2025_01: isoDate(2025, 2, 1),
  feb2025_28: isoDate(2025, 2, 28),
  mar2025_01: isoDate(2025, 3, 1),
  mar2025_15: isoDate(2025, 3, 15),
  apr2025_01: isoDate(2025, 4, 1),
  may2025_01: isoDate(2025, 5, 1),
  jun2025_01: isoDate(2025, 6, 1),
  jul2025_01: isoDate(2025, 7, 1),
  aug2025_01: isoDate(2025, 8, 1),
  sep2025_01: isoDate(2025, 9, 1),
  oct2025_01: isoDate(2025, 10, 1),
  nov2025_01: isoDate(2025, 11, 1),
  dec2025_01: isoDate(2025, 12, 1),
} as const;

// ============================================================================
// COMPANY NAME GENERATION
// ============================================================================

export const COMPANY_NAMES = [
  'Assurance Mutuelle', 'Compagnie d\'Assurance', 'Cabinet d\'Avocats',
  'Société Anonyme', 'Entreprise Générale', 'Holding Group',
  'Services Conseils', 'Tech Solutions', 'Import Export',
  'Bâtiments Travaux Publics', 'Logistics Express', 'Digital Factory'
];

// ============================================================================
// CASE DESCRIPTION TEMPLATES
// ============================================================================

export const CASE_DESCRIPTIONS = [
  'Enquête de filiation pour déterminer la paternité présumée. Le client demande une investigation complète sur les antécédents familiaux et la localisation des membres de la famille.',
  'Surveillance discrète d\'un employé suspecté d\'abus de congé maladie. Le client souhaite obtenir des preuves photographiques et vidéo.',
  'Recherche d\'une personne disparue depuis 3 mois. Derniers signalements dans la région lyonnaise. Enquête terrain nécessaire.',
  'Enquête pour fraude à l\'assurance. Le client suspecte une exagération des dommages déclarés suite à un accident de voiture.',
  'Vérification des antécédents professionnels d\'un candidat à un poste de direction. Le client souhaite confirmer les diplômes et expériences.',
  'Conflit conjugal avec suspicion d\'infidélité. Le client demande une surveillance pour confirmer ou infirmer ses soupçons.',
  'Accusation abusive de harcèlement au travail. Le client conteste les allégations et souhaite une enquête indépendante.',
  'Vol de propriété intellectuelle suspecté. Le client pense qu\'un ancien employé a utilisé des documents confidentiels.',
  'Harcèlement au travail présumé. Enquête pour établir les faits et recueillir les témoignages des collègues.',
  'Infraction au code de la route documentée. Le client souhaite des preuves suite à un accident avec délit de fuite.'
];

// ============================================================================
// NOTES TEMPLATES
// ============================================================================

export const SUBJECT_NOTES = [
  'Personne coopérative, disponible pour interviews supplémentaires',
  'Refuse de communiquer certaines informations personnelles',
  'A des antécédents judiciaires connus',
  'Déménagé récemment, adresse actuelle à vérifier',
  'Très réservé, demande confidentialité',
  'A fourni des documents d\'identité valides',
  'Situation professionnelle précaire',
  'Retraité, disponible sur rendez-vous',
  'Sans emploi fixe actuellement',
  'Exerce une profession libérale'
];

export const SPECIAL_INSTRUCTIONS = [
  'Intervention uniquement en heures ouvrables',
  'Discrétion absolue requise - sujet médiatisé',
  'Coordination avec les autorités locales',
  'Priorité haute - dossier sensible',
  'Ne pas contacter directement le sujet',
  'Documents confidentiels à manipuler avec précaution',
  'Enregistrement vidéo légal obligatoire'
];

export const LOCATION_NOTES = [
  'Zone résidentielle calme, facile d\'accès',
  'Zone industrielle - sécurité accrue',
  'Accès limité - badge nécessaire',
  'Lieu public - discrétion recommandée',
  'Stationnement disponible dans la rue',
  'Zone à haute fréquentation',
  'Accès par le service des livraisons'
];

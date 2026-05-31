export const SIDEBAR_STATES = {
    Informations: 'Informations',
    Photos: 'Photos',
    Facture: 'Facture',
    Rapport: 'Rapport',
    Mandat: 'Mandat',
    Devis: 'Devis',
    Contrat: 'Contrat',
    Reseaux: 'Reseaux',
    Pièces: 'Pièces',
    Enregistrements: 'Enregistrements',
    Clients: 'Clients',
    Paramètres: 'Paramètres',
} as const;
export type SidebarState = typeof SIDEBAR_STATES[keyof typeof SIDEBAR_STATES];
export const SIDEBAR_STATES_ARRAY = Object.values(SIDEBAR_STATES);

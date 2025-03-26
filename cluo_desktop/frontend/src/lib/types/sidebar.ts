export const SIDEBAR_STATES = {
    Informations: 'Informations',
    Photos: 'Photos',
    Facture: 'Facture',
    Rapport: 'Rapport',
    Mandat: 'Mandat',
    Devis: 'Devis',
    Reseaux: 'Reseaux',
} as const;
export type SidebarState = typeof SIDEBAR_STATES[keyof typeof SIDEBAR_STATES];
export const SIDEBAR_STATES_ARRAY = Object.values(SIDEBAR_STATES);

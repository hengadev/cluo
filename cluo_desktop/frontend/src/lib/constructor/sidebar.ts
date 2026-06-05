import {
    ReceiptEuro,
    Handshake,
    Info,
    FileText,
    Globe,
    Camera,
    UserPen,
    Paperclip,
    Mic,
    Briefcase,
    Settings,
    ShieldCheck,
} from "@lucide/svelte";

import { SIDEBAR_STATES, type SidebarState } from "$lib/types/sidebar";

export type SidebarItem = {
    icon: typeof import('@lucide/svelte').Icon;
    title: SidebarState;
    path: string;  // Route path for navigation
    disabled?: boolean;
}

export const items: SidebarItem[] = [
    { icon: Info, title: SIDEBAR_STATES.Informations, path: "/cases/:id" },
    { icon: Paperclip, title: SIDEBAR_STATES.Pièces, path: "/cases/:id/pieces" },
    { icon: Mic, title: SIDEBAR_STATES.Enregistrements, path: "/cases/:id/recordings" },
    { icon: FileText, title: SIDEBAR_STATES.Rapport, path: "/cases/:id/rapport" },
    { icon: Camera, title: SIDEBAR_STATES.Photos, path: "/cases/:id/photos" },
    { icon: UserPen, title: SIDEBAR_STATES.Devis, path: "/cases/:id/documents/estimate" },
    { icon: Handshake, title: SIDEBAR_STATES.Mandat, path: "/cases/:id/documents/mandate" },
    { icon: ShieldCheck, title: SIDEBAR_STATES.Contrat, path: "/cases/:id/documents/contract" },
    { icon: ReceiptEuro, title: SIDEBAR_STATES.Facture, path: "/cases/:id/documents/facture" },
    { icon: Briefcase, title: SIDEBAR_STATES.Clients, path: "/clients" },
    { icon: Settings, title: SIDEBAR_STATES.Paramètres, path: "/settings" },
    { icon: Globe, title: SIDEBAR_STATES.Reseaux, path: "/cases/:id/reseaux", disabled: true },
]

import {
    ReceiptEuro,
    Handshake,
    Info,
    FileText,
    Globe,
    Camera,
    UserPen,
    Users,
    Paperclip,
    Mic,
} from "@lucide/svelte";

import { SIDEBAR_STATES, type SidebarState } from "$lib/types/sidebar";

export type SidebarItem = {
    icon: typeof import('@lucide/svelte').Icon;
    title: SidebarState;
    fn: (event: MouseEvent) => void
}

function handleClick() {
    console.log("here to save the space")
}

export const items: SidebarItem[] = [
    { icon: Info, title: SIDEBAR_STATES.Informations, fn: handleClick },
    { icon: Paperclip, title: SIDEBAR_STATES.Pièces, fn: handleClick },
    { icon: Users, title: SIDEBAR_STATES.Utilisateurs, fn: handleClick },
    { icon: Mic, title: SIDEBAR_STATES.Enregistrement, fn: handleClick },
    { icon: FileText, title: SIDEBAR_STATES.Rapport, fn: handleClick },
    { icon: Camera, title: SIDEBAR_STATES.Photos, fn: handleClick },
    { icon: ReceiptEuro, title: SIDEBAR_STATES.Facture, fn: handleClick },
    { icon: Handshake, title: SIDEBAR_STATES.Mandat, fn: handleClick },
    { icon: UserPen, title: SIDEBAR_STATES.Devis, fn: handleClick },
    { icon: Globe, title: SIDEBAR_STATES.Reseaux, fn: handleClick },
]

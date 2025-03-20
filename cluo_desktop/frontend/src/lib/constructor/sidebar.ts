import {
    Home,
    Slack,
    ReceiptEuro,
    Handshake,
    Info,
    FileText,
    Globe,
    Camera,
    UserPen,
} from "@lucide/svelte";

export type SidebarItem = {
    icon: typeof import('@lucide/svelte').Icon;
    title: string;
    fn: (event: MouseEvent) => void
}

function handleClick() {
    console.log("here to save the space")
}

export const items: SidebarItem[] = [
    { icon: Home, title: "Accueil", fn: handleClick },
    { icon: Info, title: "Informations", fn: handleClick },
    { icon: Slack, title: "Images", fn: handleClick },
    { icon: ReceiptEuro, title: "Facture", fn: handleClick },
    { icon: FileText, title: "Rapport", fn: handleClick },
    { icon: Handshake, title: "Mandat", fn: handleClick },
    { icon: UserPen, title: "Devis", fn: handleClick },
    { icon: Camera, title: "Photos", fn: handleClick },
    { icon: Globe, title: "Reseaux", fn: handleClick },
]

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

export type Element = {
    image: typeof import('@lucide/svelte').Icon;
    title: string;
    onClick: (event: MouseEvent) => void
}

function handleClick() {
    console.log("here to save the space")
}

export const elements: Element[] = [
    { image: Home, title: "Accueil", onClick: handleClick },
    { image: Info, title: "Informations", onClick: handleClick },
    { image: Slack, title: "Images", onClick: handleClick },
    { image: ReceiptEuro, title: "Facture", onClick: handleClick },
    { image: FileText, title: "Rapport", onClick: handleClick },
    { image: Handshake, title: "Mandat", onClick: handleClick },
    { image: UserPen, title: "Devis", onClick: handleClick },
    { image: Camera, title: "Photos", onClick: handleClick },
    { image: Globe, title: "Reseaux", onClick: handleClick },
]

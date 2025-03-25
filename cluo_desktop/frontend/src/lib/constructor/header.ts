import {
    Plus,
    Settings,
    Bell,
    MessageSquare,
} from "@lucide/svelte";

export type HeaderItem = {
    icon: typeof import('@lucide/svelte').Icon;
    title: string;
    fn: (event: MouseEvent) => void
    bg: string
    fg: string
}

function handleClick() {
    console.log("here to save the space")
}

export const items: HeaderItem[] = [
    { icon: Plus, title: "Nouvelle affaire", fn: handleClick, bg: "transparent", fg: "dark" },
    { icon: Settings, title: "Parametres", fn: handleClick, bg: "transparent", fg: "dark" },
    { icon: Bell, title: "Notifications", fn: handleClick, bg: "transparent", fg: "dark" },
    { icon: MessageSquare, title: "Chat", fn: handleClick, bg: "transparent", fg: "dark" },
]

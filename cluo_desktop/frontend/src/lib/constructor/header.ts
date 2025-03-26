import {
    Plus,
    Settings,
    Bell,
} from "@lucide/svelte";


import type { Component } from 'svelte'

export type HeaderItem = {
    icon: typeof import('@lucide/svelte').Icon;
    title: string;
    uiComponent: Component
    bg: string
    fg: string
}

import NewCase from "$lib/custom/header/NewCase.svelte"
import SettingsPopover from "$lib/custom/header/SettingsPopover.svelte"
import Notifications from "$lib/custom/header/Notifications.svelte"

export const items: HeaderItem[] = [
    { icon: Plus, title: "Creer une affaire", uiComponent: NewCase, bg: "transparent", fg: "dark" },
    { icon: Bell, title: "Voir notifications", uiComponent: Notifications, bg: "transparent", fg: "dark" },
    { icon: Settings, title: "Parametres", uiComponent: SettingsPopover, bg: "transparent", fg: "dark" },
]

import {
    Settings,
    Bell,
    Briefcase,
    Users,
} from "@lucide/svelte";

import type { Component, Snippet } from 'svelte'

export type NavItem = {
    icon: typeof import('@lucide/svelte').Icon;
    label: string;
    href: string;
}

export type UtilityItem = {
    icon: typeof import('@lucide/svelte').Icon;
    title: string;
    uiComponent: Component<{ children: Snippet }>
    bg: string
    fg: string
}

import SettingsPopover from "$lib/custom/header/SettingsPopover.svelte"
import Notifications from "$lib/custom/header/Notifications.svelte"

export const navItems: NavItem[] = [
    { icon: Briefcase, label: "Affaires", href: "/cases" },
    { icon: Users, label: "Personnes", href: "/people" },
]

export const utilityItems: UtilityItem[] = [
    { icon: Bell, title: "Voir notifications", uiComponent: Notifications, bg: "background-alt", fg: "dark" },
    { icon: Settings, title: "Parametres", uiComponent: SettingsPopover, bg: "background-alt", fg: "dark" },
]

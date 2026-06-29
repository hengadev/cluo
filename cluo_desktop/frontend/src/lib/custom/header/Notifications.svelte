<!-- Hallmark · component: popover/notifications · genre: modern-minimal · theme: project-system
     states: default · hover · focus · active · empty · all-read
     contrast: pass (foreground/muted-foreground on background) -->
<script lang="ts">
    import { Button, Popover, Separator } from "bits-ui";
    import { Mic, Sparkles, Download, CircleAlert, Send } from "@lucide/svelte";
    import { goto } from "$app/navigation";

    import { notificationStore } from "$lib/stores/notifications.svelte";
    import type { AppNotification, NotificationKind } from "$lib/types/notifications";

    type Props = { children: import("svelte").Snippet };
    let { children }: Props = $props();

    /** Popover open state — bound so click-navigation can close the panel. */
    let open = $state(false);

    /** Lucide icon for each notification category. */
    const KIND_ICON: Record<NotificationKind, typeof Mic> = {
        transcription_completed: Mic,
        transcription_failed: Mic,
        analysis_completed: Sparkles,
        invoice_overdue: CircleAlert,
        update_available: Download,
        case_released: Send,
    };

    /** Human-readable, French relative timestamp — matches the placeholder copy. */
    function formatRelativeTime(date: Date): string {
        const seconds = Math.round((Date.now() - date.getTime()) / 1000);
        if (seconds < 60) return "À l'instant";
        const minutes = Math.round(seconds / 60);
        if (minutes < 60) return `Il y a ${minutes} min`;
        const hours = Math.round(minutes / 60);
        if (hours < 24) return `Il y a ${hours} h`;
        const days = Math.round(hours / 24);
        if (days < 7) return `Il y a ${days} j`;
        return date.toLocaleDateString("fr-FR", { day: "2-digit", month: "short" });
    }
</script>

<Popover.Root bind:open>
    <Popover.Trigger>
        <div class="relative inline-flex">
            {@render children()}
            {#if notificationStore.unreadCount > 0}
                <span class="absolute -top-0.5 -right-0.5 min-w-[16px] h-4 px-1 inline-flex items-center justify-center rounded-full bg-accent text-accent-foreground text-[10px] font-semibold tabular-nums leading-none ring-2 ring-background">
                    {notificationStore.unreadCount > 9 ? "9+" : notificationStore.unreadCount}
                </span>
            {/if}
        </div>
    </Popover.Trigger>
    <Popover.Portal>
        <Popover.Content
            class="border-border-card bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 z-30 w-[400px] rounded-card-sm border overflow-hidden"
            sideOffset={8}
            align="end"
        >
            <div class="flex items-center justify-between px-5 pt-4 pb-3">
                <div class="flex items-center gap-2">
                    <h4 class="text-sm font-semibold tracking-tight text-foreground">Notifications</h4>
                    {#if notificationStore.unreadCount > 0}
                        <span class="inline-flex items-center justify-center size-4 rounded-full bg-accent text-accent-foreground text-[10px] font-semibold tabular-nums">
                            {notificationStore.unreadCount}
                        </span>
                    {/if}
                </div>
                {#if notificationStore.unreadCount > 0}
                    <Button.Root
                        onclick={() => notificationStore.markAllRead()}
                        class="text-xs font-medium text-muted-foreground hover:text-foreground transition-interactive duration-150 cursor-pointer"
                    >
                        Tout marquer comme lu
                    </Button.Root>
                {/if}
            </div>
            <Separator.Root class="bg-border-input block h-px" />
            <div class="overflow-y-auto max-h-[360px] flex flex-col divide-y divide-border-input">
                {#if notificationStore.notifications.length === 0}
                    <div class="flex flex-col items-center justify-center py-12 px-5 text-center">
                        <p class="text-sm font-medium text-foreground">Aucune notification</p>
                        <p class="text-xs text-muted-foreground mt-1">Vous êtes à jour.</p>
                    </div>
                {:else}
                    {#each notificationStore.notifications as notif (notif.id)}
                        {@render notificationItem(notif)}
                    {/each}
                {/if}
            </div>
            <Separator.Root class="bg-border-input block h-px" />
            <div class="px-5 py-3">
                <Button.Root
                    class="text-sm font-medium w-full text-center py-1 rounded-button text-muted-foreground hover:text-foreground hover:bg-muted transition-interactive duration-150 cursor-pointer"
                >
                    Voir toutes les notifications
                </Button.Root>
            </div>
        </Popover.Content>
    </Popover.Portal>
</Popover.Root>

{#snippet notificationItem(notif: AppNotification)}
    {@const Icon = KIND_ICON[notif.kind]}
    <Button.Root
        onclick={async () => {
            if (notif.caseId) {
                open = false;
                await goto(`/cases/${notif.caseId}/informations`);
            }
            notificationStore.markRead(notif.id);
        }}
        class="flex flex-col gap-1 text-left w-full px-5 py-3.5 hover:bg-muted transition-interactive duration-150 cursor-pointer"
    >
        <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-2 min-w-0">
                <span class="size-1.5 rounded-full shrink-0 {notif.read ? 'bg-transparent' : 'bg-accent'}"></span>
                <p class="text-sm {notif.read ? 'font-medium text-foreground-alt' : 'font-semibold text-foreground'} leading-snug truncate">
                    {notif.title}
                </p>
            </div>
            <p class="text-xs text-muted-foreground shrink-0">{formatRelativeTime(notif.createdAt)}</p>
        </div>
        <div class="flex items-center gap-2 pl-3.5">
            <Icon size={14} class="text-muted-foreground shrink-0" strokeWidth={1.75} />
            <p class="text-xs text-muted-foreground leading-snug">{notif.content}</p>
        </div>
    </Button.Root>
{/snippet}

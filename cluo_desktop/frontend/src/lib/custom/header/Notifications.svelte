<!-- Hallmark · component: popover/notifications · genre: modern-minimal · theme: project-system
     states: default · hover · focus · active · empty · all-read
     contrast: pass (foreground/muted-foreground on background) -->
<script lang="ts">
    import { Button, Popover, Separator } from "bits-ui";

    type Props = { children: import("svelte").Snippet };
    let { children }: Props = $props();

    type Notification = {
        title: string;
        content: string;
        sendAt: string;
        read: boolean;
    };

    let notifications: Notification[] = $state([
        {
            title: "Nouvelle fonctionnalité : Analyse IA",
            content: "Analysez les tendances de vos affaires avec notre nouvel outil d'analyse.",
            sendAt: "À l'instant",
            read: false,
        },
        {
            title: "Rapport généré avec succès",
            content: "Le rapport de l'affaire Martin est prêt à être consulté.",
            sendAt: "Il y a 2 h",
            read: false,
        },
        {
            title: "Mise à jour disponible",
            content: "La version 2.4.0 apporte des améliorations de performance.",
            sendAt: "Hier",
            read: true,
        },
        {
            title: "Nouveau document ajouté",
            content: "Un contrat a été ajouté à l'affaire Dupont & Associés.",
            sendAt: "Il y a 2 j",
            read: true,
        },
    ]);

    const unreadCount = $derived(notifications.filter(n => !n.read).length);

    function markAllRead() {
        for (const n of notifications) n.read = true;
    }
</script>

<Popover.Root>
    <Popover.Trigger>
        {@render children()}
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
                    {#if unreadCount > 0}
                        <span class="inline-flex items-center justify-center size-4 rounded-full bg-accent text-accent-foreground text-[10px] font-semibold tabular-nums">
                            {unreadCount}
                        </span>
                    {/if}
                </div>
                {#if unreadCount > 0}
                    <Button.Root
                        onclick={markAllRead}
                        class="text-xs font-medium text-muted-foreground hover:text-foreground transition-interactive duration-150 cursor-pointer"
                    >
                        Tout marquer comme lu
                    </Button.Root>
                {/if}
            </div>
            <Separator.Root class="bg-border-input block h-px" />
            <div class="overflow-y-auto max-h-[360px] flex flex-col divide-y divide-border-input">
                {#if notifications.length === 0}
                    <div class="flex flex-col items-center justify-center py-12 px-5 text-center">
                        <p class="text-sm font-medium text-foreground">Aucune notification</p>
                        <p class="text-xs text-muted-foreground mt-1">Vous êtes à jour.</p>
                    </div>
                {:else}
                    {#each notifications as notif}
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

{#snippet notificationItem(notif: Notification)}
    <Button.Root
        class="flex flex-col gap-1 text-left w-full px-5 py-3.5 hover:bg-muted transition-interactive duration-150 cursor-pointer"
    >
        <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-2 min-w-0">
                <span class="size-1.5 rounded-full shrink-0 {notif.read ? 'bg-transparent' : 'bg-accent'}"></span>
                <p class="text-sm {notif.read ? 'font-medium text-foreground-alt' : 'font-semibold text-foreground'} leading-snug truncate">
                    {notif.title}
                </p>
            </div>
            <p class="text-xs text-muted-foreground shrink-0">{notif.sendAt}</p>
        </div>
        <p class="text-xs text-muted-foreground leading-snug pl-3.5">{notif.content}</p>
    </Button.Root>
{/snippet}

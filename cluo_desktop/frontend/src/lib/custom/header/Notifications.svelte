<script lang="ts">
    import { Button, Popover, Separator } from "bits-ui";

    type Props = { children: import("svelte").Snippet };
    let { children }: Props = $props();

    type Notification = {
        title: string;
        content: string;
        sendAt: string;
    };
    let notifications: Notification[] = [
        {
            title: "New Feature: AI Case Analysis",
            content:
                "Analyse case patterns with our new AI-powered insight tool",
            sendAt: "Just now",
        },
        {
            title: "New Feature: AI Case Analysis",
            content:
                "Analyse case patterns with our new AI-powered insight tool",
            sendAt: "Just now",
        },
        {
            title: "New Feature: AI Case Analysis",
            content:
                "Analyse case patterns with our new AI-powered insight tool",
            sendAt: "Just now",
        },
        {
            title: "New Feature: AI Case Analysis",
            content:
                "Analyse case patterns with our new AI-powered insight tool",
            sendAt: "Just now",
        },
        {
            title: "New Feature: AI Case Analysis",
            content:
                "Analyse case patterns with our new AI-powered insight tool",
            sendAt: "Just now",
        },
        {
            title: "New Feature: AI Case Analysis",
            content:
                "Analyse case patterns with our new AI-powered insight tool",
            sendAt: "Just now",
        },
        {
            title: "New Feature: AI Case Analysis",
            content:
                "Analyse case patterns with our new AI-powered insight tool",
            sendAt: "Just now",
        },
        {
            title: "New Feature: AI Case Analysis",
            content:
                "Analyse case patterns with our new AI-powered insight tool",
            sendAt: "Just now",
        },
    ];
</script>

<Popover.Root>
    <Popover.Trigger>
        {@render children()}
    </Popover.Trigger>
    <Popover.Portal>
        <Popover.Content
            class=" border-dark-10 bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 z-30 w-full max-w-[480px] rounded-[12px] border"
            sideOffset={8}
            align="end"
        >
            <div class="flex items-center justify-between px-4 py-3">
                <h4 class="text-sm font-semibold tracking-tight">
                    Notifications
                </h4>
                <Button.Root
                    class="text-xs font-medium text-foreground-alt hover:text-foreground border border-border-input rounded-input px-3 h-8 hover:bg-muted cursor-pointer transition-colors"
                >
                    Marquer comme vus
                </Button.Root>
            </div>
            <Separator.Root class="bg-border-input block h-px" />
            <div
                class="overflow-y-auto max-h-[400px] h-full flex flex-col"
            >
                {#each notifications as notif}
                    {@render notificationButton(notif)}
                {/each}
            </div>
            <Separator.Root class="bg-border-input block h-px" />
            <div class="p-3">
                <Button.Root
                    class="text-sm font-medium w-full rounded-input px-3 h-9 hover:bg-muted cursor-pointer transition-colors text-foreground-alt hover:text-foreground"
                >
                    Voir toutes les notifications
                </Button.Root>
            </div>
        </Popover.Content>
    </Popover.Portal>
</Popover.Root>

{#snippet notificationButton(notification: Notification)}
    <Button.Root
        class="flex flex-col gap-1 text-left w-full px-4 py-3 bg-transparent hover:bg-muted border-b border-border-input last:border-b-0 cursor-pointer transition-colors"
    >
        <div class="flex items-start justify-between gap-4">
            <p class="text-sm font-medium leading-snug">{notification.title}</p>
            <p class="text-xs text-foreground-alt shrink-0">{notification.sendAt}</p>
        </div>
        <p class="text-xs text-foreground-alt leading-snug">{notification.content}</p>
    </Button.Root>
{/snippet}

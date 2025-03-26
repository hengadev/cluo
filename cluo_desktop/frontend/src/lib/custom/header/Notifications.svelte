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
            <div class="flex items-center justify-between p-4">
                <div class="flex flex-col">
                    <h4
                        class="text-[17px] font-semibold leading-5 tracking-[-0.01em]"
                    >
                        Notifications
                    </h4>
                </div>
                <Button.Root
                    class="border-2 border-muted p-2 rounded-input hover:bg-[#fafafa] cursor-pointer"
                >
                    Marquer comme vus
                </Button.Root>
            </div>
            <Separator.Root class="bg-dark-10 -mx-4 block h-px" />
            <div
                class="overflow-y-scroll max-h-[400px] h-full flex flex-col items-center pb-2"
            >
                {#each notifications as notif}
                    {@render notificationButton(notif)}
                {/each}
            </div>
            <Separator.Root class="bg-dark-10 -mx-4 block h-px" />
            <div class="p-4">
                <Button.Root
                    class="border-2 border-muted p-2 w-full rounded-input hover:bg-[#fafafa] cursor-pointer"
                >
                    Voir toutes les notifications
                </Button.Root>
            </div>
        </Popover.Content>
    </Popover.Portal>
</Popover.Root>

{#snippet notificationButton(notification: Notification)}
    <Button.Root
        class="grid gap-4 text-xs w-full p-2 bg-[#fafafa] hover:bg-[#f0f0f0] border-2 border-muted p-4"
    >
        <div class="grid gap-1 text-left">
            <p class="font-semibold">{notification.title}</p>
            <p class="">{notification.content}</p>
        </div>
        <p class="text-right">{notification.sendAt}</p>
    </Button.Root>
{/snippet}

<script lang="ts">
    import { Dialog, Separator } from "bits-ui";
    import { X } from "@lucide/svelte";

    type Props = {
        children: import("svelte").Snippet;
        onConfirm?: () => void | Promise<void>;
        title?: string;
        description?: string;
        confirmLabel?: string;
        cancelLabel?: string;
    };
    let { children, onConfirm, title, description, confirmLabel = "Confirmer", cancelLabel = "Annuler" }: Props = $props();

    async function handleConfirm() {
        if (onConfirm) {
            await onConfirm();
        }
    }
</script>

<Dialog.Root>
    <Dialog.Trigger>
        {@render children()}
    </Dialog.Trigger>
    <Dialog.Portal>
        <Dialog.Overlay
            class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
        />
        <Dialog.Content
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border flex flex-col sm:max-w-[480px] md:w-full"
        >
            <div class="flex-shrink-0 px-8 pt-8 pb-6">
                <Dialog.Title class="text-base font-semibold tracking-tight">
                    {title || 'Êtes-vous sûr ?'}
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt text-sm mt-1">
                    {description || 'Prenez un moment pour examiner les détails fournis afin de vous assurer que vous comprenez les implications.'}
                </Dialog.Description>
            </div>

            <Separator.Root class="bg-border-input mx-0 !m-0 block h-px flex-shrink-0" />

            <div class="flex justify-end gap-2 px-8 py-6">
                <Dialog.Close
                    class="h-input rounded-input bg-transparent text-foreground border border-border-input hover:bg-muted focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-colors"
                >
                    {cancelLabel}
                </Dialog.Close>
                <Dialog.Close
                    onclick={handleConfirm}
                    class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/90 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-all"
                >
                    {confirmLabel}
                </Dialog.Close>
            </div>

            <Dialog.Close
                class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-6 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            >
                <div>
                    <X class="text-foreground-alt size-4" />
                    <span class="sr-only">Close</span>
                </div>
            </Dialog.Close>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

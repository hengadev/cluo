<script lang="ts">
    import { Dialog } from "bits-ui";
    import { X } from "@lucide/svelte";

    type Props = { children: import("svelte").Snippet };
    let { children }: Props = $props();
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
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border p-8 sm:max-w-[480px] md:w-full"
        >
            <Dialog.Title
                class="flex w-full text-lg font-semibold tracking-tight"
            >
                Êtes-vous sûr ?
            </Dialog.Title>
            <Dialog.Description class="text-foreground-alt !mt-2 text-sm"
                >Prenez un moment pour examiner les détails fournis afin de vous
                assurer que vous comprenez les implications.</Dialog.Description
            >
            <div class="flex justify-end gap-2 !mt-6">
                <Dialog.Close
                    class="h-input rounded-input bg-transparent text-dark hover:bg-[#fafafa] focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 py-2 text-[15px] font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] border-2 border-[#dedede] cursor-pointer"
                >
                    Cancel
                </Dialog.Close>
                <Dialog.Close
                    class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-4 py-2 text-[15px] font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
                >
                    Okay
                </Dialog.Close>
            </div>
            <Dialog.Close
                class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            >
                <div>
                    <X class="text-foreground size-5" />
                    <span class="sr-only">Close</span>
                </div>
            </Dialog.Close>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

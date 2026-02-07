<script lang="ts">
    import { Moon, Sun } from "@lucide/svelte";
    import { Button, Tooltip } from "bits-ui";
    import { theme } from "$lib/stores/theme";
</script>

<Tooltip.Provider>
    <Tooltip.Root delayDuration={100}>
        <Tooltip.Trigger
            class="rounded-10px p-3 bg-dark-50 ring-offset-background active:scale-[0.98] active:transition:all
		focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center focus-visible:ring-2 focus-visible:ring-offset-2
                    hover:bg-dark-100/50 transition-colors"
        >
            <Button.Root class="cursor-pointer" onclick={() => theme.toggle()}>
                {#if $theme === "dark"}
                    <Sun size={24} strokeWidth={1.75} />
                {:else}
                    <Moon size={24} strokeWidth={1.75} />
                {/if}
            </Button.Root>
        </Tooltip.Trigger>
        <Tooltip.Content sideOffset={8} side="bottom">
            <div
                class="rounded-input text-[1rem] align-center bg-foreground text-background font-medium border-dark-10 shadow-popover outline-hidden z-0 flex items-center justify-center border p-2"
            >
                {$theme === "dark" ? "Mode clair" : "Mode sombre"}
            </div>
        </Tooltip.Content>
    </Tooltip.Root>
</Tooltip.Provider>

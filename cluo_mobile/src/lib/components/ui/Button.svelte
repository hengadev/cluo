<script lang="ts">
    import { createEventDispatcher } from "svelte";

    interface Props {
        variant?: "primary" | "secondary" | "outline" | "ghost";
        size?: "sm" | "md" | "lg";
        disabled?: boolean;
        type?: "button" | "submit" | "reset";
        class?: string;
        children?: import("svelte").Snippet;
    }

    let {
        variant = "primary",
        size = "md",
        disabled = false,
        type = "button",
        class: className = "",
        children,
        ...restProps
    }: Props = $props();

    const dispatch = createEventDispatcher();

    const baseClasses =
        "inline-flex items-center justify-center font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 rounded-lg cursor-pointer";

    const variantClasses = {
        primary:
            "bg-foreground hover:bg-dark-800 text-background focus-visible:ring-foreground",
        secondary:
            "bg-muted hover:bg-dark-100 text-foreground focus-visible:ring-foreground",
        outline:
            "border border-border-input bg-background hover:bg-muted text-foreground focus-visible:ring-foreground",
        ghost: "hover:bg-muted text-foreground focus-visible:ring-foreground",
    };

    const sizeClasses = {
        sm: "h-9 px-3 text-sm",
        md: "h-10 px-4 py-2 text-sm",
        lg: "h-12 px-8 py-2 text-base",
    };

    const buttonClass = `${baseClasses} ${variantClasses[variant]} ${sizeClasses[size]} ${className}`;

    function handleClick(event: MouseEvent) {
        if (!disabled) {
            dispatch("click", event);
        }
    }
</script>

<button
    {type}
    {disabled}
    class={buttonClass}
    onclick={handleClick}
    {...restProps}
>
    {#if children}
        {@render children()}
    {/if}
</button>


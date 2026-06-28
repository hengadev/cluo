<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import { Eye, EyeOff } from "@lucide/svelte";

    interface Props {
        type?: "text" | "email" | "password" | "tel" | "url" | "search";
        id?: string;
        name?: string;
        placeholder?: string;
        value?: string;
        disabled?: boolean;
        readonly?: boolean;
        required?: boolean;
        class?: string;
        size?: "sm" | "md" | "lg";
        showPasswordToggle?: boolean;
    }

    let {
        type = "text",
        id,
        name,
        placeholder,
        value = $bindable(),
        disabled = false,
        readonly = false,
        required = false,
        class: className = "",
        size = "md",
        showPasswordToggle = true,
        ...restProps
    }: Props = $props();

    const dispatch = createEventDispatcher();

    let showPassword = $state(false);

    const hasToggle = $derived(type === "password" && showPasswordToggle);
    const effectiveType = $derived(hasToggle && showPassword ? "text" : type);

    const baseClasses =
        "flex w-full bg-dark-50 placeholder-dark-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-foreground focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 rounded-input";

    const sizeClasses = {
        sm: "h-9 px-3 text-base",
        md: "h-10 px-3 py-2 text-base",
        lg: "h-12 px-3 py-2 text-base",
    };

    const inputClass = $derived(
        `${baseClasses} ${sizeClasses[size]} ${hasToggle ? "pr-10" : ""} ${className}`
    );

    const iconSize = size === "sm" ? 16 : 18;

    function togglePassword() {
        showPassword = !showPassword;
    }

    function handleInput(event: Event) {
        const target = event.target as HTMLInputElement;
        dispatch("input", { value: target.value });
    }

    function handleChange(event: Event) {
        const target = event.target as HTMLInputElement;
        dispatch("change", { value: target.value });
    }

    function handleFocus(event: FocusEvent) {
        dispatch("focus", event);
    }

    function handleBlur(event: FocusEvent) {
        dispatch("blur", event);
    }
</script>

<div class="relative w-full">
    <input
        type={effectiveType}
        {id}
        {name}
        {placeholder}
        bind:value
        {disabled}
        {readonly}
        {required}
        class={inputClass}
        oninput={handleInput}
        onchange={handleChange}
        onfocus={handleFocus}
        onblur={handleBlur}
        {...restProps}
    />

    {#if hasToggle}
        <button
            type="button"
            onclick={togglePassword}
            class="absolute inset-y-0 right-0 flex w-10 items-center justify-center text-muted-foreground hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-foreground focus-visible:ring-offset-2 rounded-input cursor-pointer"
            aria-label={showPassword ? "Masquer le mot de passe" : "Afficher le mot de passe"}
            aria-pressed={showPassword}
        >
            {#if showPassword}
                <EyeOff size={iconSize} />
            {:else}
                <Eye size={iconSize} />
            {/if}
        </button>
    {/if}
</div>

<script lang="ts">
    import { createEventDispatcher } from "svelte";

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
        ...restProps
    }: Props = $props();

    const dispatch = createEventDispatcher();

    const baseClasses =
        "flex w-full bg-dark-50 placeholder-dark-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gray-400 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 rounded-input";

    const sizeClasses = {
        sm: "h-9 px-3 text-sm",
        md: "h-10 px-3 py-2 text-sm",
        lg: "h-12 px-3 py-2 text-base",
    };

    const inputClass = `${baseClasses} ${sizeClasses[size]} ${className}`;

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

<input
    {type}
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

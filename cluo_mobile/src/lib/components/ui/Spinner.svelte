<script lang="ts">
    interface Props {
        size?: "sm" | "md" | "lg";
        class?: string;
    }

    let { size = "md", class: className = "" }: Props = $props();

    const sizeClasses = {
        sm: "w-4 h-4 border-2",
        md: "w-6 h-6 border-2",
        lg: "w-8 h-8 border-3",
    };

    const spinnerClass = `${sizeClasses[size]} ${className}`;
</script>

<div
    class={`${spinnerClass} border-dark-200 border-t-dark-700 rounded-full animate-spin`}
    role="status"
    aria-label="Chargement"
></div>

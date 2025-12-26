<script lang="ts">
    import { Check, ChevronLeft, Ellipsis } from "@lucide/svelte";
    import Spinner from "$lib/components/ui/Spinner.svelte";

    import { goto } from "$app/navigation";

    // Data is passed from +page.server.ts load function
    let { data } = $props();
    const steps = data.steps;

    function goBack() {
        if (history.length > 0) history.back();
        else goto("/");
    }
</script>

<div class="min-h-screen flex flex-col gap-8 pb-24 mt-8">
    <div class="flex flex-col gap-2">
        <div class="flex items-center justify-between mb-8">
            <button onclick={goBack}>
                <ChevronLeft />
            </button>
            <button>
                <Ellipsis />
            </button>
        </div>
        <h1 class="text-dark-900 font-extrabold text-2xl">
            Traitement de l'enregistrement
        </h1>
        <p class="text-dark-600 text-base">
            Veuillez patienter pendant le traitement de votre enregistrement...
        </p>
    </div>

    <div class="flex flex-col gap-4">
        {#each steps as step, index}
            <div
                class="flex items-center gap-4 p-4 border-1 border-dark-100 rounded-2xl bg-background-alt"
            >
                <div
                    class="flex items-center justify-center w-12 h-12 rounded-full {step.status ===
                    'completed'
                        ? 'bg-green-500'
                        : 'bg-dark-50'}"
                >
                    {#if step.status === "completed"}
                        <Check class="text-white" size={24} strokeWidth={3} />
                    {:else}
                        <Spinner size="md" />
                    {/if}
                </div>

                <div class="flex-1">
                    <p
                        class="text-dark-800 font-semibold text-base {step.status ===
                        'completed'
                            ? 'line-through text-dark-500'
                            : ''}"
                    >
                        {step.title}
                    </p>
                    {#if step.status === "completed"}
                        <p class="text-green-600 text-sm">Terminé</p>
                    {:else}
                        <p class="text-dark-500 text-sm">En cours...</p>
                    {/if}
                </div>
            </div>
        {/each}
    </div>

    {#if steps.every((s) => s.status === "completed")}
        <div class="mt-4">
            <a
                href="/"
                class="flex items-center justify-center w-full px-6 py-4 bg-dark-700 hover:bg-dark-600 text-foreground rounded-xl transition-colors font-semibold no-underline"
            >
                Retour à l'accueil
            </a>
        </div>
    {/if}
</div>

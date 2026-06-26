<script lang="ts">
    import { page } from "$app/state";
    import { goto } from "$app/navigation";

    function goBack() {
        if (history.length > 1) history.back();
        else goto("/");
    }
</script>

<div class="min-h-screen flex flex-col items-center justify-center gap-6 p-6 text-center">
    <p class="text-6xl font-extrabold text-dark-200">{page.status}</p>
    <div class="flex flex-col gap-1">
        <p class="text-dark-900 font-bold text-xl">
            {page.status === 404 ? "Page introuvable" : "Une erreur est survenue"}
        </p>
        <p class="text-dark-500 text-sm">{page.error?.message ?? ""}</p>
    </div>
    <div class="flex gap-3 mt-2">
        <button
            onclick={goBack}
            class="px-5 py-3 bg-dark-800 hover:bg-dark-700 text-white rounded-xl font-semibold transition-colors"
        >
            Retour
        </button>
        <a
            href="/"
            class="px-5 py-3 border border-dark-200 text-dark-700 hover:bg-dark-50 rounded-xl font-semibold transition-colors no-underline"
        >
            Accueil
        </a>
    </div>
</div>

<script lang="ts">
	import { goto } from "$app/navigation";
	import { ArrowLeft } from "@lucide/svelte";
</script>

<div class="p-8 flex flex-col gap-6">
	<div class="flex items-center gap-4 animate-fade-in">
		<button
			class="p-2 rounded-lg hover:bg-muted transition-colors cursor-pointer"
			onclick={() => goto("/")}
			type="button"
		>
			<ArrowLeft size={20} />
		</button>
		<h1 class="text-3xl font-bold">Paramètres</h1>
	</div>

	<div class="max-w-2xl">
		<button
			class="w-full border border-border-card rounded-card p-6 bg-background hover:border-border-input-hover hover:shadow-md transition-all duration-300 cursor-pointer text-left animate-fade-in"
			style="animation-delay: 100ms;"
			onclick={() => goto("/settings/case-types")}
			type="button"
		>
			<h2 class="text-lg font-semibold text-foreground">Types d'affaire</h2>
			<p class="text-sm text-muted-foreground mt-1">
				Gérer les catégories de dossiers (surveillance, fraude, vérification, etc.)
			</p>
		</button>
	</div>
</div>

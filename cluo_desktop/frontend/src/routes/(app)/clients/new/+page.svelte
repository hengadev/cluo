<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import { ArrowLeft } from "@lucide/svelte";
	import { createClient } from "$lib/services/api";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import type { ClientType } from "$lib/types/entities";

	const toastState = getToastContext();

	let name = $state("");
	let type = $state<ClientType>("company");
	let saving = $state(false);

	const TYPE_OPTIONS: { value: ClientType; label: string }[] = [
		{ value: "person", label: "Particulier" },
		{ value: "company", label: "Entreprise" },
		{ value: "lawyer", label: "Cabinet juridique" },
		{ value: "insurance", label: "Assurance" },
		{ value: "government", label: "Administration" },
	];

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();

		const trimmed = name.trim();
		if (!trimmed) {
			toastState.add(
				TOAST_LEVELS.Warning,
				"Champ requis",
				"Le nom du client est obligatoire.",
			);
			return;
		}

		saving = true;
		try {
			const client = await createClient({ name: trimmed, type });
			toastState.add(
				TOAST_LEVELS.Info,
				"Client créé",
				`${trimmed} a été ajouté.`,
			);
			goto(`/clients/${client.id}`);
		} catch (e) {
			const msg =
				e instanceof Error ? e.message : "Erreur lors de la création";
			// Check for 409 Conflict
			if (msg.includes("409")) {
				toastState.add(
					TOAST_LEVELS.Error,
					"Conflit",
					"Un client avec ce nom existe déjà.",
				);
			} else {
				toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
			}
		} finally {
			saving = false;
		}
	}
</script>

<div class="p-8 max-w-xl flex flex-col gap-6">
	<button
		class="inline-flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
		onclick={() => goto("/people")}
	>
		<ArrowLeft size={18} />
		<span class="text-sm font-medium">Retour aux personnes</span>
	</button>

	<h1 class="text-3xl font-bold animate-fade-in">Nouveau client</h1>

	<form
		onsubmit={handleSubmit}
		class="grid gap-6 animate-fade-in"
		style="animation-delay: 100ms;"
	>
		<div class="flex flex-col gap-2">
			<label for="name" class="text-sm font-medium">Nom du client</label>
			<input
				id="name"
				type="text"
				bind:value={name}
				placeholder="Ex : AXA France, Maître Dupont..."
				required
				class="h-input rounded-input border-border-input bg-background placeholder:text-foreground-alt/50 hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2"
			/>
		</div>

		<div class="flex flex-col gap-2">
			<label for="type" class="text-sm font-medium">Type de client</label>
			<select
				id="type"
				bind:value={type}
				class="h-input rounded-input border-border-input bg-background hover:border-border-input-hover focus:ring-foreground focus:ring-offset-background focus:outline-hidden w-full px-4 text-sm focus:ring-2 focus:ring-offset-2 cursor-pointer"
			>
				{#each TYPE_OPTIONS as opt}
					<option value={opt.value}>{opt.label}</option>
				{/each}
			</select>
		</div>

		<div class="flex justify-end pt-4">
			<button
				type="submit"
				disabled={saving}
				class="h-input rounded-input bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center gap-2 px-8 text-[15px] font-semibold active:scale-[0.98] cursor-pointer transition-interactive duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
			>
				{saving ? "Création..." : "Créer le client"}
			</button>
		</div>
	</form>
</div>

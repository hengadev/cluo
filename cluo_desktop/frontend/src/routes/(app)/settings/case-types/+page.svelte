<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import {
		ArrowLeft,
		Plus,
		Pencil,
		Check,
		X,
		Trash2,
		Loader2,
	} from "@lucide/svelte";
	import {
		fetchAllCaseTypes,
		createCaseType,
		updateCaseType,
		deleteCaseType,
	} from "$lib/services/api";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import type { CaseType } from "$lib/types/entities";

	const toastState = getToastContext();

	let caseTypes: CaseType[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);

	// Inline rename state
	let editingId: string | null = $state(null);
	let editingName = $state("");
	let savingRename = $state(false);

	// New type input
	let newName = $state("");
	let creating = $state(false);

	// Delete state
	let deletingId: string | null = $state(null);
	let deleting = $state(false);

	onMount(async () => {
		await loadCaseTypes();
	});

	async function loadCaseTypes() {
		loading = true;
		error = null;
		try {
			caseTypes = await fetchAllCaseTypes();
		} catch (e) {
			error =
				e instanceof Error
					? e.message
					: "Erreur lors du chargement des types d'affaire";
		} finally {
			loading = false;
		}
	}

	function startRename(ct: CaseType) {
		editingId = ct.id;
		editingName = ct.name;
	}

	function cancelRename() {
		editingId = null;
		editingName = "";
	}

	async function saveRename() {
		if (!editingId || !editingName.trim()) return;
		savingRename = true;
		try {
			const updated = await updateCaseType(editingId, {
				name: editingName.trim(),
			});
			caseTypes = caseTypes.map((ct) =>
				ct.id === editingId ? updated : ct,
			);
			toastState.add(
				TOAST_LEVELS.Info,
				"Type mis à jour",
				`"${updated.name}" a été renommé.`,
			);
			cancelRename();
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de renommer le type",
			);
		} finally {
			savingRename = false;
		}
	}

	function handleRenameKeydown(e: KeyboardEvent) {
		if (e.key === "Enter") saveRename();
		if (e.key === "Escape") cancelRename();
	}

	async function handleCreate() {
		const name = newName.trim();
		if (!name) return;
		creating = true;
		try {
			const created = await createCaseType({ name });
			caseTypes = [...caseTypes, created];
			newName = "";
			toastState.add(
				TOAST_LEVELS.Info,
				"Type créé",
				`"${created.name}" a été ajouté.`,
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de créer le type",
			);
		} finally {
			creating = false;
		}
	}

	function handleCreateKeydown(e: KeyboardEvent) {
		if (e.key === "Enter") handleCreate();
	}

	async function confirmDelete() {
		if (!deletingId) return;
		deleting = true;
		try {
			await deleteCaseType(deletingId);
			const name = caseTypes.find((ct) => ct.id === deletingId)?.name || "";
			caseTypes = caseTypes.filter((ct) => ct.id !== deletingId);
			toastState.add(
				TOAST_LEVELS.Info,
				"Type supprimé",
				`"${name}" a été supprimé.`,
			);
		} catch (e) {
			toastState.add(
				TOAST_LEVELS.Error,
				"Erreur",
				e instanceof Error ? e.message : "Impossible de supprimer le type",
			);
		} finally {
			deleting = false;
			deletingId = null;
		}
	}
</script>

<div class="p-8 flex flex-col gap-6">
	<!-- Header -->
	<div class="flex items-center gap-4 animate-fade-in">
		<button
			class="p-2 rounded-lg hover:bg-muted transition-colors cursor-pointer"
			onclick={() => goto("/settings")}
			type="button"
		>
			<ArrowLeft size={20} />
		</button>
		<div>
			<h1 class="text-3xl font-bold">Types d'affaire</h1>
			<p class="text-sm text-muted-foreground mt-1">
				Gérer les catégories de dossiers disponibles lors de la création et l'édition d'affaires.
			</p>
		</div>
	</div>

	<div class="max-w-2xl">
		{#if loading}
			<div class="flex items-center justify-center py-12">
				<p class="text-muted-foreground">Chargement...</p>
			</div>
		{:else if error}
			<div
				class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg"
			>
				{error}
			</div>
		{:else}
			<!-- Case type list -->
			<div
				class="border border-border-card rounded-card overflow-hidden animate-fade-in"
				style="animation-delay: 100ms;"
			>
				{#each caseTypes as ct, index}
					<div
						class="flex items-center gap-3 px-5 py-3 hover:bg-muted/50 transition-colors {index > 0
							? 'border-t border-border'
							: ''}"
					>
						{#if editingId === ct.id}
							<!-- Inline rename -->
							<div class="flex items-center gap-2 flex-1">
								<input
									type="text"
									bind:value={editingName}
									onkeydown={handleRenameKeydown}
									disabled={savingRename}
									class="h-9 flex-1 rounded-input border-border-input bg-background px-3 text-sm focus:ring-foreground focus:ring-offset-background focus:outline-hidden focus:ring-2 focus:ring-offset-2"
									autofocus
								/>
								<button
									type="button"
									onclick={saveRename}
									disabled={savingRename || !editingName.trim()}
									class="p-1.5 rounded hover:bg-emerald-50 text-muted-foreground hover:text-emerald-600 transition-colors cursor-pointer disabled:opacity-50"
								>
									{#if savingRename}
										<Loader2 size={16} class="animate-spin" />
									{:else}
										<Check size={16} />
									{/if}
								</button>
								<button
									type="button"
									onclick={cancelRename}
									disabled={savingRename}
									class="p-1.5 rounded hover:bg-red-50 text-muted-foreground hover:text-red-600 transition-colors cursor-pointer"
								>
									<X size={16} />
								</button>
							</div>
						{:else}
							<!-- Display mode -->
							<span class="flex-1 text-sm text-foreground">{ct.name}</span>
							<div class="flex items-center gap-1">
								<button
									type="button"
									onclick={() => startRename(ct)}
									class="p-1.5 rounded hover:bg-muted text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
									title="Renommer"
								>
									<Pencil size={14} />
								</button>

								<ConfirmDialog
									title="Supprimer le type"
									description="Voulez-vous vraiment supprimer « {ct.name} » ? Les dossiers utilisant ce type ne seront pas affectés."
									onConfirm={confirmDelete}
								>
									<button
										type="button"
										onclick={() => (deletingId = ct.id)}
										class="p-1.5 rounded hover:bg-red-50 text-muted-foreground hover:text-red-600 transition-colors cursor-pointer"
										title="Supprimer"
									>
										<Trash2 size={14} />
									</button>
								</ConfirmDialog>
							</div>
						{/if}
					</div>
				{/each}

				{#if caseTypes.length === 0}
					<div class="px-5 py-8 text-center">
						<p class="text-sm text-muted-foreground">
							Aucun type d'affaire enregistré. Ajoutez-en un ci-dessous.
						</p>
					</div>
				{/if}

				<!-- Add new type -->
				<div
					class="flex items-center gap-3 px-5 py-3 border-t border-border bg-muted/30"
				>
					<Plus size={16} class="text-muted-foreground flex-shrink-0" />
					<input
						type="text"
						placeholder="Nouveau type d'affaire..."
						bind:value={newName}
						onkeydown={handleCreateKeydown}
						disabled={creating}
						class="h-9 flex-1 rounded-input border-transparent bg-transparent px-2 text-sm placeholder:text-muted-foreground focus:border-border-input focus:bg-background focus:ring-foreground focus:ring-offset-background focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors"
					/>
					<button
						type="button"
						onclick={handleCreate}
						disabled={creating || !newName.trim()}
						class="h-9 px-4 rounded-input bg-foreground text-background text-sm font-semibold hover:opacity-90 active:scale-[0.98] cursor-pointer disabled:opacity-50 transition-all inline-flex items-center gap-2"
					>
						{#if creating}
							<Loader2 size={14} class="animate-spin" />
						{/if}
						Ajouter
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>

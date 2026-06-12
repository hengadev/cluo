<script lang="ts">
	import { onMount } from "svelte";
	import {
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

	let editingId: string | null = $state(null);
	let editingName = $state("");
	let savingRename = $state(false);

	let newName = $state("");
	let creating = $state(false);

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
				`« ${updated.name} » a été renommé.`,
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
				`« ${created.name} » a été ajouté.`,
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
				`« ${name} » a été supprimé.`,
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

<div>
	{#if loading}
		<div class="flex items-center gap-4 py-16">
			<div class="h-px flex-1 bg-border-card"></div>
			<span class="text-[10px] font-mono uppercase tracking-[0.2em] text-muted-foreground">Chargement</span>
			<div class="h-px flex-1 bg-border-card"></div>
		</div>
	{:else if error}
		<div class="alert-error">{error}</div>
	{:else}
		<div class="grid grid-cols-2 md:grid-cols-3 gap-3">
			{#each caseTypes as ct, i (ct.id)}
				<div
					class="group relative border border-border-card rounded-card-lg overflow-hidden bg-background hover:bg-foreground hover:border-foreground transition-all duration-300 cursor-default opacity-0 animate-fade-in"
					style="animation-delay: {i * 55}ms; animation-fill-mode: forwards;"
				>
					{#if editingId === ct.id}
						<div class="p-5 flex flex-col gap-2 min-h-[130px]">
							<span class="text-[10px] font-mono uppercase tracking-[0.18em] text-muted-foreground">Renommer</span>
							<input
								type="text"
								bind:value={editingName}
								onkeydown={handleRenameKeydown}
								disabled={savingRename}
								class="font-display text-base font-semibold bg-transparent border-b border-foreground/20 focus:border-foreground pb-1 text-foreground outline-none w-full mt-1"
								autofocus
							/>
							<div class="flex gap-1 mt-auto pt-2">
								<button
									type="button"
									onclick={saveRename}
									disabled={savingRename || !editingName.trim()}
									class="p-1.5 rounded hover:bg-success/10 text-muted-foreground hover:text-success transition-colors cursor-pointer disabled:opacity-40"
								>
									{#if savingRename}
										<Loader2 size={14} class="animate-spin" />
									{:else}
										<Check size={14} />
									{/if}
								</button>
								<button
									type="button"
									onclick={cancelRename}
									disabled={savingRename}
									class="p-1.5 rounded hover:bg-destructive/10 text-muted-foreground hover:text-destructive transition-colors cursor-pointer"
								>
									<X size={14} />
								</button>
							</div>
						</div>
					{:else}
						<div class="p-5 min-h-[130px] flex flex-col justify-end select-none">
							<p class="font-display font-semibold text-foreground group-hover:text-background transition-colors duration-300 leading-snug">
								{ct.name}
							</p>
						</div>
						<div class="absolute top-3 right-3 flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
							<button
								type="button"
								onclick={() => startRename(ct)}
								class="p-1.5 rounded-md hover:bg-background/15 text-background/50 hover:text-background transition-colors cursor-pointer"
								title="Renommer"
							>
								<Pencil size={13} />
							</button>
							<ConfirmDialog
								title="Supprimer le type"
								description="Voulez-vous vraiment supprimer « {ct.name} » ? Les dossiers utilisant ce type ne seront pas affectés."
								onConfirm={confirmDelete}
							>
								<button
									type="button"
									onclick={() => (deletingId = ct.id)}
									class="p-1.5 rounded-md hover:bg-background/10 text-background/50 hover:text-destructive transition-colors cursor-pointer"
									title="Supprimer"
								>
									<Trash2 size={13} />
								</button>
							</ConfirmDialog>
						</div>
					{/if}
				</div>
			{/each}

			<!-- New type card — always last in grid -->
			<div
				class="border-2 border-dashed border-border-card rounded-card-lg p-5 flex flex-col gap-2 hover:border-foreground/20 transition-colors duration-200 opacity-0 animate-fade-in min-h-[130px]"
				style="animation-delay: {caseTypes.length * 55}ms; animation-fill-mode: forwards;"
			>
				<div class="flex items-center gap-1.5">
					<Plus size={12} class="text-muted-foreground" />
					<span class="text-[10px] font-mono uppercase tracking-[0.18em] text-muted-foreground">Nouveau type</span>
				</div>
				<input
					type="text"
					placeholder="Nom..."
					bind:value={newName}
					onkeydown={handleCreateKeydown}
					disabled={creating}
					class="font-display font-semibold text-sm bg-transparent border-b border-transparent focus:border-foreground/25 pb-1 text-foreground placeholder:text-muted-foreground/30 outline-none flex-1 w-full transition-colors"
				/>
				<button
					type="button"
					onclick={handleCreate}
					disabled={creating || !newName.trim()}
					class="mt-auto self-start h-8 px-3.5 rounded-input bg-foreground text-background text-xs font-semibold hover:opacity-90 active:scale-[0.97] cursor-pointer disabled:opacity-40 transition-interactive inline-flex items-center gap-1.5"
				>
					{#if creating}
						<Loader2 size={12} class="animate-spin" />
					{:else}
						<Plus size={12} />
					{/if}
					Créer
				</button>
			</div>
		</div>

		{#if caseTypes.length === 0}
			<p class="text-sm text-muted-foreground/50 mt-4 text-center">
				Commencez par créer votre premier type d'affaire.
			</p>
		{/if}
	{/if}
</div>

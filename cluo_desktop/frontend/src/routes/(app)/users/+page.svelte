<script lang="ts">
	import { onMount } from "svelte";
	import { fetchAllUsers } from "$lib/services/api";
	import type { AuthUser } from "$lib/types/entities";

	let users: AuthUser[] = [];
	let loading = true;
	let error: string | null = null;

	const ROLE_LABELS: Record<string, string> = {
		admin: "Administrateur",
		investigator: "Enquêteur",
		viewer: "Lecteur"
	};

	const ROLE_BADGE_CLASSES: Record<string, string> = {
		admin: "bg-red-100 text-red-800",
		investigator: "bg-blue-100 text-blue-800",
		viewer: "bg-gray-100 text-gray-800"
	};

	onMount(async () => {
		try {
			users = await fetchAllUsers();
		} catch (e) {
			error = e instanceof Error ? e.message : "Erreur inconnue";
		} finally {
			loading = false;
		}
	});
</script>

<div class="p-8 flex flex-col gap-6">
	<div class="">
		<h1 class="text-2xl font-bold">Utilisateurs</h1>
		<p class="text-muted-foreground mt-1">Gestion des utilisateurs</p>
	</div>
		{#if loading}
		<div class="flex items-center justify-center py-12">
			<p class="text-muted-foreground">Chargement...</p>
		</div>
		{:else if error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
			{error}
		</div>
		{:else if users.length === 0}
		<div class="text-center py-12">
			<p class="text-muted-foreground">Aucun utilisateur trouvé</p>
		</div>
		{:else}
		<div class="border border-border-card rounded-lg overflow-hidden">
			<table class="w-full">
				<thead class="bg-muted">
					<tr>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
							Email
						</th>
						<th class="px-6 py-3 text-left text-xs font-medium text-muted-foreground uppercase tracking-wider">
							Rôle
						</th>
					</tr>
				</thead>
				<tbody class="bg-background divide-y divide-border">
					{#each users as user}
						<tr class="hover:bg-muted/50 transition-colors">
							<td class="px-6 py-4 whitespace-nowrap">
								<div class="text-sm text-muted-foreground">{user.email}</div>
							</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {ROLE_BADGE_CLASSES[user.role] || 'bg-gray-100 text-gray-800'}">
									{ROLE_LABELS[user.role] || user.role}
								</span>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
		{/if}
</div>

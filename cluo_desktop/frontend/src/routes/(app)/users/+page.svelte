<script lang="ts">
	import { onMount } from "svelte";
	import { User } from "@lucide/svelte";
	import { userRoleBadge } from "$lib/utils/badgeVariants";
	import { fetchAllUsers } from "$lib/services/api";
	import Spinner from "$lib/components/Spinner.svelte";
	import EmptyState from "$lib/components/EmptyState.svelte";
	import type { AuthUser, UserRole } from "$lib/types/entities";

	let users: AuthUser[] = [];
	let loading = true;
	let error: string | null = null;

	const ROLE_LABELS: Record<string, string> = {
		admin: "Administrateur",
		investigator: "Enquêteur",
		viewer: "Lecteur"
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
			<Spinner size="lg" />
		</div>
		{:else if error}
		<div class="alert-error">
			{error}
		</div>
		{:else if users.length === 0}
		<div class="py-12">
			<EmptyState icon={User} message="Aucun utilisateur trouvé" />
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
					<tbody class="bg-background divide-y divide-border-input">
						{#each users as user, i}
							<tr class="hover:bg-muted/50 transition-colors animate-fade-in" style="animation-delay: {i * 50}ms;">
								<td class="px-6 py-4 whitespace-nowrap">
									<div class="text-sm text-foreground">{user.email}</div>
								</td>
							<td class="px-6 py-4 whitespace-nowrap">
								<span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full {userRoleBadge(user.role as UserRole)}">
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

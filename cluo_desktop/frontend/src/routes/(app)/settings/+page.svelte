<script lang="ts">
	import { goto } from "$app/navigation";
	import { ArrowLeft } from "@lucide/svelte";
	import { Tabs } from "bits-ui";
	import CaseTypesTab from "./CaseTypesTab.svelte";
	import UsersTab from "./UsersTab.svelte";

	const triggers = [
		{ value: "case-types", name: "Types d'affaire" },
		{ value: "users", name: "Utilisateurs" },
	];
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

	<Tabs.Root value="case-types">
		<div class="overflow-x-auto w-full border-b border-border">
			<Tabs.List class="inline-flex items-center gap-0 text-sm font-medium bg-transparent p-0 h-auto rounded-none">
				{#each triggers as trigger}
					<Tabs.Trigger
						value={trigger.value}
						class="px-4 py-2 rounded-none bg-transparent border-b-2 mb-[-1px] data-[state=active]:border-foreground data-[state=active]:text-foreground data-[state=active]:shadow-none data-[state=inactive]:border-transparent data-[state=inactive]:text-muted-foreground data-[state=inactive]:hover:bg-transparent data-[state=inactive]:hover:text-foreground transition-colors cursor-pointer whitespace-nowrap"
					>
						{trigger.name}
					</Tabs.Trigger>
				{/each}
			</Tabs.List>
		</div>

		<Tabs.Content value="case-types" class="pt-6">
			<CaseTypesTab />
		</Tabs.Content>
		<Tabs.Content value="users" class="pt-6">
			<UsersTab />
		</Tabs.Content>
	</Tabs.Root>
</div>

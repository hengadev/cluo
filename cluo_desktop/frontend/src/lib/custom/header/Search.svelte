<script lang="ts">
	import { Combobox } from "bits-ui";
	import { Search, ChevronsDown, ChevronsUp, ChevronsUpDown, Clock } from "@lucide/svelte";
	import { goto } from "$app/navigation";
	import { currentCase, recentCases } from "$lib/stores/case";
	import type { RecentCaseEntry } from "$lib/stores/case";
	import { searchAll } from "$lib/services/api";
	import type { SearchResult, Case, Client, Contact, ClientType, CaseStatus } from "$lib/types/entities";

	const TYPE_PILL: Record<'case' | 'client' | 'contact', { label: string; cls: string }> = {
		case:    { label: 'AFFAIRE',  cls: 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-300' },
		client:  { label: 'CLIENT',   cls: 'bg-violet-100 text-violet-800 dark:bg-violet-900/30 dark:text-violet-300' },
		contact: { label: 'CONTACT',  cls: 'bg-sky-100 text-sky-800 dark:bg-sky-900/30 dark:text-sky-300' },
	};

	const STATUS_LABELS: Record<CaseStatus, string> = {
		in_progress: 'En cours',
		ready: 'Prêt',
		released: 'Clôturé',
	};

	const STATUS_CLASSES: Record<CaseStatus, string> = {
		in_progress: 'bg-blue-100 text-blue-800',
		ready: 'bg-green-100 text-green-800',
		released: 'bg-purple-100 text-purple-800',
	};

	const CLIENT_TYPE_LABELS: Record<ClientType, string> = {
		person: 'Particulier',
		insurance: 'Assurance',
		lawyer: 'Cabinet juridique',
		company: 'Entreprise',
		government: 'Administration',
	};

	let searchValue = $state("");
	let searchResults = $state<SearchResult[]>([]);
	let debounceTimer: ReturnType<typeof setTimeout> | null = null;

	const showRecent = $derived(searchValue.trim() === "");
	const recentItems = $derived($recentCases);

	function handleInput(e: Event) {
		const val = (e.currentTarget as HTMLInputElement).value;
		searchValue = val;
		if (debounceTimer) clearTimeout(debounceTimer);
		if (!val.trim()) {
			searchResults = [];
			return;
		}
		debounceTimer = setTimeout(async () => {
			searchResults = await searchAll(val);
		}, 150);
	}

	function escHtml(str: string): string {
		return str
			.replace(/&/g, '&amp;')
			.replace(/</g, '&lt;')
			.replace(/>/g, '&gt;')
			.replace(/"/g, '&quot;');
	}

	function highlightText(text: string, indices?: readonly [number, number][]): string {
		if (!indices?.length || !text) return escHtml(text ?? '');
		let out = '';
		let last = 0;
		const sorted = [...indices].sort((a, b) => a[0] - b[0]);
		for (const [start, end] of sorted) {
			out += escHtml(text.slice(last, start));
			out += `<mark class="bg-yellow-200 text-yellow-900 rounded-sm not-italic">${escHtml(text.slice(start, end + 1))}</mark>`;
			last = end + 1;
		}
		out += escHtml(text.slice(last));
		return out;
	}

	function getMatchIndices(result: SearchResult, key: string): readonly [number, number][] | undefined {
		return result.matches?.find(m => m.key === key)?.indices;
	}

	function getHighlightedText(result: SearchResult): string {
		if (result.type === 'case') {
			const c = result.item as Case;
			return highlightText(c.title, getMatchIndices(result, 'title'));
		}
		if (result.type === 'client') {
			const c = result.item as Client;
			return highlightText(c.name, getMatchIndices(result, 'name'));
		}
		const c = result.item as Contact;
		const full = `${c.firstname} ${c.lastname}`;
		return highlightText(full, getMatchIndices(result, 'fullName'));
	}

	function handleValueChange(value: string | undefined) {
		if (!value) return;
		const colon = value.indexOf(':');
		const type = value.slice(0, colon);
		const id = value.slice(colon + 1);

		if (type === 'recent') {
			const entry = recentItems.find(r => r.id === id);
			if (entry) {
				currentCase.setCase(id);
				goto(`/cases/${id}`);
			}
		} else if (type === 'case') {
			const result = searchResults.find(r => r.type === 'case' && r.item.id === id);
			if (result) {
				const c = result.item as Case;
				recentCases.push({ id: c.id, title: c.title, status: c.status });
				currentCase.setCase(c.id);
				goto(`/cases/${c.id}`);
			}
		} else if (type === 'client') {
			goto(`/clients/${id}`);
		} else if (type === 'contact') {
			const result = searchResults.find(r => r.type === 'contact' && r.item.id === id);
			if (result) {
				const contact = result.item as Contact;
				goto(`/clients/${contact.clientID}`);
			}
		}

		searchValue = "";
		searchResults = [];
	}
</script>

<Combobox.Root
	type="single"
	onValueChange={handleValueChange}
	onOpenChange={(o) => {
		if (!o) {
			searchValue = "";
			searchResults = [];
			if (debounceTimer) clearTimeout(debounceTimer);
		}
	}}
>
	<div class="relative">
		<Search
			class="text-muted-foreground absolute start-3 top-1/2 size-6 -translate-y-1/2"
		/>
		<Combobox.Input
			oninput={handleInput}
			value={searchValue}
			class="h-10 rounded-input bg-dark-50 placeholder-dark-300 focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-[600px] truncate px-11 text-base transition-colors focus:ring-2 focus:ring-offset-2 sm:text-sm"
			placeholder="Rechercher dans la base de donnée"
			aria-label="Rechercher dans la base de donnée"
		/>
		<Combobox.Trigger class="absolute end-3 top-1/2 size-6 -translate-y-1/2">
			<ChevronsUpDown class="text-muted-foreground size-6" />
		</Combobox.Trigger>
	</div>

	<Combobox.Portal>
		<Combobox.Content
			class="focus-override border-muted bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=top]:slide-in-from-bottom-2 outline-hidden z-50 max-h-96 max-h-[var(--bits-combobox-content-available-height)] w-[var(--bits-combobox-anchor-width)] min-w-[var(--bits-combobox-anchor-width)] select-none rounded-xl border px-1 py-3 data-[side=bottom]:translate-y-1 data-[side=top]:-translate-y-1"
			sideOffset={10}
		>
			<Combobox.ScrollUpButton class="flex w-full items-center justify-center py-1">
				<ChevronsUp class="size-3" />
			</Combobox.ScrollUpButton>

			<Combobox.Viewport class="p-1">
				{#if showRecent}
					{#if recentItems.length === 0}
						<span class="block px-5 py-2 text-sm text-muted-foreground">
							Aucune affaire récente.
						</span>
					{:else}
						<p class="px-4 pb-1 pt-1 text-xs font-medium text-muted-foreground uppercase tracking-wide">
							Récemment ouvertes
						</p>
						{#each recentItems as entry (entry.id)}
							<Combobox.Item
								class="rounded-button data-highlighted:bg-muted outline-hidden flex h-10 w-full select-none items-center gap-3 py-3 pl-4 pr-2 text-sm"
								value="recent:{entry.id}"
								label={entry.title}
							>
								<Clock class="size-4 text-muted-foreground shrink-0" />
								<span class="truncate flex-1">{entry.title}</span>
								<span class="shrink-0 px-1.5 py-0.5 rounded-full text-xs font-medium {STATUS_CLASSES[entry.status]}">
									{STATUS_LABELS[entry.status]}
								</span>
							</Combobox.Item>
						{/each}
					{/if}
				{:else if searchResults.length === 0}
					<span class="block px-5 py-2 text-sm text-muted-foreground">
						Aucun résultat.
					</span>
				{:else}
					{#each searchResults as result (result.type + ':' + result.item.id)}
						{@const pill = TYPE_PILL[result.type]}
						<Combobox.Item
							class="rounded-button data-highlighted:bg-muted outline-hidden flex min-h-10 w-full select-none items-center gap-2 py-2 pl-3 pr-2 text-sm"
							value="{result.type}:{result.item.id}"
							label={result.type === 'case'
								? (result.item as Case).title
								: result.type === 'client'
									? (result.item as Client).name
									: `${(result.item as Contact).firstname} ${(result.item as Contact).lastname}`}
						>
							<span class="shrink-0 px-1.5 py-0.5 rounded text-[10px] font-semibold tracking-wide {pill.cls}">
								{pill.label}
							</span>

							<span class="flex-1 truncate">
								{@html getHighlightedText(result)}
							</span>

							{#if result.type === 'case'}
								{@const c = result.item as Case}
								<span class="flex shrink-0 items-center gap-1.5 text-xs text-muted-foreground">
									{#if c.city}
										<span>{c.city}</span>
									{/if}
									<span class="px-1.5 py-0.5 rounded-full font-medium {STATUS_CLASSES[c.status]}">
										{STATUS_LABELS[c.status]}
									</span>
								</span>
							{:else if result.type === 'client'}
								{@const c = result.item as Client}
								<span class="shrink-0 text-xs text-muted-foreground">
									{CLIENT_TYPE_LABELS[c.type]}
								</span>
							{:else if result.type === 'contact' && result.clientName}
								<span class="shrink-0 text-xs text-muted-foreground">
									{result.clientName}
								</span>
							{/if}
						</Combobox.Item>
					{/each}
				{/if}
			</Combobox.Viewport>

			<Combobox.ScrollDownButton class="flex w-full items-center justify-center py-1">
				<ChevronsDown class="size-3" />
			</Combobox.ScrollDownButton>
		</Combobox.Content>
	</Combobox.Portal>
</Combobox.Root>

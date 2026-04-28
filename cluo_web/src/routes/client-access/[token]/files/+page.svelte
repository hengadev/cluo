<script lang="ts">
	import type { PageData } from './$types';
	import { buttonVariants, cardVariants, cn } from '$lib/utils/design-system';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const formatDate = (date: Date) =>
		new Intl.DateTimeFormat('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' }).format(date);
</script>

<svelte:head>
	<title>Dossier client</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-background px-4 py-16">
	<div class="w-full max-w-md">

		<div class="flex flex-col items-center gap-4 mb-10">
			<div class="w-12 h-12 border border-dark-300 flex items-center justify-center">
				<span class="font-mono text-xxs tracking-widest text-foreground">A·I·R</span>
			</div>
			<div class="text-center">
				<p class="font-display font-bold text-foreground text-4xl leading-tight">Agence d'Investigations</p>
				<p class="font-display font-bold text-foreground text-4xl leading-tight">et de Recherches</p>
			</div>
		</div>

		<div class="{cardVariants({ size: 'lg' })} w-full">

		<p class="text-xs tracking-widest uppercase text-foreground-alt mb-4">Dossier client</p>

		<h1 class="font-serif text-foreground text-2xl mb-2">
			{data.summary.title}
		</h1>

		{#if data.summary.message}
			<p class="text-foreground-alt text-sm mb-6">
				{data.summary.message}
			</p>
		{/if}

		<hr class="border-border-card my-6">

		<p class="text-xs tracking-widest uppercase text-foreground-alt mb-3">Contenu</p>

		<ul class="space-y-2 mb-6">
			{#each data.summary.filesSummary as file}
				<li class="text-foreground text-sm">
					<span class="text-foreground-alt mr-2">—</span>{file}
				</li>
			{/each}
		</ul>

		<hr class="border-border-card my-6">

		<p class="text-xs tracking-widest uppercase text-foreground-alt mb-2">Disponible jusqu'au</p>
		<p class="text-foreground-alt text-sm mb-6">
			{formatDate(data.summary.expiresAt)}
		</p>

		<a
			href="/client-access/{data.token}/download"
			class={cn(
				buttonVariants({ variant: 'default', size: 'lg' }),
				'w-full text-white'
			)}
		>
			Télécharger les fichiers
		</a>
		</div>
	</div>
</div>

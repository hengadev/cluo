<script lang="ts">
	import type { PageData } from './$types';
	import { buttonVariants, cardVariants, cn } from '$lib/utils/design-system';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const errorMessages = {
		expired: 'Ce lien a expiré.',
		invalid: 'Ce lien n\'est pas valide.',
		unavailable: 'L\'accès est temporairement indisponible.'
	};

	const formatDate = (iso: string) =>
		new Intl.DateTimeFormat('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' }).format(new Date(iso));
</script>

<svelte:head>
	<title>Accès client</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-background px-4">
	<div class="w-full max-w-md">

		<!-- Agency branding (shown on both success and error) -->
		<div class="flex flex-col items-center gap-4 mb-16">
			<div class="w-12 h-12 border border-dark-300 flex items-center justify-center">
				<span class="font-mono text-xxs tracking-widest text-foreground">A·I·R</span>
			</div>
			<div class="text-center">
				<p class="font-display font-bold text-foreground text-4xl leading-tight">Agence d'Investigations</p>
				<p class="font-display font-bold text-foreground text-4xl leading-tight">et de Recherches</p>
			</div>
		</div>

		{#if data.valid}
			<!-- Landing card -->
			<div class="{cardVariants({ size: 'lg' })} w-full">
				<p class="text-xs tracking-widest uppercase text-foreground-alt mb-4">Dossier client</p>

				<h1 class="font-serif text-foreground text-2xl mb-2">
					{data.caseData.title}
				</h1>

				{#if data.caseData.description}
					<p class="text-foreground-alt text-sm mb-6">
						{data.caseData.description}
					</p>
				{/if}

				{#if data.caseData.externalReference}
					<hr class="border-border-card my-6">
					<p class="text-xs tracking-widest uppercase text-foreground-alt mb-2">Référence</p>
					<p class="text-foreground text-sm mb-6">
						{data.caseData.externalReference}
					</p>
				{/if}

				<hr class="border-border-card my-6">

				<p class="text-xs tracking-widest uppercase text-foreground-alt mb-2">Statut</p>
				<p class="text-foreground-alt text-sm mb-2 capitalize">
					{data.caseData.status === 'released' ? 'Publié' : data.caseData.status}
				</p>

				<hr class="border-border-card my-6">

				<p class="text-xs tracking-widest uppercase text-foreground-alt mb-2">Lien valide jusqu'au</p>
				<p class="text-foreground text-sm mb-6">
					{formatDate(data.caseData.tokenExpiresAt)}
				</p>

				<a
					href="/client-access/{data.token}/files"
					class={cn(
						buttonVariants({ variant: 'default', size: 'lg' }),
						'w-full text-white'
					)}
				>
					Ouvrir le dossier
				</a>

				<a
					href="/client-access/{data.token}/download"
					class={cn(
						buttonVariants({ variant: 'outline', size: 'lg' }),
						'w-full mt-3'
					)}
					download
				>
					<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2M7 10l5 5 5-5M12 15V3" />
					</svg>
					Télécharger tout
				</a>
			</div>
		{:else}
			<!-- Error card -->
			<div class="text-center">
				<h1 class="font-serif text-foreground text-3xl mb-4">
					{errorMessages[data.error]}
				</h1>
				<p class="text-foreground-alt text-sm">
					Veuillez nous contacter si vous pensez qu'il s'agit d'une erreur.
				</p>
			</div>
		{/if}
	</div>
</div>

<script lang="ts">
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const messages = {
		expired: 'Ce lien a expiré.',
		invalid: 'Ce lien n\'est pas valide.',
		unavailable: 'L\'accès est temporairement indisponible.'
	};
</script>

<svelte:head>
	<title>Accès client</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-background px-4">
	<div class="text-center max-w-md">
		<h1 class="font-serif text-foreground text-3xl mb-4">
			{messages[data.error]}
		</h1>
		<p class="text-foreground-alt text-sm">
			Veuillez nous contacter si vous pensez qu'il s'agit d'une erreur.
		</p>
	</div>
</div>

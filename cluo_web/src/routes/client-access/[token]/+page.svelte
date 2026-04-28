<script lang="ts">
	import type { PageData } from './$types';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const messages = {
		expired: 'This link has expired.',
		invalid: 'This link is not valid.',
		unavailable: 'Access is currently unavailable.'
	};
</script>

<svelte:head>
	<title>Client Access</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-background px-4">
	<div class="text-center max-w-md">
		<h1 class="font-serif text-foreground text-3xl mb-4">
			{messages[data.error]}
		</h1>
		<p class="text-foreground-alt text-sm">
			Please contact us if you believe this is an error.
		</p>
	</div>
</div>

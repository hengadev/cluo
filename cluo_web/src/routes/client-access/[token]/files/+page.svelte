<script lang="ts">
	import type { PageData } from './$types';
	import { buttonVariants, cardVariants, cn } from '$lib/utils/design-system';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();

	const formatDate = (date: Date) =>
		new Intl.DateTimeFormat('en-US', { month: 'long', day: 'numeric', year: 'numeric' }).format(date);
</script>

<svelte:head>
	<title>Case Files</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-background px-4">
	<div class="{cardVariants({ size: 'lg' })} max-w-md w-full">
		<p class="text-xs tracking-widest uppercase text-foreground-alt mb-4">Case Files</p>

		<h1 class="font-serif text-foreground text-2xl mb-2">
			{data.summary.title}
		</h1>

		{#if data.summary.message}
			<p class="text-foreground-alt text-sm mb-6">
				{data.summary.message}
			</p>
		{/if}

		<hr class="border-border-card my-6">

		<p class="text-xs tracking-widest uppercase text-foreground-alt mb-3">Included</p>

		<ul class="space-y-2 mb-6">
			{#each data.summary.filesSummary as file}
				<li class="text-foreground text-sm">
					<span class="text-foreground-alt mr-2">—</span>{file}
				</li>
			{/each}
		</ul>

		<hr class="border-border-card my-6">

		<p class="text-xs tracking-widest uppercase text-foreground-alt mb-2">Available until</p>
		<p class="text-foreground-alt text-sm mb-6">
			{formatDate(data.summary.expiresAt)}
		</p>

		<a
			href="/client-access/{data.token}/download"
			class={cn(
				buttonVariants({ variant: 'default', size: 'lg' }),
				'w-full'
			)}
		>
			Download Files
		</a>
	</div>
</div>

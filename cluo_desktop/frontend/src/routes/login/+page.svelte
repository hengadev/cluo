<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { apiPost } from '$lib/services/apiFetch';
	import { auth } from '$lib/stores/auth';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';
	import Eye from '@lucide/svelte/icons/eye';
	import EyeOff from '@lucide/svelte/icons/eye-off';

	const MOCK_USER_ROLE = import.meta.env.VITE_MOCK_USER_ROLE as string | undefined;

	onMount(() => {
		if (MOCK_USER_ROLE) {
			auth.setUser({ id: 'mock-user', email: 'dev@cluo.local', role: MOCK_USER_ROLE as 'admin' | 'investigator' | 'viewer' });
			goto('/');
		}
	});

	let email = '';
	let password = '';
	let error = '';
	let isLoading = false;
	let showPassword = false;

	async function handleLogin(e: Event) {
		e.preventDefault();
		error = '';
		isLoading = true;

		try {
			const response = await apiPost<{ user: { id: string; email: string; role: string } }>('/auth/login', {
				email,
				password
			}, { skipRefresh: true });

			auth.setUser({
				id: response.user.id,
				email: response.user.email,
				role: response.user.role as 'admin' | 'investigator' | 'viewer'
			});

			await goto('/');
		} catch (err) {
			if (err instanceof Error) {
				error = err.message;
			} else {
				error = 'Échec de la connexion. Veuillez réessayer.';
			}
		} finally {
			isLoading = false;
		}
	}

	const inputClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors disabled:opacity-50 disabled:cursor-not-allowed";
</script>

<div class="fixed top-4 right-4 z-10">
	<ThemeToggle />
</div>

<div class="flex items-center justify-center h-screen bg-background">
	<div class="w-full max-w-[400px] bg-background-alt border border-border-card rounded-card shadow-card p-8">
		<div class="text-center mb-8">
			<h1 class="font-serif text-3xl font-normal text-foreground tracking-tight">Cluo</h1>
			<p class="text-sm text-foreground-alt mt-1">Portail d'investigation</p>
		</div>

		<form onsubmit={handleLogin} class="flex flex-col gap-4">
			<div class="flex flex-col gap-1.5">
				<label for="email" class="text-sm font-medium">Email</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					placeholder="admin@cluo.fr"
					required
					disabled={isLoading}
					class={inputClass}
				/>
			</div>

			<div class="flex flex-col gap-1.5">
				<label for="password" class="text-sm font-medium">Mot de passe</label>
				<div class="relative">
					<input
						id="password"
						type={showPassword ? 'text' : 'password'}
						bind:value={password}
						placeholder="••••••••"
						required
						disabled={isLoading}
						class="{inputClass} pr-10"
					/>
					<button
						type="button"
						class="absolute inset-y-0 right-3 flex items-center text-foreground-alt hover:text-foreground transition-colors cursor-pointer"
						onclick={() => (showPassword = !showPassword)}
						aria-label={showPassword ? 'Masquer le mot de passe' : 'Afficher le mot de passe'}
					>
						{#if showPassword}
							<EyeOff size={14} />
						{:else}
							<Eye size={14} />
						{/if}
					</button>
				</div>
			</div>

			{#if error}
				<p class="text-sm text-destructive">{error}</p>
			{/if}

			<button
				type="submit"
				disabled={isLoading}
				class="h-input w-full rounded-input bg-dark text-background text-sm font-semibold shadow-mini hover:bg-dark/90 active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed transition-all cursor-pointer mt-2"
			>
				{isLoading ? 'Connexion…' : 'Se connecter'}
			</button>
		</form>
	</div>
</div>

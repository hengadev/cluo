<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { apiPost } from '$lib/services/apiFetch';
	import { auth } from '$lib/stores/auth';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';

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
</script>

<div class="theme-toggle-corner">
	<ThemeToggle />
</div>

<div class="login-container">
	<div class="login-card">
		<h1>Cluo</h1>
		<h2>Portail d'investigation</h2>

		<form onsubmit={handleLogin}>
			<div class="form-group">
				<label for="email">Email</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					placeholder="admin@clientvault.fr"
					required
					disabled={isLoading}
				/>
			</div>

			<div class="form-group">
				<label for="password">Mot de passe</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					placeholder="••••••••"
					required
					disabled={isLoading}
				/>
			</div>

			{#if error}
				<p class="alert-error">{error}</p>
			{/if}

			<button type="submit" disabled={isLoading}>
				{isLoading ? 'Connexion…' : 'Se connecter'}
			</button>
		</form>
	</div>
</div>

<style>
	.theme-toggle-corner {
		position: fixed;
		top: 1rem;
		right: 1rem;
	}

	.login-container {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 100vh;
		background: var(--background);
	}

	.login-card {
		width: 100%;
		max-width: 400px;
		padding: 2rem;
		background: var(--background-alt);
		border: 1px solid var(--border-card);
		border-radius: var(--radius-card);
		box-shadow: var(--shadow-card);
	}

	h1 {
		font-family: var(--font-serif);
		font-size: 2rem;
		font-weight: 400;
		text-align: center;
		margin: 0 0 0.25rem 0;
		color: var(--foreground);
	}

	h2 {
		font-size: 1rem;
		font-weight: 500;
		text-align: center;
		margin: 0 0 2rem 0;
		color: var(--foreground-alt);
	}

	.form-group {
		margin-bottom: 1rem;
	}

	label {
		display: block;
		margin-bottom: 0.25rem;
		font-size: 0.875rem;
		font-weight: 500;
		color: var(--foreground);
	}

	input {
		width: 100%;
		padding: 0.625rem;
		border: 1px solid var(--border-input);
		border-radius: var(--radius-input);
		font-size: 0.875rem;
		background: var(--background);
		color: var(--foreground);
	}

	input:focus {
		outline: none;
		border-color: var(--foreground);
	}

	input:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.alert-error {
		margin: 0 0 1rem 0;
		font-size: 0.875rem;
	}

	button {
		width: 100%;
		padding: 0.625rem;
		background: var(--foreground);
		color: var(--background);
		border: none;
		border-radius: var(--radius-button);
		font-size: 0.875rem;
		font-weight: 500;
		cursor: pointer;
	}

	button:hover:not(:disabled) {
		opacity: 0.9;
	}

	button:disabled {
		opacity: 0.7;
		cursor: not-allowed;
	}
</style>

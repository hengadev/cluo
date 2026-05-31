<script lang="ts">
	import { goto } from '$app/navigation';
	import { apiPost } from '$lib/services/apiFetch';
	import { auth } from '$lib/stores/auth';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';

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
				<p class="error">{error}</p>
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
		border-radius: 8px;
		box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
	}

	h1 {
		font-size: 1.5rem;
		font-weight: 700;
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
		border-radius: 4px;
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

	.error {
		margin: 0 0 1rem 0;
		padding: 0.5rem;
		background: var(--color-error-background, #fee2e2);
		color: var(--destructive);
		border-radius: 4px;
		font-size: 0.875rem;
	}

	button {
		width: 100%;
		padding: 0.625rem;
		background: var(--foreground);
		color: var(--background);
		border: none;
		border-radius: 4px;
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

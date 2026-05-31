<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { auth } from '$lib/stores/auth';
    import { Button, Input } from "$lib/components/ui";

    const MOCK_USER_ROLE = import.meta.env.VITE_MOCK_USER_ROLE as string | undefined;
    const API_URL = import.meta.env.VITE_API_URL ?? '';

    let loginError = $state<string | null>(null);
    let isLoading = $state(false);

    onMount(() => {
        if (MOCK_USER_ROLE) {
            auth.setUser({ id: 'mock-user', email: 'dev@cluo.local', role: MOCK_USER_ROLE as 'admin' | 'investigator' | 'viewer' });
            goto('/');
        }
    });

    async function handleLogin(event: Event) {
        event.preventDefault();
        loginError = null;
        isLoading = true;

        const form = event.target as HTMLFormElement;
        const data = new FormData(form);
        const email = (data.get('email') as string) ?? '';
        const password = (data.get('password') as string) ?? '';

        try {
            // Authenticate and get session cookies
            const loginRes = await fetch(`${API_URL}/auth/login`, {
                method: 'POST',
                credentials: 'include',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password }),
            });

            if (!loginRes.ok) {
                const text = await loginRes.text().catch(() => 'Authentication failed');
                throw new Error(text);
            }

            // Fetch user identity with the newly issued session cookie
            const meRes = await fetch(`${API_URL}/auth/me`, { credentials: 'include' });
            if (!meRes.ok) throw new Error('Failed to fetch user profile');

            const user = await meRes.json() as { id: string; email: string; role: string };
            auth.setUser({
                id: user.id,
                email: user.email,
                role: user.role as 'admin' | 'investigator' | 'viewer',
            });

            goto('/');
        } catch (err) {
            loginError = err instanceof Error ? err.message : 'Login failed';
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="min-h-screen bg-white">
    <!-- Mobile: Full screen card -->
    <div
        class="md:hidden min-h-screen flex items-start justify-center pt-12 p-4"
    >
        <div class="w-full max-w-sm">
            <div class="bg-white p-4 space-y-6">
                <div class="flex justify-center">
                    <div class="w-16 h-16 rounded-full bg-muted flex items-center justify-center">
                        <svg class="w-8 h-8 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                        </svg>
                    </div>
                </div>

                <div class="text-center">
                    <h1 class="text-2xl font-semibold text-foreground mb-2">Connectez-vous à votre compte</h1>
                    <p class="text-muted-foreground">Entrez vos identifiants pour continuer</p>
                </div>

                <form class="space-y-4" onsubmit={handleLogin}>
                    {#if loginError}
                        <p class="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-lg">{loginError}</p>
                    {/if}
                    <div class="space-y-2">
                        <label for="email" class="text-sm font-medium text-gray-900">Adresse e-mail</label>
                        <Input id="email" name="email" type="email" placeholder="Entrez votre adresse e-mail" size="md" required />
                    </div>
                    <div class="space-y-2">
                        <label for="password" class="text-sm font-medium text-gray-900">Mot de passe</label>
                        <Input id="password" name="password" type="password" placeholder="Entrez votre mot de passe" size="md" required />
                    </div>
                    <Button type="submit" variant="primary" size="md" class="w-full" disabled={isLoading}>
                        {isLoading ? 'Connexion...' : 'Se connecter'}
                    </Button>
                </form>
            </div>
        </div>
    </div>

    <!-- Desktop: Split screen layout -->
    <div class="hidden md:flex min-h-screen">
        <div class="flex-1 flex items-center justify-center p-8">
            <div class="w-full max-w-md">
                <div class="space-y-8">
                    <div class="flex justify-center">
                        <div class="w-20 h-20 rounded-full bg-muted flex items-center justify-center">
                            <svg class="w-10 h-10 text-muted-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                            </svg>
                        </div>
                    </div>

                    <div class="text-center">
                        <h1 class="text-3xl font-semibold text-foreground mb-3">Connectez-vous à votre compte</h1>
                        <p class="text-muted-foreground text-lg">Entrez vos identifiants pour continuer</p>
                    </div>

                    <form class="space-y-6" onsubmit={handleLogin}>
                        {#if loginError}
                            <p class="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-lg">{loginError}</p>
                        {/if}
                        <div class="space-y-3">
                            <label for="email-desktop" class="text-base font-medium text-gray-900">Adresse e-mail</label>
                            <Input id="email-desktop" name="email" type="email" placeholder="Entrez votre adresse e-mail" size="lg" required />
                        </div>
                        <div class="space-y-3">
                            <label for="password-desktop" class="text-base font-medium text-gray-900">Mot de passe</label>
                            <Input id="password-desktop" name="password" type="password" placeholder="Entrez votre mot de passe" size="lg" required />
                        </div>
                        <Button type="submit" variant="primary" size="lg" class="w-full" disabled={isLoading}>
                            {isLoading ? 'Connexion...' : 'Se connecter'}
                        </Button>
                    </form>
                </div>
            </div>
        </div>

        <div class="flex-1 bg-gray-100 relative overflow-hidden">
            <div class="absolute inset-0 flex items-center justify-center">
                <div class="text-center text-gray-600">
                    <div class="w-24 h-24 mx-auto mb-6 rounded-full bg-gray-200 flex items-center justify-center">
                        <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                        </svg>
                    </div>
                    <h2 class="text-2xl font-semibold mb-2">Bon retour</h2>
                    <p class="text-gray-500">Votre espace de travail sécurisé vous attend</p>
                </div>
            </div>
        </div>
    </div>
</div>

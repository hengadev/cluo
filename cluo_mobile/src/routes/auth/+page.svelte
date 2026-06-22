<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { auth } from '$lib/stores/auth';
    import { Button, Input } from "$lib/components/ui";
    import { apiFetchRaw } from '$lib/api/apiFetch';

    const MOCK_USER_ROLE = import.meta.env.VITE_MOCK_USER_ROLE as string | undefined;

    let loginError = $state<string | null>(null);
    let isLoading = $state(false);

    onMount(() => {
        if (MOCK_USER_ROLE) {
            auth.setUser({ id: 'mock-user', email: 'dev@cluo.local', role: MOCK_USER_ROLE as 'admin' | 'investigator' | 'viewer', name: 'John' });
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
            const loginRes = await apiFetchRaw('/auth/login', {
                method: 'POST',
                body: JSON.stringify({ email, password }),
            });

            if (!loginRes.ok) {
                let message = 'Connexion impossible. Veuillez réessayer.';
                try {
                    const data = await loginRes.json();
                    if (data?.error) message = data.error;
                } catch {
                    // Non-JSON response (e.g. proxy or server error page) — keep the generic message.
                }
                throw new Error(message);
            }

            const meRes = await apiFetchRaw('/auth/me');
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

<div class="min-h-screen bg-background">
    <!-- Mobile: Full screen card -->
    <div
        class="md:hidden min-h-screen flex items-start justify-center pt-12 p-4"
    >
        <div class="w-full max-w-sm">
            <div class="bg-background p-4 space-y-6">
                <div class="text-center">
                    <p class="font-display text-3xl font-semibold text-foreground tracking-tight">Cluo</p>
                </div>

                <div class="text-center">
                    <h1 class="text-2xl font-semibold text-foreground mb-2">Connectez-vous à votre compte</h1>
                    <p class="text-muted-foreground">Entrez vos identifiants pour continuer</p>
                </div>

                <form class="space-y-4" onsubmit={handleLogin}>
                    {#if loginError}
                        <p class="text-sm alert-error">{loginError}</p>
                    {/if}
                    <div class="space-y-2">
                        <label for="email" class="text-sm font-medium text-foreground">Adresse e-mail</label>
                        <Input id="email" name="email" type="email" placeholder="Entrez votre adresse e-mail" size="md" required />
                    </div>
                    <div class="space-y-2">
                        <label for="password" class="text-sm font-medium text-foreground">Mot de passe</label>
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
                    <div class="text-center">
                        <p class="font-display text-4xl font-semibold text-foreground tracking-tight">Cluo</p>
                    </div>

                    <div class="text-center">
                        <h1 class="text-3xl font-semibold text-foreground mb-3">Connectez-vous à votre compte</h1>
                        <p class="text-muted-foreground text-lg">Entrez vos identifiants pour continuer</p>
                    </div>

                    <form class="space-y-6" onsubmit={handleLogin}>
                        {#if loginError}
                            <p class="text-sm alert-error">{loginError}</p>
                        {/if}
                        <div class="space-y-3">
                            <label for="email-desktop" class="text-base font-medium text-foreground">Adresse e-mail</label>
                            <Input id="email-desktop" name="email" type="email" placeholder="Entrez votre adresse e-mail" size="lg" required />
                        </div>
                        <div class="space-y-3">
                            <label for="password-desktop" class="text-base font-medium text-foreground">Mot de passe</label>
                            <Input id="password-desktop" name="password" type="password" placeholder="Entrez votre mot de passe" size="lg" required />
                        </div>
                        <Button type="submit" variant="primary" size="lg" class="w-full" disabled={isLoading}>
                            {isLoading ? 'Connexion...' : 'Se connecter'}
                        </Button>
                    </form>
                </div>
            </div>
        </div>

        <div class="flex-1 bg-muted relative overflow-hidden flex items-center justify-center">
            <p class="font-display text-7xl font-semibold text-foreground/20 tracking-tight select-none">Cluo</p>
        </div>
    </div>
</div>

<script lang="ts">
    import "../app.css";
    import { page } from "$app/state";
    import { onMount } from 'svelte';
    import { dev } from '$app/environment';
    import { goto } from '$app/navigation';
    import { PUBLIC_APP_ENV } from '$env/static/public';

    import Footer from "./Footer.svelte";
    import Snackbar from "$lib/components/Snackbar.svelte";
    import { auth } from '$lib/stores/auth';

    const API_URL = import.meta.env.VITE_API_URL ?? '';
    const MOCK_MODE = import.meta.env.VITE_MOCK_MODE === 'true';
    const MOCK_USER_ROLE = import.meta.env.VITE_MOCK_USER_ROLE as string | undefined;
    const isStaging = PUBLIC_APP_ENV === 'staging';

    let { children } = $props();

    const showFooter = $derived(!page.url.pathname.startsWith("/auth"));

    onMount(async () => {
        if (!dev && 'serviceWorker' in navigator) {
            navigator.serviceWorker.register('/service-worker.js');
        }

        if (MOCK_MODE) {
            if (MOCK_USER_ROLE) {
                auth.setUser({ id: 'mock-user', email: 'dev@cluo.local', role: MOCK_USER_ROLE as 'admin' | 'investigator' | 'viewer', name: 'John' });
            }
            return;
        }

        // Skip if already heading to auth
        if (page.url.pathname.startsWith('/auth')) return;

        try {
            const res = await fetch(`${API_URL}/auth/me`, { credentials: 'include' });
            if (res.ok) {
                const user = await res.json() as { id: string; email: string; role: string };
                const namePart = user.email.split('@')[0];
                const displayName = namePart.charAt(0).toUpperCase() + namePart.slice(1);
                auth.setUser({
                    id: user.id,
                    email: user.email,
                    role: user.role as 'admin' | 'investigator' | 'viewer',
                    name: displayName,
                });
            } else {
                goto('/auth');
            }
        } catch {
            goto('/auth');
        }
    });
</script>

<svelte:head>
    {#if isStaging}
        <meta name="apple-mobile-web-app-title" content="Cluo [Staging]" />
        <link rel="apple-touch-icon" href="/icon-staging-192.png" />
    {/if}
</svelte:head>

<div class="relative px-4">
    {@render children()}
    {#if showFooter}
        <div class="fixed bottom-0 inset-x-0 z-50">
            <Footer />
        </div>
    {/if}
</div>

<Snackbar />

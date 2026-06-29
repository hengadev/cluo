<script lang="ts">
    import "../app.css";
    import { page } from "$app/state";
    import { onMount } from "svelte";
    import { dev } from "$app/environment";
    import { goto } from "$app/navigation";
    import { PUBLIC_APP_ENV } from "$env/static/public";

    import Footer from "./Footer.svelte";
    import Snackbar from "$lib/components/Snackbar.svelte";
    import { auth } from "$lib/stores/auth";
    import { currentCase } from "$lib/stores/current-case";
    import { apiFetchRaw } from "$lib/api/apiFetch";
    const MOCK_MODE = import.meta.env.VITE_MOCK_MODE === "true";
    const MOCK_USER_ROLE = import.meta.env.VITE_MOCK_USER_ROLE as
        | string
        | undefined;
    const isStaging = PUBLIC_APP_ENV === "staging";

    let { children } = $props();

    const showFooter = $derived(
        !page.url.pathname.startsWith("/auth") &&
        !page.url.pathname.startsWith("/recording/") &&
        !page.url.pathname.startsWith("/processing/")
    );

    onMount(async () => {
        if (!dev && "serviceWorker" in navigator) {
            navigator.serviceWorker.register("/service-worker.js");
        }

        if (MOCK_MODE) {
            if (MOCK_USER_ROLE) {
                auth.setUser({
                    id: "mock-user",
                    email: "dev@cluo.local",
                    role: MOCK_USER_ROLE as "admin" | "investigator" | "viewer",
                    name: "John",
                });
            }
            return;
        }

        // Skip auth for routes that don't require it (login + offline fallback)
        if (
            page.url.pathname.startsWith("/auth") ||
            page.url.pathname === "/offline"
        ) return;

        try {
            const res = await apiFetchRaw("/auth/me");
            if (res.ok) {
                const user = (await res.json()) as {
                    id: string;
                    email: string;
                    role: string;
                    name: string;
                };
                auth.setUser({
                    id: user.id,
                    email: user.email,
                    role: user.role as "admin" | "investigator" | "viewer",
                    name: user.name ?? "",
                });
            } else {
                goto("/auth");
            }
        } catch {
            // Network failure — usually offline. Don't bounce to /auth and
            // clobber whatever the service worker served (e.g. offline page);
            // only redirect when we're actually online.
            if (navigator.onLine) goto("/auth");
        }
    });
</script>

<svelte:head>
    {#if isStaging}
        <meta name="apple-mobile-web-app-title" content="Cluo [Staging]" />
        <link rel="apple-touch-icon" href="/icon-staging-180.png" sizes="180x180" />
    {/if}
</svelte:head>

<div class="min-h-screen bg-dark-900">
    <div class="bg-background rounded-b-[2rem] px-4 min-h-[calc(100dvh-11rem)]" style="padding-top: calc(env(safe-area-inset-top) + 2rem)">
        {@render children()}
    </div>
    {#if showFooter}
        <div class="fixed bottom-0 inset-x-0 z-50">
            <Footer currentCase={$currentCase} />
        </div>
    {/if}
</div>

<Snackbar />

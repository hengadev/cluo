<script lang="ts">
    import "../app.css";

    import Footer from "./Footer.svelte";

    let { children } = $props();
</script>

<div class="relative px-4">
    {@render children()}
    <div class="fixed bottom-0 inset-x-0 z-50">
        <Footer />
    </div>
</div>

<script lang="ts">
    import type { ActionData } from "./$types";
    // let { data, form } = $props();
    let { form }: ActionData = $props();

    let id: string = $state("");
    let password: string = $state("");
</script>

<form method="POST">
    <div class="field">
        <label for="id">Numero de dossier</label>
        <input type="text" name="id" id="id" bind:value={id} />
    </div>
    <div class="field">
        <label for="id">Numero de dossier</label>
        <input type="text" name="id" id="id" bind:value={password} />
    </div>
</form>
{#if form?.incorrect}
    <p>here I put something if not well done</p>
{/if}

<style>
</style>

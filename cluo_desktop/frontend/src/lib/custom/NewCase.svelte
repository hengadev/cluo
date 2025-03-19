<script lang="ts">
    import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import { Label } from "$lib/components/ui/label/index.js";
    import { Plus } from "@lucide/svelte";
    // TODO: the things that I need to start a new case
</script>

{#snippet halfField(name: string, value: string, label: string)}
    <div class="grid grid-cols-6 items-center gap-2">
        <Label for={name} class="text-center">{label}</Label>
        <Input id={name} {value} class="col-span-5" />
    </div>
{/snippet}

{#snippet fullField(name: string, value: string, label: string)}
    <div class="grid grid-cols-8 items-center">
        <Label for={name} class="text-left flex-nowrap flex-grow-1"
            >{label}</Label
        >
        <Input id={name} {value} class="col-span-7" />
    </div>
{/snippet}

<Dialog.Root>
    <Dialog.Trigger class={`${buttonVariants({ variant: "default" })}`}>
        <Button class="px-0">
            <Plus />
        </Button>
    </Dialog.Trigger>
    <Dialog.Content class="sm:max-w-[850px]">
        <Dialog.Header>
            <Dialog.Title>Creer une nouvelle enquete</Dialog.Title>
            <Dialog.Description
                >Entrer les informations concernant l'affaire. Sauvegarder quand
                vous avez fini.</Dialog.Description
            >
        </Dialog.Header>
        <div class="grid gap-8 py-4">
            <div class="grid grid-cols-2 items-center gap-4">
                {@render halfField("lastname", "DUPONT", "Nom")}
                {@render halfField("firstname", "Jean", "Prenom")}
            </div>
            {@render fullField("address1", "44 rue Aime Cesaire", "Addresse 1")}
            {@render fullField(
                "address2",
                "chez Monsieur Martinique",
                "Addresse 2",
            )}
        </div>
        <Dialog.Footer>
            <Button type="submit">Enregister les changements</Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

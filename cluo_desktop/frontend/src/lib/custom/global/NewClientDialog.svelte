<script lang="ts">
    import { Dialog, Label, Separator } from "bits-ui";
    import { X, Loader2 } from "@lucide/svelte";
    import { createClient, updateClient, ConflictError } from "$lib/services/api";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import type { Client, ClientType } from "$lib/types/entities";

    const toastState = getToastContext();

    type Props = {
        children?: import("svelte").Snippet;
        open?: boolean;
        client?: Client;
        onSaved?: (client: Client) => void;
    };
    let { children, open = $bindable(false), client, onSaved }: Props = $props();

    const isEdit = $derived(!!client);

    let name = $state(client?.name ?? "");
    let type = $state<ClientType>(client?.type ?? "company");
    let saving = $state(false);

    $effect(() => {
        if (open) {
            name = client?.name ?? "";
            type = client?.type ?? "company";
        }
    });

    const TYPE_OPTIONS: { value: ClientType; label: string }[] = [
        { value: "person", label: "Particulier" },
        { value: "company", label: "Entreprise" },
        { value: "lawyer", label: "Cabinet juridique" },
        { value: "insurance", label: "Assurance" },
        { value: "government", label: "Administration" },
    ];

    const inputClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors";
    const selectClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors cursor-pointer";

    async function handleSubmit() {
        const trimmed = name.trim();
        if (!trimmed) return;

        saving = true;
        try {
            let saved: Client;
            if (isEdit && client) {
                saved = await updateClient(client.id, { name: trimmed, type });
                toastState.add(TOAST_LEVELS.Info, "Client mis à jour", "Les modifications ont été enregistrées.");
            } else {
                saved = await createClient({ name: trimmed, type });
                toastState.add(TOAST_LEVELS.Info, "Client créé", `« ${trimmed} » a été ajouté.`);
            }
            open = false;
            onSaved?.(saved);
        } catch (e) {
            if (e instanceof ConflictError) {
                toastState.add(TOAST_LEVELS.Error, "Conflit", "Un client avec ce nom existe déjà.");
            } else {
                const msg = e instanceof Error ? e.message : "Une erreur est survenue";
                toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
            }
        } finally {
            saving = false;
        }
    }
</script>

<Dialog.Root bind:open>
    {#if children}
        <Dialog.Trigger>
            {@render children()}
        </Dialog.Trigger>
    {/if}
    <Dialog.Portal>
        <Dialog.Overlay
            class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
        />
        <Dialog.Content
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border flex flex-col sm:max-w-[480px] md:w-full"
        >
            <div class="flex-shrink-0 px-8 pt-8 pb-6">
                <Dialog.Title class="text-base font-semibold tracking-tight">
                    {isEdit ? "Modifier le client" : "Nouveau client"}
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt text-sm mt-1">
                    Les champs marqués d'un <span class="text-foreground font-medium">*</span> sont obligatoires.
                </Dialog.Description>
            </div>

            <Separator.Root class="bg-border-input mx-0 !m-0 block h-px flex-shrink-0" />

            <div class="px-8 py-6">
                <form
                    class="flex flex-col gap-4"
                    onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}
                >
                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="client-name" class="text-sm font-medium">
                            Nom du client <span class="text-foreground-alt font-normal">*</span>
                        </Label.Root>
                        <input
                            id="client-name"
                            type="text"
                            bind:value={name}
                            placeholder="Ex : AXA France, Maître Dupont..."
                            required
                            class={inputClass}
                        />
                    </div>

                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="client-type" class="text-sm font-medium text-foreground-alt">
                            Type de client
                        </Label.Root>
                        <select id="client-type" bind:value={type} class={selectClass}>
                            {#each TYPE_OPTIONS as opt}
                                <option value={opt.value}>{opt.label}</option>
                            {/each}
                        </select>
                    </div>

                    <div class="flex justify-end gap-2 pt-2">
                        <Dialog.Close
                            class="h-input rounded-input bg-transparent text-foreground border border-border-input hover:bg-muted inline-flex items-center justify-center px-5 text-sm font-semibold active:scale-[0.98] cursor-pointer transition-colors"
                        >
                            Annuler
                        </Dialog.Close>
                        <button
                            type="submit"
                            disabled={saving || !name.trim()}
                            class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/90 inline-flex items-center justify-center gap-2 px-6 text-sm font-semibold active:scale-[0.98] disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer transition-interactive"
                        >
                            {#if saving}
                                <Loader2 size={14} class="animate-spin" />
                            {/if}
                            {isEdit ? "Enregistrer" : "Créer le client"}
                        </button>
                    </div>
                </form>
            </div>

            <Dialog.Close
                class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-6 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            >
                <div>
                    <X class="text-foreground-alt size-4" />
                    <span class="sr-only">Close</span>
                </div>
            </Dialog.Close>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

<script lang="ts">
    import { Dialog, Label, Separator } from "bits-ui";
    import { X, Loader2, Mail, Phone, Briefcase } from "@lucide/svelte";
    import { createContact, updateContact } from "$lib/services/api";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import type { Contact } from "$lib/types/entities";

    const toastState = getToastContext();

    type Props = {
        children?: import("svelte").Snippet;
        open?: boolean;
        contact?: Contact;
        clientId: string;
        onSaved?: (contact: Contact) => void;
    };
    let { children, open = $bindable(false), contact, clientId, onSaved }: Props = $props();

    const isEdit = $derived(!!contact);

    let firstname = $state("");
    let lastname = $state("");
    let email = $state("");
    let phone = $state("");
    let position = $state("");
    let saving = $state(false);

    $effect(() => {
        if (open) {
            firstname = contact?.firstname ?? "";
            lastname = contact?.lastname ?? "";
            email = contact?.email ?? "";
            phone = contact?.phone ?? "";
            position = contact?.position ?? "";
        }
    });

    const inputClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors";

    async function handleSubmit() {
        const trimmedLast = lastname.trim();
        const trimmedFirst = firstname.trim();
        if (!trimmedLast || !trimmedFirst) {
            toastState.add(TOAST_LEVELS.Warning, "Champs requis", "Le nom et le prénom de l'interlocuteur sont obligatoires.");
            return;
        }
        saving = true;
        try {
            let saved: Contact;
            if (isEdit && contact) {
                saved = await updateContact(contact.id, {
                    lastname: trimmedLast,
                    firstname: trimmedFirst,
                    email: email.trim(),
                    phone: phone.trim(),
                    position: position.trim(),
                });
                toastState.add(TOAST_LEVELS.Info, "Interlocuteur mis à jour", "Les modifications ont été enregistrées.");
            } else {
                saved = await createContact({
                    clientID: clientId,
                    lastname: trimmedLast,
                    firstname: trimmedFirst,
                    email: email.trim(),
                    phone: phone.trim(),
                    position: position.trim(),
                });
                toastState.add(TOAST_LEVELS.Info, "Interlocuteur ajouté", `« ${trimmedFirst} ${trimmedLast} » a été ajouté.`);
            }
            open = false;
            onSaved?.(saved);
        } catch (e) {
            toastState.add(
                TOAST_LEVELS.Error,
                "Erreur",
                e instanceof Error ? e.message : isEdit ? "Impossible de mettre à jour l'interlocuteur" : "Impossible d'ajouter l'interlocuteur",
            );
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
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border flex flex-col sm:max-w-[520px] md:w-full"
        >
            <div class="flex-shrink-0 px-8 pt-8 pb-6">
                <Dialog.Title class="text-base font-semibold tracking-tight">
                    {isEdit ? "Modifier l'interlocuteur" : "Nouvel interlocuteur"}
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
                    <div class="grid grid-cols-2 gap-3">
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="contact-firstname" class="text-sm font-medium">
                                Prénom <span class="text-foreground-alt font-normal">*</span>
                            </Label.Root>
                            <input
                                id="contact-firstname"
                                type="text"
                                bind:value={firstname}
                                required
                                class={inputClass}
                            />
                        </div>
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="contact-lastname" class="text-sm font-medium">
                                Nom <span class="text-foreground-alt font-normal">*</span>
                            </Label.Root>
                            <input
                                id="contact-lastname"
                                type="text"
                                bind:value={lastname}
                                required
                                class={inputClass}
                            />
                        </div>
                    </div>

                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="contact-position" class="text-sm font-medium text-foreground-alt flex items-center gap-1.5">
                            <Briefcase size={13} /> Poste
                        </Label.Root>
                        <input
                            id="contact-position"
                            type="text"
                            bind:value={position}
                            class={inputClass}
                        />
                    </div>

                    <div class="grid grid-cols-2 gap-3">
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="contact-email" class="text-sm font-medium text-foreground-alt flex items-center gap-1.5">
                                <Mail size={13} /> Email
                            </Label.Root>
                            <input
                                id="contact-email"
                                type="email"
                                bind:value={email}
                                class={inputClass}
                            />
                        </div>
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="contact-phone" class="text-sm font-medium text-foreground-alt flex items-center gap-1.5">
                                <Phone size={13} /> Téléphone
                            </Label.Root>
                            <input
                                id="contact-phone"
                                type="tel"
                                bind:value={phone}
                                class={inputClass}
                            />
                        </div>
                    </div>

                    <div class="flex justify-end gap-2 pt-2">
                        <Dialog.Close
                            class="h-input rounded-input bg-transparent text-foreground border border-border-input hover:bg-muted inline-flex items-center justify-center px-5 text-sm font-semibold active:scale-[0.98] cursor-pointer transition-colors"
                        >
                            Annuler
                        </Dialog.Close>
                        <button
                            type="submit"
                            disabled={saving || !lastname.trim() || !firstname.trim()}
                            class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/90 inline-flex items-center justify-center gap-2 px-6 text-sm font-semibold active:scale-[0.98] disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer transition-interactive"
                        >
                            {#if saving}
                                <Loader2 size={14} class="animate-spin" />
                            {/if}
                            {isEdit ? "Enregistrer" : "Ajouter"}
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

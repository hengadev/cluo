<script lang="ts">
    import { Dialog, Label, Separator } from "bits-ui";
    import { X, Loader2 } from "@lucide/svelte";
    import { createCaseSubject } from "$lib/services/api";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
    import type { CaseSubject } from "$lib/types/entities";

    const toastState = getToastContext();

    type Props = { children?: import("svelte").Snippet; open?: boolean; onCreated?: (subject: CaseSubject) => void };
    let { children, open = $bindable(false), onCreated }: Props = $props();
    let saving = $state(false);

    let firstname = $state("");
    let lastname = $state("");
    let email = $state("");
    let phone = $state("");
    let address1 = $state("");
    let address2 = $state("");
    let city = $state("");
    let postalCode = $state("");
    let occupation = $state("");
    let notes = $state("");

    const inputClass = "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors";
    const textareaClass = "w-full rounded-card-sm border border-border-input bg-background px-4 py-3 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors resize-none";

    async function handleSubmit() {
        if (!firstname.trim() || !lastname.trim()) return;
        saving = true;
        try {
            const created = await createCaseSubject({
                firstname: firstname.trim(),
                lastname: lastname.trim(),
                email: email.trim() || undefined,
                phone: phone.trim() || undefined,
                address1: address1.trim() || undefined,
                address2: address2.trim() || undefined,
                city: city.trim() || undefined,
                postalCode: postalCode.trim() || undefined,
                occupation: occupation.trim() || undefined,
                notes: notes.trim() || undefined,
            });
            open = false;
            resetForm();
            toastState.add(
                TOAST_LEVELS.Info,
                "Personne ajoutée",
                `${created.firstname} ${created.lastname} a été créée.`,
            );
            onCreated?.(created);
        } catch (e) {
            toastState.add(
                TOAST_LEVELS.Error,
                "Erreur",
                e instanceof Error ? e.message : "Impossible de créer la personne.",
            );
        } finally {
            saving = false;
        }
    }

    function resetForm() {
        firstname = "";
        lastname = "";
        email = "";
        phone = "";
        address1 = "";
        address2 = "";
        city = "";
        postalCode = "";
        occupation = "";
        notes = "";
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
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] border flex flex-col max-h-[90vh] sm:max-w-[680px] md:w-full"
        >
            <!-- Header -->
            <div class="flex-shrink-0 px-8 pt-8 pb-6">
                <Dialog.Title class="text-base font-semibold tracking-tight">
                    Nouvelle personne
                </Dialog.Title>
                <Dialog.Description class="text-foreground-alt text-sm mt-1">
                    Les champs marqués d'un <span class="text-foreground font-medium">*</span> sont obligatoires.
                </Dialog.Description>
            </div>

            <Separator.Root class="bg-border-input mx-0 !m-0 block h-px flex-shrink-0" />

            <div class="flex-1 min-h-0 overflow-y-auto px-8 py-6">
                <form
                    class="flex flex-col gap-4"
                    onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}
                >
                    <!-- Firstname + Lastname -->
                    <div class="grid grid-cols-2 gap-4">
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="subject-firstname" class="text-sm font-medium">
                                Prénom <span class="text-foreground-alt font-normal">*</span>
                            </Label.Root>
                            <input
                                id="subject-firstname"
                                type="text"
                                bind:value={firstname}
                                placeholder="Prénom"
                                required
                                class={inputClass}
                            />
                        </div>
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="subject-lastname" class="text-sm font-medium">
                                Nom <span class="text-foreground-alt font-normal">*</span>
                            </Label.Root>
                            <input
                                id="subject-lastname"
                                type="text"
                                bind:value={lastname}
                                placeholder="Nom"
                                required
                                class={inputClass}
                            />
                        </div>
                    </div>

                    <!-- Email + Phone -->
                    <div class="grid grid-cols-2 gap-4">
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="subject-email" class="text-sm font-medium text-foreground-alt">Email</Label.Root>
                            <input
                                id="subject-email"
                                type="email"
                                bind:value={email}
                                placeholder="email@example.com"
                                class={inputClass}
                            />
                        </div>
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="subject-phone" class="text-sm font-medium text-foreground-alt">Téléphone</Label.Root>
                            <input
                                id="subject-phone"
                                type="tel"
                                bind:value={phone}
                                placeholder="+33 6 00 00 00 00"
                                class={inputClass}
                            />
                        </div>
                    </div>

                    <!-- Occupation -->
                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="subject-occupation" class="text-sm font-medium text-foreground-alt">Profession</Label.Root>
                        <input
                            id="subject-occupation"
                            type="text"
                            bind:value={occupation}
                            placeholder="Profession (optionnel)"
                            class={inputClass}
                        />
                    </div>

                    <!-- Address -->
                    <div class="grid grid-cols-2 gap-4">
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="subject-address1" class="text-sm font-medium text-foreground-alt">Adresse</Label.Root>
                            <input
                                id="subject-address1"
                                type="text"
                                bind:value={address1}
                                placeholder="Adresse ligne 1"
                                class={inputClass}
                            />
                        </div>
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="subject-address2" class="text-sm font-medium text-foreground-alt">Adresse (suite)</Label.Root>
                            <input
                                id="subject-address2"
                                type="text"
                                bind:value={address2}
                                placeholder="Adresse ligne 2"
                                class={inputClass}
                            />
                        </div>
                    </div>

                    <!-- City + Postal code -->
                    <div class="grid grid-cols-2 gap-4">
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="subject-city" class="text-sm font-medium text-foreground-alt">Ville</Label.Root>
                            <input
                                id="subject-city"
                                type="text"
                                bind:value={city}
                                placeholder="Paris"
                                class={inputClass}
                            />
                        </div>
                        <div class="flex flex-col gap-1.5">
                            <Label.Root for="subject-postalcode" class="text-sm font-medium text-foreground-alt">Code postal</Label.Root>
                            <input
                                id="subject-postalcode"
                                type="text"
                                bind:value={postalCode}
                                placeholder="75001"
                                class={inputClass}
                            />
                        </div>
                    </div>

                    <!-- Notes -->
                    <div class="flex flex-col gap-1.5">
                        <Label.Root for="subject-notes" class="text-sm font-medium text-foreground-alt">Notes</Label.Root>
                        <textarea
                            id="subject-notes"
                            bind:value={notes}
                            placeholder="Notes supplémentaires (optionnel)"
                            rows={3}
                            class={textareaClass}
                        ></textarea>
                    </div>

                    <!-- Footer -->
                    <div class="flex justify-end gap-2 pt-2">
                        <Dialog.Close
                            class="h-input rounded-input bg-transparent text-foreground border border-border-input hover:bg-muted inline-flex items-center justify-center px-5 text-sm font-semibold active:scale-[0.98] cursor-pointer transition-colors"
                        >
                            Annuler
                        </Dialog.Close>
                        <button
                            type="submit"
                            disabled={saving || !firstname.trim() || !lastname.trim()}
                            class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/90 inline-flex items-center justify-center gap-2 px-6 text-sm font-semibold active:scale-[0.98] disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer transition-interactive"
                        >
                            {#if saving}
                                <Loader2 size={14} class="animate-spin" />
                            {/if}
                            Ajouter la personne
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

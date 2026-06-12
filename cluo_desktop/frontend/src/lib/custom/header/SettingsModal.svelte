<script lang="ts">
    import { Dialog, Label, Separator, Switch, Select } from "bits-ui";
    import {
        X,
        Building2,
        Palette,
        BrainCircuit,
        Mail,
        FileText,
        ShieldCheck,
        RefreshCw,
        ChevronsUpDown,
        ChevronsUp,
        ChevronsDown,
        Check,
    } from "@lucide/svelte";
    import { updateDialogOpen } from "$lib/stores/update";
    import { getToastContext } from "$lib/custom/global/toast/state.svelte";
    import { TOAST_LEVELS } from "$lib/custom/global/toast/type";

    type Props = { children: import("svelte").Snippet };
    let { children }: Props = $props();

    const toastState = getToastContext();

    const SAVE_LABELS: Record<string, string> = {
        agence: "l'agence",
        apparence: "l'apparence",
        ia: "l'intelligence artificielle",
        messagerie: "la messagerie",
        documents: "les documents",
        confidentialite: "la confidentialité",
    };

    function save() {
        const label = SAVE_LABELS[activeTab];
        if (label) {
            toastState.add(TOAST_LEVELS.Info, "Paramètres enregistrés", `Les paramètres de « ${label} » ont été sauvegardés.`);
        }
    }

    let open = $state(false);
    let activeTab = $state("agence");

    // Agence
    let agencyName = $state("");
    let agencyAddress = $state("");
    let agencyCity = $state("");
    let agencyPostalCode = $state("");
    let agencyPhone = $state("");
    let agencyEmail = $state("");
    let agencyLicense = $state("");
    let agencySiret = $state("");
    let agencyWebsite = $state("");

    // Apparence
    let darkMode = $state(false);
    let compactView = $state(false);
    const languages = [
        { value: "fr", label: "Français" },
        { value: "en", label: "English" },
    ];
    let language = $state("fr");
    const selectedLanguageLabel = $derived(
        languages.find((l) => l.value === language)?.label ?? "Français"
    );

    // Intelligence artificielle
    const aiProviders = [
        { value: "anthropic", label: "Anthropic (Claude)" },
        { value: "openai", label: "OpenAI (GPT)" },
    ];
    let aiProvider = $state("anthropic");
    const selectedAiProviderLabel = $derived(
        aiProviders.find((p) => p.value === aiProvider)?.label ?? "Anthropic (Claude)"
    );
    let aiApiKey = $state("");
    let showApiKey = $state(false);
    const transcriptionModels = [
        { value: "whisper-1", label: "Whisper (OpenAI)" },
        { value: "nova-2", label: "Nova-2 (Deepgram)" },
    ];
    let transcriptionModel = $state("whisper-1");
    const selectedTranscriptionLabel = $derived(
        transcriptionModels.find((m) => m.value === transcriptionModel)?.label ?? "Whisper"
    );
    let aiTranscriptionLanguage = $state("fr");
    let aiAutoSuggest = $state(true);
    let aiChatContext = $state(true);

    // Messagerie
    let smtpHost = $state("");
    let smtpPort = $state("587");
    let smtpUser = $state("");
    let smtpPassword = $state("");
    let smtpSenderName = $state("");
    let smtpSenderEmail = $state("");
    let smtpTls = $state(true);

    // Documents
    let vatRate = $state("20");
    let defaultPaymentTerms = $state("30");
    let invoicePrefix = $state("FAC");
    let estimatePrefix = $state("DEV");
    let documentHeaderText = $state("");
    let documentFooterText = $state("");
    let autoNumbering = $state(true);

    // Confidentialité
    let autoLockEnabled = $state(false);
    const autoLockDelays = [
        { value: "5", label: "5 minutes" },
        { value: "15", label: "15 minutes" },
        { value: "30", label: "30 minutes" },
        { value: "60", label: "1 heure" },
    ];
    let autoLockDelay = $state("15");
    const selectedLockDelayLabel = $derived(
        autoLockDelays.find((d) => d.value === autoLockDelay)?.label ?? "15 minutes"
    );
    let encryptLocalStorage = $state(true);
    let analyticsEnabled = $state(false);

    const tabs = [
        { id: "agence", label: "Agence", icon: Building2 },
        { id: "apparence", label: "Apparence", icon: Palette },
        { id: "ia", label: "Intelligence artificielle", icon: BrainCircuit },
        { id: "messagerie", label: "Messagerie", icon: Mail },
        { id: "documents", label: "Documents", icon: FileText },
        { id: "confidentialite", label: "Confidentialité", icon: ShieldCheck },
        { id: "application", label: "Application", icon: RefreshCw },
    ];

    const inputClass =
        "h-input w-full rounded-card-sm border border-border-input bg-background px-4 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors";
    const textareaClass =
        "w-full rounded-card-sm border border-border-input bg-background px-4 py-3 text-sm placeholder:text-foreground-alt/40 hover:border-dark-40 focus:border-dark focus:outline-none focus:ring-2 focus:ring-foreground/10 focus:ring-offset-0 transition-colors resize-none";
</script>

<Dialog.Root bind:open>
    <Dialog.Trigger>
        {@render children()}
    </Dialog.Trigger>
    <Dialog.Portal>
        <Dialog.Overlay
            class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
        />
        <Dialog.Content
            class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full translate-x-[-50%] translate-y-[-50%] border flex max-w-[calc(100%-2rem)] sm:max-w-[860px] md:w-full"
            style="height: 620px;"
        >
            <!-- Left sidebar -->
            <aside class="flex-shrink-0 w-52 border-r border-border-input flex flex-col rounded-l-card-lg overflow-hidden">
                <div class="px-5 pt-6 pb-5 flex-shrink-0">
                    <h2 class="text-sm font-semibold tracking-tight text-foreground">Paramètres</h2>
                    <p class="text-xs text-muted-foreground mt-0.5">Configuration de l'application</p>
                </div>
                <Separator.Root class="bg-border-input block h-px flex-shrink-0" />
                <nav class="flex flex-col gap-0.5 p-3 flex-1 overflow-y-auto">
                    {#each tabs as tab}
                        {@const Icon = tab.icon}
                        <button
                            onclick={() => (activeTab = tab.id)}
                            class="flex items-center gap-2.5 px-3 py-2 rounded-button text-sm transition-interactive duration-150 cursor-pointer text-left w-full
                                {activeTab === tab.id
                                ? 'bg-surface text-foreground font-medium'
                                : 'text-muted-foreground hover:text-foreground hover:bg-surface'}"
                        >
                            <Icon class="size-4 shrink-0" />
                            <span class="truncate">{tab.label}</span>
                        </button>
                    {/each}
                </nav>
            </aside>

            <!-- Right content -->
            <div class="flex-1 min-w-0 flex flex-col">
                <div class="flex-1 min-h-0 overflow-y-auto px-8 py-7">

                    <!-- Agence -->
                    {#if activeTab === "agence"}
                        <div class="flex flex-col gap-6">
                            {@render sectionHeader("Identité de l'agence", "Informations affichées sur vos documents officiels et dans le portail client.")}
                            <div class="grid grid-cols-2 gap-4">
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="agency-name" class="text-sm font-medium text-muted-foreground">Nom de l'agence</Label.Root>
                                    <input id="agency-name" type="text" bind:value={agencyName} placeholder="Agence Dupont Investigation" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="agency-license" class="text-sm font-medium text-muted-foreground">N° d'agrément</Label.Root>
                                    <input id="agency-license" type="text" bind:value={agencyLicense} placeholder="AUT-075-2198-12-25-XXXX" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="agency-siret" class="text-sm font-medium text-muted-foreground">SIRET</Label.Root>
                                    <input id="agency-siret" type="text" bind:value={agencySiret} placeholder="000 000 000 00000" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="agency-phone" class="text-sm font-medium text-muted-foreground">Téléphone</Label.Root>
                                    <input id="agency-phone" type="tel" bind:value={agencyPhone} placeholder="+33 1 00 00 00 00" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="agency-email" class="text-sm font-medium text-muted-foreground">Email professionnel</Label.Root>
                                    <input id="agency-email" type="email" bind:value={agencyEmail} placeholder="contact@agence.fr" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="agency-website" class="text-sm font-medium text-muted-foreground">Site web</Label.Root>
                                    <input id="agency-website" type="url" bind:value={agencyWebsite} placeholder="https://agence.fr" class={inputClass} />
                                </div>
                            </div>
                            {@render sectionDivider("Adresse")}
                            <div class="grid grid-cols-2 gap-4">
                                <div class="col-span-2 flex flex-col gap-1.5">
                                    <Label.Root for="agency-address" class="text-sm font-medium text-muted-foreground">Adresse</Label.Root>
                                    <input id="agency-address" type="text" bind:value={agencyAddress} placeholder="1 rue de la Paix" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="agency-city" class="text-sm font-medium text-muted-foreground">Ville</Label.Root>
                                    <input id="agency-city" type="text" bind:value={agencyCity} placeholder="Paris" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="agency-postal" class="text-sm font-medium text-muted-foreground">Code postal</Label.Root>
                                    <input id="agency-postal" type="text" bind:value={agencyPostalCode} placeholder="75001" class={inputClass} />
                                </div>
                            </div>
                        </div>

                    <!-- Apparence -->
                    {:else if activeTab === "apparence"}
                        <div class="flex flex-col gap-6">
                            {@render sectionHeader("Apparence", "Personnalisez l'interface de l'application.")}
                            <div class="flex flex-col gap-4">
                                {@render switchRow("Mode sombre", "Utiliser le thème sombre", darkMode, (v) => (darkMode = v))}
                                {@render switchRow("Vue compacte", "Réduire l'espacement des listes pour afficher plus de contenu", compactView, (v) => (compactView = v))}
                            </div>
                            {@render sectionDivider("Langue")}
                            <div class="flex flex-col gap-1.5 max-w-xs">
                                <Label.Root class="text-sm font-medium text-muted-foreground">Langue de l'interface</Label.Root>
                                {@render selectWidget(languages, language, (v) => (language = v), selectedLanguageLabel)}
                            </div>
                        </div>

                    <!-- IA -->
                    {:else if activeTab === "ia"}
                        <div class="flex flex-col gap-6">
                            {@render sectionHeader("Intelligence artificielle", "Configurez les modèles et clés d'API utilisés pour la transcription, l'analyse et l'assistance à la rédaction.")}
                            {@render sectionDivider("Fournisseur IA")}
                            <div class="flex flex-col gap-4">
                                <div class="flex flex-col gap-1.5 max-w-xs">
                                    <Label.Root class="text-sm font-medium text-muted-foreground">Fournisseur</Label.Root>
                                    {@render selectWidget(aiProviders, aiProvider, (v) => (aiProvider = v), selectedAiProviderLabel)}
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="ai-api-key" class="text-sm font-medium text-muted-foreground">Clé d'API</Label.Root>
                                    <div class="relative">
                                        <input
                                            id="ai-api-key"
                                            type={showApiKey ? "text" : "password"}
                                            bind:value={aiApiKey}
                                            placeholder="sk-…"
                                            class="{inputClass} pr-16"
                                        />
                                        <button
                                            type="button"
                                            onclick={() => (showApiKey = !showApiKey)}
                                            class="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-muted-foreground hover:text-foreground transition-interactive cursor-pointer"
                                        >
                                            {showApiKey ? "Masquer" : "Afficher"}
                                        </button>
                                    </div>
                                    <p class="text-xs text-muted-foreground">Stockée localement, jamais transmise à des serveurs tiers.</p>
                                </div>
                            </div>
                            {@render sectionDivider("Transcription audio")}
                            <div class="flex flex-col gap-4">
                                <div class="flex flex-col gap-1.5 max-w-xs">
                                    <Label.Root class="text-sm font-medium text-muted-foreground">Modèle de transcription</Label.Root>
                                    {@render selectWidget(transcriptionModels, transcriptionModel, (v) => (transcriptionModel = v), selectedTranscriptionLabel)}
                                </div>
                                <div class="flex flex-col gap-1.5 max-w-xs">
                                    <Label.Root class="text-sm font-medium text-muted-foreground">Langue des enregistrements</Label.Root>
                                    {@render selectWidget(languages, aiTranscriptionLanguage, (v) => (aiTranscriptionLanguage = v), languages.find((l) => l.value === aiTranscriptionLanguage)?.label ?? "Français")}
                                </div>
                            </div>
                            {@render sectionDivider("Rédaction & chat")}
                            <div class="flex flex-col gap-4">
                                {@render switchRow("Suggestions automatiques", "Proposer des améliorations de texte pendant la rédaction du rapport", aiAutoSuggest, (v) => (aiAutoSuggest = v))}
                                {@render switchRow("Contexte étendu pour le chat", "Inclure automatiquement les transcriptions et pièces dans le contexte du chat IA", aiChatContext, (v) => (aiChatContext = v))}
                            </div>
                        </div>

                    <!-- Messagerie -->
                    {:else if activeTab === "messagerie"}
                        <div class="flex flex-col gap-6">
                            {@render sectionHeader("Messagerie sortante", "Configuration SMTP pour l'envoi des accès portail client depuis votre propre domaine.")}
                            <div class="grid grid-cols-2 gap-4">
                                <div class="col-span-2 flex flex-col gap-1.5">
                                    <Label.Root for="smtp-host" class="text-sm font-medium text-muted-foreground">Serveur SMTP</Label.Root>
                                    <input id="smtp-host" type="text" bind:value={smtpHost} placeholder="smtp.agence.fr" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="smtp-port" class="text-sm font-medium text-muted-foreground">Port</Label.Root>
                                    <input id="smtp-port" type="text" bind:value={smtpPort} placeholder="587" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="smtp-user" class="text-sm font-medium text-muted-foreground">Nom d'utilisateur</Label.Root>
                                    <input id="smtp-user" type="text" bind:value={smtpUser} placeholder="contact@agence.fr" class={inputClass} />
                                </div>
                                <div class="col-span-2 flex flex-col gap-1.5">
                                    <Label.Root for="smtp-password" class="text-sm font-medium text-muted-foreground">Mot de passe</Label.Root>
                                    <input id="smtp-password" type="password" bind:value={smtpPassword} placeholder="••••••••" class={inputClass} />
                                </div>
                            </div>
                            {@render sectionDivider("Expéditeur")}
                            <div class="grid grid-cols-2 gap-4">
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="smtp-sender-name" class="text-sm font-medium text-muted-foreground">Nom affiché</Label.Root>
                                    <input id="smtp-sender-name" type="text" bind:value={smtpSenderName} placeholder="Jean Dupont — Détective Privé" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="smtp-sender-email" class="text-sm font-medium text-muted-foreground">Email expéditeur</Label.Root>
                                    <input id="smtp-sender-email" type="email" bind:value={smtpSenderEmail} placeholder="contact@agence.fr" class={inputClass} />
                                </div>
                            </div>
                            <div class="flex flex-col gap-4">
                                {@render switchRow("TLS / STARTTLS", "Chiffrer la connexion au serveur SMTP (recommandé)", smtpTls, (v) => (smtpTls = v))}
                            </div>
                            <div>
                                <button
                                    type="button"
                                    class="h-input rounded-input border border-border-input bg-transparent text-sm font-medium text-foreground hover:bg-surface px-5 inline-flex items-center gap-2 transition-interactive cursor-pointer active:scale-[0.98]"
                                >
                                    <Mail class="size-4" />
                                    Envoyer un email de test
                                </button>
                            </div>
                        </div>

                    <!-- Documents -->
                    {:else if activeTab === "documents"}
                        <div class="flex flex-col gap-6">
                            {@render sectionHeader("Paramètres des documents", "Valeurs par défaut appliquées aux devis, mandats, contrats et factures.")}
                            {@render sectionDivider("Numérotation")}
                            <div class="grid grid-cols-2 gap-4">
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="inv-prefix" class="text-sm font-medium text-muted-foreground">Préfixe factures</Label.Root>
                                    <input id="inv-prefix" type="text" bind:value={invoicePrefix} placeholder="FAC" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="est-prefix" class="text-sm font-medium text-muted-foreground">Préfixe devis</Label.Root>
                                    <input id="est-prefix" type="text" bind:value={estimatePrefix} placeholder="DEV" class={inputClass} />
                                </div>
                            </div>
                            <div class="flex flex-col gap-4">
                                {@render switchRow("Numérotation automatique", "Incrémenter automatiquement le numéro à chaque nouveau document", autoNumbering, (v) => (autoNumbering = v))}
                            </div>
                            {@render sectionDivider("Facturation")}
                            <div class="grid grid-cols-2 gap-4">
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="vat-rate" class="text-sm font-medium text-muted-foreground">Taux de TVA (%)</Label.Root>
                                    <input id="vat-rate" type="number" min="0" max="100" bind:value={vatRate} placeholder="20" class={inputClass} />
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="payment-terms" class="text-sm font-medium text-muted-foreground">Délai de paiement (jours)</Label.Root>
                                    <input id="payment-terms" type="number" min="0" bind:value={defaultPaymentTerms} placeholder="30" class={inputClass} />
                                </div>
                            </div>
                            {@render sectionDivider("En-tête & pied de page")}
                            <div class="flex flex-col gap-4">
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="doc-header" class="text-sm font-medium text-muted-foreground">En-tête des documents</Label.Root>
                                    <textarea id="doc-header" bind:value={documentHeaderText} placeholder="Texte affiché en haut de chaque document (optionnel)" rows={3} class={textareaClass}></textarea>
                                </div>
                                <div class="flex flex-col gap-1.5">
                                    <Label.Root for="doc-footer" class="text-sm font-medium text-muted-foreground">Pied de page des documents</Label.Root>
                                    <textarea id="doc-footer" bind:value={documentFooterText} placeholder="Mentions légales, coordonnées bancaires, etc. (optionnel)" rows={3} class={textareaClass}></textarea>
                                </div>
                            </div>
                        </div>

                    <!-- Confidentialité -->
                    {:else if activeTab === "confidentialite"}
                        <div class="flex flex-col gap-6">
                            {@render sectionHeader("Confidentialité & sécurité", "Protégez les données sensibles de vos affaires et de vos clients.")}
                            {@render sectionDivider("Verrouillage automatique")}
                            <div class="flex flex-col gap-4">
                                {@render switchRow("Verrouillage automatique", "Verrouiller l'application après une période d'inactivité", autoLockEnabled, (v) => (autoLockEnabled = v))}
                                {#if autoLockEnabled}
                                    <div class="flex flex-col gap-1.5 max-w-xs">
                                        <Label.Root class="text-sm font-medium text-muted-foreground">Délai d'inactivité</Label.Root>
                                        {@render selectWidget(autoLockDelays, autoLockDelay, (v) => (autoLockDelay = v), selectedLockDelayLabel)}
                                    </div>
                                {/if}
                            </div>
                            {@render sectionDivider("Données")}
                            <div class="flex flex-col gap-4">
                                {@render switchRow("Chiffrement du stockage local", "Chiffrer les données en cache sur cet appareil", encryptLocalStorage, (v) => (encryptLocalStorage = v))}
                                {@render switchRow("Données d'utilisation anonymes", "Aider à améliorer l'application en envoyant des statistiques anonymes", analyticsEnabled, (v) => (analyticsEnabled = v))}
                            </div>
                            {@render sectionDivider("Export & sauvegarde")}
                            <div class="flex flex-col gap-3">
                                <p class="text-sm text-muted-foreground">Exportez l'intégralité de vos données dans un format portable pour archivage ou migration.</p>
                                <button
                                    type="button"
                                    class="h-input rounded-input border border-border-input bg-transparent text-sm font-medium text-foreground hover:bg-surface px-5 inline-flex items-center gap-2 transition-interactive cursor-pointer active:scale-[0.98] w-fit"
                                >
                                    Exporter toutes les données
                                </button>
                            </div>
                        </div>

                    <!-- Application -->
                    {:else if activeTab === "application"}
                        <div class="flex flex-col gap-6">
                            {@render sectionHeader("Application", "Informations sur la version et mises à jour logicielles.")}
                            <div class="rounded-card-sm border border-border-input bg-surface px-5 py-4 flex flex-col gap-1">
                                <span class="text-sm font-medium text-foreground">Cluo Desktop</span>
                                <span class="text-sm text-muted-foreground">Vérifiez les mises à jour pour accéder aux nouvelles fonctionnalités et corrections.</span>
                            </div>
                            <div>
                                <button
                                    type="button"
                                    onclick={() => { open = false; updateDialogOpen.set(true); }}
                                    class="h-input rounded-input border border-border-input bg-transparent text-sm font-medium text-foreground hover:bg-surface px-5 inline-flex items-center gap-2 transition-interactive cursor-pointer active:scale-[0.98]"
                                >
                                    <RefreshCw class="size-4" />
                                    Rechercher des mises à jour
                                </button>
                            </div>
                        </div>
                    {/if}

                </div>

                {#if activeTab !== "application"}
                    <div class="flex-shrink-0 border-t border-border-input px-8 py-4 flex justify-end">
                        <button
                            type="button"
                            onclick={save}
                            class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/90 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex items-center justify-center px-5 text-sm font-semibold focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer transition-interactive"
                        >
                            Enregistrer
                        </button>
                    </div>
                {/if}
            </div>

            <Dialog.Close
                class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 rounded-md focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98] cursor-pointer"
            >
                <div>
                    <X class="text-foreground-alt size-4" />
                    <span class="sr-only">Fermer</span>
                </div>
            </Dialog.Close>
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>

{#snippet sectionHeader(title: string, description: string)}
    <div>
        <h3 class="text-base font-semibold tracking-tight text-foreground">{title}</h3>
        <p class="text-sm text-muted-foreground mt-0.5">{description}</p>
    </div>
{/snippet}

{#snippet sectionDivider(label: string)}
    <div class="flex items-center gap-3">
        <span class="text-xs font-medium text-muted-foreground uppercase tracking-wider whitespace-nowrap">{label}</span>
        <div class="h-px flex-1 bg-border-input"></div>
    </div>
{/snippet}

{#snippet switchRow(label: string, description: string, value: boolean, onChange: (v: boolean) => void)}
    <div class="flex items-start justify-between gap-4">
        <div class="flex flex-col gap-0.5">
            <span class="text-sm font-medium text-foreground">{label}</span>
            <span class="text-xs text-muted-foreground">{description}</span>
        </div>
        <Switch.Root
            checked={value}
            onCheckedChange={onChange}
            class="focus-visible:ring-foreground focus-visible:ring-offset-background data-[state=checked]:bg-foreground data-[state=unchecked]:bg-dark-10 data-[state=unchecked]:shadow-mini-inset dark:data-[state=checked]:bg-foreground focus-visible:outline-hidden peer inline-flex h-[24px] min-h-[24px] w-[48px] shrink-0 cursor-pointer items-center rounded-full px-[3px] transition-colors focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
        >
            <Switch.Thumb
                class="bg-background data-[state=unchecked]:shadow-mini dark:border-background/30 dark:bg-foreground dark:shadow-popover pointer-events-none block size-[18px] shrink-0 rounded-full transition-transform data-[state=checked]:translate-x-6 data-[state=unchecked]:translate-x-0 dark:border dark:data-[state=unchecked]:border"
            />
        </Switch.Root>
    </div>
{/snippet}

{#snippet selectWidget(items: {value: string; label: string}[], value: string, onChange: (v: string) => void, displayLabel: string)}
    <Select.Root type="single" onValueChange={onChange} {items}>
        <Select.Trigger
            class="h-input rounded-input border-border-input bg-background data-placeholder:text-muted-foreground/50 inline-flex justify-between w-full select-none items-center border px-3 text-sm transition-interactive duration-150 hover:border-border-input-hover cursor-pointer"
        >
            {displayLabel}
            <ChevronsUpDown class="text-muted-foreground size-4 shrink-0" />
        </Select.Trigger>
        <Select.Portal>
            <Select.Content
                class="focus-override border-border-card bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 outline-hidden z-50 max-h-[var(--bits-select-content-available-height)] w-[var(--bits-select-anchor-width)] min-w-[var(--bits-select-anchor-width)] select-none rounded-card-sm border px-1 py-2"
                sideOffset={6}
            >
                <Select.ScrollUpButton class="flex w-full items-center justify-center py-1">
                    <ChevronsUp class="size-3" />
                </Select.ScrollUpButton>
                <Select.Viewport class="p-1">
                    {#each items as item, i (i + item.value)}
                        <Select.Item
                            class="rounded-button data-highlighted:bg-muted outline-hidden data-disabled:opacity-50 flex items-center justify-between h-9 w-full select-none py-2 pl-4 pr-2 text-sm cursor-pointer"
                            value={item.value}
                            label={item.label}
                        >
                            {#snippet children({ selected })}
                                {item.label}
                                {#if selected}
                                    <Check class="size-4 text-foreground" />
                                {/if}
                            {/snippet}
                        </Select.Item>
                    {/each}
                </Select.Viewport>
                <Select.ScrollDownButton class="flex w-full items-center justify-center py-1">
                    <ChevronsDown class="size-3" />
                </Select.ScrollDownButton>
            </Select.Content>
        </Select.Portal>
    </Select.Root>
{/snippet}

<script lang="ts">
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import {
		FileText,
		Handshake,
		ShieldCheck,
		ReceiptEuro,
		Send,
		CheckCircle,
		Plus,
		ChevronRight,
		CircleDashed,
		Loader2,
		ArrowRight,
	} from "@lucide/svelte";
	import {
		fetchDocumentWorkflow,
		sendDocument,
		acceptEstimate,
		signMandate,
		activateMandate,
		signContract,
		activateContract,
		createInvoiceFromContract,
		ConflictError,
	} from "$lib/services/api";
	import { getToastContext } from "$lib/custom/global/toast/state.svelte";
	import { TOAST_LEVELS } from "$lib/custom/global/toast/type";
	import {
		documentStatusBadge,
		documentStatusDot,
	} from "$lib/utils/badgeVariants";
	import ConfirmDialog from "$lib/custom/global/ConfirmDialog.svelte";
	import type { DocumentSummary, DocumentStatus } from "$lib/types/entities";

	const toastState = getToastContext();

	interface Props {
		caseId: string;
	}

	let { caseId }: Props = $props();

	// State
	let workflow: DocumentSummary[] = $state([]);
	let loading = $state(true);
	let error: string | null = $state(null);
	let actionInProgress: string | null = $state(null); // document id currently acting on
	let creatingInvoice = $state(false);

	// =========================================================================
	// Document type definitions for the workflow chain
	// =========================================================================

	interface DocTypeDef {
		type: string;
		label: string;
		icon: typeof FileText;
		route: string;
		createHint: string;
		autoCreated: boolean;
	}

	const DOC_TYPES: DocTypeDef[] = [
		{
			type: "estimate",
			label: "Devis",
			icon: FileText,
			route: "/cases/:id/documents/estimate",
			createHint: "Créer un devis",
			autoCreated: false,
		},
		{
			type: "mandate",
			label: "Mandat",
			icon: Handshake,
			route: "/cases/:id/documents/mandate",
			createHint: "Créé automatiquement à l'acceptation du devis",
			autoCreated: true,
		},
		{
			type: "contract",
			label: "Contrat",
			icon: ShieldCheck,
			route: "/cases/:id/documents/contract",
			createHint: "Créé automatiquement à l'activation du mandat",
			autoCreated: true,
		},
		{
			type: "invoice",
			label: "Facture",
			icon: ReceiptEuro,
			route: "/cases/:id/documents/facture",
			createHint: "Créée automatiquement à l'activation du contrat",
			autoCreated: true,
		},
	];

	// Status display
	const STATUS_LABELS: Record<string, string> = {
		draft: "Brouillon",
		sent: "Envoyé",
		signed: "Signé",
		active: "Actif",
		archived: "Archivé",
		cancelled: "Annulé",
		rejected: "Rejeté",
		expired: "Expiré",
	};

	// =========================================================================
	// Data loading
	// =========================================================================

	onMount(async () => {
		await loadWorkflow();
	});

	async function loadWorkflow() {
		if (!caseId) {
			loading = false;
			return;
		}
		loading = true;
		error = null;
		try {
			const response = await fetchDocumentWorkflow(caseId);
			workflow = response.data || [];
		} catch (e) {
			error = e instanceof Error ? e.message : "Erreur lors du chargement du workflow";
		} finally {
			loading = false;
		}
	}

	// =========================================================================
	// Helpers
	// =========================================================================

	function getDocForType(type: string): DocumentSummary | undefined {
		return workflow.find((d) => d.type === type);
	}

	function resolveRoute(route: string): string {
		return route.replace(":id", caseId);
	}

	// =========================================================================
	// Quick-action logic
	// =========================================================================

	interface QuickAction {
		label: string;
		icon: typeof Send;
		color: string;
		needsConfirm: boolean;
		confirmTitle: string;
		confirmDescription: string;
		confirmLabel: string;
		execute: () => Promise<void>;
	}

	function getQuickAction(doc: DocumentSummary): QuickAction | null {
		const { type, status, id } = doc;

		if (type === "estimate") {
			if (status === "draft") {
				return {
					label: "Envoyer",
					icon: Send,
					color: "bg-accent text-accent-foreground",
					needsConfirm: true,
					confirmTitle: "Envoyer le devis",
					confirmDescription: "Le devis sera marqué comme envoyé au client. Cette action est irréversible.",
					confirmLabel: "Envoyer",
					execute: async () => {
						await sendDocument(id, "estimate", {
							recipients: [],
							subject: `Devis`,
							message: "",
							send_email: true,
							send_sms: false,
						});
					},
				};
			}
			if (status === "sent") {
				return {
					label: "Accepter",
					icon: CheckCircle,
					color: "bg-success text-success-foreground",
					needsConfirm: true,
					confirmTitle: "Accepter le devis",
					confirmDescription: "Le devis sera marqué comme accepté et un mandat sera automatiquement créé.",
					confirmLabel: "Accepter",
					execute: async () => {
						await acceptEstimate(id);
					},
				};
			}
		}

		if (type === "mandate") {
			if (status === "draft") {
				return {
					label: "Envoyer",
					icon: Send,
					color: "bg-accent text-accent-foreground",
					needsConfirm: true,
					confirmTitle: "Envoyer le mandat",
					confirmDescription: "Le mandat sera marqué comme envoyé au client. Cette action est irréversible.",
					confirmLabel: "Envoyer",
					execute: async () => {
						await sendDocument(id, "mandate", {
							recipients: [],
							subject: `Mandat`,
							message: "",
							send_email: true,
							send_sms: false,
						});
					},
				};
			}
			if (status === "sent") {
				return {
					label: "Signer",
					icon: CheckCircle,
					color: "bg-success text-success-foreground",
					needsConfirm: true,
					confirmTitle: "Signer le mandat",
					confirmDescription: "Enregistrez la signature du mandat. Le mandat passera en état « Signé ».",
					confirmLabel: "Signer",
					execute: async () => {
						await signMandate(id, {
							signer_name: "Enquêteur",
							signer_role: "investigator",
							method: "e-sign",
						});
					},
				};
			}
			if (status === "signed") {
				return {
					label: "Activer",
					icon: ShieldCheck,
					color: "bg-success text-success-foreground",
					needsConfirm: true,
					confirmTitle: "Activer le mandat",
					confirmDescription: "Le mandat sera activé, autorisant formellement le début de l'enquête. Cette action est irréversible.",
					confirmLabel: "Activer",
					execute: async () => {
						await activateMandate(id);
					},
				};
			}
		}

		if (type === "contract") {
			if (status === "draft") {
				return {
					label: "Envoyer",
					icon: Send,
					color: "bg-accent text-accent-foreground",
					needsConfirm: true,
					confirmTitle: "Envoyer le contrat",
					confirmDescription: "Le contrat sera marqué comme envoyé au client.",
					confirmLabel: "Envoyer",
					execute: async () => {
						await sendDocument(id, "contract", {
							recipients: [],
							subject: `Contrat`,
							message: "",
							send_email: true,
							send_sms: false,
						});
					},
				};
			}
			if (status === "sent") {
				return {
					label: "Signer",
					icon: CheckCircle,
					color: "bg-success text-success-foreground",
					needsConfirm: true,
					confirmTitle: "Signer le contrat",
					confirmDescription: "Enregistrez la signature du contrat.",
					confirmLabel: "Signer",
					execute: async () => {
						await signContract(id, {
							signer_name: "Enquêteur",
							signer_role: "investigator",
							method: "e-sign",
						});
					},
				};
			}
			if (status === "signed") {
				return {
					label: "Activer",
					icon: ShieldCheck,
					color: "bg-success text-success-foreground",
					needsConfirm: true,
					confirmTitle: "Activer le contrat",
					confirmDescription: "Le contrat sera activé, mettant en vigueur l'accord commercial.",
					confirmLabel: "Activer",
					execute: async () => {
						await activateContract(id);
					},
				};
			}
		}

		if (type === "invoice") {
			if (status === "draft") {
				return {
					label: "Envoyer",
					icon: Send,
					color: "bg-accent text-accent-foreground",
					needsConfirm: true,
					confirmTitle: "Envoyer la facture",
					confirmDescription: "La facture sera marquée comme envoyée au client.",
					confirmLabel: "Envoyer",
					execute: async () => {
						await sendDocument(id, "invoice", {
							recipients: [],
							subject: `Facture`,
							message: "",
							send_email: true,
							send_sms: false,
						});
					},
				};
			}
		}

		return null;
	}

	async function executeAction(doc: DocumentSummary) {
		const action = getQuickAction(doc);
		if (!action) return;

		actionInProgress = doc.id;
		try {
			await action.execute();
			await loadWorkflow();
			toastState.add(TOAST_LEVELS.Info, "Action effectuée", "Le workflow a été mis à jour.");
		} catch (e) {
			const msg = e instanceof ConflictError
				? e.message
				: e instanceof Error
					? e.message
					: "Une erreur est survenue.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			actionInProgress = null;
		}
	}

	// Special case: create invoice from active contract
	function getCreateInvoiceAction(): QuickAction | null {
		const contract = getDocForType("contract");
		const invoice = getDocForType("invoice");
		if (!contract || contract.status !== "active" || invoice) return null;

		return {
			label: "Créer facture",
			icon: ReceiptEuro,
			color: "bg-tertiary text-background",
			needsConfirm: true,
			confirmTitle: "Créer une facture",
			confirmDescription: "Une facture sera générée à partir du contrat actif.",
			confirmLabel: "Créer",
			execute: async () => {
				const result = await createInvoiceFromContract(contract.id);
				if (result.success) {
					await goto(resolveRoute("/cases/:id/documents/facture"));
				}
			},
		};
	}

	async function executeCreateInvoiceAction() {
		const action = getCreateInvoiceAction();
		if (!action) return;
		creatingInvoice = true;
		try {
			await action.execute();
		} catch (e) {
			const msg = e instanceof ConflictError
				? e.message
				: e instanceof Error
					? e.message
					: "Une erreur est survenue.";
			toastState.add(TOAST_LEVELS.Error, "Erreur", msg);
		} finally {
			creatingInvoice = false;
		}
	}
</script>

<div class="mt-8 animate-fade-in" style="animation-delay: 500ms;">
	<div class="flex items-center justify-between mb-4">
		<h3 class="text-lg font-semibold text-foreground">Workflow des documents</h3>
		{#if !loading && !error}
			<button
				type="button"
				onclick={loadWorkflow}
				class="text-xs text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
			>
				Actualiser
			</button>
		{/if}
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-8">
			<Loader2 size={20} class="animate-spin text-muted-foreground" />
			<span class="ml-2 text-sm text-muted-foreground">Chargement du workflow...</span>
		</div>
	{:else if error}
		<div class="alert-error">
			{error}
		</div>
	{:else}
		<!-- Timeline -->
		<div class="relative">
			<!-- Connecting line -->
			<div class="absolute top-6 left-0 right-0 h-0.5 bg-border z-0 hidden lg:block"></div>

			<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 relative z-10">
				{#each DOC_TYPES as docDef, i (docDef.type)}
					{@const doc = getDocForType(docDef.type)}
					{@const action = doc ? getQuickAction(doc) : null}
					{@const createInvoiceAction = docDef.type === "invoice" ? getCreateInvoiceAction() : null}
					{@const isActing = actionInProgress === doc?.id || (docDef.type === "invoice" && creatingInvoice)}

					<div
						class="border border-border-card rounded-lg p-4 bg-background hover:shadow-md transition-all duration-300 flex flex-col min-h-[180px]"
					>
						<!-- Header: icon + status -->
						<div class="flex items-start justify-between mb-3">
							<div class="flex items-center gap-2">
								<div class="w-8 h-8 rounded-lg {doc ? 'bg-foreground/10' : 'bg-muted'} flex items-center justify-center">
									<docDef.icon
										size={16}
										class={doc ? "text-foreground" : "text-muted-foreground"}
										strokeWidth={1.5}
									/>
								</div>
								<div>
									<p class="text-sm font-semibold text-foreground">{docDef.label}</p>
									{#if doc}
										<p class="text-xs text-muted-foreground font-mono">{doc.document_ref}</p>
									{/if}
								</div>
							</div>

							<!-- Status dot/badge -->
							{#if doc}
								<span class="flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium {documentStatusBadge(doc.status as DocumentStatus)}">
									<span class="w-1.5 h-1.5 rounded-full {documentStatusDot(doc.status as DocumentStatus)}"></span>
									{STATUS_LABELS[doc.status] || doc.status}
								</span>
							{:else}
								<span class="flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium bg-muted text-muted-foreground">
									<CircleDashed size={10} />
									Absent
								</span>
							{/if}
						</div>

						<!-- Body: spacer + action -->
						<div class="flex-1"></div>

						<!-- Footer: action button or create affordance -->
						{#if doc}
							<!-- Document exists: show quick-action or navigate -->
							<div class="flex items-center gap-2 mt-2">
								{#if action}
									{#if action.needsConfirm}
										<ConfirmDialog
											title={action.confirmTitle}
											description={action.confirmDescription}
											confirmLabel={action.confirmLabel}
											onConfirm={() => executeAction(doc)}
										>
											<button
												type="button"
												disabled={isActing}
												class="h-8 rounded-md {action.color} shadow-sm hover:opacity-90 inline-flex items-center justify-center px-3 text-xs font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50 flex-1"
											>
												{#if isActing}
													<Loader2 size={12} class="mr-1 animate-spin" />
												{:else}
													<action.icon size={12} class="mr-1" />
												{/if}
												{isActing ? "..." : action.label}
											</button>
										</ConfirmDialog>
									{:else}
										<button
											type="button"
											onclick={() => executeAction(doc)}
											disabled={isActing}
											class="h-8 rounded-md {action.color} shadow-sm hover:opacity-90 inline-flex items-center justify-center px-3 text-xs font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50 flex-1"
										>
											{#if isActing}
												<Loader2 size={12} class="mr-1 animate-spin" />
											{:else}
												<action.icon size={12} class="mr-1" />
											{/if}
											{isActing ? "..." : action.label}
										</button>
									{/if}
								{/if}
								<button
									type="button"
									onclick={() => goto(resolveRoute(docDef.route))}
									class="h-8 rounded-md bg-transparent text-dark hover:bg-muted inline-flex items-center justify-center px-2 text-xs font-medium transition-colors cursor-pointer border border-border"
									title="Voir les détails"
								>
									<ChevronRight size={14} />
								</button>
							</div>
						{:else}
							<!-- Document doesn't exist: show create affordance or hint -->
							{#if docDef.type === "estimate"}
								<button
									type="button"
									onclick={() => goto(resolveRoute(docDef.route))}
									class="h-8 rounded-md bg-foreground text-background shadow-mini hover:opacity-90 inline-flex items-center justify-center px-3 text-xs font-semibold active:scale-[0.98] cursor-pointer w-full mt-2"
								>
									<Plus size={12} class="mr-1" />
									{docDef.createHint}
								</button>
							{:else if createInvoiceAction}
								<ConfirmDialog
									title={createInvoiceAction.confirmTitle}
									description={createInvoiceAction.confirmDescription}
									confirmLabel={createInvoiceAction.confirmLabel}
									onConfirm={executeCreateInvoiceAction}
								>
									<button
										type="button"
										disabled={creatingInvoice}
										class="h-8 rounded-md {createInvoiceAction.color} shadow-sm hover:opacity-90 inline-flex items-center justify-center px-3 text-xs font-semibold active:scale-[0.98] cursor-pointer disabled:opacity-50 w-full mt-2"
									>
										{#if creatingInvoice}
											<Loader2 size={12} class="mr-1 animate-spin" />
										{:else}
											<createInvoiceAction.icon size={12} class="mr-1" />
										{/if}
										{creatingInvoice ? "..." : createInvoiceAction.label}
									</button>
								</ConfirmDialog>
							{:else if docDef.autoCreated}
								<p class="text-xs text-muted-foreground italic mt-2 leading-relaxed">
									{docDef.createHint}
								</p>
							{/if}
						{/if}
					</div>

				{/each}
			</div>

			<!-- Flow arrows between columns on desktop -->
			<div class="hidden lg:flex items-center justify-between mt-2 px-[12%]">
				{#each [0, 1, 2] as i}
					<div class="flex-1 flex justify-center">
						<ArrowRight size={16} class="text-muted-foreground" />
					</div>
				{/each}
			</div>
		</div>
	{/if}
</div>

<script lang="ts">
	import { ChevronLeft, FileText, Sparkles } from "@lucide/svelte";
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";
	import TranscriptEditor from "$lib/components/TranscriptEditor.svelte";
	import { getTranscript, confirmTranscript, analyzeTranscript } from "$lib/api";
	import type { Transcript } from "$lib/types/recording";

	let transcript = $state<Transcript | null>(null);
	let isLoading = $state(true);
	let isSaving = $state(false);
	let isAnalyzing = $state(false);
	let error = $state<string | null>(null);

	// Get recording ID from URL
	let recordingId = $state("");

	onMount(() => {
		// Extract recording ID from URL path
		const pathParts = window.location.pathname.split("/");
		recordingId = pathParts[pathParts.length - 2]; // Get the ID before "transcript"

		loadTranscript();
	});

	async function loadTranscript() {
		try {
			isLoading = true;
			error = null;
			transcript = await getTranscript(recordingId);
		} catch (err) {
			error = err instanceof Error ? err.message : "Échec du chargement de la transcription";
		} finally {
			isLoading = false;
		}
	}

	async function handleSaveTranscript(editedText: string) {
		try {
			isSaving = true;
			error = null;
			await confirmTranscript(recordingId, editedText);

			// Update local state
			if (transcript) {
				transcript.text = editedText;
				transcript.isConfirmed = true;
			}
		} catch (err) {
			error = err instanceof Error ? err.message : "Échec de la sauvegarde de la transcription";
		} finally {
			isSaving = false;
		}
	}

	async function handleAnalyzeNotes() {
		try {
			isAnalyzing = true;
			error = null;
			await analyzeTranscript(recordingId);

			// Navigate to analysis page
			goto(`/recording/${recordingId}/analysis`);
		} catch (err) {
			error = err instanceof Error ? err.message : "Échec de l'analyse de la transcription";
		} finally {
			isAnalyzing = false;
		}
	}

	function goBack() {
		if (history.length > 0) history.back();
		else goto(`/recording/${recordingId}`);
	}
</script>

<div class="min-h-screen flex flex-col gap-6 pb-24 mt-8 px-4">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<button onclick={goBack} class="text-dark-700 hover:text-dark-900">
			<ChevronLeft />
		</button>
		<p class="text-dark-900 font-bold text-lg">Révision de la transcription</p>
		<div class="w-6"></div>
	</div>

	{#if isLoading}
		<div class="flex items-center justify-center p-12">
			<p class="text-dark-600">Chargement de la transcription...</p>
		</div>
	{:else if error}
		<div class="flex flex-col items-center gap-4 p-8 bg-red-50 rounded-2xl">
			<p class="text-red-700 font-semibold">Erreur</p>
			<p class="text-red-600 text-sm text-center">{error}</p>
			<button
				onclick={loadTranscript}
				class="px-6 py-3 bg-red-600 hover:bg-red-500 text-white rounded-xl transition-colors font-medium"
			>
				Réessayer
			</button>
		</div>
	{:else if transcript}
		<!-- Transcript Editor -->
		<TranscriptEditor
			text={transcript.text}
			readonly={false}
			onSave={handleSaveTranscript}
		/>

		<!-- Confirmation Status -->
		<div class="flex items-center gap-2 p-4 bg-dark-50 rounded-xl">
			<div class="w-2 h-2 rounded-full {transcript.isConfirmed ? 'bg-green-500' : 'bg-yellow-500'}"></div>
			<p class="text-dark-700 text-sm">
				{transcript.isConfirmed
					? "Transcription confirmée"
					: "Veuillez relire et confirmer la transcription"}
			</p>
		</div>

		<!-- Actions -->
		<div class="flex flex-col gap-3">
			<!-- Analyze Notes Button (only enabled after confirmation) -->
			<button
				onclick={handleAnalyzeNotes}
				disabled={!transcript.isConfirmed || isAnalyzing}
				class="flex items-center justify-center gap-2 px-6 py-4 bg-accent hover:bg-accent/90 disabled:bg-dark-200 text-accent-foreground disabled:text-dark-500 rounded-xl transition-colors font-semibold"
			>
				<Sparkles size={18} />
				<span>{isAnalyzing ? "Analyse en cours..." : "Analyser les notes"}</span>
			</button>
		</div>

		<!-- Privacy Notice -->
		<div class="flex items-center justify-center p-4 bg-dark-50 rounded-2xl mt-2">
			<p class="text-dark-600 text-sm text-center">
				La transcription et l'analyse sont traitées sur une infrastructure privée
			</p>
		</div>
	{/if}
</div>

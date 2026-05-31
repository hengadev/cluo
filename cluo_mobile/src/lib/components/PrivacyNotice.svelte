<script lang="ts">
	interface Props {
		message?: string;
		class?: string;
	}

	let {
		message = "La transcription et l'analyse sont traitées sur une infrastructure privée",
		class: className = "",
	}: Props = $props();
</script>

<div class="privacy-notice flex items-center justify-center p-4 bg-dark-50 rounded-2xl {className}">
	<p class="text-dark-600 text-sm text-center">{message}</p>
</div>

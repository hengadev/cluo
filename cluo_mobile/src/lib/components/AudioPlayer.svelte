<script lang="ts">
	import type { HTMLAudioAttributes } from "svelte/elements";
	import { Play, Pause } from "@lucide/svelte";

	interface Props {
		src: string | Blob;
		duration?: number;
		class?: HTMLAudioAttributes["class"];
	}

	let {
		src,
		duration = 0,
		class: className = "",
	}: Props = $props();

	let audioElement: HTMLAudioElement;
	let isPlaying = $state(false);
	let currentTime = $state(0);
	let progress = $state(0);

	// Format time as MM:SS
	function formatTime(seconds: number): string {
		const mins = Math.floor(seconds / 60);
		const secs = Math.floor(seconds % 60);
		return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
	}

	function togglePlay() {
		if (audioElement.paused) {
			audioElement.play();
		} else {
			audioElement.pause();
		}
	}

	function handleSeek(e: Event) {
		const target = e.target as HTMLInputElement;
		const newTime = (parseFloat(target.value) / 100) * audioElement.duration;
		audioElement.currentTime = newTime;
	}

	function handleTimeUpdate() {
		currentTime = audioElement.currentTime;
		if (audioElement.duration) {
			progress = (currentTime / audioElement.duration) * 100;
		}
	}

	function handlePlay() {
		isPlaying = true;
	}

	function handlePause() {
		isPlaying = false;
	}

	function handleEnded() {
		isPlaying = false;
		currentTime = 0;
		progress = 0;
	}

	function handleLoadedMetadata() {
		if (duration > 0) {
			// Use provided duration if available
			return;
		}
		// Otherwise use audio element duration
		duration = audioElement.duration;
	}

	// Create object URL for Blob
	const audioSrc = $derived(src instanceof Blob ? URL.createObjectURL(src) : src);
</script>

<div class="audio-player {className}">
	<audio
		bind:this={audioElement}
		src={audioSrc}
		ontimeupdate={handleTimeUpdate}
		onplay={handlePlay}
		onpause={handlePause}
		onended={handleEnded}
		onloadedmetadata={handleLoadedMetadata}
		aria-label="Lecteur audio"
	></audio>

	<div class="flex items-center gap-3">
		<button
			type="button"
			onclick={togglePlay}
			class="flex items-center justify-center w-10 h-10 bg-dark-700 hover:bg-dark-600 text-foreground rounded-full transition-colors"
			aria-label={isPlaying ? "Pause" : "Lecture"}
		>
			{#if isPlaying}
				<Pause size={16} />
			{:else}
				<Play size={16} />
			{/if}
		</button>

		<div class="flex-1 flex flex-col gap-1">
			<span class="text-dark-600 text-xs">
				{formatTime(currentTime)} / {formatTime(duration)}
			</span>
			<input
				type="range"
				min="0"
				max="100"
				value={progress}
				oninput={handleSeek}
				class="w-full h-1 bg-dark-200 rounded-lg appearance-none cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:w-3 [&::-webkit-slider-thumb]:h-3 [&::-webkit-slider-thumb]:bg-dark-700 [&::-webkit-slider-thumb]:rounded-full"
				aria-label="Rechercher dans l'audio"
			/>
		</div>
	</div>
</div>

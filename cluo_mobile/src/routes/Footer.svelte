<script lang="ts">
    import { ArrowRight, Square } from "@lucide/svelte";

    let isRecording = $state(false);
    let dragX = $state(0);
    let isDragging = $state(false);
    let containerWidth = $state(0);
    let buttonWidth = $state(60); // Approximate button width

    let mediaRecorder: MediaRecorder | null = $state(null);
    let audioChunks: Blob[] = $state([]);
    let recordingDuration = $state(0);
    let timerInterval: number | null = null;

    let containerElement: HTMLDivElement;

    // Format duration as MM:SS
    const formattedDuration = $derived(() => {
        const minutes = Math.floor(recordingDuration / 60);
        const seconds = recordingDuration % 60;
        return `${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
    });

    // Maximum distance the button can be dragged
    $effect(() => {
        if (containerElement) {
            containerWidth = containerElement.offsetWidth;
        }
    });

    const maxDrag = $derived(containerWidth - buttonWidth - 32); // 32px for padding

    function handleDragStart(e: MouseEvent | TouchEvent) {
        if (isRecording) return;
        isDragging = true;
    }

    function handleDragMove(e: MouseEvent | TouchEvent) {
        if (!isDragging || isRecording) return;

        const clientX = "touches" in e ? e.touches[0].clientX : e.clientX;
        const rect = containerElement.getBoundingClientRect();
        const newX = clientX - rect.left - buttonWidth / 2;

        // Constrain drag within bounds
        dragX = Math.max(0, Math.min(newX, maxDrag));
    }

    function handleDragEnd() {
        if (!isDragging || isRecording) return;
        isDragging = false;

        // If dragged more than 80% of the way, start recording
        if (dragX > maxDrag * 0.8) {
            startRecording();
        } else {
            // Reset position
            dragX = 0;
        }
    }

    async function startRecording() {
        try {
            dragX = 0;

            const stream = await navigator.mediaDevices.getUserMedia({
                audio: true,
            });

            mediaRecorder = new MediaRecorder(stream);
            audioChunks = [];

            mediaRecorder.ondataavailable = (e) => {
                audioChunks.push(e.data);
            };

            mediaRecorder.start();
            isRecording = true;

            // Start timer
            recordingDuration = 0;
            timerInterval = window.setInterval(() => {
                recordingDuration++;
            }, 1000);
        } catch (error) {
            console.error("Failed to start recording:", error);
            // Reset state if recording failed
            dragX = 0;
            isRecording = false;
            if (timerInterval !== null) {
                clearInterval(timerInterval);
                timerInterval = null;
            }
        }
    }

    async function stopRecording() {
        if (!mediaRecorder) return;

        // Set onstop handler before calling stop() to avoid race condition
        mediaRecorder.onstop = async () => {
            const audioBlob = new Blob(audioChunks, {
                type: mediaRecorder!.mimeType,
            });

            // Stop all tracks to release the microphone
            mediaRecorder!.stream.getTracks().forEach((track) => track.stop());

            await sendAudio(audioBlob);

            // Clear chunks after sending
            audioChunks = [];
        };

        mediaRecorder.stop();
        isRecording = false;

        // Stop timer
        if (timerInterval !== null) {
            clearInterval(timerInterval);
            timerInterval = null;
        }
    }

    async function sendAudio(blob: Blob) {
        try {
            const formData = new FormData();
            formData.append("audio", blob, "recording.webm");

            const response = await fetch("/api/audio", {
                method: "POST",
                body: formData,
            });

            if (!response.ok) {
                throw new Error(`Failed to send audio: ${response.statusText}`);
            }
        } catch (error) {
            console.error("Failed to send audio:", error);
            // TODO: Show feedback about failed upload
        }
    }
</script>

<svelte:window
    onmousemove={handleDragMove}
    onmouseup={handleDragEnd}
    ontouchmove={handleDragMove}
    ontouchend={handleDragEnd}
/>

<div
    bind:this={containerElement}
    class="relative flex justify-center items-center bg-dark-900 px-4 py-6 min-h-20 overflow-hidden"
>
    {#if !isRecording}
        <div class="absolute inset-0 flex items-center justify-center">
            <p class="text-dark-200 text-base select-none">Slide to start</p>
        </div>

        <button
            class="absolute left-4 flex bg-dark-700 p-3 rounded-2xl cursor-grab active:cursor-grabbing transition-colors touch-none z-10"
            style="transform: translateX({dragX}px); transition: {isDragging
                ? 'none'
                : 'transform 0.3s ease-out'}"
            onmousedown={handleDragStart}
            ontouchstart={handleDragStart}
        >
            <ArrowRight class="text-foreground" />
        </button>
    {:else}
        <div class="flex justify-between items-center w-full">
            <div class="w-12"></div>

            <div class="flex items-center gap-2">
                <div
                    class="w-3 h-3 bg-destructive rounded-full animate-pulse"
                ></div>
                <p class="text-dark-200 text-lg font-mono font-semibold">
                    {formattedDuration()}
                </p>
            </div>

            <button
                class="flex bg-destructive p-3 rounded-2xl hover:bg-destructive/90 transition-colors"
                onclick={stopRecording}
            >
                <Square class="text-white" fill="white" size={20} />
            </button>
        </div>
    {/if}
</div>

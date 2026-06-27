<script lang="ts">
    import { onMount } from "svelte";
    import { ArrowRight, Square, Check, X } from "@lucide/svelte";
    import { goto } from "$app/navigation";
    import AudioPlayer from "$lib/components/AudioPlayer.svelte";
    import { uploadRecording } from "$lib/api";
    import { enqueue } from "$lib/upload-queue";
    import { snackbar } from "$lib/stores/snackbar";
    import { queueCount } from "$lib/stores/upload-queue-count";
    import type { Case } from "$lib/types/case";

    interface Props {
        currentCase?: Case | null;
    }

    let { currentCase = null }: Props = $props();

    type FooterState = "idle" | "recording" | "preview";
    type ConfirmDialog = { show: true; onConfirm: () => void } | { show: false };

    let footerState: FooterState = $state("idle");
    let isRecording = $state(false);
    let dragX = $state(0);
    let isDragging = $state(false);
    let containerWidth = $state(0);
    let buttonWidth = $state(60); // Approximate button width

    let mediaRecorder: MediaRecorder | null = $state(null);
    let audioChunks: Blob[] = $state([]);
    let recordingDuration = $state(0);
    let timerInterval: number | null = null;
    let recordedBlob: Blob | null = $state(null);
    let recordingTitle: string = $state("");
    let defaultRecordingTitle: string = $state("");

    let containerElement: HTMLDivElement;

    onMount(async () => {
        try {
            const status = await navigator.permissions.query({ name: "microphone" as PermissionName });
            if (status.state === "prompt") {
                const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
                stream.getTracks().forEach((t) => t.stop());
            }
        } catch {
            // Permission API or getUserMedia not supported — will fall back to requesting at record time
        }
    });

    const formattedDuration = $derived.by(() => {
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
        if (footerState !== "idle") return;
        if (!currentCase) return;
        isDragging = true;
    }

    function handleDragMove(e: MouseEvent | TouchEvent) {
        if (!isDragging || footerState !== "idle") return;

        const clientX = "touches" in e ? e.touches[0].clientX : e.clientX;
        const rect = containerElement.getBoundingClientRect();
        const newX = clientX - rect.left - buttonWidth / 2;

        // Constrain drag within bounds
        dragX = Math.max(0, Math.min(newX, maxDrag));
    }

    function handleDragEnd() {
        if (!isDragging || footerState !== "idle") return;
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
            footerState = "recording";

            // Get supported MIME type
            const mimeType = MediaRecorder.isTypeSupported("audio/webm;codecs=opus")
                ? "audio/webm;codecs=opus"
                : MediaRecorder.isTypeSupported("audio/mp4")
                    ? "audio/mp4"
                    : "";

            if (!mimeType) {
                throw new Error("No supported audio format found");
            }

            const stream = await navigator.mediaDevices.getUserMedia({
                audio: {
                    sampleRate: 48000,
                    channelCount: 1,
                    echoCancellation: true,
                    noiseSuppression: true,
                    autoGainControl: true,
                },
            });

            mediaRecorder = new MediaRecorder(stream, { mimeType });
            audioChunks = [];

            mediaRecorder.ondataavailable = (e) => {
                if (e.data.size > 0) {
                    audioChunks.push(e.data);
                }
            };

            mediaRecorder.start(1000); // Collect data every second
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
            footerState = "idle";
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

            // Store blob for preview
            recordedBlob = audioBlob;
            defaultRecordingTitle = generateTimestampTitle();
            recordingTitle = defaultRecordingTitle;
            footerState = "preview";

            // Clear chunks
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

    function generateTimestampTitle(date: Date = new Date()): string {
        const hours = date.getHours().toString().padStart(2, "0");
        const minutes = date.getMinutes().toString().padStart(2, "0");
        return `Enregistrement ${hours}h${minutes}`;
    }

    function discardRecording() {
        recordedBlob = null;
        recordingTitle = "";
        defaultRecordingTitle = "";
        footerState = "idle";
        dragX = 0;
    }

    let confirmDialog: ConfirmDialog = $state({ show: false });
    let isUploading = $state(false);
    let lastUploadBlob: Blob | null = $state(null);

    async function keepRecording() {
        if (!recordedBlob) return;

        if (currentCase?.status === "released") {
            confirmDialog = {
                show: true,
                onConfirm: () => {
                    confirmDialog = { show: false };
                    sendAudio(recordedBlob!);
                },
            };
            return;
        }

        sendAudio(recordedBlob);
    }

    function effectiveTitle(): string {
        return recordingTitle.trim() || defaultRecordingTitle;
    }

    async function sendAudio(blob: Blob) {
        if (isUploading) return;

        const title = effectiveTitle();

        try {
            isUploading = true;
            lastUploadBlob = blob;

            const response = await uploadRecording(blob, { caseId: currentCase!.id, title });
            const recordingId = response.id;

            // Reset state and navigate to processing page
            recordedBlob = null;
            recordingTitle = "";
            defaultRecordingTitle = "";
            footerState = "idle";
            goto(`/processing/${recordingId}`);
        } catch (error) {
            console.error("Failed to send audio:", error);

            try {
                await enqueue(blob, { caseId: currentCase!.id, title });
                snackbar.show("Enregistrement mis en attente — sera envoyé dès que possible");
                queueCount.refresh();
            } catch {
                snackbar.error("Échec de l'envoi de l'enregistrement");
            }

            // Reset to idle so the investigator can continue recording
            recordedBlob = null;
            recordingTitle = "";
            defaultRecordingTitle = "";
            footerState = "idle";
            dragX = 0;
        } finally {
            isUploading = false;
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
    {#if footerState === "idle"}
        {#if currentCase}
            <div class="absolute inset-0 flex items-center justify-center">
                <p class="text-dark-200 text-base select-none">Glisser pour commencer</p>
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
            <p class="text-dark-400 text-base select-none text-center">
                Sélectionnez une affaire pour enregistrer
            </p>
        {/if}
    {:else if footerState === "recording"}
        <div class="flex justify-between items-center w-full">
            <div class="w-12"></div>

            <div class="flex items-center gap-2">
                <div
                    class="w-3 h-3 bg-destructive rounded-full animate-pulse"
                ></div>
                <p class="text-dark-200 text-lg font-mono font-semibold">
                    {formattedDuration}
                </p>
            </div>

            <button
                class="flex bg-destructive p-3 rounded-2xl hover:bg-destructive/90 transition-colors"
                onclick={stopRecording}
            >
                <Square class="text-white" fill="white" size={20} />
            </button>
        </div>
    {:else if footerState === "preview" && recordedBlob}
        <div class="flex flex-col gap-3 w-full">
            <AudioPlayer src={recordedBlob} duration={recordingDuration} />
            <input
                type="text"
                bind:value={recordingTitle}
                placeholder={defaultRecordingTitle}
                class="w-full bg-background border border-dark-200 text-dark-800 placeholder-dark-400 px-3 py-2 rounded-xl text-sm focus:outline-none focus:ring-1 focus:ring-dark-400"
            />
            <div class="flex gap-3">
                <button
                    class="flex-1 flex items-center justify-center gap-2 bg-dark-100 hover:bg-dark-200 text-dark-700 px-4 py-3 rounded-xl transition-colors"
                    onclick={discardRecording}
                >
                    <X size={18} />
                    <span class="text-sm font-medium">Annuler</span>
                </button>
                <button
                    class="flex-1 flex items-center justify-center gap-2 bg-green-600 hover:bg-green-500 text-white px-4 py-3 rounded-xl transition-colors disabled:opacity-50"
                    onclick={keepRecording}
                    disabled={isUploading}
                >
                    <Check size={18} />
                    <span class="text-sm font-medium">Conserver et envoyer</span>
                </button>
            </div>
        </div>
    {/if}
</div>

{#if confirmDialog.show}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60"
        role="dialog"
        aria-modal="true"
        aria-label="Confirmation"
    >
        <div class="bg-dark-800 rounded-2xl p-6 mx-4 max-w-sm w-full shadow-popover">
            <p class="text-foreground text-base mb-6">Cette affaire est clôturée. Continuer quand même ?</p>
            <div class="flex gap-3">
                <button
                    class="flex-1 bg-dark-100 hover:bg-dark-200 text-dark-700 px-4 py-3 rounded-xl transition-colors"
                    onclick={() => (confirmDialog = { show: false })}
                >
                    Annuler
                </button>
                <button
                    class="flex-1 bg-green-600 hover:bg-green-500 text-white px-4 py-3 rounded-xl transition-colors"
                    onclick={confirmDialog.onConfirm}
                >
                    Continuer
                </button>
            </div>
        </div>
    </div>
{/if}

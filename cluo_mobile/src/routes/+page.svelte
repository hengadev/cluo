<script lang="ts">
    let mediaRecorder: MediaRecorder | null = null;
    let audioChunks: Blob[] = [];
    let isRecording = false;

    async function startRecording() {
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
    }

    async function stopRecording() {
        if (!mediaRecorder) return;

        mediaRecorder.stop();
        isRecording = false;

        mediaRecorder.onstop = async () => {
            const audioBlob = new Blob(audioChunks, {
                type: mediaRecorder!.mimeType,
            });

            await sendAudio(audioBlob);
        };
    }

    async function sendAudio(blob: Blob) {
        const formData = new FormData();
        formData.append("audio", blob, "recording.webm");

        await fetch("/api/audio", {
            method: "POST",
            body: formData,
        });
    }
    const btn = "px-8 py-4 border-1 border-black rounded-input cursor-pointer";
</script>

<div class="min-h-screen flex justify-center items-center">
    <div>
        <button class={btn} on:click={startRecording} disabled={isRecording}>
            Start recording
        </button>

        <button class={btn} on:click={stopRecording} disabled={!isRecording}>
            Stop recording
        </button>
    </div>
</div>

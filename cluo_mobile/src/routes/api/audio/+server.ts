import type { RequestHandler } from "./$types";
import { uploadAudio } from "$lib/api/server";

export const POST: RequestHandler = async ({ request, fetch }) => {
    // Parse incoming multipart form
    const formData = await request.formData();

    const file = formData.get("audio");
    const caseId = formData.get("caseId");

    // Validate required fields
    if (!file || !(file instanceof Blob)) {
        return new Response(
            JSON.stringify({ error: "Missing or invalid 'audio' file field" }),
            { status: 400, headers: { "Content-Type": "application/json" } },
        );
    }

    if (!caseId || typeof caseId !== "string") {
        return new Response(
            JSON.stringify({ error: "Missing 'caseId' field" }),
            { status: 400, headers: { "Content-Type": "application/json" } },
        );
    }

    try {
        const mediaId = await uploadAudio(fetch, file, caseId);
        return new Response(
            JSON.stringify({ id: mediaId }),
            { status: 201, headers: { "Content-Type": "application/json" } },
        );
    } catch (err) {
        const message = err instanceof Error ? err.message : "Upload failed";
        const status = message.includes("413") ? 413
            : message.includes("400") ? 400
            : 502;
        return new Response(
            JSON.stringify({ error: message }),
            { status, headers: { "Content-Type": "application/json" } },
        );
    }
};

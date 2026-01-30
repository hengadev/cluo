import type { RequestHandler } from "./$types";
import { env } from "$env/dynamic/private";

const API_URL = env.VITE_API_URL ?? "";

export const POST: RequestHandler = async ({ request, fetch }) => {
    const formData = await request.formData();

    const res = await fetch(`${API_URL}/api/audio`, {
        method: "POST",
        body: formData,
        headers: {
            // DO NOT set Content-Type manually
            // fetch will set multipart boundaries
        },
    });

    if (!res.ok) {
        const errorText = await res.text().catch(() => "Upload failed");
        return new Response(JSON.stringify({ error: errorText }), {
            status: res.status,
            headers: { "Content-Type": "application/json" },
        });
    }

    return new Response(await res.text(), {
        status: res.status,
        headers: { "Content-Type": "application/json" },
    });
};


import type { RequestHandler } from "./$types";

export const POST: RequestHandler = async ({ request, fetch }) => {
    const formData = await request.formData();

    // TODO: I need to change that API address to make things work properly
    const res = await fetch("http://golang-api:8080/audio", {
        method: "POST",
        body: formData,
        headers: {
            // DO NOT set Content-Type manually
            // fetch will set multipart boundaries
        }
    });

    return new Response(await res.arrayBuffer(), {
        status: res.status,
        headers: res.headers
    });
};


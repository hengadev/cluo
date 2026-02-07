import type { Actions } from "./$types"
import { fail } from '@sveltejs/kit';
import { env } from "$env/dynamic/private"

// here I just use the action brother
export const action = {
    default: async ({ request }) => {
        // TODO: handle authentication
        const data = await request.formData()
        const id = data.get("id")
        const password = data.get("password")
        // do the fetch to the backend here
        const api_url = env.API_URL
        if (!api_url) {
        }
        try {
            const res = await fetch(`http://${api_url}/auth`, {
                body: JSON.stringify({ id, password }),
            })
            if (!res.ok) {
                return fail(400, { id, incorrect: true })
            }
        } catch (err) {
            console.error(err)
        }
        return { success: true }
    }
} satisfies Actions;

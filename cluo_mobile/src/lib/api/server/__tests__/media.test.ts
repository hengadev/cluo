import { describe, it, expect, vi, beforeEach } from "vitest";

// Mock the SvelteKit error function before importing the module under test
vi.mock("@sveltejs/kit", () => ({
	error: vi.fn((status: number, body: { message: string }) => {
		const err = new Error(body.message);
		(err as any).status = status;
		throw err;
	}),
}));

// Mock import.meta.env
vi.stubEnv("VITE_API_URL", "https://api.test.local");

// Import after mocks are in place
const { fetchMediaById, fetchProcessingStatus, uploadAudio } = await import("../media");

function mockFetch(response: { status: number; json?: any; text?: string }): typeof fetch {
	return vi.fn(async () => ({
		status: response.status,
		ok: response.status >= 200 && response.status < 300,
		json: async () => response.json,
		text: async () => response.text ?? JSON.stringify(response.json),
		headers: new Headers({ "Content-Type": "application/json" }),
	})) as any;
}

const SAMPLE_MEDIA = {
	id: "abc-123",
	caseId: "case-456",
	url: "https://cdn.test.local/audio.webm",
	type: "audio",
	mimeType: "audio/webm",
	fileName: "recording.webm",
	fileSize: 1024,
	caption: "Witness interview",
	isPublished: false,
	createdAt: "2025-01-15T14:30:00Z",
};

// ─── fetchMediaById ───

describe("fetchMediaById", () => {
	it("calls GET /media/{id} and returns mapped recording data", async () => {
		const fetch = mockFetch({ status: 200, json: SAMPLE_MEDIA });

		const result = await fetchMediaById(fetch, "abc-123");

		expect(fetch).toHaveBeenCalledWith(
			"https://api.test.local/media/abc-123",
			expect.objectContaining({
				headers: { Accept: "application/json" },
			}),
		);

		// Verify shape and key fields
		expect(result.id).toBe("abc-123");
		expect(result.title).toBe("Witness interview");
		expect(result.audioUrl).toBe("https://cdn.test.local/audio.webm");
		expect(result.duration).toBe("00:00");
		// SAMPLE_MEDIA has a url but isPublished=false → transcribing
		expect(result.status).toBe("transcribing");

		// Date and time are locale-formatted, just verify they are non-empty strings
		expect(result.date).toBeTruthy();
		expect(result.startTime).toBeTruthy();
		expect(typeof result.date).toBe("string");
		expect(typeof result.startTime).toBe("string");
	});

	it("uses fileName as title when caption is empty", async () => {
		const media = { ...SAMPLE_MEDIA, caption: "" };
		const fetch = mockFetch({ status: 200, json: media });

		const result = await fetchMediaById(fetch, "abc-123");

		expect(result.title).toBe("recording.webm");
	});

	it("throws a 404 SvelteKit error when API returns 404", async () => {
		const fetch = mockFetch({ status: 404, text: "Not found" });

		await expect(fetchMediaById(fetch, "nonexistent")).rejects.toThrow("Enregistrement introuvable");
		const { error } = await import("@sveltejs/kit");
		expect(error).toHaveBeenCalledWith(404, { message: "Enregistrement introuvable" });
	});

	it("throws a 500 SvelteKit error when API returns 500", async () => {
		const fetch = mockFetch({ status: 500, text: "Internal Server Error" });

		await expect(fetchMediaById(fetch, "abc-123")).rejects.toThrow("Échec du chargement de l'enregistrement");
		const { error } = await import("@sveltejs/kit");
		expect(error).toHaveBeenCalledWith(500, { message: "Échec du chargement de l'enregistrement" });
	});
});

// ─── fetchProcessingStatus ───

describe("fetchProcessingStatus", () => {
	it("calls GET /media/{id} and returns processing steps", async () => {
		const fetch = mockFetch({ status: 200, json: SAMPLE_MEDIA });

		const result = await fetchProcessingStatus(fetch, "abc-123");

		expect(fetch).toHaveBeenCalledWith(
			"https://api.test.local/media/abc-123",
			expect.objectContaining({
				headers: { Accept: "application/json" },
			}),
		);

		expect(result.recordingId).toBe("abc-123");
		expect(result.error).toBeNull();
		expect(result.steps).toHaveLength(4);

		// Step 1 always completed
		expect(result.steps[0]).toEqual({
			title: "Téléchargement audio",
			status: "completed",
		});

		// Steps 2-4 are processing when not published
		expect(result.steps[1].status).toBe("processing");
		expect(result.steps[2].status).toBe("processing");
		expect(result.steps[3].status).toBe("processing");
	});

	it("marks steps 2-4 as completed when media is published", async () => {
		const media = { ...SAMPLE_MEDIA, isPublished: true };
		const fetch = mockFetch({ status: 200, json: media });

		const result = await fetchProcessingStatus(fetch, "abc-123");

		expect(result.steps[1].status).toBe("completed");
		expect(result.steps[2].status).toBe("completed");
		expect(result.steps[3].status).toBe("completed");
	});

	it("throws 404 when API returns 404", async () => {
		const fetch = mockFetch({ status: 404, text: "Not found" });

		await expect(fetchProcessingStatus(fetch, "nonexistent")).rejects.toThrow(
			"Enregistrement introuvable",
		);
	});

	it("throws 500 when API returns 500", async () => {
		const fetch = mockFetch({ status: 500, text: "Error" });

		await expect(fetchProcessingStatus(fetch, "abc-123")).rejects.toThrow(
			"Échec du chargement de l'état du traitement",
		);
	});
});

// ─── uploadAudio ───

describe("uploadAudio", () => {
	it("calls POST /media with correct form data and returns media ID", async () => {
		const fetch = vi.fn(async () => ({
			status: 201,
			ok: true,
			json: async () => SAMPLE_MEDIA,
			text: async () => JSON.stringify(SAMPLE_MEDIA),
			headers: new Headers({ "Content-Type": "application/json" }),
		})) as any;

		const blob = new Blob(["audio-data"], { type: "audio/webm" });
		const result = await uploadAudio(fetch, blob, "case-456");

		expect(result).toBe("abc-123");
		expect(fetch).toHaveBeenCalledWith(
			"https://api.test.local/media",
			expect.objectContaining({ method: "POST" }),
		);

		// Verify FormData was sent
		const callArgs = fetch.mock.calls[0][1];
		expect(callArgs.body).toBeInstanceOf(FormData);
	});

	it("throws on API failure", async () => {
		const fetch = vi.fn(async () => ({
			status: 500,
			ok: false,
			json: async () => ({ error: "fail" }),
			text: async () => "Internal Server Error",
			headers: new Headers(),
		})) as any;

		const blob = new Blob(["audio-data"], { type: "audio/webm" });
		await expect(uploadAudio(fetch, blob, "case-456")).rejects.toThrow(
			"Échec de l'envoi (500)",
		);
	});
});

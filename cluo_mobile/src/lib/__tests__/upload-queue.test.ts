/**
 * Unit tests for the offline upload queue module.
 *
 * Uses a fake IndexedDB (fake-indexeddb) so tests run without a real browser.
 */
import { describe, it, expect, beforeEach, vi } from "vitest";
import "fake-indexeddb/auto";
import { enqueue, list, remove, flush, _reset, type QueuedUpload } from "../upload-queue";

beforeEach(() => {
	_reset();
	// Close and delete the database between tests for a clean slate
	return new Promise<void>((resolve) => {
		const req = indexedDB.deleteDatabase("cluo_recordings");
		req.onsuccess = () => resolve();
		req.onerror = () => resolve();
		req.onblocked = () => resolve();
	});
});

function makeBlob(content = "audio-data"): Blob {
	return new Blob([content], { type: "audio/webm" });
}

describe("upload queue", () => {
	describe("enqueue", () => {
		it("persists an entry and returns a queue ID", async () => {
			const id = await enqueue(makeBlob(), { caseId: "case-1", title: "Test" });
			expect(id).toBeTruthy();

			const entries = await list();
			expect(entries).toHaveLength(1);
			expect(entries[0].id).toBe(id);
		});

		it("stores the blob and metadata", async () => {
			const blob = makeBlob("hello");
			await enqueue(blob, { caseId: "c1", title: "Titre" });

			const entries = await list();
			expect(entries[0].metadata.caseId).toBe("c1");
			expect(entries[0].metadata.title).toBe("Titre");
			expect(entries[0].blob).toBeInstanceOf(Blob);
		});

		it("initialises attemptCount to 0", async () => {
			await enqueue(makeBlob(), { caseId: "c1", title: "T" });
			const entries = await list();
			expect(entries[0].attemptCount).toBe(0);
		});

		it("sets enqueuedAt to a valid ISO date", async () => {
			await enqueue(makeBlob(), { caseId: "c1", title: "T" });
			const entries = await list();
			expect(new Date(entries[0].enqueuedAt).getTime()).not.toBeNaN();
		});
	});

	describe("list", () => {
		it("returns an empty array when nothing is queued", async () => {
			const entries = await list();
			expect(entries).toEqual([]);
		});

		it("returns all queued entries", async () => {
			await enqueue(makeBlob("a"), { caseId: "c1", title: "A" });
			await enqueue(makeBlob("b"), { caseId: "c2", title: "B" });
			const entries = await list();
			expect(entries).toHaveLength(2);
		});
	});

	describe("remove", () => {
		it("clears the specified entry", async () => {
			const id = await enqueue(makeBlob(), { caseId: "c1", title: "T" });
			await remove(id);
			const entries = await list();
			expect(entries).toHaveLength(0);
		});

		it("only removes the targeted entry", async () => {
			const id1 = await enqueue(makeBlob("a"), { caseId: "c1", title: "A" });
			const id2 = await enqueue(makeBlob("b"), { caseId: "c2", title: "B" });
			await remove(id1);
			const entries = await list();
			expect(entries).toHaveLength(1);
			expect(entries[0].id).toBe(id2);
		});

		it("does not throw for a non-existent ID", async () => {
			await expect(remove("does-not-exist")).resolves.toBeUndefined();
		});
	});

	describe("flush", () => {
		it("calls uploadFn for each entry", async () => {
			await enqueue(makeBlob("a"), { caseId: "c1", title: "A" });
			await enqueue(makeBlob("b"), { caseId: "c2", title: "B" });

			const uploadFn = vi.fn().mockResolvedValue({ id: "mock" });
			const result = await flush(uploadFn);

			expect(uploadFn).toHaveBeenCalledTimes(2);
			expect(result.succeeded).toHaveLength(2);
			expect(result.failed).toHaveLength(0);
		});

		it("removes succeeded entries from the queue", async () => {
			await enqueue(makeBlob(), { caseId: "c1", title: "T" });
			const uploadFn = vi.fn().mockResolvedValue({ id: "mock" });
			await flush(uploadFn);

			const entries = await list();
			expect(entries).toHaveLength(0);
		});

		it("retains failed entries in the queue", async () => {
			await enqueue(makeBlob(), { caseId: "c1", title: "T" });
			const uploadFn = vi.fn().mockRejectedValue(new Error("network"));
			const result = await flush(uploadFn);

			expect(result.succeeded).toHaveLength(0);
			expect(result.failed).toHaveLength(1);

			const entries = await list();
			expect(entries).toHaveLength(1);
		});

		it("increments attemptCount on failure", async () => {
			await enqueue(makeBlob(), { caseId: "c1", title: "T" });
			const uploadFn = vi.fn().mockRejectedValue(new Error("fail"));
			await flush(uploadFn);

			const entries = await list();
			expect(entries[0].attemptCount).toBe(1);
		});

		it("increments attemptCount on each failed flush", async () => {
			await enqueue(makeBlob(), { caseId: "c1", title: "T" });
			const failingFn = vi.fn().mockRejectedValue(new Error("fail"));

			await flush(failingFn);
			await flush(failingFn);
			await flush(failingFn);

			const entries = await list();
			expect(entries[0].attemptCount).toBe(3);
		});

		it("is a no-op when the queue is empty", async () => {
			const uploadFn = vi.fn().mockResolvedValue({ id: "mock" });
			const result = await flush(uploadFn);

			expect(uploadFn).not.toHaveBeenCalled();
			expect(result.succeeded).toHaveLength(0);
			expect(result.failed).toHaveLength(0);
		});

		it("handles a mix of successes and failures", async () => {
			const id1 = await enqueue(makeBlob("a"), { caseId: "c1", title: "A" });
			const id2 = await enqueue(makeBlob("b"), { caseId: "c2", title: "B" });

			let callCount = 0;
			const uploadFn = vi.fn().mockImplementation(() => {
				callCount++;
				if (callCount === 1) return Promise.reject(new Error("fail"));
				return Promise.resolve({ id: "mock" });
			});

			const result = await flush(uploadFn);
			expect(result.succeeded).toHaveLength(1);
			expect(result.failed).toHaveLength(1);

			// The failed entry (first enqueued) should still be there
			const remaining = await list();
			expect(remaining).toHaveLength(1);
			expect(remaining[0].attemptCount).toBe(1);
		});

		it("is idempotent when called concurrently", async () => {
			await enqueue(makeBlob(), { caseId: "c1", title: "T" });

			const uploadFn = vi.fn().mockImplementation(
				() => new Promise((resolve) => setTimeout(() => resolve({ id: "mock" }), 50)),
			);

			// Fire two flushes in parallel
			const [r1, r2] = await Promise.all([flush(uploadFn), flush(uploadFn)]);

			// The second flush should have returned empty because the first was in progress
			const totalSucceeded = r1.succeeded.length + r2.succeeded.length;
			expect(totalSucceeded).toBe(1);

			// Only one uploadFn call should have happened
			expect(uploadFn).toHaveBeenCalledTimes(1);
		});
	});
});

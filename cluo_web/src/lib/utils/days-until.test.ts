import { describe, it, expect, vi, beforeEach } from 'vitest';
import { daysUntil } from './days-until';

describe('daysUntil', () => {
	beforeEach(() => {
		// Lock "today" to June 2, 2026 at noon (local time doesn't matter — we floor to midnight)
		vi.setSystemTime(new Date('2026-06-02T12:00:00'));
	});

	it('returns a positive integer for a future date', () => {
		expect(daysUntil('2026-06-10T00:00:00')).toBe(8);
	});

	it('returns 7 at exactly the 7-day boundary', () => {
		expect(daysUntil('2026-06-09T00:00:00')).toBe(7);
	});

	it('returns 0 on the same day', () => {
		expect(daysUntil('2026-06-02T23:59:59')).toBe(0);
	});

	it('returns a negative number for a past date', () => {
		expect(daysUntil('2026-05-25T00:00:00')).toBe(-8);
	});

	it('handles month boundaries correctly', () => {
		// June 2 → July 2 = 30 days
		expect(daysUntil('2026-07-02T00:00:00')).toBe(30);
	});

	it('handles year boundaries correctly', () => {
		// June 2 → Jan 5 next year
		expect(daysUntil('2027-01-05T00:00:00')).toBe(217);
	});
});

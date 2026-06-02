/**
 * Returns the number of whole days between today (midnight local time)
 * and the given ISO date string. Negative when the date is in the past.
 */
export function daysUntil(iso: string): number {
	const today = new Date();
	today.setHours(0, 0, 0, 0);

	const target = new Date(iso);
	target.setHours(0, 0, 0, 0);

	const diffMs = target.getTime() - today.getTime();
	return Math.floor(diffMs / (1000 * 60 * 60 * 24));
}

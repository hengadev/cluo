<script lang="ts">
	function reload() {
		window.location.reload();
	}
</script>

<svelte:head>
	<title>You're offline — Cluo</title>
</svelte:head>

<main>
	<div class="card">
		<div class="icon" aria-hidden="true">
			<svg
				width="48"
				height="48"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="1.75"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<line x1="1" y1="1" x2="23" y2="23" />
				<path d="M16.72 11.06A10.94 10.94 0 0 1 19 12.55" />
				<path d="M5 12.55a10.94 10.94 0 0 1 5.17-2.39" />
				<path d="M10.71 5.05A16 16 0 0 1 22.58 9" />
				<path d="M1.42 9a15.91 15.91 0 0 1 4.7-2.88" />
				<path d="M8.53 16.11a6 6 0 0 1 6.95 0" />
				<line x1="12" y1="20" x2="12.01" y2="20" />
			</svg>
		</div>
		<h1>You're offline</h1>
		<p>Cluo can't reach the network right now. Check your connection and try again.</p>
		<button onclick={reload}>Retry connection</button>
	</div>
</main>

<style>
	main {
		min-height: 100vh;
		min-height: 100dvh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1.5rem;
		/* Respect the notch/status bar on iOS standalone mode */
		padding-top: calc(env(safe-area-inset-top) + 1.5rem);
		font-family:
			system-ui,
			-apple-system,
			'Segoe UI',
			Roboto,
			sans-serif;
		background: #ffffff;
		color: #1a1a1a;
		box-sizing: border-box;
	}
	.card {
		max-width: 24rem;
		text-align: center;
	}
	.icon {
		display: flex;
		justify-content: center;
		color: #1a1a1a;
		opacity: 0.85;
		margin-bottom: 1.25rem;
	}
	h1 {
		font-size: 1.375rem;
		font-weight: 600;
		margin: 0 0 0.5rem;
	}
	p {
		margin: 0 0 1.75rem;
		color: #555;
		line-height: 1.5;
	}
	button {
		appearance: none;
		border: 1px solid #1a1a1a;
		background: #1a1a1a;
		color: #ffffff;
		font: inherit;
		font-weight: 500;
		padding: 0.7rem 1.5rem;
		border-radius: 9999px;
		cursor: pointer;
		transition: opacity 0.15s ease;
	}
	button:hover {
		opacity: 0.85;
	}
</style>

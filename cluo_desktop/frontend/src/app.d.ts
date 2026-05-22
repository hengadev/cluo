// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

interface ImportMetaEnv {
	readonly VITE_USE_MOCK_DATA: string;
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}

export {};

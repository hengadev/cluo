# CLAUDE.md

This file provides guidance to **Claude Code (claude.ai/code)** when working with this SvelteKit-based PWA repository.

---

## 🧭 Development Commands

### Core Commands (using pnpm)

- `pnpm dev` — Start SvelteKit development server (Vite)  
- `pnpm build` — Build production-ready app  
- `pnpm preview` — Preview production build locally  
- `pnpm prepare` — Sync SvelteKit project files  
- `pnpm check` — Run TypeScript + Svelte type checks  
- `pnpm check:watch` — Watch mode for type checks  

> Development server runs at **http://localhost:5173** by default.

---

## ⚙️ Project Architecture

### Core Stack

- **Framework**: [SvelteKit](https://kit.svelte.dev/) with TypeScript  
- **UI Library**: [Bits-UI](https://bits-ui.com/) (Radix-based primitives)  
- **Styling**: [Tailwind CSS](https://tailwindcss.com/) (minimal config, tokens-driven)  
- **Validation / Schema**: [ArkType](https://arktype.io/)  
- **Package Manager**: [pnpm](https://pnpm.io/)  
- **Backend**: Custom Golang API (hosted on VPS with custom domain)  
- **Build Tool**: [Vite](https://vitejs.dev/)  
- **PWA**: Service worker + manifest integrated with SvelteKit  

No database is used in the frontend; all persistent data flows through the backend API.

---

## 🌐 Backend Integration

All communication happens via HTTPS requests to the **Golang backend**.  
Claude should assume:

- The backend provides a REST API (JSON-based).  
- Authentication (if any) is handled via tokens/cookies set by the backend.  
- Frontend never connects to third-party APIs directly — all calls go through the backend.

**Example fetch wrapper:**

\`\`\`ts
async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
	const res = await fetch(`${import.meta.env.VITE_API_URL}/${path}`, {
		credentials: 'include',
		headers: { 'Content-Type': 'application/json' },
		...options
	});
	if (!res.ok) throw new Error(await res.text());
	return res.json();
}
\`\`\`

---

## 🧩 App Structure

| Folder | Purpose |
|--------|----------|
| `/src/routes/` | Page routes (SvelteKit file-based routing) |
| `/src/lib/components/` | Reusable UI components (Bits-UI primitives) |
| `/src/lib/stores/` | Svelte stores for global state |
| `/src/lib/schemas/` | ArkType validation schemas |
| `/src/lib/utils/` | Generic utilities and helpers |
| `/src/lib/types/` | Shared TypeScript types |
| `/src/lib/api/` | API request wrappers to the Go backend |
| `/static/` | PWA assets (icons, manifest.json, etc.) |

---

## 🎨 Design System

This app follows **a token-based minimal design system**.  
Claude should enforce the following guidelines:

- Use **Tailwind CSS variables** for colors, spacing, and typography  
- Default spacing scale: `0.25rem → 2rem`  
- Font scale follows Tailwind defaults  
- Use **Bits-UI components** whenever possible for consistency  
- Avoid introducing raw HTML elements where Bits-UI equivalents exist  
- Keep **animations subtle** and **transitions under 200ms**

**Example Button Structure:**

\`\`\`svelte
<script lang="ts">
	import { Button } from 'bits-ui';
	export let icon: string | null = null;
</script>

<Button class="gap-[var(--gap)]" style="--gap: {icon ? '1rem' : '0rem'}">
	{#if icon}<Icon name={icon} />{/if}
	<span><slot /></span>
</Button>
\`\`\`

---

## 🧱 Schema & Validation

Use **ArkType** for all client-side data validation.

**Example:**

\`\`\`ts
import { type } from 'arktype';

export const LoginSchema = type({
	email: 'string.email',
	password: 'string.min(8)'
});
\`\`\`

Schemas should be defined in `/src/lib/schemas` and reused across components or API calls for consistency.

---

## 🧠 State Management

Use **Svelte stores** for reactive global state.  
Stores should be defined in `/src/lib/stores`.

**Pattern Example:**

\`\`\`ts
// src/lib/stores/session.ts
import { writable } from 'svelte/store';

export const session = writable<{ user: string | null }>({ user: null });
\`\`\`

---

## 🔍 Visual Development Workflow

Claude should:

1. Identify changed components and corresponding routes  
2. Launch `pnpm dev` and visit affected routes  
3. Check design compliance (alignment, spacing, accessibility, responsiveness)  
4. Use Playwright or MCP commands if configured to capture screenshots or console output  
5. Verify the change visually at:
   - Mobile (375px)
   - Tablet (768px)
   - Desktop (1440px)

**Design Checklist**

- [ ] Consistent use of Bits-UI components  
- [ ] Proper spacing and layout alignment  
- [ ] Works on mobile, tablet, and desktop  
- [ ] Accessible (ARIA, keyboard navigation)  
- [ ] Smooth transitions (≤200ms)  
- [ ] No visible overflow or layout shift  

---

## 🧪 Testing Guidelines

Testing is minimal at this stage.  
Claude should prioritize **manual and visual verification** over formal unit tests.

If automated tests are added later:
- Use **Playwright** for end-to-end UI testing  
- Use **Vitest** for unit and integration testing  

---

## 🚀 Deployment Notes

- The app builds as a **static PWA** with SvelteKit’s adapter (Cloudflare, Node, or static depending on VPS setup).  
- Deployment pipeline handled via **GitHub Actions**.  
- The Golang backend is deployed separately but accessible via environment variable `VITE_API_URL`.  

**Example `.env` (local):**

\`\`\`env
VITE_API_URL=https://api.yourdomain.com
PUBLIC_APP_NAME=YourAppName
\`\`\`

---

## 🧩 Claude Agents

Recommended Claude agents for this repository:

### `/claude/agents/ui-agent.md`
Handles UI implementation using Bits-UI, Tailwind, and Svelte conventions.

### `/claude/agents/design-review-agent.md`
Performs visual reviews, accessibility checks, and layout consistency validation.

### `/claude/agents/schema-agent.md`
Validates ArkType schemas, ensures consistent type use across API calls and Svelte components.

---

## 🗂 Context References

- `/context/design-principles.md` — Design principles and spacing rules  
- `/context/api-contracts.md` — Frontend/backend API interfaces  
- `/context/pwa-config.md` — Manifest and service worker setup  
- `/context/tokens.md` — Color, font, and spacing tokens  

---

## ✅ Claude Priorities

When generating or editing code, Claude should:

1. Follow **SvelteKit + TypeScript** syntax and conventions  
2. Use **Bits-UI + Tailwind** for all visual components  
3. Validate all data structures with **ArkType**  
4. Never directly access a database — all requests go through the **Go backend**  
5. Ensure accessibility and responsive design by default  
6. Keep components **isolated, reusable, and typed**
\`\`\`
---

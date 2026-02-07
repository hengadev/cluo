# For the photo or image management

## Broad idea:

I was thinking of using AI but instead of talking about editing, we should talk about annotating the photos.

## From ChatGPT

## AI Design Prompt — Add Interaction & Stateful Behavior

> **Context**
>
> This page already exists inside a larger desktop application.
> The sidebar, topbar, typography, spacing system, color palette, icons, and general layout conventions must remain consistent with the rest of the application.
> Do **not** redesign the navigation or global UI.
>
> The goal of this iteration is **not visual redesign**, but **adding realistic, explorable interactions and stateful behavior** to the existing Image Management interface.

---

> **General Interaction Principles**
>
> Treat this page as a **real operational tool**, not a static mockup.
>
> * Any UI element that looks interactive **must behave as interactive**.
> * Avoid “always-selected” or “frozen” states unless explicitly justified.
> * When behavior cannot be implemented, visually represent **multiple states** (e.g. duplicated frames, alternate panels, or state variants).
> * Favor discoverability: users must be able to explore different paths without guessing.

---

> **Tabs & Navigation States**
>
> * Tabs such as **“Evidence Pool”** and **“Final Selection”** must be fully selectable.
> * Switching tabs must:
>
>   * Update visible content
>   * Preserve internal state per tab when appropriate
> * Show both active and inactive tab states clearly.
> * The design should make it obvious that tabs are clickable and mutually exclusive.

---

> **Image Grid & Selection Behavior**
>
> The image grid must support realistic selection workflows:
>
> * No image selected by default **or** a clearly indicated initial selection state
> * Single selection
> * Multi-selection
> * Deselection
>
> Selected images must have clear visual feedback (overlay, border, checkmark, etc.).
>
> Hover, focus, and selected states should all be visually distinct.
>
> The user must be able to:
>
> * Select different images
> * Change selection
> * Clear selection
>
> Avoid locking the interface into a single selection state.

---

> **Actions & Buttons**
>
> All action elements (e.g. **Import Photos**, **Add to Final Selection**, **Remove**, **Annotate**, **Reorder**) must imply a result.
>
> * Clicking **Import Photos** should open or suggest:
>
>   * A file picker
>   * A modal
>   * Or a placeholder interaction state
>
> Buttons should visually respond to:
>
> * Hover
> * Active
> * Disabled states (when contextually relevant)
>
> No primary action button should appear inert or decorative.

---

> **State Representation**
>
> Explicitly represent multiple UI states, such as:
>
> * Empty state (no photos)
> * Loading or importing state
> * Selection active vs inactive
> * Evidence Pool vs Final Selection
>
> If necessary, duplicate sections or frames to show how the UI changes between states.

---

> **Consistency & Scope**
>
> * Do not introduce new global patterns that conflict with other pages.
> * Reuse existing components and interaction conventions from the application wherever possible.
> * This page must feel like a natural extension of the existing system, not a standalone prototype.

---

> **Outcome Expectation**
>
> The result should allow a reviewer to:
>
> * Click through the page mentally
> * Understand what happens when they interact with elements
> * Explore alternative states without being blocked by static UI

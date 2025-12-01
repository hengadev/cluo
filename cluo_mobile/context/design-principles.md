# Design Principles for PI Mobile Recording App

## I. Core Design Philosophy & Strategy

- [ ] **Users First:** Focus on investigators and their workflow. Minimize friction in creating, reviewing, and sending case recordings.
- [ ] **Precision & Clarity:** Recordings, transcripts, and AI summaries must be clear, accurate, and easy to interpret.
- [ ] **Speed & Responsiveness:** Optimize for fast recording, playback, and AI processing feedback.
- [ ] **Simplicity:** Mobile UI should be uncluttered. Provide only relevant actions for each screen.
- [ ] **Focus & Efficiency:** Reduce cognitive load: one main action per screen where possible (Record, Send, Review).
- [ ] **Consistency:** Maintain a uniform design language for colors, typography, buttons, and input components.
- [ ] **Accessibility:** Ensure text legibility, high-contrast colors, and keyboard/assistive device compatibility.
- [ ] **Thoughtful Defaults:** Pre-select recording formats, microphone levels, and AI summarization options when possible.

## II. Design System Foundation (Tokens & Core Components)

- [ ] **Color Palette:**
  - Primary: Highlight main actions like record, send, and AI summarize.
  - Semantic Colors:  
    - Success → completed upload or processed AI summary.  
    - Warning → microphone issues or network latency.  
    - Destructive → cancel recording or delete segment.
  - Dark Mode Palette: Ensure contrast in low-light fieldwork environments.
- [ ] **Typography:**
  - Primary font: `Inter` for readability.  
  - Modular scale: H1-H4 for titles, Body Large/Medium for transcript, Body Small for timestamps.
  - Monospace for timestamps and AI-generated code/segments (`JetBrains Mono`).
- [ ] **Spacing & Layout:**
  - Base unit: 4px/8px.  
  - Consistent padding/margin for recording controls, transcripts, and action buttons.
- [ ] **Border Radius:**  
  - Small: buttons, input fields, progress bars.  
  - Medium: cards, modals, AI summaries.
- [ ] **Core Components:**
  - Buttons: primary (record/send), destructive (delete), secondary (cancel/redo), ghost (optional).
  - Input fields: text, textarea, search for cases.
  - Badges: recording status, transcript segment status.
  - Cards: AI summary, recorded segments.
  - Progress bars: recording duration, AI processing.
  - Modals/Dialogs: confirm upload, delete, AI summary settings.
  - Tabs/Navigation: Cases, Recordings, AI Summaries.
  - Tooltips: microphone tips, AI options.
  - Icons: microphone, send, play, pause, delete.
  - Avatars: optional for user identification or case-specific images.

## III. Layout, Visual Hierarchy & Structure

- [ ] **Mobile-First Grid:** Stack content vertically. Use consistent horizontal spacing.
- [ ] **Whitespace:** Separate recording controls, transcripts, and AI summaries for clarity.
- [ ] **Visual Hierarchy:**  
  - Recording controls prominent at the bottom of the screen.  
  - AI summary cards above the transcript list.  
  - Timestamps and speaker info clearly distinguishable.
- [ ] **Consistent Alignment:** Left-align transcript text, center buttons for actions.
- [ ] **Persistent Elements:**  
  - Bottom bar: main navigation (Cases, Recordings, AI Summaries).  
  - Top bar (optional): search or quick actions.

## IV. Interaction Design & Animations

- [ ] **Micro-interactions:**  
  - Button feedback on tap, hold, hover.  
  - Recording visual pulse, waveform animation during live recording.
- [ ] **Loading States:**  
  - Skeleton cards while AI summary generates.  
  - Spinners or progress bars for uploads and processing.
- [ ] **Smooth Transitions:** Expand/collapse transcript segments, AI summary updates.
- [ ] **Keyboard & Accessibility:** Ensure buttons, inputs, and modals are fully navigable.

## V. Module-Specific Design Tactics

### A. Recording Module

- [ ] **Clear Record Control:** Prominent record button; destructive cancel button.
- [ ] **Real-Time Feedback:** Audio waveform, timer, recording status badge.
- [ ] **Pause/Resume:** Allow pausing without losing data.
- [ ] **Upload & AI Summary:** Provide feedback on processing state (spinner, progress bar).
- [ ] **Undo/Redo:** Easy cancellation or redo of recording segments.

### B. Transcript & AI Summary Module

- [ ] **Transcript Readability:** Speaker labels, timestamps, scrollable container.
- [ ] **AI Summarization Cards:** Collapsible or expandable; clear summary text; actionable buttons (approve, edit, send to report).
- [ ] **Highlight & Copy:** Allow text selection and copying for case notes.
- [ ] **Status Indicators:** Badges for AI processing, user review, or completion.

### C. Case Management / Navigation Module

- [ ] **Quick Access to Cases:** List of cases with search and filter options.
- [ ] **Consistent Status Indicators:** Recordings, AI summaries, and review state.
- [ ] **Easy Navigation:** Bottom navigation bar with clear icons and labels.

## VI. CSS & Styling Architecture

- [ ] **Utility-First (Tailwind CSS):** Use defined tokens in Tailwind config (`colors`, `spacing`, `radii`).
- [ ] **Bits-UI Integration:** Use `bitsUiTheme` to maintain consistent styling across buttons, cards, inputs, badges.
- [ ] **Component Classes:** Leverage helper functions (`getButtonClasses`, `getCardClasses`) for maintainability.
- [ ] **Animations:** Use pre-defined keyframes for recording pulse, waveform, slide/fade-in components.

## VII. Best Practices

- [ ] **Iterative Testing:** Test recordings, AI summaries, and mobile UI with real users.
- [ ] **Responsive Design:** Ensure smooth experience on small phones and tablets.
- [ ] **Documentation:** Keep all component usage, colors, spacing, and tokens documented.
- [ ] **Accessibility & Privacy:** Ensure transcripts and recordings respect privacy; use encrypted communication and safe storage.


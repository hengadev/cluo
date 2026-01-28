<script lang="ts">
    import { Camera, Grid3x3, Columns, FileText, ArrowDownWideNarrow, Layers, CheckSquare } from "@lucide/svelte";

    export type ViewMode = "grid-compact" | "grid-comfortable";
    export type SortMode = "newest" | "oldest" | "filename";
    export type LayoutMode = "library" | "split" | "report";

    interface Props {
        selectMode: boolean;
        viewMode: ViewMode;
        sortMode: SortMode;
        layoutMode: LayoutMode;
        hasBurstGroups: boolean;
        onSelectModeToggle: () => void;
        onImport: () => void;
        onBurstGroupToggle: () => void;
        onViewModeChange: (mode: ViewMode) => void;
        onSortModeChange: (mode: SortMode) => void;
        onLayoutModeChange: (mode: LayoutMode) => void;
    }

    let {
        selectMode,
        viewMode,
        sortMode,
        layoutMode,
        hasBurstGroups,
        onSelectModeToggle,
        onImport,
        onBurstGroupToggle,
        onViewModeChange,
        onSortModeChange,
        onLayoutModeChange,
    }: Props = $props();

    let showViewMenu = $state(false);
    let showSortMenu = $state(false);
</script>

<div class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50">
    <div class="flex items-center gap-2 bg-background/95 backdrop-blur-sm border border-border-card rounded-full px-3 py-2 shadow-lg">
        <!-- Import -->
        <button
            class="w-10 h-10 rounded-full flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-all"
            onclick={onImport}
            title="Importer depuis la caméra"
        >
            <Camera size={18} />
        </button>

        <div class="w-px h-6 bg-border-card"></div>

        <!-- Select Mode -->
        <button
            class="w-10 h-10 rounded-full flex items-center justify-center transition-all {selectMode
                ? 'bg-primary text-primary-foreground'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted'}"
            onclick={onSelectModeToggle}
            title="Mode sélection"
        >
            <CheckSquare size={18} />
        </button>

        <!-- Burst Group -->
        <button
            class="w-10 h-10 rounded-full flex items-center justify-center transition-all {hasBurstGroups
                ? 'bg-primary text-primary-foreground'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted'}"
            onclick={onBurstGroupToggle}
            title="Groupes en mode rafale"
        >
            <Layers size={18} />
        </button>

        <div class="w-px h-6 bg-border-card"></div>

        <!-- View Mode -->
        <div class="relative">
            <button
                class="w-10 h-10 rounded-full flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-all"
                onclick={() => showViewMenu = !showViewMenu}
                title="Mode d'affichage"
            >
                {#if viewMode === "grid-compact"}
                    <Grid3x3 size={18} />
                {:else}
                    <List size={18} />
                {/if}
            </button>

            {#if showViewMenu}
                <div class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 bg-background border border-border-card rounded-lg shadow-lg py-1 min-w-[150px]">
                    <button
                        class="w-full px-4 py-2 text-left text-sm hover:bg-muted flex items-center gap-2 {viewMode === 'grid-compact'
                            ? 'text-foreground bg-muted/50'
                            : 'text-muted-foreground'}"
                        onclick={() => {
                            onViewModeChange("grid-compact");
                            showViewMenu = false;
                        }}
                    >
                        <Grid3x3 size={14} />
                        Grille compacte
                    </button>
                    <button
                        class="w-full px-4 py-2 text-left text-sm hover:bg-muted flex items-center gap-2 {viewMode === 'grid-comfortable'
                            ? 'text-foreground bg-muted/50'
                            : 'text-muted-foreground'}"
                        onclick={() => {
                            onViewModeChange("grid-comfortable");
                            showViewMenu = false;
                        }}
                    >
                        <List size={14} />
                        Grille large
                    </button>
                </div>
            {/if}
        </div>

        <!-- Sort Mode -->
        <div class="relative">
            <button
                class="w-10 h-10 rounded-full flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-muted transition-all"
                onclick={() => showSortMenu = !showSortMenu}
                title="Trier par"
            >
                <ArrowDownWideNarrow size={18} />
            </button>

            {#if showSortMenu}
                <div class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 bg-background border border-border-card rounded-lg shadow-lg py-1 min-w-[150px]">
                    <button
                        class="w-full px-4 py-2 text-left text-sm hover:bg-muted {sortMode === 'newest'
                            ? 'text-foreground bg-muted/50'
                            : 'text-muted-foreground'}"
                        onclick={() => {
                            onSortModeChange("newest");
                            showSortMenu = false;
                        }}
                    >
                        Plus récent
                    </button>
                    <button
                        class="w-full px-4 py-2 text-left text-sm hover:bg-muted {sortMode === 'oldest'
                            ? 'text-foreground bg-muted/50'
                            : 'text-muted-foreground'}"
                        onclick={() => {
                            onSortModeChange("oldest");
                            showSortMenu = false;
                        }}
                    >
                        Plus ancien
                    </button>
                    <button
                        class="w-full px-4 py-2 text-left text-sm hover:bg-muted {sortMode === 'filename'
                            ? 'text-foreground bg-muted/50'
                            : 'text-muted-foreground'}"
                        onclick={() => {
                            onSortModeChange("filename");
                            showSortMenu = false;
                        }}
                    >
                        Nom de fichier
                    </button>
                </div>
            {/if}
        </div>

        <div class="w-px h-6 bg-border-card"></div>

        <!-- Layout Mode Toggle -->
        <div class="flex bg-muted rounded-full p-1">
            <button
                class="w-9 h-9 rounded-full flex items-center justify-center transition-all {layoutMode === 'library'
                    ? 'bg-background text-foreground shadow-sm'
                    : 'text-muted-foreground hover:text-foreground'}"
                onclick={() => onLayoutModeChange("library")}
                title="Bibliothèque uniquement"
            >
                <Grid3x3 size={16} />
            </button>
            <button
                class="w-9 h-9 rounded-full flex items-center justify-center transition-all {layoutMode === 'split'
                    ? 'bg-background text-foreground shadow-sm'
                    : 'text-muted-foreground hover:text-foreground'}"
                onclick={() => onLayoutModeChange("split")}
                title="Vue partagée"
            >
                <Columns size={16} />
            </button>
            <button
                class="w-9 h-9 rounded-full flex items-center justify-center transition-all {layoutMode === 'report'
                    ? 'bg-background text-foreground shadow-sm'
                    : 'text-muted-foreground hover:text-foreground'}"
                onclick={() => onLayoutModeChange("report")}
                title="Rapport uniquement"
            >
                <FileText size={16} />
            </button>
        </div>
    </div>
</div>

<!-- Click outside to close menus -->
<svelte:window onclick={() => { showViewMenu = false; showSortMenu = false; }} />

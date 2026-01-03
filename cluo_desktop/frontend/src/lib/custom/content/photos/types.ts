export const LAYOUT_STATES = {
    List: "List",
    Grid: "Grid",
} as const;

export type LayoutState = typeof LAYOUT_STATES[keyof typeof LAYOUT_STATES];
export const LAYOUT_STATES_ARRAY = Object.values(LAYOUT_STATES);

export const SORT_STATES = {
    NewestFirst: "Newest first",
    OldestFirst: "Oldest first",
    SelectedFirst: "Selected first",
    NonSelectedFirst: "NonSelected first",
} as const;

export type SortState = typeof SORT_STATES[keyof typeof SORT_STATES];
export const SORT_STATES_ARRAY = Object.values(SORT_STATES);

export type CardType = {
    title: string;
    tags: string[];
};

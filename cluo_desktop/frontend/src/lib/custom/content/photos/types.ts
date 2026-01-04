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

export const FILTER_STATES = {
    All: "All",
    Selected: "Selected",
    NotSelected: "Not selected",
} as const;

export type FilterState = typeof FILTER_STATES[keyof typeof FILTER_STATES];
export const FILTER_STATES_ARRAY = Object.values(FILTER_STATES);

export type Image = {
    id: string;
    caseId: string;
    url: string;
    filename: string;
    filesize: number;
    caption: string;
    isPublished: boolean;
    createdAt: string;
};

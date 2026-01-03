export const LAYOUT_STATES = {
    List: "List",
    Grid: "Grid",
} as const;

export type LayoutState = typeof LAYOUT_STATES[keyof typeof LAYOUT_STATES];
export const LAYOUT_STATES_ARRAY = Object.values(LAYOUT_STATES);


export type CardType = {
    title: string;
    tags: string[];
};

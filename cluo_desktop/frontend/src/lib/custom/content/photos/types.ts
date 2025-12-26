export const VIEW_STATES = {
    List: "List",
    Grid: "Grid",
} as const;

export type ViewState = typeof VIEW_STATES[keyof typeof VIEW_STATES];
export const View_STATES_ARRAY = Object.values(VIEW_STATES);


export type CardType = {
    title: string;
    tags: string[];
};

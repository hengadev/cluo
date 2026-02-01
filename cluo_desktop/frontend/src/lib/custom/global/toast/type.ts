export const TOAST_LEVELS = {
    Info: 'Info',
    Error: 'Error',
    Alert: 'Alert',
    Warning: 'Warning',
} as const;
export type ToastLevel = typeof TOAST_LEVELS[keyof typeof TOAST_LEVELS];
export const TOAST_LEVELS_ARRAY = Object.values(TOAST_LEVELS);

export type Toast = {
    id: string;
    level: ToastLevel;
    title: string;
    message: string;
};

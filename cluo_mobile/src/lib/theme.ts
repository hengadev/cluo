export { default as theme } from 'tailwindcss/defaultTheme';

// Theme configuration for bits-ui components
export const bitsUiTheme = {
    // Button variants that match the UI design
    button: {
        base: 'inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
        variants: {
            variant: {
                default: 'bg-primary text-primary-foreground hover:bg-primary/90',
                destructive: 'bg-destructive text-destructive-foreground hover:bg-destructive/90',
                outline: 'border border-input bg-background hover:bg-accent hover:text-accent-foreground',
                secondary: 'bg-secondary text-secondary-foreground hover:bg-secondary/80',
                ghost: 'hover:bg-accent hover:text-accent-foreground',
                link: 'text-primary underline-offset-4 hover:underline',
                // Custom variants for the recording app
                record: 'bg-destructive hover:bg-destructive/90 text-destructive-foreground hover:scale-105 active:scale-95 transition-all duration-200',
                recordActive: 'bg-destructive hover:bg-destructive/90 text-destructive-foreground animate-pulse-record',
                transcript: 'bg-transcript hover:bg-muted/80 text-transcript-text border border-border',
                ai: 'bg-gradient-to-r from-primary to-primary/80 text-primary-foreground hover:from-primary/90 hover:to-primary/70',
            },
            size: {
                default: 'h-10 px-4 py-2',
                sm: 'h-9 rounded-md px-3',
                lg: 'h-11 rounded-md px-8',
                icon: 'h-10 w-10',
                // Custom sizes for recording UI
                record: 'h-16 w-16 rounded-full',
                recordSmall: 'h-12 w-12 rounded-full',
            },
        },
        defaultVariants: {
            variant: 'default',
            size: 'default',
        },
    },

    // Card component styling
    card: {
        base: 'rounded-lg border bg-card text-card-foreground shadow-sm',
        header: 'flex flex-col space-y-1.5 p-6',
        title: 'text-2xl font-semibold leading-none tracking-tight',
        description: 'text-sm text-muted-foreground',
        content: 'p-6 pt-0',
        footer: 'flex items-center p-6 pt-0',
    },

    // Input component styling
    input: {
        base: 'flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50',
        variants: {
            variant: {
                default: '',
                transcript: 'bg-muted/50 border-border resize-none focus-visible:ring-primary/20',
                search: 'pl-10 bg-muted/50 border-border focus-visible:ring-primary/20',
            },
        },
    },

    // Badge component styling
    badge: {
        base: 'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
        variants: {
            variant: {
                default: 'border-transparent bg-primary text-primary-foreground hover:bg-primary/80',
                secondary: 'border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80',
                destructive: 'border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80',
                outline: 'text-foreground',
                // Custom variants for the app
                recording: 'border-transparent bg-destructive text-destructive-foreground animate-pulse-record',
                processing: 'border-transparent bg-warning text-warning-foreground',
                success: 'border-transparent bg-success text-success-foreground',
                transcript: 'border-transparent bg-muted text-muted-foreground',
                speaker: 'border-primary/30 bg-primary/10 text-primary',
            },
        },
        defaultVariants: {
            variant: 'default',
        },
    },

    // Dialog/Modal styling
    dialog: {
        overlay: 'fixed inset-0 z-50 bg-background/80 backdrop-blur-sm data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0',
        content: 'fixed left-[50%] top-[50%] z-50 grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] sm:rounded-lg',
        header: 'flex flex-col space-y-1.5 text-center sm:text-left',
        title: 'text-lg font-semibold leading-none tracking-tight',
        description: 'text-sm text-muted-foreground',
        footer: 'flex flex-col-reverse sm:flex-row sm:justify-end sm:space-x-2',
    },

    // Tabs component styling
    tabs: {
        list: 'inline-flex h-10 items-center justify-center rounded-md bg-muted p-1 text-muted-foreground',
        trigger: 'inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1.5 text-sm font-medium ring-offset-background transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow-sm',
        content: 'mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
    },

    // Custom component themes for the recording app
    recording: {
        button: {
            base: 'btn-record',
            active: 'btn-record recording',
            icon: 'btn-record record',
            iconActive: 'btn-record recordActive',
        },
        timer: {
            base: 'recording-timer',
        },
        waveform: {
            base: 'voice-indicator',
            bar: 'voice-bar',
        },
    },

    transcript: {
        container: 'transcript-container transcript-scroll max-h-96 overflow-y-auto',
        segment: 'transcript-segment',
        timestamp: 'text-transcript-timestamp text-xs font-mono',
        speaker: 'text-transcript-speaker font-semibold text-sm',
        text: 'text-transcript leading-relaxed',
    },

    aiChat: {
        container: 'space-y-4',
        message: 'ai-message',
        userMessage: 'ai-message user',
        assistantMessage: 'ai-message assistant',
        input: 'input variant-transcript',
    },

    progress: {
        bar: 'w-full bg-secondary rounded-full h-2 overflow-hidden',
        fill: 'h-full bg-gradient-to-r from-primary to-primary/80 rounded-full transition-all duration-500 ease-out',
        circle: 'relative w-12 h-12',
        circleBg: 'absolute inset-0 rounded-full bg-secondary',
        circleFill: 'absolute inset-0 rounded-full bg-gradient-to-r from-primary to-primary/80',
    },

    // Navigation styling
    navigation: {
        base: 'flex items-center space-x-1',
        item: 'inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 data-[active]:bg-background data-[state=open]:bg-muted hover:bg-accent hover:text-accent-foreground h-10 px-4 py-2',
        itemActive: 'bg-background text-accent-foreground hover:bg-muted',
    },
} as const;

// Export color utilities for easy access
export const colors = {
    primary: {
        50: '#ecfeff',
        100: '#cffafe',
        200: '#a5f3fc',
        300: '#67e8f9',
        400: '#22d3ee',
        500: '#06b6d4',
        600: '#0891b2',
        700: '#0e7490',
        800: '#155e75',
        900: '#164e63',
        950: '#083344',
    },
    // Add other color ranges as needed
} as const;

// Helper functions for styling
export const getButtonClasses = (variant: string = 'default', size: string = 'default') => {
    const buttonConfig = bitsUiTheme.button;
    const variantClass = buttonConfig.variants.variant[variant as keyof typeof buttonConfig.variants.variant] || '';
    const sizeClass = buttonConfig.variants.size[size as keyof typeof buttonConfig.variants.size] || '';

    return `${buttonConfig.base} ${variantClass} ${sizeClass}`;
};

export const getCardClasses = (part: 'base' | 'header' | 'title' | 'description' | 'content' | 'footer' = 'base') => {
    return bitsUiTheme.card[part] || bitsUiTheme.card.base;
};

export const getBadgeClasses = (variant: string = 'default') => {
    const badgeConfig = bitsUiTheme.badge;
    const variantClass = badgeConfig.variants.variant[variant as keyof typeof badgeConfig.variants.variant] || '';

    return `${badgeConfig.base} ${variantClass}`;
};

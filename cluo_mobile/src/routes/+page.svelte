<script lang="ts">
    import { Button, Input, Checkbox } from '$lib/components/ui';
</script>

<div class="min-h-screen bg-white">
    <!-- Mobile: Full screen card -->
    <div class="md:hidden min-h-screen flex items-start justify-center pt-12 p-4">
        <div class="w-full max-w-sm">
            <!-- Mobile Login Form -->
            <div class="bg-white p-6 space-y-6">
                <!-- Logo placeholder -->
                <div class="flex justify-center">
                    <div
                        class="w-16 h-16 rounded-full bg-muted flex items-center justify-center"
                    >
                        <svg
                            class="w-8 h-8 text-muted-foreground"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                            />
                        </svg>
                    </div>
                </div>

                <!-- Header -->
                <div class="text-center">
                    <h1 class="text-2xl font-semibold text-foreground mb-2">
                        Connectez-vous à votre compte
                    </h1>
                    <p class="text-muted-foreground">
                        Entrez vos identifiants pour continuer
                    </p>
                </div>

                <!-- Form -->
                <form class="space-y-4">
                    <!-- Email field -->
                    <div class="space-y-2">
                        <label
                            for="email"
                            class="text-sm font-medium text-gray-900"
                            >Adresse e-mail</label
                        >
                        <Input
                            id="email"
                            type="email"
                            placeholder="Entrez votre adresse e-mail"
                            size="md"
                        />
                    </div>

                    <!-- Password field -->
                    <div class="space-y-2">
                        <label
                            for="password"
                            class="text-sm font-medium text-gray-900"
                            >Mot de passe</label
                        >
                        <Input
                            id="password"
                            type="password"
                            placeholder="Entrez votre mot de passe"
                            size="md"
                        />
                    </div>

                    <!-- Remember me and forgot password -->
                    <div class="flex items-center justify-between">
                        <div class="flex items-center space-x-2">
                            <Checkbox
                                id="remember"
                                name="remember"
                                size="md"
                            />
                            <label
                                for="remember"
                                class="text-sm text-muted-foreground cursor-pointer"
                            >
                                Se souvenir de moi
                            </label>
                        </div>
                        <button
                            type="button"
                            class="text-sm text-primary hover:underline cursor-pointer"
                        >
                            Mot de passe oublié ?
                        </button>
                    </div>

                    <!-- Sign in button -->
                    <Button type="submit" variant="primary" size="md" class="w-full">
                        Se connecter
                    </Button>
                </form>
            </div>
        </div>
    </div>

    <!-- Desktop: Split screen layout -->
    <div class="hidden md:flex min-h-screen">
        <!-- Left side: Login form -->
        <div class="flex-1 flex items-center justify-center p-8">
            <div class="w-full max-w-md">
                <!-- Desktop Login Form -->
                <div class="space-y-8">
                    <!-- Logo placeholder -->
                    <div class="flex justify-center">
                        <div
                            class="w-20 h-20 rounded-full bg-muted flex items-center justify-center"
                        >
                            <svg
                                class="w-10 h-10 text-muted-foreground"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                                />
                            </svg>
                        </div>
                    </div>

                    <!-- Header -->
                    <div class="text-center">
                        <h1 class="text-3xl font-semibold text-foreground mb-3">
                            Connectez-vous à votre compte
                        </h1>
                        <p class="text-muted-foreground text-lg">
                            Entrez vos identifiants pour continuer
                        </p>
                    </div>

                    <!-- Form -->
                    <form class="space-y-6">
                        <!-- Email field -->
                        <div class="space-y-3">
                            <label
                                for="email-desktop"
                                class="text-base font-medium text-gray-900"
                                >Adresse e-mail</label
                            >
                            <Input
                                id="email-desktop"
                                type="email"
                                placeholder="Entrez votre adresse e-mail"
                                size="lg"
                            />
                        </div>

                        <!-- Password field -->
                        <div class="space-y-3">
                            <label
                                for="password-desktop"
                                class="text-base font-medium text-gray-900"
                                >Mot de passe</label
                            >
                            <Input
                                id="password-desktop"
                                type="password"
                                placeholder="Entrez votre mot de passe"
                                size="lg"
                            />
                        </div>

                        <!-- Remember me and forgot password -->
                        <div class="flex items-center justify-between">
                            <div class="flex items-center space-x-3">
                                <Checkbox
                                    id="remember-desktop"
                                    name="remember-desktop"
                                    size="lg"
                                />
                                <label
                                    for="remember-desktop"
                                    class="text-base text-muted-foreground cursor-pointer"
                                >
                                    Se souvenir de moi
                                </label>
                            </div>
                            <button
                                type="button"
                                class="text-base text-primary hover:underline cursor-pointer"
                            >
                                Mot de passe oublié ?
                            </button>
                        </div>

                        <!-- Sign in button -->
                        <Button type="submit" variant="primary" size="lg" class="w-full">
                            Se connecter
                        </Button>
                    </form>

                    <!-- Sign up link -->
                    <div class="text-center pt-6 border-t border-border">
                        <p class="text-muted-foreground">
                            Pas de compte ?
                            <button
                                type="button"
                                class="text-primary hover:underline font-medium ml-1 cursor-pointer"
                            >
                                S'inscrire
                            </button>
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Right side: Placeholder background -->
        <div class="flex-1 bg-gray-100 relative overflow-hidden">
            <div class="absolute inset-0 flex items-center justify-center">
                <div class="text-center text-gray-600">
                    <div
                        class="w-24 h-24 mx-auto mb-6 rounded-full bg-gray-200 flex items-center justify-center"
                    >
                        <svg
                            class="w-12 h-12 text-gray-400"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                            />
                        </svg>
                    </div>
                    <h2 class="text-2xl font-semibold mb-2">Bon retour</h2>
                    <p class="text-gray-500">Votre espace de travail sécurisé vous attend</p>
                </div>
            </div>
        </div>
    </div>
</div>

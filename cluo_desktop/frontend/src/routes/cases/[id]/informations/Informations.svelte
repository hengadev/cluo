<script lang="ts">
    import Tags from "./tags.svelte";
    import Input from "$lib/custom/global/Input.svelte";

    import { Briefcase } from "@lucide/svelte";

    // NOTE: here is the Case thing, missing the Case Subject ID field in this
    // type Case struct {
    // 	ID                uuid.UUID  `db:"id"`
    // 	Title             string     `encx:"encrypt" db:"title_encrypted"`
    // 	Description       string     `encx:"encrypt" db:"description_encrypted"`
    // 	ClientID          uuid.UUID  `db:"client_id"`
    // 	AssignedContactID *uuid.UUID `db:"assigned_contact_id"`
    // 	ExternalReference *string    `encx:"encrypt" db:"external_reference"`
    // 	CaseType          string     `db:"case_type"`
    // 	Status            CaseStatus `encx:"encrypt" db:"status_encrypted"`
    // 	CreatedAt         time.Time  `db:"created_at"`
    // 	UpdatedAt         time.Time  `encx:"encrypt" db:"updated_at_encrypted"`
    // }
    // NOTE:
    // The information that I need :
    // - title
    // - description;
    // - client;
    // - contact;
    // - external reference if there are any;
    // - case type reference if there are any;
</script>

<div class="p-8">
    <Tags />
    <div class="flex gap-30 !mt-8">
        {@render header()}
        {@render client()}
        {@render project()}
    </div>
</div>

{#snippet client()}
    <div class="border-1 border-dark-50 rounded-input p-6 grid gap-2 animate-fade-in hover:shadow-md transition-shadow duration-300" style="animation-delay: 200ms;">
        <div class="flex justify-between items-center">
            <p class="text-dark-600">CLIENT DETAILS</p>
            <Briefcase />
        </div>
        <div class="flex items-center gap-2">
            <div class="bg-black w-10 h-10"></div>
            <div>
                <p>XYZ Legal</p>
                <p>Law Firm</p>
            </div>
        </div>
        <div class="bg-dark-50 h-[1px] w-full"></div>
        <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
                <div class="bg-black w-10 h-10"></div>
                <div>
                    <p>Jane DOE</p>
                    <p>Primary Contact (Paralegal)</p>
                </div>
            </div>
        </div>
    </div>
{/snippet}

{#snippet header()}
    <div class="grid gap-5 p-6 border-1 border-dark-50 rounded-input animate-fade-in hover:shadow-md transition-shadow duration-300" style="animation-delay: 100ms;">
        <div class="flex gap-4 items-center">
            <p class="bg-green-300 px-2 py-1 rounded-input">STATUS: OPEN</p>
            <p class="text-dark-400">Created Oct 24, 2023</p>
        </div>
        <p class="text-dark-900 font-extrabold text-3xl">Accident 5th Avenue</p>
        <div class="flex gap-4 text-xl items-center">
            <p class="text-dark-600">Case ID:</p>
            <p>#CASE-2023-849</p>
        </div>
    </div>
{/snippet}

{#snippet project()}
    <div class="border-1 border-dark-50 max-w-256 rounded-input flex-1 animate-fade-in hover:shadow-md transition-shadow duration-300" style="animation-delay: 300ms;">
        <div class="bg-dark-50 p-3">
            <p>Case details</p>
        </div>
        <div class="grid gap-8 p-4 rounded-b-input">
            <div class="grid grid-cols-2 gap-4">
                <div class="grid">
                    <p>Titre</p>
                    <Input placeholder="Affaire Ti Sonson" />
                </div>
                <div class="grid">
                    <p>Description</p>
                    <Input
                        placeholder="Une plus grosses affaires impliquant un Ti sonson"
                    />
                </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
                <div class="grid">
                    <p>Nom du client</p>
                    <Input placeholder="Cabinet DUPONT" />
                </div>
                <div class="grid">
                    <p>Nom du contact associe au client</p>
                    <Input placeholder="Maitre MARTIN" />
                </div>
            </div>
        </div>
    </div>
{/snippet}

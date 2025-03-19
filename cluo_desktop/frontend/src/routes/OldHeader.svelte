<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import {
        Dialog,
        DialogTrigger,
        DialogContent,
    } from "$lib/components/ui/dialog";
    import { Search } from "@lucide/svelte";

    let searchQuery = "";
    let searchResults: string[] = [];

    // Mock search function
    function performSearch(query: string) {
        // TODO: Replace this with your actual search logic
        const mockResults = ["Result 1", "Result 2", "Result 3"];
        searchResults = mockResults.filter((result) =>
            result.toLowerCase().includes(query.toLowerCase()),
        );
    }

    function handleSearch() {
        performSearch(searchQuery);
    }
</script>

<div class="content">
    <!-- Header with Search Bar -->
    <header class="p-4 border-b">
        <div class="flex items-center gap-2">
            <!-- Search Input -->
            <Input
                bind:value={searchQuery}
                placeholder="Search..."
                class="flex-1"
                on:keydown={(e) => e.key === "Enter" && handleSearch()}
            />

            <!-- Search Button -->
            <Dialog>
                <DialogTrigger asChild>
                    <Button on:click={handleSearch}>
                        <Search class="mr-2 h-4 w-4" />
                        Search
                    </Button>
                </DialogTrigger>

                <!-- Dialog with Search Results -->
                <DialogContent class="sm:max-w-[425px]">
                    <div class="p-4">
                        {#if searchResults.length > 0}
                            <ul>
                                {#each searchResults as result}
                                    <li class="py-2 border-b">{result}</li>
                                {/each}
                            </ul>
                        {:else}
                            <p class="text-gray-500">No results found.</p>
                        {/if}
                    </div>
                </DialogContent>
            </Dialog>
        </div>
    </header>
</div>

<style>
    .content {
        border: 2px solid red;
    }
</style>

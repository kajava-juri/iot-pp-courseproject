<script lang="ts">
    import type { Snippet } from "svelte";
    let {
        children,
        title = "",
        searchQuery = $bindable<string>(""),
    }: {
        children: Snippet;
        title?: string;
        searchQuery?: string;
    } = $props();
</script>

<div class="flex flex-col gap-6 w-full p-4">
    <div class="flex flex-col gap-2 max-h-96 overflow-y-scroll">
        <div class="flex items-center gap-2 sticky top-0 bg-white z-10 pt-2">
            <h3 class="text-2xl font-bold">{title}</h3>
            <input
                type="text"
                placeholder="Search rooms..."
                bind:value={searchQuery}
                class="px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
        </div>
        {@render children()}
    </div>
</div>

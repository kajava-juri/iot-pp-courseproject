<script lang="ts">
    import { onMount } from "svelte";
    import { fetcher } from "$utils/fetcher";
    import { tryParseRooms } from "$types/healthcare-db-types";
    import type { Room } from "$types/healthcare-db-types";
    import StatusComponent from "$components/StatusComponent.svelte";

    let allRooms: Room[] = [];
    let roomSearchQuery = "";
    let deviceSearchQuery = "";

    onMount(async () => {
        try {
            const response = await fetcher("/room");
            const parseResult = tryParseRooms(response.data);

            if (parseResult.success) {
                allRooms = parseResult.data;
                rooms = allRooms;
            } else {
                console.error("Failed to parse rooms:", parseResult.error);
            }
        } catch (error) {
            console.error("Failed to fetch rooms:", error);
        }
    });

    $: rooms = allRooms.filter((room) =>
        room.room_name.toLowerCase().includes(roomSearchQuery.toLowerCase()),
    );
</script>

<div class="flex flex-col gap-6 w-full p-4">
    <!-- Rooms Section -->
    <StatusComponent title="Rooms" bind:searchQuery={roomSearchQuery}>
        <div class="grid grid-cols-1 xl:grid-cols-5 lg:grid-cols-3 gap-3">
            {#each rooms as room}
                <div
                    class="border border-gray-200 p-4 aspect-square bg-secondary-300 flex flex-col justify-center items-center rounded-lg shadow-sm hover:shadow-md hover:border-gray-300 transition-all duration-200 cursor-pointer"
                >
                    <h2 class="text-lg font-semibold mb-2 text-gray-800">{room.room_name}</h2>
                    <p class="text-sm text-gray-600">Status: OK</p>
                </div>
            {/each}
        </div>
    </StatusComponent>

    <!-- Device Status Section -->
    <!-- <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2">
            <h1 class="text-2xl font-bold">Device Status</h1>
        </div>
        <input
            type="text"
            placeholder="Search devices..."
            bind:value={deviceSearchQuery}
            class="px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
    </div> -->
</div>

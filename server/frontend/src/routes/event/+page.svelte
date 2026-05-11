<script lang="ts">
	import { onMount } from 'svelte';
	import { TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, TableSearch } from 'flowbite-svelte';
	import { fetcher } from '$utils/fetcher';
	import { tryParseEvents } from '$types/healthcare-db-types';
	import type { Event } from '$types/healthcare-db-types';

	let events: Event[] = $state<Event[]>([]);
	let searchTerm = $state('');

	let filteredEvents = $derived.by(() => {
		if (!searchTerm) return events;

		const query = searchTerm.toLowerCase();
		return events.filter((event) =>
			event.type.toLowerCase().includes(query) ||
			String(event.device_id).includes(query) ||
			String(event.ID).includes(query)
		).sort((a, b) => b.CreatedAt.getTime() - a.CreatedAt.getTime());
	});

	onMount(async () => {
		try {
			const response = await fetcher('/event');
			const parsed = tryParseEvents(response.data);

			if (parsed.success) {
				events = parsed.data.sort((a, b) => b.CreatedAt.getTime() - a.CreatedAt.getTime());
			} else {
				console.error('Failed to parse events:', parsed.error);
			}
		} catch (error) {
			console.error('Failed to fetch events:', error);
		}
	});

	function formatDate(value: Date) {
		return value.toLocaleString();
	}
</script>

<div class="p-4 w-full">
	<h1 class="mb-4 text-2xl font-semibold">Events</h1>

	<TableSearch placeholder="Search by ID, type, or device" hoverable bind:inputValue={searchTerm}>
		<TableHead>
			<TableHeadCell>ID</TableHeadCell>
			<TableHeadCell>Type</TableHeadCell>
			<TableHeadCell>Device ID</TableHeadCell>
			<TableHeadCell>Patient</TableHeadCell>
			<TableHeadCell>Room</TableHeadCell>
			<TableHeadCell>Created At</TableHeadCell>
		</TableHead>
		<TableBody>
			{#if filteredEvents.length === 0}
				<TableBodyRow>
					<TableBodyCell colspan={6} class="text-center text-gray-500">
						No events available.
					</TableBodyCell>
				</TableBodyRow>
			{:else}
				{#each filteredEvents as event}
					<TableBodyRow>
						<TableBodyCell>{event.ID}</TableBodyCell>
						<TableBodyCell>{event.type}</TableBodyCell>
						<TableBodyCell>{event.device_id}</TableBodyCell>
						<TableBodyCell>{event.patient?.name || 'N/A'}</TableBodyCell>
						<TableBodyCell>{event.room?.room_name || 'N/A'}</TableBodyCell>
						<TableBodyCell>{formatDate(event.CreatedAt)}</TableBodyCell>
					</TableBodyRow>
				{/each}
			{/if}
		</TableBody>
	</TableSearch>
</div>

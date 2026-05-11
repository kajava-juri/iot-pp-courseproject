<script lang="ts">
  import { onMount } from 'svelte';
  import { fetcher } from '../../utils/fetcher';
  import { tryParseDevices, tryParseRooms, tryParsePatient, type Device, type Room, type Patient } from '$types/healthcare-db-types';
  import { addToast } from '$stores/toasts';
    import SearchableSelect from '$components/SearchableSelect.svelte';

  let devices: Device[] = $state([]);
  let rooms: Room[] = $state([]);
  let patients: Patient[] = $state([]);
  let editingId: number | null = $state(null);
  let editDescription = $state<string>('');
  let editRoomId = $state<string>('');
  let editPatientId = $state<string>('');

  async function loadData() {
    try {
      const devResp = (await fetcher('/device')).data;
      const parsed = tryParseDevices(devResp);
      if (parsed.success) devices = parsed.data;
      const roomResp = (await fetcher('/room')).data;
      const parsedRooms = tryParseRooms(roomResp);
      if (parsedRooms.success) rooms = parsedRooms.data;
      const patResp = (await fetcher('/patient')).data;
      const parsedPats = patResp.map((p: unknown) => tryParsePatient(p)).filter((r:any)=>r.success).map((r:any)=>r.data);
      patients = parsedPats;
    } catch (err) {
      console.error('Failed to load devices/rooms/patients', err);
    }
  }

  onMount(() => { loadData(); });

  function startEdit(d: Device) {
    editingId = d.ID;
    editDescription = d.description || '';
    editRoomId = d.room_id ? String(d.room_id) : '';
    editPatientId = d.patient_id ? String(d.patient_id) : '';
  }

  async function saveEdit() {
    if (editingId == null) return;
    const payload: any = {
      description: editDescription,
      room_id: editRoomId ? parseInt(editRoomId) : null,
      patient_id: editPatientId ? parseInt(editPatientId) : null,
    };

    try {
      const res = await fetcher(`/device/${editingId}`, 'put', payload);
      addToast('Device updated', 'success');
      editingId = null;
      await loadData();
    } catch (err:any) {
      console.error('Failed to update device', err);
      addToast(`Failed to update device: ${err?.data || err.message}`, 'error');
    }
  }

</script>

<div class="p-6">
  <h1 class="text-2xl font-bold mb-4">Devices</h1>
  <div class="space-y-4">
    {#each devices as d (d.ID)}
      <div class="p-4 border rounded flex items-start justify-between">
        <div>
          <div class="text-sm font-semibold">{d.device_name} — {d.name}</div>
          <div class="text-xs text-gray-500">Room: {d.room?.room_name || 'Unassigned'} · Patient: {d.patient?.name || 'Unassigned'}</div>
          <div class="text-xs text-gray-500">{d.description}</div>
        </div>
        <div class="ml-4">
          {#if editingId === d.ID}
            <div class="flex flex-col gap-2 w-64">
              <input class="px-2 py-1 border rounded" bind:value={editDescription} placeholder="Description" />

                <SearchableSelect
                    placeholder="Select a room"
                    options={rooms.map((room) => ({
                        label: room.room_name,
                        value: room.ID.toString(),
                    }))}
                    clearable
                    size="md"
                    bind:value={editRoomId}
                />

              <SearchableSelect
                  placeholder="Select a patient"
                  options={patients.map((patient) => ({
                      label: patient.name,
                      value: patient.ID.toString(),
                  }))}
                  clearable
                  size="md"
                  bind:value={editPatientId}
              />
              <div class="flex gap-2">
                <button class="px-3 py-1 bg-green-600 text-white rounded" on:click={saveEdit}>Save</button>
                <button class="px-3 py-1 bg-gray-200 rounded" on:click={()=> editingId = null}>Cancel</button>
              </div>
            </div>
          {:else}
            <div class="flex gap-2">
              <button class="px-3 py-1 bg-blue-600 text-white rounded" on:click={() => startEdit(d)}>Edit</button>
            </div>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>

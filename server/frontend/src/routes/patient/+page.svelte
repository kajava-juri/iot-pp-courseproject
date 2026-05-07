<script lang="ts">
    import { Label, Input, Helper, Button } from "flowbite-svelte";
    import { tryParseRooms, type RoomArray } from "$types/healthcare-db-types";
    import { fetcher } from "../../utils/fetcher";
    import { onMount } from "svelte";
    import { z } from "zod";
    import SearchableSelect from "$components/SearchableSelect.svelte";

    const patientFormSchema = z.object({
        name: z.string().trim().min(1, "Patient name is required"),
        patientID: z.string().trim().min(1, "PatientID is required"),
        roomId: z.string().trim().optional(),
    });

    let rooms = $state<RoomArray>([]);
    let selectedRoomId = $state<string>("");
    let patientName = $state<string>("");
    let patientID = $state<string>("");
    let formError = $state<string>("");
    let formErrorMap = $state<Record<string, string>>({});

    $effect(() => {
        console.log("Selected room ID:", selectedRoomId);
    });

    onMount(async () => {
        try {
            const response = (await fetcher("/room", "get")).data;
            const parsedResponse = tryParseRooms(response);
            if (!parsedResponse.success) {
                console.error(
                    "Failed to parse rooms data:",
                    parsedResponse.error,
                );
                return;
            }
            rooms = parsedResponse.data;
        } catch (error) {
            console.error("Error fetching rooms:", error);
        }
    });

    async function handleSubmit() {
        const result = patientFormSchema.safeParse({
            name: patientName,
            patientID,
            roomId: selectedRoomId,
        });

        let tempErrorMap: Record<string, string> = {};

        if (!result.success) {
            // formError = result.error.issues[0]?.message ?? "Invalid patient data";
            result.error.issues.forEach((issue) => {
                if (issue.path.length > 0) {
                    console.log(`Validation error for ${String(issue.path[0])}: ${issue.message}`);
                }
                tempErrorMap[issue.path[0] as string] = issue.message || "";
            });
            formErrorMap = tempErrorMap;
            return;
        }

        console.log("Creating patient with data:", result.data);
        try {
            const patientResult = await fetcher("/patient", "post", {
                name: patientName,
                patient_id: patientID,
                room_id: selectedRoomId ? parseInt(selectedRoomId) : undefined,
            });
        } catch (error: any) {
            console.error("Error creating patient:", error);
            formError = `Failed to create patient: ${error?.data || error.message}`;
            return;
        }
        // console.log("Patient creation result:", patientResult);
        // if (patientResult.status === 201) {
        //     // Reset form on success
        //     patientName = "";
        //     patientID = "";
        //     selectedRoomId = "";
        //     formErrorMap = {};
        // } else {
        //     formError = `Failed to create patient: ${patientResult.data}`;
        // }
    }
</script>

<div class="m-auto justify-center w-md">
    <div class="mb-6">
        <Label for="name">Patient name</Label>
        <Input id="name" placeholder="Patient name" bind:value={patientName} />
        {#if formErrorMap["name"]}
            <Helper class="mt-2" color="red">
                <span class="font-medium">{formErrorMap["name"]}</span>
            </Helper>
        {/if}
    </div>
    <div class="mb-6">
        <Label for="patient-id">PatientID</Label>
        <Input id="patient-id" placeholder="PatientID" bind:value={patientID} />
        {#if formErrorMap["patientID"]}
            <Helper class="mt-2" color="red">
                <span class="font-medium">{formErrorMap["patientID"]}</span>
            </Helper>
        {/if}
    </div>
    <div class="mb-6">
        <Label for="success" color="green" class="mb-2 block">Room</Label>
        <SearchableSelect
            placeholder="Select a room"
            options={rooms.map((room) => ({
                label: room.room_name,
                value: room.ID.toString(),
            }))}
            clearable
            size="md"
            bind:value={selectedRoomId}
        />
    </div>
    {#if formError}
        <Helper class="mt-2 my-4" color="red">
            <span class="font-medium">{formError}</span>
        </Helper>
    {/if}
    <Button type="submit" onclick={handleSubmit}>
        Create Patient
    </Button>
</div>

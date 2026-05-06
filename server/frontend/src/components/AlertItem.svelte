<script lang="ts">
    interface Props {
        PatientName: string;
        AlertType: string;
        RoomID: string;
        Timestamp: string;
    }

    let {
        PatientName,
        AlertType,
        RoomID,
        Timestamp
    }: Props = $props();
</script>

<div class="p-6 justify-center text-center flex flex-col bg-primary-700 rounded-2xl text-white font-semibold">
    <div>
        <span class="text-sm">{PatientName}</span>
        , <span class="text-sm">ROOM {RoomID}</span>
        <br/>
        <span>{AlertType}</span>
    </div>
    <div>
        <span class="text-xs text-gray-100">{Timestamp}</span>
    </div>
</div>
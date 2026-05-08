<script lang="ts">
    import { fetcher } from "../utils/fetcher";
    import type { Alert } from "$types/healthcare-db-types";
    import type { AxiosResponse } from "axios";

    interface Props {
        alert: Alert;
    }

    let { alert }: Props = $props();

    function isDeclined(a: Alert): boolean {
        return (a as Alert & { declined?: boolean }).declined === true;
    }

    function statusOf(
        a: Alert,
    ): "resolved" | "declined" | "acknowledged" | "unacknowledged" {
        if (a.resolved) return "resolved";
        if (isDeclined(a)) return "declined";
        if (a.acknowledged) return "acknowledged";
        return "unacknowledged";
    }

    let status = $derived(statusOf(alert));
    let patientLabel = $derived(
        alert.patient_id ? `Patient ${alert.patient?.name}` : "Unknown",
    );
    let roomLabel = $derived(
        alert.event?.room?.room_name ? alert.event.room.room_name.toString() : "N/A",
    );
    let alertType = $derived(alert.event?.type || "Unknown Event");
    let alertTypeText = $derived.by(() => {
        switch (alertType) {
            case "fall":
                return "Fall Detected";
            default:
                return alertType;
        }
    });
    let timestampLabel = $derived(
        new Date(alert.CreatedAt).toLocaleTimeString(),
    );

    let containerClass = $derived(
        status === "unacknowledged"
            ? "bg-primary-500"
            : status === "acknowledged"
              ? "bg-amber-700"
              : status === "declined"
                ? "bg-slate-600"
                : "bg-emerald-700",
    );

    let badgeText = $derived(
        status === "unacknowledged"
            ? "Unacknowledged"
            : status === "acknowledged"
              ? "Acknowledged"
              : status === "declined"
                ? "Declined"
                : "Resolved",
    );

    function handleAlertResponse(response: AxiosResponse) {
        if(response.status === 200) {
            console.log("Alert response successful:", response.data);
        } else {
            console.error("Alert response failed with status:", response.status);
        }
    }

    async function acknowledgeAlert() {
        const res = await fetcher(`/alert/${alert.ID}/acknowledge`, "post");
        handleAlertResponse(res);
    }

    async function declineAlert() {
        const res = await fetcher(`/alert/${alert.ID}/decline`, "post");
        handleAlertResponse(res);
    }

    async function resolveAlert() {
        const res = await fetcher(`/alert/${alert.ID}/resolve`, "post");
        handleAlertResponse(res);
    }
</script>

<div
    class={`justify-center text-center flex flex-col rounded-2xl text-white font-semibold w-xs ${containerClass}`}
>
    <div class="p-2 relative top-0 left-0 w-fit">
        <span class="text-xs rounded-full bg-black/20 uppercase p-1">
            {badgeText}
        </span>
    </div>
    <div class="px-3 pb-3">
        <div>
            <span class="text-sm">{patientLabel}</span>
            , <span class="text-sm">ROOM {roomLabel}</span>
            <br />
            <span class="uppercase">{alertTypeText}</span>
        </div>
        <div class="mt-1">
            <span class="text-xs text-gray-100">{timestampLabel}</span>
        </div>
        <div class="w-full flex justify-center gap-3 mt-3">
            {#if status === "unacknowledged"}
                <button class="px-3 py-1 bg-secondary-600 hover:bg-secondary-700 rounded-lg text-xs" onclick={acknowledgeAlert}
                    >On my way</button
                >
                <button class="px-3 py-1 bg-slate-600 hover:bg-slate-700 rounded-lg text-xs" onclick={declineAlert}
                    >Can't respond</button
                >
            {:else if status === "acknowledged"}
                <button class="px-3 py-1 bg-emerald-700 hover:bg-emerald-800 rounded-lg text-xs" onclick={resolveAlert}
                    >Resolve</button
                >
                <button class="px-3 py-1 bg-slate-700 hover:bg-slate-800 rounded-lg text-xs"
                    >View</button
                >
            {:else if status === "declined"}
                <button class="px-3 py-1 bg-slate-700 hover:bg-slate-800 rounded-lg text-xs"
                    >View</button
                >
            {:else}
                <button class="px-3 py-1 bg-slate-700 hover:bg-slate-800 rounded-lg text-xs"
                    >View</button
                >
            {/if}
        </div>
    </div>
</div>

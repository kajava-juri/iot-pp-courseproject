<script lang="ts">
  import "../app.css";
  import { onMount } from "svelte";
  import { Alert } from "flowbite-svelte";
  import { InfoCircleSolid } from "flowbite-svelte-icons";
  import Header from "$components/Header.svelte";
  import Sidemenu from "$components/Sidemenu.svelte";
  import { alertsStore } from "$stores/alerts";
  import { toasts, removeToast, addToast } from "$stores/toasts";
  import Alerts from "$components/Alerts.svelte";
  import { WsMessageStore } from "$stores/ws_store";
  import { fetcher } from "../utils/fetcher";
  import {
    tryParseAlert,
    tryParseAlerts,
    tryParseEvent,
    type Alert as DbAlert,
  } from "$types/healthcare-db-types";
  import AlertItem from "$components/AlertItem.svelte";
    import { parse } from "svelte/compiler";

  type ToastAlert = {
    toastId: string;
    alertId: number;
  };

  let { children } = $props();
  let newAlerts = $state<ToastAlert[]>([]);

  let alertsLoaded = $state<DbAlert[]>([]);

  function isDeclined(a: DbAlert): boolean {
    return (a as DbAlert & { declined?: boolean }).declined === true;
  }

  let unacknowledgedAlerts = $derived(
    alertsLoaded.filter(
      (alert) => !alert.resolved && !alert.acknowledged && !isDeclined(alert),
    ),
  );
  let acknowledgedAlerts = $derived(
    alertsLoaded.filter((alert) => !alert.resolved && alert.acknowledged),
  );
  let declinedAlerts = $derived(
    alertsLoaded.filter((alert) => !alert.resolved && isDeclined(alert)),
  );
  let resolvedAlerts = $derived(alertsLoaded.filter((alert) => alert.resolved));

  async function wsSetup() {
    WsMessageStore.subscribe((currentMessage) => {
      if (currentMessage) {
        console.log("Received WebSocket message:", currentMessage);

        // event/motion event toast
        if (currentMessage.topic.includes("/event/motion")) {

          let eventDataJson = currentMessage.payload;
          try {
            let eventData = tryParseEvent(eventDataJson);
            if (!eventData.success) {
              console.error(
                "Failed to parse motion event data from WebSocket message:",
                eventData.error,
              );
              return;
            }
            console.log("Parsed motion event data:", eventData);
            addToast(`Motion detected in ${eventData.data.room?.room_name || 'Unknown Room'}!`, "info", 4000);
          } catch (error) {
            console.error("Failed to parse motion event data:", error);
          }

          // let toastId: string = crypto.randomUUID();
          // newAlerts = [
          //   ...newAlerts,
          //   { toastId, alertId: -1 }, // alertId is not applicable for motion events
          // ];
          // setTimeout(() => {
          //   newAlerts = newAlerts.filter((item) => item.toastId !== toastId);
          // }, 4000);

        }

        if (currentMessage.topic.includes("/alert")) {
          let alertDataJson = currentMessage.payload;
          try {
            let alertData = tryParseAlert(alertDataJson);
            if (!alertData.success) {
              console.error(
                "Failed to parse alert data from WebSocket message:",
                alertData.error,
              );
              return;
            }
            console.log("Parsed alert data:", alertData.data);
            alertsStore.update((alerts) => [alertData.data, ...alerts]);
          } catch (error) {
            console.error("Failed to parse WebSocket message payload:", error);
          }
          // const alertData = {
          //   PatientName: `Patient ${currentMessage.patient_id || 'Unknown'}`,
          //   AlertType: 'Fall Detected',
          //   RoomID: 'N/A',
          //   Timestamp: new Date(currentMessage.timestamp).toLocaleTimeString()
          // };
          // alerts = [alertData, ...alerts];
          // alertsStore.update(alerts => [currentMessage, ...alerts]);
        }
      }
    });

    // Fetch unresolved alerts
    try {
      const response = (await fetcher("/alert/all/unresolved")).data;
      console.log("Fetched alerts:", response);
      const parseResult = tryParseAlerts(response);
      if (parseResult.success) {
        alertsStore.set(parseResult.data);
      } else {
        console.error("Failed to parse alerts:", parseResult.error);
      }
    } catch (error) {
      console.error("Failed to fetch alerts:", error);
    }
  }

  onMount(() => {
    wsSetup();
    let hasSeededInitialState = false;
    let knownAlertIds = new Set<number>();
    // Map to track timeout handles for each toast by their toastId
    const timeoutHandles = new Map<string, ReturnType<typeof setTimeout>>();

    // whenever alerts changes, filter out known alerts and push them to the toast alerts queue
    // for every toast alert, set a timeout that deletes it
    const unsubscribe = alertsStore.subscribe((alerts) => {
      console.log("Alerts store updated:", alerts);
      const currentIds = new Set(alerts.map((alert) => alert.ID));
      // Keep unresolved first for better operator visibility.
      const sorted = [...alerts].sort((a, b) => {
        if (a.resolved !== b.resolved) return a.resolved ? 1 : -1;
        return (
          new Date(b.CreatedAt).getTime() - new Date(a.CreatedAt).getTime()
        );
      });
      alertsLoaded = sorted;
      if (currentIds.size === 0) {
        console.log("No alerts in store, skipping toast generation");
        return;
      }

      if (!hasSeededInitialState) {
        knownAlertIds = currentIds;
        hasSeededInitialState = true;
        return;
      }

      const incomingAlerts = alerts.filter(
        (alert) => !knownAlertIds.has(alert.ID),
      );
      knownAlertIds = currentIds;

      for (const alert of incomingAlerts) {
        const toastId: string = crypto.randomUUID();
        newAlerts = [...newAlerts, { toastId, alertId: alert.ID }];
        const timeoutId = setTimeout(
          (toastId: string) => {
            newAlerts = newAlerts.filter((item) => item.toastId !== toastId);
            timeoutHandles.delete(toastId);
          },
          4000,
          toastId,
        );
        timeoutHandles.set(toastId, timeoutId);
      }
    });

    return () => {
      unsubscribe();
      for (const handle of timeoutHandles.values()) {
        clearTimeout(handle);
      }
      timeoutHandles.clear();
    };
  });

  function dismissToast(toastId: string) {
    newAlerts = newAlerts.filter((item) => item.toastId !== toastId);
  }
</script>

<div class="h-screen flex flex-col bg-white text-gray-900">
  <Header />

  <div class="flex flex-1 min-h-0">
    <Sidemenu />
    <div class="flex flex-1 w-full min-h-0 overflow-y-scroll">
      {@render children()}
    </div>
    <div class="flex flex-col w-md bg-gray-50 overflow-y-scroll min-h-0">
      <Alerts>
        {#if alertsLoaded.length === 0}
          <div class="text-gray-500 text-sm">No active alerts</div>
        {:else}
          {#if unacknowledgedAlerts.length > 0}
            <div
              class="w-full text-left text-xs font-semibold uppercase tracking-wide text-primary-500"
            >
              Unacknowledged
            </div>
            {#each unacknowledgedAlerts as alert (alert.ID)}
              <AlertItem {alert} />
            {/each}
          {/if}

          {#if acknowledgedAlerts.length > 0}
            <div
              class="w-full text-left text-xs font-semibold uppercase tracking-wide text-primary-600 mt-2"
            >
              Acknowledged
            </div>
            {#each acknowledgedAlerts as alert (alert.ID)}
              <AlertItem {alert} />
            {/each}
          {/if}

          {#if declinedAlerts.length > 0}
            <div
              class="w-full text-left text-xs font-semibold uppercase tracking-wide text-slate-600 mt-2"
            >
              Declined
            </div>
            {#each declinedAlerts as alert (alert.ID)}
              <AlertItem {alert} />
            {/each}
          {/if}

          {#if resolvedAlerts.length > 0}
            <div
              class="w-full text-left text-xs font-semibold uppercase tracking-wide text-emerald-700 mt-2"
            >
              Resolved
            </div>
            {#each resolvedAlerts as alert (alert.ID)}
              <AlertItem {alert} />
            {/each}
          {/if}
        {/if}
      </Alerts>
    </div>
  </div>

  <div
    class="pointer-events-none fixed right-4 top-20 z-50 flex w-80 flex-col gap-2"
  >
    {#each newAlerts as item (item.toastId)}
      <div class="pointer-events-auto">
        <Alert dismissable onclose={() => dismissToast(item.toastId)}>
          {#snippet icon()}<InfoCircleSolid class="h-5 w-5" />{/snippet}
          New alert added! (ID: {item.alertId})
        </Alert>
      </div>
    {/each}
    {#each $toasts as t (t.id)}
      <div class="pointer-events-auto">
        <Alert dismissable onclose={() => removeToast(t.id)}>
          {#snippet icon()}<InfoCircleSolid class="h-5 w-5" />{/snippet}
          {t.message}
        </Alert>
      </div>
    {/each}
  </div>
</div>

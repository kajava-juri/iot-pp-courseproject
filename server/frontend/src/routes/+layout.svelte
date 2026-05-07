<script lang="ts">
    import "../app.css";
    import { onMount } from 'svelte';
    import { Alert } from 'flowbite-svelte';
    import { InfoCircleSolid } from 'flowbite-svelte-icons';
    import Header from "$components/Header.svelte";
    import Sidemenu from "$components/Sidemenu.svelte";
    import { alertsStore } from '$stores/alerts';
    import Alerts from "$components/Alerts.svelte";
  import { WsSocketStore, WsMessageStore } from '$stores/ws_store.svelte.ts';
  import { fetcher } from '../utils/fetcher.ts';
    import { tryParseAlert, tryParseAlerts } from "../types/healthcare-db-types";
    import AlertItem from "$components/AlertItem.svelte";

    type ToastAlert = {
        toastId: string;
        alertId: number;
    };

    let { children } = $props();
    let newAlerts = $state<ToastAlert[]>([]);

    let alerts = $state<Array<{
        PatientName: string;
        AlertType: string;
        RoomID: string;
        Timestamp: string;
    }>>([]);

/**  onMount(async () => {
    // Subscribe to WebSocket messages
    WsMessageStore.subscribe(currentMessage => {
      if (currentMessage) {
        console.log('Received WebSocket message:', currentMessage);

        if (currentMessage.topic.endsWith('/alert/fall')) {
          let alertDataJson = currentMessage.payload;
          try {
            let alertData = tryParseAlert(alertDataJson);
            if (!alertData.success) {
              console.error('Failed to parse alert data from WebSocket message');
              return;
            }
            console.log('Parsed alert data:', alertData.data);
            const alertInfo = {
              PatientName: alertData.data.patient_id ? `Patient ${alertData.data.patient_id}` : 'Unknown',
              AlertType: alertData.data.event?.type || 'Unknown Event',
              RoomID: alertData.data.event?.device_id.toString() || 'N/A',
              Timestamp: new Date(alertData.data.CreatedAt).toLocaleTimeString()
            };
            alerts = [alertInfo, ...alerts];
            alertsStore.update(alerts => [alertData.data, ...alerts]);
          } catch (error) {
            console.error('Failed to parse WebSocket message payload:', error);
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
      const response = await fetcher('/alert/all/unresolved');
      console.log('Fetched alerts:', response);
      const parseResult = tryParseAlerts(response);
      if (parseResult.success) {
        alerts = parseResult.data.map(alert => ({
          PatientName: alert.patient_id ? `Patient ${alert.patient_id}` : 'Unknown',
          AlertType: alert.event?.type || 'Unknown Event',
          RoomID: alert.event?.device_id.toString() || 'N/A',
          Timestamp: new Date(alert.CreatedAt).toLocaleTimeString()
        }));
        alertsStore.set(parseResult.data);
      } else {
        console.error('Failed to parse alerts:', parseResult.error);
      }
    } catch (error) {
      console.error('Failed to fetch alerts:', error);
    }
  });*/

  async function wsSetup() {
    WsMessageStore.subscribe(currentMessage => {
      if (currentMessage) {
        console.log('Received WebSocket message:', currentMessage);

        if (currentMessage.topic.endsWith('/alert/fall')) {
          let alertDataJson = currentMessage.payload;
          try {
            let alertData = tryParseAlert(alertDataJson);
            if (!alertData.success) {
              console.error('Failed to parse alert data from WebSocket message');
              return;
            }
            console.log('Parsed alert data:', alertData.data);
            const alertInfo = {
              PatientName: alertData.data.patient_id ? `Patient ${alertData.data.patient_id}` : 'Unknown',
              AlertType: alertData.data.event?.type || 'Unknown Event',
              RoomID: alertData.data.event?.device_id.toString() || 'N/A',
              Timestamp: new Date(alertData.data.CreatedAt).toLocaleTimeString()
            };
            alerts = [alertInfo, ...alerts];
            alertsStore.update(alerts => [alertData.data, ...alerts]);
          } catch (error) {
            console.error('Failed to parse WebSocket message payload:', error);
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
      const response = await fetcher('/alert/all/unresolved');
      console.log('Fetched alerts:', response);
      const parseResult = tryParseAlerts(response);
      if (parseResult.success) {
        alerts = parseResult.data.map(alert => ({
          PatientName: alert.patient_id ? `Patient ${alert.patient_id}` : 'Unknown',
          AlertType: alert.event?.type || 'Unknown Event',
          RoomID: alert.event?.device_id.toString() || 'N/A',
          Timestamp: new Date(alert.CreatedAt).toLocaleTimeString()
        }));
        alertsStore.set(parseResult.data);
      } else {
        console.error('Failed to parse alerts:', parseResult.error);
      }
    } catch (error) {
      console.error('Failed to fetch alerts:', error);
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
            const currentIds = new Set(alerts.map((alert) => alert.ID));
            if (currentIds.size === 0) {
                console.log('No alerts in store, skipping toast generation');
                return;
            }

            if (!hasSeededInitialState) {
                knownAlertIds = currentIds;
                hasSeededInitialState = true;
                return;
            }

            const incomingAlerts = alerts.filter((alert) => !knownAlertIds.has(alert.ID));
            knownAlertIds = currentIds;

            for (const alert of incomingAlerts) {
                const toastId: string = crypto.randomUUID();
                newAlerts = [...newAlerts, { toastId, alertId: alert.ID }];
                const timeoutId = setTimeout((toastId: string) => {
                    newAlerts = newAlerts.filter((item) => item.toastId !== toastId);
                    timeoutHandles.delete(toastId);
                }, 4000, toastId);
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

    <div class="flex flex-1">
        <Sidemenu />
        {@render children()}
        <div class="flex flex-col w-md">
            <Alerts>
            {#each alerts as alert}
                <AlertItem
                PatientName={alert.PatientName}
                AlertType={alert.AlertType}
                RoomID={alert.RoomID}
                Timestamp={alert.Timestamp}
                />
            {/each}
            {#if alerts.length === 0}
                <div class="text-gray-500 text-sm">No active alerts</div>
            {/if}
            </Alerts>
        </div>
    </div>

        <div class="pointer-events-none fixed right-4 top-20 z-50 flex w-80 flex-col gap-2">
            {#each newAlerts as item (item.toastId)}
                <div class="pointer-events-auto">
                    <Alert dismissable onclose={() => dismissToast(item.toastId)}>
                        {#snippet icon()}<InfoCircleSolid class="h-5 w-5" />{/snippet}
                        New alert added! (ID: {item.alertId})
                    </Alert>
                </div>
            {/each}
        </div>
</div>
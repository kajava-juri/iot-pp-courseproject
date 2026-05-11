import {derived, writable} from "svelte/store";
import type { Alert } from "../types/healthcare-db-types";

export const alertsStore = writable<Alert[]>([]);

export const alertsCount = derived(alertsStore, $alerts => $alerts.filter(alert => !alert.resolved).length);
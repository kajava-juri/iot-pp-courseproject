import { writable } from 'svelte/store';
import { tryParseAlert, tryParseWSMessage, type WSMessage } from '../types/healthcare-db-types';
import { alertsStore } from './alerts';

const WsMessageStore = writable<WSMessage | null>(null);
const WsSocketStore = writable<WebSocket | null>(null);

let hostname = import.meta.env.VITE_SERVER_HOSTNAME || 'localhost';

const topicPrefix = import.meta.env.VITE_WEARABLE_IMU_TOPIC_PREFIX;

const socket = new WebSocket(`ws://${hostname}:${import.meta.env.VITE_WEB_PORT}/ws`);

socket.onopen = () => {
  console.log('WebSocket connection established');
  WsSocketStore.set(socket);

  socket.send(JSON.stringify({
    "action": "subscribe",
    "topics": [
      `${topicPrefix}/${import.meta.env.VITE_FALL_EVENT_TOPIC}`,
      `${topicPrefix}/${import.meta.env.VITE_ALERT_TOPIC}`,
      `${import.meta.env.VITE_ALERT_UPDATE_TOPIC}`,
      `${topicPrefix}/alert/vibration`
    ]
  }))
};

socket.onclose = () => {
  console.log('WebSocket connection closed');
  WsMessageStore.set(null);
  WsSocketStore.set(null);
};

socket.onmessage = (event) => {
  console.log('WebSocket message received:', event.data);
  const result = tryParseWSMessage(event.data);
  if (result.success) {
    WsMessageStore.set(result.data);
    if (result.data.topic === `${import.meta.env.VITE_ALERT_UPDATE_TOPIC}`) {
      console.log('Received alert update:', result.data.payload);
      const alertUpdate = tryParseAlert(result.data.payload);
      if (alertUpdate.success) {
        console.log('Parsed alert update:', alertUpdate.data);
        alertsStore.update(alerts => {
          const index = alerts.findIndex(alert => alert.ID === alertUpdate.data.ID);
          if (index !== -1) {
            const updatedAlerts = [...alerts];
            updatedAlerts[index] = alertUpdate.data;
            return updatedAlerts;
          }
          return alerts;
        });
      } else {
        console.error('Failed to parse alert update:', alertUpdate.error);
      }
    }
  } else {
    console.error('Failed to parse WebSocket message:', result.error);
  }
};

socket.onerror = (error) => {
  console.error('WebSocket error:', error);
};

export { WsSocketStore, WsMessageStore };
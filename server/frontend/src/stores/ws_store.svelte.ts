import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { tryParseWSMessage, type WSMessage } from '../types/healthcare-db-types';

const WsMessageStore = writable<WSMessage | null>(null);
const WsSocketStore = writable<WebSocket | null>(null);

let hostname = $state('localhost');
if (browser) {
  hostname = window.location.hostname;
}

const topicPrefix = import.meta.env.VITE_WEARABLE_IMU_TOPIC_PREFIX;

const socket = new WebSocket(`ws://${hostname}:${import.meta.env.VITE_WEB_PORT}/ws`);

socket.onopen = () => {
  console.log('WebSocket connection established');
  WsSocketStore.set(socket);

  socket.send(JSON.stringify({
    "action": "subscribe",
    "topics": [
      `${topicPrefix}/${import.meta.env.VITE_FALL_EVENT_TOPIC}`,
      `${topicPrefix}/${import.meta.env.VITE_ALERT_TOPIC}`

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
  } else {
    console.error('Failed to parse WebSocket message:', result.error);
  }
};

socket.onerror = (error) => {
  console.error('WebSocket error:', error);
};

export { WsSocketStore, WsMessageStore };
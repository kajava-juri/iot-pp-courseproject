import { writable } from 'svelte/store';
import { browser } from '$app/environment';

const wsStore = writable<WebSocket | null>(null);

let hostname = $state('localhost');
if (browser) {
  hostname = window.location.hostname;
}

const topicPrefix = import.meta.env.VITE_WEARABLE_IMU_TOPIC_PREFIX;

const socket = new WebSocket(`ws://${hostname}:${import.meta.env.VITE_WEB_PORT}/ws`);

socket.onopen = () => {
  console.log('WebSocket connection established');
  wsStore.set(socket);

  socket.send(JSON.stringify({
    "action": "subscribe",
    "topics": [`${topicPrefix}/${import.meta.env.VITE_FALL_EVENT_TOPIC}`]
  }))
};

socket.onclose = () => {
  console.log('WebSocket connection closed');
  wsStore.set(null);
};

socket.onmessage = (event) => {
  console.log('WebSocket message received:', event.data);
  wsStore.set(event.data);
};

socket.onerror = (error) => {
  console.error('WebSocket error:', error);
};

export { wsStore };
<script>
  import { onMount, onDestroy } from 'svelte';
  import mqtt from 'mqtt';

  // MQTT broker WebSocket URL (served inside Docker network via mosquitto)
  const brokerWsUrl = 'ws://localhost:9001';

  let status = 'Disconnected';
  let messages = [];
  let client;

  // Topics to subscribe to
  const topics = ['room/#', 'wearable/#'];

  onMount(() => {
    client = mqtt.connect(brokerWsUrl);

    client.on('connect', () => {
      status = 'Connected';
      topics.forEach(topic => client.subscribe(topic));
    });

    client.on('message', (topic, payload) => {
      messages = [
        { topic, payload: payload.toString(), time: new Date().toLocaleTimeString() },
        ...messages.slice(0, 49),  // keep last 50 messages
      ];
    });

    client.on('error', (err) => {
      status = `Error: ${err.message}`;
    });

    client.on('close', () => {
      status = 'Disconnected';
    });
  });

  onDestroy(() => {
    if (client) client.end();
  });
</script>

<main>
  <h1>IoT Dashboard</h1>

  <section class="status">
    <span>MQTT Status: <strong>{status}</strong></span>
  </section>

  <section class="messages">
    <h2>Recent Messages</h2>
    {#if messages.length === 0}
      <p class="empty">No messages received yet.</p>
    {:else}
      <ul>
        {#each messages as msg}
          <li>
            <code>{msg.topic}</code>: {msg.payload}
            <small>{msg.time}</small>
          </li>
        {/each}
      </ul>
    {/if}
  </section>
</main>

<style>
  main {
    font-family: sans-serif;
    max-width: 800px;
    margin: 2rem auto;
    padding: 0 1rem;
  }
  h1 {
    color: #2d6a4f;
  }
  .status {
    background: #f0f0f0;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    margin-bottom: 1rem;
  }
  .empty {
    color: #888;
    font-style: italic;
  }
  ul {
    list-style: none;
    padding: 0;
  }
  li {
    display: flex;
    gap: 1rem;
    align-items: baseline;
    border-bottom: 1px solid #eee;
    padding: 0.4rem 0;
  }
  small {
    color: #999;
    margin-left: auto;
  }
</style>

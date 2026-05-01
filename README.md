# IoT Course Project

A full-stack IoT system consisting of three Docker-based server services and three PlatformIO-based edge nodes.

## Project Structure

```
root/
├── edge_nodes/
│   ├── wearable_node/     # PlatformIO – wearable health/activity sensor
│   ├── room_node/         # PlatformIO – room environmental sensor
│   └── interface_node/    # PlatformIO – MQTT gateway / bridge node
├── server/
│   ├── svelte-ui/         # Svelte front-end (Node.js dev server)
│   ├── mosquitto/         # Eclipse Mosquitto MQTT broker
│   └── nodered/           # Node-RED flow editor
├── docs/
├── diagrams/
├── README.md
├── docker-compose.yml
└── .gitignore
```

## Server (Docker)

Three containers are orchestrated with `docker-compose.yml`:

| Service      | Port(s)        | Description                        |
|-------------|----------------|------------------------------------|
| `svelte-ui` | 5173           | Svelte front-end (Vite dev server) |
| `mosquitto` | 1883 / 9001    | MQTT broker (TCP + WebSocket)      |
| `nodered`   | 1880           | Node-RED flow editor               |

### Quick start

```bash
docker compose up --build
```

- Svelte UI  → http://localhost:5173
- Node-RED   → http://localhost:1880
- MQTT broker → `mqtt://localhost:1883`

## Edge Nodes (PlatformIO)

Each sub-folder under `edge_nodes/` is an independent PlatformIO project.

| Node             | Role                                            |
|------------------|-------------------------------------------------|
| `wearable_node`  | Wearable sensor – collects health/activity data |
| `room_node`      | Room sensor – temperature, humidity, CO₂        |
| `interface_node` | Gateway – bridges edge data to MQTT broker      |

### Build & flash

```bash
cd edge_nodes/wearable_node
pio run --target upload
```

## Requirements

- [Docker](https://docs.docker.com/get-docker/) ≥ 24
- [Docker Compose](https://docs.docker.com/compose/) (included with Docker Desktop)
- [PlatformIO Core](https://docs.platformio.org/en/latest/core/installation/) for edge nodes

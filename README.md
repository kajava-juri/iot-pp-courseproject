# IoT Course Project

A full-stack IoT system consisting of three Docker-based server services and three PlatformIO-based edge nodes.

## Project Structure

**A couple note**
- mosquitto probably will not be used because TalTech has its own MQTT broker
- nodered is there but most of the logic and message handling will be done using the Golang project

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

| Service      | Port(s)        | Description                       |
|-------------|----------------|------------------------------------|
| `svelte-ui` | 3000           | Svelte front-end (Vite dev server) |
| `mosquitto` | 1883 / 9001    | MQTT broker (TCP + WebSocket)      |
| `nodered`   | 1880           | Node-RED flow editor               |
| `Go API`    | 8081           | Serves HTTP API, Websockets and handles edge nodes MQTT messages |

### Quick start

Copy env variable and change to your environment
``` bash
cp .env.example .env
cp .env server/backend
```

**You might need to use `sudo` for docker commands**

```bash
docker compose up --build -d
```

for development and hot refresh:
```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build frontend
```

- --build   Build images before starting containers
- --d       Detached mode, dont start listening to the containers

See running service, check that everything is OK and not restarting.
``` bash
docker ps
```

- Svelte UI     -> http://localhost:3000
- Node-RED      -> http://localhost:1880
- MQTT broker   -> `mqtt://localhost:1883`
- Go API        -> http://localhost:8080/

###### PgAdmin

1. First open PgAdmin and log in


URL: http://localhost:5050

Username/Email: check in .env PGADMIN_DEFAULT_EMAIL 

Password: check in .env PGADMIN_DEFAULT_PASSWORD

2. Add server

**General**
- server name: choose whatever (e.g healthcare_server)

**Connection**
- Host name/address: check .env DB_HOST (default is `postgres`)
- Port: 5432
- Username: postgres
- Password: postgres

**Press Connect**


## Edge Nodes (PlatformIO)

Each sub-folder under `edge_nodes/` is an independent PlatformIO project.

| Node             | Role                                                   |
|------------------|--------------------------------------------------------|
| `wearable_node`  | Wearable sensor – collects health/activity data        |
| `room_node`      | Room sensor – temperature, humidity, CO₂               |
| `interface_node` | Gateway – interface with the user (OLED, relay?)       |

### Build & flash

Using PlatformIO VSCode extension.

## Requirements

- [Docker](https://docs.docker.com/get-docker/) ≥ 24
- [Docker Compose](https://docs.docker.com/compose/) (included with Docker Desktop)
- PlatformIO Core

## References

1. Svelte docker setup: https://medium.com/@balazs.csaba.diy/optimized-dockerfile-for-sveltekit-applications-from-experience-and-best-practices-99603d8d1303
2. Golang Chi guide: https://dev.to/luthfisauqi17/getting-started-with-golang-chi-a-guide-to-building-a-simple-api-210m

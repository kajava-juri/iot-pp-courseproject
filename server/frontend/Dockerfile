# credit to Balazs Csaba: https://medium.com/@balazs.csaba.diy/optimized-dockerfile-for-sveltekit-applications-from-experience-and-best-practices-99603d8d1303
# followed the recommendations in the article

ARG NODE_IMAGE=node:24.11.1-alpine3.22

# Stage 1 – build
FROM ${NODE_IMAGE} AS builder

# Build-time environment variables (from docker-compose args)
ARG VITE_WEB_PORT=8080
ARG VITE_SERVER_URL=http://localhost
ARG VITE_SERVER_HOSTNAME=localhost
ARG VITE_FALL_EVENT_TOPIC=event/fall
ARG VITE_ALERT_TOPIC=alert/fall
ARG VITE_WEARABLE_IMU_TOPIC_PREFIX=ESP16

ENV VITE_WEB_PORT=${VITE_WEB_PORT}
ENV VITE_SERVER_URL=${VITE_SERVER_URL}
ENV VITE_SERVER_HOSTNAME=${VITE_SERVER_HOSTNAME}
ENV VITE_FALL_EVENT_TOPIC=${VITE_FALL_EVENT_TOPIC}
ENV VITE_ALERT_TOPIC=${VITE_ALERT_TOPIC}
ENV VITE_WEARABLE_IMU_TOPIC_PREFIX=${VITE_WEARABLE_IMU_TOPIC_PREFIX}

WORKDIR /app

COPY package*.json ./
RUN npm i

COPY . .
RUN npm run build

# Stage 2 – runtime
FROM ${NODE_IMAGE} AS runner

# Upgrade base Alpine Linux packages and install curl (optional, useful for runtime diagnostics)
# might be used for debugging and testing backend API endpoints
RUN apk upgrade -U && apk add curl

WORKDIR /app

COPY --chown=node:node package.json package-lock.json .

# Install production dependencies only
RUN npm i --omit=dev && \
    # Clean up npm cache to reduce image size
    npm cache clean --force && \
    # Remove unnecessary files and folders from node_modules such as docs, tests, maps, git metadata
    find node_modules \( \
        -type d -empty \
        -o -iname "license*" \
        -o -name "*.md" \
        -o -name "*.txt" \
        -o -name "*.map" \
        -o -name ".git*" \
        -o -name "*.yml" \
        -o -name "*.yaml" \
        -o -name "*.json" -path "*/test/*" \
        -o -name "*.json" -path "*/tests/*" \
        -o -name "test" -type d \
        -o -name "tests" -type d \
        -o -name "__tests__" -type d \
        -o -name "coverage" -type d \
        -o -name ".nyc_output" -type d \
    \) -delete && \
    # Remove leftover tmp and cache files
    rm -rf /tmp/* /var/cache/apk/* /root/.npm && \
    # Remove globally installed npm to save space and remove vulnerability
    npm r -g npm

COPY --from=builder --chown=node:node /app/build ./build
COPY --from=builder --chown=node:node /app/package*.json ./

# Make extra sure build sourcemaps are deleted
RUN find build -name "*.map" -delete

RUN mkdir -p /data && chown -R node:node /data

# Switch to non-root user
USER node

# Set environment variable
ENV NODE_ENV=production

EXPOSE 3000

VOLUME ["/data"]

CMD ["node", "build"]

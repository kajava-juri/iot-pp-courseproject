# credit to Balazs Csaba: https://medium.com/@balazs.csaba.diy/optimized-dockerfile-for-sveltekit-applications-from-experience-and-best-practices-99603d8d1303
# followed the recommendations in the article

ARG NODE_IMAGE=node:24.11.1-alpine3.22

# Stage 1 – build
FROM ${NODE_IMAGE} AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

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
RUN npm ci --omit=dev && \
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
ENV NODE_ENV=${NODE_ENV}

EXPOSE 3000

VOLUME ["/data"]

CMD ["node", "build"]

# Optimised
FROM node:20-slim AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

# Runner
FROM node:20-alpine AS runner
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --omit=dev && npm prune --production && rm -rf /tmp/* ~/.npm
COPY --from=builder /app/build ./build
EXPOSE 3000
CMD ["node", "build"]


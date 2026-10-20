# ---- Stage 1: Build ----
FROM node:20-alpine AS builder
WORKDIR /app

# Install dependencies only when needed
COPY package*.json ./
RUN npm ci

# Copy the rest of the application code
COPY . .

# Build the Next.js application (ensure output: 'output: "trace"' or 'standalone' is set in next.config.js)
RUN npm run build


# ---- Stage 2: Runner ----
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
# Uncomment if you want to disable telemetry during runtime
# ENV NEXT_TELEMETRY_DISABLED=1

# Create a non-root system user and group
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy standalone output and static assets with correct permissions
COPY --chown=nextjs:nodejs --from=builder /app/public ./public
COPY --chown=nextjs:nodejs --from=builder /app/.next/standalone ./
COPY --chown=nextjs:nodejs --from=builder /app/.next/static ./.next/static

# Switch to the non-root user
USER nextjs

EXPOSE 3000
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Start the standalone Next.js server
CMD ["node", "server.js"]
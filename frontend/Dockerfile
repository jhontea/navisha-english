# syntax=docker/dockerfile:1
# Stage 1: Dependencies
FROM node:22-alpine AS deps

WORKDIR /app

COPY package.json package-lock.json* ./

# Cache node_modules di VPS — npm install skip jika tidak ada perubahan
RUN --mount=type=cache,target=/root/.npm \
    npm install --prefer-offline

# Stage 2: Build
FROM node:22-alpine AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

# Cache Next.js build cache — recompile hanya file yang berubah
RUN --mount=type=cache,target=/app/.next/cache \
    npm run build

# Stage 3: Runtime
FROM node:22-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3020

ENV PORT=3020
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]

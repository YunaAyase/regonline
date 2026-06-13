FROM node:22-alpine AS builder
WORKDIR /build
COPY package*.json ./
RUN npm ci
COPY . .
RUN npx nuxt build

FROM node:22-alpine
WORKDIR /app
COPY --from=builder /build/.output ./.output
COPY --from=builder /build/package*.json ./
RUN npm ci --omit=dev
ENV HOST=0.0.0.0 PORT=5000
EXPOSE 5000
CMD ["node", ".output/server/index.mjs"]

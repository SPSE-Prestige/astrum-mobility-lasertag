# ---------- base ----------
FROM node:20 AS base
WORKDIR /usr/src/app
COPY package*.json ./

# ---------- dev ----------
FROM base AS dev
ENV NODE_ENV=development
RUN npm install --legacy-peer-deps
COPY . .
CMD ["npm", "run", "dev"]

# ---------- staging ----------
FROM base AS staging
ENV NODE_ENV=staging
RUN npm ci --legacy-peer-deps
COPY . .
RUN npm run build

# ---------- build ----------
FROM base AS build
ARG NEXT_PUBLIC_STRAPI_URL=http://localhost:3100
ENV NODE_ENV=production
ENV NEXT_PUBLIC_STRAPI_URL=${NEXT_PUBLIC_STRAPI_URL}
RUN npm ci --legacy-peer-deps
COPY . .
RUN npm run build

# ---------- production ----------
FROM node:20 AS production
WORKDIR /usr/src/app

ENV NODE_ENV=production

COPY --from=build /usr/src/app/.next ./.next
COPY --from=build /usr/src/app/public ./public
COPY --from=build /usr/src/app/package*.json ./
COPY --from=build /usr/src/app/next.config.mjs ./

RUN npm ci --legacy-peer-deps

EXPOSE 3000
CMD ["sh", "-c", "npx next start -p ${PORT:-3000}"]
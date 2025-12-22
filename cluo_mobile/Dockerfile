FROM node:20.11.1-alpine AS base
WORKDIR /usr/src/app

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

COPY package.json pnpm-lock.yaml ./
COPY .npmrc* ./

FROM base AS prod-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod --frozen-lockfile

FROM base AS build
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY . .
RUN pnpm run build

FROM node:20.11.1-alpine AS runtime
WORKDIR /usr/src/app

COPY --from=prod-deps node_modules/ ./node_modules/
COPY --from=build build/ ./build/
COPY --from=build package.json ./

EXPOSE 3000

CMD ["node", "build/index.js"]


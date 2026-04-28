FROM node:20.11.1-alpine AS base
WORKDIR /usr/src/app

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

COPY package.json pnpm-lock.yaml .npmrc* ./

FROM base AS build
RUN --mount=type=cache,id=pnpm-mobile,target=/pnpm/store pnpm install --frozen-lockfile

COPY . .
RUN pnpm run build

FROM node:20.11.1-alpine AS runtime
WORKDIR /usr/src/app

COPY --from=build /usr/src/app/build/ ./build/
COPY --from=build /usr/src/app/package.json ./

EXPOSE 3000

CMD ["node", "build/index.js"]

FROM alpine:3
ARG PLUGIN_MODULE=github.com/traefik/hpauth
ARG PLUGIN_GIT_REPO=https://github.com/RashedEmon/hpauth.git
ARG PLUGIN_GIT_BRANCH=1.0.0
RUN apk update && \
    apk add git && \
    git clone ${PLUGIN_GIT_REPO} /plugins-local/src/${PLUGIN_MODULE} \
      --depth 1 --single-branch --branch ${PLUGIN_GIT_BRANCH}
FROM traefik:v3.0
COPY --from=0 /plugins-local /plugins-local
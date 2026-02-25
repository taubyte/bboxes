FROM node:lts-alpine3.22
RUN apk update && apk add bash

# Remove when as-sdk is on npm
RUN apk add --no-cache git

# So `asc` is on PATH even when project npm i fails (e.g. network timeout to Bitbucket)
RUN npm install -g assemblyscript

COPY ./utils /utils

RUN mkdir /src /out

VOLUME /src
VOLUME /out

FROM node:current-alpine3.16
RUN apk update && apk add bash

# Remove when as-sdk is on npm
RUN apk add --no-cache git

COPY ./utils /utils

RUN mkdir /src /out

VOLUME /src
VOLUME /out

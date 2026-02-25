FROM rust:alpine3.23
RUN apk update && apk add bash

RUN apk add clang

RUN apk add lld

RUN apk add musl-dev

RUN apk add binaryen

RUN apk add compiler-rt

RUN rustup target add wasm32-wasip1

COPY ./utils /utils

RUN mkdir /src /out

VOLUME /src 
VOLUME /out 
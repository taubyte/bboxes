FROM rust:alpine3.16
RUN apk update && apk add bash

RUN apk add clang

RUN apk add lld

RUN apk add musl-dev

RUN apk add binaryen

RUN apk add compiler-rt compiler-rt-static

RUN rustup target add wasm32-wasi

RUN cargo install cargo-wasi

COPY ./utils /utils

RUN mkdir /src /out

VOLUME /src 
VOLUME /out 
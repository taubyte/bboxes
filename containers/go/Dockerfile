FROM debian:stable-slim

RUN apt-get clean && apt-get update --fix-missing && \
    apt-get install -y wget gcc gcc-avr avr-libc

ENV GO_RELEASE=1.19
RUN wget https://dl.google.com/go/go${GO_RELEASE}.linux-amd64.tar.gz && \
    tar xfv go${GO_RELEASE}.linux-amd64.tar.gz -C /usr/local && \
    rm go${GO_RELEASE}.linux-amd64.tar.gz && \
    find /usr/local/go -mindepth 1 -maxdepth 1 ! -name 'src' ! -name 'VERSION' ! -name 'bin' ! -name 'pkg' -exec rm -rf {} +
ENV PATH=${PATH}:/usr/local/go/bin

ENV TINYGO_RELEASE=0.26.0
RUN wget https://github.com/tinygo-org/tinygo/releases/download/v${TINYGO_RELEASE}/tinygo${TINYGO_RELEASE}.linux-amd64.tar.gz && \
    tar xfv tinygo${TINYGO_RELEASE}.linux-amd64.tar.gz -C /usr/local && \
    rm -rf /usr/local/tinygo/src/device/{sam,stm32,nxp,nrf,avr,esp,rp} && \
    rm -rf /usr/local/tinygo/lib/{nrfx,mingw-w64,macos-minimal-sdk} && \
    rm -rf /usr/local/tinygo/src/examples && \
    rm tinygo${TINYGO_RELEASE}.linux-amd64.tar.gz
ENV PATH=${PATH}:/usr/local/tinygo/bin

RUN apt-get remove -y wget && \
    apt-get autoremove -y && \
    apt-get clean

RUN apt-get install -y bash

RUN apt-get install -y wabt

RUN apt-get install -y git git-lfs less ssh

RUN git lfs install

COPY ./go /go

COPY ./utils /utils

ENV PATH=${PATH}:/tools:/usr/bin:/bin

RUN echo $PATH

ENV GOPATH=/go

RUN mkdir /src /out

VOLUME /src
VOLUME /out
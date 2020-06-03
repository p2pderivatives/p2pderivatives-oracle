ARG ARG_CFD_GO_VERSION=0.1.21
ARG CFD_GO_ZIP=cfdgo-v${ARG_CFD_GO_VERSION}-alpine_x86_64.zip
FROM golang:1.14-alpine as dev
RUN apk update
RUN apk add make cmake gcc g++ libc-dev git unzip
ENV GO111MODULE=on

WORKDIR /p2pderivatives-oracle

# install cfd-go dependencies
ENV LD_LIBRARY_PATH=/usr/local/lib64
ARG ARG_CFD_GO_VERSION
ARG CFD_GO_ZIP
ENV CFD_GO_VERSION=${ARG_CFD_GO_VERSION}
RUN wget -O /${CFD_GO_ZIP} https://github.com/cryptogarageinc/cfd-go/releases/download/v${CFD_GO_VERSION}/${CFD_GO_ZIP} \
     && unzip -q /${CFD_GO_ZIP} -d /

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY tools tools/
COPY Makefile .

RUN make install-tools

COPY . .

FROM dev as build
RUN make oracle

FROM alpine as prod

WORKDIR /p2pdoracle

# runtime dependencies + cfd
RUN apk update
RUN apk add libstdc++
##look for shared librairies here
ENV LD_LIBRARY_PATH=/usr/local/lib64
ARG ARG_CFD_GO_VERSION
ARG CFD_GO_ZIP
ENV CFD_GO_VERSION=${ARG_CFD_GO_VERSION}
COPY --from=dev /${CFD_GO_ZIP} .
RUN unzip -q ${CFD_GO_ZIP} -d / \
    && rm /p2pdoracle/${CFD_GO_ZIP}

RUN mkdir -p /config
COPY ./test/config/default.release.yml /config/default.yml
VOLUME [ "/config" ]
RUN mkdir -p /key
VOLUME ["/key"]

COPY --from=build /p2pderivatives-oracle/bin/oracle /p2pdoracle/

ENTRYPOINT [ "/p2pdoracle/oracle" ]
CMD [ "-config", "/config", "-appname", "p2pdoracle", "-e", "default", "-migrate" ]

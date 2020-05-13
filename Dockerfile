ARG ARG_CFD_GO_VERSION=0.1.21

FROM golang:1.14-alpine as dev
RUN apk update
RUN apk add make cmake gcc g++ libc-dev git wget unzip
ENV GO111MODULE=on

WORKDIR /p2pderivatives-oracle

# install cfd-go dependencies
ENV LD_LIBRARY_PATH=/usr/local/lib64
ARG ARG_CFD_GO_VERSION
ENV CFD_GO_VERSION=${ARG_CFD_GO_VERSION}
ENV CFD_GO_ZIP=cfdgo-v${CFD_GO_VERSION}-alpine_x86_64.zip
RUN wget -O /cfd-go-v${CFD_GO_VERSION}.zip https://github.com/cryptogarageinc/cfd-go/releases/download/v${CFD_GO_VERSION}/${CFD_GO_ZIP}
RUN unzip -q /cfd-go-v${CFD_GO_VERSION}.zip -d /

RUN go get -u gotest.tools/gotestsum

COPY go.mod .
COPY go.sum .

RUN go mod download

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
ENV CFD_GO_VERSION=${ARG_CFD_GO_VERSION}
COPY --from=dev /cfd-go-v${CFD_GO_VERSION}.zip .
RUN unzip -q cfd-go-v${CFD_GO_VERSION}.zip -d /
RUN rm cfd-go-v${CFD_GO_VERSION}.zip

RUN mkdir -p /config
VOLUME [ "/config" ]
RUN mkdir -p /key
VOLUME ["/key"]

COPY --from=build /p2pderivatives-oracle/bin/oracle /p2pdoracle/

ENTRYPOINT [ "/p2pdoracle/oracle" ]
CMD [ "-config", "/config", "-appname", "p2pdoracle", "-e", "integration", "-migrate" ]
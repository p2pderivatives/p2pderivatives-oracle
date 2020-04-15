FROM golang:1.14-alpine as dev
RUN apk update
RUN apk add make gcc libc-dev git wget unzip
ENV GO111MODULE=on

WORKDIR /app

RUN go get -u github.com/jstemmer/go-junit-report

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . /app

FROM dev as build
RUN make oracle

FROM alpine as prod

RUN mkdir -p /config
VOLUME [ "/config" ]

COPY --from=build /app/bin/oracle /p2pdoracle/
WORKDIR /p2pdoracle

ENTRYPOINT [ "/p2pdoracle/oracle" ]
CMD [ "-config", "/config", "-appname", "p2pdoracle", "-e", "integration", "-migrate" ]
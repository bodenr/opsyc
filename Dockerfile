FROM golang:1.14.3-alpine3.11 AS build

WORKDIR /go/src/opsyc
COPY . .

RUN go get && \
    go build


FROM alpine:3.11.6

COPY --from=build /go/src/opsyc/opsyc /opsyc
COPY ./assets/ /assets/
ENV OPSYC_ASSETS_DIR=/assets
RUN apk add --no-cache dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--", "/opsyc"]
EXPOSE 8080:8080

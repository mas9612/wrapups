FROM golang:1.12.1 as build

RUN mkdir -p /go/src/github.com/mas9612/wrapups
WORKDIR /go/src/github.com/mas9612/wrapups
COPY . .
RUN make test && make build-server


FROM alpine:3.9.2

LABEL maintainer="Masato Yamazaki <mas9612@gmail.com>"

RUN mkdir /app
WORKDIR /app

COPY --from=build /go/src/github.com/mas9612/wrapups/wuserver .
ENTRYPOINT ["/app/wuserver"]

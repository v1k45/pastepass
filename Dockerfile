FROM golang:1.22-alpine3.19 AS builder
LABEL org.opencontainers.image.source="https://github.com/v1k45/pastepass"

# go mod cache layer
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
ARG VERSION=0.00-dev
RUN CGO_ENABLED=0 go build -o /out/pastepass -ldflags "-X 'github.com/v1k45/pastepass/config.Version=$VERSION'"


FROM alpine:3.19

ENV USER=pastepass
ENV UID=1000
ENV GID=1000
ENV APP_DIR=/app
ENV DATA_DIR=/data

RUN apk add --no-cache tini ca-certificates \
    && addgroup -g $GID -S $USER \
    && adduser -u $UID -S -G $USER $USER

WORKDIR $APP_DIR
COPY --from=builder /out/pastepass .

RUN mkdir -p $DATA_DIR && chown $UID:$GID $DATA_DIR

USER $USER
EXPOSE 8008

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["./pastepass", "-db-path", "/data/pastepass.db"]

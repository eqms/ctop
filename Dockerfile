FROM golang:1.26-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS builder

RUN apk add --no-cache make git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

# Runs as root by design: ctop needs read access to the mounted
# /var/run/docker.sock. To drop privileges, run with
# `--user <uid> --group-add <docker-gid>` matching the host socket.
FROM scratch
ENV TERM=linux
COPY --from=builder /app/ctop /ctop
ENTRYPOINT ["/ctop"]

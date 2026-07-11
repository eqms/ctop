FROM golang:1.26-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS builder

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

FROM golang:1.27-alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS builder

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

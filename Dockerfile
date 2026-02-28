FROM golang:1.22-alpine AS builder

RUN apk add --no-cache make git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM scratch
ENV TERM=linux
COPY --from=builder /app/ctop /ctop
ENTRYPOINT ["/ctop"]

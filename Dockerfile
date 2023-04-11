FROM golang:1.19 as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o karbour cmd/main.go

FROM alpine:3.17.3

WORKDIR /
COPY --from=builder /workspace/karbour .

ENTRYPOINT ["/karbour"]

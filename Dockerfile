FROM alpine:3.17.3

WORKDIR /
COPY karbour .

ENTRYPOINT ["/karbour"]

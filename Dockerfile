FROM golang:1.15-alpine as builder
WORKDIR /go/src/app
COPY . .
RUN go build ./cmd/multi-fwd

FROM alpine:3
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app/multi-fwd .
ENTRYPOINT ["/root/multi-fwd"]
CMD ["server", "-interface=0.0.0.0"]

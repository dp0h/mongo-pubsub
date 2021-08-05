FROM golang:latest as builder

WORKDIR /build/app

COPY . /build/

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o mps .

FROM scratch

COPY --from=builder /build/app/mps /app/mps

ENTRYPOINT ["/app/mps"]

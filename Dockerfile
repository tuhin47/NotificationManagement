ARG GO_VERSION=1.24.2
FROM golang:${GO_VERSION}-alpine AS builder
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
RUN apk add --no-cache ca-certificates
WORKDIR /src
COPY ./ ./
RUN go mod vendor
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor GOOS=linux go build -a -o /app .

FROM alpine:latest AS final
COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app
USER nobody:nobody

EXPOSE 8080
ENTRYPOINT ["/app"]
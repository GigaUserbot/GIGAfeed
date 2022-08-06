# Build Stage: Build bot using the alpine image, also install doppler in it
FROM golang:1.19-alpine AS builder
RUN addgroup -S gigafeeduser \
    && adduser -S -u 10000 -g gigafeeduser gigafeeduser
RUN apk add -U --no-cache ca-certificates
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o out/GIGAFeed

# Run Stage: Run bot using the bot and doppler binary copied from build stage
FROM efwfgwegfywef
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/out/GIGAFeed /
COPY --from=builder /etc/passwd /etc/passwd
USER gigafeeduser
CMD ["/GIGAFeed"]

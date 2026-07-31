# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/magicstream .

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /out/magicstream .
RUN mkdir -p media && cp -r /app/media 2>/dev/null || true
COPY --from=builder /app/media/sample.mp4 ./media/

EXPOSE 8080
CMD ["./magicstream"]

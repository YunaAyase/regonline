FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:latest
RUN apk add --no-cache tesseract-ocr tesseract-ocr-chi-sim
WORKDIR /app
COPY --from=builder /app/server ./
COPY config.yaml ./
RUN mkdir -p /app/data /app/photos
EXPOSE 8080
CMD ["/app/server"]

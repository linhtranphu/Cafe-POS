# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o print-bridge main.go

# Runtime stage
FROM alpine:latest

# Install Chromium and dependencies
RUN apk --no-cache add \
    chromium \
    chromium-chromedriver \
    font-noto \
    font-noto-cjk \
    ttf-dejavu \
    fontconfig \
    ca-certificates

# Update font cache
RUN fc-cache -f

# Set Chrome binary path for chromedp
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/bin/chromium-browser

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/print-bridge .

# Expose port
EXPOSE 3001

# Run the application
CMD ["./print-bridge"]

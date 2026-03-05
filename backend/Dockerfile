# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cafe-pos-server .

# Runtime stage
FROM alpine:latest

# Install ONLY essential system dependencies
# CHROMIUM REMOVED - HTML rendering moved to local print bridge
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    lsof \
    fontconfig \
    ttf-dejavu \
    font-noto \
    font-noto-cjk \
    font-liberation \
    font-roboto

# Update font cache
RUN fc-cache -f

# Set timezone to Vietnam
ENV TZ=Asia/Ho_Chi_Minh

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/cafe-pos-server .

# Copy embedded templates
COPY --from=builder /app/application/services/templates ./application/services/templates

# Expose port
EXPOSE 3000

# Run the application
CMD ["./cafe-pos-server"]

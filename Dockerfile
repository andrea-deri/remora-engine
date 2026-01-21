# ---------------
# Stage 1: Build
# ---------------
FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Add CA certs
RUN apk add --no-cache ca-certificates

# Copy project module and Go files
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the static executable
ARG VERSION
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X 'main.VERSION=${VERSION}' -s -w" -o remora ./cmd/monolith

# ---------------
# Stage 2: Run
# ---------------
FROM scratch

# Copy the built static executable from previous stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/remora /remora
COPY properties.env .

# Set the app start command
CMD ["/remora"]

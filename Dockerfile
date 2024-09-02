# Build Stage
####################
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum krampus.go ./

RUN go mod download
RUN go build -ldflags="-s -w -X main.Version=0.2.0" krampus.go


# Final Stage
####################
FROM gcr.io/distroless/static

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/krampus /usr/local/bin/krampus

# Command to run the executable
ENTRYPOINT [ "/usr/local/bin/krampus" ]

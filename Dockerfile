# Build Stage
####################
FROM golang:1.25-alpine AS builder

ARG VERSION=dev

WORKDIR /app

COPY go.mod go.sum krampus.go ./

RUN go mod download
RUN go build -ldflags="-s -w -X main.Version=${VERSION}" krampus.go


# Final Stage
####################
FROM gcr.io/distroless/static

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/krampus /usr/local/bin/krampus

# Command to run the executable
ENTRYPOINT [ "/usr/local/bin/krampus" ]

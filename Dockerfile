# Build Stage
####################
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

ARG VERSION=dev

WORKDIR /app

COPY go.mod go.sum krampus.go ./

RUN go mod download
<<<<<<< HEAD
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w -X main.Version=$VERSION" krampus.go
=======
RUN go build -ldflags="-s -w -X main.Version=${VERSION}" krampus.go
>>>>>>> 8ec64eab587d5aa8525eef509639b260b527444e


# Final Stage
####################
FROM gcr.io/distroless/static

# Copy the binary from the builder stage
COPY --from=builder /app/krampus /usr/local/bin/krampus

# Command to run the executable
ENTRYPOINT [ "/usr/local/bin/krampus" ]

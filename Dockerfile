
# Create the intermediate builder image.
FROM golang:1.24.0 as builder

WORKDIR /build

# Prepare dependencies in different layer to have a cache
ADD go.mod go.sum ./

RUN go mod download

# Build the application
ADD . ./

# Build the static application binary.
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o anniedad ./cmd/*.go

# Create the final small image.
FROM alpine:latest

LABEL description="This is the API gateway server"

RUN apk update && apk upgrade && apk add --no-cache ca-certificates && \
    rm -rf /var/cache/apk/* && \
    addgroup --gid 39999 anniedad && \
    adduser -h /app -s /bin/sh -G anniedad -u 39999 -D anniedad

WORKDIR /app/

# Copy from builder exactly what we need
COPY --from=builder /build/anniedad .

USER anniedad
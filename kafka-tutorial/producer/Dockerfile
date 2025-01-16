FROM ubuntu:22.04

RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && apt-get install -y curl

# Build stage

FROM golang:1.20-alpine AS builder

# set working dir

WORKDIR /app

# COPY go module files

COPY go.mod go.sum ./

RUN go mod download

COPY  . ./ 

# run application

RUN go build -o /producer *.go

# Final stage

FROM alpine:3.18

# copy compiled binary from builder stage

COPY --from=builder /producer /app/producer

CMD ["/app/producer"]
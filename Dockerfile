FROM golang:1.16 AS builder

RUN apt update \
    && apt install gettext git ca-certificates -y \
    && update-ca-certificates \
    && curl -s https://storage.googleapis.com/shipa-client/install-cloud-cli.sh | bash

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -ldflags '-extldflags "-static"' -tags 'osusergo netgo static_build' -o action .

RUN cp /build/action /usr/local/bin

#### Build result container
#FROM scratch
#
## Copy static executable
#COPY --from=builder /build/action .
#COPY --from=builder /usr/local/bin/shipa .
## Copy system files
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=builder /etc/passwd /etc/passwd


USER nobody

# Command to run when starting the container
ENTRYPOINT ["action"]
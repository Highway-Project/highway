# Build the manager binary
FROM golang:1.15 as builder

WORKDIR /highway
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

# Build
RUN make build

FROM alpine:3.12
WORKDIR /
COPY --from=builder /highway/bin/highway .

ENTRYPOINT ["/highway"]

